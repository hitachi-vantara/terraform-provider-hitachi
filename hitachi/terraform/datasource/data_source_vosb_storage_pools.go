package terraform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"
)

func DataSourceVssbStoragePools() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Pools: Obtains a list of storage pool information.",
		ReadContext: DataSourceVssbStoragePoolsRead,
		Schema:      schemaimpl.DatasourceVssbStoragePoolsSchema,
	}
}

func DataSourceVssbStoragePoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	spNames, ok := d.GetOk("storage_pool_names")
	log.WriteInfo("pool names %+v", spNames)

	if ok { // fetch storage pool by pool names

		storagePools, err := impl.GetStoragePoolsByPoolNames(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storagePools {
			eachSp := impl.ConvertStoragePoolToSchema(&sp)
			log.WriteDebug("sp: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("storage_pools", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("storage pools by pool names read successfully")

		return nil

	} else { // fetch all storage pools
		_, ok := d.GetOk("storage_pool_names")
		if !ok {
			if err := d.Set("storage_pool_names", []string{}); err != nil {
				return diag.FromErr(err)
			}
		}

		storagePools, err := impl.GetAllStoragePools(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storagePools {
			eachSp := impl.ConvertStoragePoolToSchema(&sp)
			log.WriteDebug("sp: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("storage_pools", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("all storage pools read successfully")

		return nil
	}

}
