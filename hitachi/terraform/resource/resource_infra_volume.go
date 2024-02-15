package terraform

import (
	"context"

	// "time"
	// "errors"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SyncVolumeOperation = &sync.Mutex{}

func ResourceInfraStorageVOlume() *schema.Resource {
	return &schema.Resource{
		Description:   ":meta:subcategory:VSP Storage Volume:The following request creates a volume by using the specified parity groups or pools. Specify a parity group or pool id for creating a basic volume.",
		CreateContext: resourceInfraStorageVolumeCreate,
		ReadContext:   resourceInfraStorageVolumeRead,
		UpdateContext: resourceInfraStorageVolumeUpdate,
		DeleteContext: resourceInfraStorageVolumeDelete,
		Schema:        schemaimpl.ResourceInfraVolumeSchema,
		// CustomizeDiff: customDiffFunc(),
	}
}

func resourceInfraStorageVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	SyncVolumeOperation.Lock() //??
	defer SyncVolumeOperation.Unlock()

	log.WriteInfo("starting volume create")

	volumeInfo, err := impl.CreateInfraVolume(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	volume := impl.ConvertInfraVolumeToSchema(volumeInfo)
	log.WriteDebug("Volume: %+v\n", *volume)
	volList := []map[string]interface{}{
		*volume,
	}
	if err := d.Set("volume", volList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	//d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
	d.SetId(volumeInfo.ResourceId)
	log.WriteInfo("volume created successfully")

	return nil
}

func resourceInfraStorageVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all volumes

	volume, err := impl.GetInfraSingleVolume(d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeSchma := impl.ConvertInfraVolumeToSchema(volume)
	log.WriteDebug("lun: %+v\n", *volumeSchma)

	volList := []map[string]interface{}{
		*volumeSchma,
	}

	if err := d.Set("volume", volList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volume.ResourceId)
	log.WriteInfo("volumes read successfully")

	return nil
}

func resourceInfraStorageVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume update")

	volumeInfo, err := impl.UpdateInfraVolume(d)
	if err != nil {
		return diag.FromErr(err)
	}

	volume := impl.ConvertInfraVolumeToSchema(volumeInfo)
	log.WriteDebug("volume: %+v\n", *volume)
	volList := []map[string]interface{}{
		*volume,
	}
	if err := d.Set("volume", volList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(volumeInfo.ResourceId)
	log.WriteInfo("volume updated successfully")

	return nil
}

func resourceInfraStorageVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume delete")

	err := impl.DeleteInfraVolume(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("volume deleted successfully")
	return nil
}
