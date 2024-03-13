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

func CreateInfraVolume(d *schema.ResourceData) (*terraformmodel.MtInfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageId, serial, setting, err := common.GetInfraGatewaySettings(d, nil)
	storage_serial_number = *serial

	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings , err: %v", err)
		return nil, err
	}

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraVolumeRequestFromSchema(d, setting)
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

		volumeInfo, ok := reconObj.GetVolumeByName(*storageId, name.(string))
		if ok {

			volData, err := reconObj.ReconcileVolume(*storageId, &reconcilerCreateVolRequest, &volumeInfo.ResourceId)
			if err != nil {
				log.WriteDebug("TFError| error in Create Volume, err: %v", err)
				return nil, err
			}
			terraformModelVol := terraformmodel.MtInfraVolumeInfo{}
			err = copier.Copy(&terraformModelVol, volData)
			if err != nil {
				log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
				return nil, err
			}
			return &terraformModelVol, nil
		}
	}

	volData, err := reconObj.ReconcileVolume(*storageId, &reconcilerCreateVolRequest, nil)
	if err != nil {
		log.WriteDebug("TFError| error in Create Volume, err: %v", err)
		return nil, err
	}
	terraformModelVol := terraformmodel.MtInfraVolumeInfo{}
	err = copier.Copy(&terraformModelVol, volData)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	return &terraformModelVol, nil
}

func GetInfraVolumes(d *schema.ResourceData) (*[]terraformmodel.MtInfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageId, serial, setting, err := common.GetInfraGatewaySettings(d, nil)
	storage_serial_number = *serial
	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings, err: %v", err)
		return nil, err
	}

	startLdevID := d.Get("start_ldev_id").(int)

	if startLdevID <= 0 {
		startLdevID = 0
		// return nil, nil, fmt.Errorf("start_ldev_id must be greater than or equal to 0")
	}

	endLdevID := d.Get("end_ldev_id").(int)
	if endLdevID <= 0 {
		endLdevID = 65000
		// return nil, nil, fmt.Errorf("end_ldev_id must be greater than or equal to 0")
	}

	if endLdevID < startLdevID {
		return nil, fmt.Errorf("end_ldev_id must be greater than or equal to start_ldev_id")
	}

	isUndefindLdev := d.Get("undefined_ldev").(bool)

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_BEGIN), setting.Address)

	if setting.PartnerId == nil {

		response, err := reconObj.GetVolumesFromLdevIds(*storageId, &startLdevID, &endLdevID)
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
				if p.EmulationType == "NOT DEFINED" {
					result.Data = append(result.Data, p)
				}
			}
		} else {
			result = *response

		}

		// Converting reconciler to terraform
		terraformResponse := terraformmodel.MTInfraVolumes{}

		err = copier.Copy(&terraformResponse, result)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

		return &terraformResponse.Data, nil
	}

	mtResponse, err := reconObj.GetVolumesByPartnerSubscriberID(*storageId, startLdevID, endLdevID)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, err
	}
	var result reconcilermodel.MTVolumes
	if isUndefindLdev {
		result.Path = mtResponse.Path
		result.Message = mtResponse.Message
		for _, p := range mtResponse.Data {
			if p.EmulationType == "NOT DEFINED" {
				result.Data = append(result.Data, p)
			}
		}
	} else {
		result = *mtResponse

	}
	terraformMtResponse := terraformmodel.MTInfraVolumes{}
	err = copier.Copy(&terraformMtResponse, result)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

	return &terraformMtResponse.Data, nil

}

func GetInfraVolume(d *schema.ResourceData) (*terraformmodel.MtInfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageId, serial, setting, err := common.GetInfraGatewaySettings(d, nil)
	storage_serial_number = *serial
	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings, err: %v", err)
		return nil, err
	}

	var ldev_id int

	ldevID, _ := d.GetOk("ldev_id")
	if ldevID.(int) != -1 {
		ldev_id = ldevID.(int)
		if ldev_id < 0 {
			return nil, fmt.Errorf("ldev_id must be greater than or equal to 0 is %s", ldevID)
		}
	}

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_BEGIN), setting.Address)
	procResponse, mtResponse, err := reconObj.GetVolumeByLDevId(*storageId, ldev_id)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, err
	}
	terraformMtResponse := terraformmodel.MtInfraVolumeInfo{}
	if procResponse != nil {
		err = copier.Copy(&terraformMtResponse, procResponse)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

	} else if mtResponse != nil {

		err = copier.Copy(&terraformMtResponse, mtResponse)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

	return &terraformMtResponse, nil
}

func GetInfraSingleVolume(d *schema.ResourceData) (*terraformmodel.MtInfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId, _, address, err := common.GetValidateStorageIDFromSerialResource(d, nil)
	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings , err: %v", err)
		return nil, err
	}

	if serial == "" {
		serial, err = common.GetSerialFromStorageId(*address, *storageId)
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

	vol_id := d.State().ID

	storageSetting, err := cache.GetInfraSettingsFromCache(*address)
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
	response, err := reconObj.GetVolumeByID(*storageId, vol_id)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.MtInfraVolumeInfo{}

	err = copier.Copy(&terraformResponse, response)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)

	return &terraformResponse, nil
}

// func ConvertInfraVolumeToSchema(pg *terraformmodel.MtInfraVolumeInfo) *map[string]interface{} {
// 	var pga []string
// 	sp := map[string]interface{}{
// 		"storage_serial_number":          storage_serial_number,
// 		"resource_id":                    pg.ResourceId,
// 		"deduplication_compression_mode": pg.DeduplicationCompressionMode,
// 		"emulation_type":                 pg.EmulationType,
// 		"format_or_shred_rate":           pg.FormatOrShredRate,
// 		"ldev_id":                        pg.LdevId,
// 		"name":                           pg.Name,
// 		"parity_group_id":                append(pga, pg.ParityGroupId),
// 		"pool_id":                        pg.PoolId,
// 		"resource_group_id":              pg.ResourceGroupId,
// 		"status":                         pg.Status,
// 		"total_capacity":                 pg.TotalCapacity,
// 		"used_capacity":                  pg.UsedCapacity,
// 		"used_capacity_in_mb":            common.BytesToMegabytes(pg.UsedCapacity),
// 		"virtual_storage_device_id":      pg.VirtualStorageDeviceId,
// 		"stripe_size":                    pg.StripeSize,
// 		"type":                           pg.Type,
// 		"path_count":                     pg.PathCount,
// 		"provision_type":                 pg.ProvisionType,
// 		"is_command_device":              pg.IsCommandDevice,
// 		"logical_unit_id_hex_format":     pg.LogicalUnitIdHexFormat,
// 		"virtual_logical_unit_id":        pg.VirtualLogicalUnitId,
// 		"naa_id":                         pg.NaaId,
// 		"dedup_compression_progress":     pg.DedupCompressionProgress,
// 		"dedup_compression_status":       pg.DedupCompressionStatus,
// 		"is_alua_enabled":                pg.IsALUA,
// 		"is_dynamic_pool_volume":         pg.IsDynamicPoolVolume,
// 		"is_journal_pool_volume":         pg.IsJournalPoolVolume,
// 		"is_pool_volume":                 pg.IsPoolVolume,
// 		"pool_name":                      pg.PoolName,
// 		"quorum_disk_id":                 pg.QuorumDiskId,
// 		"is_in_gad_pair":                 pg.IsInGadPair,
// 		"is_vvol":                        pg.IsVVol,
// 		"total_capacity_in_mb":           common.BytesToMegabytes(pg.TotalCapacity),
// 	}

// 	return &sp
// }

func ConvertPartnersInfraVolumeToSchema(pg *terraformmodel.MtInfraVolumeInfo) *map[string]interface{} {
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
		"used_capacity_in_mb":            common.BytesToMegabytes(pg.UsedCapacity),
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
		"total_capacity_in_mb":           common.BytesToMegabytes(pg.TotalCapacity),
		"partner_id":                     pg.PartnerId,
		"subscriber_id":                  pg.SubscriberId,
		"entitlement_status":             pg.EntitlementStatus,
		"free_capacity_in_mb":            common.BytesToMegabytes(pg.TotalCapacity - pg.UsedCapacity),
	}

	return &sp
}

func CreateInfraVolumeRequestFromSchema(d *schema.ResourceData, setting *reconcilermodel.InfraGwSettings) (*terraformmodel.InfraVolumeTypes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.InfraVolumeTypes{}

	name, ok := d.GetOk("name")

	if ok {
		createInput.Name = name.(string)
	}

	pool_id := d.Get("pool_id").(int)
	if pool_id != -1 {

		createInput.PoolID = &pool_id
	}

	lun_id := d.Get("ldev_id").(int)
	if lun_id != -1 {
		createInput.LunId = &lun_id
	}

	resourceGroupId := d.Get("resource_group_id").(int)
	if resourceGroupId >= 0 {
		createInput.ResourceGroupId = &resourceGroupId
	}

	paritygroup_id, ok := d.GetOk("paritygroup_id")
	if ok {
		createInput.ParityGroupId = paritygroup_id.(string)

	}

	capacity, ok := d.GetOk("size_gb")

	if ok {
		inString := common.GbToMbString(capacity.(float64))
		createInput.Capacity = inString
	}

	system, ok := d.GetOk("system")
	if ok {
		createInput.System = system.(string)

	} else {
		createInput.System = reconcilermodel.DefaultSystemSerialNumber
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

	storageId, _, setting, err := common.GetInfraGatewaySettings(d, nil)
	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings , err: %v", err)
		return err
	}

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}

	volumeId := d.State().ID

	_, err = reconObj.ReconcileVolume(*storageId, nil, &volumeId)
	if err != nil {
		log.WriteDebug("TFError| error in ReconcileVolume Delete volume, err: %v", err)
		return err
	}

	return nil
}

func UpdateInfraVolume(d *schema.ResourceData) (*terraformmodel.MtInfraVolumeInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	volumeID := d.State().ID

	storageId, serial, setting, err := common.GetInfraGatewaySettings(d, nil)
	storage_serial_number = *serial
	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings , err: %v", err)
		return nil, err
	}

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraVolumeRequestFromSchema(d, setting)
	if err != nil {
		return nil, err
	}

	reconcilerCreateVolRequest := reconcilermodel.CreateVolumeParams{}
	err = copier.Copy(&reconcilerCreateVolRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	volData, err := reconObj.ReconcileVolume(*storageId, &reconcilerCreateVolRequest, &volumeID)
	if err != nil {
		log.WriteDebug("TFError| error in Update Volume, err: %v", err)
		return nil, err
	}

	terraformModelVolume := terraformmodel.MtInfraVolumeInfo{}
	err = copier.Copy(&terraformModelVolume, volData)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	return &terraformModelVolume, nil
}
