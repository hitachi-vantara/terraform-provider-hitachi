package terraform

import (
	"context"
	"encoding/json"
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
		Description:   "VSP One SDS Block Volume: CRUD operations of a volume.",
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

	_, ok := d.GetOk("compute_nodes")
	if !ok {
		if err := d.Set("compute_nodes", []string{}); err != nil { // needed to get rid of '/* of string */' in the output
			return diag.FromErr(err)
		}
	}

	d.Set("volume", nil) // clear old state first
	if err := d.Set("volume", volList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	// verify
	volData := d.Get("volume")
	volJSON, err := json.MarshalIndent(volData, "", "  ")
	if err != nil {
		log.WriteDebug("[ERROR] Failed to marshal volume data: %s", err)
	} else {
		log.WriteDebug("[DEBUG] Volume data: %s", string(volJSON))
	}

	d.SetId("")
	d.SetId(volumeData.ID)
	log.WriteInfo("volume created successfully")

	// // Always refresh the resource state
	// resourceReadVolume(ctx, d, m)

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
		if strings.Contains(err.Error(), "The request could not be executed") {
			log.WriteDebug("TFError| error deleting volume, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED_MSG))
			err = fmt.Errorf("%v", mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED_MSG))
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}
	return nil
}

func resourceReadVolume(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbVolumeNodesRead(ctx, d, m)
}

// CustomDiff is intended for schema-based validations only.
// Terraform expects CustomDiff to be deterministic, stateless, and side-effect free.
// Avoid making API calls here, as they may introduce performance issues,
// non-deterministic behavior, or failures during `terraform plan`.
//
// All backend validations (e.g., checking if volume exists, validating compute node names, etc.)
// should be moved to the resource's Create/Update functions instead.
//
// See: https://developer.hashicorp.com/terraform/plugin/framework/resources/customize-diff
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

	// fix for 'volume' not updated in console output
	d.SetNewComputed("volume")

	return nil
}
