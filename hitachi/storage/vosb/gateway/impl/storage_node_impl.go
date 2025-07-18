package vssbstorage

import (

	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vosb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// Add Storage Node
func (psm *vssbStorageManager) AddStorageNode(
	configurationFile string,
	exportedConfigurationFile string,
	setupUserPassword string,
	expectedCloudProvider string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	err = psm.doAddStorageNode(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider)
	if err != nil {
		log.WriteDebug("TFError| error in AddStorageNode, err: %v", err)
	}

	return
}

// GetStorageNodes gets nodes information of vssb storage
func (psm *vssbStorageManager) GetStorageNodes() (*vssbmodel.StorageNodes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageNodes vssbmodel.StorageNodes
	apiSuf := "objects/storage-nodes"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storageNodes)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storageNodes, nil
}

// GetStorageNodes gets node information for a specific node of vssb storage
func (psm *vssbStorageManager) GetStorageNode(id string) (*vssbmodel.StorageNode, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storageNode vssbmodel.StorageNode
	apiSuf := fmt.Sprintf("objects/storage-nodes/%s", id)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storageNode)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &storageNode, nil
}

// validateStorageNodeParameters validates the configuration parameters based on cloud provider
func validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider string) error {
	if expectedCloudProvider != "baremetal" {
		// For non-baremetal cloud providers: exported_configuration_file is required, others should not be given
		if exportedConfigurationFile == "" {
			return fmt.Errorf("exported_configuration_file is required when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
		if setupUserPassword != "" {
			return fmt.Errorf("setup_user_password should not be provided when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
		if configurationFile != "" {
			return fmt.Errorf("configuration_file should not be provided when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
	} else {
		// For baremetal: exported_configuration_file must not be given, others are required
		if exportedConfigurationFile != "" {
			return fmt.Errorf("exported_configuration_file should not be provided when expected_cloud_provider is 'baremetal'")
		}
		if setupUserPassword == "" {
			return fmt.Errorf("setup_user_password is required when expected_cloud_provider is 'baremetal'")
		}
		if configurationFile == "" {
			return fmt.Errorf("configuration_file is required when expected_cloud_provider is 'baremetal'")
		}
	}
	return nil
}

func (psm *vssbStorageManager) doAddStorageNode(
	configurationFile string,
	exportedConfigurationFile string,
	setupUserPassword string,
	expectedCloudProvider string) (err error){

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	
	// Validate parameter combinations based on cloud provider
	err = validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider)
	if err != nil {
		log.WriteError(err)
		return
	}
	
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	
	// Add expected cloud provider field
	// cloudProviderField, err := writer.CreateFormField("expectedCloudProvider")
	// if err != nil {
	// 	log.WriteError(err)
	// 	return
	// }
	// cloudProviderField.Write([]byte(expectedCloudProvider))

	if expectedCloudProvider != "baremetal" {
		// For non-baremetal cloud providers: only add exported configuration file
		if exportedConfigurationFile != "" {
			var exportedConfigField io.Writer
			exportedConfigField, err = writer.CreateFormFile("exportedConfigurationFile", filepath.Base(exportedConfigurationFile))
			if err != nil {
				log.WriteError(err)
				return err
			}
			var exportedFd *os.File
			exportedFd, err = os.Open(exportedConfigurationFile)
			if err != nil {
				log.WriteError(err)
				return err
			}
			defer exportedFd.Close()
			_, err = io.Copy(exportedConfigField, exportedFd)
			if err != nil {
				log.WriteError(err)
				return err
			}
		}
	} else {
		// For baremetal: add setupUserPassword and configurationFile
		if setupUserPassword != "" {
			var formField io.Writer
			formField, err = writer.CreateFormField("setupUserPassword")
			if err != nil {
				log.WriteError(err)
				return err
			}
			formField.Write([]byte(setupUserPassword))
		}

		if configurationFile != "" {
			var fw io.Writer
			fw, err = writer.CreateFormFile("configurationFile", filepath.Base(configurationFile))
			if err != nil {
				log.WriteError(err)
				return err
			}
			var fd *os.File
			fd, err = os.Open(configurationFile)
			if err != nil {
				log.WriteError(err)
				return err
			}
			defer fd.Close()
			_, err = io.Copy(fw, fd)
			if err != nil {
				log.WriteError(err)
				return err
			}
		}
	}

	writer.Close()

	////////////////////////////////////////////////////////////////

	psm.storageSetting.ContentType = writer.FormDataContentType()

	apiSuf := "objects/storage-nodes"
	affRes, err := httpmethod.PostCallFormExt(psm.storageSetting, apiSuf, form)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v\n", apiSuf, err)
		return
	}
	log.WriteDebug("affRes= %v", affRes)
	return

}
