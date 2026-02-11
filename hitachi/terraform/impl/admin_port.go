package terraform

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	provmanager "terraform-provider-hitachi/hitachi/storage/admin/provisioner"
	provimpl "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Datasources -------------------

func DatasourceAdminOnePortRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	v, ok := d.GetOk("port_id")
	if !ok {
		return diag.Errorf("port_id must be specified")
	}
	portID := v.(string)

	// call provisioner directly
	provObj, err := getPortProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	portInfo, err := provObj.GetPortByID(portID)
	if err != nil {
		log.WriteDebug("failed to get port %v: %v", portID, err)
		return diag.FromErr(fmt.Errorf("failed to get port %v: %v", portID, err))
	}

	log.WriteDebug("port %+v", portInfo)
	if portInfo == nil {
		return nil
	}

	if err := d.Set("port_info", []map[string]interface{}{convertOnePortInfoToSchema(portInfo)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set port_info: %w", err))
	}

	// Set the resource ID so Terraform shows computed fields
	d.SetId(portInfo.ID)

	return nil
}

func DatasourceAdminMultiplePortsRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	params := gwymodel.GetPortParams{}

	if v, ok := d.GetOk("protocol"); ok {
		val := v.(string)
		params.Protocol = &val
	}

	log.WriteDebug("Params: %+v", params)

	// call provisioner directly
	provObj, err := getPortProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	ports, err := provObj.GetPorts(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ports_info", convertMultiplePortInfosListToSchema(ports)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port_count", ports.Count); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// ------------------- Resources -------------------

// ResourceAdminPortRead reads the current state of a port
func ResourceAdminPortRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	v, ok := d.GetOk("port_id")
	if !ok {
		return diag.Errorf("port_id is required for read operation")
	}
	portID := v.(string)

	// Get reconciler manager
	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get port information through reconciler
	portInfo, err := recObj.ReconcileReadAdminPort(portID)
	if err != nil {
		log.WriteDebug("failed to get port %v: %v", portID, err)
		return diag.FromErr(fmt.Errorf("failed to get port %v: %v", portID, err))
	}

	if err := d.Set("port_info", []map[string]interface{}{convertOnePortInfoToSchema(portInfo)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set port_info: %w", err))
	}

	// Set the resource ID so Terraform shows computed fields
	d.SetId(portInfo.ID)

	return nil
}

// ResourceAdminPortUpdate updates port configuration
func ResourceAdminPortUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting port update")

	serial := d.Get("serial").(int)

	v, ok := d.GetOk("port_id")
	if !ok {
		return diag.Errorf("port_id is required for update operation")
	}
	portID := v.(string)

	// Build update parameters from schema
	params := buildUpdatePortParams(d)
	log.WriteDebug("Params: %+v", params)

	// Get reconciler manager
	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the port through reconciler
	err = recObj.ReconcileUpdateAdminPort(portID, *params)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo(fmt.Sprintf("port %s updated successfully", portID))

	// Read the updated port information
	return ResourceAdminPortRead(d)
}

// ------------------- Helpers -------------------

func buildUpdatePortParams(d *schema.ResourceData) *gwymodel.UpdatePortParams {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	params := &gwymodel.UpdatePortParams{}

	// --------------------
	// Helper functions
	// --------------------
	// getRawPath returns the raw config value
	getRawPath := func(d *schema.ResourceData, flatKey string) cty.Value {
		rc := d.GetRawConfig() // THIS IS cty.Value

		parts := strings.Split(flatKey, ".")

		cur := rc
		for _, p := range parts {

			// list index?
			if idx, err := strconv.Atoi(p); err == nil {
				cur = cur.Index(cty.NumberIntVal(int64(idx)))
				continue
			}

			// attribute
			cur = cur.GetAttr(p)
		}

		return cur
	}

	setBool := func(flatKey string) *bool {
		log.WriteDebug("Key: %s", flatKey)

		raw := getRawPath(d, flatKey)

		// Case 1: User explicitly set value in config
		if raw.IsKnown() && !raw.IsNull() {
			b := raw.True() // correct API
			log.WriteDebug("Explicit %s = %v", flatKey, b)
			return &b
		}

		// Case 2: Update — changed
		if d.HasChange(flatKey) {
			newVal := d.Get(flatKey).(bool)
			log.WriteDebug("Changed %s = %v", flatKey, newVal)
			return &newVal
		}

		return nil
	}

	setString := func(flatKey string) *string {
		log.WriteDebug("Key: %s", flatKey)

		raw := getRawPath(d, flatKey)

		if raw.IsKnown() && !raw.IsNull() {
			s := raw.AsString()
			log.WriteDebug("Explicit %s = '%s'", flatKey, s)
			return &s
		}

		if d.HasChange(flatKey) {
			newVal := d.Get(flatKey).(string)
			log.WriteDebug("Changed %s = '%s'", flatKey, newVal)
			return &newVal
		}

		return nil
	}

	setInt := func(flatKey string) *int {
		log.WriteDebug("Key: %s", flatKey)

		raw := getRawPath(d, flatKey)

		if raw.IsKnown() && !raw.IsNull() {
			i64, _ := raw.AsBigFloat().Int64()
			i := int(i64)
			log.WriteDebug("Explicit %s = %d", flatKey, i)
			return &i
		}

		if d.HasChange(flatKey) {
			newVal := d.Get(flatKey).(int)
			log.WriteDebug("Changed %s = %d", flatKey, newVal)
			return &newVal
		}

		return nil
	}

	// --------------------
	// TOP-LEVEL FIELDS
	// --------------------
	params.PortSpeed = setString("port_speed")
	params.PortSecurity = setBool("port_security")

	// --------------------
	// FC INFORMATION
	// --------------------
	if _, ok := d.GetOk("fc_information"); ok {

		fc := &gwymodel.FCInformationParam{}
		setAny := false

		if s := setString("fc_information.0.al_pa"); s != nil {
			fc.AlPa = s
			setAny = true
		}

		if b := setBool("fc_information.0.fabric_switch_setting"); b != nil {
			log.WriteDebug("Inside if: fabric_switch_setting b: %+v", b)
			fc.FabricSwitchSetting = b
			setAny = true
		}

		if s := setString("fc_information.0.connection_type"); s != nil {
			fc.ConnectionType = s
			setAny = true
		}

		if setAny {
			params.FcInformationParam = fc
		}
	}

	// --------------------
	// iSCSI INFORMATION
	// --------------------
	if _, ok := d.GetOk("iscsi_information"); ok {

		iscsi := &gwymodel.IscsiInformationParam{}
		setAny := false

		// leaf-level fields
		for _, key := range []string{
			"vlan_use",
			"selective_ack",
			"delayed_ack",
			"isns_server_mode",
		} {
			k := "iscsi_information.0." + key
			if b := setBool(k); b != nil {
				switch key {
				case "vlan_use":
					iscsi.VlanUse = b
				case "selective_ack":
					iscsi.SelectiveAck = b
				case "delayed_ack":
					iscsi.DelayedAck = b
				case "isns_server_mode":
					iscsi.IsnsServerMode = b
				}
				setAny = true
			}
		}

		for _, key := range []string{
			"add_vlan_id",
			"delete_vlan_id",
			"tcp_port",
			"keep_alive_timer",
			"isns_server_port",
		} {
			k := "iscsi_information.0." + key
			if i := setInt(k); i != nil {
				switch key {
				case "add_vlan_id":
					iscsi.AddVlanID = i
				case "delete_vlan_id":
					iscsi.DeleteVlanID = i
				case "tcp_port":
					iscsi.TCPPort = i
				case "keep_alive_timer":
					iscsi.KeepAliveTimer = i
				case "isns_server_port":
					iscsi.IsnsServerPort = i
				}
				setAny = true
			}
		}

		for _, key := range []string{
			"ip_mode",
			"window_size",
			"mtu_size",
			"isns_server_ip_address",
		} {
			k := "iscsi_information.0." + key
			if s := setString(k); s != nil {
				switch key {
				case "ip_mode":
					iscsi.IPMode = s
				case "window_size":
					iscsi.WindowSize = s
				case "mtu_size":
					iscsi.MTUSize = s
				case "isns_server_ip_address":
					iscsi.IsnsServerIP = s
				}
				setAny = true
			}
		}

		// IPv4 nested
		if _, ok := d.GetOk("iscsi_information.0.ipv4_information"); ok {

			ip4 := &gwymodel.IPv4Information{}
			ip4Set := false

			for _, key := range []string{"address", "subnet_mask", "default_gateway"} {
				k := "iscsi_information.0.ipv4_information.0." + key
				if s := setString(k); s != nil {
					switch key {
					case "address":
						ip4.Address = *s
					case "subnet_mask":
						ip4.SubnetMask = *s
					case "default_gateway":
						ip4.DefaultGateway = *s
					}
					ip4Set = true
				}
			}

			if ip4Set {
				iscsi.IPv4Info = ip4
				setAny = true
			}
		}

		// IPv6 nested
		if _, ok := d.GetOk("iscsi_information.0.ipv6_information"); ok {

			ip6 := &gwymodel.IPv6Information{}
			ip6Set := false

			for _, key := range []string{
				"linklocal", "linklocal_address", "linklocal_address_status",
				"global", "global_address", "global_address_status",
				"default_gateway",
			} {
				k := "iscsi_information.0.ipv6_information.0." + key
				if s := setString(k); s != nil {
					switch key {
					case "linklocal":
						ip6.Linklocal = *s
					case "linklocal_address":
						ip6.LinklocalAddress = *s
					case "linklocal_address_status":
						ip6.LinklocalAddressStatus = *s
					case "global":
						ip6.Global = *s
					case "global_address":
						ip6.GlobalAddress = *s
					case "global_address_status":
						ip6.GlobalAddressStatus = *s
					case "default_gateway":
						ip6.DefaultGateway = *s
					}
					ip6Set = true
				}
			}

			if ip6Set {
				iscsi.IPv6Info = ip6
				setAny = true
			}
		}

		if setAny {
			params.IscsiInformationParam = iscsi
		}
	}

	// --------------------
	// NVMe/TCP INFORMATION
	// --------------------
	if _, ok := d.GetOk("nvme_tcp_information"); ok {

		nvme := &gwymodel.NvmeTcpInformationParam{}
		setAny := false

		for _, key := range []string{"vlan_use", "selective_ack", "delayed_ack"} {
			k := "nvme_tcp_information.0." + key
			if b := setBool(k); b != nil {
				switch key {
				case "vlan_use":
					nvme.VlanUse = b
				case "selective_ack":
					nvme.SelectiveAck = b
				case "delayed_ack":
					nvme.DelayedAck = b
				}
				setAny = true
			}
		}

		for _, key := range []string{
			"add_vlan_id", "delete_vlan_id", "tcp_port", "discovery_tcp_port",
		} {
			k := "nvme_tcp_information.0." + key
			if i := setInt(k); i != nil {
				switch key {
				case "add_vlan_id":
					nvme.AddVlanID = i
				case "delete_vlan_id":
					nvme.DeleteVlanID = i
				case "tcp_port":
					nvme.TCPPort = i
				case "discovery_tcp_port":
					nvme.DiscoveryTCPPort = i
				}
				setAny = true
			}
		}

		for _, key := range []string{"ip_mode", "window_size", "mtu_size"} {
			k := "nvme_tcp_information.0." + key
			if s := setString(k); s != nil {
				switch key {
				case "ip_mode":
					nvme.IPMode = s
				case "window_size":
					nvme.WindowSize = s
				case "mtu_size":
					nvme.MTUSize = s
				}
				setAny = true
			}
		}

		// IPv4 nested
		if _, ok := d.GetOk("nvme_tcp_information.0.ipv4_information"); ok {
			ip4 := &gwymodel.IPv4Information{}
			ip4Set := false

			for _, key := range []string{"address", "subnet_mask", "default_gateway"} {
				k := "nvme_tcp_information.0.ipv4_information.0." + key
				if s := setString(k); s != nil {
					switch key {
					case "address":
						ip4.Address = *s
					case "subnet_mask":
						ip4.SubnetMask = *s
					case "default_gateway":
						ip4.DefaultGateway = *s
					}
					ip4Set = true
				}
			}

			if ip4Set {
				nvme.IPv4Info = ip4
				setAny = true
			}
		}

		// IPv6 nested
		if _, ok := d.GetOk("nvme_tcp_information.0.ipv6_information"); ok {

			ip6 := &gwymodel.IPv6Information{}
			ip6Set := false

			for _, key := range []string{
				"linklocal", "linklocal_address", "linklocal_address_status",
				"global", "global_address", "global_address_status",
				"default_gateway",
			} {
				k := "nvme_tcp_information.0.ipv6_information.0." + key
				if s := setString(k); s != nil {
					switch key {
					case "linklocal":
						ip6.Linklocal = *s
					case "linklocal_address":
						ip6.LinklocalAddress = *s
					case "linklocal_address_status":
						ip6.LinklocalAddressStatus = *s
					case "global":
						ip6.Global = *s
					case "global_address":
						ip6.GlobalAddress = *s
					case "global_address_status":
						ip6.GlobalAddressStatus = *s
					case "default_gateway":
						ip6.DefaultGateway = *s
					}
					ip6Set = true
				}
			}

			if ip6Set {
				nvme.IPv6Info = ip6
				setAny = true
			}
		}

		if setAny {
			params.NvmeTcpInformationParam = nvme
		}
	}

	return params
}

func getPortProvisionerManager(serial int) (provmanager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := provmodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	provObj, err := provimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get provisioner manager: %v", err)
		return nil, fmt.Errorf("failed to get provisioner manager: %w", err)
	}

	return provObj, nil
}

func convertOnePortInfoToSchema(p *gwymodel.PortInfo) map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	m := map[string]interface{}{
		"id":                p.ID,
		"protocol":          p.Protocol,
		"port_wwn":          p.PortWwn,
		"port_iscsi_name":   p.PortIscsiName,
		"port_speed":        p.PortSpeed,
		"actual_port_speed": p.ActualPortSpeed,
		"port_security":     p.PortSecurity,
	}

	// Handle FC Information if present
	if p.FcInformation != nil {
		fcInfo := []map[string]interface{}{{
			"al_pa":                  p.FcInformation.AlPa,
			"fabric_switch_setting":  p.FcInformation.FabricSwitchSetting,
			"connection_type":        p.FcInformation.ConnectionType,
			"sfp_data_transfer_rate": p.FcInformation.SfpDataTransferRate,
			"port_mode":              p.FcInformation.PortMode,
		}}
		m["fc_information"] = fcInfo
	}

	// Handle iSCSI Information if present
	if p.IscsiInformation != nil {
		iscsiInfo := map[string]interface{}{
			"vlan_use":               p.IscsiInformation.VlanUse,
			"vlan_id":                p.IscsiInformation.VlanId,
			"ip_mode":                p.IscsiInformation.IpMode,
			"is_ipv6_updating":       p.IscsiInformation.IsIpv6Updating,
			"selective_ack":          p.IscsiInformation.SelectiveAck,
			"delayed_ack":            p.IscsiInformation.DelayedAck,
			"mtu_size":               p.IscsiInformation.MtuSize,
			"link_mtu_size":          p.IscsiInformation.LinkMtuSize,
			"virtual_port_enabled":   p.IscsiInformation.VirtualPortEnabled,
			"tcp_port":               p.IscsiInformation.TcpPort,
			"window_size":            p.IscsiInformation.WindowSize,
			"keep_alive_timer":       p.IscsiInformation.KeepAliveTimer,
			"isns_server_mode":       p.IscsiInformation.IsnsServerMode,
			"isns_server_ip_address": p.IscsiInformation.IsnsServerIpAddress,
			"isns_server_port":       p.IscsiInformation.IsnsServerPort,
		}

		// Handle IPv4 Information if present
		if p.IscsiInformation.Ipv4Information != nil {
			ipv4Info := []map[string]interface{}{{
				"address":         p.IscsiInformation.Ipv4Information.Address,
				"subnet_mask":     p.IscsiInformation.Ipv4Information.SubnetMask,
				"default_gateway": p.IscsiInformation.Ipv4Information.DefaultGateway,
			}}
			iscsiInfo["ipv4_information"] = ipv4Info
		}

		// Handle IPv6 Information if present
		if p.IscsiInformation.Ipv6Information != nil {
			ipv6Info := []map[string]interface{}{{
				"linklocal":                p.IscsiInformation.Ipv6Information.Linklocal,
				"linklocal_address":        p.IscsiInformation.Ipv6Information.LinklocalAddress,
				"linklocal_address_status": p.IscsiInformation.Ipv6Information.LinklocalAddressStatus,
				"global":                   p.IscsiInformation.Ipv6Information.Global,
				"global_address":           p.IscsiInformation.Ipv6Information.GlobalAddress,
				"global_address_status":    p.IscsiInformation.Ipv6Information.GlobalAddressStatus,
				"default_gateway":          p.IscsiInformation.Ipv6Information.DefaultGateway,
			}}
			iscsiInfo["ipv6_information"] = ipv6Info
		}

		m["iscsi_information"] = []map[string]interface{}{iscsiInfo}
	}

	// Handle NVMe-TCP Information if present
	if p.NvmeTcpInformation != nil {
		nvmeInfo := map[string]interface{}{
			"vlan_use":             p.NvmeTcpInformation.VlanUse,
			"vlan_id":              p.NvmeTcpInformation.VlanId,
			"ip_mode":              p.NvmeTcpInformation.IpMode,
			"is_ipv6_updating":     p.NvmeTcpInformation.IsIpv6Updating,
			"selective_ack":        p.NvmeTcpInformation.SelectiveAck,
			"delayed_ack":          p.NvmeTcpInformation.DelayedAck,
			"mtu_size":             p.NvmeTcpInformation.MtuSize,
			"link_mtu_size":        p.NvmeTcpInformation.LinkMtuSize,
			"virtual_port_enabled": p.NvmeTcpInformation.VirtualPortEnabled,
			"tcp_port":             p.NvmeTcpInformation.TcpPort,
			"discovery_tcp_port":   p.NvmeTcpInformation.DiscoveryTcpPort,
			"window_size":          p.NvmeTcpInformation.WindowSize,
		}

		// Handle IPv4 Information if present
		if p.NvmeTcpInformation.Ipv4Information != nil {
			ipv4Info := []map[string]interface{}{{
				"address":         p.NvmeTcpInformation.Ipv4Information.Address,
				"subnet_mask":     p.NvmeTcpInformation.Ipv4Information.SubnetMask,
				"default_gateway": p.NvmeTcpInformation.Ipv4Information.DefaultGateway,
			}}
			nvmeInfo["ipv4_information"] = ipv4Info
		}

		// Handle IPv6 Information if present
		if p.NvmeTcpInformation.Ipv6Information != nil {
			ipv6Info := []map[string]interface{}{{
				"linklocal":                p.NvmeTcpInformation.Ipv6Information.Linklocal,
				"linklocal_address":        p.NvmeTcpInformation.Ipv6Information.LinklocalAddress,
				"linklocal_address_status": p.NvmeTcpInformation.Ipv6Information.LinklocalAddressStatus,
				"global":                   p.NvmeTcpInformation.Ipv6Information.Global,
				"global_address":           p.NvmeTcpInformation.Ipv6Information.GlobalAddress,
				"global_address_status":    p.NvmeTcpInformation.Ipv6Information.GlobalAddressStatus,
				"default_gateway":          p.NvmeTcpInformation.Ipv6Information.DefaultGateway,
			}}
			nvmeInfo["ipv6_information"] = ipv6Info
		}

		m["nvme_tcp_information"] = []map[string]interface{}{nvmeInfo}
	}

	log.WriteDebug("Convert: %+v", m)
	return m
}

func convertMultiplePortInfosListToSchema(ports *gwymodel.PortInfoList) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Defensive check - now using Data field from API response
	if ports == nil || len(ports.Data) == 0 {
		return nil
	}

	portList := make([]map[string]interface{}, len(ports.Data))
	for i, p := range ports.Data {
		m := convertOnePortInfoToSchema(&p)
		portList[i] = m
	}

	return portList
}
