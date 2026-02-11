package terraform

import (
	"context"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source for listing pools
func DatasourceAdminPools() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to list pools in VSP One Block.",
		Schema:      schemaimpl.DataSourceAdminPoolsSchema(),
		ReadContext: dataSourceAdminPoolsRead,
	}
}

// Data source for getting pool info by ID
func DatasourceAdminPool() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to get pool information by ID in VSP One Block.",
		Schema:      schemaimpl.DataSourceAdminPoolSchema(),
		ReadContext: dataSourceAdminPoolRead,
	}
}

func dataSourceAdminPoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DataSourceAdminPoolsRead(ctx, d, m)
}

func dataSourceAdminPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DataSourceAdminPoolRead(ctx, d, m)
}
