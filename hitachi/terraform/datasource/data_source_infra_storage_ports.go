package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraStoragePorts() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Ports:The following request obtains information about ports.",
		ReadContext: dataSourceInfraStoragePortsRead,
		Schema:      schemaimpl.DataInfraStoragePortsSchema,
	}
}

func dataSourceInfraStoragePortsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storage_id, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	// fetch all storage ports
	port_id := d.Get("port_id").(string)

	if storage_id != nil {

		response, mtResponse, err := impl.GetInfraStoragePorts(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}

		if mtResponse == nil {
			for _, sp := range *response {
				eachSp := impl.ConvertInfraStoragePortToSchema(&sp)
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
		} else {
			for _, sp := range *mtResponse {
				eachSp := impl.ConvertInfraMTStoragePortToSchema(&sp)
				log.WriteDebug("it: %+v\n", *eachSp)
				spList = append(spList, *eachSp)
			}

			if err := d.Set("partner_ports", spList); err != nil {
				return diag.FromErr(err)
			}
			d.Set("total_port_count", len(spList))
			if port_id == "" {
				d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
			} else {
				for _, item := range *mtResponse {
					element := &item
					d.SetId(element.ResourceId)
					d.Set("resource_id", element.ResourceId)
					d.Set("port_id", element.PortInfo.PortId)
					break
				}
			}
		}
	} else {
		if port_id == "" {
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
		}
	}
	log.WriteInfo("storage port read successfully")
	return nil

}
