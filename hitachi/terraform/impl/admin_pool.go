package terraform

import (
	"context"
	"encoding/json"
	"strconv"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	recmanager "terraform-provider-hitachi/hitachi/storage/admin/reconciler"
	recimpl "terraform-provider-hitachi/hitachi/storage/admin/reconciler/impl"
	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Resource CRUD Operations -------------------

func ResourceAdminPoolCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Build create parameters
	params, err := buildCreateAdminPoolParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	poolID, err := recObj.ReconcileCreateAdminPool(params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the resource ID to the pool ID
	d.SetId(strconv.Itoa(poolID))

	log.WriteInfo("VSP One pool created successfully with ID: %d", poolID)
	return ResourceAdminPoolRead(d)
}

func ResourceAdminPoolRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid pool ID format: %s", d.Id())
	}

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	pool, timedOut, err := recObj.ReconcileReadAdminPool(poolID)
	if err != nil {
		// If pool not found, remove from state
		if IsNotFoundError(err) {
			log.WriteWarn("VSP One pool %d not found, removing from state", poolID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Set attributes from the retrieved pool
	if err := setAdminPoolResourceData(d, pool); err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("VSP One pool read successfully: %d", poolID)

	// Add diagnostic warning if polling timed out
	if timedOut {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Pool capacity consistency polling timed out",
				Detail:   "The pool expansion operation is taking longer than expected and has not yet completed. Please continue running the Storage Pool data source to monitor the final pool capacity size. While the capacity is not yet finalized, the storage pool remains available for provisioning.",
			},
		}
	}

	return nil
}

func ResourceAdminPoolUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid pool ID format: %s", d.Id())
	}

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// If there are no relevant changes, short-circuit to a no-op
	// This prevents calling backend update with empty params and avoids 400 errors
	if !(d.HasChange("threshold_warning") || d.HasChange("threshold_depletion") || d.HasChange("drive_configuration") || d.HasChange("name")) {
		log.WriteInfo("No changes detected for VSP One pool; skipping update")
		return ResourceAdminPoolRead(d)
	}

	// Check if drive configuration was changed for expansion
	if d.HasChange("drive_configuration") {
		expandParams, err := buildExpandAdminPoolParamsFromDriveConfig(d)
		if err != nil {
			return diag.FromErr(err)
		}

		if expandParams != nil {
			err = recObj.ReconcileExpandAdminPool(poolID, *expandParams)
			if err != nil {
				return diag.FromErr(err)
			}
			log.WriteInfo("VSP One pool expanded successfully: %d", poolID)
		}
	}

	// Build update parameters for other changes
	params, err := buildUpdateAdminPoolParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// If there are other changes, apply them
	if params != nil {
		err = recObj.ReconcileUpdateAdminPool(poolID, *params)
		if err != nil {
			return diag.FromErr(err)
		}
		log.WriteInfo("VSP One pool updated successfully: %d", poolID)
	}

	return ResourceAdminPoolRead(d)
}

func ResourceAdminPoolDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid pool ID format: %s", d.Id())
	}

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	err = recObj.ReconcileDeleteAdminPool(poolID)
	if err != nil {
		// If pool not found, consider it already deleted
		if IsNotFoundError(err) {
			log.WriteWarn("VSP One pool %d not found, considering deleted", poolID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("VSP One pool deleted successfully: %d", poolID)
	return nil
}

// ------------------- Helper Functions -------------------

func getReconcilerManagerForAdminPool(serial int) (recmanager.AdminStorageManager, error) {
	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := recmodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	return recimpl.NewEx(setting)
}

func buildCreateAdminPoolParams(d *schema.ResourceData) (gwymodel.CreateAdminPoolParams, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	params := gwymodel.CreateAdminPoolParams{
		Name:                d.Get("name").(string),
		IsEncryptionEnabled: d.Get("encryption").(bool),
	}

	// Build drives array from schema
	drivesInput := d.Get("drive_configuration").([]interface{})
	var drives []gwymodel.CreateAdminPoolDrive

	for _, driveInterface := range drivesInput {
		driveData := driveInterface.(map[string]interface{})
		drive := gwymodel.CreateAdminPoolDrive{
			DriveTypeCode:   driveData["drive_type_code"].(string),
			DataDriveCount:  driveData["data_drive_count"].(int),
			RaidLevel:       driveData["raid_level"].(string),
			ParityGroupType: driveData["parity_group_type"].(string),
		}
		drives = append(drives, drive)
	}

	params.Drives = drives

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Create Pool Params: %v", string(b))

	return params, nil
}

func buildUpdateAdminPoolParams(d *schema.ResourceData) (*gwymodel.UpdateAdminPoolParams, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	hasChanges := false
	params := gwymodel.UpdateAdminPoolParams{}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		hasChanges = true
	}

	if d.HasChange("threshold_warning") {
		params.ThresholdWarning = d.Get("threshold_warning").(int)
		hasChanges = true
	}

	if d.HasChange("threshold_depletion") {
		params.ThresholdDepletion = d.Get("threshold_depletion").(int)
		hasChanges = true
	}

	if !hasChanges {
		return nil, nil
	}

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Update Pool Params: %v", string(b))

	return &params, nil
}

func buildExpandAdminPoolParamsFromDriveConfig(d *schema.ResourceData) (*gwymodel.ExpandAdminPoolParams, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Get the old and new drive configurations
	oldDrivesRaw, newDrivesRaw := d.GetChange("drive_configuration")

	oldDrives := oldDrivesRaw.([]interface{})
	newDrives := newDrivesRaw.([]interface{})

	// If no new drives were added, no expansion needed
	if len(newDrives) <= len(oldDrives) {
		log.WriteInfo("No new drives added; no expansion needed")
		return nil, nil
	}

	// Extract only the newly added drives (drives beyond the original count)
	var additionalDrives []gwymodel.ExpandAdminPoolDrive
	for i := len(oldDrives); i < len(newDrives); i++ {
		driveMap := newDrives[i].(map[string]interface{})
		drive := gwymodel.ExpandAdminPoolDrive{
			DriveTypeCode:   driveMap["drive_type_code"].(string),
			DataDriveCount:  driveMap["data_drive_count"].(int),
			RaidLevel:       driveMap["raid_level"].(string),
			ParityGroupType: driveMap["parity_group_type"].(string),
		}
		additionalDrives = append(additionalDrives, drive)
	}

	if len(additionalDrives) == 0 {
		return nil, nil
	}

	params := gwymodel.ExpandAdminPoolParams{
		AdditionalDrives: additionalDrives,
	}

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Expand Pool Params from Drive Config: %v", string(b))

	return &params, nil
}

func setAdminPoolResourceData(d *schema.ResourceData, pool *gwymodel.AdminPool) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Convert the pool data to the format expected by AdminPoolInfoSchema
	poolData := map[string]interface{}{
		"pool_id":                         pool.ID,
		"name":                            pool.Name,
		"status":                          pool.Status,
		"encryption_status":               pool.Encryption,
		"total_capacity":                  pool.TotalCapacity,
		"effective_capacity":              pool.EffectiveCapacity,
		"used_capacity":                   pool.UsedCapacity,
		"free_capacity":                   pool.FreeCapacity,
		"number_of_volumes":               pool.NumberOfVolumes,
		"number_of_tiers":                 pool.NumberOfTiers,
		"number_of_drive_types":           pool.NumberOfDriveTypes,
		"contains_capacity_saving_volume": pool.ContainsCapacitySavingVolume,
		"config_status":                   pool.ConfigStatus,
	}

	// Set drives information
	if len(pool.Drives) > 0 {
		var drives []map[string]interface{}
		for _, drive := range pool.Drives {
			driveMap := map[string]interface{}{
				"drive_type":             drive.DriveType,
				"drive_interface":        drive.DriveInterface,
				"drive_rpm":              drive.DriveRpm,
				"drive_capacity":         drive.DriveCapacity,
				"display_drive_capacity": drive.DisplayDriveCapacity,
				"total_capacity":         drive.TotalCapacity,
				"number_of_drives":       drive.NumberOfDrives,
				"raid_level":             drive.RaidLevel,
				"parity_group_type":      drive.ParityGroupType,
			}

			if len(drive.Locations) > 0 {
				driveMap["locations"] = drive.Locations
			}

			drives = append(drives, driveMap)
		}
		poolData["drives"] = drives
	}

	// Set capacity management
	capacityManage := map[string]interface{}{
		"used_capacity_rate":  pool.CapacityManage.UsedCapacityRate,
		"threshold_warning":   pool.CapacityManage.ThresholdWarning,
		"threshold_depletion": pool.CapacityManage.ThresholdDepletion,
	}
	poolData["capacity_manage"] = []interface{}{capacityManage}

	// Set saving effects
	savingEffects := map[string]interface{}{
		"efficiency_data_reduction":                  pool.SavingEffects.EfficiencyDataReduction,
		"efficiency_fmd_saving":                      pool.SavingEffects.EfficiencyFmdSaving,
		"pre_capacity_fmd_saving":                    pool.SavingEffects.PreCapacityFmdSaving,
		"post_capacity_fmd_saving":                   pool.SavingEffects.PostCapacityFmdSaving,
		"is_total_efficiency_support":                pool.SavingEffects.IsTotalEfficiencySupport,
		"total_efficiency_status":                    pool.SavingEffects.TotalEfficiencyStatus,
		"data_reduction_without_system_data_status":  pool.SavingEffects.DataReductionWithoutSystemDataStatus,
		"software_saving_without_system_data_status": pool.SavingEffects.SoftwareSavingWithoutSystemDataStatus,
		"data_reduction_without_system_data":         pool.SavingEffects.DataReductionWithoutSystemData,
		"software_saving_without_system_data":        pool.SavingEffects.SoftwareSavingWithoutSystemData,
	}

	if pool.SavingEffects.CalculationStartTime != nil {
		savingEffects["calculation_start_time"] = pool.SavingEffects.CalculationStartTime.String()
	}
	if pool.SavingEffects.CalculationEndTime != nil {
		savingEffects["calculation_end_time"] = pool.SavingEffects.CalculationEndTime.String()
	}

	poolData["saving_effects"] = []interface{}{savingEffects}

	// Set the data attribute as a list containing the pool data
	if err := d.Set("data", []interface{}{poolData}); err != nil {
		return err
	}

	log.WriteDebug("Pool data set successfully in data attribute")
	return nil
}

// ------------------- Data Source Operations -------------------

func DataSourceAdminPoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Build list parameters
	params := gwymodel.AdminPoolListParams{}

	if nameFilter, ok := d.GetOk("name_filter"); ok {
		name := nameFilter.(string)
		params.Name = &name
	}

	if statusFilter, ok := d.GetOk("status_filter"); ok {
		status := statusFilter.(string)
		params.Status = &status
	}

	if configStatusFilter, ok := d.GetOk("config_status_filter"); ok {
		configStatus := configStatusFilter.(string)
		params.ConfigStatus = &configStatus
	}

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	poolListResponse, err := recObj.GetAdminPoolList(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if poolListResponse == nil || len(poolListResponse.Data) == 0 {
		log.WriteInfo("No VSP One pools found")
		d.SetId("")
		if err := d.Set("data", []interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("pool_counts", 0); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	// Flatten the pool list response
	pools := flattenAdminPoolListResponse(poolListResponse)
	if err := d.Set("data", pools); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pool_counts", len(poolListResponse.Data)); err != nil {
		return diag.FromErr(err)
	}

	// Set a meaningful ID for the data source
	d.SetId(strconv.Itoa(serial) + "_pools")

	log.WriteInfo("VSP One pools list retrieved successfully: %d pools", len(poolListResponse.Data))
	return nil
}

func DataSourceAdminPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolID := d.Get("pool_id").(int)

	recObj, err := getReconcilerManagerForAdminPool(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	pool, err := recObj.GetAdminPoolInfo(poolID)
	if err != nil {
		return diag.FromErr(err)
	}

	if pool == nil {
		return diag.Errorf("VSP One pool %d not found", poolID)
	}

	// Set attributes from the retrieved pool
	if err := setAdminPoolResourceData(d, pool); err != nil {
		return diag.FromErr(err)
	}

	// Set a meaningful ID for the data source
	d.SetId(strconv.Itoa(serial) + "_pool_" + strconv.Itoa(poolID))

	log.WriteInfo("VSP One pool info retrieved successfully: %d", poolID)
	return nil
}

// flattenAdminPoolListResponse flattens the admin pool list response for the Terraform schema
func flattenAdminPoolListResponse(resp *gwymodel.AdminPoolListResponse) []map[string]interface{} {
	if resp == nil || len(resp.Data) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(resp.Data))
	for i, pool := range resp.Data {
		poolMap := map[string]interface{}{
			"pool_id":                         pool.ID,
			"name":                            pool.Name,
			"status":                          pool.Status,
			"encryption_status":               pool.Encryption,
			"total_capacity":                  pool.TotalCapacity,
			"effective_capacity":              pool.EffectiveCapacity,
			"used_capacity":                   pool.UsedCapacity,
			"free_capacity":                   pool.FreeCapacity,
			"number_of_volumes":               pool.NumberOfVolumes,
			"number_of_tiers":                 pool.NumberOfTiers,
			"number_of_drive_types":           pool.NumberOfDriveTypes,
			"contains_capacity_saving_volume": pool.ContainsCapacitySavingVolume,
		}

		// Add capacity management
		capacityManage := map[string]interface{}{
			"used_capacity_rate":  pool.CapacityManage.UsedCapacityRate,
			"threshold_warning":   pool.CapacityManage.ThresholdWarning,
			"threshold_depletion": pool.CapacityManage.ThresholdDepletion,
		}
		poolMap["capacity_manage"] = []interface{}{capacityManage}

		// Add saving effects
		savingEffects := map[string]interface{}{
			"efficiency_data_reduction":                  pool.SavingEffects.EfficiencyDataReduction,
			"efficiency_fmd_saving":                      pool.SavingEffects.EfficiencyFmdSaving,
			"pre_capacity_fmd_saving":                    pool.SavingEffects.PreCapacityFmdSaving,
			"post_capacity_fmd_saving":                   pool.SavingEffects.PostCapacityFmdSaving,
			"is_total_efficiency_support":                pool.SavingEffects.IsTotalEfficiencySupport,
			"total_efficiency_status":                    pool.SavingEffects.TotalEfficiencyStatus,
			"data_reduction_without_system_data_status":  pool.SavingEffects.DataReductionWithoutSystemDataStatus,
			"software_saving_without_system_data_status": pool.SavingEffects.SoftwareSavingWithoutSystemDataStatus,
			"data_reduction_without_system_data":         pool.SavingEffects.DataReductionWithoutSystemData,
			"software_saving_without_system_data":        pool.SavingEffects.SoftwareSavingWithoutSystemData,
		}

		if pool.SavingEffects.CalculationStartTime != nil {
			savingEffects["calculation_start_time"] = pool.SavingEffects.CalculationStartTime.String()
		}
		if pool.SavingEffects.CalculationEndTime != nil {
			savingEffects["calculation_end_time"] = pool.SavingEffects.CalculationEndTime.String()
		}

		poolMap["saving_effects"] = []interface{}{savingEffects}

		// Add config status
		poolMap["config_status"] = pool.ConfigStatus

		// Add drives information
		if len(pool.Drives) > 0 {
			var drives []map[string]interface{}
			for _, drive := range pool.Drives {
				driveMap := map[string]interface{}{
					"drive_type":             drive.DriveType,
					"drive_interface":        drive.DriveInterface,
					"drive_rpm":              drive.DriveRpm,
					"drive_capacity":         drive.DriveCapacity,
					"display_drive_capacity": drive.DisplayDriveCapacity,
					"total_capacity":         drive.TotalCapacity,
					"number_of_drives":       drive.NumberOfDrives,
					"raid_level":             drive.RaidLevel,
					"parity_group_type":      drive.ParityGroupType,
				}

				if len(drive.Locations) > 0 {
					driveMap["locations"] = drive.Locations
				}

				drives = append(drives, driveMap)
			}
			poolMap["drives"] = drives
		}

		result[i] = poolMap
	}
	return result
}

// IsNotFoundError checks if the error is a "not found" error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// Check for common "not found" error patterns
	errStr := err.Error()
	return contains(errStr, "not found") ||
		contains(errStr, "404") ||
		contains(errStr, "does not exist")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
