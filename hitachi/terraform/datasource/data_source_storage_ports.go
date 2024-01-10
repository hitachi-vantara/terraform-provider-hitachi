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

func DataSourceStoragePorts() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Ports:The following request obtains information about ports.",
		ReadContext: dataSourceStoragePortsRead,
		Schema:      schemaimpl.StoragePortsSchema,
	}
}

func dataSourceStoragePortsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// serial := d.Get("serial").(int)
	portId := d.Get("port_id").(string)

	// fetch all storage ports
	if portId == "" {
		storagePorts, err := impl.GetStoragePorts(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storagePorts {
			eachSp := impl.ConvertStoragePortToSchema(&sp)
			log.WriteDebug("it: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}
		if err := d.Set("total_port_count", len(spList)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("ports", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("storage ports read successfully")

		return nil

	} else { // fetch port by portId
		storagePort, err := impl.GetStoragePortByPortId(d)
		if err != nil {
			return diag.FromErr(err)
		}

		sp := impl.ConvertStoragePortToSchema(storagePort)
		log.WriteDebug("storage port: %+v\n", *sp)

		spList := []map[string]interface{}{
			*sp,
		}

		if err := d.Set("total_port_count", len(spList)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ports", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(storagePort.PortId)
		log.WriteInfo("storage port read successfully")

		return nil
	}

}
