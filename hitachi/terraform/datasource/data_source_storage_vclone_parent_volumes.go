package terraform

import (
	"context"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceVspVirtualCloneParentVolume() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource to show volumes that are acting as parents for virtual clones (TIA) in Hitachi VSP storage.",
		ReadContext: datasourceVspVirtualCloneParentVolumeRead,
		Schema:      schemaimpl.DatasourceVspVirtualCloneParentVolumeSchema(),
	}
}

func datasourceVspVirtualCloneParentVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DatasourceVspVirtualCloneParentVolumeRead(d)
}
