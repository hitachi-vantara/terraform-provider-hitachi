package terraform

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Datasources -------------------

func DatasourceAdminOneIscsiTargetRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverId := d.Get("server_id").(int)
	portId := d.Get("port_id").(string)

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	iscsiTargetInfo, err := provObj.GetIscsiTargetByPort(serverId, portId)
	if err != nil {
		log.WriteDebug("failed to get iscsi_target serverId:%v portId:%v err:%v", serverId, portId, err)
		return diag.FromErr(fmt.Errorf("failed to get iscsi_target serverId:%v portId:%v err:%v", serverId, portId, err))
	}

	log.WriteDebug("iscsi_target %+v", iscsiTargetInfo)
	if iscsiTargetInfo == nil {
		// no matching resource — clear state if needed (return empty diags)
		d.SetId("")
		return nil
	}

	if err := d.Set("iscsi_target_info", []map[string]interface{}{convertOneIscsiTargetInfoToSchema(iscsiTargetInfo)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set iscsi_target_info: %w", err))
	}

	// Use serverId/portId as the Terraform state ID
	d.SetId(fmt.Sprintf("%d/%s", serverId, portId))

	return nil
}

func DatasourceAdminMultipleIscsiTargetsRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverId := d.Get("server_id").(int)

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	iscsiTargets, err := provObj.GetIscsiTargets(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	if iscsiTargets == nil {
		// No data returned
		d.SetId("")
		return nil
	}

	if err := d.Set("iscsi_targets_info", convertMultipleIscsiTargetInfosListToSchema(iscsiTargets)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iscsi_targets_count", iscsiTargets.Count); err != nil {
		return diag.FromErr(err)
	}

	// Use a time-based id so datasource re-evaluates when requested
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// ------------------- Resource -------------------

func ResourceAdminIscsiTargetRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var diags diag.Diagnostics

	serial := d.Get("serial").(int)
	serverId := d.Get("server_id").(int)
	portId := d.Get("port_id").(string)

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	iscsiTargetInfo, err := provObj.GetIscsiTargetByPort(serverId, portId)
	if err != nil {
		log.WriteDebug("failed to get iscsi_target serverId:%v portId:%v err:%v", serverId, portId, err)
		return diag.FromErr(fmt.Errorf("failed to get iscsi_target serverId:%v portId:%v err:%v", serverId, portId, err))
	}

	log.WriteDebug("iscsi_target %+v", iscsiTargetInfo)
	if iscsiTargetInfo == nil {
		// Resource not found - clear state
		d.SetId("")
		return diags
	}

	if err := d.Set("iscsi_target_info", []map[string]interface{}{convertOneIscsiTargetInfoToSchema(iscsiTargetInfo)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set iscsi_target_info: %w", err))
	}

	// ensure state id follows serverId/portId format
	d.SetId(fmt.Sprintf("%d/%s", serverId, portId))

	log.WriteInfo("iscsi_target read successfully")
	return diags
}

func ResourceAdminIscsiTargetDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// No backend delete — just remove from state
	d.SetId("")
	log.WriteInfo("iscsi_target deleted from state")
	return nil
}

func ResourceAdminIscsiTargetCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting iscsi_target update")

	serial := d.Get("serial").(int)
	serverId := d.Get("server_id").(int)
	portId := d.Get("port_id").(string)
	targetIscsiName := d.Get("target_iscsi_name").(string)

	if targetIscsiName == "" {
		// nothing to do, re-read
		return ResourceAdminIscsiTargetRead(d)
	}

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// perform the rename (server-side change)
	if targetIscsiName != "" {
		if err := provObj.ChangeIscsiTargetName(serverId, portId, targetIscsiName); err != nil {
			return diag.FromErr(err)
		}
	}

	// set state id
	d.SetId(fmt.Sprintf("%d/%s", serverId, portId))

	log.WriteInfo("iscsi_target updated successfully")
	// return read to populate computed fields
	return ResourceAdminIscsiTargetRead(d)
}

// ------------------- Helpers -------------------
// convert a single IscsiTargetInfoByPort to schema
func convertOneIscsiTargetInfoToSchema(v *gwymodel.IscsiTargetInfoByPort) map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	m := map[string]interface{}{
		"port_id": v.PortID,
	}

	if v.TargetIscsiName != nil {
		m["target_iscsi_name"] = *v.TargetIscsiName
	} else {
		// keep it absent or set to empty string (schema has Computed=true, it's fine to omit)
		// We'll set nil behavior by leaving it out; Terraform will treat it as null/not-present.
	}

	log.WriteDebug("Convert single iscsi target: %+v", m)
	return m
}

// convert the gwymodel.IscsiTargetInfoList to []map[string]interface{} for Terraform
func convertMultipleIscsiTargetInfosListToSchema(iscsiTargets *gwymodel.IscsiTargetInfoList) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if iscsiTargets == nil || len(iscsiTargets.Data) == 0 {
		return nil
	}

	// Sort by PortID for deterministic order
	sort.Slice(iscsiTargets.Data, func(i, j int) bool {
		return strings.Compare(iscsiTargets.Data[i].PortID, iscsiTargets.Data[j].PortID) < 0
	})

	out := make([]map[string]interface{}, len(iscsiTargets.Data))
	for i, v := range iscsiTargets.Data {
		out[i] = convertOneIscsiTargetInfoToSchema(&v)
	}

	log.WriteDebug("Convert list to schema: %+v", out)
	return out
}
