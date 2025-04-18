package vssbstorage

import (
	// "errors"
	// "fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"

	"github.com/jinzhu/copier"
)

// ChangeUserPassword is used to change the password of a user
func (psm *vssbStorageManager) ChangeUserPassword(userID string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_USER_PASSWORD_BEGIN), userID)
	user := vssbgatewaymodel.ChangeUserPasswordReq{
		CurrentPassword: reqBody.CurrentPassword,
		NewPassword:     reqBody.NewPassword,
	}
	userInfo, err := gatewayObj.ChangeUserPassword(userID, &user)
	if err != nil {
		log.WriteDebug("TFError| failed to call ChangeUserPassword err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CHANGE_USER_PASSWORD_FAILED), userID)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_USER_PASSWORD_END), userID)

	provUserInfo := vssbmodel.User{}
	err = copier.Copy(&provUserInfo, userInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| Prov userInfo: %v", provUserInfo)

	return &provUserInfo, nil
}
