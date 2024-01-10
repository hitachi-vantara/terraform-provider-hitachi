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

func DataSourceInfraParityGroups() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Parity Groups:The following request obtains information about Parity Groups.",
		ReadContext: dataSourceInfraParityGroupsRead,
		Schema:      schemaimpl.DataInfraParityGroupsSchema,
	}
}

func dataSourceInfraParityGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all storage ports

	parityGroups, err := impl.GetInfraParityGroups(d)
	if err != nil {
		return diag.FromErr(err)
	}

	spList := []map[string]interface{}{}
	for _, sp := range *parityGroups {
		eachSp := impl.ConvertInfraGwParityGroupToSchema(&sp)
		log.WriteDebug("it: %+v\n", *eachSp)
		spList = append(spList, *eachSp)
	}

	if err := d.Set("parity_groups", spList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.WriteInfo("storage ports read successfully")

	return nil

}
