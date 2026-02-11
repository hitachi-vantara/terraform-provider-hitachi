package terraform

import (
	"context"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DatasourceAdminVolumeServerConnections lists all volume-server connections with optional filters.
func DatasourceAdminVolumeServerConnections() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource to list volume-server connections in VSP One storage with optional filters.",
		ReadContext: datasourceAdminVolumeServerConnectionsRead,
		Schema:      schemaimpl.DatasourceAdminMultipleVolumeServerConnectionsSchema(),
	}
}

func datasourceAdminVolumeServerConnectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DatasourceAdminMultipleVolumeServerConnectionsRead(d)
}
