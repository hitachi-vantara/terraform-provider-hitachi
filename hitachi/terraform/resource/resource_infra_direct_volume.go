package terraform

import (
	"context"
	"fmt"
	"strconv"
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
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	// datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

var SyncVolumeOperation = &sync.Mutex{}

func ResourceInfraStorageVolume() *schema.Resource {
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

	storage_id, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {
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
	} else {
		serial := d.Get("serial").(int)

		logicalUnit, err := impl.CreateLun(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		lun := impl.ConvertLunToSchema(logicalUnit, serial)
		log.WriteDebug("lun: %+v\n", *lun)
		lunList := []map[string]interface{}{
			*lun,
		}
		if err := d.Set("volume", lunList); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		//d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
		d.SetId(strconv.Itoa(logicalUnit.LdevID))
		log.WriteInfo("lun created successfully")
	}
	log.WriteInfo("volume created successfully")

	return nil
}

func resourceInfraStorageVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// fetch all volumes
	storage_id, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)
	if storage_id != nil {

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
	} else {
		return datasourceimpl.DataSourceStorageLunRead(ctx, d, m)
	}
	log.WriteInfo("volumes read successfully")

	return nil
}

func resourceInfraStorageVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume update")
	storage_id, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {
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
	} else {
		serial := d.Get("serial").(int)

		logicalUnit, err := impl.UpdateLun(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		terraformModelLun := terraformmodel.LogicalUnit{}
		err = copier.Copy(&terraformModelLun, logicalUnit)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil
		}

		lun := impl.ConvertLunToSchema(&terraformModelLun, serial)
		log.WriteDebug("lun: %+v\n", *lun)
		lunList := []map[string]interface{}{
			*lun,
		}
		if err := d.Set("volume", lunList); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
		d.SetId(strconv.Itoa(logicalUnit.LdevID))
	}

	log.WriteInfo("volume updated successfully")

	return nil
}

func resourceInfraStorageVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume delete")
	storage_id, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)
	if storage_id != nil {
		err := impl.DeleteInfraVolume(d)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		log.WriteInfo("volume deleted successfully")
		return nil
	} else {
		err := impl.DeleteLun(d)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		log.WriteInfo("lun deleted successfully")
		return nil
	}

}

func InfraVolumeDIffValidate(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// providerConfig := d.Get("provider_config")
	// var storageId string

	storageSettings := m.(*terraformmodel.AllStorageTypes)
	directSerial := d.Get("serial").(int)

	sanStorageSetting, sanError := cache.GetSanSettingsFromCache(strconv.Itoa(directSerial))

	storageId, infraError := common.GetValidateStorageIDFromSerial(d)

	if sanError != nil && len(storageSettings.VspStorageSystem) > 0 {
		return sanError
	} else if infraError != nil && len(storageSettings.InfraGwInfo) > 0 {
		return infraError
	}

	if storageId != nil && sanStorageSetting != nil {
		return fmt.Errorf("found Same Serial number in both the provider %v", directSerial)
	}

	if storageId != nil {

		// The resource is configured with hitachi_infrastructure_gateway_provider
		return ValidateInfraVolumeDIff(d, *storageId)

	} else if sanStorageSetting != nil {
		return ValidateSanStorageVolumeDIff(d)
	}
	return nil
}

func ValidateInfraVolumeDIff(d *schema.ResourceDiff, storageId string) error {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return err
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
	var volOk bool = true

	name, _ := d.GetOk("name")

	if d.Id() == "" {
		_, volOk = reconObj.GetVolumeByName(storageId, name.(string))
	}

	mandatoryIntFields := []string{"pool_id"}
	missingFields := []string{}

	for _, field := range mandatoryIntFields {
		data, _ := d.GetOk(field)
		if data.(int) == -1 && !volOk {
			missingFields = append(missingFields, field)
		}
	}

	size, _ := d.GetOk("size_gb")

	if size.(float64) <= 0 && !volOk {
		missingFields = append(missingFields, "size_gb")
	}
	system, ok := d.GetOk("system")

	if ok {
		found, _, err := reconObj.FindUcpSystemByName(system.(string))
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		if !(*found) {
			return fmt.Errorf("provided system does not exist %s", system.(string))
		}
	}

	serial := common.GetSerialStringFromDiff(d)

	if serial != "" && !volOk {

		_, err := reconObj.FindStorageSystemByNameAndSerial(reconcilermodel.DefaultSystemName, serial)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("mandatory fields missing for new volume creation: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
