package terraform

import (
	"context"
	"strconv"

	// "fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	common "terraform-provider-hitachi/hitachi/terraform/common"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraVolume() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Volumess:The following request obtains information about Volumes.",
		ReadContext: DataSourceInfraVolumeRead,
		Schema:      schemaimpl.DataInfraVolumeSchema,
	}
}

func DataSourceInfraVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	storage_id, _, err := common.GetValidateStorageIDFromSerialResource(d, m)

	if err != nil {

		log.WriteDebug("Error in get storage ID %s", err)

		return diag.FromErr(err)
	}

	if storage_id != nil {

		porcvolume, mtVolume, err := impl.GetInfraVolume(d)
		if err != nil {
			return diag.FromErr(err)
		}

		if mtVolume != nil {

			volumeSchma := impl.ConvertPartnersInfraVolumeToSchema(mtVolume)
			log.WriteDebug("volume: %+v\n", *volumeSchma)

			volList := []map[string]interface{}{
				*volumeSchma,
			}

			if err := d.Set("subscriber_volume", volList); err != nil {
				return diag.FromErr(err)
			}
			d.Set("volume", nil)

			d.SetId(mtVolume.ResourceId)

		} else if porcvolume != nil {
			volumeSchma := impl.ConvertInfraVolumeToSchema(porcvolume)
			log.WriteDebug("volume: %+v\n", *volumeSchma)

			volList := []map[string]interface{}{
				*volumeSchma,
			}

			if err := d.Set("volume", volList); err != nil {
				return diag.FromErr(err)
			}
			d.Set("subscriber_volume", nil)
			d.SetId(porcvolume.ResourceId)

		}

		log.WriteInfo("volume read successfully")

	} else {
		serial := d.Get("serial").(int)

		logicalUnit, err := impl.GetLun(d)
		if err != nil {
			return diag.FromErr(err)
		}

		lun := impl.ConvertLunToSchema(logicalUnit, serial)
		log.WriteDebug("lun: %+v\n", *lun)

		lunList := []map[string]interface{}{
			*lun,
		}
		if err := d.Set("volume", lunList); err != nil {
			return diag.FromErr(err)
		}

		// always run
		// d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.SetId(strconv.Itoa(logicalUnit.LdevID))
		log.WriteInfo("lun read successfully")
	}
	// fetch all volumes
	return nil
}
