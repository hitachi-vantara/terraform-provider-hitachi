package sanstorage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

func CheckGwyAsyncResponse(resJSONString *string) (*sanmodel.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug|resJSONString: %s", *resJSONString)

	var jobResponse sanmodel.JobResponse

	// example:
	// "errorSource" : "/ConfigurationManager/v1/objects/ldevs",
	// "message" : "An unsupported parameter is specified in the request body, or the hierarchy of the specified parameter is invalid (attribute = capacityInGB)",
	// "solution" : "Check whether the specified attribute is correct, and whether the attribute is specified in the correct hierarchy.",
	// "messageId" : "KART40038-E",
	// "detailCode" : "40038E-0"

	// check if response is error
	var gatewayError sanmodel.GatewayError
	err := json.Unmarshal([]byte(*resJSONString), &gatewayError)
	if err != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err)
	}

	if gatewayError.ErrorSource != "" {
		// failure
		jobResponse = sanmodel.JobResponse{
			Error: gatewayError,
		}
	} else {
		// job started
		err := json.Unmarshal([]byte(*resJSONString), &jobResponse)
		if err != nil {
			log.WriteDebug("TFError| error in Unmarshal call, err: %v", err)
			return nil, fmt.Errorf("failed to unmarshal json response: %+v", err)
		}
	}

	return &jobResponse, nil
}

func CheckJobStatus(storageSetting sanmodel.StorageDeviceSettings, jobID string) (*sanmodel.JobResponse, error) {
	// curl -k -v -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session $TOKEN" -X GET
	// https://mgmtIP/ConfigurationManager/v1/objects/jobs/21739
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, "objects/jobs/"+jobID)

	resJSONString, err := utils.HTTPGet(url, &headers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in HTTPGet call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)

	var jobResponse sanmodel.JobResponse

	err2 := json.Unmarshal([]byte(resJSONString), &jobResponse)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err2)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return &jobResponse, nil
}

func WaitForJobToComplete(storageSetting sanmodel.StorageDeviceSettings, jobResponse *sanmodel.JobResponse) (string, *sanmodel.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var FIRST_WAIT_TIME time.Duration = 1 // in sec
	MAX_RETRY_COUNT := 6

	// if r.status_code != http.client.ACCEPTED {
	// 	panic(errors.New("Exception")) //requests.HTTPError(r)
	// }

	var jobResult *sanmodel.JobResponse
	var err error
	status := "Initializing"
	retryCount := 1
	waitTime := FIRST_WAIT_TIME

	for status != "Completed" {
		if retryCount > MAX_RETRY_COUNT {
			err = fmt.Errorf("Exception: %v", "Timeout Error! Operation was not completed.")
			log.WriteError(err)
			log.WriteDebug("TFError| error in MAX RETRY COUNT condition, err: %v", err)
			return "", nil, err
		}

		time.Sleep(waitTime * time.Second)

		jobResult, err = CheckJobStatus(storageSetting, strconv.Itoa(jobResponse.JobID))
		if err != nil {
			log.WriteDebug("TFError| error in CheckJobStatus call, err: %v", err)
			return "", nil, err
		}

		status = jobResult.Status
		double_time := waitTime * 2
		if double_time < 120 {
			waitTime = double_time
		} else {
			waitTime = 120
		}
		retryCount += 1
	}

	// at this point, job status is completed

	if jobResult.State == "Failed" {
		log.WriteDebug("TFDebug|Error! SSB code : %+v", jobResult.Error.ErrorCode)
		//	fmt.Errorf("Exception: %v, %v", "Job Error!", jobResult.text))
		return jobResult.State, jobResult, nil
	}

	// otherwise state is Succeeded
	// log.WriteDebug("Async job was succeeded. affected resource : %v" + jobResult.AffectedResources[0])
	log.WriteDebug("TFDebug|Async job %v succeeded.", jobResponse.JobID)
	return jobResult.State, jobResult, nil
}

func CheckResponseAndWaitForJob(storageSetting sanmodel.StorageDeviceSettings, resJSONString *string) (*sanmodel.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pjobResponse, err := CheckGwyAsyncResponse(resJSONString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckGwyAsyncResponse call, err: %v", err)
		return pjobResponse, err
	}

	if pjobResponse.JobID == 0 {
		// job didn't start
		return pjobResponse, fmt.Errorf(pjobResponse.Error.Message)
	}

	state, job, err := WaitForJobToComplete(storageSetting, pjobResponse)
	if err != nil {
		log.WriteDebug("TFError| error in WaitForJobToComplete call, err: %v", err)
		return job, err
	}

	log.WriteDebug("TFDebug|Final JOB: %+v", job)

	if state != "Succeeded" {
		return job, CheckJobGatewayErrorResponse(job.Error, err)
	}

	return job, nil
}

func CheckJobGatewayErrorResponse(gwyError sanmodel.GatewayError, httpErr error) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	errmsg := ""
	if httpErr != nil {
		errmsg = fmt.Sprintf("%s\n", httpErr.Error())
		log.WriteDebug("TFError| error in HTTP response, err: %v", httpErr)
	}

	gwyJSON, err := json.MarshalIndent(gwyError, "", "  ")
	if err != nil {
		gwyJSON = []byte(fmt.Sprintf("%+v", gwyError))
	}
	errmsg = fmt.Sprintf("%s%s", errmsg, string(gwyJSON))
	if errmsg == "" {
		errmsg = "Failed but got no error message in response"
	}
	return fmt.Errorf("%s", errmsg)
}

func CheckHttpErrorResponse(resJSONString string, httpErr error) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	errmsg := ""
	if httpErr != nil {
		errmsg = fmt.Sprintf("%s\n", httpErr.Error())
		log.WriteDebug("TFError| error in HTTP response, err: %v", httpErr)
	}
	errmsg = fmt.Sprintf("%s%+v", errmsg, resJSONString)
	if errmsg == "" {
		errmsg = "Failed but got no error message in response"
	}
	return fmt.Errorf("%s", errmsg)
}
