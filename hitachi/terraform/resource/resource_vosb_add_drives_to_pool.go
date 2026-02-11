package terraform

import (
	"context"
	"fmt"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncAddDrivesToPoolOperation = &sync.Mutex{}

func ResourceVssbAddDrivesToPool() *schema.Resource {
	return &schema.Resource{
		Description:   "VSP One SDS Block: Add drives to a storage pool.",
		CreateContext: resourceVssbAddDrivesToPoolCreate,
		ReadContext:   resourceVssbAddDrivesToPoolRead,
		UpdateContext: resourceVssbAddDrivesToPoolUpdate,
		DeleteContext: resourceVssbAddDrivesToPoolDelete,
		Schema:        schemaimpl.ResourceVssbStoragePoolSchema,
		CustomizeDiff: validateAddDrivesToPoolInputs,
	}
}

func resourceVssbAddDrivesToPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncAddDrivesToPoolOperation.Lock()
	defer syncAddDrivesToPoolOperation.Unlock()

	log.WriteInfo("Starting AddDrivesToPool operation")

	err := impl.AddDrivesToStoragePool(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to add drives to storage pool: %w", err))
	}

	dv, ok := d.GetOk("drive_ids")
	log.WriteDebug("drive_ids: %v, ok: %v", dv, ok)
	if !ok {
		if err := d.Set("drive_ids", []string{}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("add-drives-%s", d.Get("storage_pool_name").(string)))
	d.Set("status", "Drives added successfully")
	log.WriteInfo("Drives added to storage pool successfully")
	return nil
}

func resourceVssbAddDrivesToPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Re-run the create logic
	return resourceVssbAddDrivesToPoolCreate(ctx, d, m)
}

func resourceVssbAddDrivesToPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No real deletion here; it's a one-time action
	d.SetId("")
	return nil
}

func resourceVssbAddDrivesToPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No backend state to read
	return nil
}

func validateAddDrivesToPoolInputs(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	addAllOffline, aOk := diff.GetOk("add_all_offline_drives")
	driveIds, dOk := diff.GetOk("drive_ids")

	log.WriteDebug("Validating inputs: add_all_offline_drives=%v, drive_ids=%v", addAllOffline, driveIds)

	if aOk && dOk {
		if addAllOffline.(bool) && len(driveIds.([]interface{})) > 0 {
			log.WriteDebug("'add_all_offline_drives' cannot be true when 'drive_ids' are specified. Use one or the other")
			return fmt.Errorf("'add_all_offline_drives' cannot be true when 'drive_ids' are specified. Use one or the other")
		}
	}

	if !addAllOffline.(bool) && len(driveIds.([]interface{})) == 0 {
		log.WriteDebug("Either 'add_all_offline_drives' must be true or at least one 'drive_ids' must be provided")
		return fmt.Errorf("you must specify either 'add_all_offline_drives' or provide at least one 'drive_ids'")
	}

	return nil
}
