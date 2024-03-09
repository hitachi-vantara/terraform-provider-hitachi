package terraform

import (
	"context"

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

func DataSourceInfraStorageDevices() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Devices:The following request obtains information about storage devices.",
		ReadContext: DataSourceInfraStorageDevicesRead,
		Schema:      schemaimpl.DataInfraStorageDevicesSchema,
	}
}

func GetSerialString(d *schema.ResourceData) string {
	serial_number := -1
	serial := ""

	sid, okId := d.GetOk("serial")
	if okId {
		serial_number = sid.(int)
		serial = strconv.Itoa(serial_number)
	}
	return serial
}

func DataSourceInfraStorageDevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := GetSerialString(d)

	// fetch all storage devices

	var response *[]terraform.InfraStorageDeviceInfo

	var err error
	list := []map[string]interface{}{}

	if serial == "" {
		var mtResponse *terraform.InfraMTStorageDevices
		response, mtResponse, err = impl.GetInfraStorageDevices(d)
		if err != nil {
			return diag.FromErr(err)
		}
		if mtResponse == nil {
			for _, item := range *response {
				eachItem := impl.ConvertInfraStorageDeviceToSchema(&item)
				log.WriteDebug("it: %+v\n", *eachItem)
				list = append(list, *eachItem)
			}
			if err := d.Set("storage_devices", list); err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, item := range mtResponse.Data {
				eachItem := impl.ConvertPartnersInfraStorageDeviceToSchema(&item)
				log.WriteDebug("it: %+v\n", *eachItem)
				list = append(list, *eachItem)
			}
			if err := d.Set("partner_storage_devices", list); err != nil {
				return diag.FromErr(err)
			}
		}

		log.WriteDebug("storageDevices: %+v\n", response)
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	} else {
		var mtResponse *terraform.InfraMTStorageDevice
		response, mtResponse, err = impl.GetInfraStorageDevice(d, serial)
		if err != nil {
			return diag.FromErr(err)
		}
		if mtResponse == nil {
			for _, item := range *response {
				eachItem := impl.ConvertInfraStorageDeviceToSchema(&item)
				log.WriteDebug("it: %+v\n", *eachItem)
				list = append(list, *eachItem)
			}
			if err := d.Set("storage_devices", list); err != nil {
				return diag.FromErr(err)
			}
			for _, item := range *response {
				element := &item
				d.SetId(element.ResourceId)
				break
			}
		} else {
			eachItem := impl.ConvertPartnersInfraStorageDeviceToSchema(mtResponse)
			list = append(list, *eachItem)
			if err := d.Set("partner_storage_devices", list); err != nil {
				return diag.FromErr(err)
			}
			d.SetId(mtResponse.Storage.ResourceId)
		}
	}

	return nil

}
