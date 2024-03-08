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
	storage_id, _, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {

		// fetch all volumes

		volumes, err := impl.GetInfraVolumes(d)

		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *volumes {
			eachSp := impl.ConvertPartnersInfraVolumeToSchema(&sp)
			log.WriteDebug("it: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("volumes", spList); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("volumes read successfully")

	} else {
		serial := d.Get("serial").(int)

		logicalUnits, err := impl.GetRangeOfLuns(d)
		if err != nil {
			return diag.FromErr(err)
		}

		lunList := []map[string]interface{}{}

		for _, lun := range *logicalUnits {
			eachLun := impl.ConvertLunToSchema(&lun, serial)
			log.WriteDebug("eachLun: %+v\n", &eachLun)
			lunList = append(lunList, *eachLun)
		}

		if err := d.Set("volumes", lunList); err != nil {
			return diag.FromErr(err)
		}

		// always run
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		//d.SetId(strconv.Itoa(logicalUnits.LdevID))
		log.WriteInfo("range of luns read successfully")

	}
	return nil
}
