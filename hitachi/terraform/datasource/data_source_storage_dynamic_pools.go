package terraform

import (
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

func DataSourceStorageDynamicPools() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Dynamic Pools: returns all dynamic pools information",
		ReadContext: DataSourceStorageDynamicPoolsRead,
		Schema:      schemaimpl.DataDynamicPoolsSchema,
	}
}

func DataSourceStorageDynamicPoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// fetch all dynamic pool info
	var dynamicPools []terraformmodel.DynamicPool

	dynamicPoolsSource, err := impl.GetDynamicPools(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = copier.Copy(&dynamicPools, dynamicPoolsSource)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return diag.FromErr(err)
	}

	dpList := []map[string]interface{}{}
	for _, dp := range dynamicPools {
		eachDp := impl.ConvertDynamicPoolToSchema(&dp, serial)
		log.WriteDebug("dp: %+v\n", *eachDp)
		dpList = append(dpList, *eachDp)
	}

	if err := d.Set("dynamic_pools", dpList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.WriteInfo("all dynamic pools read successfully")

	return nil
}
