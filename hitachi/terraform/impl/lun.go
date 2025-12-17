package terraform

import (
	"fmt"
	"strings"
	"strconv"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	// terrcommon "terraform-provider-hitachi/hitachi/terraform/common"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetLun gets a lun
func GetLun(d *schema.ResourceData) (*sangatewaymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// check if this is getting executed from "data"
	lunID := d.Get("ldev_id").(int)

	// check if this is getting executed from "resource"
	_, lunOk := d.GetOk("ldev_id")
	if !lunOk {
		lunFromState := d.State().ID
		log.WriteDebug("TFDebug| lunFromState from state: %s", lunFromState)
		if lunFromState != "" {
			lun, err := strconv.Atoi(lunFromState)
			if err != nil {
				log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
				return nil, err
			}
			lunID = lun
		}
	}

	// vnext
	// var lunID int

	// finalLdev, err := terrcommon.ExtractLdevFields(d, "ldev_id", "ldev_hex")
	// if err != nil {
	// 	log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), lunID)
	// 	return nil, err
	// }
	// if finalLdev != nil {
	// 	lunID = *finalLdev
	// } else {
	// 	lunFromState := d.State().ID
	// 	log.WriteDebug("TFDebug| lunFromState from state: %s", lunFromState)
	// 	if lunFromState != "" {
	// 		lun, err := strconv.Atoi(lunFromState)
	// 		if err != nil {
	// 			log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
	// 			return nil, err
	// 		}
	// 		lunID = lun
	// 	} else {
	// 		return nil, fmt.Errorf("state ID for ldev_id is missing")
	// 	}
	// }

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

	startLdevID := d.Get("start_ldev_id").(int)
	if startLdevID < 0 {
		return nil, fmt.Errorf("start_ldev_id must be greater than or equal to 0")
	}

	endLdevID := d.Get("end_ldev_id").(int)
	if endLdevID < 0 {
		return nil, fmt.Errorf("end_ldev_id must be greater than or equal to 0")
	}

	if endLdevID < startLdevID {
		return nil, fmt.Errorf("end_ldev_id must be greater than or equal to start_ldev_id")
	}

	// vnext
	// startLdevID := 0
	// finalStartLdev, err := terrcommon.ExtractLdevFields(d, "start_ldev_id", "start_ldev_hex")
	// if err != nil {
	// 	return nil, err
	// }
	// if finalStartLdev != nil {
	// 	startLdevID = *finalStartLdev
	// }

	// endLdevID := 0
	// finalEndLdev, err := terrcommon.ExtractLdevFields(d, "end_ldev_id", "end_ldev_hex")
	// if err != nil {
	// 	return nil, err
	// }
	// if finalEndLdev != nil {
	// 	endLdevID = *finalEndLdev
	// }

	// if startLdevID < 0 {
	// 	return nil, fmt.Errorf("start_ldev_id/start_ldev_hex must be greater than or equal to 0")
	// }

	// if endLdevID < 0 {
	// 	return nil, fmt.Errorf("end_ldev_id/end_ldev_hex must be greater than or equal to 0")
	// }

	// if endLdevID < startLdevID {
	// 	return nil, fmt.Errorf("end_ldev_id/end_ldev_hex must be greater than or equal to start_ldev_id/start_ldev_hex")
	// }

	isUndefindLdev := d.Get("undefined_ldev").(bool)

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
	luns, err := reconObj.GetRangeOfLuns(startLdevID, endLdevID, isUndefindLdev)
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

	pldevID := CheckSchemaIfLunGet(d)
	if pldevID != nil {
		lun, err := reconObj.GetLun(*pldevID)
		if err != nil {
			return nil, err
		}
		if lun.ByteFormatCapacity == "" {
			// does not exist, or in the process of being deleted
			return nil, fmt.Errorf("volume does not exist")
		}
		return lun, nil
	}

	reconcilerCreateLunRequest, err := CreateLunRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	lun, err := reconObj.SetLun(reconcilerCreateLunRequest)
	if err != nil {
		log.WriteDebug("TFError| error in SetLun, err: %v", err)
		return nil, err
	}

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

	// size_gb → ByteFormatCapacity
	if size_gb, ok := d.GetOk("size_gb"); ok {
		createInput.ByteFormatCapacity = utils.ConvertFloatSizeToSmartUnit(size_gb.(float64))
	}

	ldevId, ok := d.GetOk("ldev_id")
	if ok {
		lid := ldevId.(int)
		createInput.LdevID = &lid
	}

	// vnext
	// finalLdev, err := terrcommon.ExtractLdevFields(d, "ldev_id", "ldev_hex")
	// if err != nil {
	// 	return nil, err
	// }
	// if finalLdev != nil {
	// 	createInput.LdevID = finalLdev
	// }

	name, ok := d.GetOk("name")
	if ok {
		label := name.(string)
		createInput.Name = &label
	}

	dedup_mode, ok := d.GetOk("capacity_saving")
	if ok {
		dedup := dedup_mode.(string)
		createInput.DataReductionMode = &dedup
	}

	if v, ok := d.GetOk("is_data_reduction_share_enabled"); ok {
		val := v.(bool)
		createInput.IsDataReductionSharedVolumeEnabled = &val
	}

	if v, ok := d.GetOk("is_compression_acceleration_enabled"); ok {
		val := v.(bool)
		createInput.IsCompressionAccelerationEnabled = &val
	}

	pool_id := d.Get("pool_id").(int)
	pool_name := d.Get("pool_name").(string)
	paritygroup_id := d.Get("paritygroup_id").(string)
	external_paritygroup_id := d.Get("external_paritygroup_id").(string)
	log.WriteDebug("Pool ID=%v Pool Name=%v PG=%v ExPG=%v\n", pool_id, pool_name, paritygroup_id, external_paritygroup_id)

	if pool_id >= 0 {
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

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
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

	poolId := -1

	for _, pool := range *pools {
		if strings.EqualFold(pool.PoolName, poolName) {
			poolId = pool.PoolID
			break

		}
	}

	return &poolId, nil
}

func CheckSchemaIfLunGet(d *schema.ResourceData) *int {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	fields := []string{
		"size_gb",
		"name",
		// "capacity_saving",
		"pool_id",
		"pool_name",
		"paritygroup_id",
	}

	for _, f := range fields {
		if _, ok := d.GetOk(f); ok {
			return nil
		}
	}

	ldevId, ok := d.GetOk("ldev_id")
	if ok {
		lid := ldevId.(int)
		return &lid
	}
	return nil
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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevID, setting.Serial)

	err = reconObj.DeleteLun(lunID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevID, setting.Serial)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevID, setting.Serial)

	return nil
}

func ConvertLunToSchema(logicalUnit *sangatewaymodel.LogicalUnit, serial int) *map[string]interface{} {

	lun := map[string]interface{}{
		"storage_serial_number": serial,

		// --- BASIC IDENTIFIERS ---
		"ldev_id":              logicalUnit.LdevID,
		// vnext
		// "ldev_hex":             utils.IntToHexString(logicalUnit.LdevID),
		"virtual_ldev_id":      logicalUnit.VirtualLdevID,
		"clpr_id":              logicalUnit.ClprID,
		"emulation_type":       logicalUnit.EmulationType,
		"byte_format_capacity": logicalUnit.ByteFormatCapacity,
		"block_capacity":       logicalUnit.BlockCapacity,

		// --- ATTRIBUTES ---
		"attributes": logicalUnit.Attributes,
		"label":      logicalUnit.Label,
		"status":     logicalUnit.Status,

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

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), reconcilerUpdateLunRequest.LdevID, setting.Serial)
	lun, err := reconObj.UpdateLun(reconcilerUpdateLunRequest)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), reconcilerUpdateLunRequest.LdevID, setting.Serial)
		return nil, err
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

	if d.HasChange("capacity_saving") {
		dedup := d.Get("capacity_saving").(string)
		if dedup == "" {
			dedup = "disabled"
		}
		updateInput.DataReductionMode = &dedup
	}

	if v, ok := d.GetOk("data_reduction_process_mode"); ok {
		val := v.(string)
		updateInput.DataReductionProcessMode = &val
	}

	if d.HasChange("is_compression_acceleration_enabled") {
		if v, ok := d.GetOk("is_compression_acceleration_enabled"); ok {
			val := v.(bool)
			updateInput.IsCompressionAccelerationEnabled = &val
		}
	}

	if v, ok := d.GetOk("is_alua_enabled"); ok {
		val := v.(bool)
		updateInput.IsAluaEnabled = &val
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
