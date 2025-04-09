package terraform

import (
	"context"
	"fmt"

	// "time"
	// "errors"
	"sync"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncChapUserOperation = &sync.Mutex{}

func ResourceVssbStorageChapUser() *schema.Resource {
	return &schema.Resource{
		Description:   ":meta:subcategory:VSS Block iSCSI Target CHAP User:The following request sets the CHAP user.",
		CreateContext: resourceVssbChapUserCreate,
		ReadContext:   resourceVssbChapUserRead,
		UpdateContext: resourceVssbChapUserUpdate,
		DeleteContext: resourceVssbChapUserDelete,
		Schema:        schemaimpl.ResourceVssbChapUserSchema,
		CustomizeDiff: resourceChapUserResourceCustomDiff,
	}
}

func resourceVssbChapUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting chap user resource delete")

	err := impl.DeleteVssbChapUserResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("chap user resource deleted successfully")
	return nil
}

func resourceVssbChapUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChapUserOperation.Lock()
	defer syncChapUserOperation.Unlock()

	log.WriteInfo("starting chap user creation")
	chapUser, err := impl.CreateVssbChapUser(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cu := impl.ConvertVssbChapUserToSchema(chapUser)
	log.WriteDebug("cu: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_users", cuList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(chapUser.ID)
	log.WriteInfo("chap user  created successfully")
	return nil
}

func resourceVssbChapUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbChapUsersRead(ctx, d, m)
}

func resourceVssbChapUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChapUserOperation.Lock()
	defer syncChapUserOperation.Unlock()

	log.WriteInfo("starting chap user update")
	chapUser, err := impl.UpdateVssbChapUser(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cu := impl.ConvertVssbChapUserToSchema(chapUser)
	log.WriteDebug("cu: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_users", cuList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(chapUser.ID)
	log.WriteInfo("chap user updated successfully")
	return nil
}

func resourceChapUserResourceCustomDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_block_address").(string)

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
		log.WriteDebug("TFError| error in Reconciler NewEx, err: %v", err)
		return err
	}

	name, ok := d.GetOk("target_chap_user_name")
	if !ok {

		log.WriteDebug("target_chap_user_name: %s", name.(string))
		return fmt.Errorf("target_chap_user_name")
	}
	tcun := name.(string)
	//data := &schema.ResourceData{}
	id := d.Id()

	if id != "" {
		chapUser, err := reconObj.GetChapUserInfoById(id)
		if err != nil {
			return fmt.Errorf("Error getting chap user by id")
		}
		log.WriteDebug("chap user id not set: %v", chapUser)
	} else {
		chapUser, err := reconObj.GetChapUserInfoByName(tcun)
		if err != nil {
			return fmt.Errorf("Error getting chap user by name")
		}
		log.WriteDebug("chap user id set id : %v, chapUser %v", id, chapUser)

	}
	//computeNodeCheck := d.GetRawConfig().GetAttr("compute_nodes").IsNull()
	return nil
}
