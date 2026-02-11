package terraform

import (
    "context"

    commonlog "terraform-provider-hitachi/hitachi/common/log"
    impl "terraform-provider-hitachi/hitachi/terraform/impl"
    schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceVspSnapshotRange() *schema.Resource {
    return &schema.Resource{
        Description: "Datasource to show snapshot info within a specific P-VOL LDEV range in Hitachi VSP storage.",
        ReadContext: datasourceVspSnapshotRangeRead,
        Schema:      schemaimpl.DatasourceVspSnapshotRangeSchema(),
    }
}

func datasourceVspSnapshotRangeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    return impl.DatasourceVspSnapshotRangeRead(d)
}
