package terraform

import (
	"context"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAdminServerPath() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage server paths in Hitachi storage system",
		CreateContext: resourceAdminServerPathCreate,
		ReadContext:   resourceAdminServerPathRead,
		UpdateContext: resourceAdminServerPathUpdate,
		DeleteContext: resourceAdminServerPathDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: schemaimpl.ResourceAdminServerPathSchema,
	}
}

func resourceAdminServerPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerPathCreate(d)
}

func resourceAdminServerPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerPathRead(d)
}

func resourceAdminServerPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerPathUpdate(d)
}

func resourceAdminServerPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerPathDelete(d)
}
