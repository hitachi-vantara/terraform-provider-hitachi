package terraform

import (
	"context"
	"errors"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraUcpSystems() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Devices:The following request obtains information about storage devices.",
		ReadContext: DataSourceInfraUcpSystemRead,
		Schema:      schemaimpl.DataInfraUcpSystemSchema,
	}
}

func DataSourceInfraUcpSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial_number").(string)
	name := d.Get("name").(string)

	if serial != "" && name != "" {
		err := errors.New("both serial_number  and name are not allowed. Either serial_number or name or none of them can be specified")
		return diag.FromErr(err)
	}

	// fetch ucp systems
	response, err := impl.GetInfraUcpSystems(d)
	if err != nil {
		return diag.FromErr(err)
	}

	list := []map[string]interface{}{}
	for _, item := range *response {
		eachItem := impl.ConvertInfraUcpSystemToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("systems", list); err != nil {
		return diag.FromErr(err)
	}

	log.WriteDebug("systems: %+v\n", response)
	if serial != "" && name != "" {
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			d.Set("name", element.Name)
			d.Set("serial_number", element.SerialNumber)
			break
		}
	}
	log.WriteInfo("storage devices read successfully")

	return nil

}
