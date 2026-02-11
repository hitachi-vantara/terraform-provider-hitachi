package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetChapUsers Get the information about CHAP users
func (psm *sanStorageManager) GetChapUsers(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargetChapUsers sanmodel.IscsiTargetChapUsers
	apiSuf := fmt.Sprintf("objects/chap-users?portId=%v&hostGroupNumber=%v", portID, iscsiTargetNumber)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargetChapUsers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTargetChapUsers, nil
}

// GetChapUser Get information about a specific CHAP user
func (psm *sanStorageManager) GetChapUser(portID string, iscsiTargetNumber int, chapUserName string, wayOfChapUser string) (*sanmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var iscsiTargetChapUser sanmodel.IscsiTargetChapUser
	apiSuf := fmt.Sprintf("objects/chap-users/%v,%v,%v,%v", portID, iscsiTargetNumber, wayOfChapUser, chapUserName)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &iscsiTargetChapUser)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &iscsiTargetChapUser, nil
}

// SetChapUserName Set the CHAP user name for the iSCSI target
func (psm *sanStorageManager) SetChapUserName(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBody := sanmodel.SetChapUserNameReq{
		ChapUserName:    chapUserName,
		PortID:          portID,
		HostGroupNumber: iscsiTargetNumber,
		WayOfChapUser:   wayOfChapUser,
	}

	_, err := httpmethod.PostCall(psm.storageSetting, "objects/chap-users", reqBody)
	if err != nil {
		log.WriteDebug("TFError| error in SetChapUserName - objects/chap-users API call, err: %v", err)
		return err
	}

	return nil
}

// SetChapUserSecret Set a secret for the CHAP user
func (psm *sanStorageManager) SetChapUserSecret(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName, chapUserPassword string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	reqBody := sanmodel.SetChapUserSecretReq{
		ChapPassword: chapUserPassword,
	}

	apiSuf := fmt.Sprintf("objects/chap-users/%v,%v,%v,%v", portID, iscsiTargetNumber, wayOfChapUser, chapUserName)
	_, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}

// DeleteChapUser Delete the CHAP user from the iSCSI target
func (psm *sanStorageManager) DeleteChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/chap-users/%v,%v,%v,%v", portID, iscsiTargetNumber, wayOfChapUser, chapUserName)

	_, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}

	return nil
}
