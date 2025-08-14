package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/message-catalog"
	provisionermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"

	"github.com/jinzhu/copier"
)

func (psm *vssbStorageManager) RestoreConfigurationDefinitionFile(createConfigParam *provisionermodel.CreateConfigurationFileParam) error {
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

	gwyConfigParam := vssbgatewaymodel.CreateConfigurationFileParam{}
	err = copier.Copy(&gwyConfigParam, createConfigParam)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from provisioner to gateway structure, err: %v", err)
		return err
	}

	err = gatewayObj.RestoreConfigurationDefinitionFile(&gwyConfigParam)
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
