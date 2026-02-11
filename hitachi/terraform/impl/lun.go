package terraform

import (
	"fmt"
	"strconv"
	"strings"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetLun gets a lun
func GetLun(d *schema.ResourceData) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	var lunID int

	finalLdev, err := terrcommon.ExtractLdevFields(d, "ldev_id", "ldev_id_hex")
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), lunID)
		return nil, err
	}
	if finalLdev != nil {
		lunID = *finalLdev
	} else {
		lunFromState := d.State().ID
		log.WriteDebug("TFDebug| lunFromState from state: %s", lunFromState)
		if lunFromState != "" {
			lun, err := strconv.Atoi(lunFromState)
			if err != nil {
				log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
				return nil, err
			}
			lunID = lun
		} else {
			return nil, fmt.Errorf("state ID for ldev_id is missing")
		}
	}

	log.WriteDebug("TFDebug| lunID: %v", lunID)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_BEGIN), lunID)
	lun, err := reconObj.GetLun(lunID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), lunID)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_END), lunID)

	return lun, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func GetRangeOfLuns(d *schema.ResourceData) (*[]sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	startLdevID := 0
	finalStartLdev, err := terrcommon.ExtractLdevFields(d, "start_ldev_id", "start_ldev_id_hex")
	if err != nil {
		return nil, err
	}
	if finalStartLdev != nil {
		startLdevID = *finalStartLdev
	}

	endLdevID := 0
	finalEndLdev, err := terrcommon.ExtractLdevFields(d, "end_ldev_id", "end_ldev_id_hex")
	if err != nil {
		return nil, err
	}
	if finalEndLdev != nil {
		endLdevID = *finalEndLdev
	}

	if startLdevID < 0 {
		return nil, fmt.Errorf("start_ldev_id/start_ldev_id_hex must be greater than or equal to 0")
	}

	if endLdevID < 0 {
		return nil, fmt.Errorf("end_ldev_id/end_ldev_id_hex must be greater than or equal to 0")
	}

	if endLdevID < startLdevID {
		return nil, fmt.Errorf("end_ldev_id/end_ldev_id_hex must be greater than or equal to start_ldev_id/start_ldev_id_hex")
	}

	isUndefindLdev := d.Get("undefined_ldev").(bool)

	// Get new parameters
	filterOption := ""
	if v, ok := d.GetOk("filter_option"); ok {
		f, err := normalizeFilterOption(v.(string))
		if err != nil {
			return nil, err
		}
		filterOption = f
	}

	// Extract include_detail_info and include_cache_info parameters
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, add all detail types except class
		detailInfoTypes = append(detailInfoTypes, "FMC", "externalVolume", "virtualSerialNumber", "savingInfo", "qos", "nguId")
		log.WriteDebug("TFDebug| GetLuns: include_detail_info=true, added detail types")
	}
	if includeCacheInfo, ok := d.GetOk("include_cache_info"); ok && includeCacheInfo.(bool) {
		// When includeCacheInfo is true, add class
		detailInfoTypes = append(detailInfoTypes, "class")
		log.WriteDebug("TFDebug| GetLuns: include_cache_info=true, added class")
	}

	detailInfoType := ""
	if len(detailInfoTypes) > 0 {
		detailInfoType = strings.Join(detailInfoTypes, ",")
		log.WriteDebug("TFDebug| GetLuns: final detailInfoType: %s", detailInfoType)
	}

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_BEGIN), startLdevID, endLdevID)

	if strings.EqualFold(filterOption, "mainframe") {
		rawLuns, err := reconObj.GetRangeOfLuns(startLdevID, endLdevID, isUndefindLdev, "", detailInfoType)
		if err != nil {
			log.WriteDebug("TFError| error fetching luns for mainframe filter, err: %v", err)
			return nil, err
		}
		filtered := make([]sangatewaymodel.LogicalUnit, 0)
		for _, lu := range *rawLuns {
			et := lu.EmulationType
			if strings.HasPrefix(et, "3390-A") || strings.HasPrefix(et, "3390-V") {
				// Checked as prefix as MF volumes would have EmulationType as "3390-V-CVS"
				filtered = append(filtered, lu)
			}
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)
		return &filtered, nil
	}

	luns, err := reconObj.GetRangeOfLuns(startLdevID, endLdevID, isUndefindLdev, filterOption, detailInfoType)
	if err != nil {
		log.WriteDebug("TFError| error in GetRangeOfLuns terraform call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)

	return luns, nil
}

// CreateLun creates a lun
func CreateLun(d *schema.ResourceData) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	if IsReadExistingMode(d) {
		lun, err := GetLun(d)
		if err != nil {
			return nil, err
		}
		if lun == nil {
			return nil, fmt.Errorf("volume not found on storage %d", serial)
		}
		if lun.EmulationType == "NOT DEFINED" || lun.ByteFormatCapacity == "" {
			return nil, fmt.Errorf("volume %v does not exist on storage %d", lun.LdevID, serial)
		}
		log.WriteDebug("IsReadExistingMode: Returning from reading Lun in create.")
		return lun, nil
	}

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	// Decide request type before entering reconciler flow.
	req, err := CreateLunRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	lun, err := reconObj.SetLun(req)
	if err != nil {
		log.WriteDebug("TFError| error in SetLun, err: %v", err)
		return nil, err
	}

	// Decide whether a post-create label update is required.
	// For mainframe volumes created on a parity group, some arrays require
	// creating the LDEV without a label and then setting the label afterward.
	_, isMainframe := d.GetOk("cylinder")
	pool_id := d.Get("pool_id").(int)
	pool_name := d.Get("pool_name").(string)
	paritygroup_id := d.Get("paritygroup_id").(string)
	external_paritygroup_id := d.Get("external_paritygroup_id").(string)

	if isMainframe && paritygroup_id != "" && pool_id < 0 && pool_name == "" && external_paritygroup_id == "" {
		if nameVal, ok := d.GetOk("name"); ok && nameVal.(string) != "" {
			// Ensure resource ID is set so update helpers can derive LDEV
			d.SetId(strconv.Itoa(*lun))
			label := nameVal.(string)
			updateReq := reconcilermodel.UpdateLunRequest{
				LdevID: lun,
				Name:   &label,
			}
			_, err := reconObj.UpdateLun(&updateReq)
			if err != nil {
				log.WriteDebug("TFError| error updating label for parity-group mainframe volume, err: %v", err)
				return nil, err
			}
		}
	}

	// For parity-group mainframe volumes, storage may require formatting after create.
	// Option A: invoke NORMAL format (mapped to API FMT) immediately post-create.
	if isMainframe && paritygroup_id != "" && pool_id < 0 && pool_name == "" && external_paritygroup_id == "" {
		// ensure we have the ldev id
		if lun != nil {
			opType := "QFMT"
			// Determine if data reduction is enabled on the newly created LUN
			isForce := false
			currentLun, err := reconObj.GetLun(*lun)
			if err != nil {
				log.WriteDebug("TFError| error fetching LUN to determine data reduction mode: %v", err)
				return nil, err
			}
			if currentLun.DataReductionMode != "" && currentLun.DataReductionMode != "disabled" {
				isForce = true
			}

			fmtReq := reconcilermodel.FormatLdevRequest{
				OperationType:              &opType,
				IsDataReductionForceFormat: &isForce,
			}

			_, err = reconObj.FormatLdev(*lun, fmtReq)
			if err != nil {
				log.WriteDebug("TFError| error invoking FormatLdev after create for parity-group mainframe LDEV: %v", err)
				// Attempt to restore LDEV to normal state if the format failed.
				if unErr := reconObj.UnblockLun(*lun); unErr != nil {
					log.WriteDebug("TFError| error attempting UnblockLun after failed format: %v", unErr)
				}
				return nil, err
			}
		}
	}

	// Single fetch of the LUN after creation (and optional update)
	logicalUnit, err := reconObj.GetLun(*lun)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun, err: %v", err)
		return nil, err
	}

	return logicalUnit, nil
}

func CreateLunRequestFromSchema(d *schema.ResourceData) (*reconcilermodel.LunRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := reconcilermodel.LunRequest{}

	// Determine mainframe intent.
	// Mainframe vs block is driven ONLY by cylinder presence.
	_, isMainframe := d.GetOk("cylinder")

	// size_gb → ByteFormatCapacity
	if size_gb, ok := d.GetOk("size_gb"); ok {
		createInput.ByteFormatCapacity = utils.ConvertFloatSizeToSmartUnit(size_gb.(float64))
	}

	// Mainframe-specific fields
	if isMainframe {
		if v, ok := d.GetOk("emulation_type"); ok {
			emulationType := v.(string)
			createInput.EmulationType = &emulationType
		} else {
			// If not specified, default to OPEN-V for mainframe volumes.
			// (Block volumes do not support emulation_type, and mainframe intent is driven by cylinder.)
			defaultEmulationType := "OPEN-V"
			createInput.EmulationType = &defaultEmulationType
		}
		if v, ok := d.GetOk("cylinder"); ok {
			c := v.(int)
			createInput.Cylinder = &c
		}
		if v, ok := d.GetOk("ssid"); ok {
			ssid := v.(string)
			createInput.Ssid = &ssid
		}
		if v, ok := d.GetOk("mp_blade_id"); ok {
			mp := v.(int)
			createInput.MpBladeID = &mp
		}
		if v, ok := d.GetOk("clpr_id"); ok {
			clpr := v.(int)
			createInput.ClprID = &clpr
		}
		if v, ok := d.GetOk("is_tse_volume"); ok {
			b := v.(bool)
			createInput.IsTseVolume = &b
		}
		if v, ok := d.GetOk("is_ese_volume"); ok {
			b := v.(bool)
			createInput.IsEseVolume = &b
		}
	}

	finalLdev, err := terrcommon.ExtractLdevFields(d, "ldev_id", "ldev_id_hex")
	if err != nil {
		return nil, err
	}
	if finalLdev != nil {
		if isMainframe {
			return nil, fmt.Errorf("ldev_id/ldev_id_hex is not supported for mainframe volume creation")
		}
		createInput.LdevID = finalLdev
	}

	// placement fields (read early so name handling can inspect them)
	pool_id := d.Get("pool_id").(int)
	pool_name := d.Get("pool_name").(string)
	paritygroup_id := d.Get("paritygroup_id").(string)
	external_paritygroup_id := d.Get("external_paritygroup_id").(string)

	name, ok := d.GetOk("name")
	if ok {
		// For mainframe volumes created on a parity group, skip sending the
		// label on create (the storage requires a follow-up update to set it).
		skipLabel := isMainframe && paritygroup_id != "" && pool_id < 0 && pool_name == "" && external_paritygroup_id == "" && name.(string) != ""
		if !skipLabel {
			label := name.(string)
			createInput.Name = &label
		}
	}

	// Data reduction / compression acceleration are block-only settings.
	// Mainframe (cylinder-driven) creates must not send these fields at all.
	if !isMainframe {
		dedup_mode, ok := d.GetOk("capacity_saving")
		if ok {
			dedup := dedup_mode.(string)
			createInput.DataReductionMode = &dedup
		}

		if v, ok := d.GetOk("is_data_reduction_shared_volume_enabled"); ok {
			val := v.(bool)
			createInput.IsDataReductionSharedVolumeEnabled = &val
		}

		if v, ok := d.GetOk("is_compression_acceleration_enabled"); ok {
			val := v.(bool)
			createInput.IsCompressionAccelerationEnabled = &val
		}
	}

	log.WriteDebug("Pool ID=%v Pool Name=%v PG=%v ExPG=%v\n", pool_id, pool_name, paritygroup_id, external_paritygroup_id)

	if pool_id >= -1 {
		// pool_id_int := pool_id.(int)
		createInput.PoolID = &pool_id
	} else if pool_name != "" {
		ppid, err := GetPoolIdFromPoolName(d, pool_name)
		createInput.PoolID = ppid
		if err != nil {
			return nil, fmt.Errorf("could not find a pool with name %v", pool_name)
		}
	} else if paritygroup_id != "" {
		createInput.ParityGroupID = &paritygroup_id
	} else if external_paritygroup_id != "" {
		createInput.ExternalParityGroupID = &external_paritygroup_id
	}

	// For parity-group mainframe (3390-V), send capacity in blocks instead of cylinder.
	if isMainframe && createInput.Cylinder != nil && createInput.ParityGroupID != nil {
		blocks := calculate3390LdevBlocks(int64(*createInput.Cylinder))
		createInput.BlockCapacity = &blocks
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func calculate3390LdevBlocks(requestedCylinders int64) int64 {
	const blocksPerCyl int64 = 1740
	const boundaryBlocks int64 = 77952
	if requestedCylinders <= 0 {
		return 0
	}

	rawBlocks := requestedCylinders * blocksPerCyl
	boundaryUnits := (rawBlocks + boundaryBlocks - 1) / boundaryBlocks
	return boundaryUnits * boundaryBlocks
}

func GetPoolIdFromPoolName(d *schema.ResourceData, poolName string) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	//log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevID, setting.Serial)

	pools, err := reconObj.GetPools()
	if err != nil {
		//log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevID, setting.Serial)
		return nil, err
	}
	//log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevID, setting.Serial)

	poolId := -999

	for _, pool := range *pools {
		if strings.EqualFold(pool.PoolName, poolName) {
			poolId = pool.PoolID
			break

		}
	}

	return &poolId, nil
}

func normalizeFilterOption(s string) (string, error) {
	m := map[string]string{
		"defined":         "defined",
		"undefined":       "undefined",
		"dpvolume":        "dpVolume",
		"lumapped":        "luMapped",
		"luunmapped":      "luUnmapped",
		"externalvolume":  "externalVolume",
		"mappednamespace": "mappedNamespace",
		"mainframe":       "mainframe",
	}
	key := strings.ToLower(strings.TrimSpace(s))
	if v, ok := m[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("invalid filter_option: %s", s)
}

// DeleteLun deletes a lun
func DeleteLun(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	ldevID, ok := d.GetOk("ldev_id")
	log.WriteDebug("ldevID: %+v", ldevID)
	lunID := 0
	isLdevIdSetFromState := false
	if !ok {
		lunFromState := d.State().ID
		if lunFromState != "" {
			lun, err := strconv.Atoi(lunFromState)
			if err != nil {
				log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
				return err
			}
			lunID = lun
			isLdevIdSetFromState = true
		} else {
			volume, ok := d.GetOk("volume")
			if !ok {
				return fmt.Errorf("no volume data in resource")
			}
			log.WriteDebug("volume: %+v", volume.([]map[string]interface{})[0])
			ldevID, ok = volume.([]map[string]interface{})[0]["ldev_id"]
			if !ok {
				return fmt.Errorf("found no ldev_id in info")
			}
			log.WriteDebug("volume ldevID: %+v", ldevID)
		}
	}

	if !isLdevIdSetFromState {
		lunID = ldevID.(int)
	}

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), lunID, setting.Serial)

	err = reconObj.DeleteLun(lunID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), lunID, setting.Serial)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), lunID, setting.Serial)

	return nil
}

func ConvertLunToSchema(logicalUnit *sangatewaymodel.LogicalUnit, serial int) *map[string]interface{} {

	lun := map[string]interface{}{
		"storage_serial_number": serial,

		// --- BASIC IDENTIFIERS ---
		"ldev_id":               logicalUnit.LdevID,
		"ldev_id_hex":           utils.IntToHexString(logicalUnit.LdevID),
		"virtual_ldev_id":       logicalUnit.VirtualLdevID,
		"virtual_serial_number": logicalUnit.VirtualSerialNumber,
		"virtual_model":         logicalUnit.VirtualModel,
		"clpr_id":               logicalUnit.ClprID,
		"emulation_type":        logicalUnit.EmulationType,
		"byte_format_capacity":  logicalUnit.ByteFormatCapacity,
		"block_capacity":        logicalUnit.BlockCapacity,
		"cylinder":              logicalUnit.Cylinder,

		// --- ATTRIBUTES ---
		"attributes": logicalUnit.Attributes,
		"label":      logicalUnit.Label,
		"status":     logicalUnit.Status,

		"parent_ldev_id":     logicalUnit.ParentLdevId,
		"parent_ldev_id_hex": utils.IntToHexString(logicalUnit.ParentLdevId),

		// --- CORE INFO ---
		"mpblade_id":                 logicalUnit.MpBladeID,
		"ssid":                       logicalUnit.Ssid,
		"pool_id":                    logicalUnit.PoolID,
		"num_of_used_block":          logicalUnit.NumOfUsedBlock,
		"is_full_allocation_enabled": logicalUnit.IsFullAllocationEnabled,
		"resource_group_id":          logicalUnit.ResourceGroupID,

		// --- DATA REDUCTION ---
		"data_reduction_status":        logicalUnit.DataReductionStatus,
		"data_reduction_mode":          logicalUnit.DataReductionMode,
		"data_reduction_process_mode":  logicalUnit.DataReductionProcessMode,
		"data_reduction_progress_rate": logicalUnit.DataReductionProgressRate,

		// --- ALUA ---
		"is_alua_enabled": logicalUnit.IsAluaEnabled,

		// --- NAA ---
		"naa_id": logicalUnit.NaaID,

		// --- COMPRESSION ACCEL ---
		"is_compression_acceleration_enabled": logicalUnit.IsCompressionAccelerationEnabled,
		"compression_acceleration_status":     logicalUnit.CompressionAccelerationStatus,

		// --- RAID ---
		"raid_level":                 logicalUnit.RaidLevel,
		"raid_type":                  logicalUnit.RaidType,
		"num_of_parity_groups":       logicalUnit.NumOfParityGroups,
		"parity_group_ids":           logicalUnit.ParityGroupIds,
		"drive_type":                 logicalUnit.DriveType,
		"drive_byte_format_capacity": logicalUnit.DriveByteFormatCapacity,
		"drive_block_capacity":       logicalUnit.DriveBlockCapacity,

		// --- CAPACITIES ---
		"total_capacity_in_mb": logicalUnit.TotalCapacityInMB,
		"free_capacity_in_mb":  logicalUnit.FreeCapacityInMB,
		"used_capacity_in_mb":  logicalUnit.UsedCapacityInMB,

		// --- PORT COUNT ---
		"num_ports": logicalUnit.NumOfPorts,

		// --- NEW: COMPOSING / SNAPSHOT POOLS ---
		"composing_pool_id": logicalUnit.ComposingPoolId,
		"snapshot_pool_id":  logicalUnit.SnapshotPoolId,

		// --- NEW: EXTERNAL VOLUME FIELDS ---
		"external_vendor_id":        logicalUnit.ExternalVendorId,
		"external_product_id":       logicalUnit.ExternalProductId,
		"external_volume_id":        logicalUnit.ExternalVolumeId,
		"external_volume_id_string": logicalUnit.ExternalVolumeIdString,

		"num_of_external_ports": logicalUnit.NumOfExternalPorts,

		// --- QUORUM ---
		"quorum_disk_id":               logicalUnit.QuorumDiskId,
		"quorum_storage_serial_number": logicalUnit.QuorumStorageSerialNumber,
		"quorum_storage_type_id":       logicalUnit.QuorumStorageTypeId,

		// --- NAMESPACE / SUBSYSTEM ---
		"namespace_id":     logicalUnit.NamespaceID,
		"nvm_subsystem_id": logicalUnit.NvmSubsystemId,

		// --- RELOCATION / TIERING ---
		"is_relocation_enabled": logicalUnit.IsRelocationEnabled,
		"tier_level":            logicalUnit.TierLevel,

		// Tier-level usage
		"used_capacity_per_tier_level1": logicalUnit.UsedCapacityPerTierLevel1,
		"used_capacity_per_tier_level2": logicalUnit.UsedCapacityPerTierLevel2,
		"used_capacity_per_tier_level3": logicalUnit.UsedCapacityPerTierLevel3,

		"tier_level_for_new_page_allocation": logicalUnit.TierLevelForNewPageAllocation,

		// --- OPERATION TYPE ---
		"operation_type":                    logicalUnit.OperationType,
		"preparing_operation_progress_rate": logicalUnit.PreparingOperationProgressRate,
	}

	// --- INTERNAL PORT LIST ---
	ports := []map[string]interface{}{}
	for _, p := range logicalUnit.Ports {
		ports = append(ports, map[string]interface{}{
			"port_id":          p.PortID,
			"hostgroup_number": p.HostGroupNumber,
			"hostgroup_name":   p.HostGroupName,
			"lun":              p.Lun,
		})
	}
	lun["ports"] = ports

	// --- EXTERNAL PORT LIST ---
	extPorts := []map[string]interface{}{}
	for _, ep := range logicalUnit.ExternalPorts {
		extPorts = append(extPorts, map[string]interface{}{
			"port_id":          ep.PortID,
			"hostgroup_number": ep.HostGroupNumber,
			"lun":              ep.Lun,
			"wwn":              ep.Wwn,
		})
	}
	lun["external_ports"] = extPorts

	return &lun
}

// UpdateLun updates a lun
func UpdateLun(d *schema.ResourceData) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Mainframe vs block is driven ONLY by cylinder presence.
	_, isMainframe := d.GetOk("cylinder")

	reconcilerUpdateLunRequest, err := UpdateLunRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	// For mainframe volumes, ESE enablement must be updated via a dedicated action.
	if isMainframe && d.HasChange("is_ese_volume") {
		isEse := d.Get("is_ese_volume").(bool)
		if reconcilerUpdateLunRequest.LdevID == nil {
			return nil, fmt.Errorf("ldev_id is missing; cannot update is_ese_volume")
		}
		if err := reconObj.SetEseVolume(*reconcilerUpdateLunRequest.LdevID, isEse); err != nil {
			log.WriteDebug("TFError| error in SetEseVolume reconciler call, err: %v", err)
			return nil, err
		}
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), reconcilerUpdateLunRequest.LdevID, setting.Serial)
	lun, err := reconObj.UpdateLun(reconcilerUpdateLunRequest)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), reconcilerUpdateLunRequest.LdevID, setting.Serial)
		return nil, err
	}

	// If `volume_format_type` is configured, invoke the format action after successful update.
	// Used GetOk rather than HasChange so a re-apply will re-run the format even if the
	// value hasn't changed (this enables users to re-run a format by reapplying).
	if v, ok := d.GetOk("volume_format_type"); ok {
		volFmt := v.(string)
		// treat empty or "NONE" as explicit no-op
		if volFmt == "" || strings.EqualFold(volFmt, "NONE") {
			log.WriteDebug("volume_format_type is NONE/empty; skipping format action")
		} else {
			if reconcilerUpdateLunRequest.LdevID == nil {
				return nil, fmt.Errorf("ldev_id is missing; cannot perform format")
			}

			// map Terraform values to API operation types
			var opType string
			switch strings.ToUpper(volFmt) {
			case "QUICK":
				opType = "QFMT"
			case "NORMAL":
				opType = "FMT"
			default:
				return nil, fmt.Errorf("invalid volume_format_type: %s", volFmt)
			}

			// Determine if data reduction is enabled on the volume (DP)
			// fetch current lun to inspect data reduction mode
			currentLun, err := reconObj.GetLun(*reconcilerUpdateLunRequest.LdevID)
			if err != nil {
				log.WriteDebug("TFError| error fetching LUN to determine data reduction mode: %v", err)
				return nil, err
			}

			if currentLun.Status != "BLK" {
				log.WriteDebug("TFDebug| blocking current volume as volume is not in BLK status before format")
				if err := reconObj.BlockLun(*reconcilerUpdateLunRequest.LdevID); err != nil {
					log.WriteDebug("TFError| error blocking LUN before format: %v", err)
					return nil, err
				}
				log.WriteDebug("TFDebug| successfully blocked the volume before format")
			}

			isForce := false
			if opType == "FMT" {
				if currentLun.DataReductionMode != "" && currentLun.DataReductionMode != "disabled" {
					isForce = true
				}
			}

			fmtReq := reconcilermodel.FormatLdevRequest{
				OperationType:              &opType,
				IsDataReductionForceFormat: &isForce,
			}

			_, err = reconObj.FormatLdev(*reconcilerUpdateLunRequest.LdevID, fmtReq)
			if err != nil {
				log.WriteDebug("TFError| error invoking FormatLdev reconciler call, err: %v", err)
				// Try to restore LDEV state to normal if the format failed.
				if unErr := reconObj.UnblockLun(*reconcilerUpdateLunRequest.LdevID); unErr != nil {
					log.WriteDebug("TFError| error attempting UnblockLun after failed format: %v", unErr)
				}
				return nil, err
			}
		}
	}

	logicalUnit, err := reconObj.GetLun(*lun)
	if err != nil {
		log.WriteDebug("TFError| error in GetLun, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), reconcilerUpdateLunRequest.LdevID, setting.Serial)

	return logicalUnit, nil
}

func UpdateLunRequestFromSchema(d *schema.ResourceData) (*reconcilermodel.UpdateLunRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Input Res: %+v", d)
	log.WriteDebug("Input volume: %+v", d.Get("volume"))
	log.WriteDebug("Input State: %+v", d.Get("state"))
	log.WriteDebug("Input Diff: %+v", d.Get("diff"))

	updateInput := reconcilermodel.UpdateLunRequest{}

	// Mainframe vs block is driven ONLY by cylinder presence.
	_, isMainframe := d.GetOk("cylinder")

	if d.HasChange("size_gb") {
		old, new := d.GetChange("size_gb")
		oldVal := old.(float64)
		newVal := new.(float64)
		if newVal <= oldVal {
			return nil, fmt.Errorf("new size_gb (%.2f) must be greater than old (%.2f)", newVal, oldVal)
		}

		diffGB := newVal - oldVal
		expandSizeStr := utils.ConvertFloatSizeToSmartUnit(diffGB)
		updateInput.ByteFormatCapacity = &expandSizeStr
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateInput.Name = &name
	}

	if !isMainframe {
		if d.HasChange("capacity_saving") {
			dedup := d.Get("capacity_saving").(string)
			if dedup == "" {
				dedup = "disabled"
			}
			updateInput.DataReductionMode = &dedup
		}
	}

	if !isMainframe {
		if v, ok := d.GetOk("data_reduction_process_mode"); ok {
			val := v.(string)
			updateInput.DataReductionProcessMode = &val
		}
	}

	if !isMainframe {
		if d.HasChange("is_compression_acceleration_enabled") {
			if v, ok := d.GetOk("is_compression_acceleration_enabled"); ok {
				val := v.(bool)
				updateInput.IsCompressionAccelerationEnabled = &val
			}
		}
	}

	if !isMainframe {
		if v, ok := d.GetOk("is_alua_enabled"); ok {
			val := v.(bool)
			updateInput.IsAluaEnabled = &val
		}
	}

	pldevID, err := getLdevIdFromSchema(d)
	if err != nil {
		return nil, err
	}
	updateInput.LdevID = pldevID

	log.WriteDebug("updateInput: %+v", updateInput)
	return &updateInput, nil
}

func getLdevIdFromSchema(d *schema.ResourceData) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	idStr := d.Id()
	if idStr == "" {
		return nil, fmt.Errorf("resource ID is empty; cannot determine ldev_id")
	}

	ldevID, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid resource ID '%s': %v", idStr, err)
	}

	log.WriteDebug("ldevID derived from resource ID: %d", ldevID)
	return &ldevID, nil
}

func IsReadExistingMode(raw interface{}) bool {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var getter func(string) interface{}

	// Detect the type and map the methods
	switch v := raw.(type) {
	case *schema.ResourceData:
		getter = v.Get
	case *schema.ResourceDiff:
		getter = v.Get
	default:
		return false // Unsupported type
	}

	// Extract attributes using the assigned getter
	ldevID := getter("ldev_id").(int)
	ldevHex := getter("ldev_id_hex").(string)
	ldevProvided := (ldevID != 0) || (ldevHex != "")

	sizeGB := getter("size_gb").(float64)
	cylinder := getter("cylinder").(int)
	poolID := getter("pool_id").(int)
	poolName := getter("pool_name").(string)
	parityGroup := getter("paritygroup_id").(string)
	extParityGroup := getter("external_paritygroup_id").(string)
	name := getter("name").(string)

	hasCapacity := (sizeGB > 0) || (cylinder > 0)
	hasPlacement := (poolID >= -1) || (poolName != "") || (parityGroup != "") || (extParityGroup != "") || (name != "")

	log.WriteDebug("sizeGB=%v cylinder=%v poolID=%v poolName='%s' parityGroup='%s' extParityGroup='%s' name='%s'",
		sizeGB, cylinder, poolID, poolName, parityGroup, extParityGroup, name)

	// Read-Only mode check
	return ldevProvided && !hasCapacity && !hasPlacement
}
