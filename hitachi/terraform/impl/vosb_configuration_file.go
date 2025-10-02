package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
)

func CreateDownloadConfigurationDefinitionFile(d *schema.ResourceData) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

	downloadOnly := d.Get("download_existconfig_only").(bool)
	createOnly := d.Get("create_only").(bool)
	downloadPath := d.Get("download_path").(string)

	doDownload := true
	doCreate := true

	if downloadOnly && createOnly {
		log.WriteDebug("`create_only` is ignored when `download_existconfig_only` is true\n")
		doDownload = true
		doCreate = false
	}

	if !downloadOnly && !createOnly {
		doDownload = true
		doCreate = true
	}

	if !downloadOnly && createOnly {
		doDownload = false
		doCreate = true
	}

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return "", err
	}
	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_DOWNLOAD_CONFIG_BEGIN))

	createConfigParam := expandCreateConfigurationFileParam(d)

	finalPath, err := reconObj.ReconcileConfigurationDefinitionFile(doCreate, doDownload, downloadPath, createConfigParam)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_DOWNLOAD_CONFIG_FAILED))
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_DOWNLOAD_CONFIG_END))
	return finalPath, nil
}

func expandAddressSetting(list []interface{}) []reconcilermodel.AddressSetting {
	var result []reconcilermodel.AddressSetting
	for _, v := range list {
		m := v.(map[string]interface{})
		addr := reconcilermodel.AddressSetting{
			Index:                    m["index"].(int),
			ControlPortIPv4Address:   m["control_port_ipv4_address"].(string),
			InternodePortIPv4Address: m["internode_port_ipv4_address"].(string),
			ComputePortIPv4Address:   m["compute_port_ipv4_address"].(string),
		}
		if v6, ok := m["compute_port_ipv6_address"]; ok && v6 != nil {
			addr.ComputePortIPv6Address = v6.(string)
		}
		result = append(result, addr)
	}
	return result
}

func expandCreateConfigurationFileParam(d *schema.ResourceData) *reconcilermodel.CreateConfigurationFileParam {
	raw := d.Get("create_configuration_file_param").([]interface{})
	if len(raw) == 0 || raw[0] == nil {
		return &reconcilermodel.CreateConfigurationFileParam{}
	}
	m := raw[0].(map[string]interface{})
	param := &reconcilermodel.CreateConfigurationFileParam{
		ExportFileType:        m["export_file_type"].(string),
		MachineImageID:        m["machine_image_id"].(string),
		NumberOfDrives:        m["number_of_drives"].(int),
		RecoverSingleDrive:    m["recover_single_drive"].(bool),
		DriveID:               m["drive_id"].(string),
		RecoverSingleNode:     m["recover_single_node"].(bool),
		NodeID:                m["node_id"].(string),
		TemplateS3Url:         m["template_s3_url"].(string),
	}
	if v, ok := m["address_setting"]; ok && v != nil {
		param.AddressSetting = expandAddressSetting(v.([]interface{}))
	}
	return param
}
