package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
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

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return
	}

	err = provObj.AddStorageNode(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider, vmConfigFileS3URI)
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_NODES_FAILED))
		return
	}

	return
}

// GetStorageNodes gets nodes information of vssb storage
func (psm *vssbStorageManager) GetStorageNodes() (*vssbmodel.StorageNodes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_BEGIN))
	provStorageNodes, err := provObj.GetStorageNodes()
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_NODES_FAILED))
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStorageNodes := vssbmodel.StorageNodes{}
	err = copier.Copy(&reconStorageNodes, provStorageNodes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_END))

	return &reconStorageNodes, nil
}

// GetStorageNode gets node information of vssb storage with node auth settings
func (psm *vssbStorageManager) GetStorageNode(nodeName string) (*vssbmodel.StorageNode, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_BEGIN), nodeName)
	allStorageNodes, err := psm.GetStorageNodes()
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_NODE_FAILED), nodeName)
		return nil, err
	}

	nodeId := ""
	log.WriteDebug("20250612 nodes: %+v", *allStorageNodes)
	for _, node := range allStorageNodes.Data {
		log.WriteDebug("node: %v", node)
		if strings.EqualFold(nodeName, node.Name) {
			nodeId = node.ID
		}
	}
	if nodeId == "" {
		log.WriteDebug("TFError| error getting specified node id from list of all nodes")
		log.WriteError(mc.GetMessage(mc.ERR_GET_NODE_FAILED), nodeName)
		return nil, fmt.Errorf("The specified node name could not be found.")
	}
	provnode, err := provObj.GetStorageNode(nodeId)
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_NODE_FAILED), nodeId)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconnode := vssbmodel.StorageNode{}
	err = copier.Copy(&reconnode, provnode)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_END), nodeId)

	return &reconnode, nil
}
