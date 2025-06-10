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

	finalPath, err := reconObj.ReconcileConfigurationDefinitionFile(doCreate, doDownload, downloadPath)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_DOWNLOAD_CONFIG_FAILED))
		log.WriteDebug("TFError| error in Restoring configuration definitionfile, err: %v", err)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_DOWNLOAD_CONFIG_END))
	return finalPath, nil
}
