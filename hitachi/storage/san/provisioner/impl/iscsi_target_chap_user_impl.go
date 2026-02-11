package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

func (psm *sanStorageManager) GetChapUser(portID string, iscsiTargetNumber int, chapUserName string, wayOfChapUser string) (*sanmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	iscsiTargetChapUser, err := gatewayObj.GetChapUser(portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	log.WriteDebug("iscsiTargetChapUser = %v", iscsiTargetChapUser)
	log.WriteDebug("err = %v", err)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	provChapUser := sanmodel.IscsiTargetChapUser{}
	err = copier.Copy(&provChapUser, iscsiTargetChapUser)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provChapUser, nil
}

func (psm *sanStorageManager) GetChapUsers(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTargetChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN), portID, iscsiTargetNumber)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	iscsiTargetChapUsers, err := gatewayObj.GetChapUsers(portID, iscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSERS_FAILED), portID, iscsiTargetNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_END), portID, iscsiTargetNumber)

	provChapUsers := sanmodel.IscsiTargetChapUsers{}

	err = copier.Copy(&provChapUsers, iscsiTargetChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("Provisioner ChapUsers : %v", provChapUsers)
	return &provChapUsers, nil
}

func (psm *sanStorageManager) CreateChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string, chapUserSecret string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERNAME_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	err = gatewayObj.SetChapUserName(portID, iscsiTargetNumber, wayOfChapUser, chapUserName)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_ISCSITARGET_CHAPUSERNAME_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERNAME_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
	err = gatewayObj.SetChapUserSecret(portID, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERSECRET_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_CHAPUSER_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	return nil
}

func (psm *sanStorageManager) DeleteChapUser(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = gatewayObj.DeleteChapUser(portID, iscsiTargetNumber, wayOfChapUser, chapUserName)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser)

	return nil
}

func (psm *sanStorageManager) ChangeChapUserSecret(portID string, iscsiTargetNumber int, wayOfChapUser, chapUserName string, chapUserSecret string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERSECRET_BEGIN), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
	err = gatewayObj.SetChapUserSecret(portID, iscsiTargetNumber, wayOfChapUser, chapUserName, chapUserSecret)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_ISCSITARGET_CHAPUSERSECRET_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
		log.WriteError(mc.GetMessage(mc.ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)

		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_ISCSITARGET_CHAPUSERSECRET_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)
	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSITARGET_CHAPUSER_END), portID, iscsiTargetNumber, chapUserName, wayOfChapUser, chapUserSecret)

	return nil
}
