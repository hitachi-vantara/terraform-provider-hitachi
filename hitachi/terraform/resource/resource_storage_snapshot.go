package terraform

import (
	"context"
	"fmt"

	// "strconv"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncSnapshotOperation = &sync.Mutex{}

func ResourceVspSnapshot() *schema.Resource {
	return &schema.Resource{
		Description:   `Vsp Snapshot Resource: Manages Thin Image and Thin Image Advanced pairs. Supports actions such as create, split, resync, restore, clone, and vclone.`,
		CreateContext: resourceVspSnapshotCreate,
		ReadContext:   resourceVspSnapshotRead,
		UpdateContext: resourceVspSnapshotUpdate,
		DeleteContext: resourceVspSnapshotDelete,
		Schema:        schemaimpl.ResourceVspSnapshotSchema(),
		CustomizeDiff: resourceVspSnapshotCustomizeDiff,
	}
}

func resourceVspSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncSnapshotOperation.Lock()
	defer syncSnapshotOperation.Unlock()

	return impl.ResourceVspSnapshotApply(d)
}

func resourceVspSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceVspSnapshotRead(d)
}

func resourceVspSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncSnapshotOperation.Lock()
	defer syncSnapshotOperation.Unlock()

	return impl.ResourceVspSnapshotApply(d)
}

func resourceVspSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncSnapshotOperation.Lock()
	defer syncSnapshotOperation.Unlock()

	return impl.ResourceVspSnapshotDelete(d)
}

func resourceVspSnapshotCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	log := commonlog.GetLogger()

	// --- 1. CREATE PHASE VALIDATION ---
	if d.Id() == "" {
		action := d.Get("state").(string)
		_, muOk := d.GetOk("mirror_unit_id")

		if action != "" && action != "create" && !muOk {
			err := fmt.Errorf("mirror_unit_id must be specified")
			log.WriteError("TFError| %v", err)
			return err
		}
	}

	// --- 2. UPDATE PHASE VALIDATION (Immutability) ---
	if d.Id() != "" {
		immutableFields := []string{
			"serial",
			"snapshot_group_name",
			"snapshot_pool_id",
			"pvol_ldev_id",
			"pvol_ldev_id_hex",
			"is_consistency_group",
			"is_clone",
			"auto_clone",
			"can_cascade",
			"copy_speed",
			"is_data_reduction_force_copy",
		}

		for _, field := range immutableFields {
			if d.HasChange(field) {
				old, new := d.GetChange(field)

				// FIX: Allow the change if the new value is a zero-value (Destroying)
				// This prevents the "cannot change from X to 0" error during destroy.
				if isZeroValue(new) {
					continue
				}

				err := fmt.Errorf("%s is immutable: cannot change from %v to %v. Recreate resource if needed", field, old, new)
				log.WriteError("TFError| %v", err)
				return err
			}
		}

		// Fix for mirror_unit_id as well
		if d.HasChange("mirror_unit_id") {
			old, new := d.GetChange("mirror_unit_id")
			if !isZeroValue(new) {
				err := fmt.Errorf("mirror_unit_id is immutable: cannot change from %v to %v", old, new)
				log.WriteError("TFError| %v", err)
				return err
			}
		}
	}

	d.SetNewComputed("snapshot")
	d.SetNewComputed("vclone")
	d.SetNewComputed("additional_info")
	return nil
}

// Helper function to detect zero values across types
func isZeroValue(v interface{}) bool {
	switch t := v.(type) {
	case int:
		return t == 0
	case string:
		return t == ""
	case bool:
		// Careful: bool is tricky. If false is a valid "new" state,
		// you might need d.GetRawConfig() or check if the resource is being deleted.
		// For snapshots, if the ID is being cleared, we usually care more about the LDEV IDs.
		return false
	}
	return v == nil
}
