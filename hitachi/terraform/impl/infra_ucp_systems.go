package terraform

import (
	"fmt"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetInfraUcpSystems(d *schema.ResourceData) (*[]terraformmodel.InfraUcpSystemInfo, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_UCP_SYSTEMS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetUcpSystems()
	if err != nil {
		log.WriteDebug("TFError| error getting GetUcpSystems, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_UCP_SYSTEMS_FAILED), setting.Address)
		return nil, err
	}

	serial := d.Get("serial_number").(string)
	name := d.Get("name").(string)

	var result model.UcpSystem
	if serial != "" {
		found := false
		for _, ucp := range reconResponse.Data {
			if ucp.SerialNumber == serial {
				result.Path = reconResponse.Path
				result.Message = reconResponse.Message
				result.Data = ucp
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("UCP System with serial number %s not found", serial)
			log.WriteDebug("UCP System with serial number %s not found", serial)
			return nil, err
		}
	}
	if name != "" {
		found := false
		for _, ucp := range reconResponse.Data {
			if ucp.Name == name {
				result.Path = reconResponse.Path
				result.Message = reconResponse.Message
				result.Data = ucp
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("UCP System with name %s not found", name)
			return nil, err
		}
	}
	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraUcpSystems{}
	if serial != "" || name != "" {
		err = copier.Copy(&terraformResponse, &result)
	} else {
		err = copier.Copy(&terraformResponse, reconResponse)
	}
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_UCP_SYSTEMS_END), setting.Address)

	return &terraformResponse.Data, nil
}

func GetInfraUcpSystem(d *schema.ResourceData, serial string) (*[]terraformmodel.InfraUcpSystemInfo, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetStorageDevice("")
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraStorageDevice, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_DEVICES_FAILED), setting.Address)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraUcpSystems{}
	err = copier.Copy(&terraformResponse, reconResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_END), setting.Address)

	return &terraformResponse.Data, nil
}

func ConvertInfraUcpSystemToSchema(pg *terraformmodel.InfraUcpSystemInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"resource_id":         pg.ResourceId,
		"name":                pg.Name,
		"resource_state":      pg.ResourceState,
		"serial_number":       pg.SerialNumber,
		"gateway_address":     pg.GatewayAddress,
		"model":               pg.Model,
		"vcenter":             pg.Vcenter,
		"zone":                pg.Zone,
		"vcenter_resource_id": pg.VcenterResourceId,
		"region":              pg.Region,
		"workload_type":       pg.WorkloadType,
		"result_status":       pg.ResultStatus,
		"result_message":      pg.ResultMessage,
		"plugin_registered":   pg.PluginRegistered,
		"linked":              pg.Linked,
	}

	geoInformation := []map[string]interface{}{}
	a := map[string]interface{}{
		"geo_location": pg.GeoInformation.GeoLocation,
		"country":      pg.GeoInformation.Country,
		"latitude":     pg.GeoInformation.Latitude,
		"longitude":    pg.GeoInformation.Longitude,
		"zipcode":      pg.GeoInformation.Zipcode,
	}
	geoInformation = append(geoInformation, a)
	sp["geo_information"] = geoInformation

	computeDevices := []map[string]interface{}{}
	for _, device := range pg.ComputeDevices {
		p := map[string]interface{}{
			"resource_id":          device.ResourceId,
			"bmc_address":          device.BmcAddress,
			"bmc_firmware_version": device.BmcFirmwareVersion,
			"bios_version":         device.BiosVersion,
			"resource_state":       device.ResourceState,
			"model":                device.Model,
			"serial":               device.Serial,
			"is_management":        device.IsManagement,
			"health_status":        device.HealthStatus,
			"gateway_address":      device.GatewayAddress,
		}
		computeDevices = append(computeDevices, p)
	}
	sp["compute_devices"] = computeDevices

	storageDevices := []map[string]interface{}{}
	for _, device := range pg.StorageDevices {
		p := map[string]interface{}{
			"serial_number":     device.SerialNumber,
			"resource_id":       device.ResourceId,
			"address":           device.Address,
			"model":             device.Model,
			"microcode_version": device.MicrocodeVersion,
			"resource_state":    device.ResourceState,
			"health_state":      device.HealthState,
			"ucp_systems":       device.UcpSystems,
			"svp_ip":            device.SvpIp,
			"gateway_address":   device.GatewayAddress,
		}
		storageDevices = append(storageDevices, p)
	}
	sp["storage_devices"] = storageDevices

	ethernetSwitches := []map[string]interface{}{}
	for _, device := range pg.EthernetSwitches {
		p := map[string]interface{}{
			"resource_id":      device.ResourceId,
			"address":          device.Address,
			"name":             device.Name,
			"model":            device.Model,
			"serial_number":    device.SerialNumber,
			"firmware_version": device.FirmwareVersion,
			"resource_state":   device.ResourceState,
			"health_status":    device.HealthStatus,
			"gateway_address":  device.GatewayAddress,
			"is_management":    device.IsManagement,
		}
		ethernetSwitches = append(ethernetSwitches, p)
	}
	sp["ethernet_switches"] = ethernetSwitches

	fibreChannelSwitches := []map[string]interface{}{}
	for _, device := range pg.FibreChannelSwitches {
		p := map[string]interface{}{
			"resource_id":      device.ResourceId,
			"address":          device.Address,
			"model":            device.Model,
			"serial_number":    device.SerialNumber,
			"firmware_version": device.FirmwareVersion,
			"resource_state":   device.ResourceState,
			"health_state":     device.HealthState,
			"switch_name":      device.SwitchName,
			"gateway_address":  device.GatewayAddress,
		}
		fibreChannelSwitches = append(fibreChannelSwitches, p)
	}
	sp["fibre_channel_switches"] = fibreChannelSwitches

	return &sp
}
