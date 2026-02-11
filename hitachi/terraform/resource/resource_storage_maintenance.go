package terraform

import (
	"context"
	"fmt"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceStorageMaintenance() *schema.Resource {
	return &schema.Resource{
		Description:   "Invoke an appliance-wide stop-format action. This resource performs the action on create and (optionally) update; it does not perform any action on delete.",
		CreateContext: resourceStorageMaintenanceCreate,
		ReadContext:   resourceStorageMaintenanceRead,
		UpdateContext: resourceStorageMaintenanceUpdate,
		DeleteContext: resourceStorageMaintenanceDelete,
		Schema:        schemaimpl.ResourceStorageMaintenanceSchema,
	}
}

func resourceStorageMaintenanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Always attempt to invoke if `should_stop_all_volume_format` is true (default).
	if v, ok := d.GetOk("should_stop_all_volume_format"); ok && v.(bool) {
		if err := impl.StopAllVolumeFormat(d); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}

	serial := d.Get("serial").(int)
	// Use timestamped id so re-creating or re-running update generates a new id.
	d.SetId(fmt.Sprintf("stop-format-invoked/%d-%d", serial, time.Now().UnixNano()))

	return resourceStorageMaintenanceRead(ctx, d, m)
}

func resourceStorageMaintenanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// This resource is effectively an action trigger; there is no remote
	// persistent object to refresh. Keep state as-is unless serial is missing.
	if _, ok := d.GetOk("serial"); !ok {
		d.SetId("")
	}
	return nil
}

func resourceStorageMaintenanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Only invoke during update if `should_stop_all_volume_format` is present and true. If the
	// attribute is unset/false, skip invoking.
	if v, ok := d.GetOk("should_stop_all_volume_format"); ok && v.(bool) {
		if err := impl.StopAllVolumeFormat(d); err != nil {
			return diag.FromErr(err)
		}
		serial := d.Get("serial").(int)
		d.SetId(fmt.Sprintf("stop-format-invoked/%d-%d", serial, time.Now().UnixNano()))
	}

	return resourceStorageMaintenanceRead(ctx, d, m)
}

func resourceStorageMaintenanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Resource intentionally performs no destroy action; simply remove state.
	d.SetId("")
	return nil
}
