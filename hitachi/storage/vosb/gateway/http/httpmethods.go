package vssbstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

func getCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := ""
	if apiSuf == "configuration/version" {
		url = GetUrlWithoutVersion(storageSetting.ClusterAddress, apiSuf)
	} else {
		url = GetUrl(storageSetting.ClusterAddress, apiSuf)
	}

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	resJSONString, err := utils.HTTPGet(url, nil, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponse(resJSONString, err)
		log.WriteError(err)
		return err
	}

	err2 := json.Unmarshal([]byte(resJSONString), output)
	if err2 != nil {
		log.WriteError(err)
		return fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return nil
}


func PostCallAsyncForm(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	var pHeader *map[string]string
	var header map[string]string
	if storageSetting.ContentType != "" {
		header = make(map[string]string)
		header["Content-Type"] = storageSetting.ContentType
		pHeader = &header
	}
	
	// var form *bytes.Buffer
	form, ok := reqBody.(*bytes.Buffer)
	if ok {
		// If reqBody is a bytes.Buffer, we can use it directly
		fmt.Printf("ok= %v\n", ok)
		reqBodyInBytes = form.Bytes()
	}

	jobString, err := utils.HTTPPostForm(url, pHeader, reqBodyInBytes, *form, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponseExt(jobString, err)
		log.WriteError(err)
		return nil, err
	}
			fmt.Printf("jobString= %v\n", jobString)

	return &jobString, err
}

func PostCallAsync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPPost(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func postCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PostCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	if len(job.AffectedResources) < 1 {
		return nil, nil
	}
	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func PatchCallAsync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPPatch(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func patchCallSync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	resJSONString, err := utils.HTTPPatch(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponse(resJSONString, err)
		log.WriteError(err)
		return err
	}

	err2 := json.Unmarshal([]byte(resJSONString), output)
	if err2 != nil {
		return fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return nil
}

func patchCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PatchCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	if len(job.AffectedResources) < 1 {
		return nil, nil
	}
	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func DeleteCallAsync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPDeleteWithBody(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func deleteCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := DeleteCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	if len(job.AffectedResources) < 1 {
		return nil, nil
	}
	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func downloadFileCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, toFilePath string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	return utils.HTTPDownloadFile(url, toFilePath, nil, &httpBasicAuth)
}

func postCallForm(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PostCallAsyncForm(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// TODO sng: CheckResponseAndWaitForJob needs handle the case where job state is not "Succeeded"
	// and return no http errot,
	// we then need to get the error from the job.error
	// and a step further to get details from the job events
	job, err := CheckResponseAndWaitForJobExt(storageSetting, jobString)
	if err != nil {

		// TODO sng: if the error msg is "" here
		// we need handle the case where job state is "Failed"

		if err.Error() == "" && job.State == "Failed" {
			return nil, CheckJobErrorInEvents(storageSetting, job)
		}

		log.WriteError(err)
		return nil, err
	}
	if len(job.AffectedResources) < 1 {
		return nil, nil
	}
	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func CheckJobErrorInEvents(storageSetting vssbmodel.StorageDeviceSettings, jobResponse *vssbmodel.JobResponse) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if jobResponse.Error.Message == "" {
		log.WriteDebug("TFDebug|failed Job %s has no job.error.message", jobResponse.JobID)
		return nil
	}

	// use the jobId to get the related events
	time.Sleep(1 * time.Second)
	jobId := jobResponse.JobID
	log.WriteDebug("TFDebug|Job %s failed with error: %s", jobId, jobResponse.Error.Message)

	// check the events for this jobID for more details
	type Event struct {
		EventName string `json:"eventName"`
		Message string `json:"message"`
	}

	type Events struct {
		Data []Event `json:"data"`
	}

	output := jobResponse
	fmt.Printf("state %s\n", output.State)
	if output.State == "Failed" {
		// get the events
		output2 := &Events{}	
		apiSuf := "objects/event-logs?severity=Error"

		err := getCall(storageSetting, apiSuf, output2)
		if err != nil {
			log.WriteDebug("TFDebug|failed to fetch events for Job %s", jobResponse.JobID)
			return nil
		}

		for _, event := range output2.Data {
			// fmt.Printf("Event: %s, Message: %s\n", event.EventName, event.Message)
			if strings.Contains(event.Message, output.JobID) {
				fmt.Printf("Event: %s\n", event.EventName)
				fmt.Printf("Message: %s\n", event.Message)
				return fmt.Errorf("job failed with event message: %s", event.Message)
			}
		}
	}

	return fmt.Errorf("job %s failed with error in jobResponse: %s", jobId, jobResponse.Error.Message)
}