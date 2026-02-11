package admin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	config "terraform-provider-hitachi/hitachi/common/config"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func CheckGwyAsyncResponse(resJSONString *string) (*model.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug|resJSONString: %s", *resJSONString)

	var asyncStatusArray []model.AsyncCommandStatus
	var asyncStatusSingle model.AsyncCommandStatus
	var gatewayError model.GatewayError
	var jobResponse model.JobResponse

	raw := []byte(*resJSONString)

	// 1. Try AsyncCommandStatus array
	if err := json.Unmarshal(raw, &asyncStatusArray); err == nil && len(asyncStatusArray) > 0 {
		njobs := len(asyncStatusArray)
		if njobs == 1 {
			// single job
			parts := strings.Split(asyncStatusArray[0].StatusResource, "/")
			jobid, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse JobID from StatusResource: %w", err)
			}
			jobResponse = model.JobResponse{
				JobID: jobid,
			}
			return &jobResponse, nil
		} else if njobs > 1 {
			// multiple jobs
			mjobs := []int{}
			for _, v := range asyncStatusArray {
				parts := strings.Split(v.StatusResource, "/")
				jobid, err := strconv.Atoi(parts[len(parts)-1])
				if err != nil {
					return nil, fmt.Errorf("failed to parse JobID from StatusResource: %w", err)
				}
				mjobs = append(mjobs, jobid)
			}
			jobResponse = model.JobResponse{
				MultipleJobs: mjobs,
			}
			return &jobResponse, nil
		} else {
			return nil, fmt.Errorf("found no StatusResource")
		}
	}

	// 2. Try AsyncCommandStatus single
	if err := json.Unmarshal(raw, &asyncStatusSingle); err == nil && asyncStatusSingle.StatusResource != "" {
		parts := strings.Split(asyncStatusSingle.StatusResource, "/")
		jobid, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse JobID from StatusResource: %w", err)
		}

		jobResponse = model.JobResponse{
			JobID: jobid,
		}
		return &jobResponse, nil
	}

	// 3. Try GatewayError
	if err := json.Unmarshal(raw, &gatewayError); err == nil {
		if gatewayError.ErrorSource != "" {
			jobResponse = model.JobResponse{
				Error: gatewayError,
			}
			return &jobResponse, nil
		}
	}

	// 4. Try JobResponse
	if err := json.Unmarshal(raw, &jobResponse); err != nil {
		log.WriteDebug("TFError| failed to unmarshal JSON response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal JSON response: %+v", err)
	}

	return &jobResponse, nil
}

func CheckJobStatus(storageSetting model.StorageDeviceSettings, jobID string) (*model.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	headers, err := GetAuthTokenHeader(storageSetting.MgmtIP, storageSetting.Username, storageSetting.Password)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in GetAuthTokenHeader call, err: %v", err)
		return nil, err
	}

	url := GetUrl(storageSetting.MgmtIP, "objects/command-status/"+jobID)

	resJSONString, err := utils.HTTPGet(url, &headers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in HTTPGet call, err: %v", err)
		return nil, err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)

	var jobResponse model.JobResponse

	err2 := json.Unmarshal([]byte(resJSONString), &jobResponse)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err2)
		return nil, fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return &jobResponse, nil
}

func CheckResponseAndWaitForJob(storageSetting model.StorageDeviceSettings, resJSONString *string) (*model.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pjobResponse, err := CheckGwyAsyncResponse(resJSONString)
	if err != nil {
		log.WriteDebug("TFError| error in CheckGwyAsyncResponse call, err: %v", err)
		return pjobResponse, err
	}

	if pjobResponse.JobID == 0 && len(pjobResponse.MultipleJobs) == 0 {
		if pjobResponse.Progress == "" && pjobResponse.Status == "" && pjobResponse.ErrorResource == "" && pjobResponse.ErrorMessage == "" {
			log.WriteDebug("No command status found or JobID and no error or no status/progress found in response, treat this as success.")
			return pjobResponse, nil
		} else {
			// job didn't start
			return pjobResponse, fmt.Errorf(pjobResponse.Error.Message)
		}
	}

	// for single job
	if pjobResponse.JobID != 0 && len(pjobResponse.MultipleJobs) == 0 {
		_, job, err := WaitForJobToComplete(storageSetting, pjobResponse)
		if err != nil {
			log.WriteDebug("TFError| error in WaitForJobToComplete call, err: %v", err)
			return job, err
		}

		log.WriteDebug("TFDebug|Final JOB: %+v", job)

		// if state != "normal" {
		// 	return job, CheckJobGatewayErrorResponse(job.Error, err)
		// }
		return job, nil
	}

	// for multiple jobs
	if pjobResponse.JobID == 0 && len(pjobResponse.MultipleJobs) > 0 {
		aggregatedJobResponse, err := WaitForMultipleJobsToComplete(storageSetting, pjobResponse.MultipleJobs)
		if err != nil {
			log.WriteDebug("TFError| error in WaitForMultipleJobsToComplete call, err: %v", err)
			return aggregatedJobResponse, err
		}

		log.WriteDebug("TFDebug|Final JOB: %+v", aggregatedJobResponse)
		return aggregatedJobResponse, nil
	}

	return pjobResponse, nil
}

func CheckJobGatewayErrorResponse(gwyError model.GatewayError, httpErr error) error {
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

// --- common helpers ---

const (
	firstWaitTime = 1
	maxWaitTime   = 120
)

func calculateMaxRetryCount(apiTimeoutSec int) int {
	waitTime := firstWaitTime
	totalWait := 0
	count := 0

	for totalWait < apiTimeoutSec {
		totalWait += waitTime
		count++
		if waitTime < maxWaitTime {
			waitTime *= 2
			if waitTime > maxWaitTime {
				waitTime = maxWaitTime
			}
		}
	}
	return count
}

func getApiTimeoutSec() int {
	apiTimeoutSec := config.DEFAULT_API_TIMEOUT
	if config.ConfigData != nil && config.ConfigData.APITimeout > 0 {
		apiTimeoutSec = config.ConfigData.APITimeout
	}
	return apiTimeoutSec
}

func monitorJob(storageSetting model.StorageDeviceSettings, jobId int, maxRetry int, log commonlog.ILogger) (*model.JobResponse, error) {
	progress := "processing"
	retryCount := 1
	waitTime := firstWaitTime
	var jobResult *model.JobResponse
	var err error

	for progress != "completed" {
		if retryCount > maxRetry {
			return nil, fmt.Errorf("Timeout Error! Job %d was not completed.", jobId)
		}

		log.WriteDebug("JobID=%d | attempt=%d | waitTime=%v", jobId, retryCount, waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		jobResult, err = CheckJobStatus(storageSetting, strconv.Itoa(jobId))
		if err != nil {
			return nil, fmt.Errorf("JobID=%d: %w", jobId, err)
		}

		progress = jobResult.Progress
		if waitTime*2 < maxWaitTime {
			waitTime *= 2
		} else {
			waitTime = maxWaitTime
		}
		retryCount++
	}
	return jobResult, nil
}

func WaitForJobToComplete(storageSetting model.StorageDeviceSettings, jobResponse *model.JobResponse) (string, *model.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiTimeoutSec := getApiTimeoutSec()
	maxRetry := calculateMaxRetryCount(apiTimeoutSec)

	log.WriteDebug("APITimeout=%v MAX_RETRY_COUNT=%v jobId=%v", apiTimeoutSec, maxRetry, jobResponse.JobID)

	jobResult, err := monitorJob(storageSetting, jobResponse.JobID, maxRetry, log)
	if err != nil {
		log.WriteError(err)
		return "", nil, err
	}

	if jobResult.Status != "normal" {
		log.WriteDebug("TFDebug|Error! SSB code : %+v", jobResult.ErrorCode)
		return jobResult.Status, jobResult, fmt.Errorf(jobResult.ErrorMessage)
	}

	log.WriteDebug("TFDebug|Async job %v succeeded.", jobResponse.JobID)
	return jobResult.Status, jobResult, nil
}

func WaitForMultipleJobsToComplete(storageSetting model.StorageDeviceSettings, jobIds []int) (*model.JobResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiTimeoutSec := getApiTimeoutSec()
	maxRetry := calculateMaxRetryCount(apiTimeoutSec)

	log.WriteDebug("APITimeout=%v MAX_RETRY_COUNT=%v for %d jobs", apiTimeoutSec, maxRetry, len(jobIds))

	type jobResultWrapper struct {
		jobId int
		resp  *model.JobResponse
		err   error
	}

	resultCh := make(chan jobResultWrapper, len(jobIds))
	var wg sync.WaitGroup

	for _, jobId := range jobIds {
		wg.Add(1)
		go func(jId int) {
			defer wg.Done()
			resp, err := monitorJob(storageSetting, jId, maxRetry, log)
			resultCh <- jobResultWrapper{jobId: jId, resp: resp, err: err}
		}(jobId)
	}

	wg.Wait()
	close(resultCh)

	aggregated := &model.JobResponse{
		Status:            "normal",
		MultipleJobs:      jobIds,
		Progress:          "completed",
		AffectedResources: []string{},
		OperationDetails:  []model.OperationDetail{},
	}

	var errorMessages []string
	var errorResources []string
	var pollingErrors []string
	var uniqueSSB1s []string
	var uniqueSSB2s []string
	firstJobId := 0
	firstFailedProgress := ""

	for res := range resultCh {
		if res.resp != nil {
			aggregated.AffectedResources = append(aggregated.AffectedResources, res.resp.AffectedResources...)
			aggregated.OperationDetails = append(aggregated.OperationDetails, res.resp.OperationDetails...)

			if res.resp.Status != "normal" {
				aggregated.Status = "error" // not "failed"

				if firstJobId == 0 {
					firstJobId = res.resp.JobID
					firstFailedProgress = res.resp.Progress
				}

				// Deduplicate error message
				if res.resp.ErrorMessage != "" && !containsString(errorMessages, res.resp.ErrorMessage) {
					errorMessages = append(errorMessages, res.resp.ErrorMessage)
				}

				// Deduplicate error resource
				if res.resp.ErrorResource != "" && !containsString(errorResources, res.resp.ErrorResource) {
					errorResources = append(errorResources, res.resp.ErrorResource)
				}

				// Deduplicate SSB1 and SSB2 codes
				if res.resp.ErrorCode.SSB1 != "" && !containsString(uniqueSSB1s, res.resp.ErrorCode.SSB1) {
					uniqueSSB1s = append(uniqueSSB1s, res.resp.ErrorCode.SSB1)
				}
				if res.resp.ErrorCode.SSB2 != "" && !containsString(uniqueSSB2s, res.resp.ErrorCode.SSB2) {
					uniqueSSB2s = append(uniqueSSB2s, res.resp.ErrorCode.SSB2)
				}
			}
		}

		if res.err != nil {
			aggregated.Status = "error"
			if firstJobId == 0 {
				firstJobId = res.jobId
				firstFailedProgress = "processing"
			}
			pollingErrors = append(pollingErrors, res.err.Error())
		}
	}

	if aggregated.Status == "error" {
		aggregated.JobID = firstJobId
		aggregated.Progress = firstFailedProgress

		// Combine SSBs if any
		if len(uniqueSSB1s) > 0 || len(uniqueSSB2s) > 0 {
			aggregated.ErrorCode = model.ErrorCode{
				SSB1: strings.Join(uniqueSSB1s, ","),
				SSB2: strings.Join(uniqueSSB2s, ","),
			}
		}

		if len(errorMessages) > 0 {
			aggregated.ErrorMessage = strings.Join(errorMessages, "; ")
		}
		if len(errorResources) > 0 {
			aggregated.ErrorResource = strings.Join(errorResources, "; ")
		}

		// return the errorMessage as err (if present)
		if aggregated.ErrorMessage != "" {
			return aggregated, fmt.Errorf(aggregated.ErrorMessage)
		}
	}

	if len(pollingErrors) > 0 {
		return aggregated, fmt.Errorf("Polling errors: %s", strings.Join(pollingErrors, "; "))
	}

	return aggregated, nil
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
