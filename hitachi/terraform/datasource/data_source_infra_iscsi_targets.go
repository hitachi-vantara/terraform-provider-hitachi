package terraform

import (
	"context"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraIscsiTargets() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage iSCSI Targets:The following request obtains information about iSCSI Targets.",
		ReadContext: dataSourceInfraIscsiTargetsRead,
		Schema:      schemaimpl.DataInfraIscsiTargetsSchema,
	}
}

func dataSourceInfraIscsiTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	response, err := impl.GetInfraIscsiTargetsById(d)
	if err != nil {
		return diag.FromErr(err)
	}

	list := []map[string]interface{}{}
	for _, item := range *response {
		eachItem := impl.ConvertInfraIscsiTargetToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("iscsi_targets", list); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("total_iscsi_target_count", len(list))

	log.WriteInfo("iscsiTargets read successfully")
	return nil

}
