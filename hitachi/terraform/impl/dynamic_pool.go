package terraform

import (
	"fmt"
	"strconv"
	"strings"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	utils "terraform-provider-hitachi/hitachi/common/utils"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetDynamicPools(d *schema.ResourceData) (*[]terraformmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Extract pool_type filter (DP/HTI)
	poolType := ""
	if v, ok := d.GetOk("pool_type"); ok {
		poolType = v.(string)
		log.WriteDebug("TFDebug| GetDynamicPools: pool_type parameter extracted: %s", poolType)
	} else {
		log.WriteDebug("TFDebug| GetDynamicPools: NO pool_type parameter found, returning all pool types")
	}

	// Extract is_mainframe parameter
	var isMainframe *bool
	if v, ok := d.GetOk("is_mainframe"); ok {
		val := v.(bool)
		isMainframe = &val
		log.WriteDebug("TFDebug| GetDynamicPools: is_mainframe parameter extracted: %t", *isMainframe)
	} else {
		log.WriteDebug("TFDebug| GetDynamicPools: NO is_mainframe parameter found, returning all pools")
	}

	// Extract include_detail_info and include_cache_info parameters
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, add all detail types except class
		detailInfoTypes = append(detailInfoTypes, "FMC", "tierPhysicalCapacity", "efficiency", "formattedCapacity", "autoAddPoolVol", "tierDiskType")
		log.WriteDebug("TFDebug| GetDynamicPools: include_detail_info=true, added detail types")
	}
	if includeCacheInfo, ok := d.GetOk("include_cache_info"); ok && includeCacheInfo.(bool) {
		// When includeCacheInfo is true, add class
		detailInfoTypes = append(detailInfoTypes, "class")
		log.WriteDebug("TFDebug| GetDynamicPools: include_cache_info=true, added class")
	}

	detailInfoType := ""
	if len(detailInfoTypes) > 0 {
		detailInfoType = strings.Join(detailInfoTypes, ",")
		log.WriteDebug("TFDebug| GetDynamicPools: final detailInfoType: %s", detailInfoType)
	} else {
		log.WriteDebug("TFDebug| GetDynamicPools: NO detail info parameters set")
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_BEGIN), setting.Serial)
	log.WriteDebug("TFDebug| GetDynamicPools: Calling GetDynamicPools with is_mainframe: %v, pool_type: '%s', detailInfoType: '%s'", isMainframe, poolType, detailInfoType)
	dynamicPools, err := reconObj.GetDynamicPools(isMainframe, poolType, detailInfoType)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOLS_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelDynamicPools := []terraformmodel.DynamicPool{}
	err = copier.Copy(&terraformModelDynamicPools, dynamicPools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_END), setting.Serial)

	return &terraformModelDynamicPools, nil
}

func GetDynamicPoolByName(d *schema.ResourceData) (*terraformmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolName := d.Get("pool_name").(string)

	// Extract include_detail_info and include_cache_info parameters
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, add all detail types except class
		detailInfoTypes = append(detailInfoTypes, "FMC", "tierPhysicalCapacity", "efficiency", "formattedCapacity", "autoAddPoolVol", "tierDiskType")
		log.WriteDebug("TFDebug| GetDynamicPoolByName: include_detail_info=true, added detail types")
	}
	if includeCacheInfo, ok := d.GetOk("include_cache_info"); ok && includeCacheInfo.(bool) {
		// When includeCacheInfo is true, add class
		detailInfoTypes = append(detailInfoTypes, "class")
		log.WriteDebug("TFDebug| GetDynamicPoolByName: include_cache_info=true, added class")
	}

	detailInfoType := ""
	if len(detailInfoTypes) > 0 {
		detailInfoType = strings.Join(detailInfoTypes, ",")
		log.WriteDebug("TFDebug| GetDynamicPoolByName: final detailInfoType: %s", detailInfoType)
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

	pools, err := reconObj.GetDynamicPools(nil, "", detailInfoType)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOLS_FAILED), setting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_END), setting.Serial)

	poolId := -1

	for _, pool := range *pools {
		if strings.EqualFold(pool.PoolName, poolName) {
			poolId = pool.PoolID
			break

		}
	}

	if poolId == -1 {
		return nil, fmt.Errorf("could not find pool with pool name %v", poolName)

	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_BEGIN), poolId, setting.Serial)
	dynamicPool, err := reconObj.GetDynamicPoolById(poolId)
	if err != nil {
		return nil, err
	}

	terraformModelDynamicPool := terraformmodel.DynamicPool{}
	err = copier.Copy(&terraformModelDynamicPool, dynamicPool)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_END), poolId, setting.Serial)

	return &terraformModelDynamicPool, nil
}

func GetDynamicPoolById(d *schema.ResourceData) (*terraformmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	poolId := d.Get("pool_id").(int)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_BEGIN), poolId, setting.Serial)
	dynamicPool, err := reconObj.GetDynamicPoolById(poolId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_DYNAMIC_POOL_ID_FAILED), poolId, setting.Serial)
		return nil, err
	}

	terraformModelDynamicPool := terraformmodel.DynamicPool{}
	err = copier.Copy(&terraformModelDynamicPool, dynamicPool)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOL_ID_END), poolId, setting.Serial)

	return &terraformModelDynamicPool, nil
}

func ConvertDynamicPoolToSchema(dynamicPool *terraformmodel.DynamicPool, serial int) *map[string]interface{} {
	lun := map[string]interface{}{
		"storage_serial_number":                   serial,
		"pool_id":                                 dynamicPool.PoolID,
		"pool_status":                             dynamicPool.PoolStatus,
		"used_capacity_rate":                      dynamicPool.UsedCapacityRate,
		"used_physical_capacity_rate":             dynamicPool.UsedPhysicalCapacityRate,
		"snapshot_count":                          dynamicPool.SnapshotCount,
		"pool_name":                               dynamicPool.PoolName,
		"available_volume_capacity":               dynamicPool.AvailableVolumeCapacity,
		"available_physical_volume_capacity":      dynamicPool.AvailablePhysicalVolumeCapacity,
		"total_pool_capacity":                     dynamicPool.TotalPoolCapacity,
		"total_physical_capacity":                 dynamicPool.TotalPhysicalCapacity,
		"num_of_ldevs":                            dynamicPool.NumOfLdevs,
		"first_ldev_id":                           dynamicPool.FirstLdevID,
		"first_ldev_id_hex":                       utils.IntToHexString(dynamicPool.FirstLdevID),
		"warning_threshold":                       dynamicPool.WarningThreshold,
		"depletion_threshold":                     dynamicPool.DepletionThreshold,
		"virtual_volume_capacity_rate":            dynamicPool.VirtualVolumeCapacityRate,
		"is_mainframe":                            dynamicPool.IsMainframe,
		"is_shrinking":                            dynamicPool.IsShrinking,
		"located_volume_count":                    dynamicPool.LocatedVolumeCount,
		"total_located_capacity":                  dynamicPool.TotalLocatedCapacity,
		"blocking_mode":                           dynamicPool.BlockingMode,
		"total_reserved_capacity":                 dynamicPool.TotalReservedCapacity,
		"reserved_volume_count":                   dynamicPool.ReservedVolumeCount,
		"pool_type":                               dynamicPool.PoolType,
		"duplication_number":                      dynamicPool.DuplicationNumber,
		"effective_capacity":                      dynamicPool.EffectiveCapacity,
		"data_reduction_accelerate_comp_capacity": dynamicPool.DataReductionAccelerateCompCapacity,
		"data_reduction_capacity":                 dynamicPool.DataReductionCapacity,
		"data_reduction_before_capacity":          dynamicPool.DataReductionBeforeCapacity,
		"data_reduction_accelerate_comp_rate":     dynamicPool.DataReductionAccelerateCompRate,
		"compression_rate":                        dynamicPool.CompressionRate,
		"duplication_rate":                        dynamicPool.DuplicationRate,
		"data_reduction_rate":                     dynamicPool.DataReductionRate,
		"snapshot_used_capacity":                  dynamicPool.SnapshotUsedCapacity,
		"suspend_snapshot":                        dynamicPool.SuspendSnapshot,
		"formatted_capacity":                      dynamicPool.FormattedCapacity,
		"auto_add_pool_vol":                       dynamicPool.AutoAddPoolVol,
	}

	// Add nested data reduction objects
	dataReductionAccelerateComp := []map[string]interface{}{
		{
			"is_reduction_capacity_available": dynamicPool.DataReductionAccelerateCompIncludingSystemData.IsReductionCapacityAvailable,
			"reduction_capacity":              dynamicPool.DataReductionAccelerateCompIncludingSystemData.ReductionCapacity,
			"is_reduction_rate_available":     dynamicPool.DataReductionAccelerateCompIncludingSystemData.IsReductionRateAvailable,
		},
	}
	// Add reduction_rate only if isReductionRateAvailable is true
	if dynamicPool.DataReductionAccelerateCompIncludingSystemData.IsReductionRateAvailable {
		dataReductionAccelerateComp[0]["reduction_rate"] = dynamicPool.DataReductionAccelerateCompIncludingSystemData.ReductionRate
	}
	lun["data_reduction_accelerate_comp_including_system_data"] = dataReductionAccelerateComp

	dataReduction := []map[string]interface{}{
		{
			"is_reduction_capacity_available": dynamicPool.DataReductionIncludingSystemData.IsReductionCapacityAvailable,
			"reduction_capacity":              dynamicPool.DataReductionIncludingSystemData.ReductionCapacity,
			"is_reduction_rate_available":     dynamicPool.DataReductionIncludingSystemData.IsReductionRateAvailable,
		},
	}
	// Add reduction_rate only if isReductionRateAvailable is true
	if dynamicPool.DataReductionIncludingSystemData.IsReductionRateAvailable {
		dataReduction[0]["reduction_rate"] = dynamicPool.DataReductionIncludingSystemData.ReductionRate
	}
	lun["data_reduction_including_system_data"] = dataReduction

	capacitiesExcluding := []map[string]interface{}{
		{
			"used_virtual_volume_capacity": dynamicPool.CapacitiesExcludingSystemData.UsedVirtualVolumeCapacity,
			"compressed_capacity":          dynamicPool.CapacitiesExcludingSystemData.CompressedCapacity,
			"deduped_capacity":             dynamicPool.CapacitiesExcludingSystemData.DedupedCapacity,
			"reclaimed_capacity":           dynamicPool.CapacitiesExcludingSystemData.ReclaimedCapacity,
			"system_data_capacity":         dynamicPool.CapacitiesExcludingSystemData.SystemDataCapacity,
			"pre_used_capacity":            dynamicPool.CapacitiesExcludingSystemData.PreUsedCapacity,
			"pre_compressed_capacity":      dynamicPool.CapacitiesExcludingSystemData.PreCompressedCapacity,
			"pre_dedupred_capacity":        dynamicPool.CapacitiesExcludingSystemData.PreDedupredCapacity,
		},
	}
	lun["capacities_excluding_system_data"] = capacitiesExcluding

	// Add efficiency field if available
	if dynamicPool.Efficiency.IsCalculated {
		efficiency := []map[string]interface{}{
			{
				"is_calculated":          dynamicPool.Efficiency.IsCalculated,
				"total_ratio":            dynamicPool.Efficiency.TotalRatio,
				"compression_ratio":      dynamicPool.Efficiency.CompressionRatio,
				"snapshot_ratio":         dynamicPool.Efficiency.SnapshotRatio,
				"provisioning_rate":      dynamicPool.Efficiency.ProvisioningRate,
				"calculation_start_time": dynamicPool.Efficiency.CalculationStartTime,
				"calculation_end_time":   dynamicPool.Efficiency.CalculationEndTime,
			},
		}

		// Add dedup_and_compression if available
		if dynamicPool.Efficiency.DedupAndCompression.TotalRatio != "" {
			dedupComp := []map[string]interface{}{
				{
					"total_ratio":       dynamicPool.Efficiency.DedupAndCompression.TotalRatio,
					"compression_ratio": dynamicPool.Efficiency.DedupAndCompression.CompressionRatio,
					"dedupe_ratio":      dynamicPool.Efficiency.DedupAndCompression.DedupeRatio,
					"reclaim_ratio":     dynamicPool.Efficiency.DedupAndCompression.ReclaimRatio,
				},
			}
			efficiency[0]["dedupe_and_compression"] = dedupComp
		}

		// Add accelerated_compression if available
		if dynamicPool.Efficiency.AcceleratedCompression.TotalRatio != "" {
			accelComp := []map[string]interface{}{
				{
					"total_ratio":       dynamicPool.Efficiency.AcceleratedCompression.TotalRatio,
					"compression_ratio": dynamicPool.Efficiency.AcceleratedCompression.CompressionRatio,
					"reclaim_ratio":     dynamicPool.Efficiency.AcceleratedCompression.ReclaimRatio,
				},
			}
			efficiency[0]["accelerated_compression"] = accelComp
		}

		lun["efficiency"] = efficiency
	}

	return &lun
}
