package terraform

import (
	"context"
	"fmt"
	"strings"

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
		Description: `VSP Storage iSCSI Target CHAP User: Sets the CHAP user name for the iSCSI target. Two types of CHAP user names can be set: the CHAP user name of the iSCSI target side and the CHAP user name of the host (iSCSI initiator) that connects to the iSCSI target.`,
		Importer: &schema.ResourceImporter{
			StateContext: importVspIscsiChapUser,
		},
		CreateContext: resourceStorageIscsiChapUserCreate,
		ReadContext:   resourceStorageIscsiChapUserRead,
		UpdateContext: resourceStorageIscsiChapUserUpdate,
		DeleteContext: resourceStorageIscsiChapUserDelete,
		Schema:        schemaimpl.ResourceIscsiChapUserSchema,
		CustomizeDiff: resourceStorageIscsiChapUserCustomizeDiff,
	}
}

func resourceStorageIscsiChapUserCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Create-only required fields (import should not require config to provide them).
	if d.Id() == "" {
		serialRaw, ok := d.GetOk("serial")
		if !ok || serialRaw.(int) < 1 {
			return fmt.Errorf("serial must be specified and must be >= 1")
		}
		portIDRaw, ok := d.GetOk("port_id")
		if !ok || strings.TrimSpace(portIDRaw.(string)) == "" {
			return fmt.Errorf("port_id must be specified")
		}
		targetNumRaw, ok := d.GetOkExists("iscsi_target_number")
		if !ok || targetNumRaw.(int) < 0 {
			return fmt.Errorf("iscsi_target_number must be specified and must be >= 0")
		}
		chapTypeRaw, ok := d.GetOk("chap_user_type")
		if !ok || strings.TrimSpace(chapTypeRaw.(string)) == "" {
			return fmt.Errorf("chap_user_type must be specified")
		}
		t := strings.ToLower(strings.TrimSpace(chapTypeRaw.(string)))
		if t != "target" && t != "initiator" {
			return fmt.Errorf("chap_user_type must be one of: target, initiator")
		}
		chapNameRaw, ok := d.GetOk("chap_user_name")
		if !ok || strings.TrimSpace(chapNameRaw.(string)) == "" {
			return fmt.Errorf("chap_user_name must be specified")
		}
	} else {
		// For update/imported resources, lookup keys must be known (from state or config).
		if d.Get("serial").(int) < 1 {
			return fmt.Errorf("serial must be known (import or config must provide it)")
		}
		if strings.TrimSpace(d.Get("port_id").(string)) == "" {
			return fmt.Errorf("port_id must be known (import or config must provide it)")
		}
		if _, ok := d.GetOkExists("iscsi_target_number"); !ok {
			return fmt.Errorf("iscsi_target_number must be known (import or config must provide it)")
		}
		chapType := strings.ToLower(strings.TrimSpace(d.Get("chap_user_type").(string)))
		if chapType != "target" && chapType != "initiator" {
			return fmt.Errorf("chap_user_type must be known and must be one of: target, initiator")
		}
		if strings.TrimSpace(d.Get("chap_user_name").(string)) == "" {
			return fmt.Errorf("chap_user_name must be known (import or config must provide it)")
		}

		// The lookup keys identify the CHAP user; changing them would retarget the resource.
		if d.HasChange("serial") || d.HasChange("port_id") || d.HasChange("iscsi_target_number") || d.HasChange("chap_user_type") || d.HasChange("chap_user_name") {
			return fmt.Errorf("serial, port_id, iscsi_target_number, chap_user_type, and chap_user_name are immutable after create/import")
		}
	}

	// Always refresh output.
	if err := d.SetNewComputed("chap_user"); err != nil {
		return err
	}
	return nil
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
		return diag.FromErr(err)
	}

	it := impl.ConvertIscsiTargetChapUserToSchema(iscsiChapUser, serial)
	log.WriteDebug("iscsi target chap user : %+v\n", *it)
	itList := []map[string]interface{}{
		*it,
	}
	if err := d.Set("chap_user", itList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("iscsi_target_number", iscsiChapUser.HostGroupNumber); err != nil {
		return diag.FromErr(err)
	}
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
		return diag.FromErr(err)
	}

	cu := impl.ConvertIscsiTargetChapUserToSchema(iscsiTargetChapUser, serial)
	log.WriteDebug("iscsiTarget: %+v\n", *cu)
	cuList := []map[string]interface{}{
		*cu,
	}
	if err := d.Set("chap_user", cuList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("iscsi_target_number", iscsiTargetChapUser.HostGroupNumber); err != nil {
		return diag.FromErr(err)
	}
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
	log.WriteInfo("iscsi target chap user deleted successfully")
	return nil
}
