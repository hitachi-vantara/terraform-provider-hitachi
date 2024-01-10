package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetChapUsers Get Chap Users information
func (psm *sanStorageManager) GetChapUsers(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN), portID, iscsiTargetNumber)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	provChapUsersInfo, err := provObj.GetChapUsers(portID, iscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetChapUsers provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSERS_FAILED), portID, iscsiTargetNumber)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerChapUsers := sanmodel.IscsiTargetChapUsers{}
	err = copier.Copy(&reconcilerChapUsers, provChapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("reconcilerChapUsers %v\n", reconcilerChapUsers)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_END), portID, iscsiTargetNumber)
	return &reconcilerChapUsers, nil
}

// GetChapUser Get Chap User information
func (psm *sanStorageManager) GetChapUser(portID string, iscsiTargetNumber int, chapUserName, wayOfChapUser string) (*sanmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	provChapUsersInfo, err := provObj.GetChapUser(portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	if err != nil {
		log.WriteDebug("TFError| error in GetChapUser provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerChapUser := sanmodel.IscsiTargetChapUser{}
	err = copier.Copy(&reconcilerChapUser, provChapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	return &reconcilerChapUser, nil
}

// ReconcileChapUser will reconcile and call Create/Update Chap User
func (psm *sanStorageManager) ReconcileChapUser(createInput *sanmodel.ChapUserRequest) (*sanmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	portID := createInput.PortID
	iscsiTargetNumber := createInput.IscsiTargetNumber
	chapUserName := createInput.ChapUserName
	wayOfChapUser := createInput.WayOfChapUser
	chapUserSecret := createInput.ChapUserSecret

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	provChapUserInfo, err := provObj.GetChapUser(portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	if err != nil {
		log.WriteDebug("TFError| error in GetChapUser provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

		// chap user  (user name, port id, iscsiTargetNumber, wayOfChapUser) does not exist, call provisioner to create chap user
		err := provObj.CreateChapUser(portID, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
		if err != nil {
			log.WriteInfo(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
			return nil, err
		} else {
			// provChapUserInfo is nill if chap user does not exit, populate it after the creation of the chap user
			provChapUserInfo, err = provObj.GetChapUser(portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
			if err != nil {
				log.WriteDebug("TFError| error in GetChapUser provisioner call, err: %v", err)
				return nil, err
			}
		}
	}
	if provChapUserInfo != nil {
		// chap user  (user name, port id, iscsiTargetNumber, wayOfChapUser) exists, update chap user secret
		err := provObj.ChangeChapUserSecret(portID, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
		if err != nil {
			log.WriteDebug("TFError| error in ChangeChapUserSecret provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
			return nil, err
		}

	}
	// Converting Prov to Reconciler
	reconcilerChapUser := sanmodel.IscsiTargetChapUser{}
	err = copier.Copy(&reconcilerChapUser, provChapUserInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	return &reconcilerChapUser, nil
}

// DeleteChapUser Delete Chap User
func (psm *sanStorageManager) DeleteChapUser(portID string, iscsiTargetNumber int, chapUserName string, wayOfChapUser string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = provObj.DeleteChapUser(portID, iscsiTargetNumber, wayOfChapUser, chapUserName)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteChapUser provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	return nil
}
