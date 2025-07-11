package telemetry

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	config "terraform-provider-hitachi/hitachi/common/config"
	diskcache "terraform-provider-hitachi/hitachi/common/diskcache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	vosbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	"time"

	"github.com/google/uuid"
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
func UpdateTelemetryStats(status string, elapsedTime float64, storageSettingInt interface{}, outputForModelOrVersion interface{}) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	statsMutex.Lock()
	defer statsMutex.Unlock()

	task := getMethodTaskExecutionInfo(status, elapsedTime, storageSettingInt, outputForModelOrVersion)

	err := updateLocalUsagesStat(task)
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
func getMethodTaskExecutionInfo(status string, elapsedTime float64, storageSettingInt interface{}, outputForModelOrVersion interface{}) MethodTaskExecution {
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
		storageModel = getSanStorageModel(storageSettingInt, outputForModelOrVersion)
	} else if connectionType == "vosb" {
		storageType = "sds_block"
		storageSetting := storageSettingInt.(vosbmodel.StorageDeviceSettings)
		terraformResourceMethod = storageSetting.TerraformResourceMethod
		// storageSerial = getUUIDFromIP(storageSetting.ClusterAddress).String()
		// storageSerial = storageSetting.ClusterAddress
		storageModel = getVosbStorageVersion(storageSettingInt, outputForModelOrVersion)
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

// updateLocalUsagesStat updates the local usages.json file
func updateLocalUsagesStat(task MethodTaskExecution) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	usages := &UsagesTelemetry{
		ExecutionStats: make(map[string]ExecutionStat),
	}

	if err := loadUsagesStatFromLocalFile(usages); err != nil {
		log.WriteError("Error loading usages telemetry stats from file: %v", err)
	}

	updateExecutionStats(usages, task)
	updateSanStorageSystems(usages, task)
	updateSdsBlockSystems(usages, task)

	if err := saveUsagesToLocalFile(usages); err != nil {
		log.WriteError("Error saving usages telemetry stats to file: %v", err)
		return fmt.Errorf("failed to save usages telemetry stats: %w", err)
	}

	return nil
}

// updateExecutionStats updates the execution stats in the usages.json
func updateExecutionStats(usages *UsagesTelemetry, task MethodTaskExecution) {
	moduleName := task.TerraformTfType + "." + task.GatewayMethodName
	stat := usages.ExecutionStats[moduleName]

	if task.Status == "success" {
		stat.Success++
	} else if task.Status == "failure" {
		stat.Failure++
	}

	total := stat.Success + stat.Failure
	if total > 0 {
		stat.AverageTime = math.Round(((stat.AverageTime*float64(total-1)+task.ElapsedTime)/float64(total))*100) / 100
	}

	usages.ExecutionStats[moduleName] = stat
}

// updateSanStorageSystems updates the list of SAN storage systems in the usages.json
func updateSanStorageSystems(usages *UsagesTelemetry, task MethodTaskExecution) {
	if task.ConnectionType != "san" || task.StorageSerial == "" {
		return
	}

	for _, system := range usages.SanStorageSystems {
		if system.StorageSerial == task.StorageSerial {
			return
		}
	}

	usages.SanStorageSystems = append(usages.SanStorageSystems, SanStorageSystem{
		StorageModel:  task.StorageModel,
		StorageSerial: task.StorageSerial,
	})
}

// updateSdsBlockSystems updates the list of SDS block systems in the usages.json
func updateSdsBlockSystems(usages *UsagesTelemetry, task MethodTaskExecution) {
	if task.ConnectionType != "vosb" || task.StorageModel == "" {
		return
	}

	for _, system := range usages.SdsBlockSystems {
		if system.ClusterAddress == task.StorageSerial {
			return
		}
	}

	usages.SdsBlockSystems = append(usages.SdsBlockSystems, SdsBlockSystem{
		ClusterAddress: task.StorageSerial,
		Version:        task.StorageModel,
	})
}

// loadUsagesStatFromLocalFile reads the usages telemetry stats from a JSON file.
func loadUsagesStatFromLocalFile(usages *UsagesTelemetry) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	file, err := os.Open(TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // It's okay if the file doesn't exist
		}
		return fmt.Errorf("failed to open usages telemetry file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(usages)
	if err != nil {
		return fmt.Errorf("failed to decode usages telemetry JSON: %w", err)
	}
	return nil
}

// saveUsagesToLocalFile saves the UsagesTelemetry struct to a JSON file.
func saveUsagesToLocalFile(usages *UsagesTelemetry) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	file, err := os.OpenFile(TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open usages telemetry file: %w", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(usages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal usages telemetry data: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write usages telemetry data: %w", err)
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

	if config.ConfigData == nil || config.ConfigData.AWS_URL == "" {
		log.WriteDebug("AWS URL is not configured. Skipping telemetry data sending.")
		return nil
	}

	log.WriteDebug("Sending telemetry data: %+v", telemetryPayload)

	// err = sendPOSTRequestToAWS(AWS_URL, telemetryPayload)
	err = sendPOSTRequestToAWS(config.ConfigData.AWS_URL, telemetryPayload)
	if err != nil {
		log.WriteError("Error sending telemetry stat details: %v", err)
		return err
	}

	log.WriteDebug("Telemetry stat details sent successfully to AWS: %v %v", telemetryPayload.ModuleName, telemetryPayload.OperationName)
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
		"User-Agent":   "terraform",
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

	awsTimeout := time.Duration(config.DEFAULT_AWS_TIMEOUT) * time.Second
	if config.ConfigData != nil && config.ConfigData.AWSTimeout > 0 {
		awsTimeout = time.Duration(config.ConfigData.AWSTimeout) * time.Second
	}

	// HTTP client with timeout, no proxy, and cert verification skipped (like validate_certs=False)
	tr := &http.Transport{
		Proxy: nil, // disables proxy
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING: disable in production!
		},
	}
	client := &http.Client{
		Timeout:   awsTimeout,
		Transport: tr,
	}

	log.WriteDebug("Sending POST request to AWS")
	log.WriteDebug("URL: %s", url)

	prettyBody := &bytes.Buffer{}
	if err := json.Indent(prettyBody, reqBodyInBytes, "", "  "); err != nil {
		log.WriteDebug("Raw request body: %s", string(reqBodyInBytes))
	} else {
		log.WriteDebug("Request Body:\n%s", prettyBody.String())
	}

	log.WriteDebug("Request Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			log.WriteDebug("  %s: %s", key, value)
		}
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.WriteError("Error sending HTTP request to telemetry AWS:\n  URL: %s\n  Error: %+v", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.WriteError("Received non-OK response from telemetry AWS: %d %s - %s", resp.StatusCode, resp.Status, string(bodyBytes))
		return fmt.Errorf("received non-OK response from telemetry AWS: %d %s", resp.StatusCode, resp.Status)
	}

	// log.WriteDebug("Request sent successfully to telemetry AWS.")

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

func getSanStorageModel(storageSettingInt interface{}, outputForModelOrVersion interface{}) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting := storageSettingInt.(sanmodel.StorageDeviceSettings)

	// Use global cache if already populated
	if sanStorageModel != "" {
		log.WriteDebug("from global cache var sanStorageModel: %v", sanStorageModel)
		return sanStorageModel
	}

	key := storageSetting.MgmtIP + ":StorageSystemInfo"
	var cachedInfo sanmodel.StorageSystemInfo
	found, _ := diskcache.Get(key, &cachedInfo)
	if found {
		sanStorageModel = cachedInfo.Model
		log.WriteDebug("from disk cache sanStorageModel: %v", sanStorageModel)
		return sanStorageModel
	}

	// Fallback to API response if provided
	if infoFromAPI, ok := outputForModelOrVersion.(*sanmodel.StorageSystemInfo); ok && infoFromAPI != nil {
		log.WriteDebug("sanmodel.StorageSystemInfo from API: %+v", infoFromAPI)
		sanStorageModel = infoFromAPI.Model
		return sanStorageModel
	}

	return ""
}

func getVosbStorageVersion(storageSettingInt interface{}, outputForModelOrVersion interface{}) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting := storageSettingInt.(vosbmodel.StorageDeviceSettings)

	// Use global cache if already populated
	if vosbStorageVersion != "" {
		log.WriteDebug("from global cache var vosbStorageVersion: %v", vosbStorageVersion)
		return vosbStorageVersion
	}

	// read from disk cache first
	key := storageSetting.ClusterAddress + ":StorageVersionInfo"
	var cachedInfo vosbmodel.StorageVersionInfo
	found, _ := diskcache.Get(key, &cachedInfo)
	if found {
		vosbStorageVersion = cachedInfo.ApiVersion
		log.WriteDebug("from disk cache vosbStorageVersion: %v", vosbStorageVersion)
		return vosbStorageVersion
	}

	// for first time http call of vosb GetStorageVersionInfo(), get apiVersion from output of api call
	if versionInfo, ok := outputForModelOrVersion.(*vosbmodel.StorageVersionInfo); ok {
		log.WriteDebug("vosbmodel.StorageVersionInfo: %+v", versionInfo)
		vosbStorageVersion = versionInfo.ApiVersion
		return vosbStorageVersion
	}

	return ""
}

func IsUserConsentExist() bool {
	if _, err := os.Stat(TERRAFORM_USER_CONSENT_FILE); os.IsNotExist(err) {
		return false
	}
	return true
}
