package terraform

import (

	// "fmt"

	// "time"

	"context"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageIscsiTarget() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage ISCSI Targets: The following request gets information about iSCSI targets of the port.",
		ReadContext: DataSourceStorageIscsiTargetRead,
		Schema:      schemaimpl.DataIscsiTargetSchema,
	}
}

func DataSourceStorageIscsiTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	iscsiTarget, err := impl.GetIscsiTarget(d)
	if err != nil {
		return diag.FromErr(err)
	}

	it := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
	log.WriteDebug("iscsiTarget: %+v\n", *it)

	itList := []map[string]interface{}{
		*it,
	}

	if err := d.Set("iscsitarget", itList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(iscsiTarget.PortID + strconv.Itoa(iscsiTarget.IscsiTargetNumber))
	log.WriteInfo("iscsiTarget read successfully")

	return nil
}
