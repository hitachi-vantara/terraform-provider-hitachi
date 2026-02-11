package terraform

import (
	"context"
	"sync"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Mutex to prevent concurrent operations
var syncServerHbaOperation = &sync.Mutex{}

func ResourceAdminServerHba() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage server HBAs in VSP One storage. Adds and removes HBA information (WWN or iSCSI name) for a server.",
		CreateContext: resourceAdminServerHbaCreate,
		ReadContext:   resourceAdminServerHbaRead,
		DeleteContext: resourceAdminServerHbaDelete,
		Schema:        schemaimpl.ResourceAdminServerHbaSchema(),
	}
}

func resourceAdminServerHbaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncServerHbaOperation.Lock()
	defer syncServerHbaOperation.Unlock()

	return impl.ResourceAdminServerHbaCreate(d)
}

func resourceAdminServerHbaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerHbaRead(d)
}

func resourceAdminServerHbaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncServerHbaOperation.Lock()
	defer syncServerHbaOperation.Unlock()

	return impl.ResourceAdminServerHbaDelete(d)
}
