package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetAllVolumes gets all available volume details
func (psm *vssbStorageManager) GetAllVolumes(computeNodeName string) (*vssbmodel.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
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
	provVolumes, err := provObj.GetAllVolumes(computeNodeName)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllVolumes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerVolumes := vssbmodel.Volumes{}
	err = copier.Copy(&reconcilerVolumes, provVolumes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))
	return &reconcilerVolumes, nil
}

// GetVolumeDetail based on volume name
func (psm *vssbStorageManager) GetVolumeDetails(volumeName string) (*vssbmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
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
	provNodes, err := provObj.GetVolumeDetails(volumeName)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeDetails provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcileNodes := vssbmodel.Volume{}
	err = copier.Copy(&reconcileNodes, provNodes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))
	return &reconcileNodes, nil
}

// ReconcileVolume
func (psm *vssbStorageManager) ReconcileVolume(postData *vssbmodel.CreateVolume) (*vssbmodel.Volume, error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_BEGIN), *postData.Name)
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
	provVolume, err := provObj.GetVolumeDetailsByName(*postData.Name)
	if err != nil {
		log.WriteInfo("TFError| No Volume found so creating a new volume, err: %v", err)

	}
	if provVolume == nil {
		var nickname string = ""
		if *postData.NickName == "" {
			nickname = *postData.Name
		} else {
			nickname = *postData.NickName
		}
		poolDetails, err := provObj.GetStoragePoolByPoolName(*postData.PoolName)
		if err != nil {
			log.WriteDebug("TFError| error in GetVolumeDetails provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
			return nil, err
		}
		poolId := poolDetails.ID
		volumeAdd, err := provObj.CreateVolume(*postData.Name, nickname, poolId, *postData.CapacityInGB)
		if err != nil {
			log.WriteDebug("TFError| error in GetVolumeDetails provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
			return nil, err
		}
		log.WriteDebug("Volume created on this ID: %v", *volumeAdd)
	} else {
		if *postData.NickName != "" {
			if provVolume.NickName != *postData.NickName {
				err := provObj.UpdateVolumeNickName(provVolume.ID, *postData.NickName)
				if err != nil {
					log.WriteDebug("TFError| error in UpdateVolumeNickName provisioner call, err: %v", err)
					log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
					return nil, err
				}
			}
		}

		var capacityMB int32 = int32(*postData.CapacityInGB * 1024)
		if int32(provVolume.TotalCapacity) < capacityMB {
			additionalSize := capacityMB - int32(provVolume.TotalCapacity)
			err := provObj.ExpandVolume(provVolume.ID, &additionalSize)
			if err != nil {
				log.WriteDebug("TFError| error in ExpandVolume provisioner call, err: %v", err)
				log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
				return nil, err
			}

		}
	}
	provNodesAfterCreate, err := provObj.GetVolumeDetails(*postData.Name)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeDetails provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
		return nil, err
	}
	if postData.ComputeNodes != nil {

		err = psm.AddRemoveVolumeToComputeNodes(provNodesAfterCreate, &postData.ComputeNodes)
	}
	if err != nil {
		log.WriteDebug("TFError| error in Add compute node psm call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
		return nil, err
	}

	//Fetch the final volume details after creating/updating the the volume
	provNodesAfterAddCN, err := provObj.GetVolumeDetails(*postData.Name)
	if err != nil {
		log.WriteDebug("TFError| error in GetVolumeDetails provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcileNodes := vssbmodel.Volume{}
	err = copier.Copy(&reconcileNodes, provNodesAfterAddCN)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), *postData.Name)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_END), *postData.Name)
	return &reconcileNodes, nil
}

func (psm *vssbStorageManager) AddRemoveVolumeToComputeNodes(provNode *provisonermodel.Volume, computeNodeNames *[]string) error {
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
		return err
	}
	computeNodeSlices := []string{}
	for _, node := range provNode.ComputeNodes {
		computeNodeSlices = append(computeNodeSlices, node.Nickname)
	}
	added, removed, _, _ := utils.GetStringSliceDiff(computeNodeSlices, *computeNodeNames)
	//Aattching compute nodes to the volume

	for _, computeNode := range added {
		log.WriteInfo(mc.GetMessage(mc.INFO_ADD_VOLUME_TO_COMPUTE_NODE_BEGIN), computeNode)
		serverID, err := provObj.GetComputeNodeIdByName(computeNode)
		if err != nil {
			log.WriteDebug("TFError| error in getcomputeNodename provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_ADD_VOLUME_TO_COMPUTE_NODE_FAILED), computeNode)
			return err
		}
		res, err := provObj.AddVolumeToComputeNode(provNode.ID, serverID)
		if err != nil {
			log.WriteDebug("TFError| error in AddVolumeToComputeNode provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_ADD_VOLUME_TO_COMPUTE_NODE_FAILED), computeNode)
			return err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_ADD_VOLUME_TO_COMPUTE_NODE_END))

		log.WriteDebug("Volume added on this ID: %v", res)
	}
	//Deattaching the computenodes from the volume

	for _, computeNode := range removed {
		log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_VOLUME_FROM_COMPUTE_NODE_BEGIN), computeNode)
		serverID, err := provObj.GetComputeNodeIdByName(computeNode)
		if err != nil {
			log.WriteDebug("TFError| error in getcomputeNodename provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_REMOVE_VOLUME_FROM_COMPUTE_NODE_FAILED), computeNode)
			return err
		}
		err = provObj.RemoveVolumeFromComputeNode(provNode.ID, serverID)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteVolumeToComputeNode provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_REMOVE_VOLUME_FROM_COMPUTE_NODE_FAILED), computeNode)
			return err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_VOLUME_FROM_COMPUTE_NODE_END), computeNode)

		log.WriteDebug("Volume removed on this ID: %v", computeNode)
	}

	return nil
}

func (psm *vssbStorageManager) DeleteVolumeResource(volumeID *string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_BEGIN), volumeID)
	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	err = provObj.DeleteVolume(*volumeID)
	if err != nil {
		log.WriteDebug("TFError| error in Deletevolume call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED), volumeID)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_END), volumeID)
	return nil
}
