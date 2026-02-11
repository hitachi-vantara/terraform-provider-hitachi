package sanstorage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func CheckGwyAsyncResponse(resJSONString *string) (*sanmodel.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("resJSONString: %s", *resJSONString)

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
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, "objects/jobs/"+jobID)

	resJSONString, err := utils.HTTPGet(url, &headers)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("resJSONString: %s", resJSONString)

	var jobResponse sanmodel.JobResponse

	err2 := json.Unmarshal([]byte(resJSONString), &jobResponse)
	if err2 != nil {
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

	// Increase timeout for long-running format (normal / FMT) operations.
	// Detect format invoke from the original request URL and inspect the
	// request body to distinguish FMT (normal) from QFMT (quick).
	if jobResponse != nil && jobResponse.Request.RequestURL != "" {
		// Detect format invoke requests. Cover both appliance-level and
		// per-LDEV endpoints such as:
		//   /ConfigurationManager/v1/objects/ldevs/{id}/actions/format/invoke
		if strings.Contains(jobResponse.Request.RequestURL, "format/invoke") || (strings.Contains(jobResponse.Request.RequestURL, "ldevs") && strings.Contains(jobResponse.Request.RequestURL, "format/invoke")) {
			body := jobResponse.Request.RequestBody
			if body != "" {
				ub := strings.ToUpper(body)
				// If the request explicitly asked for QFMT, keep default timeout.
				// Otherwise if it asked for FMT (normal format) increase initial interval.
				if strings.Contains(ub, "FMT") {
					// Normal format (FMT) can take much longer; increase initial interval to reduce polling.
					FIRST_WAIT_TIME = 5
				}
			}
		}
	}

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
			return "", nil, err
		}

		time.Sleep(waitTime * time.Second)

		jobResult, err = CheckJobStatus(storageSetting, strconv.Itoa(jobResponse.JobID))
		if err != nil {
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
		log.WriteDebug("Error! SSB code : %+v", jobResult.Error.ErrorCode)
		//	fmt.Errorf("Exception: %v, %v", "Job Error!", jobResult.text))
		return jobResult.State, jobResult, nil
	}

	// otherwise state is Succeeded
	// log.WriteDebug("Async job was succeeded. affected resource : %v" + jobResult.AffectedResources[0])
	log.WriteDebug("Async job %v succeeded.", jobResponse.JobID)
	return jobResult.State, jobResult, nil
}

func CheckResponseAndWaitForJob(storageSetting sanmodel.StorageDeviceSettings, resJSONString *string) (*sanmodel.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pjobResponse, err := CheckGwyAsyncResponse(resJSONString)
	if err != nil {
		return pjobResponse, err
	}

	if pjobResponse.JobID == 0 {
		// job didn't start
		return pjobResponse, fmt.Errorf(pjobResponse.Error.Message)
	}

	state, job, err := WaitForJobToComplete(storageSetting, pjobResponse)
	if err != nil {
		return job, err
	}

	log.WriteDebug("Final JOB: %+v", job)

	if state != "Succeeded" {
		return job, fmt.Errorf(job.Error.Message)
	}

	return job, nil
}
