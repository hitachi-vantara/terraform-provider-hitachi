package terraform

import (
	"context"
	"fmt"
	"strings"

	// "time"
	// "errors"

	"sync"
	cache "terraform-provider-hitachi/hitachi/common/cache"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncCreateVolOpertation = &sync.Mutex{}

func ResourceVssbStorageCreateVolume() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block Compute Node:Creates a volume.",
		CreateContext: resourceCreateVolume,
		Schema:        schemaimpl.ResourceVolumeSchema,
		ReadContext:   resourceReadVolume,
		UpdateContext: resourceCreateVolume,
		DeleteContext: resourceDeleteVolume,
		CustomizeDiff: resourceMyResourceCustomDiff,
	}
}

func resourceCreateVolume(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncCreateVolOpertation.Lock()
	defer syncCreateVolOpertation.Unlock()

	log.WriteInfo("starting volume create")

	volumeData, err := impl.CreateVolume(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	volume := impl.ConvertVssbVolumesToSchema(volumeData)
	log.WriteDebug("volume: %+v\n", *volume)
	volList := []map[string]interface{}{
		*volume,
	}
	if err := d.Set("volume", volList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(volumeData.ID)
	log.WriteInfo("volume created successfully")

	return nil
}

func resourceDeleteVolume(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncCreateVolOpertation.Lock()
	defer syncCreateVolOpertation.Unlock()

	log.WriteInfo("starting volume Delete")

	err := impl.DeleteVolume(d)
	if err != nil {
		if err.Error() == "The request could not be executed." {
			log.WriteDebug("TFError| error deleting volume, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED_MSG))
			err = fmt.Errorf(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED_MSG))
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}
	return nil
}

func resourceReadVolume(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbVolumeNodesRead(ctx, d, m)
}

func resourceMyResourceCustomDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	vssbAddr := d.Get("vosb_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in Reconciler NewEx, err: %v", err)
		return err
	}

	name, ok := d.GetOk("name")
	if !ok {

		log.WriteDebug("name: %s", name.(string))
		return fmt.Errorf("name is required")
	}
	capacity, capacityOk := d.GetOk("capacity_gb")
	var capacityMB int
	if capacityOk {
		capacityMB = int(capacity.(float64) * 1024)
	}
	poolName, storagePoolOk := d.GetOk("storage_pool")
	volname := name.(string)
	volData, volOk := reconObj.GetVolumeDetails(volname)
	fmt.Printf("Existing volume found, going for update functionality: %s\n", volname)
	if volOk != nil {
		notpresent := []string{}
		if !capacityOk {
			notpresent = append(notpresent, "capacity_gb")
		}
		if !storagePoolOk {
			notpresent = append(notpresent, "storage_pool")
		}

		if len(notpresent) > 0 {
			return fmt.Errorf("parameters are required for the new volume creation:  %v", strings.Join(notpresent, ", "))

		}

	}

	poolDetails, err := reconObj.GetStoragePoolByPoolName(poolName.(string))
	if err != nil {
		return fmt.Errorf("storage_pool is details not found. Provide correct pool name")
	}
	if volOk == nil {
		if poolDetails.ID != volData.PoolId {
			return fmt.Errorf("can't change the pool once volume is created, Provide the correct pool name which volume was created with")
		}
		if capacityMB < volData.TotalCapacity {
			return fmt.Errorf("volume capacity can't be less than existing volume capacity")

		}
	}
	computeNodes := d.Get("compute_nodes")
	computeNodeCheck := d.GetRawConfig().GetAttr("compute_nodes").IsNull()
	if !computeNodeCheck {
		noNodes := []string{}
		for _, node := range computeNodes.([]interface{}) {
			_, err := reconObj.GetComputeNodeInformationByName(node.(string), "")
			if err != nil {
				noNodes = append(noNodes, node.(string))
			}
		}
		if len(noNodes) > 0 {
			return fmt.Errorf("no compute node found for then given compute node names: %s", strings.Join(noNodes, ", "))
		}
	}

	log.WriteDebug("T: %v", volData)
	return nil
}
