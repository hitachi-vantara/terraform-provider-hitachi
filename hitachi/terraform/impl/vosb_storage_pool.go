package terraform

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
	"strings"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
)

func GetAllStoragePools(d *schema.ResourceData) (*[]terraformmodel.StoragePool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_BEGIN))
	reconStoragePools, err := reconObj.GetAllStoragePools()
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_POOLS_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStoragePools := []terraformmodel.StoragePool{}
	err = copier.Copy(&terraformStoragePools, reconStoragePools.Data)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_END))

	return &terraformStoragePools, nil
}

func GetStoragePoolsByPoolNames(d *schema.ResourceData) (*[]terraformmodel.StoragePool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)
	spNames, ok := d.GetOk("storage_pool_names")
	storagePoolNames := make([]string, 0)
	if ok {
		poolOpt := spNames.([]interface{})
		poolNames := make([]string, len(poolOpt))
		for index, value := range poolOpt {
			switch typedValue := value.(type) {
			case string:
				poolNames[index] = typedValue
			}
		}
		storagePoolNames = append(storagePoolNames, poolNames...)
	}

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	names := strings.Join(storagePoolNames, ",")

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_BEGIN), names)
	reconStoragePools, err := reconObj.GetStoragePoolsByPoolNames(storagePoolNames)
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_POOL_FAILED), names)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStoragePools := []terraformmodel.StoragePool{}
	err = copier.Copy(&terraformStoragePools, reconStoragePools.Data)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_END), names)

	return &terraformStoragePools, nil
}

func AddDrivesToStoragePool(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_DRIVES_STORAGE_POOL_BEGIN))

	updateInput, err := ConvertVssbStoragePoolFromSchema(d)
	if err != nil {
		return err
	}

	reconStoragePoolResource := reconcilermodel.StoragePoolResource{}
	err = copier.Copy(&reconStoragePoolResource, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return err
	}

	err = reconObj.AddDrivesToStoragePool(&reconStoragePoolResource)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_ADD_DRIVES_STORAGE_POOL_FAILED))
		log.WriteDebug("TFError| error in Updating ComputeNode - ReconcileComputeNode , err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_DRIVES_STORAGE_POOL_END))

	return nil
}

func ConvertStoragePoolToSchema(storagePool *terraformmodel.StoragePool) *map[string]interface{} {
	sp := map[string]interface{}{
		"pool_id":                     storagePool.ID,
		"pool_name":                   storagePool.Name,
		"protection_domain_id":        storagePool.ProtectionDomainId,
		"status_summary":              storagePool.StatusSummary,
		"status":                      storagePool.Status,
		"total_capacity":              storagePool.TotalCapacity,
		"total_raw_capacity":          storagePool.TotalRawCapacity,
		"used_capacity":               storagePool.UsedCapacity,
		"free_capacity":               storagePool.FreeCapacity,
		"total_physical_capacity":     storagePool.TotalPhysicalCapacity,
		"meta_data_physical_capacity": storagePool.MetaDataPhysicalCapacity,
		"reserved_physical_capacity":  storagePool.ReservedPhysicalCapacity,
		"usable_physical_capacity":    storagePool.UsablePhysicalCapacity,
		"blocked_physical_capacity":   storagePool.BlockedPhysicalCapacity,
	}

	cm := []map[string]interface{}{}
	w := map[string]interface{}{
		"used_capacity_rate":                     storagePool.CapacityManage.UsedCapacityRate,
		"maximum_reserve_rate":                   storagePool.CapacityManage.MaximumReserveRate,
		"threshold_warning":                      storagePool.CapacityManage.ThresholdWarning,
		"threshold_depletion":                    storagePool.CapacityManage.ThresholdDepletion,
		"threshold_storage_controller_depletion": storagePool.CapacityManage.ThresholdStorageControllerDepletion,
	}
	cm = append(cm, w)
	sp["capacity_manage"] = cm

	se := []map[string]interface{}{}
	e := map[string]interface{}{
		"efficiency_data_reduction":                        storagePool.SavingEffect.EfficiencyDataReduction,
		"pre_capacity_data_reduction":                      storagePool.SavingEffect.PreCapacityDataReduction,
		"post_capacity_data_reduction":                     storagePool.SavingEffect.PostCapacityDataReduction,
		"total_efficiency_status":                          storagePool.SavingEffect.TotalEfficiencyStatus,
		"data_reduction_without_system_data_status":        storagePool.SavingEffect.DataReductionWithoutSystemDataStatus,
		"total_efficiency":                                 storagePool.SavingEffect.TotalEfficiency,
		"data_reduction_without_system_data":               storagePool.SavingEffect.DataReductionWithoutSystemData,
		"pre_capacity_data_reduction_without_system_data":  storagePool.SavingEffect.PreCapacityDataReductionWithoutSystemData,
		"post_capacity_data_reduction_without_system_data": storagePool.SavingEffect.PostCapacityDataReductionWithoutSystemData,
		"calculation_start_time":                           storagePool.SavingEffect.CalculationStartTime.Format("2006-01-02 15:04:05"),
		"calculation_end_time":                             storagePool.SavingEffect.CalculationEndTime.Format("2006-01-02 15:04:05"),
	}
	se = append(se, e)
	sp["saving_effects"] = se

	sp["number_of_volumes"] = storagePool.NumberOfVolumes
	sp["redundant_policy"] = storagePool.RedundantPolicy
	sp["redundant_type"] = storagePool.RedundantType
	sp["data_redundancy"] = storagePool.DataRedundancy
	sp["storage_controller_capacities_general_status"] = storagePool.StorageControllerCapacitiesGeneralStatus
	sp["total_volume_capacity"] = storagePool.TotalVolumeCapacity
	sp["provisioned_volume_capacity"] = storagePool.ProvisionedVolumeCapacity
	sp["other_volume_capacity"] = storagePool.OtherVolumeCapacity
	sp["temporary_volume_capacity"] = storagePool.TemporaryVolumeCapacity
	sp["rebuild_capacity_policy"] = storagePool.RebuildCapacityPolicy
	sp["rebuild_capacity_status"] = storagePool.RebuildCapacityStatus

	rc := []map[string]interface{}{}
	r := map[string]interface{}{
		"number_of_tolerable_drive_failures": storagePool.RebuildCapacityResourceSetting.NumberOfTolerableDriveFailures,
	}
	rc = append(rc, r)
	sp["rebuild_capacity_resource_setting"] = rc

	rr := []map[string]interface{}{}
	rr1 := map[string]interface{}{
		"number_of_drives": storagePool.RebuildableResources.NumberOfDrives,
	}
	rr = append(rr, rr1)
	sp["rebuildable_resources"] = rr

	return &sp
}

func ConvertVssbStoragePoolFromSchema(d *schema.ResourceData) (*terraformmodel.StoragePoolResource, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	input := terraformmodel.StoragePoolResource{}

	storagePoolName := d.Get("storage_pool_name")
	input.StoragePoolName = storagePoolName.(string)

	addAllOfflineDrives, ok := d.GetOk("add_all_offline_drives")
	if ok {
		aaod := addAllOfflineDrives.(bool)
		input.AddAllOfflineDrives = aaod
	}

	rawDriveIds := d.Get("drive_ids")

	// If 'drive_ids' is not provided, set input.DriveIds to nil
	if rawDriveIds == nil {
		input.DriveIds = nil
	} else {
		// Convert it to a slice of strings
		var driveIds []string
		for _, v := range rawDriveIds.([]interface{}) {
			if str, ok := v.(string); ok {
				driveIds = append(driveIds, str)
			} else {
				return nil, fmt.Errorf("invalid value in drive_ids list: expected string, got %T", v)
			}
		}
		// Assign the populated driveIds slice to input.DriveIds
		input.DriveIds = driveIds
	}

	// Check if 'AddOfflineDrives' and 'DriveIds' are both provided, which is invalid
	if input.AddAllOfflineDrives && len(input.DriveIds) > 0 {
		return nil, fmt.Errorf("invalid input: 'AddOfflineDrives' cannot be true when 'DriveIds' are provided. Set one, not both")
	}

	log.WriteDebug("input: %+v", input)
	return &input, nil
}
