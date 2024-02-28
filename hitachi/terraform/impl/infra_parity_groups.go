package terraform

import (
	"strconv"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetInfraParityGroups(d *schema.ResourceData) (*[]terraformmodel.InfraParityGroupInfo, error) {
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

	parityGroupIdsMap := map[string]string{}
	ids, ok := d.GetOk("parity_group_ids")
	if ok {
		pgIds := ids.([]interface{})
		for _, value := range pgIds {
			switch typedValue := value.(type) {
			case string:
				parityGroupIdsMap[typedValue] = typedValue
			}
		}
		log.WriteDebug("TFDebug| parity group filter will be apply")
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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_PARITY_GROUPS_BEGIN), setting.Address)
	reconParityGroups, err := reconObj.GetParityGroups(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraParityGroups, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_PARITY_GROUPS_FAILED), setting.Address)
		return nil, err
	}

	var result model.ParityGroups
	if len(parityGroupIdsMap) > 0 {
		result.Path = reconParityGroups.Path
		result.Message = reconParityGroups.Message
		for _, p := range reconParityGroups.Data {
			_, ok := parityGroupIdsMap[p.ParityGroupId]
			if ok {
				result.Data = append(result.Data, p)
			}
		}
	}

	// Converting reconciler to terraform
	terraformStoragePorts := terraformmodel.InfraParityGroups{}

	if len(parityGroupIdsMap) > 0 {
		err = copier.Copy(&terraformStoragePorts, &result)
	} else {
		err = copier.Copy(&terraformStoragePorts, reconParityGroups)
	}

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_PARITY_GROUPS_END), setting.Address)

	return &terraformStoragePorts.Data, nil
}

func ConvertInfraGwParityGroupToSchema(pg *terraformmodel.InfraParityGroupInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"storage_serial_number":      storage_serial_number,
		"resource_id":                pg.ResourceId,
		"parity_group_id":            pg.ParityGroupId,
		"free_capacity":              utils.TransformSizeToUnit(uint64(pg.FreeCapacity)),
		"resource_group_id":          pg.ResourceGroupId,
		"total_capacity":             utils.TransformSizeToUnit(uint64(pg.TotalCapacity)),
		"ldev_ids":                   pg.LdevIds,
		"raid_level":                 pg.RaidLevel,
		"drive_type":                 pg.DriveType,
		"copyback_mode":              pg.CopybackMode,
		"status":                     pg.Status,
		"is_pool_array_group":        pg.IsPoolArrayGroup,
		"is_accelerated_compression": pg.IsAcceleratedCompression,
		"is_encryption_enabled":      pg.IsEncryptionEnabled,
	}

	return &sp
}
