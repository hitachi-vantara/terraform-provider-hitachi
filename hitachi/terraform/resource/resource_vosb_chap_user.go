package terraform

import (
	"context"
	"fmt"
	"strings"

	// "fmt"

	// "time"
	// "errors"
	"sync"

	// cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	// reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncChapUserOperation = &sync.Mutex{}

func ResourceVssbStorageChapUser() *schema.Resource {
	return &schema.Resource{
		Description: "VSP One SDS Block iSCSI Target CHAP User: The following request sets the CHAP user.",
		Importer: &schema.ResourceImporter{
			StateContext: importVosbChapUser,
		},
		CreateContext: resourceVssbChapUserCreate,
		ReadContext:   resourceVssbChapUserRead,
		UpdateContext: resourceVssbChapUserUpdate,
		DeleteContext: resourceVssbChapUserDelete,
		Schema:        schemaimpl.ResourceVssbChapUserSchema,
	}
}

func resourceVssbChapUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr, _ := d.Get("vosb_address").(string)
	if strings.TrimSpace(vssbAddr) == "" {
		return diag.FromErr(fmt.Errorf("vosb_address is required to delete a CHAP user"))
	}

	log.WriteInfo("starting chap user resource delete")

	err := impl.DeleteVssbChapUserResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.WriteInfo("chap user resource deleted successfully")
	return nil
}

func resourceVssbChapUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChapUserOperation.Lock()
	defer syncChapUserOperation.Unlock()

	vssbAddr, _ := d.Get("vosb_address").(string)
	if strings.TrimSpace(vssbAddr) == "" {
		return diag.FromErr(fmt.Errorf("vosb_address is required to create a CHAP user"))
	}

	name, _ := d.Get("target_chap_user_name").(string)
	if strings.TrimSpace(name) == "" {
		return diag.FromErr(fmt.Errorf("target_chap_user_name is required to create a CHAP user"))
	}
	secret, _ := d.Get("target_chap_user_secret").(string)
	if strings.TrimSpace(secret) == "" {
		return diag.FromErr(fmt.Errorf("target_chap_user_secret is required to create a CHAP user"))
	}

	log.WriteInfo("starting chap user creation")
	chapUser, err := impl.CreateVssbChapUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cu := impl.ConvertVssbChapUserToSchema(chapUser)
	log.WriteDebug("cu: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_users", cuList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(chapUser.ID)
	log.WriteInfo("chap user  created successfully")
	return nil
}

func resourceVssbChapUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	chapUser, err := impl.GetVssbChapUserByName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cu := impl.ConvertVssbChapUserToSchema(chapUser)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_users", cuList); err != nil {
		return diag.FromErr(err)
	}
	// Populate selectors/IDs so post-import operations don't fail.
	_ = d.Set("chap_user_id", chapUser.ID)
	_ = d.Set("target_chap_user_name", chapUser.TargetChapUserName)
	// Do NOT set secrets here; they are write-only.

	d.SetId(chapUser.ID)
	return nil
}

func resourceVssbChapUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChapUserOperation.Lock()
	defer syncChapUserOperation.Unlock()

	vssbAddr, _ := d.Get("vosb_address").(string)
	if strings.TrimSpace(vssbAddr) == "" {
		return diag.FromErr(fmt.Errorf("vosb_address is required to update a CHAP user"))
	}

	name, _ := d.Get("target_chap_user_name").(string)
	if strings.TrimSpace(name) == "" {
		return diag.FromErr(fmt.Errorf("target_chap_user_name is required to update a CHAP user"))
	}
	secret, _ := d.Get("target_chap_user_secret").(string)
	if strings.TrimSpace(secret) == "" {
		return diag.FromErr(fmt.Errorf("target_chap_user_secret is required to update a CHAP user"))
	}

	log.WriteInfo("starting chap user update")
	chapUser, err := impl.UpdateVssbChapUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cu := impl.ConvertVssbChapUserToSchema(chapUser)
	log.WriteDebug("cu: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_users", cuList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(chapUser.ID)
	log.WriteInfo("chap user updated successfully")
	return nil
}
