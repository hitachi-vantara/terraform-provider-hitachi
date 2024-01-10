package terraform

import (
	"context"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//terraform "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraHostGroups() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Host Group:The following request gets information about host groups of the ports.",
		ReadContext: dataSourceInfraHostGroupsRead,
		Schema:      schemaimpl.DataInfraHostGroupsSchema,
	}
}

func dataSourceInfraHostGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	response, err := impl.GetInfraHostGroupsByPortIds(d)
	if err != nil {
		return diag.FromErr(err)
	}

	list := []map[string]interface{}{}
	for _, item := range *response {
		eachItem := impl.ConvertInfraHostGroupToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("hostgroups", list); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("total_hostgroup_count", len(list))

	log.WriteInfo("host groups read successfully")
	return nil

}
