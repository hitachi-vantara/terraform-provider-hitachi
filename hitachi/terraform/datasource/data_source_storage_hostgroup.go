package terraform

import (
	"context"
	"fmt"
	// "strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageHostGroup() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Host Group: The following request gets information about host group of the port.",
		ReadContext: DataSourceStorageHostGroupRead,
		Schema:      schemaimpl.DataHostGroupSchema,
	}
}

func DataSourceStorageHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	hostgroup, err := impl.GetHostGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	hg := impl.ConvertHostGroupToSchema(hostgroup, serial)
	log.WriteDebug("hg: %+v\n", *hg)

	hgList := []map[string]interface{}{
		*hg,
	}

	if err := d.Set("hostgroup", hgList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s,%d,%s", hostgroup.PortID, hostgroup.HostGroupNumber, hostgroup.HostGroupName))
	log.WriteInfo("hg read successfully")

	return nil
}
