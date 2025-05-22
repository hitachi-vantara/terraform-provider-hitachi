package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	diskcache "terraform-provider-hitachi/hitachi/common/diskcache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	vosbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	"time"
)

var statsMutex sync.Mutex

// terraform restarts the plugin binary everytime, so these globals in memory are only good for one run
var global_Consent *UserConsent
var sanStorageModel string
var vosbStorageVersion string

func init() {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Ensure the directory exists, create if it doesn't
	if _, err := os.Stat(TERRAFORM_TELEMETRY_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(TERRAFORM_TELEMETRY_DIR, os.ModePerm)
		if err != nil {
			log.WriteError("Failed to create directory: %v\n", err)
			os.Exit(1) // Exit with failure if the directory creation fails
		}
	}
}

// UpdateTelemetryStats updates the execution stats based on the method call
func UpdateTelemetryStats(status string, elapsedTime float64, storageSettingInt interface{}, thirdArg interface{}) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	statsMutex.Lock()
	defer statsMutex.Unlock()

	task := getMethodTaskExecutionInfo(status, elapsedTime, storageSettingInt, thirdArg)

	err := updateLocalJsonFileStats(task)
	if err != nil {
		log.WriteError("Error updating JSON file stats: %v", err)
	}

	// Send execution details to the API after saving the stats
	err = sendTelemetryStatsToAWS(task)
	if err != nil {
		log.WriteError("Error sending telemetry stats: %v", err)
	}
}

// fill up initial data
func getMethodTaskExecutionInfo(status string, elapsedTime float64, storageSettingInt interface{}, thirdArg interface{}) MethodTaskExecution {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	siteId := global_Consent.SiteID

	// Capture call stack
	connectionType, methodName := GetGatewayCallStackInfo()

	terraformResourceMethod := ""
	storageModel := ""
	storageSerial := ""
	storageType := ""
	if connectionType == "san" {
		storageType = "vsp"
		storageSetting := storageSettingInt.(sanmodel.StorageDeviceSettings)
		terraformResourceMethod = storageSetting.TerraformResourceMethod
		storageSerial = strconv.Itoa(storageSetting.Serial)
		storageModel = getSanStorageModel(storageSettingInt, thirdArg)
	} else if connectionType == "vosb" {
		storageType = "sds_block"
		storageSetting := storageSettingInt.(vosbmodel.StorageDeviceSettings)
		terraformResourceMethod = storageSetting.TerraformResourceMethod
		storageSerial = getUUIDFromIP(storageSetting.ClusterAddress).String()
		storageModel = getVosbStorageVersion(storageSettingInt, thirdArg)
	}

	log.WriteDebug("terraformResourceMethod: %+v", terraformResourceMethod)
	terraformResourceName := connectionType + "." + terraformResourceMethod
	terraformTfType := getTerraformProviderMapValue(terraformResourceName)

	task := MethodTaskExecution{
		TerraformTfType:       terraformTfType,
		TerraformResourceName: terraformResourceName,
		GatewayMethodName:     methodName,
		Status:                status,
		ElapsedTime:           math.Round(elapsedTime*100) / 100,
		ConnectionType:        connectionType,
		StorageModel:          storageModel,
		StorageSerial:         storageSerial,
		StorageType:           storageType,
		SiteId:                siteId,
	}

	log.WriteDebug("Method Task Execution: %+v", task)
	return task
}

// updateLocalJsonFileStats updates the stats for a specific method in the JSON file
func updateLocalJsonFileStats(task MethodTaskExecution) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Initialize json stat
	executionJsonStat := make(map[string]*ExecutionJsonFileStats)

	// First, load the existing json stats from the JSON file
	err := loadStatsFromLocalJSONFile(&executionJsonStat)
	if err != nil {
		log.WriteError("Error loading json stats from file: %v", err)
	}

	// Retrieve json stats for the method or initialize if not present
	moduleName := task.TerraformTfType + "." + task.GatewayMethodName
	stats := executionJsonStat[moduleName]
	if stats == nil {
		stats = &ExecutionJsonFileStats{}
	}

	if task.Status == "success" {
		stats.Success++
	} else if task.Status == "failure" {
		stats.Failure++
	}

	averageTime := (stats.AverageTime*float64(stats.Success+stats.Failure-1) + task.ElapsedTime) / float64(stats.Success+stats.Failure)
	stats.AverageTime = math.Round(averageTime*100) / 100

	executionJsonStat[moduleName] = stats

	err = saveToLocalJSONFile(executionJsonStat)
	if err != nil {
		log.WriteError("Error saving execution stats to file: %v", err)
		return fmt.Errorf("failed to save execution stats: %w", err)
	}
	return nil
}

// loadStatsFromLocalJSONFile reads the execution stats
func loadStatsFromLocalJSONFile(executionJsonStat *map[string]*ExecutionJsonFileStats) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	file, err := os.Open(TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // If the file doesn't exist, it's okay to continue
		}
		return fmt.Errorf("failed to open stats file: %w", err)
	}
	defer file.Close()

	// Decode the JSON data into the executionJsonStat map
	decoder := json.NewDecoder(file)
	err = decoder.Decode(executionJsonStat)
	if err != nil {
		return fmt.Errorf("failed to decode stats from file: %w", err)
	}
	return nil
}

// saveToLocalJSONFile saves the execution json stats to a JSON file
func saveToLocalJSONFile(executionJsonStat map[string]*ExecutionJsonFileStats) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	file, err := os.OpenFile(TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Marshal the TaskStats data into JSON format
	statsJSON, err := json.MarshalIndent(executionJsonStat, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	_, err = file.Write(statsJSON)
	if err != nil {
		return fmt.Errorf("failed to write stats to file: %w", err)
	}

	return nil
}

// sendTelemetryStatsToAWS sends telemetry stats to an external API
func sendTelemetryStatsToAWS(task MethodTaskExecution) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	operationStatus := 1 // success
	if task.Status != "success" {
		operationStatus = 0
	}

	telemetryPayload := TelemetryPayload{
		ModuleName:      task.TerraformTfType,
		OperationName:   task.GatewayMethodName,
		OperationStatus: operationStatus,
		ConnectionType:  task.ConnectionType,
		ProcessTime:     task.ElapsedTime,
		StorageModel:    task.StorageModel,
		StorageSerial:   task.StorageSerial,
		StorageType:     task.StorageType,
		Site:            task.SiteId,
	}

	// Save the telemetry data to the JSON file
	err := saveTelemetryStatsToJSONFile(telemetryPayload)
	if err != nil {
		log.WriteError("Error saving execution stats to file: %v", err)
		return fmt.Errorf("failed to save execution stats: %w", err)
	}

	log.WriteDebug("Sending telemetry data: %+v", telemetryPayload)

	err = sendPOSTRequestToAWS(APIG_URL, telemetryPayload)
	if err != nil {
		log.WriteError("Error sending telemetry stat details: %v", err)
		return err
	}

	log.WriteDebug("Telemetry stat details sent successfully to API.")
	return nil
}

// sendPOSTRequestToAWS sends a POST request to a specified URL
func sendPOSTRequestToAWS(url string, data interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Convert the request body data into JSON format
	reqBodyInBytes, err := json.Marshal(data)
	if err != nil {
		log.WriteError("Error marshalling data to JSON: %v", err)
		return err
	}

	// Define the headers for the request
	headers := map[string]string{
		"Content-Type": "application/json",
		"user-agent":   "terraform",
	}

	// Create a new HTTP request with the specified URL, method, and body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyInBytes))
	if err != nil {
		log.WriteError("Error creating new HTTP request: %v", err)
		return err
	}

	// Set the headers for the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Optional: Add Basic Authentication headers if required
	httpBasicAuth := utils.HttpBasicAuthentication{}
	httpBasicAuth.SetAuthHeaders(req)

	// Create an HTTP client with a timeout (e.g., 30 seconds)
	client := &http.Client{
		Timeout: 30 * time.Second, // Set the timeout to 30 seconds
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.WriteError("Error sending HTTP request to telemetry AWS: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WriteError("Received non-OK response from telemetry AWS: %d %s", resp.StatusCode, resp.Status)
		return fmt.Errorf("received non-OK response from telemetry AWS: %d %s", resp.StatusCode, resp.Status)
	}

	log.WriteDebug("Request sent successfully to telemetry AWS.")

	return nil
}

// saveTelemetryStatsToJSONFile saves the telemetry payload to a JSON file, for debugging
func saveTelemetryStatsToJSONFile(telemetryPayload TelemetryPayload) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	file, err := os.OpenFile(TERRAFORM_TELEMETRY_STATS_AWS_FILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Marshal the telemetryPayload into JSON format
	statsJSON, err := json.MarshalIndent(telemetryPayload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry to JSON: %w", err)
	}

	// add a newline before appending the new data for clarity
	statsJSON = append(statsJSON, '\n')

	log.WriteDebug("Telemetry statsJSON saved: %+v", string(statsJSON))

	_, err = file.Write(statsJSON)
	if err != nil {
		return fmt.Errorf("failed to write telemetry to file: %w", err)
	}

	return nil
}

func readUserConsentFile() UserConsent {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var consent UserConsent

	data, err := os.ReadFile(TERRAFORM_USER_CONSENT_FILE)
	if err != nil {
		log.WriteDebug("Error reading consent file: %v. Treating it as no consent", err)
		return consent
	}

	err = json.Unmarshal(data, &consent)
	if err != nil {
		log.WriteDebug("Error unmarshalling JSON: %v. Treating it as no consent", err)
		return consent
	}

	return consent
}

// CheckTelemetryConsent checks user consent
func CheckTelemetryConsent() bool {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if global_Consent == nil {
		consent := readUserConsentFile()
		global_Consent = &consent
	}
	return global_Consent.UserConsentAccepted
}

// for vosb, get the UUID based on an IP address
func getUUIDFromIP(ip string) uuid.UUID {
	// Use a predefined namespace for UUIDv5 or create your own (e.g., DNS namespace)
	namespace := uuid.NameSpaceDNS
	// Generate UUIDv5 based on the IP address using SHA-1 internally
	return uuid.NewSHA1(namespace, []byte(ip))
}

// getTerraformProviderMapValue
// input ex:   vosb.resource.resourceVosbChangeUserPasswordCreate
// returns ex: vosb.resource.hitachi_vosb_change_user_password
func getTerraformProviderMapValue(input string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("input: %+v", input)
	if input == "" {
		return ""
	}

	suffixes := []string{"Create", "Update", "Read", "Delete", "CustomDiff"}

	// Split into three parts
	parts := strings.SplitN(input, ".", 3)
	log.WriteDebug("parts: %+v", parts)
	if len(parts) < 3 {
		return ""
	}

	// Remove the suffix from the third part if it matches any of the defined suffixes
	for _, suffix := range suffixes {
		if strings.HasSuffix(parts[2], suffix) {
			parts[2] = parts[2][:len(parts[2])-len(suffix)]
			break // Only remove the first matching suffix
		}
	}

	resourceName := parts[0] + "." + parts[1] + "." + strings.Title(parts[2])
	log.WriteDebug("resourceName: %+v", resourceName)

	if val, exists := TerraformProviderMap[resourceName]; exists {
		return val
	}

	return ""
}

func getSanStorageModel(storageSettingInt interface{}, thirdArg interface{}) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting := storageSettingInt.(sanmodel.StorageDeviceSettings)

	// use global
	if sanStorageModel != "" {
		return sanStorageModel
	}

	// read from disk cache
	key := storageSetting.MgmtIP + ":StorageSystemInfo"
	storageInfo := sanmodel.StorageSystemInfo{}
	found, _ := diskcache.Get(key, &storageInfo)
	if found {
		sanStorageModel = storageInfo.Model
		return sanStorageModel
	}

	// for first time http call of san GetStorageSystemInfo(), get storage model from output of api call
	if storageInfo, ok := thirdArg.(*sanmodel.StorageSystemInfo); ok {
		log.WriteDebug("sanmodel.StorageSystemInfo: %+v", storageInfo)
		sanStorageModel = storageInfo.Model
		return sanStorageModel
	}

	return ""
}

func getVosbStorageVersion(storageSettingInt interface{}, thirdArg interface{}) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting := storageSettingInt.(vosbmodel.StorageDeviceSettings)

	// use global
	if vosbStorageVersion != "" {
		return vosbStorageVersion
	}

	// read from disk cache first
	key := storageSetting.ClusterAddress + ":StorageVersionInfo"
	versionInfo := vosbmodel.StorageVersionInfo{}
	found, err := diskcache.Get(key, &versionInfo)
	log.WriteDebug("diskcache.Get Found=%v, ERR: %+v, versionInfo:%+v", found, err, versionInfo)

	if found {
		vosbStorageVersion = versionInfo.ApiVersion
		return vosbStorageVersion
	}

	// for first time http call of vosb GetStorageVersionInfo(), get apiVersion from output of api call
	if versionInfo, ok := thirdArg.(*vosbmodel.StorageVersionInfo); ok {
		log.WriteDebug("StorageVersionInfo: %+v", versionInfo)
		vosbStorageVersion = versionInfo.ApiVersion
		return vosbStorageVersion
	}

	return ""
}

// // not used yet
// func ValidateTerraformUserConsent() string {
// 	if _, err := os.Stat(TERRAFORM_USER_CONSENT_FILE); os.IsNotExist(err) {
// 		return USER_CONSENT_MISSING
// 	}
// 	return ""
// }
