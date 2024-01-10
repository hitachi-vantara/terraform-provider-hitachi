package terraform

import (
	"context"
	"errors"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraform "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraStoragePools() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Pools:The following request obtains information about pools.",
		ReadContext: dataSourceInfraStoragePoolsRead,
		Schema:      schemaimpl.DataInfraStoragePoolsSchema,
	}
}

func dataSourceInfraStoragePoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all storage ports
	pool_id := -1

	pool_name := d.Get("pool_name").(string)
	pid, okId := d.GetOk("pool_id")
	if okId {
		pool_id = pid.(int)
	}

	if pool_name != "" && pool_id != -1 {
		err := errors.New("both name  and pool_id are not allowed. Either name or pool_id or none of them can be specified")
		return diag.FromErr(err)
	}

	var response *[]terraform.InfraStoragePoolInfo
	var err error
	list := []map[string]interface{}{}

	response, err = impl.GetInfraGwStoragePools(d)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, item := range *response {
		eachItem := impl.ConvertInfraGwStoragePoolToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("storage_pools", list); err != nil {
		return diag.FromErr(err)
	}

	log.WriteDebug("storageDevices: %+v\n", response)
	if pool_name == "" && pool_id == -1 {
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			d.Set("resource_id", element.ResourceId)
			d.Set("pool_name", element.Name)
			d.Set("pool_id", element.PoolId)
			break
		}
	}
	log.WriteInfo("storage pools read successfully")

	return nil

}
