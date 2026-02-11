package terraform

import (
	"context"

	// "time"
	// "errors"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncIscsiChapUserOperation = &sync.Mutex{}

func ResourceStorageIscsiChapUser() *schema.Resource {
	return &schema.Resource{
		Description:   `VSP Storage iSCSI Target CHAP User: Sets the CHAP user name for the iSCSI target. Two types of CHAP user names can be set: the CHAP user name of the iSCSI target side and the CHAP user name of the host (iSCSI initiator) that connects to the iSCSI target.`,
		CreateContext: resourceStorageIscsiChapUserCreate,
		ReadContext:   resourceStorageIscsiChapUserRead,
		UpdateContext: resourceStorageIscsiChapUserUpdate,
		DeleteContext: resourceStorageIscsiChapUserDelete,
		Schema:        schemaimpl.ResourceIscsiChapUserSchema,
		// CustomizeDiff: customDiffFunc(),
	}
}

// resourceStorageIscsiChapUserCreate .
func resourceStorageIscsiChapUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncIscsiChapUserOperation.Lock()
	defer syncIscsiChapUserOperation.Unlock()

	log.WriteInfo("starting iscsi chap user create")

	serial := d.Get("serial").(int)

	iscsiChapUser, err := impl.CreateIscsiTargetChapUser(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	it := impl.ConvertIscsiTargetChapUserToSchema(iscsiChapUser, serial)
	log.WriteDebug("iscsi target chap user : %+v\n", *it)
	itList := []map[string]interface{}{
		*it,
	}
	if err := d.Set("chap_user", itList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("iscsi_target_number", iscsiChapUser.HostGroupNumber)
	d.SetId(iscsiChapUser.ChapUserID)
	log.WriteInfo("iscsi target chap user created successfully")

	return nil
}

func resourceStorageIscsiChapUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceStorageChapUserRead(ctx, d, m)
}

// resourceStorageIscsiTargetUpdate .
func resourceStorageIscsiChapUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting iscsi target chap user update")

	serial := d.Get("serial").(int)

	iscsiTargetChapUser, err := impl.UpdateIscsiTargetChapUser(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cu := impl.ConvertIscsiTargetChapUserToSchema(iscsiTargetChapUser, serial)
	log.WriteDebug("iscsiTarget: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_user", cuList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("iscsi_target_number", iscsiTargetChapUser.HostGroupNumber)
	d.SetId(iscsiTargetChapUser.ChapUserID)
	log.WriteInfo("iscsi target updated successfully")

	return nil
}

// resourceStorageIscsiChapUserDelete
func resourceStorageIscsiChapUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting iscsi target chap user delete")

	err := impl.DeleteIscsiTargetChapUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("iscsi target chap user deleted successfully")
	return nil
}
