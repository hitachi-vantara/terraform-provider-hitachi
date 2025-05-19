package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
)

var syncAddDrivesToPoolOperation = &sync.Mutex{}

func ResourceVssbAddDrivesToPool() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block: Add drives to a storage pool.",
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
	addAllOffline := diff.Get("add_all_offline_drives").(bool)
	driveIds := diff.Get("drive_ids").([]interface{})

	if addAllOffline && len(driveIds) > 0 {
		return fmt.Errorf("'add_all_offline_drives' cannot be true when 'drive_ids' are specified. Use one or the other.")
	}

	if !addAllOffline && len(driveIds) == 0 {
		return fmt.Errorf("you must specify either 'add_all_offline_drives' or provide at least one 'drive_ids'")
	}

	return nil
}
