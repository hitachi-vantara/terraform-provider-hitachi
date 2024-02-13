package terraform

import (
	// "encoding/json"

	"fmt"
	"strconv"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconcilermodel "terraform-provider-hitachi/hitachi/infra_gw/model"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	// model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func CreateInfraVolume(d *schema.ResourceData) (*terraformmodel.InfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
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

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}
	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	log.WriteDebug("Found settings for address %s>>>>>>>>>>>>>>>>>", address, setting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraVolumeRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	reconcilerCreateVolRequest := reconcilermodel.CreateVolumeParams{}
	err = copier.Copy(&reconcilerCreateVolRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	name, ok := d.GetOk("name")
	// Check if the volume exists
	if ok {

		volumeInfo, ok := reconObj.GetVolumeByName(storageId, name.(string))
		if ok {

			volData, err := reconObj.ReconcileVolume(storageId, &reconcilerCreateVolRequest, &volumeInfo.ResourceId)
			if err != nil {
				log.WriteDebug("TFError| error in Create Volume, err: %v", err)
				return nil, err
			}
			terraformModelVol := terraformmodel.InfraVolumeInfo{VolumeInfo: *volData}
			return &terraformModelVol, nil
		}
	}

	volData, err := reconObj.ReconcileVolume(storageId, &reconcilerCreateVolRequest, nil)
	if err != nil {
		log.WriteDebug("TFError| error in Create Volume, err: %v", err)
		return nil, err
	}

	terraformModelVol := terraformmodel.InfraVolumeInfo{VolumeInfo: *volData}

	return &terraformModelVol, nil
}

func GetInfraVolumes(d *schema.ResourceData) (*[]terraformmodel.InfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
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
	d.Set("serial", storage_serial_number)

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

	isUndefindLdev := d.Get("undefined_ldev").(bool)

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_BEGIN), setting.Address)
	response, err := reconObj.GetVolumes(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, err
	}

	var result reconcilermodel.Volumes
	if isUndefindLdev {
		result.Path = response.Path
		result.Message = response.Message
		for _, p := range response.Data {
			if p.LdevId >= startLdevID && p.LdevId <= endLdevID {
				if p.EmulationType == "NOT DEFINED" {
					result.Data = append(result.Data, p)
				}
			}
		}
	} else {
		result.Path = response.Path
		result.Message = response.Message
		for _, p := range response.Data {
			if p.LdevId >= startLdevID && p.LdevId <= endLdevID {
				if p.EmulationType != "NOT DEFINED" {
					result.Data = append(result.Data, p)
				}
			}
		}
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraVolumes{}

	err = copier.Copy(&terraformResponse, result)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

	return &terraformResponse.Data, nil
}

func GetInfraVolume(d *schema.ResourceData) (*[]terraformmodel.InfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
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
	d.Set("serial", storage_serial_number)

	ldevID := d.Get("ldev_id").(int)
	if ldevID < 0 {
		return nil, fmt.Errorf("ldev_id must be greater than or equal to 0")
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_BEGIN), setting.Address)
	response, err := reconObj.GetVolumes(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, err
	}

	var result reconcilermodel.Volumes

	result.Path = response.Path
	result.Message = response.Message
	for _, p := range response.Data {
		if p.LdevId == ldevID {
			result.Data = append(result.Data, p)
			break
		}
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraVolumes{}

	err = copier.Copy(&terraformResponse, result)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

	return &terraformResponse.Data, nil
}

func ConvertInfraVolumeToSchema(pg *terraformmodel.InfraVolumeInfo) *map[string]interface{} {
	var pga []string
	sp := map[string]interface{}{
		"storage_serial_number":          storage_serial_number,
		"resource_id":                    pg.ResourceId,
		"deduplication_compression_mode": pg.DeduplicationCompressionMode,
		"emulation_type":                 pg.EmulationType,
		"format_or_shred_rate":           pg.FormatOrShredRate,
		"ldev_id":                        pg.LdevId,
		"name":                           pg.Name,
		"parity_group_id":                append(pga, pg.ParityGroupId),
		"pool_id":                        pg.PoolId,
		"resource_group_id":              pg.ResourceGroupId,
		"status":                         pg.Status,
		"total_capacity":                 pg.TotalCapacity,
		"used_capacity":                  pg.UsedCapacity,
		"virtual_storage_device_id":      pg.VirtualStorageDeviceId,
		"stripe_size":                    pg.StripeSize,
		"type":                           pg.Type,
		"path_count":                     pg.PathCount,
		"provision_type":                 pg.ProvisionType,
		"is_command_device":              pg.IsCommandDevice,
		"logical_unit_id_hex_format":     pg.LogicalUnitIdHexFormat,
		"virtual_logical_unit_id":        pg.VirtualLogicalUnitId,
		"naa_id":                         pg.NaaId,
		"dedup_compression_progress":     pg.DedupCompressionProgress,
		"dedup_compression_status":       pg.DedupCompressionStatus,
		"is_alua_enabled":                pg.IsALUA,
		"is_dynamic_pool_volume":         pg.IsDynamicPoolVolume,
		"is_journal_pool_volume":         pg.IsJournalPoolVolume,
		"is_pool_volume":                 pg.IsPoolVolume,
		"pool_name":                      pg.PoolName,
		"quorum_disk_id":                 pg.QuorumDiskId,
		"is_in_gad_pair":                 pg.IsInGadPair,
		"is_vvol":                        pg.IsVVol,
	}

	return &sp
}

func CreateInfraVolumeRequestFromSchema(d *schema.ResourceData) (*terraformmodel.InfraVolumeTypes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.InfraVolumeTypes{}

	name, ok := d.GetOk("name")
	if ok {
		createInput.Name = name.(string)
	}

	pool_id, ok := d.GetOk("pool_id")
	if ok {
		createInput.PoolID = pool_id.(int)
	} else {
		return nil, fmt.Errorf("pool_id is mandatory for new volume creation")
	}

	lun_id, ok := d.GetOk("lun_id")
	if ok {
		createInput.LunId = lun_id.(int)
	}

	resourceGroupId, ok := d.GetOk("resource_group_id")
	if ok {
		createInput.ResourceGroupId = resourceGroupId.(int)
	}

	paritygroup_id, ok := d.GetOk("parity_group_id")
	if ok {
		createInput.ParityGroupId = paritygroup_id.(string)
	} else {

		return nil, fmt.Errorf("paritygroup_id is mandatory for new volume creation")
	}

	capacity, ok := d.GetOk("capacity")
	if ok {
		createInput.Capacity = capacity.(string)
	} else {

		return nil, fmt.Errorf("capacity is mandatory for new volume creation")
	}
	system, ok := d.GetOk("system")
	if ok {
		createInput.System = system.(string)

	} else {

		return nil, fmt.Errorf("system is mandatory for new volume creation")
	}
	deduplicationCompressionMode, ok := d.GetOk("deduplication_compression_mode")
	if ok {
		createInput.DeduplicationCompressionMode = deduplicationCompressionMode.(string)
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func DeleteInfraVolume(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return err
		}
		d.Set("storage_id", storageId)
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return err
	}

	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}

	volumeId := d.Id()

	_, err = reconObj.ReconcileVolume(storageId, nil, &volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in ReconcileVolume Delete volume, err: %v", err)
		return err
	}

	return nil
}

func UpdateInfraVolume(d *schema.ResourceData) (*terraformmodel.InfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)
	volumeID := d.Id()

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
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

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraVolumeRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	reconcilerCreateVolRequest := reconcilermodel.CreateVolumeParams{}
	err = copier.Copy(&reconcilerCreateVolRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	volData, err := reconObj.ReconcileVolume(storageId, &reconcilerCreateVolRequest, &volumeID)
	if err != nil {
		log.WriteDebug("TFError| error in Update Volume, err: %v", err)
		return nil, err
	}

	terraformModelLun := terraformmodel.InfraVolumeInfo{VolumeInfo: *volData}

	return &terraformModelLun, nil
}
