package terraform

import (
	"context"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source for listing servers
func DataSourceAdminServerList() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to list servers in VSP One Block.",
		Schema:      schemaimpl.DataSourceAdminServerListSchema,
		ReadContext: dataSourceAdminServerListRead,
	}
}

// Data source for getting server info by ID
func DataSourceAdminServerInfo() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to get server information by ID in VSP One Block.",
		Schema:      schemaimpl.DataSourceAdminServerInfoSchema,
		ReadContext: dataSourceAdminServerInfoRead,
	}
}

func dataSourceAdminServerListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DataSourceAdminServerListRead(ctx, d, m)
}

func dataSourceAdminServerInfoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DataSourceAdminServerInfoRead(ctx, d, m)
}

// Data source for getting server path info
func DatasourceAdminServerPath() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to get server path information in VSP One Block.",
		Schema:      schemaimpl.DataSourceAdminServerPathSchema,
		ReadContext: dataSourceAdminServerPathRead,
	}
}

func dataSourceAdminServerPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DataSourceAdminServerPathRead(ctx, d, m)
}
