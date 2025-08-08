package terraform

import (
	"context"
	"fmt"

	// "fmt"

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

var syncIscsiTargetOperation = &sync.Mutex{}

func ResourceStorageIscsiTarget() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage iSCSI Target: The following request creates a iSCSI target and the iSCSI name for the port. The host mode and the host mode option can also be specified at the same time when the iSCSI target is created.",

		CreateContext: resourceStorageIscsiTargetCreate,
		ReadContext:   resourceStorageIscsiTargetRead,
		UpdateContext: resourceStorageIscsiTargetUpdate,
		DeleteContext: resourceStorageIscsiTargetDelete,
		Schema:        schemaimpl.ResourceIscsiTargetSchema,
		// CustomizeDiff: customDiffFunc(),
	}
}

// resourceStorageIscsiTargetCreate .
func resourceStorageIscsiTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncIscsiTargetOperation.Lock()
	defer syncIscsiTargetOperation.Unlock()

	log.WriteInfo("starting iscsi target create")

	serial := d.Get("serial").(int)

	iscsiTarget, err := impl.CreateIscsiTarget(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	it := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
	log.WriteDebug("iscsi target: %+v\n", *it)
	itList := []map[string]interface{}{
		*it,
	}
	if err := d.Set("iscsitarget", itList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("iscsi_target_number", iscsiTarget.IscsiTargetNumber)
	createID := fmt.Sprintf("%s%d", iscsiTarget.PortID, iscsiTarget.IscsiTargetNumber)
	d.SetId(createID)
	log.WriteInfo("iscsi target created successfully")

	return nil
}

func resourceStorageIscsiTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceStorageIscsiTargetRead(ctx, d, m)
}

// resourceStorageIscsiTargetUpdate .
func resourceStorageIscsiTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting iscsi target update")

	serial := d.Get("serial").(int)

	iscsiTarget, err := impl.UpdateIscsiTarget(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hg := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
	log.WriteDebug("iscsiTarget: %+v\n", *hg)
	hgList := []map[string]interface{}{
		*hg,
	}
	if err := d.Set("iscsitarget", hgList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("iscsi_target_number", iscsiTarget.IscsiTargetNumber)
	updatedID := fmt.Sprintf("%s%d", iscsiTarget.PortID, iscsiTarget.IscsiTargetNumber)
	d.SetId(updatedID)
	log.WriteInfo("iscsi target updated successfully")

	return nil
}

// resourceStorageIscsiTargetDelete
func resourceStorageIscsiTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting iscsc target delete")

	err := impl.DeleteIscsiTarget(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("iscsi target deleted successfully")
	return nil
}
