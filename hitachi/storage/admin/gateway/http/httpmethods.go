package admin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func getCall(storageSetting model.StorageDeviceSettings, apiSuf string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		return err
	}

	url := GetUrl(storageSetting.MgmtIP, apiSuf)
	resJSONString, err := utils.HTTPGet(url, &headers)
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

func PostCallAsync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	jobString, err := utils.HTTPPost(url, &headers, reqBodyInBytes)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func postCallNoRetry(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	resIds := GetOperationDetailResourceIDs(job)
	return resIds, nil
}

func postCall(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	return utils.ExecuteWithRetry("postCall", func() (*string, error) {
		return postCallNoRetry(storageSetting, apiSuf, reqBody)
	})
}

// This is for both sync and async Patch
func PatchCallSyncAsync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	jobString, err := utils.HTTPPatch(url, &headers, reqBodyInBytes)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func patchCallAsyncNoRetry(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PatchCallSyncAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	resIds := GetOperationDetailResourceIDs(job)
	return resIds, nil
}

func patchCallAsync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	return utils.ExecuteWithRetry("patchCallAsync", func() (*string, error) {
		return patchCallAsyncNoRetry(storageSetting, apiSuf, reqBody)
	})
}

func patchCallSync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	resJSONString, err := PatchCallSyncAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// only affectedResources and operationDetails are populated
	var jobResponse model.JobResponse
	err2 := json.Unmarshal([]byte(*resJSONString), &jobResponse)
	if err2 != nil {
		log.WriteError(err)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	resIds := GetOperationDetailResourceIDs(&jobResponse)
	return resIds, nil
}

func DeleteCallAsync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	jobString, err := utils.HTTPDeleteWithBody(url, &headers, reqBodyInBytes)
	if err != nil {
		err := CheckHttpErrorResponse(jobString, err)
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func deleteCallNoRetry(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	resIds := GetOperationDetailResourceIDs(job)
	return resIds, nil
}

func deleteCall(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	return utils.ExecuteWithRetry("deleteCall", func() (*string, error) {
		return deleteCallNoRetry(storageSetting, apiSuf, reqBody)
	})
}

func GetAffectedResourceSuffix(job *model.JobResponse) *string {
	log := commonlog.GetLogger()

	if job == nil || len(job.AffectedResources) < 1 {
		return nil
	}

	var suffixes []string
	for _, res := range job.AffectedResources {
		parts := strings.Split(res, "/")
		if len(parts) > 0 {
			suffixes = append(suffixes, parts[len(parts)-1])
		}
	}

	affRes := strings.Join(suffixes, ",") // concatenate with comma
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes
}

func GetOperationDetailResourceIDs(job *model.JobResponse) *string {
	log := commonlog.GetLogger()

	if job == nil || len(job.OperationDetails) < 1 {
		return nil
	}

	var resourceIDs []string
	for _, detail := range job.OperationDetails {
		if detail.ResourceID != "" {
			resourceIDs = append(resourceIDs, detail.ResourceID)
		}
	}

	if len(resourceIDs) == 0 {
		return nil
	}

	sort.Strings(resourceIDs)
	concatenated := strings.Join(resourceIDs, ",") // join with comma
	log.WriteDebug("TFDebug|Concatenated ResourceIDs=%+v\n", concatenated)

	return &concatenated
}
