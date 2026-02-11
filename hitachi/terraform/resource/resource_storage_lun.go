package terraform

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncLunOperation = &sync.Mutex{}

func ResourceStorageLun() *schema.Resource {
	return &schema.Resource{
		Description: `VSP Storage Volume: The following request creates a volume by using the specified parity groups or pools. Specify a parity group or pool id for creating a basic volume.`,
		Importer: &schema.ResourceImporter{
			StateContext: importStorageLunState,
		},
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
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

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

	d.SetId(strconv.Itoa(logicalUnit.LdevID))
	log.WriteInfo("lun read successfully")

	return nil

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
	isUpdate := !isCreate

	// -----------------------------------------------------
	// Validate LDEV ID / HEX
	// -----------------------------------------------------

	// Retrieve both values once
	ldevIDValue, ldevIDProvided := d.GetOk("ldev_id")
	ldevHexValue, ldevHexProvided := d.GetOk("ldev_id_hex")

	// Only one allowed
	if ldevIDProvided && ldevHexProvided {
		return fmt.Errorf("only one of ldev_id or ldev_id_hex may be specified, not both")
	}

	if impl.IsReadExistingMode(d) {
		log.WriteInfo("Validation| Detected 'Read Existing Volume' mode. Bypassing creation checks.")
		d.SetNewComputed("volume")
		return nil
	}

	if isUpdate {
		var newLdevID *int

		// --- Handle ldev_id ---
		if ldevIDProvided {
			id, ok := ldevIDValue.(int)
			if !ok {
				return fmt.Errorf("invalid type for ldev_id")
			}
			newLdevID = &id
		}

		// --- Handle ldev_id_hex ---
		if newLdevID == nil && ldevHexProvided {
			hexStr, ok := ldevHexValue.(string)
			if !ok {
				return fmt.Errorf("invalid type for ldev_id_hex")
			}

			parsed, err := utils.HexStringToInt(hexStr)
			if err != nil {
				return fmt.Errorf("failed to parse ldev_id_hex: %v", err)
			}
			newLdevID = &parsed
		}

		// Neither provided → do nothing
		if newLdevID != nil {
			// Parse ID from TF state string
			idFromState, err := strconv.Atoi(currentID)
			if err != nil {
				return fmt.Errorf("failed to parse current resource ID (%s): %v", currentID, err)
			}

			// Compare values
			if *newLdevID != idFromState {
				return fmt.Errorf(
					"LDEV value (%d) does not match the provisioned LDEV ID (%d). "+
						"Hitachi LDEVs cannot be reassigned after provisioning.",
					*newLdevID, idFromState,
				)
			}
		}
	}

	// -----------------------------------------------------
	// Mainframe vs Block input validation
	// -----------------------------------------------------
	cylinderVal, hasCylinder := d.GetOk("cylinder")
	var cylinder int
	if hasCylinder {
		cylinder = cylinderVal.(int)
	}

	emulationType := ""
	_, emulationProvided := d.GetOk("emulation_type")
	if emulationProvided {
		emulationType = d.Get("emulation_type").(string)
	}
	is3390A := emulationProvided && strings.EqualFold(emulationType, "3390-A")
	is3390V := emulationProvided && strings.EqualFold(emulationType, "3390-V")

	// size_gb is the block-capacity input in this resource.
	_, hasSizeGB := d.GetOk("size_gb")

	_, hasSSIDOnly := d.GetOk("ssid")
	_, hasMpBlade := d.GetOk("mp_blade_id")
	_, hasClpr := d.GetOk("clpr_id")
	// bools have defaults, so treat them as mainframe only when explicitly set true.
	isTse := d.Get("is_tse_volume").(bool)
	isEse := d.Get("is_ese_volume").(bool)

	// Mainframe volume is indicated ONLY by cylinder.
	isMainframe := hasCylinder
	if !isCreate && isMainframe {
		// On updates, only ESE enablement is configurable among mainframe-specific fields.
		// Other mainframe fields are immutable after creation.
		if d.HasChange("cylinder") || d.HasChange("emulation_type") || d.HasChange("ssid") || d.HasChange("mp_blade_id") || d.HasChange("clpr_id") || d.HasChange("is_tse_volume") {
			return fmt.Errorf("for mainframe volumes, only is_ese_volume can be updated; cylinder/emulation_type/ssid/mp_blade_id/clpr_id/is_tse_volume are immutable")
		}
	}

	// Mainframe-only fields are not valid unless cylinder is set.
	if !isMainframe {
		if hasSSIDOnly || hasMpBlade || hasClpr || isTse || isEse || emulationProvided {
			return fmt.Errorf("mainframe-only fields (cylinder/emulation_type/ssid/mp_blade_id/clpr_id/is_tse_volume/is_ese_volume) cannot be specified unless cylinder is set")
		}
	}

	// Capacity requirement:
	// - On create: must specify exactly one of cylinder (mainframe) or size_gb (block)
	// - On update: only validate what is provided (do not force respecifying capacity)
	if isCreate {
		if isMainframe {
			if cylinder <= 0 {
				return fmt.Errorf("cylinder must be >= 1")
			}
			if hasSizeGB {
				return fmt.Errorf("size_gb cannot be specified for mainframe volumes; use cylinder")
			}
		} else {
			if !hasSizeGB {
				return fmt.Errorf("either size_gb (block volume) or cylinder (mainframe volume) must be specified")
			}

			// volume_format_type is update-only
			if isCreate {
				if v, ok := d.GetOk("volume_format_type"); ok {
					vs := v.(string)
					// Default NONE is allowed on create; disallow QUICK/NORMAL on create
					if !strings.EqualFold(vs, "NONE") && vs != "" {
						return fmt.Errorf("volume_format_type cannot be specified as QUICK or NORMAL during create; it may only be used on update")
					}
				}
			}
			// Block volumes do not support emulation_type.
			if emulationProvided {
				return fmt.Errorf("emulation_type is not supported for block volumes")
			}
		}
	} else {
		// Update: validate mainframe vs block constraints only when fields are present.
		if isMainframe {
			if cylinder <= 0 {
				return fmt.Errorf("cylinder must be >= 1")
			}
			if hasSizeGB {
				return fmt.Errorf("size_gb cannot be specified for mainframe volumes; use cylinder")
			}
		} else {
			// Block volumes do not support emulation_type.
			if emulationProvided {
				return fmt.Errorf("emulation_type is not supported for block volumes")
			}
		}
	}

	// Validate size_gb only when used.
	if hasSizeGB {
		sizeRaw := d.Get("size_gb")
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
	}

	// -----------------------------------------------------
	// Pool or Parity Group validation
	// -----------------------------------------------------

	// pool_id: default -999
	poolID := -999
	if v, ok := d.GetOkExists("pool_id"); ok {
		poolID = v.(int)
	}
	hasPoolID := poolID >= -1

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

	// -----------------------------------------------------
	// Mainframe-specific validations (cylinder-driven)
	// -----------------------------------------------------
	if isMainframe {
		// LDEV ID/HEX are not supported for mainframe volume creation.
		if ldevIDProvided || ldevHexProvided {
			return fmt.Errorf("ldev_id and ldev_id_hex are not supported for mainframe volumes")
		}
		// external_paritygroup_id not supported in mainframe.
		if hasExternalParityGroup {
			return fmt.Errorf("external_paritygroup_id is not supported for mainframe volumes")
		}

		// emulation_type is supported for mainframe volumes, including OPEN-V.

		// 3390-A: dynamic pool mainframe volume
		if is3390A {
			if !isDpPool {
				return fmt.Errorf("for emulation_type=3390-A, one of pool_id or pool_name must be specified")
			}
			if hasParityGroup {
				return fmt.Errorf("for emulation_type=3390-A, paritygroup_id must not be specified")
			}
		}

		// 3390-V: parity group mainframe volume
		if is3390V {
			if !hasParityGroup {
				return fmt.Errorf("for emulation_type=3390-V, paritygroup_id must be specified")
			}
			if isDpPool {
				return fmt.Errorf("for emulation_type=3390-V, pool_id/pool_name must not be specified")
			}
		}

		// clpr_id is supported only for pool-based volumes.
		if _, ok := d.GetOk("clpr_id"); ok && !isDpPool {
			return fmt.Errorf("clpr_id can only be specified when creating pool-based volumes (pool_id or pool_name)")
		}

		// TSE and ESE are mutually exclusive.
		isTse := d.Get("is_tse_volume").(bool)
		isEse := d.Get("is_ese_volume").(bool)
		if isTse && isEse {
			return fmt.Errorf("is_tse_volume and is_ese_volume cannot both be true")
		}

		// Validate mainframe label/name restrictions when name is provided.
		if v, ok := d.GetOk("name"); ok {
			if err := validateMainframeVolumeLabel(v.(string)); err != nil {
				return err
			}
		}
	}

	// Flags (validated differently for mainframe vs block)
	capacitySaving := d.Get("capacity_saving").(string)
	isShareEnabled := d.Get("is_data_reduction_shared_volume_enabled").(bool)
	isAccelerationEnabled := d.Get("is_compression_acceleration_enabled").(bool)
	_, hasAlua := d.GetOk("is_alua_enabled")
	_, hasDRProcessMode := d.GetOk("data_reduction_process_mode")

	if isMainframe {
		// Keep mainframe validation isolated: block-only knobs are not supported here.
		if capacitySaving != "disabled" || isShareEnabled || isAccelerationEnabled || hasDRProcessMode || hasAlua {
			return fmt.Errorf("capacity_saving, is_data_reduction_shared_volume_enabled, is_compression_acceleration_enabled, data_reduction_process_mode, and is_alua_enabled are not supported for mainframe volumes")
		}
	} else {
		// Block-specific validations
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
			// if capacitySaving == "disabled" {
			// 	if isShareEnabled || isAccelerationEnabled || hasDRProcessMode {
			// 		return fmt.Errorf("data_reduction_process_mode, is_data_reduction_shared_volume_enabled=true, is_compression_acceleration_enabled=true can only be used when capacity_saving is not 'disabled'")
			// 	}
			// }
		}
	}

	// Fix console output of "volume"
	d.SetNewComputed("volume")

	return nil
}

// importStorageLunState supports importing either "<serial>/<ldev>" or just "<ldev>".
// If only <ldev> is provided, the resource config must include `serial`.
func importStorageLunState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// id provided to import
	id := d.Id()
	// If id already contains '/', parse serial and ldev
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		serialStr := parts[0]
		ldevStr := parts[1]
		serial, err := strconv.Atoi(serialStr)
		if err != nil {
			return nil, fmt.Errorf("invalid serial in import id: %s", serialStr)
		}
		// set serial in the config state so provider has it
		if err := d.Set("serial", serial); err != nil {
			return nil, err
		}
		d.SetId(ldevStr)
		return []*schema.ResourceData{d}, nil
	}

	// No serial in id; ensure serial is present in resource config
	if _, ok := d.GetOk("serial"); !ok {
		return nil, fmt.Errorf("import requires either '<serial>/<ldev>' or resource config must include 'serial'")
	}

	// id is the ldev id; leave d.Id() as-is
	return []*schema.ResourceData{d}, nil
}

func validateMainframeVolumeLabel(label string) error {
	// 1 to 32 characters.
	if len(label) < 1 || len(label) > 32 {
		return fmt.Errorf("name must be between 1 and 32 characters for mainframe volumes")
	}
	// Cannot start or end with space.
	if strings.HasPrefix(label, " ") || strings.HasSuffix(label, " ") {
		return fmt.Errorf("name cannot start or end with a space for mainframe volumes")
	}

	allowed := func(r rune) bool {
		if r >= 'a' && r <= 'z' {
			return true
		}
		if r >= 'A' && r <= 'Z' {
			return true
		}
		if r >= '0' && r <= '9' {
			return true
		}
		if r == ' ' {
			return true
		}
		switch r {
		case '!', '#', '$', '%', '&', '\'', '(', ')', '+', ',', '-', '.', ':', '=', '@', '[', ']', '^', '_', '`', '{', '}', '~', '/', '\\':
			return true
		default:
			return false
		}
	}

	for _, r := range label {
		if !allowed(r) {
			return fmt.Errorf("name contains an invalid character for mainframe volumes")
		}
	}

	return nil
}
