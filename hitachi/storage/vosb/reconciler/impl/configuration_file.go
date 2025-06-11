package vssbstorage

import (
	"os"
	"path/filepath"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisionerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisionermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
)

func (psm *vssbStorageManager) RestoreConfigurationDefinitionFile(createConfigParam *reconcilermodel.CreateConfigurationFileParam) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisionermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provisionerObj, err := provisionerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_RESTORE_CONFIG_BEGIN))

	provConfigParam := provisionermodel.CreateConfigurationFileParam{}
	err = copier.Copy(&provConfigParam, createConfigParam)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to provisioner structure, err: %v", err)
		return err
	}

	err = provisionerObj.RestoreConfigurationDefinitionFile(&provConfigParam)
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

	objStorage := provisionermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provisionerObj, err := provisionerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("error in NewEx call, err: %v", err)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DOWNLOAD_CONFIG_BEGIN))

	filePath, err := provisionerObj.DownloadConfigurationFile(toFilePath)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DOWNLOAD_CONFIG_FAILED))
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DOWNLOAD_CONFIG_END))

	return filePath, nil
}

func (psm *vssbStorageManager) ReconcileConfigurationDefinitionFile(doCreate bool, doDownload bool, downloadPath string, createConfigParam *reconcilermodel.CreateConfigurationFileParam) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if doCreate {
		err := psm.RestoreConfigurationDefinitionFile(createConfigParam)
		if err != nil {
			return "", err
		}
	}

	if doDownload {
		if filepath.Ext(downloadPath) == "" && !isDir(downloadPath) {
			downloadPath += ".tar.gz"
		}

		finalPath, err := psm.DownloadConfigurationFile(downloadPath)
		if err != nil {
			return "", err
		}
		return finalPath, nil
	}

	return "", nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
