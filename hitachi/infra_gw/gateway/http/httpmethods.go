package infra_gw

import (
	"encoding/json"
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	//sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetCall(storageSetting model.InfraGwSettings, apiSuf string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return err
	}

	url := GetUrl(storageSetting.Address, apiSuf)
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

func PostCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
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
	url := GetUrl(storageSetting.Address, apiSuf)

	// TODO: uncomment following when you need to work on lock and unlock resources
	// if reqBody == nil {
	// 	reqBodyInBytes = []byte{}
	// }

	taskString, err := utils.HTTPPost(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPost call, err: %v", err)
		return nil, err
	}

	return &taskString, err
}

func PostCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := PostCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in PostCallAsync call, err: %v", err)
		return nil, err
	}

	var response model.Response
	err = json.Unmarshal([]byte(*taskString), &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	task, err := CheckResponseAndWaitForTask(storageSetting, taskString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForTask call, task: %v err: %v", task, err)
		return nil, err
	}

	return &response.Data.ResourceId, nil
}

func PatchCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
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
	url := GetUrl(storageSetting.Address, apiSuf)

	taskString, err := utils.HTTPPatch(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPatch call, err: %v", err)
		return nil, err
	}

	return &taskString, err
}

func PatchCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := PatchCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in PatchCallAsync call, err: %v", err)
		return nil, err
	}

	var response model.Response
	err = json.Unmarshal([]byte(*taskString), &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	task, err := CheckResponseAndWaitForTask(storageSetting, taskString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForTask call, task: %v err: %v", task, err)
		return nil, err
	}

	return &response.Data.ResourceId, nil
}

func DeleteCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
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
	url := GetUrl(storageSetting.Address, apiSuf)

	taskString, err := utils.HTTPDeleteWithBody(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPDeleteWithBody call, err: %v", err)
		return nil, err
	}

	return &taskString, err
}

func DeleteCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := DeleteCallAsync(storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteCallAsync call, err: %v", err)
		return nil, err
	}

	var response model.Response
	err = json.Unmarshal([]byte(*taskString), &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	task, err := CheckResponseAndWaitForTask(storageSetting, taskString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForTask call, task: %v err: %v", task, err)
		return nil, err
	}

	return &response.Data.ResourceId, nil
}
