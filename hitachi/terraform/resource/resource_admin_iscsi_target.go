package terraform

import (
	"context"
	"sync"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Mutex to prevent concurrent create operations
var syncAdminIscsiTargetOperation = &sync.Mutex{}

func ResourceAdminIscsiTarget() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage iscsiTargets in VSP One storage.",
		CreateContext: resourceAdminIscsiTargetCreate,
		ReadContext:   resourceAdminIscsiTargetRead,
		UpdateContext: resourceAdminIscsiTargetUpdate,
		DeleteContext: resourceAdminIscsiTargetDelete,
		Schema:        schemaimpl.ResourceAdminIscsiTargetSchema(),
	}
}

func resourceAdminIscsiTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncAdminIscsiTargetOperation.Lock()
	defer syncAdminIscsiTargetOperation.Unlock()

	return impl.ResourceAdminIscsiTargetCreate(d)
}

func resourceAdminIscsiTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminIscsiTargetRead(d)
}

func resourceAdminIscsiTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAdminIscsiTargetCreate(ctx, d, m)
}

func resourceAdminIscsiTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminIscsiTargetDelete(d)
}
