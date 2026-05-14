package terraform

import (
	"context"
	"fmt"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAdminPool() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage VSP One pools in Hitachi storage system",
		CreateContext: resourceAdminPoolCreate,
		ReadContext:   resourceAdminPoolRead,
		UpdateContext: resourceAdminPoolUpdate,
		DeleteContext: resourceAdminPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importAdminPool,
		},
		Schema:        schemaimpl.ResourceAdminPoolSchema(),
		CustomizeDiff: resourceAdminPoolCustomiseDiff,
	}
}

func resourceAdminPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminPoolCreate(d)
}

func resourceAdminPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminPoolRead(d)
}

func resourceAdminPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminPoolUpdate(d)
}

func resourceAdminPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceAdminPoolDelete(d)
}

// resourceAdminPoolCustomiseDiff performs custom validation for VSP One pool configuration
func resourceAdminPoolCustomiseDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Create-only requirements (kept out of the schema to allow import with minimal config).
	if diff.Id() == "" {
		if serialRaw, ok := diff.GetOk("serial"); !ok || serialRaw.(int) < 1 {
			return fmt.Errorf("serial must be specified and must be >= 1")
		}
		if nameRaw, ok := diff.GetOk("name"); !ok || nameRaw.(string) == "" {
			return fmt.Errorf("name must be specified")
		}
		drivesRaw, ok := diff.GetOk("drive_configuration")
		if !ok {
			return fmt.Errorf("drive_configuration must be specified")
		}
		drives := drivesRaw.([]interface{})
		if len(drives) == 0 {
			return fmt.Errorf("drive_configuration must include at least one block")
		}
		if err := validateDriveBlocks(drives); err != nil {
			return err
		}
	}

	// Validate thresholds
	if err := validateThresholds(ctx, diff, v); err != nil {
		return err
	}

	// Validate encryption immutability
	if err := validateEncryptionImmutability(ctx, diff, v); err != nil {
		return err
	}

	// Validate drive configuration changes
	if err := validateDriveConfigurationChanges(ctx, diff, v); err != nil {
		return err
	}

	if err := diff.SetNewComputed("data"); err != nil {
		return err
	}
	return nil
}

// validateThresholds validates that threshold warning and depletion are properly configured
func validateThresholds(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	warningRaw, warningExists := diff.GetOkExists("threshold_warning")
	depletionRaw, depletionExists := diff.GetOkExists("threshold_depletion")

	// If one threshold is specified, both must be specified
	if (warningExists && !depletionExists) || (!warningExists && depletionExists) {
		return fmt.Errorf("both threshold_warning and threshold_depletion must be specified together")
	}

	// If both are specified, validate the relationship
	if warningExists && depletionExists {
		warning := warningRaw.(int)
		depletion := depletionRaw.(int)

		if depletion < warning {
			return fmt.Errorf("threshold_depletion (%d) must be greater than or equal to threshold_warning (%d)", depletion, warning)
		}
	}

	return nil
}

// validateEncryptionImmutability validates that encryption setting is not changed after creation
func validateEncryptionImmutability(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Immutable after creation.
	if diff.Id() != "" && diff.HasChange("encryption") {
		return fmt.Errorf("encryption setting is immutable; recreate the resource to change encryption")
	}
	return nil
}

// validateDriveConfigurationChanges validates drive configuration modifications
func validateDriveConfigurationChanges(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	if !diff.HasChange("drive_configuration") {
		return nil
	}

	oldDrivesRaw, newDrivesRaw := diff.GetChange("drive_configuration")
	oldDrives, _ := oldDrivesRaw.([]interface{})
	newDrives, _ := newDrivesRaw.([]interface{})

	// Validate blocks whenever drive_configuration is set/changed.
	if err := validateDriveBlocks(newDrives); err != nil {
		return err
	}

	// NOTE: For imported resources, Terraform won't have any prior drive_configuration in state.
	// We treat the first configured drive_configuration as an explicit expansion request
	// (i.e., "add these drives"). After the first apply, append-only semantics apply.

	// If drives are being removed, that's not allowed
	if len(newDrives) < len(oldDrives) {
		return fmt.Errorf("removing drive_configuration blocks is not allowed; only adding new blocks is permitted")
	}

	// If no new drives were added, ensure existing blocks are unchanged
	if len(newDrives) <= len(oldDrives) {
		// Verify immutability of existing drive blocks
		for i := 0; i < len(newDrives) && i < len(oldDrives); i++ {
			nd := newDrives[i].(map[string]interface{})
			od := oldDrives[i].(map[string]interface{})

			if nd["drive_type_code"] != od["drive_type_code"] ||
				nd["data_drive_count"] != od["data_drive_count"] ||
				nd["raid_level"] != od["raid_level"] ||
				nd["parity_group_type"] != od["parity_group_type"] {
				return fmt.Errorf("modifying existing drive_configuration blocks is not allowed; only adding new blocks is permitted")
			}
		}
	}

	return nil
}

func validateDriveBlocks(blocks []interface{}) error {
	for i, b := range blocks {
		m, ok := b.(map[string]interface{})
		if !ok {
			return fmt.Errorf("drive_configuration[%d] must be an object", i)
		}

		raidLevel, _ := m["raid_level"].(string)
		dataDriveCount, _ := m["data_drive_count"].(int)

		switch raidLevel {
		case "RAID5":
			if dataDriveCount < 5 {
				return fmt.Errorf("drive_configuration[%d].data_drive_count must be >= 5 for RAID5", i)
			}
		case "RAID6":
			if dataDriveCount < 9 {
				return fmt.Errorf("drive_configuration[%d].data_drive_count must be >= 9 for RAID6", i)
			}
		}
	}
	return nil
}
