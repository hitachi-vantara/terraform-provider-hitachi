package terraform

import (

	// "time"

	"context"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/jinzhu/copier"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageParityGroups() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Parity Group:The following request obtains information about all parity groups.",
		ReadContext: DataSourceStorageParityGroupsRead,
		Schema:      schemaimpl.DataParityGroupsSchema,
	}
}

func DataSourceStorageParityGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

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

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.WriteInfo("all parity group read successfully")

	return nil
}
