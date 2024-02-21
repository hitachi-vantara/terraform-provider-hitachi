package terraform

import (
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

func GetInfraGwStoragePorts(d *schema.ResourceData) (*[]terraformmodel.InfraStoragePortInfo, error) {
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

	port_id := d.Get("port_id").(string)

	log.WriteDebug("addr : %v, storage_id : %v", address, storageId)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_PORTS_BEGIN), setting.Address)
	reconStoragePorts, err := reconObj.GetStoragePorts(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePorts, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_PORTS_FAILED), setting.Address)
		return nil, err
	}

	var result model.StoragePort
	if port_id != "" {
		for _, port := range reconStoragePorts.Data {
			if port.PortId == port_id {
				result.Path = reconStoragePorts.Path
				result.Message = reconStoragePorts.Message
				result.Data = port
				break
			}
		}
	}

	// Converting reconciler to terraform
	terraformStoragePorts := terraformmodel.InfraStoragePorts{}

	if port_id != "" {
		err = copier.Copy(&terraformStoragePorts, &result)
	} else {
		err = copier.Copy(&terraformStoragePorts, reconStoragePorts)
	}

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_PORTS_END), setting.Address)

	return &terraformStoragePorts.Data, nil
}

func ConvertInfraGwStoragePortToSchema(storagePort *terraformmodel.InfraStoragePortInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"port_id":             storagePort.PortId,
		"type":                storagePort.Type,
		"speed":               storagePort.Speed,
		"resource_group_id":   storagePort.ResourceGroupId,
		"wwn":                 storagePort.Wwn,
		"resource_id":         storagePort.ResourceId,
		"attribute":           storagePort.Attribute,
		"connection_type":     storagePort.ConnectionType,
		"fabric_on":           storagePort.FabricOn,
		"mode":                storagePort.Mode,
		"is_security_enabled": storagePort.IsSecurityEnabled,
	}

	return &sp
}
