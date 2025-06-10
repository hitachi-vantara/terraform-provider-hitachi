package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vosb/gateway/http"
)

// Restores a configuration definition file.
func (psm *vssbStorageManager) RestoreConfigurationDefinitionFile() error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/configuration-file/actions/create/invoke"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		return err
	}
	return nil
}

// Downloads a restored configuration definition file.
// You can run this API only for the cluster master node (primary).
func (psm *vssbStorageManager) DownloadConfigurationFile(toFilePath string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/configuration-file/download"

	filePath, err := httpmethod.DownloadFileCall(psm.storageSetting, apiSuf, toFilePath)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteInfo("Configuration file saved at: %s", filePath)
	return filePath, nil
}
