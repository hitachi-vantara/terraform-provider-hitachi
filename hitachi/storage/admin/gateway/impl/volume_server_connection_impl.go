package admin

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) GetVolumeServerConnections(params model.GetVolumeServerConnectionsParams) (*model.VolumeServerConnectionsResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var connList model.VolumeServerConnectionsResponse
	q := url.Values{}

	if params.ServerId != nil {
		q.Add("serverId", fmt.Sprintf("%d", *params.ServerId))
	}
	if params.ServerNickname != nil {
		q.Add("serverNickname", *params.ServerNickname)
	}
	if params.StartVolumeId != nil {
		q.Add("startVolumeId", fmt.Sprintf("%d", *params.StartVolumeId))
	}
	if params.Count != nil {
		q.Add("count", fmt.Sprintf("%d", *params.Count))
	}

	log.WriteDebug("TFDebug | QueryParams:%+v", q)

	apiSuf := fmt.Sprintf("objects/volume-server-connections?%s", q.Encode())
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &connList)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to get volume-server connections: %w", err))
		return nil, err
	}

	return &connList, nil
}

func (psm *adminStorageManager) GetOneVolumeServerConnection(volumeId, serverId int) (*model.VolumeServerConnectionDetail, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var connList model.VolumeServerConnectionDetail
	apiSuf := fmt.Sprintf("objects/volume-server-connections/%d,%d", volumeId, serverId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &connList)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to get volume-server connections: %w", err))
		return nil, err
	}

	return &connList, nil
}

func (psm *adminStorageManager) AttachVolumeToServers(params model.AttachVolumeServerConnectionParam) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/volume-server-connections"
	connIDs, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to create volume-server connection: %w", err))
		return "", err
	}

	log.WriteDebug("TFDebug | created connection IDs = %s", *connIDs)
	return *connIDs, nil
}

func (psm *adminStorageManager) DetachVolumeToServers(volumeId, serverId int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/volume-server-connections/%d,%d", volumeId, serverId)

	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to delete volume-server connection: %w", err))
		return err
	}

	log.WriteDebug("TFDebug | deleted volume-server-connection %d,%d", volumeId, serverId)
	return nil
}
