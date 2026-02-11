package terraform

import (
	"context"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceVspMultipleSnapshotGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource to show snapshots (Thin Image pairs) in snapshot groups in Hitachi VSP storage.",
		ReadContext: datasourceVspMultipleSnapshotGroupsRead,
		Schema:      schemaimpl.DatasourceVspMultipleSnapshotGroupsSchema(),
	}
}

func datasourceVspMultipleSnapshotGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return impl.DatasourceVspMultipleSnapshotGroupsRead(d)
}
