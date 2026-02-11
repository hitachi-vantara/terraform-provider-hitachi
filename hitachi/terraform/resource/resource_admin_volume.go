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

// Mutex to prevent concurrent create operations
var syncVolumeOperation = &sync.Mutex{}

func ResourceAdminVolume() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage volumes in VSP One storage.",
		CreateContext: resourceAdminVolumeCreate,
		ReadContext:   resourceAdminVolumeRead,
		UpdateContext: resourceAdminVolumeUpdate,
		DeleteContext: resourceAdminVolumeDelete,
		Schema:        schemaimpl.ResourceAdminVolumeSchema(),
		CustomizeDiff: resourceAdminVolumeCustomizeDiff,
	}
}

func resourceAdminVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeOperation.Lock()
	defer syncVolumeOperation.Unlock()

	return impl.ResourceAdminVolumeCreate(d)
}

func resourceAdminVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminVolumeRead(d)
}

func resourceAdminVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeOperation.Lock()
	defer syncVolumeOperation.Unlock()

	return impl.ResourceAdminVolumeUpdate(d)
}

func resourceAdminVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeOperation.Lock()
	defer syncVolumeOperation.Unlock()

	return impl.ResourceAdminVolumeDelete(d)
}

func resourceAdminVolumeCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := ValidateNicknameParamValues(d); err != nil {
		return err
	}
	if err := ValidateDataReductionSettingsValues(d); err != nil {
		return err
	}
	if err := ValidateVolumeIDValues(d); err != nil {
		return err
	}

	d.SetNewComputed("volumes_info")
	d.SetNewComputed("volume_count")
	return nil
}

// minimalDiff-like type for testing
type volumeDiff interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Id() string
	HasChange(string) bool
}

func ValidateVolumeIDValues(d volumeDiff) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	resourceID := d.Id()
	rawVolID, volIDExists := d.GetOk("volume_id")
	rawVolHex, volHexExists := d.GetOk("volume_id_hex") // New: Get volume_id_hex
	_, numVolsExists := d.GetOk("number_of_volumes")
	isCreate := resourceID == ""

	log.WriteDebug("ValidateVolumeIDValues called with:")
	log.WriteDebug("  Resource ID: %q", resourceID)
	log.WriteDebug("  isCreate: %v", isCreate)
	log.WriteDebug("  volume_id exists: %v", volIDExists)
	log.WriteDebug("  volume_id_hex exists: %v", volHexExists) // New: Log volume_id_hex existence
	log.WriteDebug("  number_of_volumes exists: %v", numVolsExists)
	log.WriteDebug("  HasChange(volume_id): %v", d.HasChange("volume_id"))
	log.WriteDebug("  HasChange(volume_id_hex): %v", d.HasChange("volume_id_hex")) // New: Log volume_id_hex change

	// --- Skip validation for refresh (no change in volume_id or volume_id_hex) ---
	if resourceID != "" && !d.HasChange("volume_id") && !d.HasChange("volume_id_hex") {
		log.WriteInfo("Skipping validation — no change to volume_id/volume_id_hex (likely refresh)")
		return nil
	}

	// --- CREATE ---
	if isCreate {
		if volIDExists || volHexExists {
			log.WriteError("Invalid: volume_id/volume_id_hex specified during create (only valid for update).")
			return fmt.Errorf("volume_id or volume_id_hex cannot be specified during create; only valid for update operations")
		}
		// no volume_id/volume_id_hex, number_of_volumes is optional -> valid
		log.WriteInfo("Create validation passed.")
		return nil
	}

	// --- UPDATE ---
	if !isCreate {
		// New: Check for mutual exclusivity of volume_id and volume_id_hex
		if volIDExists && volHexExists {
			log.WriteError("Invalid: both volume_id and volume_id_hex cannot both be specified during update (must be mutually exclusive).")
			return fmt.Errorf("only one of volume_id or volume_id_hex can be specified during updates")
		}

		// must have either volume_id OR volume_id_hex
		if !volIDExists && !volHexExists {
			log.WriteError("Invalid: one of volume_id or volume_id_hex must be specified during update (one is required).")
			return fmt.Errorf("one of volume_id or volume_id_hex must be specified for update operations")
		}

		// conflict: both ID/Hex and number_of_volumes given
		if numVolsExists {
			log.WriteError("Invalid: number_of_volumes cannot be specified during update.")
			return fmt.Errorf("number_of_volumes cannot be specified during update")
		}

		// Determine the volume ID to validate against state IDs
		var volIDp *int
		var volHexStrp *string

		if volIDExists {
			volId := rawVolID.(int)
			volIDp = &volId
		} else {
			volHex := rawVolHex.(string)
			volHexStrp = &volHex
		}

		finalLdev, err := utils.ParseLdev(volIDp, volHexStrp)
		if err != nil {
			return err
		}
		volIDStr := strconv.Itoa(finalLdev)

		// must match one of the state IDs
		stateIDs := strings.Split(resourceID, ",")
		found := false
		for _, sid := range stateIDs {
			sid = strings.TrimSpace(sid)
			if sid == volIDStr { // Compare with the decimal ID string
				found = true
				break
			}
		}
		if !found {
			log.WriteError("Volume ID %s (from volume_id or volume_id_hex) not found in current state IDs: %v", volIDStr, stateIDs)
			return fmt.Errorf("volume ID %s (from volume_id or volume_id_hex) not found in current state IDs (%v)", volIDStr, stateIDs)
		}

		log.WriteDebug("Volume ID %s matched in current state IDs", volIDStr)
		log.WriteInfo("Update validation passed.")
		return nil
	}

	log.WriteInfo("Validation completed (no action).")
	return nil
}

func ValidateNicknameParamValues(d volumeDiff) error {
	paramsRaw, ok := d.GetOk("nickname_param")
	if !ok || paramsRaw == nil {
		return nil
	}

	paramList, ok := paramsRaw.([]interface{})
	if !ok || len(paramList) == 0 {
		return nil
	}

	m, ok := paramList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("nickname_param[0] must be a map")
	}

	baseName, _ := m["base_name"].(string)
	startNumber, _ := m["start_number"].(int)
	numberOfDigits, _ := m["number_of_digits"].(int)

	if startNumber == -1 && numberOfDigits > 0 {
		return fmt.Errorf("'number_of_digits' is specified but 'start_number' is not — please specify start_number when using number_of_digits")
	}

	totalLength := len(baseName)
	if startNumber != -1 {
		if numberOfDigits > 0 {
			totalLength += numberOfDigits
		} else {
			totalLength++ // default 1 digit suffix
		}
	}

	if totalLength > 32 {
		return fmt.Errorf("nickname length exceeds 32 characters: base_name (%d) + suffix digits = %d", len(baseName), totalLength)
	}

	return nil
}

func ValidateDataReductionSettingsValues(d volumeDiff) error {
	savingSetting, _ := d.Get("capacity_saving").(string)
	// isShareEnabled, _ := d.Get("is_data_reduction_share_enabled").(bool)
	// Commented out as this combination is required for DRS Volumes
	// if isShareEnabled && savingSetting == "DISABLE" {
	// 	return fmt.Errorf("is_data_reduction_share_enabled can only be true if capacity_saving is not DISABLE")
	// }

	// Handle compression_acceleration
	if _, ok := d.GetOk("compression_acceleration"); ok {
		isCreate := d.Id() == ""

		if isCreate {
			// During Create, compression_acceleration is output-only
			return fmt.Errorf("compression_acceleration cannot be specified during create; it is only valid for updates")
		}

		// Optional: sanity check — can only apply if savingSetting != DISABLE
		if savingSetting == "DISABLE" {
			return fmt.Errorf("compression_acceleration can only be used when capacity_saving is not DISABLE (current: %s)", savingSetting)
		}
	}

	return nil
}
