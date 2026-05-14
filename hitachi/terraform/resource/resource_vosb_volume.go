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
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncCreateVolOpertation = &sync.Mutex{}

func ResourceVssbStorageCreateVolume() *schema.Resource {
	return &schema.Resource{
		Description: "VSP One SDS Block Volume: CRUD operations of a volume.",
		Importer: &schema.ResourceImporter{
			StateContext: importVosbVolume,
		},
		CreateContext: resourceCreateVolume,
		Schema:        schemaimpl.ResourceVolumeSchema,
		ReadContext:   resourceReadVolume,
		UpdateContext: resourceCreateVolume,
		DeleteContext: resourceDeleteVolume,
		CustomizeDiff: resourceMyResourceCustomDiff,
	}
}

func validateVssbVolumeComputeNodes(d *schema.ResourceData) error {
	computeNodesRaw, ok := d.GetOk("compute_nodes")
	if !ok {
		return nil
	}
	computeNodes, ok := computeNodesRaw.([]interface{})
	if !ok || len(computeNodes) == 0 {
		return nil
	}

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
		return err
	}

	noNodes := []string{}
	for _, node := range computeNodes {
		nodeName, ok := node.(string)
		if !ok {
			continue
		}
		_, err := reconObj.GetComputeNodeInformationByName(nodeName, "")
		if err != nil {
			noNodes = append(noNodes, nodeName)
		}
	}
	if len(noNodes) > 0 {
		return fmt.Errorf("no compute node found for then given compute node names: %s", strings.Join(noNodes, ", "))
	}
	return nil
}

func resourceCreateVolume(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncCreateVolOpertation.Lock()
	defer syncCreateVolOpertation.Unlock()

	if d.Id() == "" {
		log.WriteInfo("starting volume create")
	} else {
		log.WriteInfo("starting volume update")
	}
	if err := validateVssbVolumeComputeNodes(d); err != nil {
		return diag.FromErr(err)
	}

	volumeData, err := impl.CreateVolume(d)
	if err != nil {
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

	if err := d.Set("volume", nil); err != nil { // clear old state first
		return diag.FromErr(err)
	}
	if err := d.Set("volume", volList); err != nil {
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
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	volumeNode, err := impl.GetVssbVolumeNode(d)
	if err != nil {
		// Treat not-found as a removed resource.
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	volume := impl.ConvertVssbVolumesToSchema(volumeNode)
	log.WriteDebug("vol: %+v\n", *volume)
	volumeList := []map[string]interface{}{
		*volume,
	}
	if err := d.Set("volume", volumeList); err != nil {
		return diag.FromErr(err)
	}
	// Keep selectors populated in state for post-import updates.
	if volumeNode != nil {
		_ = d.Set("name", volumeNode.Name)
	}

	d.SetId(volumeNode.ID)
	log.WriteInfo("all vssb volume read successfully")
	return nil
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

	// Only do deterministic schema validations here (no API calls).

	currentID := d.Id()
	isCreate := currentID == ""
	isUpdate := !isCreate

	// Enforce create-only required fields.
	// (Schema keeps them Optional so imports/updates don't require respecifying them.)
	if isCreate {
		if v, ok := d.GetOk("vosb_address"); !ok || strings.TrimSpace(v.(string)) == "" {
			return fmt.Errorf("vosb_address is required for create")
		}
		if v, ok := d.GetOk("name"); !ok || strings.TrimSpace(v.(string)) == "" {
			return fmt.Errorf("name is required for create")
		}
		if v, ok := d.GetOk("storage_pool"); !ok || strings.TrimSpace(v.(string)) == "" {
			return fmt.Errorf("storage_pool is required for create")
		}
		cap, ok := d.GetOk("capacity_gb")
		if !ok {
			return fmt.Errorf("capacity_gb is required for create")
		}
		if cap.(float64) <= 0 {
			return fmt.Errorf("capacity_gb must be greater than zero")
		}
	}

	// Name is required for update/apply as well (backend reconciler dereferences it).
	if v, ok := d.GetOk("name"); !ok || strings.TrimSpace(v.(string)) == "" {
		return fmt.Errorf("name must be specified (set it in config or import the resource so it is populated in state)")
	}

	// storage_pool is create-only (pool migration is not supported here)
	if isUpdate && d.HasChange("storage_pool") {
		return fmt.Errorf("storage_pool is create-only and cannot be changed")
	}

	// Prevent capacity shrink; allow expand.
	if d.HasChange("capacity_gb") {
		oldRaw, newRaw := d.GetChange("capacity_gb")
		oldVal, okOld := oldRaw.(float64)
		newVal, okNew := newRaw.(float64)
		if okNew && newVal <= 0 {
			return fmt.Errorf("capacity_gb must be greater than zero")
		}
		if okOld && okNew && oldVal > 0 && newVal < oldVal {
			return fmt.Errorf("capacity_gb cannot be decreased (old %.2f, new %.2f)", oldVal, newVal)
		}
	}

	// Mutually exclusive controller selectors (kept deterministic)
	if _, scOk := d.GetOk("storage_controller_name"); scOk {
		if _, fdOk := d.GetOk("fault_domain_id"); fdOk {
			return fmt.Errorf("storage_controller_name and fault_domain_id are mutually exclusive")
		}
	}

	// fix for 'volume' not updated in console output
	if err := d.SetNewComputed("volume"); err != nil {
		return err
	}

	return nil
}
