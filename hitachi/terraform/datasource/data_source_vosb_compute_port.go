package terraform

import (
	"context"
	"fmt"

	// "fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbComputePort() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Port:Obtains a list of ports information.",
		ReadContext: DataSourceVssbComputePortRead,
		Schema:      schemaimpl.DataSourceVssbComputePortSchema,
	}
}

func DataSourceVssbComputePortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	port := d.Get("name").(string)

	if port == "" {
		err := fmt.Errorf("port is not specified, please specify a port name")
		return diag.FromErr(err)
	} else {
		portInfo, err := impl.GetVssbComputePortByName(d, "", "")
		if err != nil {
			return diag.FromErr(err)
		}

		log.WriteDebug("port Info: %+v\n", *portInfo)
		//sp := impl.ConvertVssbPortDetailSettingsToSchema(portInfo)
		sp := impl.ConvertVssbIscsiPortAuthSettingsToSchema(portInfo)

		log.WriteDebug("port: %+v\n", *sp)

		spList := []map[string]interface{}{
			*sp,
		}

		if err := d.Set("compute_port", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(portInfo.Port.ID)
		log.WriteInfo("port read successfully")
		return nil
	}
}
