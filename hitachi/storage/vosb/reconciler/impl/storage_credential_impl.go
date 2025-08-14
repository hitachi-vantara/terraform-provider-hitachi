package vssbstorage

import (
	// "fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
)

// ChangeUserPassword is used to change the password of a user
func (psm *vssbStorageManager) ChangeUserPassword(userID string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_USER_PASSWORD_BEGIN), userID)
	user := provisonermodel.ChangeUserPasswordReq{
		CurrentPassword: reqBody.CurrentPassword,
		NewPassword:     reqBody.NewPassword,
	}
	userInfo, err := provObj.ChangeUserPassword(userID, &user)
	if err != nil {
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
