package terraform

import (
	"context"
	"sync"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Mutex to prevent concurrent update operations
var syncPortOperation = &sync.Mutex{}

func ResourceAdminPort() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage ports in VSP One storage.",
		CreateContext: resourceAdminPortCreate,
		ReadContext:   resourceAdminPortRead,
		UpdateContext: resourceAdminPortUpdate,
		DeleteContext: resourceAdminPortDelete,
		Schema:        schemaimpl.ResourceAdminPortSchema(),
		CustomizeDiff: resourceAdminPortCustomizeDiff,
	}
}

func resourceAdminPortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminPortRead(d)
}

func resourceAdminPortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For hardware port management, "create" means:
	// 1. Verify the port exists
	// 2. Read current state
	// 3. Apply any configuration changes
	// This allows terraform apply to work on first run without separate import

	// First, read the current port state to verify it exists
	if diags := impl.ResourceAdminPortRead(d); diags.HasError() {
		return diags
	}

	// If port exists and we have configuration to apply, update it
	syncPortOperation.Lock()
	defer syncPortOperation.Unlock()

	return impl.ResourceAdminPortUpdate(d)
}

func resourceAdminPortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncPortOperation.Lock()
	defer syncPortOperation.Unlock()

	return impl.ResourceAdminPortUpdate(d)
}

func resourceAdminPortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Dummy delete implementation - no actual API call performed
	// This is a placeholder to satisfy Terraform's validation requirements
	d.SetId("")
	return nil
}

func resourceAdminPortCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	d.SetNewComputed("port_info")
	return nil
}
