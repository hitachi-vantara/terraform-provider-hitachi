package terraform

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	// "time"
	// "errors"
	"sync"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"

	reconcilermodel "terraform-provider-hitachi/hitachi/infra_gw/model"
	common "terraform-provider-hitachi/hitachi/terraform/common"

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
		CustomizeDiff: InfraVolumeDIffValidate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceInfraStorageVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	SyncVolumeOperation.Lock() //??
	defer SyncVolumeOperation.Unlock()

	_, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

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

func InfraVolumeDIffValidate(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialStringFromDiff(d)
	storageId := d.Get("storage_id").(string)

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
		return err
	} else if serial == "" && storageId == "" { 
		err := errors.New("either serial or storage_id can't be empty. Please specify one")
        return err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return err
		}

	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return err
	}

	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	if setting.PartnerId != nil {
		subId, ok := d.GetOk("subscriber_id")
		if ok {
			subIdw := subId.(string)
			setting.SubscriberId = &subIdw
		}
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}
	var volok bool = true

	name, _ := d.GetOk("name")

	if d.Id() == "" {
		_, volok = reconObj.GetVolumeByName(storageId, name.(string))
	}

	mandatoryFields := []string{"pool_id", "parity_group_id", "capacity", "system"}
	missingFields := []string{}

	for _, field := range mandatoryFields {
		_, ok := d.GetOk(field)
		if !ok && !volok {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("mandatory fields missing for new volume creation: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
