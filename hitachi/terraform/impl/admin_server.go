package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	recmanager "terraform-provider-hitachi/hitachi/storage/admin/reconciler"
	recimpl "terraform-provider-hitachi/hitachi/storage/admin/reconciler/impl"
	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Resource CRUD Operations -------------------

func ResourceAdminServerCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Build create parameters
	params, err := buildCreateAdminServerParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Build add hostgroups parameters
	addHgParams, err := buildAddHostGroupsToServerParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverID, err := recObj.ReconcileCreateAdminServer(params, addHgParams)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%d", serial, serverID))

	log.WriteInfo("VSP One server created successfully with ID: %d", serverID)
	return ResourceAdminServerRead(d)
}

func ResourceAdminServerRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Extract server ID from resource ID
	serverID, err := extractServerIDFromResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverInfo, err := recObj.ReconcileReadAdminServer(serverID)
	if err != nil {
		// If server not found, remove from state
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Server not found",
			Detail:   fmt.Sprintf("Server with ID %d not found, removing from state", serverID),
		}}
	}

	// Set all computed attributes
	if err := setAdminServerAttributes(d, serverInfo); err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("VSP One server read successfully: %d", serverID)
	return nil
}

func ResourceAdminServerUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Extract server ID from resource ID
	serverID, err := extractServerIDFromResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Build update parameters
	params := buildUpdateAdminServerParams(d)

	// Build add hostgroups parameters
	addHgParams, err := buildAddHostGroupsToServerParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = recObj.ReconcileUpdateAdminServer(serverID, params, addHgParams)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("VSP One server updated successfully: %d", serverID)
	return ResourceAdminServerRead(d)
}

func ResourceAdminServerDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Extract server ID from resource ID
	serverID, err := extractServerIDFromResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Build delete parameters
	params := gwymodel.DeleteAdminServerParams{
		KeepLunConfig: d.Get("keep_lun_config").(bool),
	}

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Delete all paths before deleting the server
	log.WriteInfo("Deleting all paths for server %d before server deletion", serverID)
	err = deleteAllServerPaths(recObj, serverID)
	if err != nil {
		log.WriteWarn("Failed to delete some paths for server %d: %v", serverID, err)
		// Continue with server deletion even if path deletion fails
	}

	err = recObj.ReconcileDeleteAdminServer(serverID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Clear the resource ID
	d.SetId("")
	log.WriteInfo("VSP One server deleted successfully: %d", serverID)
	return nil
}

// ------------------- Data Source Operations -------------------

func DataSourceAdminServerListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	params := gwymodel.AdminServerListParams{}

	if nickname, ok := d.GetOk("nickname"); ok {
		nick := nickname.(string)
		params.Nickname = &nick
	}

	if hbaWwn, ok := d.GetOk("hba_wwn"); ok {
		wwn := hbaWwn.(string)
		params.HbaWwn = &wwn
	}
	if iscsiName, ok := d.GetOk("iscsi_name"); ok {
		iscsi := iscsiName.(string)
		params.IscsiName = &iscsi
	}

	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Params: %v", string(b))

	// call provisioner directly
	provObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return diag.FromErr(err)
	}

	resp, err := provObj.GetAdminServerList(params)
	if err != nil {
		log.WriteError("Failed to get VSP One server list: %v", err)
		return diag.FromErr(err)
	}

	log.WriteInfo("Successfully retrieved server list: %+v", resp)

	// Set the ID for the resource
	d.SetId(fmt.Sprintf("admin_server_list_%d", serial))

	// Set the attributes in the schema
	if err := d.Set("data", flattenAdminServerListResponse(resp)); err != nil {
		log.WriteError("Failed to set data attribute: %v", err)
		return diag.FromErr(err)
	}
	if err := d.Set("server_count", len(resp.Data)); err != nil {
		log.WriteError("Failed to set server_count attribute: %v", err)
		return diag.FromErr(err)
	}

	return nil
}

func DataSourceAdminServerInfoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)

	// call provisioner directly
	provObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		log.WriteDebug("TFError| error in getReconcilerManagerForAdminServer, err: %v", err)
		return diag.FromErr(err)
	}

	resp, err := provObj.GetAdminServerInfo(serverID)
	if err != nil {
		log.WriteError("Failed to get VSP One server info: %v", err)
		return diag.FromErr(err)
	}

	log.WriteInfo("Successfully retrieved server info: %+v", resp)
	// Set the ID for the resource
	d.SetId(fmt.Sprintf("admin_server_info_%d_%d", serial, serverID))

	// Set the attributes in the schema
	if err := d.Set("data", []map[string]interface{}{flattenAdminServerInfoResponse(resp)}); err != nil {
		log.WriteError("Failed to set info attribute: %v", err)
		return diag.FromErr(err)
	}
	return nil
}

// ------------------- Helper Functions -------------------

func buildCreateAdminServerParams(d *schema.ResourceData) (gwymodel.CreateAdminServerParams, error) {
	params := gwymodel.CreateAdminServerParams{
		ServerNickname: d.Get("server_nickname").(string),
		IsReserved:     d.Get("is_reserved").(bool),
	}

	// Protocol and OsType are required if not reserved
	if !params.IsReserved {
		if protocol, ok := d.GetOk("protocol"); ok {
			params.Protocol = protocol.(string)
		} else {
			return params, fmt.Errorf("protocol is required when is_reserved is false")
		}

		if osType, ok := d.GetOk("os_type"); ok {
			params.OsType = osType.(string)
		} else {
			return params, fmt.Errorf("os_type is required when is_reserved is false")
		}
	}

	// OS type options are optional
	if osTypeOptions, ok := d.GetOk("os_type_options"); ok {
		optionsList := osTypeOptions.([]interface{})
		params.OsTypeOptions = make([]int, len(optionsList))
		for i, opt := range optionsList {
			params.OsTypeOptions[i] = opt.(int)
		}
	}

	return params, nil
}

func buildUpdateAdminServerParams(d *schema.ResourceData) gwymodel.UpdateAdminServerParams {
	params := gwymodel.UpdateAdminServerParams{}

	if d.HasChange("server_nickname") {
		params.Nickname = d.Get("server_nickname").(string)
	}

	if d.HasChange("os_type") {
		params.OsType = d.Get("os_type").(string)
	}

	if d.HasChange("os_type_options") {
		if osTypeOptions, ok := d.GetOk("os_type_options"); ok {
			optionsList := osTypeOptions.([]interface{})
			params.OsTypeOptions = make([]int, len(optionsList))
			for i, opt := range optionsList {
				params.OsTypeOptions[i] = opt.(int)
			}
		}
	}

	return params
}

// for add/sync hostgroups
func buildAddHostGroupsToServerParams(d *schema.ResourceData) (gwymodel.AddHostGroupsToServerParam, error) {
	params := gwymodel.AddHostGroupsToServerParam{}

	// If host_groups not defined at all
	if _, ok := d.GetOk("host_groups"); !ok {
		return params, nil
	}

	// Check if host_groups changed in Terraform plan
	if !d.HasChange("host_groups") {
		// No changes → return empty (nothing to add)
		return params, nil
	}

	oldRaw, newRaw := d.GetChange("host_groups")
	oldList := convertToHostGroupList(oldRaw)
	newList := convertToHostGroupList(newRaw)

	// Build a set of old host groups for quick lookup
	oldSet := make(map[string]struct{}, len(oldList))
	for _, oldHg := range oldList {
		key := buildHostGroupKey(oldHg.PortID, oldHg.HostGroupID, oldHg.HostGroupName)
		oldSet[key] = struct{}{}
	}

	// Collect only the new ones that are not in the old set
	for _, newHg := range newList {
		key := buildHostGroupKey(newHg.PortID, newHg.HostGroupID, newHg.HostGroupName)
		if _, exists := oldSet[key]; !exists {
			params.HostGroups = append(params.HostGroups, newHg)
		}
	}

	return params, nil
}

func convertToHostGroupList(raw interface{}) []gwymodel.HostGroupForAddToServerParam {
	list := []gwymodel.HostGroupForAddToServerParam{}

	if raw == nil {
		return list
	}

	items, ok := raw.([]interface{})
	if !ok {
		return list
	}

	for _, item := range items {
		if item == nil {
			continue
		}

		m := item.(map[string]interface{})
		hg := gwymodel.HostGroupForAddToServerParam{}

		if portID, ok := m["port_id"].(string); ok && portID != "" {
			hg.PortID = portID
		}

		if id, ok := m["host_group_id"].(int); ok && id > 0 {
			hg.HostGroupID = &id
		}

		if name, ok := m["host_group_name"].(string); ok && name != "" {
			hg.HostGroupName = &name
		}

		list = append(list, hg)
	}
	return list
}

// Create a unique key for comparison (port + identifier)
func buildHostGroupKey(portID string, id *int, name *string) string {
	if id != nil {
		return fmt.Sprintf("%s-id-%d", portID, *id)
	}
	if name != nil {
		return fmt.Sprintf("%s-name-%s", portID, *name)
	}
	return portID
}

func extractServerIDFromResourceID(resourceID string) (int, error) {
	if resourceID == "" {
		return 0, fmt.Errorf("resource ID is empty")
	}

	// Resource ID format: "serial-serverid"
	parts := strings.Split(resourceID, "-")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid resource ID format: %s", resourceID)
	}

	serverID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid server ID in resource ID: %s", parts[1])
	}

	return serverID, nil
}

func setAdminServerAttributes(d *schema.ResourceData, serverInfo *gwymodel.AdminServerInfo) error {
	// Create the data block data
	serverInfoData := map[string]interface{}{
		"server_id":                     serverInfo.ID,
		"nickname":                      serverInfo.Nickname,
		"protocol":                      serverInfo.Protocol,
		"os_type":                       serverInfo.OsType,
		"os_type_options":               serverInfo.OsTypeOptions,
		"total_capacity":                serverInfo.TotalCapacity,
		"used_capacity":                 serverInfo.UsedCapacity,
		"number_of_volumes":             serverInfo.NumberOfVolumes,
		"number_of_paths":               serverInfo.NumberOfPaths,
		"paths":                         flattenAdminServerPaths(serverInfo.Paths),
		"is_inconsistent":               serverInfo.IsInconsistent,
		"is_reserved":                   serverInfo.IsReserved,
		"has_non_fullmesh_lu_paths":     serverInfo.HasNonFullmeshLuPaths,
		"has_unaligned_os_types":        serverInfo.HasUnalignedOsTypes,
		"has_unaligned_os_type_options": serverInfo.HasUnalignedOsTypeOptions,
	}

	// Set the data block as a list with one item
	if err := d.Set("data", []map[string]interface{}{serverInfoData}); err != nil {
		return fmt.Errorf("failed to set data: %w", err)
	}

	return nil
}

func getReconcilerManagerForAdminServer(serial int) (recmanager.AdminStorageManager, error) {
	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := recmodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	return recimpl.NewEx(setting)
}

// flattenAdminServerListResponse flattens the admin server list response for the Terraform schema
func flattenAdminServerListResponse(resp *gwymodel.AdminServerListResponse) []map[string]interface{} {
	if resp == nil || len(resp.Data) == 0 {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, len(resp.Data))
	for i, server := range resp.Data {
		item := map[string]interface{}{
			"server_id":                server.ID,
			"nickname":                 server.Nickname,
			"protocol":                 server.Protocol,
			"os_type":                  server.OsType,
			"total_capacity":           server.TotalCapacity,
			"used_capacity":            server.UsedCapacity,
			"number_of_paths":          server.NumberOfPaths,
			"is_inconsistent":          server.IsInconsistent,
			"modification_in_progress": server.ModificationInProgress,
			"is_reserved":              server.IsReserved,
			"has_unaligned_os_types":   server.HasUnalignedOsTypes,
		}
		result[i] = item
	}
	return result
}

// flattenAdminServerInfoResponse flattens the admin server info response for the Terraform schema
func flattenAdminServerInfoResponse(resp *gwymodel.AdminServerInfo) map[string]interface{} {
	if resp == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"nickname":                      resp.Nickname,
		"protocol":                      resp.Protocol,
		"os_type":                       resp.OsType,
		"os_type_options":               resp.OsTypeOptions,
		"total_capacity":                resp.TotalCapacity,
		"used_capacity":                 resp.UsedCapacity,
		"number_of_volumes":             resp.NumberOfVolumes,
		"number_of_paths":               resp.NumberOfPaths,
		"paths":                         flattenAdminServerPaths(resp.Paths),
		"is_inconsistent":               resp.IsInconsistent,
		"is_reserved":                   resp.IsReserved,
		"has_non_fullmesh_lu_paths":     resp.HasNonFullmeshLuPaths,
		"has_unaligned_os_types":        resp.HasUnalignedOsTypes,
		"has_unaligned_os_type_options": resp.HasUnalignedOsTypeOptions,
	}
}

// flattenAdminServerPaths flattens the paths slice for the Terraform schema
func flattenAdminServerPaths(paths []gwymodel.AdminServerPath) []map[string]interface{} {
	result := make([]map[string]interface{}, len(paths))
	for i, path := range paths {
		result[i] = map[string]interface{}{
			"hba_wwn":    path.HbaWwn,
			"iscsi_name": path.IscsiName,
			"port_ids":   path.PortIds,
		}
	}
	return result
}

// buildSetAdminServerPathParams builds parameters for setting server path
func buildSetAdminServerPathParams(pathMap map[string]interface{}) (gwymodel.SetAdminServerPathParams, error) {
	params := gwymodel.SetAdminServerPathParams{}

	// Validate that either hba_wwn or iscsi_name is specified
	hbaWwn, hasHbaWwn := pathMap["hba_wwn"]
	iscsiName, hasIscsiName := pathMap["iscsi_name"]

	if !hasHbaWwn && !hasIscsiName {
		return params, fmt.Errorf("either hba_wwn or iscsi_name must be specified")
	}

	if hasHbaWwn && hasIscsiName {
		hbaWwnStr := hbaWwn.(string)
		iscsiNameStr := iscsiName.(string)
		if hbaWwnStr != "" && iscsiNameStr != "" {
			return params, fmt.Errorf("only one of hba_wwn or iscsi_name can be specified, not both")
		}
	}

	// Set the values
	if hasHbaWwn {
		params.HbaWwn = hbaWwn.(string)
	}
	if hasIscsiName {
		params.IscsiName = iscsiName.(string)
	}

	// Handle port_ids
	if portIds, ok := pathMap["port_ids"]; ok && portIds != nil {
		portIdsList := portIds.([]interface{})
		params.PortIds = make([]string, len(portIdsList))
		for i, portId := range portIdsList {
			params.PortIds[i] = portId.(string)
		}
	}

	return params, nil
}

// buildDeleteAdminServerPathParams builds parameters for deleting server path
func buildDeleteAdminServerPathParams(hbaWwn, iscsiName, portId string) (gwymodel.DeleteAdminServerPathParams, error) {
	params := gwymodel.DeleteAdminServerPathParams{
		PortId: portId,
	}

	// Validate that either hba_wwn or iscsi_name is specified
	if hbaWwn == "" && iscsiName == "" {
		return params, fmt.Errorf("either hba_wwn or iscsi_name must be specified")
	}

	if hbaWwn != "" && iscsiName != "" {
		return params, fmt.Errorf("only one of hba_wwn or iscsi_name can be specified, not both")
	}

	// Set the values
	if hbaWwn != "" {
		params.HbaWwn = hbaWwn
	}
	if iscsiName != "" {
		params.IscsiName = iscsiName
	}

	return params, nil
}

// deleteAllServerPaths deletes all paths for a server before deletion
func deleteAllServerPaths(recObj recmanager.AdminStorageManager, serverID int) error {
	log := commonlog.GetLogger()

	// Get current server info to get existing paths
	serverInfo, err := recObj.GetAdminServerInfo(serverID)
	if err != nil {
		log.WriteWarn("Failed to get server info for path deletion: %v", err)
		return nil // Continue with server deletion even if we can't get paths
	}

	if serverInfo == nil || len(serverInfo.Paths) == 0 {
		log.WriteInfo("No paths found for server %d, skipping path deletion", serverID)
		return nil
	}

	log.WriteInfo("Deleting %d paths for server %d before server deletion", len(serverInfo.Paths), serverID)

	// Delete each path individually
	for _, path := range serverInfo.Paths {
		for _, portId := range path.PortIds {
			deleteParams, err := buildDeleteAdminServerPathParams(path.HbaWwn, path.IscsiName, portId)
			if err != nil {
				log.WriteError("Failed to build delete path parameters: %v", err)
				continue // Continue with other paths
			}

			err = recObj.ReconcileDeleteAdminServerPath(serverID, deleteParams)
			if err != nil {
				log.WriteError("Failed to delete path for server %d: %v", serverID, err)
				// Continue with other paths - don't fail the entire operation
			} else {
				log.WriteInfo("Successfully deleted path for server %d (port: %s)", serverID, portId)
			}
		}
	}

	return nil
}

// getPathKey creates a unique key for a path based on hba_wwn or iscsi_name
func getPathKey(pathMap map[string]interface{}) string {
	if hbaWwn, ok := pathMap["hba_wwn"]; ok && hbaWwn.(string) != "" {
		return "hba:" + hbaWwn.(string)
	}
	if iscsiName, ok := pathMap["iscsi_name"]; ok && iscsiName.(string) != "" {
		return "iscsi:" + iscsiName.(string)
	}
	return ""
}

// parsePathKey parses a path key back to hba_wwn and iscsi_name
func parsePathKey(key string) (string, string) {
	if strings.HasPrefix(key, "hba:") {
		return strings.TrimPrefix(key, "hba:"), ""
	}
	if strings.HasPrefix(key, "iscsi:") {
		return "", strings.TrimPrefix(key, "iscsi:")
	}
	return "", ""
}

// extractPortIds extracts port IDs from a path map
func extractPortIds(pathMap map[string]interface{}) []string {
	if portIds, ok := pathMap["port_ids"]; ok && portIds != nil {
		portIdsList := portIds.([]interface{})
		result := make([]string, len(portIdsList))
		for i, portId := range portIdsList {
			result[i] = portId.(string)
		}
		return result
	}
	return []string{}
}

// getRemovedPorts returns ports that are in oldPorts but not in newPorts
func getRemovedPorts(oldPorts, newPorts []string) []string {
	newPortsMap := make(map[string]bool)
	for _, port := range newPorts {
		newPortsMap[port] = true
	}

	var removed []string
	for _, port := range oldPorts {
		if !newPortsMap[port] {
			removed = append(removed, port)
		}
	}
	return removed
}

func DataSourceAdminServerPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverIDInput := d.Get("server_id").(int)
	hbaWwn := d.Get("hba_wwn").(string)
	iscsiName := d.Get("iscsi_name").(string)
	portId := d.Get("port_id").(string)

	// Validation: Either hba_wwn or iscsi_name must be specified
	if hbaWwn == "" && iscsiName == "" {
		return diag.FromErr(fmt.Errorf("either hba_wwn or iscsi_name must be specified"))
	}

	if hbaWwn != "" && iscsiName != "" {
		return diag.FromErr(fmt.Errorf("only one of hba_wwn or iscsi_name can be specified, not both"))
	}

	// Build parameters
	params := gwymodel.AdminServerPathParams{
		ServerID:  serverIDInput,
		HbaWwn:    hbaWwn,
		IscsiName: iscsiName,
		PortId:    portId,
	}

	// Get reconciler manager
	provObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		log.WriteDebug("TFError| error in getReconcilerManagerForAdminServer, err: %v", err)
		return diag.FromErr(err)
	}

	resp, err := provObj.GetAdminServerPath(params)
	if err != nil {
		log.WriteError("Failed to get VSP One server path: %v", err)
		return diag.FromErr(err)
	}

	log.WriteInfo("Successfully retrieved server path: %+v", resp)

	// Set the ID for the resource
	d.SetId(fmt.Sprintf("admin_server_path_%d_%s", serverIDInput, resp.ID))

	// Set the attributes in the schema
	pathData := map[string]interface{}{
		"id":         resp.ID,
		"server_id":  resp.ServerID,
		"hba_wwn":    resp.HbaWwn,
		"iscsi_name": resp.IscsiName,
		"port_id":    resp.PortId,
	}

	if err := d.Set("data", []map[string]interface{}{pathData}); err != nil {
		log.WriteError("Failed to set data attribute: %v", err)
		return diag.FromErr(err)
	}

	return nil
}

// ------------------- Server Path Resource CRUD Operations -------------------

func ResourceAdminServerPathCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)
	hbaWwn := d.Get("hba_wwn").(string)
	iscsiName := d.Get("iscsi_name").(string)

	// Validation: Either hba_wwn or iscsi_name must be specified
	if hbaWwn == "" && iscsiName == "" {
		return diag.FromErr(fmt.Errorf("either hba_wwn or iscsi_name must be specified"))
	}

	if hbaWwn != "" && iscsiName != "" {
		return diag.FromErr(fmt.Errorf("only one of hba_wwn or iscsi_name can be specified, not both"))
	}

	// Get port IDs
	portIdsInterface := d.Get("port_ids").([]interface{})
	portIds := make([]string, len(portIdsInterface))
	for i, portId := range portIdsInterface {
		portIds[i] = portId.(string)
	}

	// Build parameters
	params := gwymodel.SetAdminServerPathParams{
		HbaWwn:    hbaWwn,
		IscsiName: iscsiName,
		PortIds:   portIds,
	}

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	err = recObj.ReconcileSetAdminServerPath(serverID, params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create server path: %w", err))
	}

	// Create resource ID
	pathIdentifier := hbaWwn
	if pathIdentifier == "" {
		pathIdentifier = iscsiName
	}
	d.SetId(fmt.Sprintf("%d-%d-%s", serial, serverID, pathIdentifier))

	log.WriteInfo("Server path created successfully for server %d", serverID)
	return ResourceAdminServerPathRead(d)
}

func ResourceAdminServerPathRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)
	hbaWwn := d.Get("hba_wwn").(string)
	iscsiName := d.Get("iscsi_name").(string)

	// We need to get server info to check if the path exists
	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverInfo, err := recObj.GetAdminServerInfo(serverID)
	if err != nil {
		// If server not found, remove from state
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Server not found",
			Detail:   fmt.Sprintf("Server with ID %d not found, removing path from state", serverID),
		}}
	}

	// Check if the path exists in the server info
	pathExists := false
	for _, path := range serverInfo.Paths {
		if (hbaWwn != "" && path.HbaWwn == hbaWwn) || (iscsiName != "" && path.IscsiName == iscsiName) {
			pathExists = true
			break
		}
	}

	if !pathExists {
		// Path no longer exists, remove from state
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Path not found",
			Detail:   fmt.Sprintf("Path not found on server %d, removing from state", serverID),
		}}
	}

	log.WriteInfo("Server path read successfully for server %d", serverID)
	return nil
}

func ResourceAdminServerPathUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)

	// Handle port_ids changes
	if d.HasChange("port_ids") {
		oldPortIds, newPortIds := d.GetChange("port_ids")
		oldPortsList := oldPortIds.([]interface{})
		newPortsList := newPortIds.([]interface{})

		// Convert to string slices
		oldPorts := make([]string, len(oldPortsList))
		newPorts := make([]string, len(newPortsList))
		for i, port := range oldPortsList {
			oldPorts[i] = port.(string)
		}
		for i, port := range newPortsList {
			newPorts[i] = port.(string)
		}

		hbaWwn := d.Get("hba_wwn").(string)
		iscsiName := d.Get("iscsi_name").(string)

		recObj, err := getReconcilerManagerForAdminServer(serial)
		if err != nil {
			return diag.FromErr(err)
		}

		// Remove old ports that are not in new ports
		removedPorts := getRemovedPorts(oldPorts, newPorts)
		for _, portId := range removedPorts {
			deleteParams, err := buildDeleteAdminServerPathParams(hbaWwn, iscsiName, portId)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to build delete parameters: %w", err))
			}

			err = recObj.ReconcileDeleteAdminServerPath(serverID, deleteParams)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to delete path port %s: %w", portId, err))
			}
		}

		// Add new ports (this will handle all ports, existing ones will be ignored by the API)
		setParams := gwymodel.SetAdminServerPathParams{
			HbaWwn:    hbaWwn,
			IscsiName: iscsiName,
			PortIds:   newPorts,
		}

		err = recObj.ReconcileSetAdminServerPath(serverID, setParams)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update server path: %w", err))
		}
	}

	log.WriteInfo("Server path updated successfully for server %d", serverID)
	return ResourceAdminServerPathRead(d)
}

func ResourceAdminServerPathDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)
	hbaWwn := d.Get("hba_wwn").(string)
	iscsiName := d.Get("iscsi_name").(string)

	// Get port IDs
	portIdsInterface := d.Get("port_ids").([]interface{})

	recObj, err := getReconcilerManagerForAdminServer(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Delete all ports for this path
	for _, portIdInterface := range portIdsInterface {
		portId := portIdInterface.(string)
		deleteParams, err := buildDeleteAdminServerPathParams(hbaWwn, iscsiName, portId)
		if err != nil {
			log.WriteError("Failed to build delete parameters for port %s: %v", portId, err)
			continue
		}

		err = recObj.ReconcileDeleteAdminServerPath(serverID, deleteParams)
		if err != nil {
			log.WriteError("Failed to delete path port %s: %v", portId, err)
			// Continue with other ports
		} else {
			log.WriteInfo("Successfully deleted path port %s for server %d", portId, serverID)
		}
	}

	// Clear the resource ID
	d.SetId("")
	log.WriteInfo("Server path deleted successfully for server %d", serverID)
	return nil
}
