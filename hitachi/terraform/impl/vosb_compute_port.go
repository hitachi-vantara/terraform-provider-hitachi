package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	"fmt"
	"strings"

	// "io/ioutil"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

var authModeMap = map[string]string{
	"chap":                              "CHAP",
	"chapcomplyingwithinitiatorsetting": "CHAPComplyingWithInitiatorSetting",
	"none":                              "None",
}

func GetVssbComputePortByName(d *schema.ResourceData, vssbaddr string, inputPortName string) (*terraformmodel.PortDetailSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr, ok := d.Get("vosb_block_address").(string)
	if !ok {
		vssbAddr = vssbaddr
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

	portName, ok := d.Get("name").(string)
	if !ok {
		if inputPortName != "" {
			portName = inputPortName
			ok = true
		}
	}

	if ok {
		var portId string = ""
		var found bool
		allPorts, err := GetVssbAllPorts(d, vssbAddr)
		if err != nil {
			log.WriteDebug("TFError| error getting all ports")
			return nil, err
		}

		for _, p := range *allPorts {
			if p.Nickname == portName {
				found = true
				portId = p.ID
				break
			}
		}
		if !found {
			err = fmt.Errorf("did not find port name %v", portName)
			return nil, err
		}
		if portId == "" {
			err = fmt.Errorf("did not find port Id for port name %v", portName)
			return nil, err
		} else {

			return GetPortInfoByID(setting, portId)
		}
	} else {
		err = fmt.Errorf("port name is missing, please specify port name, in the name field of the input")
		return nil, err
	}
}

func AllowChapUsersToAccessComputePort(d *schema.ResourceData) (*terraformmodel.PortDetailSettings, error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_block_address").(string)

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

	portInfo, err := GetVssbComputePortByName(d, "", "")
	if err != nil {
		return nil, err
	}

	auth_mode, ok := d.Get("authentication_settings").(string)
	log.WriteInfo("authentication_settings %+v", auth_mode)
	if !ok {
		err := fmt.Errorf("authentication_settings is not specified, please specify authentication_settings, one of the following: CHAP, CHAPComplyingWithInitiatorSetting, None")
		return nil, err
	}
	authmode := authModeMap[strings.ToLower(auth_mode)]
	if authmode == "" {
		err := fmt.Errorf("invalid authentication_settings specified %v, valid authentication_settings is one of the following: CHAP, CHAPComplyingWithInitiatorSetting, None", auth_mode)
		return nil, err

	}

	target_chap_users, ok := d.Get("target_chap_users").([]interface{})
	if !ok {
		err := fmt.Errorf("target_chap_users is not specified, please specify target_chap_users, empty list[] to remove chap users from the port")
		return nil, err
	}

	targetchapusers := []string{}
	if len(target_chap_users) > 0 {
		log.WriteDebug("target_chap_users: %+v\n", target_chap_users)
		for _, cu := range target_chap_users {
			targetchapusers = append(targetchapusers, cu.(string))
		}
	}

	err = reconObj.AllowChapUsersToAccessComputePort(portInfo.Port.ID, authmode, targetchapusers)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePorts, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_PORTS_FAILED))
		return nil, err
	}

	portDetailInfo, err := reconObj.GetPortInfoByID(portInfo.Port.ID)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPortInfoByID, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_PORTS_FAILED))
		return nil, err
	}

	terraformComputePort := terraformmodel.PortDetailSettings{}
	err = copier.Copy(&terraformComputePort, portDetailInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	return &terraformComputePort, nil
}

func ConvertVssbIscsiPortAuthSettingsToSchema(storagePort *terraformmodel.PortDetailSettings) *map[string]interface{} {
	sp := map[string]interface{}{}

	pas := []map[string]interface{}{}
	dataPas := map[string]interface{}{
		"auth_mode":              storagePort.AuthSettings.AuthMode,
		"is_discovery_chap_auth": storagePort.AuthSettings.IsDiscoveryChapAuth,
		"is_mutual_chap_auth":    storagePort.AuthSettings.IsMutualChapAuth,
	}
	pas = append(pas, dataPas)
	sp["port_auth_settings"] = pas

	cus := []map[string]interface{}{}
	for _, cu := range storagePort.ChapUsers.Data {
		dataCus := map[string]interface{}{
			"chap_user_id":             cu.ID,
			"target_chap_user_name":    cu.TargetChapUserName,
			"initiator_chap_user_name": cu.InitiatorChapUserName,
		}
		cus = append(cus, dataCus)
	}
	sp["chap_users"] = cus

	return &sp
}

func ConvertVssbPortDetailSettingsToSchema(storagePort *terraformmodel.PortDetailSettings) *map[string]interface{} {
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

	cus := []map[string]interface{}{}
	for _, cu := range storagePort.ChapUsers.Data {
		dataCus := map[string]interface{}{
			"chap_user_id":             cu.ID,
			"target_chap_user_name":    cu.TargetChapUserName,
			"initiator_chap_user_name": cu.InitiatorChapUserName,
		}
		cus = append(cus, dataCus)
	}
	sp["chap_users"] = cus

	return &sp
}

func GetVssbAllPorts(d *schema.ResourceData, vssbaddr string) (*[]terraformmodel.StoragePortVssb, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr, ok := d.Get("vosb_block_address").(string)
	if !ok {
		vssbAddr = vssbaddr

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

func GetVssbPortById(d *schema.ResourceData) (*terraformmodel.PortDetailSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}
	portId := d.Get("name").(string)

	return GetPortInfoByID(setting, portId)
}

func GetPortInfoByID(setting reconcilermodel.StorageDeviceSettings, portId string) (*terraformmodel.PortDetailSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)

		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portId)

	terraformPort := terraformmodel.PortDetailSettings{}
	portDetails, err := reconObj.GetPortInfoByID(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPort, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, err
	}

	err = copier.Copy(&terraformPort, portDetails)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_END), portId)

	return &terraformPort, nil

}
