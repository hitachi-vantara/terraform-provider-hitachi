package terraform

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	recmanager "terraform-provider-hitachi/hitachi/storage/san/reconciler"
	recimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	recmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Snapshot Datasources -------------------

func DatasourceVspSnapshotRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	input := recmodel.SnapshotGetMultipleInput{}

	if v, ok := d.GetOk("snapshot_group_name"); ok {
		val := v.(string)
		input.SnapshotGroupName = &val
	}

	if v, ok := d.GetOkExists("mirror_unit_id"); ok {
		val := v.(int)
		input.MuNumber = &val
	}

	finalPvol, err := terrcommon.ExtractLdevFields(d, "pvol_ldev_id", "pvol_ldev_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalPvol != nil {
		input.PvolLdevID = finalPvol
	}

	finalSvol, err := terrcommon.ExtractLdevFields(d, "svol_ldev_id", "svol_ldev_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalSvol != nil {
		input.SvolLdevID = finalSvol
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshots, err := reconObj.ReconcileGetMultipleSnapshots(input)
	if err != nil {
		log.WriteDebug("failed to get snapshots: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get snapshots: %v", err))
	}

	if err := d.Set("snapshots", convertSnapshotsToSchema(snapshots.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshots: %w", err))
	}

	if err := d.Set("snapshot_count", len(snapshots.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_count: %w", err))
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func DatasourceVspSnapshotRangeRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	params := recmodel.SnapshotGetMultipleRangeInput{}

	finalStartPvol, err := terrcommon.ExtractLdevFields(d, "start_pvol_ldev_id", "start_pvol_ldev_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalStartPvol != nil {
		params.StartPvolLdevID = finalStartPvol
	}

	finalEndPvol, err := terrcommon.ExtractLdevFields(d, "end_pvol_ldev_id", "end_pvol_ldev_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalEndPvol != nil {
		params.EndPvolLdevID = finalEndPvol
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshots, err := reconObj.ReconcileGetMultipleSnapshotsRange(params)
	if err != nil {
		log.WriteDebug("failed to get snapshot range: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get snapshot range: %v", err))
	}

	if err := d.Set("snapshots", convertSnapshotsToSchema(snapshots.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshots: %w", err))
	}

	if err := d.Set("snapshot_count", len(snapshots.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_count: %w", err))
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// ------------------- Snapshot Resource -------------------

func ResourceVspSnapshotRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	input, err := getSnapshotReconcilerInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := reconObj.ReconcileReadExistingSnapshotVclone(input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Update the state with all returned metadata
	updateSnapshotResourceState(d, result)

	return nil
}

func ResourceVspSnapshotDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	input, err := getSnapshotReconcilerInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	input.Action = utils.Ptr("delete")

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = reconObj.ReconcileSnapshotVclone(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// ResourceVspSnapshotApply handles both Create and Update logic
func ResourceVspSnapshotApply(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	input, err := getSnapshotReconcilerInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := reconObj.ReconcileSnapshotVclone(input)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSnapshotResourceState(d, result)

	return nil
}

// ------------------- Helpers -------------------

func updateSnapshotResourceState(d *schema.ResourceData, result *recmodel.ReconcileSnapshotResult) {
	log := commonlog.GetLogger()

	if result == nil {
		return
	}

	// 1. Handle Snapshot list (ensure it's never nil)
	if result.Snapshot != nil {
		d.Set("snapshot", convertSnapshotsToSchema([]gwymodel.Snapshot{*result.Snapshot}))
		// Set ID from the snapshot
		d.SetId(result.Snapshot.SnapshotID)
	} else {
		// Use an empty slice to avoid the "of object" null error
		d.Set("snapshot", []map[string]interface{}{})
	}

	// 2. Handle vClone list (ensure it's never nil)
	if result.VcloneFamily != nil {
		d.Set("vclone", convertVcloneFamilyToSchema(result.VcloneFamily))

		// CRITICAL FIX: If the Snapshot is gone, we MUST provide a persistent ID.
		// If d.Id() is currently empty, use the S-VOL LDEV ID.
		if d.Id() == "" {
			d.SetId(strconv.Itoa(result.VcloneFamily.LdevID))
		}
	} else {
		d.Set("vclone", []map[string]interface{}{})
	}

	d.Set("additional_info", convertUniversalInfoToSchema(result.UniversalInfo))

	// 3. Safety Check: If after all logic the ID is STILL empty,
	// Terraform will throw "Root object absent". We must force an ID.
	if d.Id() == "" {
		log.WriteDebug("No state id")
		d.SetId("No snapshot found")
	}
}

func convertSnapshotsToSchema(snapshots []gwymodel.Snapshot) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	result := make([]map[string]interface{}, 0, len(snapshots))
	for _, s := range snapshots {
		m := map[string]interface{}{
			"pvol_ldev_id":                   s.PvolLdevID,
			"pvol_ldev_id_hex":               utils.IntToHexString(s.PvolLdevID),
			"mirror_unit_id":                 s.MuNumber,
			"svol_ldev_id":                   s.SvolLdevID,
			"svol_ldev_id_hex":               utils.IntToHexString(s.SvolLdevID),
			"snapshot_pool_id":               s.SnapshotPoolID,
			"snapshot_id":                    s.SnapshotID,
			"snapshot_group_name":            s.SnapshotGroupName,
			"primary_or_secondary":           s.PrimaryOrSecondary,
			"status":                         s.Status,
			"is_redirect_on_write":           s.IsRedirectOnWrite,
			"is_consistency_group":           s.IsConsistencyGroup,
			"is_written_in_svol":             s.IsWrittenInSvol,
			"is_clone":                       s.IsClone,
			"can_cascade":                    s.CanCascade,
			"snapshot_data_read_only":        s.SnapshotDataReadOnly,
			"pvol_processing_status":         s.PvolProcessingStatus,
			"svol_processing_status":         s.SvolProcessingStatus,
			"retention_period_hours":         s.RetentionPeriod,
			"is_virtual_clone_volume":        s.IsVirtualCloneVolume,
			"is_virtual_clone_parent_volume": s.IsVirtualCloneParentVolume,
		}

		// Handle pointer fields/optionals
		if s.ConcordanceRate != nil {
			m["concordance_rate"] = *s.ConcordanceRate
		}
		if s.ProgressRate != nil {
			m["progress_rate"] = *s.ProgressRate
		}
		if s.SplitTime != nil {
			m["split_time"] = *s.SplitTime
		}

		result = append(result, m)
	}

	return result
}

func convertSnapshotRangeListToSchema(snapshots []gwymodel.SnapshotAll) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if len(snapshots) == 0 {
		return nil
	}

	// Sort by SnapshotReplicationID (pvol,mu) for consistent UI display
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].SnapshotReplicationID < snapshots[j].SnapshotReplicationID
	})

	snapshotList := make([]map[string]interface{}, len(snapshots))
	for i, s := range snapshots {
		m := map[string]interface{}{
			"snapshot_replication_id": s.SnapshotReplicationID,
			"pvol_ldev_id":            s.PvolLdevID,
			"pvol_ldev_id_hex":        utils.IntToHexString(s.PvolLdevID),
			"mirror_unit_id":          s.MuNumber,
			"snapshot_group_name":     s.SnapshotGroupName,
			"snapshot_pool_id":        s.SnapshotPoolID,
			"svol_ldev_id":            s.SvolLdevID,
			"svol_ldev_id_hex":        utils.IntToHexString(s.SvolLdevID),
			"consistency_group_id":    s.ConsistencyGroupID,
			"status":                  s.Status,
			"is_redirect_on_write":    s.IsRedirectOnWrite,
			"is_clone":                s.IsClone,
			"can_cascade":             s.CanCascade,
		}

		// Optional/Pointer fields
		if s.ConcordanceRate != nil {
			m["concordance_rate"] = *s.ConcordanceRate
		}
		if s.SplitTime != nil {
			m["split_time"] = *s.SplitTime
		}

		snapshotList[i] = m
	}

	return snapshotList
}

func convertUniversalInfoToSchema(info *recmodel.SnapshotUniversalInfo) []map[string]interface{} {
	if info == nil {
		return []map[string]interface{}{}
	}
	m := map[string]interface{}{
		"serial":                 info.StorageSerial,
		"is_thin_image_advanced": info.IsThinImageAdvanced,
		"snapshot_pool_type":     info.SnapshotPoolType,
		"pvol_attributes":        info.PvolAttributes,
		"svol_attributes":        info.SvolAttributes,
	}
	return []map[string]interface{}{m}
}

func convertVcloneFamilyToSchema(v *gwymodel.SnapshotFamily) []map[string]interface{} {
	if v == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		// Mapping ParentLdevID to pvol and the object's LdevID to svol
		"pvol_ldev_id":     v.ParentLdevID,
		"svol_ldev_id":     v.LdevID,
		"pvol_ldev_id_hex": utils.IntToHexString(v.ParentLdevID),
		"svol_ldev_id_hex": utils.IntToHexString(v.LdevID),

		// Attribute Flags
		"is_virtual_clone_volume":        v.IsVirtualCloneVolume,
		"is_virtual_clone_parent_volume": v.IsVirtualCloneParentVolume,

		// Metadata
		"split_time": v.SplitTime,
		"pool_id":    v.PoolID,
	}

	return []map[string]interface{}{m}
}

func getSnapshotReconcilerInput(d *schema.ResourceData) (recmodel.SnapshotReconcilerInput, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pvol, mu, err := getSnapshotIdentifiers(d)
	if err != nil {
		return recmodel.SnapshotReconcilerInput{}, err
	}

	// Capture the current configuration value (the "new" value in the plan)
	_, newSvol := d.GetChange("svol_ldev_id")
	svolID := newSvol.(int)

	input := recmodel.SnapshotReconcilerInput{
		PvolLdevID:               pvol,
		MuNumber:                 mu,
		Action:                   terrcommon.GetStringPointer(d, "state"),
		SnapshotGroupName:        terrcommon.GetStringPointer(d, "snapshot_group_name"),
		SnapshotPoolID:           terrcommon.GetIntPointer(d, "snapshot_pool_id"),
		SvolLdevID:               &svolID,
		IsConsistencyGroup:       terrcommon.GetBoolPointer(d, "is_consistency_group"),
		AutoSplit:                terrcommon.GetBoolPointer(d, "auto_split"),
		IsClone:                  terrcommon.GetBoolPointer(d, "is_clone"),
		ClonesAutomation:         terrcommon.GetBoolPointer(d, "auto_clone"),
		CanCascade:               terrcommon.GetBoolPointer(d, "can_cascade"),
		IsDataReductionForceCopy: terrcommon.GetBoolPointer(d, "is_data_reduction_force_copy"),
		CopySpeed:                terrcommon.GetStringPointer(d, "copy_speed"),
		RetentionPeriod:          terrcommon.GetIntPointer(d, "retention_period_hours"),
		DefragOperation:          terrcommon.GetStringPointer(d, "defrag_operation"),
	}

	// --- State Fallback Logic ---

	// 1. Check Snapshot Group Name
	if input.SnapshotGroupName == nil || *input.SnapshotGroupName == "" {
		if val, ok := d.GetOk("snapshot_group_name"); ok {
			groupName := val.(string)
			input.SnapshotGroupName = &groupName
			log.WriteDebug("[INFO] SnapshotGroupName missing in input, using '%s' from state", groupName)
		}
	}

	// 2. Check Snapshot Pool ID
	if input.SnapshotPoolID == nil {
		if val, ok := d.GetOk("snapshot_pool_id"); ok {
			poolID := val.(int) // adjust type (int/int64) based on your model
			input.SnapshotPoolID = &poolID
			log.WriteDebug("[INFO] SnapshotPoolID missing in input, using '%d' from state", poolID)
		}
	}

	log.WriteDebug("Snapshot Input: %+v", input)
	return input, nil
}

// getSnapshotIdentifiers handles the logic for extracting and validating pvol/mu
// from either the resource ID (state) or the user configuration (input).
func getSnapshotIdentifiers(d *schema.ResourceData) (*int, *int, error) {
	resourceID := d.Id()

	// pvol := terrcommon.GetIntPointer(d, "pvol_ldev_id")
	var pvol *int

	mu := terrcommon.GetIntPointer(d, "mirror_unit_id")

	finalPvol, err := terrcommon.ExtractLdevFields(d, "pvol_ldev_id", "pvol_ldev_id_hex")
	if err != nil {
		return nil, nil, err
	}
	if finalPvol != nil {
		pvol = finalPvol
	} else {
		return nil, nil, fmt.Errorf("pvol_ldev_id is required")
	}

	if resourceID != "" {
		// Existing resource: Parse "pvol,mu" from State ID
		parts := strings.Split(resourceID, ",")
		if len(parts) == 2 {
			statePvol, _ := strconv.Atoi(parts[0])
			stateMu, _ := strconv.Atoi(parts[1])

			// 1. Validate P-VOL Consistency
			if pvol != nil && *pvol != statePvol {
				return nil, nil, fmt.Errorf("pvol_ldev_id conflict: config has %d but state has %d. "+
					"A snapshot resource cannot be moved to a different P-VOL. Please recreate the resource", *pvol, statePvol)
			}
			// If input was nil, fallback to state value
			if pvol == nil {
				pvol = &statePvol
			}

			// 2. Validate MU Consistency
			if mu != nil && *mu != stateMu {
				return nil, nil, fmt.Errorf("mirror_unit_id conflict: config has %d but state has %d. "+
					"Snapshot Mirror Units cannot be changed after creation", *mu, stateMu)
			}
			// If input was nil, fallback to state value
			if mu == nil {
				mu = &stateMu
			}
		}
	}

	return pvol, mu, nil
}

func getReconcilerManagerSan(serial int) (recmanager.SanStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := recmodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	recObj, err := recimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get reconciler manager: %v", err)
		return nil, fmt.Errorf("failed to get reconciler manager: %w", err)
	}

	return recObj, nil
}
