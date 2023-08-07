package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	"strconv"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// "terraform-provider-hitachi/hitachi/common/utils"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	//common "terraform-provider-hitachi/hitachi/terraform/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetDynamicPools(d *schema.ResourceData) (*[]terraformmodel.DynamicPool, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DYNAMIC_POOLS_BEGIN), setting.Serial)
	dynamicPools, err := reconObj.GetDynamicPools()
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
		"storage_serial_number":              serial,
		"pool_id":                            dynamicPool.PoolID,
		"pool_status":                        dynamicPool.PoolStatus,
		"used_capacity_rate":                 dynamicPool.UsedCapacityRate,
		"used_physical_capacity_rate":        dynamicPool.UsedPhysicalCapacityRate,
		"snapshot_count":                     dynamicPool.SnapshotCount,
		"pool_name":                          dynamicPool.PoolName,
		"available_volume_capacity":          dynamicPool.AvailableVolumeCapacity,
		"available_physical_volume_capacity": dynamicPool.AvailablePhysicalVolumeCapacity,
		"total_pool_capacity":                dynamicPool.TotalPoolCapacity,
		"total_physical_capacity":            dynamicPool.TotalPhysicalCapacity,
		"num_of_ldevs":                       dynamicPool.NumOfLdevs,
		"first_ldev_id":                      dynamicPool.FirstLdevID,
		"warning_threshold":                  dynamicPool.WarningThreshold,
		"depletion_threshold":                dynamicPool.DepletionThreshold,
		"virtual_volume_capacity_rate":       dynamicPool.VirtualVolumeCapacityRate,
		"is_mainframe":                       dynamicPool.IsMainframe,
		"is_shrinking":                       dynamicPool.IsShrinking,
		"located_volume_count":               dynamicPool.LocatedVolumeCount,
		"total_located_capacity":             dynamicPool.TotalLocatedCapacity,
		"blocking_mode":                      dynamicPool.BlockingMode,
		"total_reserved_capacity":            dynamicPool.TotalReservedCapacity,
		"reserved_volume_count":              dynamicPool.ReservedVolumeCount,
		"pool_type":                          dynamicPool.PoolType,
		"duplication_number":                 dynamicPool.DuplicationNumber,
		"duplication_rate":                   dynamicPool.DuplicationRate,
		"data_reduction_rate":                dynamicPool.DataReductionRate,
		"snapshot_used_capacity":             dynamicPool.SnapshotUsedCapacity,
		"suspend_snapshot":                   dynamicPool.SuspendSnapshot,
	}

	return &lun
}
