package sanstorage

import (
	"fmt"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"time"
)

// GetLun to get Lun information
func (psm *sanStorageManager) GetLun(ldevID int) (*sanmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var logicalUnit sanmodel.LogicalUnit
	apiSuf := fmt.Sprintf("objects/ldevs/%d", ldevID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &logicalUnit)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &logicalUnit, nil
}

// GetAllLun to get all Lun information
func (psm *sanStorageManager) GetAllLun() (*sanmodel.LogicalUnits, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var logicalUnits sanmodel.LogicalUnits
	apiSuf := "objects/ldevs"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &logicalUnits)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &logicalUnits, nil
}

// CreateLun is use to create new lun
func (psm *sanStorageManager) CreateLun(reqBody sanmodel.CreateLunRequestGwy) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/ldevs"
	affRes, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	arr := strings.Split(*affRes, ",")
	sldevID := arr[0]
	ldevID, _ := strconv.Atoi(sldevID)
	log.WriteDebug("TFDebug | ldevID= %d", ldevID)
	return &ldevID, nil
}

// UpdateLun used to update lun
func (psm *sanStorageManager) UpdateLun(reqBody sanmodel.UpdateLunRequestGwy, ldevID int) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d", ldevID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &ldevID, nil
}

// ExpandLun used to expand lun
func (psm *sanStorageManager) ExpandLun(reqBody sanmodel.ExpandLunRequestGwy, ldevID int) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d/actions/expand/invoke", ldevID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &ldevID, nil
}

// DeleteLun is used to delete lun
func (psm *sanStorageManager) DeleteLun(ldevID int, capacitySaving bool) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// additional delete req body
	reqBody := map[string]bool{}
	if capacitySaving {
		reqBody = map[string]bool{
			"isDataReductionDeleteForceExecute": true,
		}
	}

	apiSuf := fmt.Sprintf("objects/ldevs/%d", ldevID)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	// Deleting data on a DP volume for which the capacity saving function (compression or deduplication) is enabled takes time.
	// Use the status of the target resource rather than the status of the job to check whether the volume has been deleted.
	err = psm.WaitUntilLunIsDeleted(psm.storageSetting, ldevID)
	if err != nil {
		log.WriteDebug("TFError| error in calling WaitUntilLunIsDeleted, err: %v", err)
		return err
	}

	return nil
}

// WaitUntilLunIsDeleted is internal fun call from delete lun. It wait untill Lun is deleted
func (psm *sanStorageManager) WaitUntilLunIsDeleted(storageSetting sanmodel.StorageDeviceSettings, ldevID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var FIRST_WAIT_TIME time.Duration = 1 // in sec
	MAX_RETRY_COUNT := 10

	var err error
	retryCount := 1
	waitTime := FIRST_WAIT_TIME

	for err != nil {
		if retryCount > MAX_RETRY_COUNT {
			err = fmt.Errorf("exception: %v", "timeout error! operation was not completed.")
			log.WriteError(err)
			log.WriteDebug("TFError| error in condition retryCount > MAX_RETRY_COUNT, err: %v", err)
			return err
		}

		time.Sleep(waitTime * time.Second)

		lun, err := psm.GetLun(ldevID)
		if err != nil {
			// deletion could have been done
			log.WriteDebug("TFError| error in GetLun, err: %v", err)
			return err
		}
		log.WriteDebug("TFDebug|Lun Deletion: %+v\n", lun)

		double_time := waitTime * 2
		if double_time < 45 {
			waitTime = double_time
		} else {
			waitTime = 45
		}
		retryCount += 1
	}
	return nil
}
