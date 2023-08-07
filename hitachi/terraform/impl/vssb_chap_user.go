package terraform

import (
	// "encoding/json"
	// "context"
	//"fmt"
	// "io/ioutil"
	// "time"

	"errors"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetAllVssbChapUsers(d *schema.ResourceData) (*terraformmodel.ChapUsers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))

	reconChapUsers, err := reconObj.GetAllChapUsers()
	log.WriteDebug("TFError| Recon chapUsersInfo: %v", reconChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllChapUsers, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))

		return nil, err
	}

	// Converting reconciler to terraform
	terraformChapUsers := terraformmodel.ChapUsers{}
	err = copier.Copy(&terraformChapUsers, reconChapUsers)

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END))
	log.WriteDebug("TFError| Terraform chapUsersInfo: %v", terraformChapUsers)

	return &terraformChapUsers, nil
}

func ConvertVssbChapUserToSchema(chapUser *terraformmodel.ChapUser) *map[string]interface{} {

	vssbcu := map[string]interface{}{
		"chap_user_id":             chapUser.ID,
		"target_chap_user_name":    chapUser.TargetChapUserName,
		"initiator_chap_user_name": chapUser.InitiatorChapUserName,
	}

	return &vssbcu
}

func GetVssbChapUserById(d *schema.ResourceData) (*terraformmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	chapUserId := d.Get("target_chap_user").(string)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), chapUserId)

	reconChapUser, err := reconObj.GetChapUserInfoById(chapUserId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVssbChapUserById, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), chapUserId)

		return nil, err
	}

	// Converting reconciler to terraform
	terraformChapUser := terraformmodel.ChapUser{}
	err = copier.Copy(&terraformChapUser, reconChapUser)

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END), chapUserId)

	return &terraformChapUser, nil
}

func GetVssbChapUserByName(d *schema.ResourceData) (*terraformmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	chapUserName, ok := d.Get("target_chap_user").(string)
	if !ok {
		chapUserName, ok = d.Get("target_chap_user_name").(string)
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_CHAP_USER_BEGIN), chapUserName)

	reconChapUser, err := reconObj.GetChapUserInfoByName(chapUserName)
	log.WriteDebug("TFError| Recon chapUsersInfo: %v", reconChapUser)
	if err != nil {
		log.WriteDebug("TFError| error getting GetVssbChapUserByName, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_CHAP_USER_FAILED), chapUserName)

		return nil, err
	}

	// Converting reconciler to terraform
	terraformChapUser := terraformmodel.ChapUser{}
	err = copier.Copy(&terraformChapUser, reconChapUser)

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_END), chapUserName)
	log.WriteDebug("TFError| Terraform chapUsersInfo: %v", terraformChapUser)

	return &terraformChapUser, nil
}

func DeleteVssbChapUserResource(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}

	chapUserId := d.State().ID
	if chapUserId == "" {
		return errors.New("id field is empty, it is required to delete the chap user resource")
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_BEGIN), chapUserId)
	err = reconObj.DeleteChapUser(chapUserId)
	if err != nil {
		log.WriteDebug("TFError| error getting DeleteComputeNodeResource, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_CHAP_USER_FAILED), chapUserId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_CHAP_USER_END), chapUserId)

	return nil
}

func CreateVssbChapUser(d *schema.ResourceData) (*terraformmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}
	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}
	createInput, err := CreateVssbChapUserFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_BEGIN), createInput.TargetChapUserName)
	reconcilerCreateChapUser := reconcilermodel.ChapUserReq{}
	err = copier.Copy(&reconcilerCreateChapUser, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	chapUser, err := reconObj.ReconcileChapUser(&reconcilerCreateChapUser)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_CHAP_USER_FAILED), createInput.TargetChapUserName)
		log.WriteDebug("TFError| error in Creating Chap User - ReconcileChapUser , err: %v", err)
		return nil, err
	}

	terraformModelChapUser := terraformmodel.ChapUser{}
	err = copier.Copy(&terraformModelChapUser, chapUser)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_CHAP_USER_END), createInput.TargetChapUserName)
	return &terraformModelChapUser, nil
}

func CreateVssbChapUserFromSchema(d *schema.ResourceData) (*terraformmodel.ChapUserReq, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.ChapUserReq{}

	tcuname, ok := d.GetOk("target_chap_user_name")
	if ok {
		tname := tcuname.(string)
		createInput.TargetChapUserName = tname
	}

	tsecret, ok := d.GetOk("target_chap_user_secret")
	if ok {
		ts := tsecret.(string)
		createInput.TargetChapSecret = ts
	}

	icuname, ok := d.GetOk("initiator_chap_user_name")
	if ok {
		iname := icuname.(string)
		createInput.InitiatorChapUserName = iname
	}

	isecret, ok := d.GetOk("initiator_chap_user_secret")
	if ok {
		is := isecret.(string)
		createInput.InitiatorChapSecret = is
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func UpdateVssbChapUser(d *schema.ResourceData) (*terraformmodel.ChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}
	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}
	updateInput, err := CreateVssbChapUserFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_BEGIN), updateInput.TargetChapUserName)
	reconcilerCreateChapUser := reconcilermodel.ChapUserReq{}
	err = copier.Copy(&reconcilerCreateChapUser, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	// Add ID for Chap User
	chapUserId := d.State().ID
	if chapUserId == "" {
		return nil, errors.New("unable to find chap user id, update operation failed")
	}
	reconcilerCreateChapUser.ID = chapUserId

	chapUser, err := reconObj.ReconcileChapUser(&reconcilerCreateChapUser)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_CHAP_USER_FAILED), updateInput.TargetChapUserName)
		log.WriteDebug("TFError| error in Updating ComputeNode - ReconcileComputeNode , err: %v", err)
		return nil, err
	}

	terraformModelChapUser := terraformmodel.ChapUser{}
	err = copier.Copy(&terraformModelChapUser, chapUser)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_CHAP_USER_END), updateInput.TargetChapUserName)
	return &terraformModelChapUser, nil
}
