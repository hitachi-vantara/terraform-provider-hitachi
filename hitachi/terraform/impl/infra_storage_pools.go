package terraform

import (
	"errors"
	"strconv"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetInfraGwStoragePools(d *schema.ResourceData) (*[]terraformmodel.InfraStoragePoolInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	if serial == "" && storageId == "" {
		err := errors.New("both serial and storage_id can't be empty. Please specify one")
		return nil, err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
		return nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, err
		}
		d.Set("storage_id", storageId)
	}
	if serial == "" {
		serial, err = common.GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	}

	pool_name := d.Get("pool_name").(string)
	pool_id := -1
	pid, okId := d.GetOk("pool_id")
	if okId {
		pool_id = pid.(int)
	}

	log.WriteDebug("addr : %v, storage_id : %v pool_name: %v pool_id: %v", address, storageId, pool_name, pool_id)

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_BEGIN), setting.Address)
	reconStoragePools, err := reconObj.GetStoragePools(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), setting.Address)
		return nil, err
	}

	var result model.StoragePool
	if pool_name != "" {
		for _, pool := range reconStoragePools.Data {
			if pool.Name == pool_name {
				result.Path = reconStoragePools.Path
				result.Message = reconStoragePools.Message
				result.Data = pool
				break
			}
		}
	}
	if pool_id != -1 {
		for _, pool := range reconStoragePools.Data {
			if pool.PoolId == pool_id {
				result.Path = reconStoragePools.Path
				result.Message = reconStoragePools.Message
				result.Data = pool
				break
			}
		}
	}

	// Converting reconciler to terraform
	terraformStoragePools := terraformmodel.InfraStoragePools{}

	if pool_name != "" || pool_id != -1 {
		err = copier.Copy(&terraformStoragePools, &result)
	} else {
		err = copier.Copy(&terraformStoragePools, reconStoragePools)
	}
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_END), setting.Address)

	return &terraformStoragePools.Data, nil
}

func GetInfraGwStoragePool(storageId, poolId string) (*[]terraformmodel.InfraStoragePoolInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_STORAGE_POOLS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetStoragePool(storageId, poolId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraGwStoragePool, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GW_GET_STORAGE_POOLS_FAILED), setting.Address)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraStoragePools{}
	err = copier.Copy(&terraformResponse, reconResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_STORAGE_POOLS_END), setting.Address)

	return &terraformResponse.Data, nil
}

func ConvertInfraGwStoragePoolToSchema(storagePool *terraformmodel.InfraStoragePoolInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"storage_serial_number":            storage_serial_number,
		"resource_id":                      storagePool.ResourceId,
		"pool_id":                          storagePool.PoolId,
		"ldev_ids":                         storagePool.LdevIds,
		"name":                             storagePool.Name,
		"depletion_threshold_rate":         storagePool.DepletionThresholdRate,
		"free_capacity":                    storagePool.FreeCapacity,
		"free_capacity_in_units":           storagePool.FreeCapacityInUnits,
		"replication_depletion_alert_rate": storagePool.ReplicationDepletionAlertRate,
		"replication_usage_rate":           storagePool.ReplicationUsageRate,
		"resource_group_id":                storagePool.ResourceGroupId,
		"status":                           storagePool.Status,
		"subscription_limit_rate":          storagePool.SubscriptionLimitRate,
		"subscription_rate":                storagePool.SubscriptionRate,
		"subscription_warning_rate":        storagePool.SubscriptionWarningRate,
		"total_capacity":                   storagePool.TotalCapacity,
		"total_capacity_in_unit":           storagePool.TotalCapacityInUnit,
		"type":                             storagePool.Type,
		"virtual_volume_count":             storagePool.VirtualVolumeCount,
		"warning_threshold_rate":           storagePool.WarningThresholdRate,
		"deduplication_enabled":            storagePool.DeduplicationEnabled,
	}

	dpVolumes := []map[string]interface{}{}
	for _, pool := range storagePool.DpVolumes {
		p := map[string]interface{}{
			"logical_unit_id": pool.LogicalUnitId,
			"size":            pool.Size,
		}
		dpVolumes = append(dpVolumes, p)
	}
	sp["dp_volumes"] = dpVolumes

	return &sp
}
