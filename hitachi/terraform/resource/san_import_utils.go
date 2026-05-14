package terraform

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importExpectNParts(importID string, n int) ([]string, error) {
	parts := strings.Split(importID, "/")
	if len(parts) != n {
		return nil, fmt.Errorf("invalid import id %q: expected %d '/'-separated parts", importID, n)
	}
	for i := range parts {
		if parts[i] == "" {
			return nil, fmt.Errorf("invalid import id %q: empty segment", importID)
		}
	}
	return parts, nil
}

func importParseSerial(serialStr string) (int, error) {
	serial, err := strconv.Atoi(serialStr)
	if err != nil {
		return 0, fmt.Errorf("invalid serial %q in import id: %w", serialStr, err)
	}
	return serial, nil
}

// importVspSnapshot supports:
//   - <serial>/<pvol_ldev_id>
//   - <serial>/<pvol_ldev_id>,<mirror_unit_id>
//
// It sets the native snapshot ID into state (either "<pvol>" or "<pvol>,<mu>").
func importVspSnapshot(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := strings.TrimSpace(d.Id())
	if importID == "" {
		return nil, fmt.Errorf("import id is empty")
	}

	parts, err := importExpectNParts(importID, 2)
	if err != nil {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<pvol_ldev_id>' or '<serial>/<pvol_ldev_id>,<mirror_unit_id>'", importID)
	}

	serial, err := importParseSerial(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, err
	}

	snapshotPart := strings.TrimSpace(parts[1])
	commaParts := strings.Split(snapshotPart, ",")
	if len(commaParts) != 1 && len(commaParts) != 2 {
		return nil, fmt.Errorf("invalid snapshot portion %q in import id; expected '<pvol_ldev_id>' or '<pvol_ldev_id>,<mirror_unit_id>'", snapshotPart)
	}

	pvolStr := strings.TrimSpace(commaParts[0])
	pvol, err := strconv.Atoi(pvolStr)
	if err != nil {
		return nil, fmt.Errorf("invalid pvol_ldev_id %q in import id: %w", pvolStr, err)
	}

	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if err := d.Set("pvol_ldev_id", pvol); err != nil {
		return nil, err
	}

	if len(commaParts) == 2 {
		muStr := strings.TrimSpace(commaParts[1])
		mu, err := strconv.Atoi(muStr)
		if err != nil {
			return nil, fmt.Errorf("invalid mirror_unit_id %q in import id: %w", muStr, err)
		}
		if err := d.Set("mirror_unit_id", mu); err != nil {
			return nil, err
		}
		// Use the native resource ID format.
		d.SetId(fmt.Sprintf("%d,%d", pvol, mu))
	} else {
		// Use the native resource ID format (without MU).
		d.SetId(fmt.Sprintf("%d", pvol))
	}

	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importVspHostGroup supports:
//   - <serial>/<port_id>,<hostgroup_name> (hostgroup_number resolved via API)
//   - <serial>/<port_id>,<hostgroup_number>,<hostgroup_name>
func importVspHostGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := strings.TrimSpace(d.Id())
	if importID == "" {
		return nil, fmt.Errorf("import id is empty")
	}
	// Command-line-only import:
	// <serial>/<port_id>,<hostgroup_name> OR <serial>/<port_id>,<hostgroup_number>,<hostgroup_name>
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<port_id>,<hostgroup_name>' or '<serial>/<port_id>,<hostgroup_number>,<hostgroup_name>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	hostGroupPart := strings.TrimSpace(parts[1])
	commaParts := strings.Split(hostGroupPart, ",")
	if len(commaParts) != 2 && len(commaParts) != 3 {
		return nil, fmt.Errorf("invalid hostgroup portion %q in import id; expected '<port_id>,<hostgroup_name>' or '<port_id>,<hostgroup_number>,<hostgroup_name>'", hostGroupPart)
	}

	portID := strings.TrimSpace(commaParts[0])
	if portID == "" {
		return nil, fmt.Errorf("invalid port_id %q in import id", commaParts[0])
	}

	hgNum := -1
	hgName := ""
	if len(commaParts) == 3 {
		parsedPortID, parsedHgNum, parsedHgName, err := terrcommon.ParseHostGroupFromID(hostGroupPart)
		if err != nil {
			return nil, fmt.Errorf("invalid hostgroup portion %q in import id: %w", parts[1], err)
		}
		portID = parsedPortID
		hgNum = parsedHgNum
		hgName = parsedHgName
	} else {
		// <port_id>,<hostgroup_name> form
		hgName = strings.TrimSpace(commaParts[1])
		if hgName == "" {
			return nil, fmt.Errorf("invalid hostgroup_name %q in import id", commaParts[1])
		}
	}

	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if err := d.Set("port_id", portID); err != nil {
		return nil, err
	}
	if err := d.Set("hostgroup_name", hgName); err != nil {
		return nil, err
	}

	if hgNum == -1 {
		// Resolve hostgroup_number by listing hostgroups on the array.
		all, err := impl.GetAllHostGroups(d)
		if err != nil {
			return nil, err
		}
		if all == nil || len(all.HostGroups) == 0 {
			return nil, fmt.Errorf("no hostgroups found on storage (serial=%d)", serial)
		}
		matches := make([]int, 0, 2)
		for i := range all.HostGroups {
			hg := &all.HostGroups[i]
			if strings.TrimSpace(hg.PortID) != portID {
				continue
			}
			if strings.TrimSpace(hg.HostGroupName) != hgName {
				continue
			}
			matches = append(matches, hg.HostGroupNumber)
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("hostgroup %q not found on port %q (serial=%d)", hgName, portID, serial)
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("hostgroup_name %q is not unique on port %q; import using '<serial>/<port_id>,<hostgroup_number>,<hostgroup_name>'", hgName, portID)
		}
		hgNum = matches[0]
	}

	if err := d.Set("hostgroup_number", hgNum); err != nil {
		return nil, err
	}
	// Set the native resource ID format used by Read/Update.
	d.SetId(fmt.Sprintf("%s,%d,%s", portID, hgNum, hgName))
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importVspSnapshotGroup expects: <serial>/<snapshot_group_name>
// Both serial and snapshot_group_name are required by Read to locate the group on storage.
func importVspSnapshotGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := strings.TrimSpace(d.Id())
	if importID == "" {
		return nil, fmt.Errorf("import id is empty")
	}
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<snapshot_group_name>'", importID)
	}
	serial, err := importParseSerial(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, err
	}
	groupName := strings.TrimSpace(parts[1])
	if groupName == "" {
		return nil, fmt.Errorf("invalid import id %q: snapshot_group_name must not be empty", importID)
	}
	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if err := d.Set("snapshot_group_name", groupName); err != nil {
		return nil, err
	}
	// Native ID for the resource is the group name.
	d.SetId(groupName)
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

// importVspIscsiChapUser expects: <serial>/<port_id>/<iscsi_target_number>/<chap_user_type>/<chap_user_name>
func importVspIscsiChapUser(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := strings.TrimSpace(d.Id())
	if importID == "" {
		return nil, fmt.Errorf("invalid import id: empty")
	}
	parts := strings.Split(importID, "/")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid import id %q; expected '<serial>/<port_id>/<iscsi_target_number>/<chap_user_type>/<chap_user_name>'", importID)
	}
	serial, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid serial %q in import id: %w", parts[0], err)
	}
	portID := strings.TrimSpace(parts[1])
	itNum, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return nil, fmt.Errorf("invalid iscsi_target_number %q in import id: %w", parts[2], err)
	}
	chapType := strings.ToLower(strings.TrimSpace(parts[3]))
	if chapType != "target" && chapType != "initiator" {
		return nil, fmt.Errorf("invalid chap_user_type %q in import id (must be 'target' or 'initiator')", parts[3])
	}
	chapName := strings.TrimSpace(parts[4])
	if chapName == "" {
		return nil, fmt.Errorf("invalid chap_user_name %q in import id", parts[4])
	}

	if err := d.Set("serial", serial); err != nil {
		return nil, err
	}
	if err := d.Set("port_id", portID); err != nil {
		return nil, err
	}
	if err := d.Set("iscsi_target_number", itNum); err != nil {
		return nil, err
	}
	if err := d.Set("chap_user_type", chapType); err != nil {
		return nil, err
	}
	if err := d.Set("chap_user_name", chapName); err != nil {
		return nil, err
	}

	// Resource Read will reset the ID to ChapUserID.
	d.SetId("import")
	return []*schema.ResourceData{d}, nil
}
