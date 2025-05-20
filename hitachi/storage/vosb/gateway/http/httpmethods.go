package vssbstorage

import (
	"encoding/json"
	"fmt"
	"strings"

	// "strconv"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

func GetCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
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
		log.WriteError(err)
		log.WriteDebug("TFError| error in HTTPGet call, err: %v", err)
		return err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)
	err2 := json.Unmarshal([]byte(resJSONString), output)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal, err: %v", err2)
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
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPPost(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPost call, err: %v", err)
		return nil, err
	}

	return &jobString, err
}

func PostCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PostCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in PostCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
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
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPPatch(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPatch call, err: %v", err)
		return nil, err
	}

	return &jobString, err
}

func PatchCallSync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return err
	}

	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	resJSONString, err := utils.HTTPPatch(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPatch call, err: %v", err)
		return err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)
	err2 := json.Unmarshal([]byte(resJSONString), output)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal, err: %v", err2)
		return fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return nil
}

func PatchCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := PatchCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in PatchCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
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
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.ClusterAddress, apiSuf)

	httpBasicAuth := utils.HttpBasicAuthentication{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
	}

	jobString, err := utils.HTTPDeleteWithBody(url, nil, reqBodyInBytes, &httpBasicAuth)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPDeleteWithBody call, err: %v", err)
		return nil, err
	}

	return &jobString, err
}

func DeleteCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	jobString, err := DeleteCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
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
