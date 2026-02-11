package terraform

import (
    "context"

    commonlog "terraform-provider-hitachi/hitachi/common/log"
    impl "terraform-provider-hitachi/hitachi/terraform/impl"
    schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceVspSnapshotFamily() *schema.Resource {
    return &schema.Resource{
        Description: "Datasource to show the snapshot family (parents and cascaded clones) for a specific LDEV in Hitachi VSP storage.",
        ReadContext: datasourceVspSnapshotFamilyRead,
        Schema:      schemaimpl.DatasourceVspSnapshotFamilySchema(),
    }
}

func datasourceVspSnapshotFamilyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    log := commonlog.GetLogger()
    log.WriteEnter()
    defer log.WriteExit()

    return impl.DatasourceVspSnapshotFamilyRead(d)
}