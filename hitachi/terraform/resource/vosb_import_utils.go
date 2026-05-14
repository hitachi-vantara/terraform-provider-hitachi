package terraform

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importVosbAddressAndKey(d *schema.ResourceData, meta interface{}) (addr string, key string, err error) {
	importID := ""
	if d != nil {
		importID = d.Id()
	}
	if importID == "" {
		return "", "", fmt.Errorf("invalid import id %q", importID)
	}
	parts := strings.SplitN(strings.TrimSpace(importID), "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid import id %q; expected '<vosb_address>/<key>'", importID)
	}
	addr = strings.TrimSpace(parts[0])
	key = strings.TrimSpace(parts[1])
	if addr == "" || key == "" {
		return "", "", fmt.Errorf("invalid import id %q; expected '<vosb_address>/<key>'", importID)
	}
	return addr, key, nil
}

// importVosbWithAddress expects: <vosb_address>/<resource_id>
func importVosbWithAddress(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addr, key, err := importVosbAddressAndKey(d, meta)
	if err != nil {
		return nil, err
	}
	if err := d.Set("vosb_address", addr); err != nil {
		return nil, err
	}
	d.SetId(key)
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

func importVosbWithLookupField(fieldName string) schema.StateContextFunc {
	// Many VOSB resources read by name (or other string key) rather than d.Id().
	// During import we must set that lookup field so Read can find the object.
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		addr, key, err := importVosbAddressAndKey(d, meta)
		if err != nil {
			return nil, err
		}
		if err := d.Set("vosb_address", addr); err != nil {
			return nil, err
		}
		if err := d.Set(fieldName, key); err != nil {
			return nil, err
		}
		// Preserve a stable pre-Read ID; Read will replace it with backend UUID (if applicable).
		d.SetId(key)
		return []*schema.ResourceData{d}, nil
	}
}

// importVosbVolume supports:
//  1. <vosb_address>/<volume_name>
func importVosbVolume(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addr, key, err := importVosbAddressAndKey(d, meta)
	if err != nil {
		return nil, err
	}
	addr = strings.TrimSpace(addr)
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("invalid import id %q", d.Id())
	}
	if err := d.Set("vosb_address", addr); err != nil {
		return nil, err
	}
	if strings.HasPrefix(key, "id:") {
		return nil, fmt.Errorf("invalid import id %q; importing by UUID is not supported, expected '<vosb_address>/<volume_name>'", d.Id())
	}

	// Default: treat key as a volume name.
	if err := d.Set("name", key); err != nil {
		return nil, err
	}
	// Preserve a stable pre-Read ID; Read will replace it with backend UUID.
	d.SetId(key)
	return []*schema.ResourceData{d}, nil
}

// importVosbComputeNode expects: <vosb_address>/<compute_node_name>
var importVosbComputeNode = importVosbWithLookupField("compute_node_name")

// importVosbStorageNode expects: <vosb_address>/<node_name>
var importVosbStorageNode = importVosbWithLookupField("node_name")

// importVosbChapUser expects: <vosb_address>/<target_chap_user_name>
var importVosbChapUser = importVosbWithLookupField("target_chap_user_name")
