package terraform

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	recmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- SnapshotGroups Datasources -------------------

func DatasourceVspSnapshotGroupRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	groupName := d.Get("snapshot_group_name").(string)

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := reconObj.ReconcileGetSnapshotGroup(groupName)
	if err != nil {
		log.WriteDebug("failed to get snapshot group %s: %v", groupName, err)
		return diag.FromErr(fmt.Errorf("failed to get snapshot group: %v", err))
	}

	if err := d.Set("snapshot_group", convertSnapshotGroupsToSchema([]gwymodel.SnapshotGroup{*group})); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_group: %w", err))
	}

	if err := d.Set("snapshot_count", len(group.Snapshots)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_count: %w", err))
	}

	d.SetId(group.SnapshotGroupName)
	return nil
}

func DatasourceVspMultipleSnapshotGroupsRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	includePairs := d.Get("include_pairs").(bool)

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetches all groups with pair details
	groups, err := reconObj.ReconcileGetMultipleSnapshotGroups(includePairs)
	if err != nil {
		log.WriteDebug("failed to get snapshot group list: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get snapshot groups: %v", err))
	}

	schemaData := convertSnapshotGroupsToSchema(groups.Data)
	if err := d.Set("snapshot_groups", schemaData); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_groups: %w", err))
	}

	// Count total snapshot groups
	groupCount := len(groups.Data)

	// Sum up all snapshot pairs across all groups
	totalSnapshotPairs := 0
	for _, group := range groups.Data {
		// SnapshotGroup contains a slice of Snapshots (the member pairs)
		totalSnapshotPairs += len(group.Snapshots)
	}

	// Set the count of the groups themselves
	if err := d.Set("snapshotgroup_count", groupCount); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshotgroup_count: %w", err))
	}

	// Set the total count of all snapshot pairs found within those groups
	if err := d.Set("snapshot_count", totalSnapshotPairs); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_count: %w", err))
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

// ------------------- Snapshot Group Resource -------------------

func ResourceVspSnapshotGroupRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	groupName := d.Get("snapshot_group_name").(string)

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := reconObj.ReconcileGetSnapshotGroup(groupName)
	if err != nil {
		// Handle 404: If the group is gone, remove it from Terraform state
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			log.WriteInfo("Snapshot group %s not found on storage, removing from state", groupName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("snapshot_group", convertSnapshotGroupsToSchema([]gwymodel.SnapshotGroup{*group})); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set snapshot_group: %w", err))
	}

	d.SetId(group.SnapshotGroupName)
	return nil
}

func ResourceVspSnapshotGroupApply(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	input := recmodel.SnapshotGroupReconcilerInput{
		SnapshotGroupName: terrcommon.GetStringPointer(d, "snapshot_group_name"),
		Action:            terrcommon.GetStringPointer(d, "state"),
		AutoSplit:         terrcommon.GetBoolPointer(d, "auto_split"),
		CopySpeed:         terrcommon.GetStringPointer(d, "copy_speed"),
		RetentionPeriod:   terrcommon.GetIntPointer(d, "retention_period_hours"),
	}

	// ReconcileSnapshotGroupVFamily now returns (SnapshotGroup, []SnapshotFamily, error)
	group, vclones, err := reconObj.ReconcileSnapshotGroupVFamily(input)
	if err != nil {
		return diag.FromErr(err)
	}

	return updateSnapshotGroupResourceState(d, group, vclones)
}

func ResourceVspSnapshotGroupDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	groupName := d.Id()

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Depending on your requirements, 'delete' might remove the pairs or just the group label
	input := recmodel.SnapshotGroupReconcilerInput{
		SnapshotGroupName: &groupName,
		Action:            utils.Ptr("delete"),
	}

	_, err = reconObj.ReconcileSnapshotGroup(input)
	if err != nil {
		// If it's already gone, don't fail the destroy
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// ------------------- Helper: Conversion -------------------

func updateSnapshotGroupResourceState(d *schema.ResourceData, group *gwymodel.SnapshotGroup, vclones []gwymodel.SnapshotFamily) diag.Diagnostics {
	// We must ensure the ID is set to prevent "Root object absent"
	groupName := d.Get("snapshot_group_name").(string)
	d.SetId(groupName)

	// 1. Handle standard Snapshot Group data
	if group != nil {
		if err := d.Set("snapshot_group", convertSnapshotGroupsToSchema([]gwymodel.SnapshotGroup{*group})); err != nil {
			return diag.FromErr(err)
		}
	} else {
		d.Set("snapshot_group", []map[string]interface{}{})
	}

	// 2. Handle promoted vClone data
	if len(vclones) > 0 {
		if err := d.Set("vclones", convertSnapshotFamiliesToSchema(vclones)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		d.Set("vclones", []map[string]interface{}{})
	}

	return nil
}

func convertSnapshotGroupsToSchema(groups []gwymodel.SnapshotGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		m := map[string]interface{}{
			"snapshot_group_name": g.SnapshotGroupName,
			"snapshot_group_id":   g.SnapshotGroupID,
			"snapshots":           convertSnapshotsToSchema(g.Snapshots),
		}
		result = append(result, m)
	}
	return result
}

func convertSnapshotFamiliesToSchema(families []gwymodel.SnapshotFamily) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(families))

	for _, v := range families {
		m := map[string]interface{}{
			"pvol_ldev_id":     v.ParentLdevID,
			"svol_ldev_id":     v.LdevID,
			"pvol_ldev_id_hex": utils.IntToHexString(v.ParentLdevID),
			"svol_ldev_id_hex": utils.IntToHexString(v.LdevID),

			// Attribute Flags
			"is_virtual_clone_volume":        v.IsVirtualCloneVolume,
			"is_virtual_clone_parent_volume": v.IsVirtualCloneParentVolume,

			// Metadata
			"split_time": v.SplitTime,
			"pool_id":    v.PoolID,
		}
		result = append(result, m)
	}

	return result
}
