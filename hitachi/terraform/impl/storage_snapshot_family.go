package terraform

import (
	"fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceVspSnapshotFamilyRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	finalLdev, err := terrcommon.ExtractLdevFields(d, "ldev_id", "ldev_id_hex")
	if err != nil {
		return diag.FromErr(err)
	}
	if finalLdev == nil {
		return diag.FromErr(fmt.Errorf("either ldev_id or ldev_id_hex must be specified"))
	}
	ldevID := *finalLdev

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	families, err := reconObj.ReconcileGetFamily(ldevID)
	if err != nil {
		log.WriteError("TFError| Failed to fetch snapshot family for LDEV %d: %v", ldevID, err)
		return diag.FromErr(err)
	}

	if len(families) > 0 {
		members := make([]map[string]interface{}, 0, len(families))
		for _, f := range families {
			row := map[string]interface{}{
				"ldev_id":                        f.LdevID,
				"ldev_id_hex":                    utils.IntToHexString(f.LdevID),
				"snapshot_group_name":            f.SnapshotGroupName,
				"primary_or_secondary":           f.PrimaryOrSecondary,
				"status":                         f.Status,
				"pvol_ldev_id":                   f.PvolLdevID,
				"pvol_ldev_id_hex":               utils.IntToHexString(f.PvolLdevID),
				"svol_ldev_id":                   f.SvolLdevID,
				"svol_ldev_id_hex":               utils.IntToHexString(f.SvolLdevID),
				"mirror_unit_id":                 f.MuNumber,
				"pool_id":                        f.PoolID,
				"is_virtual_clone_volume":        f.IsVirtualCloneVolume,
				"is_virtual_clone_parent_volume": f.IsVirtualCloneParentVolume,
				"split_time":                     f.SplitTime,
				"parent_ldev_id":                 f.ParentLdevID,
				"snapshot_group_id":              f.SnapshotGroupID,
				"snapshot_id":                    f.SnapshotID,
			}
			members = append(members, row)
		}

		if err := d.Set("family_members", members); err != nil {
			return diag.FromErr(err)
		}
		d.Set("total_members", len(members))
	} else {
		d.Set("total_members", 0)
		d.Set("family_members", []interface{}{})
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

func DatasourceVspVirtualCloneParentVolumeRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	parentLdevIds, err := reconObj.ReconcileGetVirtualCloneParentVolumes()
	if err != nil {
		return diag.FromErr(err)
	}

	if len(parentLdevIds) > 0 {
		// Set the plain integer slice directly to the schema
		if err := d.Set("parent_volumes", parentLdevIds); err != nil {
			return diag.FromErr(err)
		}
		// Convert the integer slice to a slice of Hex strings
		parentLdevIdsHex := make([]string, 0, len(parentLdevIds))
		for _, id := range parentLdevIds {
			parentLdevIdsHex = append(parentLdevIdsHex, utils.IntToHexString(id))
		}
		if err := d.Set("parent_volumes_hex", parentLdevIdsHex); err != nil {
			return diag.FromErr(err)
		}
		d.Set("parent_count", len(parentLdevIds))
	} else {
		d.Set("parent_volumes", []int{})
		d.Set("parent_count", 0)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}
