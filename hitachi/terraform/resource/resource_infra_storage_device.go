package terraform

import (
	"context"
	"fmt"
	"sync"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconcilermodel "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"

	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncInfraStorageDeviceOperation = &sync.Mutex{}

func ResourceInfraStorageDevice() *schema.Resource {
	return &schema.Resource{
		Description:   `:meta:subcategory:VSP Storage Device:The following request adds a storage device.`,
		CreateContext: resourceInfraStorageDeviceCreate,

		ReadContext:   resourceInfraStorageDeviceRead,
		UpdateContext: resourceInfraStorageDeviceUpdate,
		DeleteContext: resourceInfraStorageDeviceDelete,
		Schema:        schemaimpl.ResourceInfraStorageDeviceSchema,
		CustomizeDiff: ValidateInfraStorageDIff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		//CustomizeDiff: resourceMyResourceCustomDiffInfraHostGroup,
	}
}

func resourceInfraStorageDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncInfraStorageDeviceOperation.Lock() //??
	defer syncInfraStorageDeviceOperation.Unlock()

	log.WriteInfo("starting Infra Storage Device create")

	//serial := d.Get("serial").(int)

	response, err := impl.CreateInfraStorageDevice(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	list := []map[string]interface{}{}
	for _, item := range *response {
		eachItem := impl.ConvertInfraStorageDeviceToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("storage_devices", list); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	for _, item := range *response {
		element := &item
		d.SetId(element.ResourceId)
		/*
			d.Set("hostgroup_name", element.HostGroupName)
			d.Set("hostgroup_number", element.HostGroupId)
			d.Set("port", element.Port)
		*/
		break
	}
	log.WriteInfo("Infra Storage Device created successfully")

	return nil
}

func resourceInfraStorageDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceInfraStorageDevicesRead(ctx, d, m)
}

func resourceInfraStorageDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting Infra Storage Device update")

	//serial := d.Get("serial").(int)

	response, err := impl.UpdateInfraStorageDevice(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	list := []map[string]interface{}{}
	for _, item := range *response {
		eachItem := impl.ConvertInfraStorageDeviceToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("storage_devices", list); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	for _, item := range *response {
		element := &item
		d.SetId(element.ResourceId)
		break
	}

	log.WriteInfo("Infra Storage Device updated successfully")

	return nil
}

func resourceInfraStorageDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting Infra Storage Device delete")

	err := impl.DeleteInfraStorageDevice(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("Infra Storage Device deleted successfully")
	return nil
}

func ValidateInfraStorageDIff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {

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

	system, ok := d.GetOk("system")

	if ok {

		found, _, err := reconObj.FindUcpSystemByName(system.(string))
		if err != nil {
			log.WriteDebug("TFError| error in FindUcpSystemByName, err: %v", err)
			return err
		}
		if !(*found) {
			log.WriteDebug("TFError| provided system does not exist %s", system.(string))
			return fmt.Errorf("provided system does not exist %s", system.(string))
		}
	}

	return nil
}
