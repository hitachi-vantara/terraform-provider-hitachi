package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetComputeNode gets compute node server
func (psm *vssbStorageManager) GetComputeNode(serverID string) (*vssbmodel.ComputeNodeWithPathDetails, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_BEGIN), serverID)
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
	provServer, err := provObj.GetComputeNode(serverID)
	if err != nil {
		log.WriteDebug("TFError| error in GetComputeNode provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_FAILED), serverID)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerServer := vssbmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&reconcilerServer, provServer)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_END), serverID)
	return &reconcilerServer, nil
}

// GetAllComputeNodes
func (psm *vssbStorageManager) GetAllComputeNodes() (*vssbmodel.Servers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_BEGIN))
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
	provServers, err := provObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllComputeNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_SERVERS_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerServers := vssbmodel.Servers{}
	err = copier.Copy(&reconcilerServers, provServers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_END))
	return &reconcilerServers, nil
}

// ReconcileComputeNode .
func (psm *vssbStorageManager) ReconcileComputeNode(inputCompute *vssbmodel.ComputeResource) (*vssbmodel.ComputeNodeWithPathDetails, error) {
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
	existingResource, err := psm.GetComputeNodeInformationByName(inputCompute.Name, inputCompute.ID)
	if err != nil {
		log.WriteDebug("TFError| error in GetComputeNodeInformationByName provisioner call, err: %v", err)
		return nil, err
	}

	// CREATE RESOURCE
	if existingResource == nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_BEGIN), inputCompute.Name)

		if (inputCompute.OsType == "") {
			return nil, fmt.Errorf("OsType is required when adding compute node")
		}

		if (!strings.EqualFold(inputCompute.OsType, "Linux") && !strings.EqualFold(inputCompute.OsType, "VMware") && !strings.EqualFold(inputCompute.OsType, "Windows")) {
			return nil, fmt.Errorf("OsType is invalid")
		}
		// Converting Reconciler to Prov
		provResource := provisonermodel.ComputeResource{}
		err = copier.Copy(&provResource, inputCompute)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to prov structure, err: %v", err)
			return nil, err
		}
		err := provObj.CreateComputeResource(&provResource)
		if err != nil {
			log.WriteDebug("TFError| error in CreateComputeResource provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_SERVER_FAILED), inputCompute.Name)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_END), inputCompute.Name)
	} else {
		// UPDATE RESOURCE
		err := psm.UpdateComputeNodeResource(inputCompute, existingResource)
		if err != nil {
			log.WriteDebug("TFError| error in UpdateComputeNodeResource provisioner call, err: %v", err)
			return nil, err
		}
	}

	// Read resource after all operations
	provisionerResource, err := psm.GetComputeNodeInformationByName(inputCompute.Name, inputCompute.ID)
	if err != nil {
		log.WriteDebug("TFError| error in GetComputeNodeInformationByName provisioner call, err: %v", err)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerResource := vssbmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&reconcilerResource, provisionerResource)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	return &reconcilerResource, nil
}

// GetComputeNodeInformationByName .
func (psm *vssbStorageManager) GetComputeNodeInformationByName(computeName string, id string) (*provisonermodel.ComputeNodeWithPathDetails, error) {
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
	var existingResource *provisonermodel.ComputeNodeWithPathDetails = nil
	if id == "" && computeName != "" {
		serverId, err := provObj.GetComputeNodeIdByName(computeName)
		if err != nil {
			log.WriteDebug("TFError| error in GetComputeNodeIdByName provisioner call, err: %v", err)
			return nil, err
		}

		if serverId != "" {
			existingResource, err = provObj.GetComputeNode(serverId)
			if err != nil {
				log.WriteDebug("TFError| error in GetComputeNode provisioner call, err: %v", err)
				return nil, err
			}
		}
	} else if id != "" {
		existingResource, err = provObj.GetComputeNode(id)
		if err != nil {
			log.WriteDebug("TFError| error in GetComputeNode provisioner call, err: %v", err)
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid name and Id")
	}
	return existingResource, nil
}

// UpdateComputeNodeResource .
func (psm *vssbStorageManager) UpdateComputeNodeResource(inputCompute *vssbmodel.ComputeResource, existingCompute *provisonermodel.ComputeNodeWithPathDetails) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("TFDebug| updating existing compute node %#v", *existingCompute)
	log.WriteDebug("TFDebug| input compute node %#v", *inputCompute)
	err := psm.ReconcileComputeNodeNameAndOsType(inputCompute, existingCompute)
	if err != nil {
		log.WriteDebug("TFError| error in ReconcileComputeNodeNameAndOsType, err: %v", err)
		return err
	}
	// Update connection when available for add, remove, update
	if (inputCompute.IscsiConnections != nil) || (inputCompute.IscsiConnections == nil && len(existingCompute.ComputePaths.Data) > 0) {
		err = psm.ReconcileComputeNodeIscsiConnections(inputCompute, existingCompute)
		if err != nil {
			log.WriteDebug("TFError| error in ReconcileComputeNodeIscsiConnections, err: %v", err)
			return err
		}
	}
	if inputCompute.FcConnections != nil && len(inputCompute.FcConnections) > 0 {
		err := psm.AddRemoveFCPorts(inputCompute, existingCompute)
		if err != nil {
			log.WriteDebug("TFError| error in ReconcileComputeNodeIscsiConnections, err: %v", err)
			return err
		}
	}
	return nil
}

func (psm *vssbStorageManager) AddRemoveFCPorts(inputCompute *vssbmodel.ComputeResource, existingCompute *provisonermodel.ComputeNodeWithPathDetails) error {
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
	existingHbaNames := []string{}
	for _, path := range existingCompute.Node.Paths {
		if "FC" == path.Protocol {
			existingHbaNames = append(existingHbaNames, path.HbaName)
		}
	}
	newHbaNames := []string{}
	for _, hbaname := range inputCompute.FcConnections {
		newHbaNames = append(newHbaNames, hbaname.HostWWN)

	}
	added, removed, _, _ := utils.GetStringSliceDiff(existingHbaNames, newHbaNames)
	log.WriteDebug("TFDebug| added FC %#v", added)
	log.WriteDebug("TFDebug| removed FC %#v", removed)

	err = provObj.AddFCportsWWN(added, existingCompute.Node.ID)
	if err != nil {
		log.WriteDebug("TFError| error in Add of FC ports, err: %v", err)
		return err
	}

	err = provObj.RemoveFCportsWWN(removed, existingCompute.Node.ID)
	if err != nil {
		log.WriteDebug("TFError| error in Remove of FC ports, err: %v", err)
		return err
	}

	return nil

}

// ReconcileComputeNodeNameAndOsType used to reconcile - update computenode for name ans os type
func (psm *vssbStorageManager) ReconcileComputeNodeNameAndOsType(inputCompute *vssbmodel.ComputeResource, existingCompute *provisonermodel.ComputeNodeWithPathDetails) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	var err error = nil
	if !strings.EqualFold(inputCompute.Name, existingCompute.Node.Nickname) || (!strings.EqualFold(inputCompute.OsType, existingCompute.Node.OsType) && inputCompute.OsType != "") {
		provMgr, _ := provisonerimpl.New(psm.storageSetting.Username, psm.storageSetting.Password, psm.storageSetting.ClusterAddress)
		// If Only Name need to update
		if inputCompute.Name != "" && inputCompute.OsType == "" {
			log.WriteDebug("TFDebug| updating only name")
			err = provMgr.UpdateComputeNode(existingCompute.Node.ID, inputCompute.Name, existingCompute.Node.OsType)
		} else {
			// If both Name and OsType need to update
			log.WriteDebug("TFDebug| updating name and os type")
			err = provMgr.UpdateComputeNode(existingCompute.Node.ID, inputCompute.Name, inputCompute.OsType)
		}
	}
	return err
}

func GetComputeNodeIscsiInitiator(iqn string, inputCompute *vssbmodel.ComputeResource) (*vssbmodel.IscsiConnector, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	var computeNodeIscsiInitiator *vssbmodel.IscsiConnector = nil
	for _, connection := range inputCompute.IscsiConnections {
		if strings.EqualFold(iqn, connection.IscsiInitiator) {
			computeNodeIscsiInitiator = &connection
			break
		}
	}
	if computeNodeIscsiInitiator == nil {
		log.WriteDebug("TFDebug| iscsi initiator not found")
		return computeNodeIscsiInitiator, fmt.Errorf("Compute Node iSCSI Initiator IQN %s not found.", iqn)
	} else {
		return computeNodeIscsiInitiator, nil
	}

}

func FindIScsiConnectionPorts(iqn string, inputCompute *vssbmodel.ComputeResource) (ports []string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	ports = []string{}
	for _, connection := range inputCompute.IscsiConnections {
		if strings.EqualFold(iqn, connection.IscsiInitiator) {
			ports = connection.PortNames
		}
	}
	if len(ports) == 0 {
		err = fmt.Errorf("IScsi Initiator IQN %s is not found in the compute node configuration.", iqn)
	}
	return ports, err
}

// ReconcileComputeNodeIscsiConnections .
func (psm *vssbStorageManager) ReconcileComputeNodeIscsiConnections(inputCompute *vssbmodel.ComputeResource, existingCompute *provisonermodel.ComputeNodeWithPathDetails) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	var newInitiatorIqns []string
	for _, connection := range inputCompute.IscsiConnections {
		newInitiatorIqns = append(newInitiatorIqns, connection.IscsiInitiator)
	}

	var existingInitiatorIqns []string
	var existingHostWwns []string
	for _, computeNodePath := range existingCompute.Node.Paths {
		// assuming computeNodePath.HbaName carries IQN or WWN value
		if utils.IsIqn(computeNodePath.HbaName) {
			existingInitiatorIqns = append(existingInitiatorIqns, computeNodePath.HbaName)
		} else {
			existingHostWwns = append(existingHostWwns, computeNodePath.HbaName)
		}
	}
	// Remove duplicate iqn if any
	existingInitiatorIqns = utils.RemoveDuplicateFromStringArray(existingInitiatorIqns)

	// get provisioner manager
	provMgr, err := provisonerimpl.New(psm.storageSetting.Username, psm.storageSetting.Password, psm.storageSetting.ClusterAddress)
	if err != nil {
		log.WriteDebug("TFError| error in provisioner New, err: %v", err)
		return err
	}

	newIqns, removedIqns, commonIqns, err := utils.GetStringSliceDiff(existingInitiatorIqns, newInitiatorIqns)
	if err != nil {
		log.WriteDebug("TFError| error in GetStringSliceDiff, err: %v", err)
		return err
	}

	// Remove
	for _, removedIqn := range removedIqns {
		// Get Initiator ID for IQN
		hbaId, err := provMgr.GetInitiatorIdByServerId(existingCompute.Node.ID, removedIqn)
		if err != nil {
			log.WriteDebug("TFError| failed to call GetInitiatorIdByServerId err: %+v", err)
			return err
		}
		if hbaId != "" {
			// Delete IQN from Compute Node
			err = provMgr.DeleteIscsiHbaFromComputeNode(existingCompute.Node.ID, hbaId)
			if err != nil {
				log.WriteDebug("TFError| error in DeleteIscsiHbaFromComputeNode, err: %v", err)
				return err
			}
		}
	}

	// Add New IQN and its Ports
	for _, newIqn := range newIqns {
		// Register New IQN
		err = provMgr.AddIscsiHbaToComputeNode(existingCompute.Node.ID, newIqn)
		if err != nil {
			log.WriteDebug("TFError| error in AddIscsiHbaToComputeNode, err: %v", err)
			return err
		}
		// Get Ports Name for IQN
		portNames := GetPortsForIqn(newIqn, inputCompute)
		if portNames != nil {
			// Get Initiator ID for IQN
			initiatorId, err := provMgr.GetInitiatorIdByServerId(existingCompute.Node.ID, newIqn)
			if err != nil {
				log.WriteDebug("TFError| failed to call GetInitiatorIdByServerId err: %+v", err)
				return err
			}
			// Get Port Id Array from Name
			portIds, err := provMgr.GetPortsIdsByName(portNames)
			if err != nil {
				log.WriteDebug("TFError| failed to call GetPortsIdsByName err: %+v", err)
				return err
			}
			// For Each Port Id - Add Path of InitiatorId and Port ID
			for _, portId := range portIds {
				err = provMgr.AddStoragePortToComputeNodeHbaByHbaIdAndPortId(existingCompute.Node.ID, initiatorId, portId)
				if err != nil {
					log.WriteDebug("TFError| failed to call AddStoragePortToComputeNodeHbaByHbaIdAndPortId err: %+v", err)
					return err
				}
			}
		}
	}

	// Reconcile for newly Added or Removed Ports form IQN
	for _, commonIqn := range commonIqns {
		initiator, err := GetComputeNodeIscsiInitiator(commonIqn, inputCompute)
		if err != nil {
			log.WriteDebug("TFError| error in GetComputeNodeIscsiInitiator, err: %v", err)
			return err
		}
		err = psm.ReconcileComputeNodeIscsiConnectionStoragePorts(initiator, existingCompute)
		if err != nil {
			log.WriteDebug("TFError| error in ReconcileComputeNodeIscsiConnectionStoragePorts, err: %v", err)
			return err
		}
	}
	return nil
}

// GetPortsForIqn .
func GetPortsForIqn(iqn string, inputCompute *vssbmodel.ComputeResource) []string {
	for _, connection := range inputCompute.IscsiConnections {
		if connection.IscsiInitiator == iqn {
			return connection.PortNames
		}
	}
	return nil
}

func GetMapHbaNameHbaId(existingCompute *provisonermodel.ComputeNodeWithPathDetails) (mapHbaNameHbaId map[string]string) {
	for _, path := range existingCompute.ComputePaths.Data {
		mapHbaNameHbaId[path.IScsiInitiatorIqn] = path.IScsiInitiatorId
	}
	return mapHbaNameHbaId
}

func GetMapPortNamePortId(existingCompute *provisonermodel.ComputeNodeWithPathDetails) (mapPortNamePortId map[string]string) {
	for _, path := range existingCompute.ComputePaths.Data {
		mapPortNamePortId[path.PortName] = path.PortId
	}
	return mapPortNamePortId
}

func FindComputeNodeHbaIdByHbaName(hbaName string, existingCompute *provisonermodel.ComputeNodeWithPathDetails) (foundHbaId string, err error) {
	foundHbaId = ""
	for _, path := range existingCompute.ComputePaths.Data {
		if strings.EqualFold(hbaName, path.IScsiInitiatorIqn) {
			foundHbaId = path.IScsiInitiatorId
		}
	}
	if foundHbaId == "" {
		return foundHbaId, fmt.Errorf("Compute Node HBA %s is not found", hbaName)
	}
	return foundHbaId, nil
}

// FindStoragePort .
func FindStoragePort(portName string, storagePorts *[]provisonermodel.StoragePort) (foundPort *provisonermodel.StoragePort, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	foundPort = nil
	for _, port := range *storagePorts {
		if strings.EqualFold(portName, port.Nickname) {
			foundPort = &port
			break
		}
	}
	if foundPort == nil {
		return nil, fmt.Errorf("Storage Port %s not found.", portName)
	}
	return foundPort, nil
}

// ReconcileComputeNodeIscsiConnectionStoragePorts .
func (psm *vssbStorageManager) ReconcileComputeNodeIscsiConnectionStoragePorts(iscsiConnector *vssbmodel.IscsiConnector, existingCompute *provisonermodel.ComputeNodeWithPathDetails) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := provisonerimpl.New(psm.storageSetting.Username, psm.storageSetting.Password, psm.storageSetting.ClusterAddress)
	if err != nil {
		log.WriteDebug("TFError| error in provisioner New, err: %v", err)
		return err
	}
	// get current ports of the hba
	var existingPortNames []string
	for _, path := range existingCompute.ComputePaths.Data {
		if strings.EqualFold(path.IScsiInitiatorIqn, iscsiConnector.IscsiInitiator) {
			existingPortNames = append(existingPortNames, path.PortNickname)
		}
	}

	addedPorts, removedPorts, _, err := utils.GetStringSliceDiff(existingPortNames, iscsiConnector.PortNames)
	if err != nil {
		log.WriteDebug("TFError| error in GetStringSliceDiff, err: %v", err)
		return err
	}
	// Get Initiator ID for IQN
	hbaId, err := provObj.GetInitiatorIdByServerId(existingCompute.Node.ID, iscsiConnector.IscsiInitiator)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetInitiatorIdByServerId err: %+v", err)
		return err
	}

	// add storage port to the compute node HBA
	storagePortsObj, err := provObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts, err: %v", err)
		return err
	}
	// remove storage port from the compute node HBA
	for _, portName := range removedPorts {
		// verify if the storage port in TF is existed in the VSS-B storage
		storagePort, err := FindStoragePort(portName, &storagePortsObj.Data)
		if err != nil {
			log.WriteDebug("TFError| error in FindStoragePort, err: %v", err)
			return err
		}
		err = provObj.RemoveStoragePortFromComputeNodeHbaByHbaIdAndPortId(existingCompute.Node.ID, hbaId, storagePort.ID)
		if err != nil {
			log.WriteDebug("TFError| error in RemoveStoragePortFromComputeNodeHbaByHbaIdAndPortId, err: %v", err)
			return err
		}
	}
	//Add Port
	for _, portName := range addedPorts {
		// verify if the storage port in TF is existed in the VSS-B storage
		storagePort, err := FindStoragePort(portName, &storagePortsObj.Data)
		if err != nil {
			log.WriteDebug("TFError| error in FindStoragePort, err: %v", err)
			return err
		}
		err = provObj.AddStoragePortToComputeNodeHbaByHbaIdAndPortId(existingCompute.Node.ID, hbaId, storagePort.ID)
		if err != nil {
			log.WriteDebug("TFError| error in AddStoragePortToComputeNodeHbaByHbaIdAndPortId, err: %v", err)
			return err
		}
	}
	return nil
}

// DeleteComputeNodeResource is used to delete the compute node
func (psm *vssbStorageManager) DeleteComputeNodeResource(serverId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_BEGIN), serverId)
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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_BEGIN), serverId)
	err = provObj.DeleteComputeNodeResource(serverId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteComputeNodeResource provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_SERVER_FAILED), serverId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_END), serverId)

	return nil
}
