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

// Mutex to prevent concurrent attach/detach operations.
var syncVolumeServerConnectionOperation = &sync.Mutex{}

// ResourceAdminVolumeServerConnection manages the attachment and detachment of volumes to servers in VSP One storage.
func ResourceAdminVolumeServerConnection() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage volume-server connections (attach or detach volumes to/from servers) in VSP One storage.",
		CreateContext: resourceAdminVolumeServerConnectionCreate,
		ReadContext:   resourceAdminVolumeServerConnectionRead,
		UpdateContext: resourceAdminVolumeServerConnectionUpdate,
		DeleteContext: resourceAdminVolumeServerConnectionDelete,
		Schema:        schemaimpl.ResourceAdminVolumeServerConnectionSchema(),
		CustomizeDiff: resourceAdminVolumeServerConnectionCustomizeDiff,
	}
}

func resourceAdminVolumeServerConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeServerConnectionOperation.Lock()
	defer syncVolumeServerConnectionOperation.Unlock()

	return impl.ResourceAdminVolumeServerConnectionCreate(d)
}

func resourceAdminVolumeServerConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeServerConnectionOperation.Lock()
	defer syncVolumeServerConnectionOperation.Unlock()

	return impl.ResourceAdminVolumeServerConnectionUpdate(d)
}

func resourceAdminVolumeServerConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminVolumeServerConnectionRead(d)
}

func resourceAdminVolumeServerConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeServerConnectionOperation.Lock()
	defer syncVolumeServerConnectionOperation.Unlock()

	return impl.ResourceAdminVolumeServerConnectionDelete(d)
}

func resourceAdminVolumeServerConnectionCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	servHasChanged := d.HasChange("server_ids")

	// 1. Get existence checks for the two mutually exclusive fields
	_, volIDsExists := d.GetOk("volume_ids")
	_, volHexsExists := d.GetOk("volume_id_hexes")

	// Check if either field is being changed (only for logging)
	idsHasChange := d.HasChange("volume_ids")
	hexsHasChange := d.HasChange("volume_id_hexes")

	log.WriteDebug("ValidateVolumeIDValues called with:")
	log.WriteDebug("  volume_ids exists: %v", volIDsExists)
	log.WriteDebug("  volume_id_hexes exists: %v", volHexsExists)
	log.WriteDebug("  HasChange(volume_ids): %v", idsHasChange)
	log.WriteDebug("  HasChange(volume_id_hexes): %v", hexsHasChange)
	log.WriteDebug("  HasChange(server_ids): %v", servHasChanged)

	if !servHasChanged && !idsHasChange && !hexsHasChange {
		return nil
	}

	// --- Core Validation: Mutual Exclusivity (XOR Logic) ---

	// 1. Check for mutual exclusivity: Cannot have both volume_ids and volume_id_hexes
	if volIDsExists && volHexsExists {
		log.WriteError("Invalid: both volume_ids and volume_id_hexes specified (must be mutually exclusive).")
		return fmt.Errorf("volume_ids and volume_id_hexes cannot both be specified")
	}

	// 2. Check for presence: Must have either volume_ids OR volume_id_hexes
	if !volIDsExists && !volHexsExists {
		log.WriteError("Invalid: neither volume_ids nor volume_id_hexes specified (one is required).")
		return fmt.Errorf("one of volume_ids or volume_id_hexes must be specified")
	}

	log.WriteInfo("Volume input validation passed.")

	// Mark computed fields for refresh
	d.SetNewComputed("connections_info")
	d.SetNewComputed("connections_count")

	return nil
}
