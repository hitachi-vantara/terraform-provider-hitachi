package infra_gw

import (
	"encoding/json"
	"fmt"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"

	//sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

func CheckAsyncResponse(resJSONString *string) (*model.TaskResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug|resJSONString: %s", *resJSONString)

	var taskResponse model.TaskResponse

	err := json.Unmarshal([]byte(*resJSONString), &taskResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err)
	}

	return &taskResponse, nil
}

func CheckTaskStatus(storageSetting model.InfraGwSettings, taskId string) (*model.TaskResponse, error) {
	// curl -k -v -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Bearer $TOKEN" -X GET
	// https://<address>/porcelain/v2/tasks/<task_id>

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.Address, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	url := GetUrl(storageSetting.Address, "/tasks/"+taskId, false)

	resJSONString, err := utils.HTTPGet(url, &headers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in HTTPGet call, err: %v", err)
		return nil, err
	}

	// log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)

	var taskResponse model.TaskResponse

	err2 := json.Unmarshal([]byte(resJSONString), &taskResponse)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err2)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return &taskResponse, nil
}

func WaitForTaskToComplete(storageSetting model.InfraGwSettings, taskResponse *model.TaskResponse) (string, *model.TaskResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var FIRST_WAIT_TIME time.Duration = 1 // in sec
	MAX_RETRY_COUNT := 8

	// if r.status_code != http.client.ACCEPTED {
	// 	panic(errors.New("Exception")) //requests.HTTPError(r)
	// }

	var taskResult *model.TaskResponse
	var err error
	status := "Initializing"
	retryCount := 1
	waitTime := FIRST_WAIT_TIME

	for status != "Success" {
		if retryCount > MAX_RETRY_COUNT {
			err = fmt.Errorf("exception: %s, Timeout Error! Operation was not completed.", status)
			log.WriteError(err)
			log.WriteDebug("TFError| error in MAX RETRY COUNT condition, err: %v", err)
			return "", nil, err
		}

		time.Sleep(waitTime * time.Second)

		taskResult, err = CheckTaskStatus(storageSetting, taskResponse.Data.TaskId)
		if err != nil {
			log.WriteDebug("TFError| error in CheckTaskStatus call, err: %v", err)
			return "", nil, err
		}

		status = taskResult.Data.Status
		if status == "Failed" {
			log.WriteDebug("TFDebug|Error! SSB code : %+v", taskResult)
			return taskResult.Data.Status, taskResult, nil
		}
		double_time := waitTime * 2
		if double_time < 180 {
			waitTime = double_time
		} else {
			waitTime = 180
		}
		retryCount += 1
	}

	if taskResult.Data.Status != "Success" {
		log.WriteDebug("TFDebug|Error! SSB code : %+v", taskResult)
		return taskResult.Data.Status, taskResult, nil
	}

	log.WriteDebug("TFDebug|Async task %v succeeded.", taskResponse.Data.TaskId)
	return taskResult.Data.Status, taskResult, nil
}

func CheckResponseAndWaitForTask(storageSetting model.InfraGwSettings, resJSONString *string) (*model.TaskResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	taskResponse, err := CheckAsyncResponse(resJSONString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckGwyAsyncResponse call, err: %v", err)
		return taskResponse, err
	}
	state, task, err := WaitForTaskToComplete(storageSetting, taskResponse)
	if err != nil {
		log.WriteDebug("TFError| error in WaitForTaskToComplete call, err: %v", err)
		return task, err
	}

	log.WriteDebug("TFDebug|Final Task: %+v", task)

	if state != "Success" {
		if task.Data.Events == nil || len(task.Data.Events) == 0 {
			return task, fmt.Errorf("task %s", func() string {
				if len(task.Message) > 0 {
					return fmt.Sprintf("failed with unknown reason: %v", task.Message)
				}
				return "failed with unknown reason"
			}())
		}

		return task, fmt.Errorf(task.Data.Events[0].Description)
	}

	return task, nil
}
