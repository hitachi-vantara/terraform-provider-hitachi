package admin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	provmanager "terraform-provider-hitachi/hitachi/storage/admin/provisioner"
	provimpl "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
)

func (psm *adminStorageManager) ReconcileReadAdminVolumes(volumeIDs []int) ([]gwymodel.VolumeInfoByID, []int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if len(volumeIDs) == 0 {
		return nil, nil, fmt.Errorf("volume IDs are empty")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {                   
		return nil, nil, err
	}

	var mu sync.Mutex
	var volumeInfos []gwymodel.VolumeInfoByID
	var existingIDs []int
	var firstErr error

	errs := utils.RunConcurrentOperations("Read", volumeIDs,
		func(volID int, idx int) error {
			volumeInfo, err := provObj.GetVolumeByID(volID)
			if err != nil {
				if IsNotFoundError(err) {
					log.WriteWarn("[Read #%d] Volume %d not found", idx+1, volID)
					return nil
				}
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("failed to fetch volume %d: %w", volID, err)
				}
				mu.Unlock()
				return err
			}

			if volumeInfo != nil {
				mu.Lock()
				volumeInfos = append(volumeInfos, *volumeInfo)
				existingIDs = append(existingIDs, volID)
				mu.Unlock()
			}
			return nil
		})

	for _, e := range errs {
		if e != nil && firstErr == nil {
			firstErr = e
		}
	}

	if len(existingIDs) == 0 {
		return nil, nil, fmt.Errorf("no existing volumes found")
	}

	sort.Slice(volumeInfos, func(i, j int) bool { return volumeInfos[i].ID < volumeInfos[j].ID })
	sort.Ints(existingIDs)

	return volumeInfos, existingIDs, firstErr
}

func (psm *adminStorageManager) ReconcileDeleteAdminVolumes(volumeIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if len(volumeIDs) == 0 {
		return fmt.Errorf("no volume IDs provided")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	errs := utils.RunConcurrentOperations("Delete", volumeIDs,
		func(volID int, idx int) error {
			err := provObj.DeleteVolume(volID)
			if err != nil {
				if IsNotFoundError(err) {
					log.WriteWarn("[Delete #%d] Volume %d not found or already deleted", idx+1, volID)
					return nil
				}
				return fmt.Errorf("volume %d: %v", volID, err)
			}
			return nil
		})

	var allErrs []string
	for _, e := range errs {
		if e != nil {
			allErrs = append(allErrs, e.Error())
		}
	}

	if len(allErrs) > 0 {
		errMsg := fmt.Sprintf("one or more volumes failed to delete: %s", strings.Join(allErrs, "; "))
		log.WriteError(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	log.WriteInfo("All volumes deleted successfully.")
	return nil
}

func (psm *adminStorageManager) ReconcileCreateAdminVolumes(params gwymodel.CreateVolumeParams) ([]int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume creation reconciliation")

	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Params: %v", string(b))

	volumeIDsStr, err := psm.createVolumes(params)
	if err != nil {
		log.WriteError("%v", err)
		return nil, err
	}

	ids, err := parseVolumeIDs(volumeIDsStr)
	if err != nil {
		return nil, err
	}

	log.WriteInfo("volumes created successfully: %v", ids)
	return ids, nil
}

func (psm *adminStorageManager) ReconcileUpdateAdminVolume(volumeID int, params gwymodel.CreateVolumeParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume update reconciliation: ID=%d", volumeID)

	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Update Params: %v", string(b))

	err := psm.updateVolume(volumeID, params)
	if err != nil {
		log.WriteError("failed to update volume %d: %v", volumeID, err)
		return err
	}

	log.WriteInfo(fmt.Sprintf("volume %d updated successfully", volumeID))
	return nil
}

// ------------------- Helpers -------------------
func (psm *adminStorageManager) createVolumes(params gwymodel.CreateVolumeParams) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Creating volumes with params: %+v", params)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return "", err
	}

	volumeIDsStr, err := provObj.CreateVolume(params)
	if err != nil {
		log.WriteError("volume creation failed: %v", err)
		return "", fmt.Errorf("volume creation failed: %w", err)
	}

	if volumeIDsStr == "" {
		return "", fmt.Errorf("no volumes were created")
	}

	log.WriteInfo("volumes created successfully: %v", volumeIDsStr)
	return volumeIDsStr, nil

}

func (psm *adminStorageManager) updateVolume(volumeID int, params gwymodel.CreateVolumeParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting update for Volume ID %d", volumeID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Retrieve current volume info ---
	tfVol, err := provObj.GetVolumeByID(volumeID)
	if err != nil {
		log.WriteError("Failed to retrieve volume info for ID %d: %v", volumeID, err)
		return fmt.Errorf("failed to retrieve volume info for ID %d: %w", volumeID, err)
	}

	// Pretty-print parameters for debugging
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Update Params: %v", string(b))

	// --- Step 1: Capacity Expansion ---
	if err := psm.expandVolume(params, tfVol); err != nil {
		log.WriteError("Failed to expand Volume ID %d: %v", volumeID, err)
		return err
	}

	// --- Step 2: Compression Settings ---
	if err := psm.updateVolumeReductionSettings(params, tfVol); err != nil {
		log.WriteError("Failed to update reduction settings for Volume ID %d: %v", volumeID, err)
		return err
	}

	// --- Step 3: Nickname update ---
	if err := psm.updateVolumeNicknameIfNeeded(params, tfVol); err != nil {
		log.WriteError("Failed to update nickname for Volume ID %d: %v", volumeID, err)
		return err
	}

	log.WriteInfo("Successfully updated Volume ID %d", volumeID)
	return nil
}

func (psm *adminStorageManager) expandVolume(params gwymodel.CreateVolumeParams, tfVol *gwymodel.VolumeInfoByID) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	requestedCap := params.Capacity
	currentCap := tfVol.TotalCapacity
	increment := requestedCap - currentCap

	if increment <= 0 {
		if requestedCap < currentCap {
			errMsg := fmt.Sprintf("Decreasing capacity is not allowed: requested %d MiB < current %d MiB", requestedCap, currentCap)
			log.WriteError(errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		// Equal capacity (no change)
		log.WriteInfo("Volume ID %d already at requested capacity (%d MiB), skipping expansion.", tfVol.ID, currentCap)
		return nil
	}

	log.WriteInfo("Expanding Volume ID %d by %d MiB (from %d → %d).", tfVol.ID, increment, currentCap, requestedCap)

	expandParams := gwymodel.ExpandVolumeParams{
		Capacity: increment,
	}

	err = provObj.ExpandVolume(tfVol.ID, expandParams)
	if err != nil {
		log.WriteError("Failed to expand Volume ID %d: %v", tfVol.ID, err)
		return err
	}

	log.WriteInfo("Successfully expanded Volume ID %d by %d MiB.", tfVol.ID, increment)
	return nil
}

func (psm *adminStorageManager) updateVolumeReductionSettings(params gwymodel.CreateVolumeParams, tfVol *gwymodel.VolumeInfoByID) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Skip update entirely if SavingSetting is nil ---
	if params.SavingSetting == nil {
		log.WriteInfo("Skipping reduction settings update for Volume ID %d: SavingSetting is nil.", tfVol.ID)
		return nil
	}

	// Extract requested saving
	requestedSaving := *params.SavingSetting

	// Extract requested acceleration (optional)
	var requestedAccel *bool
	if params.CompressionAcceleration != nil {
		requestedAccel = params.CompressionAcceleration
	}

	// Current values
	currentSaving := tfVol.SavingSetting
	currentAccel := tfVol.CompressionAcceleration // may be nil

	// --- Skip if both values match or unchanged ---
	sameSaving := currentSaving == requestedSaving
	sameAccel := (requestedAccel == nil && currentAccel == nil) ||
		(requestedAccel != nil && currentAccel != nil && *requestedAccel == *currentAccel)

	if sameSaving && sameAccel {
		log.WriteInfo("Volume ID %d already has the requested reduction settings (Saving: %v, Accel: %v), skipping update.",
			tfVol.ID, requestedSaving,
			func() any {
				if currentAccel == nil {
					return "(nil)"
				}
				return *currentAccel
			}())
		return nil
	}

	// --- Proceed with update ---
	log.WriteInfo("Updating reduction settings for Volume ID %d (Saving: %v → %v, Accel: %v → %v).",
		tfVol.ID,
		currentSaving, requestedSaving,
		func() any {
			if currentAccel == nil {
				return "(nil)"
			}
			return *currentAccel
		}(),
		func() any {
			if requestedAccel == nil {
				return "(no change)"
			}
			return *requestedAccel
		}(),
	)

	// Build params for provisioner
	updateParams := gwymodel.UpdateVolumeReductionParams{
		SavingSetting: params.SavingSetting,
	}

	// Only include CompressionAcceleration if explicitly provided
	if requestedAccel != nil {
		updateParams.CompressionAcceleration = requestedAccel
	}

	if err := provObj.UpdateVolumeReductionSettings(tfVol.ID, updateParams); err != nil {
		log.WriteError("Failed to update reduction settings for Volume ID %d: %v", tfVol.ID, err)
		return err
	}

	log.WriteInfo("Successfully updated reduction settings for Volume ID %d.", tfVol.ID)
	return nil
}

func (psm *adminStorageManager) updateVolumeNicknameIfNeeded(params gwymodel.CreateVolumeParams, tfVol *gwymodel.VolumeInfoByID) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	newNickname := extractRequestedNickname(params)
	oldNickname := ""
	if tfVol.Nickname != nil {
		oldNickname = *tfVol.Nickname
	}

	// --- Skip if no nickname change requested ---
	if newNickname == "" {
		log.WriteInfo("No nickname update requested for Volume ID %d (NicknameParam.BaseName is empty).", tfVol.ID)
		return nil
	}

	if newNickname == oldNickname {
		log.WriteInfo("Volume ID %d already has the requested nickname '%s'; skipping update.", tfVol.ID, newNickname)
		return nil
	}

	// --- Proceed with nickname update ---
	log.WriteInfo("Updating nickname for Volume ID %d from '%s' to '%s'.", tfVol.ID, oldNickname, newNickname)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	updateParams := gwymodel.UpdateVolumeNicknameParams{
		Nickname: newNickname,
	}

	if err := provObj.UpdateVolumeNickname(tfVol.ID, updateParams); err != nil {
		log.WriteError("Failed to update nickname for Volume ID %d: %v", tfVol.ID, err)
		return fmt.Errorf("failed to update nickname: %w", err)
	}

	log.WriteInfo("Successfully updated nickname for Volume ID %d to '%s'.", tfVol.ID, newNickname)
	return nil
}

func (psm *adminStorageManager) getProvisionerManager() (provmanager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	setting := provmodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get provisioner manager: %v", err)
		return nil, fmt.Errorf("failed to get provisioner manager: %w", err)
	}

	return provObj, nil
}

func extractRequestedNickname(params gwymodel.CreateVolumeParams) string {
	baseName := params.NicknameParam.BaseName
	if baseName == "" {
		return ""
	}

	startNum := -1
	numDigits := 0
	if params.NicknameParam.StartNumber != nil {
		startNum = *params.NicknameParam.StartNumber
	}
	if params.NicknameParam.NumberOfDigits != nil {
		numDigits = *params.NicknameParam.NumberOfDigits
	}

	// No suffix case (StartNumber == -1)
	if startNum == -1 {
		return baseName
	}

	// With suffix
	if numDigits > 0 {
		return fmt.Sprintf("%s%0*d", baseName, numDigits, startNum)
	}

	// Default single-digit suffix
	return fmt.Sprintf("%s%d", baseName, startNum)
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Common 404 / not found patterns
	notFoundPhrases := []string{
		"404",                    // HTTP 404 Not Found
		"not found",              // generic message
		"does not exist",         // explicit existence message
		"no such object",         // alternate REST wording
		"resource not available", // alternative API wording
		"kart70006-e",            // Hitachi error code for "not exist"
		"unmanaged resource",     // unmanaged / deleted
	}

	for _, phrase := range notFoundPhrases {
		if strings.Contains(errStr, strings.ToLower(phrase)) {
			return true
		}
	}

	return false
}

// parseVolumeIDs converts a comma-separated string like "1,2,3" into a slice of ints: []int{1, 2, 3}.
func parseVolumeIDs(volumeIDsStr string) ([]int, error) {
	idParts := strings.Split(volumeIDsStr, ",")

	ids := make([]int, 0, len(idParts))
	for _, part := range idParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		idInt, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid volume ID '%s': %v", part, err)
		}
		ids = append(ids, idInt)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no volume ids parsed from string: '%s'", volumeIDsStr)
	}

	return ids, nil
}
