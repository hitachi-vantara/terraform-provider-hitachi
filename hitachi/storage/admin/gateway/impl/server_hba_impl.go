package admin

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// GetServerHBAs retrieves all HBAs for a specific server
func (psm *adminStorageManager) GetServerHBAs(serverID int) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var serverHBAList gwymodel.ServerHBAList

	apiSuf := fmt.Sprintf("objects/servers/%d/hbas", serverID)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &serverHBAList)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| Server HBAs retrieved successfully for serverID: %d", serverID)
	return &serverHBAList, nil
}

// GetServerHBAByWwn retrieves a specific HBA by server ID and HBA WWN
func (psm *adminStorageManager) GetServerHBAByWwn(serverID int, hbaWwn string) (*gwymodel.ServerHBA, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var serverHBA gwymodel.ServerHBA

	apiSuf := fmt.Sprintf("objects/servers/%d/hbas/%s", serverID, hbaWwn)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &serverHBA)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| Server HBA retrieved successfully for serverID: %d, hbaWwn: %s", serverID, hbaWwn)
	return &serverHBA, nil
}

// CreateServerHBAs adds HBA information to a server
func (psm *adminStorageManager) CreateServerHBAs(serverID int, params gwymodel.CreateServerHBAParams) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%d/hbas", serverID)

	// Async call to create HBAs
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| Server HBAs created successfully for serverID: %d", serverID)

	// After successful creation, get the updated HBA list
	return psm.GetServerHBAs(serverID)
}

// DeleteServerHBA removes HBA information from a server
func (psm *adminStorageManager) DeleteServerHBA(serverID int, initiatorName string) (*gwymodel.ServerHBAList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/servers/%d/hbas/%s", serverID, initiatorName)

	// Async call to delete HBA
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	log.WriteDebug("TFDebug| Server HBA deleted successfully for serverID: %d, initiatorName: %s", serverID, initiatorName)

	// After successful deletion, get the updated HBA list
	return psm.GetServerHBAs(serverID)
}
