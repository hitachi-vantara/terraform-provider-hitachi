package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/message-catalog"
)

func (psm *vssbStorageManager) RestoreConfigurationDefinitionFile() error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_CONFIG_BEGIN))

	err = gatewayObj.RestoreConfigurationDefinitionFile()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_RESTORE_CONFIG_FAILED))
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_CONFIG_END))

	return nil
}

func (psm *vssbStorageManager) DownloadConfigurationFile(toFilePath string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("error in NewEx call, err: %v", err)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DOWNLOAD_CONFIG_BEGIN))

	filePath, err := gatewayObj.DownloadConfigurationFile(toFilePath)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DOWNLOAD_CONFIG_FAILED))
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DOWNLOAD_CONFIG_END))

	return filePath, nil
}
