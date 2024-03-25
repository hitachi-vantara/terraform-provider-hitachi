package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetAllChapUsers get all chap users
func (psm *vssbStorageManager) GetAllChapUsers() (*vssbmodel.ChapUsers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))
	provChapUsers, err := provObj.GetAllChapUsers()
	log.WriteDebug("TFError| Prov chapUsersInfo: %v", provChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllChapUsers provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerChapUsers := vssbmodel.ChapUsers{}
	err = copier.Copy(&reconcilerChapUsers, provChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END))
	log.WriteDebug("TFError| Recon chapUsersInfo: %v", reconcilerChapUsers)
	return &reconcilerChapUsers, nil
}

// GetChapUserInfoById is used to get the chap user info by chap user Id
func (psm *vssbStorageManager) GetChapUserInfoById(chapUserId string) (*vssbmodel.ChapUser, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), chapUserId)
	provChapUser, err := provObj.GetChapUserInfoById(chapUserId)
	if err != nil {
		log.WriteDebug("TFError| error in GetChapUserInfoById provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), chapUserId)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerChapUser := vssbmodel.ChapUser{}
	err = copier.Copy(&reconcilerChapUser, provChapUser)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_END), chapUserId)
	return &reconcilerChapUser, nil
}

// GetChapUserInfoByName is used to get the chap user info by target chap user name
func (psm *vssbStorageManager) GetChapUserInfoByName(targetChapUserName string) (*vssbmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), targetChapUserName)
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), targetChapUserName)
	provChapUser, err := provObj.GetChapUserInfoByName(targetChapUserName)
	if err != nil {
		log.WriteDebug("TFError| error in GetChapUserInfoByName provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), targetChapUserName)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcilerChapUser := vssbmodel.ChapUser{}
	err = copier.Copy(&reconcilerChapUser, provChapUser)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_END), targetChapUserName)
	return &reconcilerChapUser, nil
}

// CreateChapUser is used to create a chap user
func (psm *vssbStorageManager) CreateChapUser(chapUserResource *vssbmodel.ChapUserReq) error {
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
		return err
	}
	createReq := provisonermodel.ChapUserReq{
		TargetChapUserName:    chapUserResource.TargetChapUserName,
		TargetChapSecret:      chapUserResource.TargetChapSecret,
		InitiatorChapUserName: chapUserResource.InitiatorChapUserName,
		InitiatorChapSecret:   chapUserResource.InitiatorChapSecret,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_BEGIN), chapUserResource.TargetChapUserName)
	err = provObj.CreateChapUser(&createReq)
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

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_BEGIN), chapUserId)

	err = provObj.DeleteChapUser(chapUserId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_CHAP_USER_FAILED), chapUserId)
		log.WriteDebug("TFError| failed to call DeleteChapUser err: %+v", err)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_END), chapUserId)

	return nil
}

// UpdateChapUser is used to edit a chap user
func (psm *vssbStorageManager) UpdateChapUser(reqBody *vssbmodel.ChapUserReq) error {
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
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_BEGIN), reqBody.ID, reqBody.TargetChapUserName)
	chapUser := provisonermodel.ChapUserReq{
		ID:                    reqBody.ID,
		TargetChapUserName:    reqBody.TargetChapUserName,
		TargetChapSecret:      reqBody.TargetChapSecret,
		InitiatorChapUserName: reqBody.InitiatorChapUserName,
		InitiatorChapSecret:   reqBody.InitiatorChapSecret,
	}
	err = provObj.UpdateChapUserById(chapUser.ID, &chapUser)
	if err != nil {
		log.WriteDebug("TFError| failed to call UpdateChapUser err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_CHAP_USER_FAILED), reqBody.ID, reqBody.TargetChapUserName)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_END), reqBody.ID, reqBody.TargetChapUserName)
	return nil

}

// GetExistingChapUserInformation .
func (psm *vssbStorageManager) GetExistingChapUserInformation(chapUserName string, id string) (*provisonermodel.ChapUser, error) {
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
	var existingResource *provisonermodel.ChapUser = nil

	if id != "" {
		existingResource, err = provObj.GetChapUserInfoById(id)
		if err != nil {
			log.WriteDebug("TFError| error in GetChapUserInfoById provisioner call, err: %v", err)
			return nil, err
		}
	} else if id == "" && chapUserName != "" {
		existingResource, err = provObj.GetChapUserInfoByName(chapUserName)
		if err != nil {
			log.WriteDebug("TFError| error in GetChapUserInfoByName provisioner call, err: %v", err)
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("invalid chap user name  and chap user id")
	}
	return existingResource, nil
}

// ReconcileChapUser .
func (psm *vssbStorageManager) ReconcileChapUser(inputChapUser *vssbmodel.ChapUserReq) (*vssbmodel.ChapUser, error) {
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
	existingResource, err := psm.GetExistingChapUserInformation(inputChapUser.TargetChapUserName, inputChapUser.ID)
	if err != nil {
		// The chap user does not exist
		log.WriteDebug("TFError| error in GetExistingChapUserInformation provisioner call, err: %v", err)
		existingResource = nil

	}
	log.WriteDebug("TFDebug| existingResource: %v", existingResource)

	// CREATE RESOURCE
	if existingResource == nil {
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_BEGIN), inputChapUser.TargetChapUserName)

		// Converting Reconciler to Prov
		provResource := provisonermodel.ChapUserReq{}
		err = copier.Copy(&provResource, inputChapUser)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to prov structure, err: %v", err)
			return nil, err
		}
		err := provObj.CreateChapUser(&provResource)
		if err != nil {
			log.WriteDebug("TFError| error in CreateChapUser provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_CHAP_USER_FAILED), inputChapUser.TargetChapUserName)
			return nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_END), inputChapUser.TargetChapUserName)
	} else {
		// update chap user if input id and existing id are same
		if inputChapUser.ID == existingResource.ID {
			if existingResource.TargetChapUserName == inputChapUser.TargetChapUserName {
				// if chap user id and target chap user name are same, should not pass
				// target chap user name in the request body of the PATCH call
				inputChapUser.TargetChapUserName = ""
			}
			err := psm.UpdateChapUser(inputChapUser)
			if err != nil {
				log.WriteDebug("TFError| error in UpdateChapUser provisioner call, err: %v", err)
				return nil, err
			}

		}
	}

	// Read resource after all operations
	provisionerResource, err := psm.GetChapUserInfoById(existingResource.ID)
	if err != nil {
		log.WriteDebug("TFError| error in GetChapUserInfoByName provisioner call, err: %v", err)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerResource := vssbmodel.ChapUser{}
	err = copier.Copy(&reconcilerResource, provisionerResource)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	return &reconcilerResource, nil
}

/*
func HasResourceChanged(existing *provisonermodel.ChapUser, req *vssbmodel.ChapUserReq) (bool, error) {

	if existing.ID != req.ID {
		return false, fmt.Errorf("something wrong, existing id and requested id did not match, but try to update")
	}

	if existing.TargetChapUserName

	return true, nil
}
*/

func (psm *vssbStorageManager) GetChapUsersAllowedToAccessPort(portID string) (*vssbmodel.ChapUsers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))
	provChapUsers, err := provObj.GetChapUsersAllowedToAccessPort(portID)
	log.WriteDebug("TFError| Prov chapUsersInfo: %v", provChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllChapUsers provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerChapUsers := vssbmodel.ChapUsers{}
	err = copier.Copy(&reconcilerChapUsers, provChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END))
	log.WriteDebug("TFError| Recon chapUsersInfo: %v", reconcilerChapUsers)
	return &reconcilerChapUsers, nil
}
