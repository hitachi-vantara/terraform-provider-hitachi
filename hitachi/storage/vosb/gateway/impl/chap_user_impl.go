package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vosb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
)

// GetAllChapUsers get all chap users
func (psm *vssbStorageManager) GetAllChapUsers() (*vssbmodel.ChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var chapUsers vssbmodel.ChapUsers
	apiSuf := "objects/chap-users"
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &chapUsers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &chapUsers, nil
}

// CreateChapUser is used to create a chap user
func (psm *vssbStorageManager) CreateChapUser(reqBody *vssbmodel.ChapUserReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/chap-users"
	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeleteChapUser is used to delete a chap user by chap user id
func (psm *vssbStorageManager) DeleteChapUser(chapUserID string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/chap-users/%s", chapUserID)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// GetChapUserInfo is used to get the chap user information
func (psm *vssbStorageManager) GetChapUserInfo(chapUserID string) (*vssbmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var chapUserInfo vssbmodel.ChapUser
	apiSuf := fmt.Sprintf("objects/chap-users/%s", chapUserID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &chapUserInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &chapUserInfo, nil
}

// UpdateChapUser is used to edit a chap user
func (psm *vssbStorageManager) UpdateChapUser(chapUserID string, reqBody *vssbmodel.ChapUserReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/chap-users/%s", chapUserID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// UpdatePortAuthSettings is used to edit the authentication settings for the compute port for the target operation.
func (psm *vssbStorageManager) UpdatePortAuthSettings(portID string, reqBody *vssbmodel.PortAuthSettings) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s", portID)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	log.WriteInfo("UpdatePortAuthSettings Successful")
	return nil
}

// GetChapUsersAllowedToAccessPort is used to obtain a list of information about a CHAP user who is allowed to access the compute port.
func (psm *vssbStorageManager) GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var chapUsers vssbmodel.ChapUsers
	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s/chap-users", portID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &chapUsers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &chapUsers, nil
}

// AllowChapUserToAccessPort is used to allows a CHAP user to access the compute port.
func (psm *vssbStorageManager) AllowChapUserToAccessPort(portID string, reqBody *vssbmodel.ChapUserIdReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s/chap-users", portID)

	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, reqBody)

	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// DeletePortAccessForChapUser used to cancel compute port access permission for a CHAP user.
func (psm *vssbStorageManager) DeletePortAccessForChapUser(portId, chapUserId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s/chap-users/%s", portId, chapUserId)
	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// GetChapUserInfoAllowedToAccessPort is used to obtain a list of information about a CHAP user who is allowed to access the compute port.
func (psm *vssbStorageManager) GetChapUserInfoAllowedToAccessPort(portId, chapUserId string) (*vssbmodel.ChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var chapUsers vssbmodel.ChapUsers
	apiSuf := fmt.Sprintf("objects/port-auth-settings/%s/chap-users/%s", portId, chapUserId)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &chapUsers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &chapUsers, nil
}
