package terraform

import (

	// "time"

	"context"
	"fmt"
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
		Description: "VSP Storage Dynamic Pool:The following request gets information items such as the pool status, the pool usage rate, and the pool threshold.",
		ReadContext: DataSourceStorageDynamicPoolsRead,
		Schema:      schemaimpl.DataDynamicPoolsSchema,
	}
}

func DataSourceStorageDynamicPoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if !okId && !okName {
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

		log.WriteInfo("all dynamic pool read successfully")
	} else {
		var dynamicPool *terraformmodel.DynamicPool
		var err error
		if pid >= 0 {
			// fetch dynamic pool info by pool id

			dynamicPool, err = impl.GetDynamicPoolById(d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if pName != "" {
			// fetch dynamic pool info by pool name

			dynamicPool, err = impl.GetDynamicPoolByName(d)
			if err != nil {
				return diag.FromErr(err)
			}
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

	}

	return nil
}
