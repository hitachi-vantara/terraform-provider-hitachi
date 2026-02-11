package sanstorage

import (
	"fmt"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
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

// GetRangeOfLunsWithOptions to get Lun information with options
func (psm *sanStorageManager) GetRangeOfLunsWithOptions(startLdevID int, endLdevID int, isUndefinedLdev bool, filterOption string, detailInfoType string) (*[]sanmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var logicalUnits sanmodel.LogicalUnits

	// Build API endpoint with query parameters
	apiSuf := "objects/ldevs"
	params := make([]string, 0)

	// Add range parameters
	if startLdevID >= 0 {
		params = append(params, fmt.Sprintf("startLdevId=%d", startLdevID))
	}
	if endLdevID >= 0 && endLdevID >= startLdevID {
		params = append(params, fmt.Sprintf("endLdevId=%d", endLdevID))
	}

	// Add filterOption parameter
	if filterOption != "" {
		params = append(params, fmt.Sprintf("ldevOption=%s", filterOption))
	}

	// Add detailInfoType parameter
	if detailInfoType != "" {
		params = append(params, fmt.Sprintf("detailInfoType=%s", detailInfoType))
	}

	// Add isUndefinedLdev parameter
	if isUndefinedLdev {
		params = append(params, "isUndefinedLdev=true")
	}

	if len(params) > 0 {
		apiSuf += "?" + strings.Join(params, "&")
	}

	log.WriteDebug("TFDebug| API call: %s", apiSuf)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &logicalUnits)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| API response: received %d volumes", len(logicalUnits.ListOfLun))
	if len(logicalUnits.ListOfLun) > 0 {
		firstVol := logicalUnits.ListOfLun[0]
		log.WriteDebug("TFDebug| First volume: ldevId=%d, cylinder=%d, status=%s, attributes=%v",
			firstVol.LdevID, firstVol.Cylinder, firstVol.Status, firstVol.Attributes)
	}

	return &logicalUnits.ListOfLun, nil
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

// SetEseVolume toggles ESE enablement for a mainframe volume.
func (psm *sanStorageManager) SetEseVolume(ldevID int, isEseVolume bool) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d/actions/set-ese/invoke", ldevID)
	reqBody := sanmodel.SetEseVolumeRequestGwy{
		Parameters: sanmodel.SetEseVolumeParameters{
			IsEseVolume: isEseVolume,
		},
	}
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
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

// FormatLdev invokes the format action on an LDEV
func (psm *sanStorageManager) FormatLdev(reqBody sanmodel.FormatLdevRequestGwy, ldevID int) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d/actions/format/invoke", ldevID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &ldevID, nil
}

// BlockLun changes LDEV status to 'blk' via change-status action
func (psm *sanStorageManager) BlockLun(ldevID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d/actions/change-status/invoke", ldevID)
	reqBody := sanmodel.ChangeStatusRequestGwy{
		Parameters: sanmodel.ChangeStatusParameters{
			Status: "blk",
		},
	}

	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// UnblockLun changes LDEV status back to 'nml' via change-status action
func (psm *sanStorageManager) UnblockLun(ldevID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ldevs/%d/actions/change-status/invoke", ldevID)
	reqBody := sanmodel.ChangeStatusRequestGwy{
		Parameters: sanmodel.ChangeStatusParameters{
			Status: "nml",
		},
	}

	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// StopAllVolumeFormat invokes the appliance-wide stop-format action which attempts
// to stop normal format operations for all volumes.
func (psm *sanStorageManager) StopAllVolumeFormat() error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "services/ldev-service/actions/stop-format/invoke"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
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

	retryCfg := utils.GetRetryConfig()

	MAX_RETRY_COUNT := retryCfg.MaxRetries
	var err error
	retryCount := 1
	waitTime := time.Duration(retryCfg.Delay)

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
