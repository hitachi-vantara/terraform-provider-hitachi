package sanstorage

import (
	"encoding/json"
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	telemetry "terraform-provider-hitachi/hitachi/common/telemetry"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// New wrapper functions for telemetry
func GetCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
	// Wrap the original function using telemetry.WrapMethod
	wrappedFunc := telemetry.WrapMethod("GetCall", getCall)
	// Type assert to the original function signature
	wrappedHTTPMethod := wrappedFunc.(func(sanmodel.StorageDeviceSettings, string, interface{}) error)
	// Now call the wrapped function
	return wrappedHTTPMethod(storageSetting, apiSuf, output)
}

func PostCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	// Wrap the original function using telemetry.WrapMethod
	wrappedFunc := telemetry.WrapMethod("PostCall", postCall)
	// Type assert to the original function signature
	wrappedHTTPMethod := wrappedFunc.(func(sanmodel.StorageDeviceSettings, string, interface{}) (*string, error))
	// Now call the wrapped function
	return wrappedHTTPMethod(storageSetting, apiSuf, reqBody)
}

func PatchCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	// Wrap the original function using telemetry.WrapMethod
	wrappedFunc := telemetry.WrapMethod("PatchCall", patchCall)
	// Type assert to the original function signature
	wrappedHTTPMethod := wrappedFunc.(func(sanmodel.StorageDeviceSettings, string, interface{}) (*string, error))
	// Now call the wrapped function
	return wrappedHTTPMethod(storageSetting, apiSuf, reqBody)
}

func DeleteCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	// Wrap the original function using telemetry.WrapMethod
	wrappedFunc := telemetry.WrapMethod("DeleteCall", deleteCall)
	// Type assert to the original function signature
	wrappedHTTPMethod := wrappedFunc.(func(sanmodel.StorageDeviceSettings, string, interface{}) (*string, error))
	// Now call the wrapped function
	return wrappedHTTPMethod(storageSetting, apiSuf, reqBody)
}

///////////////
func getCall(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return err
	}

	url := GetUrl(storageSetting.MgmtIP, apiSuf)
	resJSONString, err := utils.HTTPGet(url, &headers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPGet call, err: %v", err)
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

func PostCallAsync(storageSetting sanmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	// TODO: uncomment following when you need to work on lock and unlock resources
	// if reqBody == nil {
	// 	reqBodyInBytes = []byte{}
	// }

	jobString, err := utils.HTTPPost(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPost call, err: %v", err)
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
		log.WriteDebug("TFError| error in PostCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
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
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	jobString, err := utils.HTTPPatch(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPatch call, err: %v", err)
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
		log.WriteDebug("TFError| error in PatchCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
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
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.MgmtIP, apiSuf)

	jobString, err := utils.HTTPDeleteWithBody(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPDeleteWithBody call, err: %v", err)
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
		log.WriteDebug("TFError| error in DeleteCallAsync call, err: %v", err)
		return nil, err
	}

	job, err := CheckResponseAndWaitForJob(storageSetting, jobString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForJob call, err: %v", err)
		return nil, err
	}

	sarr := strings.Split(job.AffectedResources[0], "/")
	affRes := sarr[len(sarr)-1]
	log.WriteDebug("TFDebug|affRes=%+v\n", affRes)

	return &affRes, nil
}
