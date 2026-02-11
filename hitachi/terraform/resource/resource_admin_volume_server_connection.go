package terraform

import (
	"context"
	"sync"

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
	// Check if either input list has changed
	volChanged := d.HasChange("volume_ids")
	servChanged := d.HasChange("server_ids")

	if !volChanged && !servChanged {
		return nil
	}

	// Mark computed fields for refresh
	d.SetNewComputed("connections_info")
	d.SetNewComputed("connections_count")

	return nil
}
