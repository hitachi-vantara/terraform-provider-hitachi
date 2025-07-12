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
	setupUserPassword string) (err error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	err = psm.doAddStorageNode(configurationFile, setupUserPassword)
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

func (psm *vssbStorageManager) doAddStorageNode(
	configurationFile string,
	setupUserPassword string) (err error){

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	formField, err := writer.CreateFormField("setupUserPassword")
	if err != nil {
		log.WriteError(err)
		return
	}
	formField.Write([]byte(setupUserPassword))

	fw, err := writer.CreateFormFile("configurationFile", filepath.Base(configurationFile))
	if err != nil {
		log.WriteError(err)
		return
	}
	fd, err := os.Open(configurationFile)
	if err != nil {
		log.WriteError(err)
		return
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.WriteError(err)
		return
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
