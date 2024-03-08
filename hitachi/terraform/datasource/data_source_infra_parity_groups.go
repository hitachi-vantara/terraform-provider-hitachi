package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
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

	storage_id,_, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)
	// fetch all storage ports

	if storage_id != nil {
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
	} else {
		serial := d.Get("serial").(int)
		var parityGroups []terraformmodel.ParityGroup

		parityGroupSource, err := impl.GetParityGroups(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = copier.Copy(&parityGroups, parityGroupSource)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return diag.FromErr(err)
		}

		pgList := []map[string]interface{}{}
		for _, pg := range parityGroups {
			eachPg := impl.ConvertParityGroupToSchema(&pg, serial)
			log.WriteDebug("pg: %+v\n", *eachPg)
			pgList = append(pgList, *eachPg)
		}

		if err := d.Set("parity_groups", pgList); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.WriteInfo("storage ports read successfully")

	return nil

}
