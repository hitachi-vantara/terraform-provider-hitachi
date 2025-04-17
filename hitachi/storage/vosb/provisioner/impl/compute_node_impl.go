package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetComputeNode gets compute node server
func (psm *vssbStorageManager) GetComputeNode(serverID string) (*vssbmodel.ComputeNodeWithPathDetails, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_BEGIN), serverID)
	serverInfo, err := gatewayObj.GetComputeNode(serverID)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_FAILED), serverID)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_END), serverID)

	provServer := vssbmodel.Server{}
	err = copier.Copy(&provServer, serverInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	pathsInfo, err := gatewayObj.GetPathsInfoForComputeNode(serverID)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetPathsInfoForComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_FAILED), serverID)
		return nil, err
	}

	initatorInfo, err := gatewayObj.GetInitiatorsInformationForComputeNode(serverID)
	for i := 0; i < len(initatorInfo.Data); i++ {
		for j := 0; j < len(provServer.Paths); j++ {
			if initatorInfo.Data[i].Name == provServer.Paths[j].HbaName {
				provServer.Paths[j].Protocol = initatorInfo.Data[i].Protocol
				break
			}
		}
	}
	if err != nil {
		log.WriteDebug("TFError| failed to call GetInitiatorsInformationForComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_FAILED), serverID)
		return nil, err
	}

	nodeWithPaths := vssbmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&nodeWithPaths.Node, provServer)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	err = copier.Copy(&nodeWithPaths.ComputePaths, pathsInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &nodeWithPaths, nil
}

// CreateComputeResource use to create compute node resource
func (psm *vssbStorageManager) CreateComputeResource(computeResource *vssbmodel.ComputeResource) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_BEGIN), computeResource.Name)
	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	createReq := vssbgatewaymodel.ComputeNodeCreateReq{
		ServerNickname: computeResource.Name,
		OsType:         computeResource.OsType,
	}

	err = gatewayObj.RegisterComputeNode(&createReq)
	if err != nil {
		log.WriteDebug("TFError| failed to call RegisterComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_SERVER_FAILED), computeResource.Name)
		return err
	}

	serverId, err := psm.GetComputeNodeIdByName(computeResource.Name)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetComputeNodeIdByName err: %+v", err)
		return err
	}
	if serverId == "" {
		log.WriteDebug("TFError| server id is blank")
		return fmt.Errorf("unable to find compute node id")
	}

	for _, connection := range computeResource.IscsiConnections {
		registerInitiatorReq := vssbgatewaymodel.RegisterInitiator{
			Protocol:  "iSCSI",
			IscsiName: connection.IscsiInitiator,
		}
		if connection.IscsiInitiator == "" {
			log.WriteDebug("TFError| Iscsi initiator is blank")
			return fmt.Errorf("please enter valid iscsi initiator")
		}
		err = gatewayObj.RegisterInitiatorInfoForComputeNode(serverId, &registerInitiatorReq)
		if err != nil {
			log.WriteDebug("TFError| failed to call RegisterInitiatorInfoForComputeNode err: %+v", err)
			return err
		}
		initiatorId, err := psm.GetInitiatorIdByServerId(serverId, connection.IscsiInitiator)
		if err != nil {
			log.WriteDebug("TFError| failed to call GetInitiatorsIdByServerId err: %+v", err)
			return err
		}
		if connection.PortNames != nil {
			portsId, err := psm.GetPortsIdsByName(connection.PortNames)
			if err != nil {
				log.WriteDebug("TFError| failed to call GetPortsIdsByName err: %+v", err)
				return err
			}
			for _, portId := range portsId {
				registerPath := vssbgatewaymodel.ComputeNodePathReq{
					HbaId:  initiatorId,
					PortId: portId,
				}
				err = gatewayObj.AddPathInfoToComputeNode(serverId, &registerPath)
				if err != nil {
					log.WriteDebug("TFError| failed to call AddPathInfoToComputeNode err: %+v", err)
					return err
				}
			}
		}
	}

	// register FcConnections
	for _, connection := range computeResource.FcConnections {
		registerHbaReq := vssbgatewaymodel.RegisterHba{
			Protocol:    "FC",
			HbaWwn:      strings.ToLower(connection.HostWWN),
			IsTargetAny: false,
		}
		if connection.HostWWN == "" {
			log.WriteDebug("TFError| Host WWN is blank")
			return fmt.Errorf("please enter valid host wwn")
		}
		err = gatewayObj.RegisterHbaInfoForComputeNode(serverId, &registerHbaReq)
		if err != nil {
			log.WriteDebug("TFError| failed to call RegisterHbaInfoForComputeNode err: %+v", err)
			return err
		}
		err = gatewayObj.ConfigureHbaPortsForComputeNode(serverId)
		if err != nil {
			log.WriteDebug("TFError| failed to call ConfigureHbaPortsForComputeNode err: %+v", err)
			return err
		}

		// FIXME - uncomment following and comment above function ConfigureHbaPortsForComputeNode to not configure fc ports full meshed

		// initiatorId, err := psm.GetInitiatorIdByServerId(serverId, connection.HostWWN)
		// if err != nil {
		// 	log.WriteDebug("TFError| failed to call GetInitiatorsIdByServerId err: %+v", err)
		// 	return err
		// }
		// if connection.PortNames != nil {
		// 	portsId, err := psm.GetPortsIdsByName(connection.PortNames)
		// 	if err != nil {
		// 		log.WriteDebug("TFError| failed to call GetPortsIdsByName err: %+v", err)
		// 		return err
		// 	}
		// 	for _, portId := range portsId {
		// 		registerPath := vssbgatewaymodel.ComputeNodePathReq{
		// 			HbaId:  initiatorId,
		// 			PortId: portId,
		// 		}
		// 		err = gatewayObj.AddPathInfoToComputeNode(serverId, &registerPath)
		// 		if err != nil {
		// 			log.WriteDebug("TFError| failed to call AddPathInfoToComputeNode err: %+v", err)
		// 			return err
		// 		}
		// 	}
		// }
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_SERVER_END), computeResource.Name)
	return nil
}

func (psm *vssbStorageManager) AddFCportsWWN(wwns []string, serverId string) error {
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
		return err
	}
	for _, wwn := range wwns {
		registerHbaReq := vssbgatewaymodel.RegisterHba{
			Protocol:    "FC",
			HbaWwn:      strings.ToLower(wwn),
			IsTargetAny: false,
		}
		err = gatewayObj.RegisterHbaInfoForComputeNode(serverId, &registerHbaReq)
		if err != nil {
			log.WriteDebug("TFError| failed to call RegisterHbaInfoForComputeNode err: %+v", err)
			return err
		}
		err = gatewayObj.ConfigureHbaPortsForComputeNode(serverId)
		if err != nil {
			log.WriteDebug("TFError| failed to call ConfigureHbaPortsForComputeNode err: %+v", err)
			return err
		}

	}
	return nil
}

func (psm *vssbStorageManager) RemoveFCportsWWN(wwns []string, serverId string) error {
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
		return err
	}
	for _, wwn := range wwns {
		initaterId, err := psm.GetInitiatorIdByServerId(serverId, wwn)

		if err != nil {
			log.WriteDebug("TFError| failed to get initiater details from wwns id, err: %+v", err)
			return err
		}
		err = gatewayObj.DeleteInitiatorInfoForComputeNode(serverId, initaterId)
		if err != nil {
			log.WriteDebug("TFError| failed to de-attach WWNS from ComputeNode err: %+v", err)
			return err
		}
		// err = gatewayObj.ConfigureHbaPortsForComputeNode(serverId)
		// if err != nil {
		// 	log.WriteDebug("TFError| failed to call ConfigureHbaPortsForComputeNode err: %+v", err)
		// 	return err
		// }
	}
	return nil
}

// GetPortsIdsByName .
func (psm *vssbStorageManager) GetPortsIdsByName(ports []string) ([]string, error) {
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

	storagePorts, err := gatewayObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStoragePorts err: %+v", err)
		return nil, err
	}
	var portIds []string = nil
	if storagePorts.Data != nil {
		for _, port := range storagePorts.Data {
			for _, inputPort := range ports {
				if port.Nickname == inputPort {
					portIds = append(portIds, port.ID)
				}
			}
			if len(ports) == len(portIds) {
				break
			}
		}
	}
	if len(ports) != len(portIds) {
		return nil, fmt.Errorf("port id is not available, please enter valid port ids")
	}
	return portIds, nil
}

// GetInitiatorIdByServerId .
func (psm *vssbStorageManager) GetInitiatorIdByServerId(serverId string, initiatorName string) (string, error) {
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
		return "", err
	}

	initiators, err := gatewayObj.GetInitiatorsInformationForComputeNode(serverId)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetInitiatorsInformationForComputeNode err: %+v", err)
		return "", err
	}
	if initiators.Data != nil {
		for _, initiator := range initiators.Data {
			if initiator.Name == initiatorName {
				return initiator.ID, nil
			}
		}
	}
	return "", nil
}

// GetComputeNodeIdByName .
func (psm *vssbStorageManager) GetComputeNodeIdByName(computeNodeName string) (string, error) {
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
		return "", err
	}

	servers, err := gatewayObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllComputeNodes err: %+v", err)
		return "", err
	}
	for _, server := range servers.Data {
		if server.Nickname == computeNodeName {
			return server.ID, nil
		}
	}
	return "", nil
	//return "", fmt.Errorf("node not found")
}

// DeleteComputeNodeResource is used to delete the compute node
func (psm *vssbStorageManager) DeleteComputeNodeResource(serverId string) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_BEGIN), serverId)

	err = gatewayObj.DeleteComputeNode(serverId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_SERVER_FAILED), serverId)
		log.WriteDebug("TFError| failed to call DeleteComputeNode err: %+v", err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_END), serverId)

	return nil
}

// GetComputeResourceInfo .
func (psm *vssbStorageManager) GetComputeResourceInfo(computeNodeName string) (*vssbmodel.ComputeResourceOutput, error) {
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
	var responseCompute vssbmodel.ComputeResourceOutput
	servers, err := gatewayObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllComputeNodes err: %+v", err)
		return &responseCompute, err
	}
	for _, server := range servers.Data {
		if server.Nickname == computeNodeName {
			responseCompute.Name = server.Nickname
			responseCompute.OsType = server.OsType
			responseCompute.ID = server.ID
			break
		}
	}

	// Use to determine need to Create or Update compute node
	if responseCompute.ID == "" {
		log.WriteDebug("TFDebug| server id is blank, unable to find compute node")
		return &responseCompute, nil
	}
	initiators, err := gatewayObj.GetInitiatorsInformationForComputeNode(responseCompute.ID)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetInitiatorsInformationForComputeNode err: %+v", err)
		return &responseCompute, err
	}
	storagePorts, err := gatewayObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStoragePorts err: %+v", err)
		return &responseCompute, err
	}
	if initiators.Data != nil {
		for _, initiator := range initiators.Data {
			var finalPortList []string = nil
			for _, initiatorPortId := range initiator.PortIDs {
				for _, storagePort := range storagePorts.Data {
					if initiatorPortId == storagePort.ID {
						finalPortList = append(finalPortList, storagePort.Nickname)
						continue
					}
				}
			}
			responseCompute.IscsiConnection = append(responseCompute.IscsiConnection, vssbmodel.IscsiConnector{IscsiInitiator: initiator.Name, PortNames: finalPortList})
		}
	}
	return &responseCompute, nil
}

// GetAllComputeNodes get all compute node servers
func (psm *vssbStorageManager) GetAllComputeNodes() (*vssbmodel.Servers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_BEGIN))
	serversInfo, err := gatewayObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllComputeNodes err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_SERVERS_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_END))

	provServers := vssbmodel.Servers{}
	err = copier.Copy(&provServers, serversInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provServers, nil
}

// GetConnectionInfoBtwnVolumeAndServerByServerID is used to get the connection information between server and volume by serverID
func (psm *vssbStorageManager) GetConnectionInfoBtwnVolumeAndServerByServerID(serverID string) (*vssbmodel.VolumeServerConnectionsInfo, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CONNECTION_BY_SERVER_BEGIN), serverID)
	ConnectionInfo, err := gatewayObj.GetConnectionInfoBtwnVolumeAndServerByServerID(serverID)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetConnectionInfoBtwnVolumeAndServerByServerID err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CONNECTION_BY_SERVER_FAILED), serverID)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CONNECTION_BY_SERVER_END), serverID)

	provConnections := vssbmodel.VolumeServerConnectionsInfo{}
	err = copier.Copy(&provConnections, ConnectionInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provConnections, nil
}

// UpdateComputeNode .
func (psm *vssbStorageManager) UpdateComputeNode(computeNodeId, computeNodeName, osType string) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_NAME_AND_OS_TYPE_BEGIN), computeNodeName, osType)
	node := vssbgatewaymodel.ComputeNodeInformation{
		Nickname: computeNodeName,
		OsType:   osType,
	}
	err = gatewayObj.EditComputeNode(computeNodeId, &node)

	if err != nil {
		log.WriteDebug("TFError| failed to call EditComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_NAME_AND_OS_TYPE_FAILED), computeNodeName, osType)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_NAME_AND_OS_TYPE_END), computeNodeName, osType)
	return nil
}

// AddStoragePortToComputeNodeHbaByHbaIdAndPortId .
func (psm *vssbStorageManager) AddStoragePortToComputeNodeHbaByHbaIdAndPortId(computeNodeId, computeNodeHbaId, storagePortId string) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_COMPUTE_PATH_INFO_BEGIN), computeNodeHbaId, storagePortId)
	registerPath := vssbgatewaymodel.ComputeNodePathReq{
		HbaId:  computeNodeHbaId,
		PortId: storagePortId,
	}
	err = gatewayObj.AddPathInfoToComputeNode(computeNodeId, &registerPath)
	if err != nil {
		log.WriteDebug("TFError| failed to call AddPathInfoToComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_COMPUTE_PATH_INFO_FAILED), computeNodeHbaId, storagePortId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_COMPUTE_PATH_INFO_END), computeNodeHbaId, storagePortId)
	return nil
}

// RemoveStoragePortFromComputeNodeHbaByHbaIdAndPortId .
func (psm *vssbStorageManager) RemoveStoragePortFromComputeNodeHbaByHbaIdAndPortId(computeNodeId, computeNodeHbaId, storagePortId string) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_COMPUTE_PATH_INFO_BEGIN), computeNodeHbaId, storagePortId)
	registerPath := vssbgatewaymodel.ComputeNodePathReq{
		HbaId:  computeNodeHbaId,
		PortId: storagePortId,
	}
	err = gatewayObj.DeleteComputeNodePath(computeNodeId, &registerPath)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteComputeNodePath err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_REMOVE_COMPUTE_PATH_INFO_FAILED), computeNodeHbaId, storagePortId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_COMPUTE_PATH_INFO_END), computeNodeHbaId, storagePortId)
	return nil
}

// AddIscsiHbaToComputeNode .
func (psm *vssbStorageManager) AddIscsiHbaToComputeNode(computeNodeId, iqn string) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_COMPUTE_IQN_INFO_BEGIN), iqn)
	registerInitiatorReq := vssbgatewaymodel.RegisterInitiator{
		Protocol:  "iSCSI",
		IscsiName: iqn,
	}
	err = gatewayObj.RegisterInitiatorInfoForComputeNode(computeNodeId, &registerInitiatorReq)
	if err != nil {
		log.WriteDebug("TFError| failed to call RegisterInitiatorInfoForComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_COMPUTE_IQN_INFO_FAILED), iqn)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_COMPUTE_IQN_INFO_END), iqn)
	return nil
}

// DeleteIscsiHbaFromComputeNode .
func (psm *vssbStorageManager) DeleteIscsiHbaFromComputeNode(computeNodeId, iqnId string) error {
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
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_COMPUTE_IQN_INFO_BEGIN), iqnId)
	err = gatewayObj.DeleteInitiatorInfoForComputeNode(computeNodeId, iqnId)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteInitiatorInfoForComputeNode err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_REMOVE_COMPUTE_IQN_INFO_FAILED), iqnId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_COMPUTE_IQN_INFO_END), iqnId)
	return nil
}
