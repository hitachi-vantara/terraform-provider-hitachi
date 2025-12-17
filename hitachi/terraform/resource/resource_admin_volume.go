package terraform

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
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
	_, numVolsExists := d.GetOk("number_of_volumes")
	isCreate := resourceID == ""

	log.WriteDebug("ValidateVolumeIDValues called with:")
	log.WriteDebug("  Resource ID: %q", resourceID)
	log.WriteDebug("  isCreate: %v", isCreate)
	log.WriteDebug("  volume_id exists: %v", volIDExists)
	log.WriteDebug("  number_of_volumes exists: %v", numVolsExists)
	log.WriteDebug("  HasChange(volume_id): %v", d.HasChange("volume_id"))

	// --- Skip validation for refresh (no change in volume_id) ---
	if resourceID != "" && !d.HasChange("volume_id") {
		log.WriteInfo("Skipping validation — no change to volume_id (likely refresh)")
		return nil
	}

	// --- CREATE ---
	if isCreate {
		if volIDExists {
			log.WriteError("Invalid: volume_id specified during create (only valid for update).")
			return fmt.Errorf("volume_id cannot be specified during create; it is only valid for update operations")
		}
		// no volume_id, number_of_volumes is optional -> valid
		log.WriteInfo("Create validation passed.")
		return nil
	}

	// --- UPDATE ---
	if !isCreate {
		// must have volume_id
		if !volIDExists {
			log.WriteError("Invalid: volume_id missing during update (required when ID exists).")
			return fmt.Errorf("volume_id must be specified for update operations")
		}
		// conflict: both given
		if numVolsExists {
			log.WriteError("Invalid: both volume_id and number_of_volumes specified during update.")
			return fmt.Errorf("volume_id and number_of_volumes cannot both be specified during update")
		}

		// volume_id must match one of the state IDs
		volID := rawVolID.(int)
		stateIDs := strings.Split(resourceID, ",")
		found := false
		for _, sid := range stateIDs {
			sid = strings.TrimSpace(sid)
			if sid == strconv.Itoa(volID) {
				found = true
				break
			}
		}
		if !found {
			log.WriteError("volume_id %d not found in current state IDs: %v", volID, stateIDs)
			return fmt.Errorf("volume_id %d not found in current state IDs (%v)", volID, stateIDs)
		}

		log.WriteDebug("volume_id %d matched in current state IDs", volID)
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
	isShareEnabled, _ := d.Get("is_data_reduction_share_enabled").(bool)

	if isShareEnabled && savingSetting == "DISABLE" {
		return fmt.Errorf("is_data_reduction_share_enabled can only be true if capacity_saving is not DISABLE")
	}

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
