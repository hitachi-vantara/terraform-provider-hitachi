package vssbstorage

import (
	"encoding/json"
	"fmt"
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
