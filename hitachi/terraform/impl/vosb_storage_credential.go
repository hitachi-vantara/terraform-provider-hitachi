package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
)

// ConvertVssbStorageUserToSchema converts a userInfo struct to a pointer to a map[string]interface{} for Terraform
func ConvertVssbStorageUserToSchema(userInfo *terraformmodel.User) map[string]interface{} {
	// Format time as a string in the desired format (e.g., "2006-01-02 15:04:05")
	expirationTimeStr := userInfo.PasswordExpirationTime.Format("2006-01-02 15:04:05")

	// Construct user_groups as a slice of maps
	userGroups := make([]interface{}, len(userInfo.UserGroups))
	for i, group := range userInfo.UserGroups {
		userGroups[i] = map[string]interface{}{
			"user_group_id":        group.UserGroupId,
			"user_group_object_id": group.UserGroupObjectId,
		}
	}

	// Construct privileges as a slice of maps
	privileges := make([]interface{}, len(userInfo.Privileges))
	for i, privilege := range userInfo.Privileges {
		privileges[i] = map[string]interface{}{
			"scope":      privilege.Scope,
			"role_names": privilege.RoleNames,
		}
	}

	// Create the final map and return it
	return map[string]interface{}{
		"user_id":                  userInfo.UserId,
		"user_object_id":           userInfo.UserObjectId,
		"password_expiration_time": expirationTimeStr, // Convert time.Time to string
		"is_enabled":               userInfo.IsEnabled,
		"user_groups":              userGroups, // List of maps for user_groups
		"is_built_in":              userInfo.IsBuiltIn,
		"authentication":           userInfo.Authentication,
		"role_names":               userInfo.RoleNames,
		"is_enabled_console_login": userInfo.IsEnabledConsoleLogin,
		"vps_id":                   userInfo.VpsId,
		"privileges":               privileges, // List of maps for privileges
	}
}

func CreateVssbChangeUserPasswordInputFromSchema(d *schema.ResourceData) (*terraformmodel.ChangeUserPasswordReq, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	input := terraformmodel.ChangeUserPasswordReq{}

	tsecret_current, ok := d.GetOk("current_password")
	if ok {
		ts := tsecret_current.(string)
		input.CurrentPassword = ts
	}

	tsecret_new, ok := d.GetOk("new_password")
	if ok {
		ts := tsecret_new.(string)
		input.NewPassword = ts
	}

	return &input, nil
}

func ChangeVssbUserPassword(d *schema.ResourceData) (*terraformmodel.User, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_block_address").(string)

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

	userId := d.Get("user_id").(string)
	inputReq, err := CreateVssbChangeUserPasswordInputFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_USER_PASSWORD_BEGIN), userId)
	reconChangeUserPasswordReq := reconcilermodel.ChangeUserPasswordReq{}
	err = copier.Copy(&reconChangeUserPasswordReq, inputReq)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	userInfo, err := reconObj.ChangeUserPassword(userId, &reconChangeUserPasswordReq)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CHANGE_USER_PASSWORD_FAILED), userId)
		log.WriteDebug("TFError| error in Updating ComputeNode - ReconcileComputeNode , err: %v", err)
		return nil, err
	}

	terraformModelUserInfo := terraformmodel.User{}
	err = copier.Copy(&terraformModelUserInfo, userInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_USER_PASSWORD_END), userId)
	return &terraformModelUserInfo, nil
}
