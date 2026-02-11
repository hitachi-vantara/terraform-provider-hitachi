package terraform

import (
	"context"
	"fmt"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAdminServer() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage VSP One servers in Hitachi storage system",
		CreateContext: resourceAdminServerCreate,
		ReadContext:   resourceAdminServerRead,
		UpdateContext: resourceAdminServerUpdate,
		DeleteContext: resourceAdminServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:        schemaimpl.ResourceAdminServerSchema,
		CustomizeDiff: resourceAdminServerCustomizeDiff,
	}
}

func resourceAdminServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerCreate(d)
}

func resourceAdminServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerRead(d)
}

func resourceAdminServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerUpdate(d)
}

func resourceAdminServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminServerDelete(d)
}

func resourceAdminServerCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateHostGroups(ctx, d, meta); err != nil {
		return err
	}

	return nil
}

func validateHostGroups(ctx context.Context, d *schema.ResourceDiff, _ interface{}) error {
	hostGroups := d.Get("host_groups").([]interface{})
	for i, raw := range hostGroups {
		if raw == nil {
			continue
		}
		hg := raw.(map[string]interface{})
		idSet := hg["host_group_id"] != nil && hg["host_group_id"].(int) != 0
		nameSet := hg["host_group_name"] != nil && hg["host_group_name"].(string) != ""

		if !idSet && !nameSet {
			return fmt.Errorf("host_groups[%d]: either host_group_id or host_group_name must be specified", i)
		}
		if idSet && nameSet {
			return fmt.Errorf("host_groups[%d]: only one of host_group_id or host_group_name can be specified", i)
		}
	}
	return nil
}
