package terraform

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		CustomizeDiff: resourceStorageIscsiTargetCustomDiff,
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
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	iscsiTarget, err := impl.GetIscsiTarget(d)
	if err != nil {
		return diag.FromErr(err)
	}

	it := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
	log.WriteDebug("iscsiTarget: %+v\n", *it)

	itList := []map[string]interface{}{
		*it,
	}

	if err := d.Set("iscsitarget", itList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(iscsiTarget.PortID + strconv.Itoa(iscsiTarget.IscsiTargetNumber))
	log.WriteInfo("iscsiTarget read successfully")

	return nil
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

func resourceStorageIscsiTargetCustomDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Fix console output
	d.SetNewComputed("iscsitarget")

	return validateLunIscsi(ctx, d, meta)
}

func validateLunIscsi(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	lunSet, ok := d.GetOk("lun")
	if !ok {
		return nil
	}

	for _, item := range lunSet.(*schema.Set).List() {
		m := item.(map[string]interface{})

		// Retrieve actual values
		idVal := m["ldev_id"].(int)      // default = -1
		hexVal := m["ldev_id_hex"].(string) // default = ""

		hasID := idVal != -1   // -1 = not provided
		hasHex := hexVal != "" // empty string = not provided

		// Both provided → invalid
		if hasID && hasHex {
			return fmt.Errorf("only one of ldev_id or ldev_id_hex may be specified in each lun block")
		}

		// Neither provided → invalid
		if !hasID && !hasHex {
			return fmt.Errorf("one of ldev_id or ldev_id_hex must be specified in each lun block")
		}
	}

	return nil
}
