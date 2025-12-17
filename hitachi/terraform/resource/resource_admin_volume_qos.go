package terraform

import (
	"context"
	"fmt"
	"sync"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Mutex to prevent concurrent set operation
var syncVolumeQosOperation = &sync.Mutex{}

func ResourceAdminVolumeQos() *schema.Resource {

	return &schema.Resource{
		Description:   "Manage volume QoS settings in VSP One storage.",
		CreateContext: resourceAdminVolumeQosCreate,
		ReadContext:   resourceAdminVolumeQosRead,
		UpdateContext: resourceAdminVolumeQosUpdate,
		DeleteContext: resourceAdminVolumeQosNoop,
		Schema:        schemaimpl.StorageVolumeSetQosSchema,
		CustomizeDiff: resourceAdminVolumeQosCustomizeDiff,
	}
}

func resourceAdminVolumeQosCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeQosOperation.Lock()
	defer syncVolumeQosOperation.Unlock()
	return impl.ResourceAdminVolumeQosCreate(d)
}

func resourceAdminVolumeQosRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeQosOperation.Lock()
	defer syncVolumeQosOperation.Unlock()
	return impl.ResourceAdminVolumeQosRead(d)
}

func resourceAdminVolumeQosNoop(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// TODO need to implement update
func resourceAdminVolumeQosUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	syncVolumeQosOperation.Lock()
	defer syncVolumeQosOperation.Unlock()

	return impl.ResourceAdminVolumeQosUpdate(d)
}

func resourceAdminVolumeQosCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Validate threshold block
	if v, ok := d.Get("threshold").([]interface{}); ok && len(v) > 0 && v[0] != nil {
		m, _ := v[0].(map[string]interface{})
		if m != nil {
			if b, ok := m["is_upper_iops_enabled"].(bool); ok {
				uv := m["upper_iops"].(int)
				if b && uv == 0 {
					return fmt.Errorf("upper_iops must be non-zero when is_upper_iops_enabled is true")
				}
				if !b && uv != 0 {
					return fmt.Errorf("upper_iops must be 0 when is_upper_iops_enabled is false")
				}
				// lower_iops check for comparison
				if lb, ok := m["is_lower_iops_enabled"].(bool); ok {
					lv := m["lower_iops"].(int)
					if b && lb && uv != 0 && lv != 0 && uv <= lv {
						return fmt.Errorf("upper_iops must be greater than lower_iops when both are enabled and non-zero")
					}
				}
			}
			if b, ok := m["is_upper_transfer_rate_enabled"].(bool); ok {
				uv := m["upper_transfer_rate"].(int)
				if b && uv == 0 {
					return fmt.Errorf("upper_transfer_rate must be non-zero when is_upper_transfer_rate_enabled is true")
				}
				if !b && uv != 0 {
					return fmt.Errorf("upper_transfer_rate must be 0 when is_upper_transfer_rate_enabled is false")
				}
				// lower_transfer_rate check for comparison
				if lb, ok := m["is_lower_transfer_rate_enabled"].(bool); ok {
					lv := m["lower_transfer_rate"].(int)
					if b && lb && uv != 0 && lv != 0 && uv <= lv {
						return fmt.Errorf("upper_transfer_rate must be greater than lower_transfer_rate when both are enabled and non-zero")
					}
				}
			}
			if b, ok := m["is_lower_iops_enabled"].(bool); ok {
				v := m["lower_iops"].(int)
				if b && v == 0 {
					return fmt.Errorf("lower_iops must be non-zero when is_lower_iops_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("lower_iops must be 0 when is_lower_iops_enabled is false")
				}
			}
			if b, ok := m["is_lower_transfer_rate_enabled"].(bool); ok {
				v := m["lower_transfer_rate"].(int)
				if b && v == 0 {
					return fmt.Errorf("lower_transfer_rate must be non-zero when is_lower_transfer_rate_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("lower_transfer_rate must be 0 when is_lower_transfer_rate_enabled is false")
				}
			}
			if b, ok := m["is_response_priority_enabled"].(bool); ok {
				v := m["response_priority"].(int)
				if b && v == 0 {
					return fmt.Errorf("response_priority must be non-zero when is_response_priority_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("response_priority must be 0 when is_response_priority_enabled is false")
				}
			}
		}
	}
	// Validate alert_setting block
	if v, ok := d.Get("alert_setting").([]interface{}); ok && len(v) > 0 && v[0] != nil {
		m, _ := v[0].(map[string]interface{})
		if m != nil {
			if b, ok := m["is_upper_alert_enabled"].(bool); ok {
				v := m["upper_alert_allowable_time"].(int)
				if b && v == 0 {
					return fmt.Errorf("upper_alert_allowable_time must be non-zero when is_upper_alert_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("upper_alert_allowable_time must be 0 when is_upper_alert_enabled is false")
				}
			}
			if b, ok := m["is_lower_alert_enabled"].(bool); ok {
				v := m["lower_alert_allowable_time"].(int)
				if b && v == 0 {
					return fmt.Errorf("lower_alert_allowable_time must be non-zero when is_lower_alert_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("lower_alert_allowable_time must be 0 when is_lower_alert_enabled is false")
				}
			}
			if b, ok := m["is_response_alert_enabled"].(bool); ok {
				v := m["response_alert_allowable_time"].(int)
				if b && v == 0 {
					return fmt.Errorf("response_alert_allowable_time must be non-zero when is_response_alert_enabled is true")
				}
				if !b && v != 0 {
					return fmt.Errorf("response_alert_allowable_time must be 0 when is_response_alert_enabled is false")
				}
			}
		}
	}
	return nil
}
