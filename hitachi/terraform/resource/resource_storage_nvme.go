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

// Mutex to synchronize NVM Subsystem operations (Creation/Deletion/Updates)
var syncNvmSubsystemOperation = &sync.Mutex{}

func ResourceVspNvmSubsystem() *schema.Resource {
	return &schema.Resource{
		Description:   `Vsp NVM Subsystem Resource: Manages NVMe subsystems, including port assignments, Host NQN registrations, and Namespace mappings.`,
		CreateContext: resourceVspNvmSubsystemCreate,
		ReadContext:   resourceVspNvmSubsystemRead,
		UpdateContext: resourceVspNvmSubsystemUpdate,
		DeleteContext: resourceVspNvmSubsystemDelete,
		Schema:        schemaimpl.ResourceVspNvmSubsystemSchema(),
		CustomizeDiff: resourceVspNvmSubsystemCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: importVspNvmSubsystemState,
		},
	}
}

func resourceVspNvmSubsystemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncNvmSubsystemOperation.Lock()
	defer syncNvmSubsystemOperation.Unlock()

	// Deferring all logic to implementation layer (Apply)
	return impl.ResourceVspNvmSubsystemApply(d)
}

func resourceVspNvmSubsystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return impl.ResourceVspNvmSubsystemRead(d)
}

func resourceVspNvmSubsystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncNvmSubsystemOperation.Lock()
	defer syncNvmSubsystemOperation.Unlock()

	return impl.ResourceVspNvmSubsystemApply(d)
}

func resourceVspNvmSubsystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncNvmSubsystemOperation.Lock()
	defer syncNvmSubsystemOperation.Unlock()

	return impl.ResourceVspNvmSubsystemDelete(d)
}

func resourceVspNvmSubsystemCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	log := commonlog.GetLogger()

	// -------------------------------------------------
	// 1. Cross-field validation
	// -------------------------------------------------

	// host_mode_options requires host_mode
	if d.Get("host_mode_options") != nil {
		hostModeOptions := d.Get("host_mode_options").([]interface{})
		hostMode := d.Get("host_mode").(string)

		if len(hostModeOptions) > 0 && hostMode == "" {
			return fmt.Errorf("host_mode must be set when host_mode_options are specified")
		}
	}

	// -------------------------------------------------
	// 2. Namespace LDEV Validation (ldev_id vs ldev_id_hex)
	// -------------------------------------------------
	if nsVal, ok := d.GetOk("namespaces"); ok {
		for _, item := range nsVal.(*schema.Set).List() {
			m := item.(map[string]interface{})

			// Retrieve actual values based on your schema defaults
			idVal := m["ldev_id"].(int)         // Default is -1
			hexVal := m["ldev_id_hex"].(string) // Default is ""

			hasID := idVal != -1
			hasHex := hexVal != ""

			// Both provided → invalid
			if hasID && hasHex {
				return fmt.Errorf("only one of ldev_id or ldev_id_hex may be specified in each namespace block")
			}

			// Neither provided → invalid
			if !hasID && !hasHex {
				return fmt.Errorf("one of ldev_id or ldev_id_hex must be specified in each namespace block")
			}
		}
	}

	// -------------------------------------------------
	// 3. Prevent meaningless drift in computed output
	// -------------------------------------------------

	if err := d.SetNewComputed("nvm_subsystems"); err != nil {
		log.WriteError("Failed to set nvm_subsystems computed: %v", err)
		return err
	}

	return nil
}

// importVspNvmSubsystemState supports importing "<serial>/<nvm_subsystem_id>".
func importVspNvmSubsystemState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// ID provided to import
	id := d.Id()

	if !strings.Contains(id, "/") {
		return nil, fmt.Errorf("import requires '<serial>/<nvm_subsystem_id>'")
	}

	parts := strings.SplitN(id, "/", 2)

	serialStr := parts[0]
	subsystemStr := parts[1]

	serial, err := strconv.Atoi(serialStr)
	if err != nil {
		return nil, fmt.Errorf("invalid serial in import id: %s", serialStr)
	}

	subsystemID, err := strconv.Atoi(subsystemStr)
	if err != nil {
		return nil, fmt.Errorf("invalid nvm_subsystem_id in import id: %s", subsystemStr)
	}

	// Set serial into state so provider can use it
	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}

	if err := d.Set("nvm_subsystem_id", subsystemID); err != nil {
		return nil, err
	}

	// Keep ID as just subsystem ID (consistent with your LUN pattern)
	d.SetId(subsystemStr)

	return []*schema.ResourceData{d}, nil
}
