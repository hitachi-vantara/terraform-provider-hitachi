package terraform

import (
	"context"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DatasourceAdminVolumeServerConnection retrieves information about one volume-server connection.
func DatasourceAdminVolumeServerConnection() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource to show a specific volume-server connection in VSP One storage.",
		ReadContext: datasourceAdminVolumeServerConnectionRead,
		Schema:      schemaimpl.DatasourceAdminOneVolumeServerConnectionSchema(),
	}
}

func datasourceAdminVolumeServerConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DatasourceAdminOneVolumeServerConnectionRead(d)
}
