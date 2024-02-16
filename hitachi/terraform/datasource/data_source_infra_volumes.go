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

func DataSourceInfraVolumes() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Parity Groups:The following request obtains information about Parity Groups.",
		ReadContext: DataSourceInfraVolumesRead,
		Schema:      schemaimpl.DataInfraVolumesSchema,
	}
}

func DataSourceInfraVolumesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all volumes

	volumes, mtVolumes, err := impl.GetInfraVolumes(d)

	if err != nil {
		return diag.FromErr(err)
	}

	spList := []map[string]interface{}{}
	if mtVolumes == nil {
		for _, sp := range *volumes {
			eachSp := impl.ConvertInfraVolumeToSchema(&sp)
			log.WriteDebug("it: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("volumes", spList); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("volumes read successfully")

	} else {
		for _, sp := range *mtVolumes {
			eachSp := impl.ConvertPartnersInfraVolumeToSchema(&sp)
			log.WriteDebug("it: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("partner_volumes", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("volumes read successfully")
	}

	return nil

}
