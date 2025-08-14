package terraform

import (
	"context"
	"fmt"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageDynamicPool() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Dynamic Pool: returns specific dynamic pool information",
		ReadContext: DataSourceStorageDynamicPoolRead,
		Schema:      schemaimpl.DataDynamicPoolSchema,
	}
}

func DataSourceStorageDynamicPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	pid := 0
	poolId, okId := d.GetOk("pool_id")
	if okId {
		pid = poolId.(int)
	}

	pName := ""
	poolName, okName := d.GetOk("pool_name")
	if okName {
		pName = poolName.(string)
	}

	if okId == okName {
		err := fmt.Errorf("either pool_id or pool_name is required for dynamic pool datasource")
		return diag.FromErr(err)

	}

	var dynamicPool *terraformmodel.DynamicPool
	var err error
	if okId && pid >= 0 {
		// fetch dynamic pool info by pool id

		dynamicPool, err = impl.GetDynamicPoolById(d)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("pool_name", dynamicPool.PoolName)
	}
	if pName != "" {
		// fetch dynamic pool info by pool name

		dynamicPool, err = impl.GetDynamicPoolByName(d)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("pool_id", dynamicPool.PoolID)
	}

	dpList := []map[string]interface{}{}

	dp := impl.ConvertDynamicPoolToSchema(dynamicPool, serial)
	log.WriteDebug("dp: %+v\n", *dp)
	dpList = append(dpList, *dp)

	if err := d.Set("dynamic_pools", dpList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(dynamicPool.PoolID), 10))

	log.WriteInfo("dynamic pool read successfully")

	return nil
}
