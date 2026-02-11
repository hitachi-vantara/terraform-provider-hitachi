package terraform

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Data Sources -------------------

func DatasourceAdminOneVolumeServerConnectionRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	volumeId := d.Get("volume_id").(int)
	serverId := d.Get("server_id").(int)

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	connections, err := provObj.GetOneVolumeServerConnection(volumeId, serverId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get connections: %v", err))
	}

	if err := d.Set("connection_info", []map[string]interface{}{convertVolumeServerConnectionToSchema(connections)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set connection_info: %w", err))
	}

	// Use a stable ID composed of serial + volume + server
	d.SetId(fmt.Sprintf("%d-%d-%d", serial, volumeId, serverId))

	return nil
}

func DatasourceAdminMultipleVolumeServerConnectionsRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Validate that at least one of server_id or server_nickname is provided
	serverId, hasServerId := d.GetOk("server_id")
	serverNickname, hasServerNickname := d.GetOk("server_nickname")

	if !hasServerId && !hasServerNickname {
		return diag.Errorf("either 'server_id' or 'server_nickname' must be specified")
	}

	startVolumeId := d.Get("start_volume_id").(int)
	requestedCount := d.Get("requested_count").(int)

	// Build parameters for API call
	params := gwymodel.GetVolumeServerConnectionsParams{
		StartVolumeId: &startVolumeId,
		Count:         &requestedCount,
	}

	if hasServerId {
		id := serverId.(int)
		params.ServerId = &id
	}
	if hasServerNickname {
		name := serverNickname.(string)
		params.ServerNickname = &name
	}

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	connections, err := provObj.GetVolumeServerConnections(params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get connections: %v", err))
	}

	if err := d.Set("connections_info", convertMultipleVolumeServerConnectionsToSchema(connections)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connections_count", connections.Count); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total_count", connections.TotalCount); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

// ------------------- Resource -------------------

func ResourceAdminVolumeServerConnectionCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting attach (volume-server connection create)")

	serial := d.Get("serial").(int)
	params, err := buildAttachVolumeServerConnectionsParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = provObj.AttachVolumeToServers(*params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to attach volumes to servers: %v", err))
	}

	setResourceIDFromInputs(d)

	log.WriteInfo("attach completed successfully")
	return ResourceAdminVolumeServerConnectionRead(d)
}

func ResourceAdminVolumeServerConnectionRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	idStr := d.Id()
	if idStr == "" {
		return diag.Errorf("connection ID is empty")
	}

	connectionPairs := parseCompositeConnectionID(idStr)
	if len(connectionPairs) == 0 {
		return diag.Errorf("invalid connection ID format: %s", idStr)
	}

	reconObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	results, existing, err := reconObj.ReconcileReadVolumeServerConnections(connectionPairs)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(existing) == 0 {
		d.SetId("")
		return nil
	}

	if err := d.Set("connections_info", convertVolumeServerConnectionsListToSchema(results)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set connections_info: %w", err))
	}

	newID := buildCompositeConnectionID(results)
	d.SetId(newID)

	log.WriteInfo("Volume-server connections read successfully via reconciler")
	return nil
}

func ResourceAdminVolumeServerConnectionDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting volume-server connection delete via reconciler")

	serial := d.Get("serial").(int)
	idStr := d.Id()
	if idStr == "" {
		return diag.Errorf("connection ID is empty")
	}

	connectionPairs := parseCompositeConnectionID(idStr)
	if len(connectionPairs) == 0 {
		log.WriteWarn("No valid connections found in ID")
		d.SetId("")
		return nil
	}

	reconObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	err = reconObj.ReconcileDeleteVolumeServerConnections(connectionPairs)
	if err != nil {
		return diag.FromErr(fmt.Errorf("one or more detach operations failed: %v", err))
	}

	d.SetId("")
	log.WriteInfo("All volume-server connections deleted successfully via reconciler")
	return nil
}

func ResourceAdminVolumeServerConnectionUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting update for volume-server connections")

	serial := d.Get("serial").(int)
	reconObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse desired state from schema
	rawVolIDs, okVol := d.GetOk("volume_ids")
	rawSrvIDs, okSrv := d.GetOk("server_ids")
	if !okVol || !okSrv {
		return diag.Errorf("'volume_ids' and 'server_ids' must both be specified")
	}

	volumeIDs := make([]int, len(rawVolIDs.([]interface{})))
	for i, v := range rawVolIDs.([]interface{}) {
		volumeIDs[i] = v.(int)
	}
	serverIDs := make([]int, len(rawSrvIDs.([]interface{})))
	for i, v := range rawSrvIDs.([]interface{}) {
		serverIDs[i] = v.(int)
	}

	// Build desired pair list
	desiredPairs := make([]recmodel.VolumeServerPair, 0)
	for _, v := range volumeIDs {
		for _, s := range serverIDs {
			desiredPairs = append(desiredPairs, recmodel.VolumeServerPair{
				VolumeID: v,
				ServerID: s,
			})
		}
	}

	// Parse existing state from ID
	existingPairs := parseCompositeConnectionID(d.Id())

	// 🔹 Check if there are actual changes (ignore order)
	if !sameVolumeServerPairs(existingPairs, desiredPairs) {
		if err := reconObj.ReconcileUpdateVolumeServerConnections(existingPairs, desiredPairs); err != nil {
			return diag.FromErr(err)
		}
	}

	setResourceIDFromInputs(d)

	log.WriteInfo("Volume-server connection update completed successfully")
	return ResourceAdminVolumeServerConnectionRead(d)
}

// ------------------- Helpers -------------------

func buildAttachVolumeServerConnectionsParams(d *schema.ResourceData) (*gwymodel.AttachVolumeServerConnectionParam, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	rawVolIDs, okVol := d.GetOk("volume_ids")
	rawSrvIDs, okSrv := d.GetOk("server_ids")

	if !okVol || !okSrv {
		return nil, fmt.Errorf("'volume_ids' and 'server_ids' must both be specified")
	}

	// Convert volume_ids
	volListRaw := rawVolIDs.([]interface{})
	if len(volListRaw) == 0 {
		return nil, fmt.Errorf("'volume_ids' must contain at least one entry")
	}
	volumeIDs := make([]int, len(volListRaw))
	for i, v := range volListRaw {
		volumeIDs[i] = v.(int)
	}

	// Convert server_ids
	srvListRaw := rawSrvIDs.([]interface{})
	if len(srvListRaw) == 0 {
		return nil, fmt.Errorf("'server_ids' must contain at least one entry")
	}
	serverIDs := make([]int, len(srvListRaw))
	for i, v := range srvListRaw {
		serverIDs[i] = v.(int)
	}

	params := &gwymodel.AttachVolumeServerConnectionParam{
		VolumeIds: volumeIDs,
		ServerIds: serverIDs,
	}

	log.WriteDebug("Attach params: %+v", params)
	return params, nil
}

func convertVolumeServerConnectionToSchema(conn *gwymodel.VolumeServerConnectionDetail) map[string]interface{} {
	if conn == nil {
		return nil
	}

	// Convert LUN list
	luns := make([]map[string]interface{}, 0, len(conn.Luns))
	for _, l := range conn.Luns {
		luns = append(luns, map[string]interface{}{
			"lun":     l.Lun,
			"port_id": l.PortId,
		})
	}

	// Build Terraform-compatible map
	return map[string]interface{}{
		"id":        conn.Id,
		"volume_id": conn.VolumeId,
		"server_id": conn.ServerId,
		"luns":      luns,
	}
}

func convertVolumeServerConnectionsListToSchema(resp []gwymodel.VolumeServerConnectionDetail) []map[string]interface{} {
	// Sort for stable output (by ServerId then VolumeId)
	sort.Slice(resp, func(i, j int) bool {
		if resp[i].ServerId == resp[j].ServerId {
			return resp[i].VolumeId < resp[j].VolumeId
		}
		return resp[i].ServerId < resp[j].ServerId
	})

	// Convert each connection
	connections := make([]map[string]interface{}, len(resp))
	for i, conn := range resp {
		connections[i] = convertVolumeServerConnectionToSchema(&conn)
	}

	return connections
}

func convertMultipleVolumeServerConnectionsToSchema(resp *gwymodel.VolumeServerConnectionsResponse) []map[string]interface{} {
	if resp == nil || len(resp.Data) == 0 {
		return nil
	}

	// Sort for stable output (by ServerId then VolumeId)
	sort.Slice(resp.Data, func(i, j int) bool {
		if resp.Data[i].ServerId == resp.Data[j].ServerId {
			return resp.Data[i].VolumeId < resp.Data[j].VolumeId
		}
		return resp.Data[i].ServerId < resp.Data[j].ServerId
	})

	// Convert each connection
	connections := make([]map[string]interface{}, len(resp.Data))
	for i, conn := range resp.Data {
		connections[i] = convertVolumeServerConnectionToSchema(&conn)
	}

	return connections
}

// Build composite ID like "vol1,server1:vol2,server1:vol1,server2"
func buildCompositeConnectionID(connections []gwymodel.VolumeServerConnectionDetail) string {
	if len(connections) == 0 {
		return ""
	}

	ids := make([]string, len(connections))
	for i, c := range connections {
		ids[i] = c.Id
	}

	return strings.Join(ids, ":")
}

func parseCompositeConnectionID(idStr string) []recmodel.VolumeServerPair {
	if idStr == "" {
		return nil
	}

	parts := strings.Split(idStr, ":")
	pairs := make([]recmodel.VolumeServerPair, 0, len(parts))

	for _, part := range parts {
		ids := strings.Split(strings.TrimSpace(part), ",")
		if len(ids) != 2 {
			continue
		}

		volID, err1 := strconv.Atoi(strings.TrimSpace(ids[0]))
		srvID, err2 := strconv.Atoi(strings.TrimSpace(ids[1]))
		if err1 == nil && err2 == nil {
			pairs = append(pairs, recmodel.VolumeServerPair{VolumeID: volID, ServerID: srvID})
		}
	}

	return pairs
}

// sameVolumeServerPairs returns true if two slices contain the same pairs, regardless of order or duplicates.
func sameVolumeServerPairs(a, b []recmodel.VolumeServerPair) bool {
	if len(a) != len(b) {
		return false
	}

	setA := make(map[string]struct{}, len(a))
	for _, p := range a {
		key := fmt.Sprintf("%d-%d", p.VolumeID, p.ServerID)
		setA[key] = struct{}{}
	}

	for _, p := range b {
		key := fmt.Sprintf("%d-%d", p.VolumeID, p.ServerID)
		if _, exists := setA[key]; !exists {
			return false
		}
	}

	return true
}

func setResourceIDFromInputs(d *schema.ResourceData) {
	// Read the input lists from the schema
	volRaw, okVol := d.GetOk("volume_ids")
	servRaw, okServ := d.GetOk("server_ids")

	if !okVol || !okServ {
		// If either list is missing, don't set the ID
		return
	}

	// Convert []interface{} → []int
	volList := make([]int, 0)
	for _, v := range volRaw.([]interface{}) {
		volList = append(volList, v.(int))
	}

	servList := make([]int, 0)
	for _, s := range servRaw.([]interface{}) {
		servList = append(servList, s.(int))
	}

	// Sort to ensure deterministic order
	sort.Ints(volList)
	sort.Ints(servList)

	// Build vol,serv pairs like vol1,serv1:vol1,serv2:vol2,serv1:vol2,serv2
	pairs := make([]string, 0, len(volList)*len(servList))
	for _, volID := range volList {
		for _, servID := range servList {
			pairs = append(pairs, fmt.Sprintf("%d,%d", volID, servID))
		}
	}

	resourceID := strings.Join(pairs, ":")
	d.SetId(resourceID)
}
