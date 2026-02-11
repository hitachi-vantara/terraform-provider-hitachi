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
	expectedCloudProvider string,
	vmConfigFileS3URI string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if expectedCloudProvider == "azure" {
		err = psm.doAddStorageNodeAzure(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
		if err != nil {
			log.WriteDebug("TFError| error from doAddStorageNodeAzure, err: %v", err)
		}
		return
	} else if expectedCloudProvider == "google" {
		err = psm.doAddStorageNodeGoogle(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
		if err != nil {
			log.WriteDebug("TFError| error from doAddStorageNodeGoogle, err: %v", err)
		}
		return
	}

	err = psm.doAddStorageNode(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
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
func validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI string) error {
	if expectedCloudProvider == "aws" {
		// For AWS: configuration_file and vm_configuration_file_s3_uri are required, others should not be given
		if configurationFile == "" {
			return fmt.Errorf("configuration_file is required when expected_cloud_provider is 'aws'")
		}
		if vmConfigFileS3URI == "" {
			return fmt.Errorf("vm_configuration_file_s3_uri is required when expected_cloud_provider is 'aws'")
		}
		if setupUserPassword != "" {
			return fmt.Errorf("setup_user_password should not be provided when expected_cloud_provider is 'aws'")
		}
		if exportedConfigurationFile != "" {
			return fmt.Errorf("exported_configuration_file should not be provided when expected_cloud_provider is 'aws'")
		}
	} else if expectedCloudProvider == "azure" {
		if exportedConfigurationFile == "" {
			return fmt.Errorf("exported_configuration_file is required when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
		if setupUserPassword != "" {
			return fmt.Errorf("setup_user_password should not be provided when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
		if configurationFile != "" {
			return fmt.Errorf("configuration_file should not be provided when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
		if vmConfigFileS3URI != "" {
			return fmt.Errorf("vm_configuration_file_s3_uri should not be provided when expected_cloud_provider is '%s'", expectedCloudProvider)
		}
	} else if expectedCloudProvider == "baremetal" {
		// For baremetal: exported_configuration_file and vm_configuration_file_s3_uri must not be given, others are required
		if exportedConfigurationFile != "" {
			return fmt.Errorf("exported_configuration_file should not be provided when expected_cloud_provider is 'baremetal'")
		}
		if vmConfigFileS3URI != "" {
			return fmt.Errorf("vm_configuration_file_s3_uri should not be provided when expected_cloud_provider is 'baremetal'")
		}
		if setupUserPassword == "" {
			return fmt.Errorf("setup_user_password is required when expected_cloud_provider is 'baremetal'")
		}
		if configurationFile == "" {
			return fmt.Errorf("configuration_file is required when expected_cloud_provider is 'baremetal'")
		}
	} else 	if expectedCloudProvider == "google" {
		// For GPC: no additional parameters are required
		if configurationFile != "" {
			return fmt.Errorf("configuration_file should not be provided when expected_cloud_provider is 'google'")
		}
		if exportedConfigurationFile != "" {
			return fmt.Errorf("exported_configuration_file should not be provided when expected_cloud_provider is 'google'")
		}
		if setupUserPassword != "" {
			return fmt.Errorf("setup_user_password should not be provided when expected_cloud_provider is 'google'")
		}
		if vmConfigFileS3URI != "" {
			return fmt.Errorf("vm_configuration_file_s3_uri should not be provided when expected_cloud_provider is 'google'")
		}
	}

	return nil
}

func (psm *vssbStorageManager) doAddStorageNode(
	configurationFile string,
	exportedConfigurationFile string,
	setupUserPassword string,
	expectedCloudProvider string,
	vmConfigFileS3URI string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Validate parameter combinations based on cloud provider
	err = validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
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

	if expectedCloudProvider == "aws" {
		// For AWS: add configuration file and S3 URI
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
		if vmConfigFileS3URI != "" {
			var s3URIField io.Writer
			s3URIField, err = writer.CreateFormField("vmConfigurationFileS3Uri")
			if err != nil {
				log.WriteError(err)
				return err
			}
			s3URIField.Write([]byte(vmConfigFileS3URI))
		}
	} else if expectedCloudProvider == "baremetal" {
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

func (psm *vssbStorageManager) doAddStorageNodeAzure(
	configurationFile string,
	exportedConfigurationFile string,
	setupUserPassword string,
	expectedCloudProvider string,
	vmConfigFileS3URI string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Validate parameter combinations based on cloud provider
	err = validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
	if err != nil {
		log.WriteError(err)
		return
	}

	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)

	if expectedCloudProvider == "azure" {
		var exportedConfigField io.Writer
		exportedConfigField, err = writer.CreateFormFile("exportedConfigurationFile", filepath.Base(exportedConfigurationFile))
		if err != nil {
			log.WriteError(err)
			return err
		}
		var exportedFd *os.File
		exportedFd, err = os.OpenFile(exportedConfigurationFile, os.O_RDONLY, 0)
		if err != nil {
			log.WriteError(err)
			return err
		}
		defer exportedFd.Close()

		// Read binary file in 32KB chunks to handle large files efficiently
		buffer := make([]byte, 32*1024)
		for {
			n, readErr := exportedFd.Read(buffer)
			if n > 0 {
				_, writeErr := exportedConfigField.Write(buffer[:n])
				if writeErr != nil {
					log.WriteError(writeErr)
					return writeErr
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				log.WriteError(readErr)
				return readErr
			}
		}
	}

	writer.Close()

	////////////////////////////////////////////////////////////////

	log.WriteDebug("Binary form data prepared, size: %d bytes", form.Len())

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

func (psm *vssbStorageManager) doAddStorageNodeGoogle(
	configurationFile string,
	exportedConfigurationFile string,
	setupUserPassword string,
	expectedCloudProvider string,
	vmConfigFileS3URI string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Validate parameter combinations based on cloud provider
	err = validateStorageNodeParameters(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
	if err != nil {
		log.WriteError(err)
		return
	}

	////////////////////////////////////////////////////////////////
	// For Google: No parameters, empty body, just POST with Content-Length: 0
	////////////////////////////////////////////////////////////////

	psm.storageSetting.ContentType = "application/json" // could also be empty, but safe default

	apiSuf := "objects/storage-nodes"

	affRes, err := httpmethod.PostCallNoBodyExt(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v\n", apiSuf, err)
		return
	}

	log.WriteDebug("affRes= %v", affRes)
	return
}
