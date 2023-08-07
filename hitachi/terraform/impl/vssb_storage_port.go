package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetVssbPort(d *schema.ResourceData) (*terraformmodel.PortWithAuthSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

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

	portId := d.Get("port_id").(string)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portId)

	reconPort, reconPortAuthSetting, err := reconObj.GetPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePort, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformPort := terraformmodel.PortWithAuthSettings{}
	err = copier.Copy(&terraformPort.Port, reconPort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	err = copier.Copy(&terraformPort.AuthSettings, reconPortAuthSetting)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_END), portId)

	return &terraformPort, nil
}

func GetVssbStoragePorts(d *schema.ResourceData) (*[]terraformmodel.StoragePortVssb, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_BEGIN))

	reconStoragePorts, err := reconObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePorts, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_PORTS_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStoragePorts := terraformmodel.StoragePortsVssb{}
	err = copier.Copy(&terraformStoragePorts, reconStoragePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_PORTS_END))

	return &terraformStoragePorts.Data, nil
}

func ConvertVssbStoragePortToSchema(storagePort *terraformmodel.StoragePortVssb) *map[string]interface{} {
	sp := map[string]interface{}{
		//"id":                    storagePort.ID,
		"protocol":              storagePort.Protocol,
		"type":                  storagePort.Type,
		"nickname":              storagePort.Nickname,
		"name":                  storagePort.Name,
		"configured_port_speed": storagePort.ConfiguredPortSpeed,
		"port_speed":            storagePort.PortSpeed,
		"por_speed_duplex":      storagePort.PortSpeedDuplex,
		"protection_domain_id":  storagePort.ProtectionDomainID,
		"storage_node_id":       storagePort.StorageNodeID,
		"interface_name":        storagePort.InterfaceName,
		"status_summary":        storagePort.StatusSummary,
		"status":                storagePort.Status,
	}

	if storagePort.Protocol == "FC" {
		fc := []map[string]interface{}{}
		dataFc := map[string]interface{}{
			"connection_type":        storagePort.FcInformation.ConnectionType,
			"sfp_data_transfer_rate": storagePort.FcInformation.SfpDataTransferRate,
			"physical_wwn":           storagePort.FcInformation.PhysicalWwn,
		}
		fc = append(fc, dataFc)
		sp["fc_information"] = fc
	}

	if storagePort.Protocol == "iSCSI" {

		ipv4 := []map[string]interface{}{}
		dataIpv4 := map[string]interface{}{
			"address":         storagePort.IscsiInformation.Ipv4Information.Address,
			"subnet_mask":     storagePort.IscsiInformation.Ipv4Information.SubnetMask,
			"default_gateway": storagePort.IscsiInformation.Ipv4Information.DefaultGateway,
		}
		ipv4 = append(ipv4, dataIpv4)

		ipv6 := []map[string]interface{}{}
		dataIpv6 := map[string]interface{}{
			"linklocal_address_mode": storagePort.IscsiInformation.Ipv6Information.LinklocalAddressMode,
			"linklocal_address":      storagePort.IscsiInformation.Ipv6Information.LinklocalAddress,
			"global_address_mode":    storagePort.IscsiInformation.Ipv6Information.GlobalAddressMode,
			"global_address_1":       storagePort.IscsiInformation.Ipv6Information.GlobalAddress1,
			"subnet_prefix_length_1": storagePort.IscsiInformation.Ipv6Information.SubnetPrefixLength1,
			"default_gateway":        storagePort.IscsiInformation.Ipv6Information.DefaultGateway,
		}
		ipv6 = append(ipv6, dataIpv6)

		isnsServers := []map[string]interface{}{}
		for _, server := range storagePort.IscsiInformation.IsnsServers {
			dataIsnsServer := map[string]interface{}{
				"index":       server.Index,
				"server_name": server.ServerName,
				"port":        server.Port,
			}
			isnsServers = append(isnsServers, dataIsnsServer)
		}

		iscsiInfo := map[string]interface{}{
			"ip_mode":                storagePort.IscsiInformation.IPMode,
			"delayed_ack":            storagePort.IscsiInformation.DelayedAck,
			"mtu_size":               storagePort.IscsiInformation.MtuSize,
			"mac_address":            storagePort.IscsiInformation.MacAddress,
			"is_isns_client_enabled": storagePort.IscsiInformation.IsIsnsClientEnabled,
			"ipv4_information":       ipv4,
			"ipv6_information":       ipv6,
			"isns_servers":           isnsServers,
		}

		finalIscsiInfo := []map[string]interface{}{}
		finalIscsiInfo = append(finalIscsiInfo, iscsiInfo)

		sp["iscsi_information"] = finalIscsiInfo
	}

	return &sp
}

func ConvertVssbPortWithAuthSettingsToSchema(storagePort *terraformmodel.PortWithAuthSettings) *map[string]interface{} {
	sp := map[string]interface{}{
		//"id":                    storagePort.Port.ID,
		"protocol":              storagePort.Port.Protocol,
		"type":                  storagePort.Port.Type,
		"nickname":              storagePort.Port.Nickname,
		"name":                  storagePort.Port.Name,
		"configured_port_speed": storagePort.Port.ConfiguredPortSpeed,
		"port_speed":            storagePort.Port.PortSpeed,
		"por_speed_duplex":      storagePort.Port.PortSpeedDuplex,
		"protection_domain_id":  storagePort.Port.ProtectionDomainID,
		"storage_node_id":       storagePort.Port.StorageNodeID,
		"interface_name":        storagePort.Port.InterfaceName,
		"status_summary":        storagePort.Port.StatusSummary,
		"status":                storagePort.Port.Status,
	}

	if storagePort.Port.Protocol == "FC" {
		fc := []map[string]interface{}{}
		dataFc := map[string]interface{}{
			"connection_type":        storagePort.Port.FcInformation.ConnectionType,
			"sfp_data_transfer_rate": storagePort.Port.FcInformation.SfpDataTransferRate,
			"physical_wwn":           storagePort.Port.FcInformation.PhysicalWwn,
		}
		fc = append(fc, dataFc)
		sp["fc_information"] = fc
	}

	if storagePort.Port.Protocol == "iSCSI" {
		ipv4 := []map[string]interface{}{}
		dataIpv4 := map[string]interface{}{
			"address":         storagePort.Port.IscsiInformation.Ipv4Information.Address,
			"subnet_mask":     storagePort.Port.IscsiInformation.Ipv4Information.SubnetMask,
			"default_gateway": storagePort.Port.IscsiInformation.Ipv4Information.DefaultGateway,
		}
		ipv4 = append(ipv4, dataIpv4)

		ipv6 := []map[string]interface{}{}
		dataIpv6 := map[string]interface{}{
			"linklocal_address_mode": storagePort.Port.IscsiInformation.Ipv6Information.LinklocalAddressMode,
			"linklocal_address":      storagePort.Port.IscsiInformation.Ipv6Information.LinklocalAddress,
			"global_address_mode":    storagePort.Port.IscsiInformation.Ipv6Information.GlobalAddressMode,
			"global_address_1":       storagePort.Port.IscsiInformation.Ipv6Information.GlobalAddress1,
			"subnet_prefix_length_1": storagePort.Port.IscsiInformation.Ipv6Information.SubnetPrefixLength1,
			"default_gateway":        storagePort.Port.IscsiInformation.Ipv6Information.DefaultGateway,
		}
		ipv6 = append(ipv6, dataIpv6)

		isnsServers := []map[string]interface{}{}
		for _, server := range storagePort.Port.IscsiInformation.IsnsServers {
			dataIsnsServer := map[string]interface{}{
				"index":       server.Index,
				"server_name": server.ServerName,
				"port":        server.Port,
			}
			isnsServers = append(isnsServers, dataIsnsServer)
		}

		iscsiInfo := map[string]interface{}{
			"ip_mode":                storagePort.Port.IscsiInformation.IPMode,
			"delayed_ack":            storagePort.Port.IscsiInformation.DelayedAck,
			"mtu_size":               storagePort.Port.IscsiInformation.MtuSize,
			"mac_address":            storagePort.Port.IscsiInformation.MacAddress,
			"is_isns_client_enabled": storagePort.Port.IscsiInformation.IsIsnsClientEnabled,
			"ipv4_information":       ipv4,
			"ipv6_information":       ipv6,
			"isns_servers":           isnsServers,
		}

		finalIscsiInfo := []map[string]interface{}{}
		finalIscsiInfo = append(finalIscsiInfo, iscsiInfo)

		sp["iscsi_information"] = finalIscsiInfo
	}

	pas := []map[string]interface{}{}
	dataPas := map[string]interface{}{
		"auth_mode":              storagePort.AuthSettings.AuthMode,
		"is_discovery_chap_auth": storagePort.AuthSettings.IsDiscoveryChapAuth,
		"is_mutual_chap_auth":    storagePort.AuthSettings.IsMutualChapAuth,
	}
	pas = append(pas, dataPas)
	sp["port_auth_settings"] = pas

	return &sp
}
