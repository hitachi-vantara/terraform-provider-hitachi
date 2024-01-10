package vssbstorage

import (
	"errors"
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vssb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetAllChapUsers get all chap users
func (psm *vssbStorageManager) GetAllChapUsers() (*vssbmodel.ChapUsers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))
	chapUsersInfo, err := gatewayObj.GetAllChapUsers()
	log.WriteDebug("TFError| Gateway chapUsersInfo: %v", chapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END))

	provChapUsers := vssbmodel.ChapUsers{}
	err = copier.Copy(&provChapUsers, chapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| Prov chapUsersInfo: %v", provChapUsers)

	return &provChapUsers, nil
}

// CreateChapUser is used to create a chap user
func (psm *vssbStorageManager) CreateChapUser(chapUserResource *vssbmodel.ChapUserReq) error {
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

	createReq := vssbgatewaymodel.ChapUserReq{
		TargetChapUserName:    chapUserResource.TargetChapUserName,
		TargetChapSecret:      chapUserResource.TargetChapSecret,
		InitiatorChapUserName: chapUserResource.InitiatorChapUserName,
		InitiatorChapSecret:   chapUserResource.InitiatorChapSecret,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_BEGIN), chapUserResource.TargetChapUserName)
	err = gatewayObj.CreateChapUser(&createReq)
	if err != nil {
		log.WriteDebug("TFError| failed to call CreateChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_CHAP_USER_FAILED), chapUserResource.TargetChapUserName)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_END), chapUserResource.TargetChapUserName)
	return nil
}

// DeleteChapUser is used to delete a chap user by chap user
func (psm *vssbStorageManager) DeleteChapUser(chapUserId string) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_BEGIN), chapUserId)

	err = gatewayObj.DeleteChapUser(chapUserId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_CHAP_USER_FAILED), chapUserId)
		log.WriteDebug("TFError| failed to call DeleteChapUser err: %+v", err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_END), chapUserId)

	return nil
}

// GetChapUserInfoById is used to get the chap user info by chap user Id
func (psm *vssbStorageManager) GetChapUserInfoById(chapUserId string) (*vssbmodel.ChapUser, error) {
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
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), chapUserId)
	chapUser, err := gatewayObj.GetChapUserInfo(chapUserId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), chapUserId)
		log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
		return nil, err
	} else {
		provChapUser := vssbmodel.ChapUser{}
		err = copier.Copy(&provChapUser, chapUser)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_END), chapUserId)
		return &provChapUser, nil
	}

}

// GetChapUserInfoByName is used to get the chap user info by target chap user name
func (psm *vssbStorageManager) GetChapUserInfoByName(chapUserTargetName string) (*vssbmodel.ChapUser, error) {
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

	chapUsers, err := gatewayObj.GetAllChapUsers()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
		return nil, err
	}

	var chapUserId string
	found := false
	for _, chapUser := range chapUsers.Data {
		if chapUser.TargetChapUserName == chapUserTargetName {
			chapUserId = chapUser.ID
			found = true
			break
		}
	}

	if found {
		log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), chapUserTargetName)
		chapUser, err := gatewayObj.GetChapUserInfo(chapUserId)
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), chapUserTargetName)
			log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
			return nil, err
		} else {
			provChapUser := vssbmodel.ChapUser{}
			err = copier.Copy(&provChapUser, chapUser)
			if err != nil {
				log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
				return nil, err
			}
			log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_END), chapUserTargetName)
			return &provChapUser, nil
		}

	} else {
		err_text := fmt.Sprintf("Target Chap User %s does not exit", chapUserTargetName)
		err = errors.New(err_text)
		return nil, err
	}
}

// UpdateChapUserById is used to edit a chap user
func (psm *vssbStorageManager) UpdateChapUserById(chapUserID string, reqBody *vssbmodel.ChapUserReq) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_BEGIN), chapUserID, reqBody.TargetChapUserName)
	chapUser := vssbgatewaymodel.ChapUserReq{
		TargetChapUserName:    reqBody.TargetChapUserName,
		TargetChapSecret:      reqBody.TargetChapSecret,
		InitiatorChapUserName: reqBody.InitiatorChapUserName,
		InitiatorChapSecret:   reqBody.InitiatorChapSecret,
	}
	err = gatewayObj.UpdateChapUser(chapUserID, &chapUser)
	if err != nil {
		log.WriteDebug("TFError| failed to call UpdateChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_CHAP_USER_FAILED), chapUserID, reqBody.TargetChapUserName)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_END), chapUserID, reqBody.TargetChapUserName)
	return nil

}

// UpdateChapUserInfoByName is used to get the chap user info by target chap user name
func (psm *vssbStorageManager) UpdateChapUserInfoByName(reqBody *vssbmodel.ChapUserReq) error {
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

	chapUsers, err := gatewayObj.GetAllChapUsers()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
		return err
	}

	var chapUserId string
	found := false
	for _, chapUser := range chapUsers.Data {
		if chapUser.TargetChapUserName == reqBody.TargetChapUserName {
			chapUserId = chapUser.ID
			found = true
			break
		}
	}

	if found {
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_BEGIN), chapUserId, reqBody.TargetChapUserName)
		chapUser := vssbgatewaymodel.ChapUserReq{
			TargetChapUserName:    reqBody.TargetChapUserName,
			TargetChapSecret:      reqBody.TargetChapSecret,
			InitiatorChapUserName: reqBody.InitiatorChapUserName,
			InitiatorChapSecret:   reqBody.InitiatorChapSecret,
		}
		err = gatewayObj.UpdateChapUser(chapUserId, &chapUser)
		if err != nil {
			log.WriteDebug("TFError| failed to call UpdateChapUser err: %+v", err)
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_CHAP_USER_FAILED), chapUserId, reqBody.TargetChapUserName)
			return err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_END), chapUserId, reqBody.TargetChapUserName)
		return nil

	} else {
		err_text := fmt.Sprintf("Target Chap User %s does not exit", reqBody.TargetChapUserName)
		err = errors.New(err_text)
		return err
	}
}

func (psm *vssbStorageManager) GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))
	chapUsersInfo, err := gatewayObj.GetChapUsersAllowedToAccessPort(portID)
	log.WriteDebug("TFError| Gateway chapUsersInfo: %v", chapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllChapUsers err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END))

	provChapUsers := vssbmodel.ChapUsers{}
	err = copier.Copy(&provChapUsers, chapUsersInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| Prov chapUsersInfo: %v", provChapUsers)

	return &provChapUsers, nil

}
