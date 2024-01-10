package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vssb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
)

// GetComputeNode gets compute node server
func (psm *vssbStorageManager) GetComputeNode(serverID string) (*vssbmodel.Server, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var server vssbmodel.Server
	apiSuf := fmt.Sprintf("objects/servers/%s", serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &server)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &server, nil
}

// GetAllComputeNodes get all compute node servers
func (psm *vssbStorageManager) GetAllComputeNodes() (*vssbmodel.Servers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var servers vssbmodel.Servers
	apiSuf := "objects/servers"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &servers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &servers, nil
}

// RegisterComputeNode used to register compute node
func (psm *vssbStorageManager) RegisterComputeNode(reqBody *vssbmodel.ComputeNodeCreateReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/")
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeleteComputeNode  is used for Deleting server for id
func (psm *vssbStorageManager) DeleteComputeNode(serverID string) (err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s", serverID)
	_, err = httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// EditComputeNode used to set nick name and os type
func (psm *vssbStorageManager) EditComputeNode(serverID string, reqBody *vssbmodel.ComputeNodeInformation) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s", serverID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// RegisterInitiatorInfoForComputeNode used to register initiator information to server
func (psm *vssbStorageManager) RegisterInitiatorInfoForComputeNode(serverID string, reqBody *vssbmodel.RegisterInitiator) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/hbas", serverID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// RegisterHbaInfoForComputeNode used to register Hba information to server
func (psm *vssbStorageManager) RegisterHbaInfoForComputeNode(serverID string, reqBody *vssbmodel.RegisterHba) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/hbas", serverID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// ConfigureHbaPortsForComputeNode used to mesh fc ports
func (psm *vssbStorageManager) ConfigureHbaPortsForComputeNode(serverID string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/paths", serverID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeleteInitiatorInfoForComputeNode used to Delete initiator information to server
func (psm *vssbStorageManager) DeleteInitiatorInfoForComputeNode(serverID, hbaID string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/hbas/%s", serverID, hbaID)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// GetInitiatorInformationForComputeNode get single initiator information
func (psm *vssbStorageManager) GetInitiatorInformationForComputeNode(serverID, hbaID string) (*vssbmodel.Initiator, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var info vssbmodel.Initiator
	apiSuf := fmt.Sprintf("objects/servers/%s/hbas/%s", serverID, hbaID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &info)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &info, nil
}

// GetInitiatorsInformationForComputeNode get all initiator information for compute node servers
func (psm *vssbStorageManager) GetInitiatorsInformationForComputeNode(serverID string) (*vssbmodel.Initiators, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var infos vssbmodel.Initiators
	apiSuf := fmt.Sprintf("objects/servers/%s/hbas", serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &infos)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &infos, nil
}

// GetPathInfoForComputeNode is used to get path info for compute node
func (psm *vssbStorageManager) GetPathInfoForComputeNode(serverID string, reqBody *vssbmodel.ComputeNodePathReq) (*vssbmodel.ComputeNodePath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pathInfo vssbmodel.ComputeNodePath
	apiSuf := fmt.Sprintf("objects/servers/%s/paths/%s,%s", serverID, reqBody.HbaId, reqBody.PortId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pathInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &pathInfo, nil
}

// GetPathsInfoForComputeNode is used to get paths info for compute nodes
func (psm *vssbStorageManager) GetPathsInfoForComputeNode(serverID string) (*vssbmodel.ComputeNodePaths, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pathsInfo vssbmodel.ComputeNodePaths
	apiSuf := fmt.Sprintf("objects/servers/%s/paths", serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pathsInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &pathsInfo, nil
}

// AddPathInfoToComputeNode is used to add path info to compute node
func (psm *vssbStorageManager) AddPathInfoToComputeNode(serverID string, reqBody *vssbmodel.ComputeNodePathReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/paths", serverID)
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeleteComputeNodePath is used to delete compute node path
func (psm *vssbStorageManager) DeleteComputeNodePath(serverID string, reqBody *vssbmodel.ComputeNodePathReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%s/paths/%s,%s", serverID, reqBody.HbaId, reqBody.PortId)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// GetConnectionInfoBtwnVolumeAndServerByVolumeID is used to get the connection information between server and volume by volumeID
func (psm *vssbStorageManager) GetConnectionInfoBtwnVolumeAndServerByVolumeID(volumeID string) (*vssbmodel.VolumeServerConnectionsInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var connectionInfo vssbmodel.VolumeServerConnectionsInfo
	apiSuf := fmt.Sprintf("objects/volume-server-connections?volumeId=%s", volumeID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &connectionInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &connectionInfo, nil
}

// GetConnectionInfoBtwnVolumeAndServerByServerID is used to get the connection information between server and volume by serverID
func (psm *vssbStorageManager) GetConnectionInfoBtwnVolumeAndServerByServerID(serverID string) (*vssbmodel.VolumeServerConnectionsInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var connectionsInfo vssbmodel.VolumeServerConnectionsInfo
	apiSuf := fmt.Sprintf("objects/volume-server-connections?serverId=%s", serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &connectionsInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &connectionsInfo, nil
}

// GetConnectionInfoBtwnVolumeAndServerBoth is used to get the connection information between server and volume both
func (psm *vssbStorageManager) GetConnectionInfoBtwnVolumeAndServerBoth(volumeID, serverID string) (*vssbmodel.VolumeServerConnectionInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var connectionInfo vssbmodel.VolumeServerConnectionInfo
	apiSuf := fmt.Sprintf("objects/volume-server-connections/%s,%s", volumeID, serverID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &connectionInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &connectionInfo, nil
}

// SetPathBtwnVolumeAndServer is used to set path between server and volume
func (psm *vssbStorageManager) SetPathBtwnVolumeAndServer(reqBody *vssbmodel.SetPathVolumeServerReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volume-server-connections")
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// ReleaseMultipleConnectionsBtwnVolumeAndServer is used to release multiple connections between server and volume
func (psm *vssbStorageManager) ReleaseMultipleConnectionsBtwnVolumeAndServer(reqBody *vssbmodel.ReleaseMultiConVolumeServerReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volume-server-connections/actions/release/invoke")
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// ReleaseConnectionBtwnVolumeAndServer is used to release connection between server and volume
func (psm *vssbStorageManager) ReleaseConnectionBtwnVolumeAndServer(reqBody *vssbmodel.SetPathVolumeServerReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volume-server-connections/%s,%s", reqBody.VolumeId, reqBody.ServerId)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}
