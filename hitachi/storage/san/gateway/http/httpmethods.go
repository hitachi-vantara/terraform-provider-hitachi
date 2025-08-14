package sanstorage

import (
	"encoding/json"
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

func getCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
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

func PostCallAsync(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	// TODO: uncomment following when you need to work on lock and unlock resources
	// if reqBody == nil {
	// 	reqBodyInBytes = []byte{}
	// }

	jobString, err := utils.HTTPPost(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func postCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func PatchCallAsync(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func patchCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}

func DeleteCallAsync(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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
		log.WriteError(err)
		return nil, err
	}

	return &jobString, err
}

func deleteCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}
