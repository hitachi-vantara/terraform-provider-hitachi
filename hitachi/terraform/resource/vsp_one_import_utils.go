package terraform

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func trimImportID(raw string) (string, error) {
	id := strings.TrimSpace(raw)
	if id == "" {
		return "", fmt.Errorf("import id is empty")
	}
	return id, nil
}

func setIntIfParse(d *schema.ResourceData, key string, raw string) {
	if v, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
		_ = d.Set(key, v)
	}
}

// importAdminVolume expects: <serial>/<volume_id[,volume_id...]>
func importAdminVolume(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID, err := trimImportID(d.Id())
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<volume_id[,volume_id...]>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	volumeIDsPart := strings.TrimSpace(parts[1])
	if volumeIDsPart == "" {
		return nil, fmt.Errorf("invalid import id %q: missing volume id list", importID)
	}

	// Normalize volume ID list to match provider's stable ID format (sorted ints).
	rawParts := strings.Split(volumeIDsPart, ",")
	ids := make([]int, 0, len(rawParts))
	for _, p := range rawParts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid volume_id %q in import id: %w", p, err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("import id %q did not contain any volume ids", importID)
	}
	sort.Ints(ids)
	partsOut := make([]string, len(ids))
	for i, id := range ids {
		partsOut[i] = strconv.Itoa(id)
	}

	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if len(ids) == 1 {
		// Setting the selector is safe for single-volume resources and improves UX.
		if err := d.Set("volume_id", ids[0]); err != nil {
			return nil, err
		}
	}
	d.SetId(strings.Join(partsOut, ","))

	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importAdminPort expects: <serial>/<port_id>
func importAdminPort(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID, err := trimImportID(d.Id())
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<port_id>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	portID := strings.TrimSpace(parts[1])
	if portID == "" {
		return nil, fmt.Errorf("invalid port_id %q in import id", parts[1])
	}

	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if err := d.Set("port_id", portID); err != nil {
		return nil, err
	}

	// Preserve the resource's native ID format (port ID).
	d.SetId(portID)
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importAdminServer expects: <serial>/<server_id>
func importAdminServer(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID, err := trimImportID(d.Id())
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<server_id>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	serverID, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid server_id %q in import id: %w", parts[1], err)
	}
	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	// Preserve the resource's native ID format: "serial-serverid".
	d.SetId(fmt.Sprintf("%d-%d", serial, serverID))
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importAdminPool expects: <serial>/<pool_id>
func importAdminPool(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID, err := trimImportID(d.Id())
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<pool_id>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	poolID := strings.TrimSpace(parts[1])
	if poolID == "" {
		return nil, fmt.Errorf("invalid pool_id %q in import id", parts[1])
	}
	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}

	// Preserve native ID format: pool ID only.
	setIntIfParse(d, "pool_id", poolID)
	d.SetId(poolID)
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}
