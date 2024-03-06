package terraform

import (
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

func GetInfraStoragePorts(d *schema.ResourceData) (*[]terraformmodel.InfraStoragePortInfo, *[]terraformmodel.InfraMTStoragePortInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageId, setting, err := common.GetInfraGatewaySettings(d, nil)

	if err != nil {
		log.WriteDebug("TFError| error in GetInfraGatewaySettings, err: %v", err)
		return nil, nil, err
	}

	d.Set("serial", storage_serial_number)

	port_id := d.Get("port_id").(string)

	reconObj, err := reconimpl.NewEx(*setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_PORTS_BEGIN), setting.Address)

	if setting.PartnerId == nil {
		reconStoragePorts, err := reconObj.GetStoragePorts(*storageId)
		if err != nil {
			log.WriteDebug("TFError| error getting GetStoragePorts, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_PORTS_FAILED), setting.Address)
			return nil, nil, err
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
			return nil, nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_PORTS_END), setting.Address)

		return &terraformStoragePorts.Data, nil, nil
	}

	mtResponse, err := reconObj.GetStoragePortsByPartnerIdOrSubscriberId(*storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVolumes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_VOLUMES_FAILED), setting.Address)
		return nil, nil, err
	}
	terraformMtResponse := terraformmodel.InfraMTStoragePorts{}
	err = copier.Copy(&terraformMtResponse, mtResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_VOLUMES_END), setting.Address)
	return nil, &terraformMtResponse.Data, nil
}

func ConvertInfraStoragePortToSchema(storagePort *terraformmodel.InfraStoragePortInfo) *map[string]interface{} {
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

func ConvertInfraMTStoragePortToSchema(storagePort *terraformmodel.InfraMTStoragePortInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"port_id":             storagePort.PortId,
		"port_type":           storagePort.PortType,
		"speed":               storagePort.Speed,
		"resource_group_id":   storagePort.ResourceGroupId,
		"wwn":                 storagePort.Wwn,
		"attribute":           storagePort.Attribute,
		"connection_type":     storagePort.ConnectionType,
		"fabric_on":           storagePort.FabricOn,
		"mode":                storagePort.Mode,
		"is_security_enabled": storagePort.IsSecurityEnabled,
	}

	return &sp
}
