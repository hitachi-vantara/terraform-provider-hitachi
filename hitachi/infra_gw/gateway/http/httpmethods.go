package infra_gw

import (
	"encoding/json"
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	//sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetCall(storageSetting model.InfraGwSettings, apiSuf string, reqHeaders *map[string]string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return err
	}

	if reqHeaders != nil {
		for key, value := range *reqHeaders {
			headers[key] = value
		}
	}

	url := GetUrl(storageSetting.Address, apiSuf, storageSetting.V3API)
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

func PostCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, reqHeaders *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	if reqHeaders != nil {
		for key, value := range *reqHeaders {
			headers[key] = value
		}
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.Address, apiSuf, storageSetting.V3API)

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

func PostCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, redHeader *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := PostCallAsync(storageSetting, apiSuf, reqBody, redHeader)
	if err != nil {
		log.WriteDebug("TFError| error in PostCallAsync call, err: %v", err)
		return nil, err
	}
	

	return MakeFinalResponse(storageSetting, taskString)
}

func PatchCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, reqHeaders *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	if reqHeaders != nil {
		for key, value := range *reqHeaders {
			headers[key] = value
		}
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.Address, apiSuf, storageSetting.V3API)

	taskString, err := utils.HTTPPatch(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPPatch call, err: %v", err)
		return nil, err
	}

	return &taskString, err
}

func PatchCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, reqHeaders *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := PatchCallAsync(storageSetting, apiSuf, reqBody, reqHeaders)
	if err != nil {
		log.WriteDebug("TFError| error in PatchCallAsync call, err: %v", err)
		return nil, err
	}

	return MakeFinalResponse(storageSetting, taskString)
}

func DeleteCallAsync(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, reqHeaders *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	if reqHeaders != nil {
		for key, value := range *reqHeaders {
			headers[key] = value
		}
	}

	reqBodyInBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|reqBodyInBytes: %s\n", string(reqBodyInBytes))
	url := GetUrl(storageSetting.Address, apiSuf, storageSetting.V3API)

	taskString, err := utils.HTTPDeleteWithBody(url, &headers, reqBodyInBytes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in utils.HTTPDeleteWithBody call, err: %v", err)
		return nil, err
	}

	return &taskString, err
}

func DeleteCall(storageSetting model.InfraGwSettings, apiSuf string, reqBody interface{}, reqHeaders *map[string]string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskString, err := DeleteCallAsync(storageSetting, apiSuf, reqBody, reqHeaders)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteCallAsync call, err: %v", err)
		return nil, err
	}

	return MakeFinalResponse(storageSetting, taskString)
}

func MakeFinalResponse(storageSetting model.InfraGwSettings, taskString *string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.Response
	err := json.Unmarshal([]byte(*taskString), &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in Marshal call, err: %v", err)
		return nil, err
	}

	if response.Data.TaskId == "" {

		var basicResponse model.BasicResponse
		err := json.Unmarshal([]byte(*taskString), &basicResponse)
		if err != nil {
			log.WriteError(err)
			log.WriteDebug("TFError| error in Marshal call, err: %v", err)
			return nil, err
		}
		if basicResponse.TaskId == "" {
			return taskString, nil
		}
		response.Data.TaskId = basicResponse.TaskId
		response.Data.ResourceId = basicResponse.ResourceId
		response.Data.State = basicResponse.State

		reString, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		taskString = func(s string) *string { return &s }(string(reString))
	}

	task, err := CheckResponseAndWaitForTask(storageSetting, taskString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckResponseAndWaitForTask call, task: %v err: %v", task, err)
		return nil, err
	}

	if response.Data.ResourceId == "" {

		for _, attr := range task.Data.AdditionalAttributes {
			if attr.Type == "resource" {
				log.WriteDebug("TFDebug|resource: %+v", attr)
				response.Data.ResourceId = attr.Id
			}
		}
	}

	return &response.Data.ResourceId, nil
}
