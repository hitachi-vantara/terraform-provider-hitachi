package terraform

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// utils "terraform-provider-hitachi/hitachi/common/utils"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncLunOperation = &sync.Mutex{}

func ResourceStorageLun() *schema.Resource {
	return &schema.Resource{
		Description:   `VSP Storage Volume: The following request creates a volume by using the specified parity groups or pools. Specify a parity group or pool id for creating a basic volume.`,
		CreateContext: resourceStorageLunCreate,
		ReadContext:   resourceStorageLunRead,
		UpdateContext: resourceStorageLunUpdate,
		DeleteContext: resourceStorageLunDelete,
		Schema:        schemaimpl.ResourceLunSchema,
		CustomizeDiff: resourceStorageLunValidation,
	}
}

func resourceStorageLunCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncLunOperation.Lock() //??
	defer syncLunOperation.Unlock()

	log.WriteInfo("starting lun create")

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

	return nil
}

func resourceStorageLunRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceStorageLunRead(ctx, d, m)
}

func resourceStorageLunUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting lun update")

	serial := d.Get("serial").(int)

	logicalUnit, err := impl.UpdateLun(d)
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

	d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
	d.SetId(strconv.Itoa(logicalUnit.LdevID))
	log.WriteInfo("lun updated successfully")

	return nil
}

func resourceStorageLunDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting lun delete")

	err := impl.DeleteLun(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("lun deleted successfully")
	return nil
}

func resourceStorageLunValidation(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	currentID := d.Id()
	isCreate := currentID == ""
	// isUpdate := !isCreate

	// Prevent LDEV ID change during update
	if currentID != "" {
		// Resource is in UPDATE mode

		ldevRaw, exists := d.GetOk("ldev_id")
		if exists {
			newLdevID, ok := ldevRaw.(int)
			if !ok {
				return fmt.Errorf("invalid type for ldev_id")
			}

			// Convert the resource ID into int to compare
			idFromState, err := strconv.Atoi(currentID)
			if err != nil {
				return fmt.Errorf("failed to parse current resource ID (%s): %v", currentID, err)
			}

			// LDEV ID mismatch check
			if newLdevID != idFromState {
				return fmt.Errorf(
					"ldev_id (%d) does not match the provisioned LDEV ID (%d). "+
						"Hitachi LDEVs cannot be reassigned after provisioning.",
					newLdevID, idFromState,
				)
			}
		}
	}

	// -----------------------------------------------------
	// Validate size_gb input
	// -----------------------------------------------------
	sizeRaw := d.Get("size_gb")
	if sizeRaw == nil {
		return fmt.Errorf("size_gb is required")
	}

	newSizeGB, ok := sizeRaw.(float64)
	if !ok {
		return fmt.Errorf("invalid type for size_gb")
	}
	if newSizeGB <= 0 {
		return fmt.Errorf("size_gb must be greater than zero, got %.2f", newSizeGB)
	}

	// Prevent size decrease
	oldRaw, newRaw := d.GetChange("size_gb")
	if oldRaw != nil && newRaw != nil {
		oldVal, okOld := oldRaw.(float64)
		newVal, okNew := newRaw.(float64)
		if okOld && okNew && newVal < oldVal {
			return fmt.Errorf("new size_gb value must be >= old value: %.2f", oldVal)
		}
	}

	// Validate against real array capacity (state field)
	if volSizeMB, ok := d.Get("total_capacity_in_mb").(int); ok && volSizeMB > 0 {
		actualSizeGB := float64(volSizeMB) / 1024.0
		if newSizeGB < actualSizeGB {
			return fmt.Errorf(
				"requested size_gb (%.2f GB) is smaller than the provisioned size (%.2f GB). "+
					"Hitachi volumes cannot be shrunk.",
				newSizeGB, actualSizeGB,
			)
		}
	}

	// -----------------------------------------------------
	// Pool or Parity Group validation
	// -----------------------------------------------------

	// pool_id: default -1 means "not provided"
	poolID := -1
	if v, ok := d.GetOk("pool_id"); ok {
		poolID = v.(int)
	}
	hasPoolID := poolID >= 0

	// pool_name
	poolName, hasPoolName := d.Get("pool_name").(string)
	if hasPoolName && poolName == "" {
		hasPoolName = false
	}

	// paritygroup_id
	parityGroupID, hasParityGroup := d.Get("paritygroup_id").(string)
	if hasParityGroup && parityGroupID == "" {
		hasParityGroup = false
	}

	// external_paritygroup_id
	externalParityGroupID, hasExternalParityGroup := d.Get("external_paritygroup_id").(string)
	if hasExternalParityGroup && externalParityGroupID == "" {
		hasExternalParityGroup = false
	}

	// exactly one of these must be specified
	count := 0
	if hasPoolID {
		count++
	}
	if hasPoolName {
		count++
	}
	if hasParityGroup {
		count++
	}

	if hasExternalParityGroup {
		count++
	}

	if count != 1 {
		return fmt.Errorf("exactly one of pool_id, pool_name, paritygroup_id, or external_paritygroup_id must be specified")
	}

	isDpPool := hasPoolID || hasPoolName
	isPG := hasParityGroup || hasExternalParityGroup

	// Flags
	capacitySaving := d.Get("capacity_saving").(string)
	isShareEnabled := d.Get("is_data_reduction_shared_volume_enabled").(bool)
	isAccelerationEnabled := d.Get("is_compression_acceleration_enabled").(bool)
	_, hasAlua := d.GetOk("is_alua_enabled")
	_, hasDRProcessMode := d.GetOk("data_reduction_process_mode")

	if isPG {
		// Parity-Group restrictions
		if capacitySaving != "disabled" || isShareEnabled || isAccelerationEnabled || hasDRProcessMode || hasAlua {
			return fmt.Errorf("data reduction and ALUA settings are only supported for DP-pool volumes (either pool_id or pool_name specified)")
		}
	} else if isDpPool {
		// DP-pool restrictions below
		if isCreate {
			if hasDRProcessMode || hasAlua {
				return fmt.Errorf("data_reduction_process_mode, is_alua_enabled cannot be specified during create; it is only valid for updates")
			}
		}
		if capacitySaving == "disabled" {
			if isShareEnabled || isAccelerationEnabled || hasDRProcessMode {
				return fmt.Errorf("data_reduction_process_mode, is_data_reduction_shared_volume_enabled=true, is_compression_acceleration_enabled=true can only be used when capacity_saving is not 'disabled'")

			}
		}
	}

	// Fix console output of "volume"
	d.SetNewComputed("volume")

	return nil
}
