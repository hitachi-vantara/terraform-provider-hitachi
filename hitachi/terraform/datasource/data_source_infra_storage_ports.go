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

func DataSourceInfraStoragePorts() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Ports:The following request obtains information about ports.",
		ReadContext: dataSourceInfraStoragePortsRead,
		Schema:      schemaimpl.InfraStoragePortsSchema,
	}
}

func dataSourceInfraStoragePortsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all storage ports
	port_id := d.Get("port_id").(string)

	response, err := impl.GetInfraGwStoragePorts(d)
	if err != nil {
		return diag.FromErr(err)
	}

	spList := []map[string]interface{}{}
	for _, sp := range *response {
		eachSp := impl.ConvertInfraGwStoragePortToSchema(&sp)
		log.WriteDebug("it: %+v\n", *eachSp)
		spList = append(spList, *eachSp)
	}

	if err := d.Set("ports", spList); err != nil {
		return diag.FromErr(err)
	}
	d.Set("total_port_count", len(spList))
	if port_id == "" {
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			d.Set("resource_id", element.ResourceId)
			d.Set("port_id", element.PortId)
			break
		}
	}
	log.WriteInfo("storage ports read successfully")

	return nil

}
