package terraform

import (
	"context"
	"errors"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbStoragePorts() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSS Block Storage Port:Obtains a list of storage ports information.",
		ReadContext: DataSourceVssbStoragePortsRead,
		Schema:      schemaimpl.DataVssbStoragePortSchema,
	}
}

func DataSourceVssbStoragePortsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	port_name := d.Get("port_name").(string)

	if port_name == "" {

		storagePorts, err := impl.GetVssbStoragePorts(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storagePorts {
			eachSp := impl.ConvertVssbStoragePortToSchema(&sp)
			log.WriteDebug("storage port: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("ports", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("all vssb storage ports read successfully")

		return nil
	} else {
		if !utils.IsValidPortName(port_name) {
			err := errors.New("port name can not exceed 255 characters.")
			return diag.FromErr(err)
		}

		port, err := impl.GetVssbPort(d)

		if err != nil {
			return diag.FromErr(err)
		}

		pwas := impl.ConvertVssbPortWithAuthSettingsToSchema(port)
		pList := []map[string]interface{}{
			*pwas,
		}
		if err := d.Set("ports", pList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(port.Port.ID)
		log.WriteInfo("port read successfully")
		return nil
	}

}
