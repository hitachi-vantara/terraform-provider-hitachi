package terraform

import (
	"context"
	"fmt"
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
		Description: "Manage ports in VSP One storage.",
		Importer: &schema.ResourceImporter{
			StateContext: importAdminPort,
		},
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
	return nil
}

func resourceAdminPortCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if d.Id() == "" {
		// Create path requires explicit selectors.
		serialRaw, ok := d.GetOk("serial")
		if !ok || serialRaw.(int) < 1 {
			return fmt.Errorf("serial must be specified and must be >= 1")
		}
		portRaw, ok := d.GetOk("port_id")
		if !ok || portRaw.(string) == "" {
			return fmt.Errorf("port_id must be specified")
		}
	} else {
		// Update/imported resources can rely on state-populated values.
		serial := d.Get("serial").(int)
		if serial < 1 {
			return fmt.Errorf("serial must be known (import or config must provide it)")
		}
		portID := d.Get("port_id").(string)
		if portID == "" {
			return fmt.Errorf("port_id must be known (import or config must provide it)")
		}
	}

	if err := d.SetNewComputed("port_info"); err != nil {
		return err
	}
	return nil
}
