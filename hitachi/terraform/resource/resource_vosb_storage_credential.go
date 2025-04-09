package terraform

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
)

var syncChangeUserPasswordOperation = &sync.Mutex{}

// password change is a one-time operation

func ResourceVssbChangeUserPassword() *schema.Resource {
	return &schema.Resource{
		Description:   ":meta:subcategory:VOS Block Change Storage User Password:The following request changes user password.",
		CreateContext: resourceVssbChangeUserPasswordCreate,
		UpdateContext: resourceVssbChangeUserPasswordUpdate,
		DeleteContext: resourceVssbChangeUserPasswordDelete,
		ReadContext:   resourceVssbChangeUserPasswordRead,
		Schema:        schemaimpl.ResourceVssbChangeUserPasswordSchema,
	}
}

// Create method will handle the password change
func resourceVssbChangeUserPasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChangeUserPasswordOperation.Lock()
	defer syncChangeUserPasswordOperation.Unlock()

	log.WriteInfo("starting change user password")
	userInfo, err := impl.ChangeVssbUserPassword(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	userInfoMap := impl.ConvertVssbStorageUserToSchema(userInfo)
	log.WriteDebug("userInfoMap: %+v\n", userInfoMap)

	storageUserList := []interface{}{userInfoMap}
	if err := d.Set("storage_user", storageUserList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(userInfo.UserId)
	d.Set("status", "Password changed successfully")
	log.WriteInfo("password changed successfully")
	return nil
}

func resourceVssbChangeUserPasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceVssbChangeUserPasswordCreate(ctx,d, m)
}

func resourceVssbChangeUserPasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceVssbChangeUserPasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
