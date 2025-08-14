package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"

	"github.com/jinzhu/copier"
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

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return
	}

	err = gatewayObj.AddStorageNode(configurationFile, exportedConfigurationFile, setupUserPassword, expectedCloudProvider)
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

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_BEGIN))
	storageNodes, err := gatewayObj.GetStorageNodes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStorageNodes, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_NODES_FAILED))
		return nil, err
	}

	provStorageNodes := vssbmodel.StorageNodes{}
	err = copier.Copy(&provStorageNodes, storageNodes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_END))

	return &provStorageNodes, nil
}

// GetStorageNode gets node information for a specific node of vssb storage
func (psm *vssbStorageManager) GetStorageNode(nodeId string) (*vssbmodel.StorageNode, error) {
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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_BEGIN), nodeId)
	storageNode, err := gatewayObj.GetStorageNode(nodeId)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStorageNode, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_NODE_FAILED), nodeId)
		return nil, err
	}

	provStorageNode := vssbmodel.StorageNode{}
	err = copier.Copy(&provStorageNode, storageNode)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_END), nodeId)

	return &provStorageNode, nil
}
