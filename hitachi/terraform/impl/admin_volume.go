package terraform

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	provmanager "terraform-provider-hitachi/hitachi/storage/admin/provisioner"
	provimpl "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
	recmanager "terraform-provider-hitachi/hitachi/storage/admin/reconciler"
	recimpl "terraform-provider-hitachi/hitachi/storage/admin/reconciler/impl"
	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Datasources -------------------

func DatasourceAdminOneVolumeRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	v, ok := d.GetOkExists("volume_id")
	if !ok {
		return diag.Errorf("volume_id must be specified")
	}
	volumeID := v.(int)

	// call provisioner directly
	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeInfo, err := provObj.GetVolumeByID(volumeID)
	if err != nil {
		log.WriteDebug("failed to get volume %v: %v", volumeID, err)
		return diag.FromErr(fmt.Errorf("failed to get volume %v: %v", volumeID, err))
	}

	log.WriteDebug("volume %+v", volumeInfo)
	if volumeInfo == nil {
		return nil
	}

	if err := d.Set("volume_info", []map[string]interface{}{convertOneVolumeInfoToSchema(volumeInfo)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set volume_info: %w", err))
	}

	// Set the resource ID so Terraform shows computed fields
	d.SetId(fmt.Sprintf("%d", volumeInfo.ID))

	return nil
}

func DatasourceAdminMultipleVolumesRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	params := gwymodel.GetVolumeParams{}

	if v, ok := d.GetOkExists("pool_id"); ok {
		val := v.(int)
		params.PoolID = &val
	}
	if v, ok := d.GetOk("pool_name"); ok {
		val := v.(string)
		params.PoolName = &val
	}
	if v, ok := d.GetOkExists("server_id"); ok {
		val := v.(int)
		params.ServerID = &val
	}
	if v, ok := d.GetOk("server_nickname"); ok {
		val := v.(string)
		params.ServerNickname = &val
	}
	if v, ok := d.GetOk("nickname"); ok {
		val := v.(string)
		params.Nickname = &val
	}
	if v, ok := d.GetOkExists("min_total_capacity_mb"); ok {
		val := int64(v.(int))
		params.MinTotalCapacity = &val
	}
	if v, ok := d.GetOkExists("max_total_capacity_mb"); ok {
		val := int64(v.(int))
		params.MaxTotalCapacity = &val
	}
	if v, ok := d.GetOkExists("min_used_capacity_mb"); ok {
		val := int64(v.(int))
		params.MinUsedCapacity = &val
	}
	if v, ok := d.GetOkExists("max_used_capacity_mb"); ok {
		val := int64(v.(int))
		params.MaxUsedCapacity = &val
	}
	if v, ok := d.GetOkExists("start_volume_id"); ok {
		val := v.(int)
		params.StartVolumeID = &val
	}
	if v, ok := d.GetOkExists("requested_count"); ok {
		val := v.(int)
		params.Count = &val
	}

	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Params: %v", string(b))

	// call provisioner directly
	provObj, err := getProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	volumes, err := provObj.GetVolumes(params)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("volumes_info", convertMultipleVolumeInfosListToSchema(volumes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("volume_count", volumes.Count); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total_count", volumes.TotalCount); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// ------------------- Resource -------------------

func ResourceAdminVolumeRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	idStr := d.Id()
	if idStr == "" {
		return diag.Errorf("volume IDs in state are empty")
	}

	idParts := strings.Split(idStr, ",")
	var ids []int
	for _, part := range idParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idInt, err := strconv.Atoi(part)
		if err == nil {
			ids = append(ids, idInt)
		}
	}

	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	volInfos, existingIDs, err := recObj.ReconcileReadAdminVolumes(ids)
	if err != nil {
		if strings.Contains(err.Error(), "no existing volumes found") {
			d.SetId("") // Terraform will recreate them
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "All volumes missing",
				Detail:   err.Error(),
			}}
		}
		d.SetId("") // this will do terraform update, and create all missing vols
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	// Warn if desired number > existing volumes
	desiredNumber := len(ids) // default to number of IDs in state
	if len(existingIDs) < desiredNumber {
		diags = diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Some volumes are missing",
			Detail:   fmt.Sprintf("The number of existing volumes (%d) is less than the desired number (%d).", len(existingIDs), desiredNumber),
		}}
		d.SetId("") // this will do terraform update, and create missing vols
		// proceed to update volumes_info
	}

	// Normal update of computed fields
	if err := d.Set("volumes_info", convertVolumesInfoToSchema(volInfos)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set volume_info: %w", err))
	}
	if err := d.Set("volume_count", len(volInfos)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set volume_count: %w", err))
	}

	setResourceIDFromVolumeIDs(d, existingIDs)

	log.WriteInfo("volumes read successfully")
	return diags
}

func ResourceAdminVolumeDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	idStr := d.Id()
	if idStr == "" {
		return diag.Errorf("resource ID is empty")
	}

	// Parse IDs
	var ids []int
	for _, part := range strings.Split(idStr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idInt, err := strconv.Atoi(part)
		if err != nil {
			log.WriteWarn("invalid volume ID in state: %s", part)
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "Invalid volume ID",
				Detail:   fmt.Sprintf("Invalid volume ID in state: %s", part),
			}}
		}
		ids = append(ids, idInt)
	}

	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := recObj.ReconcileDeleteAdminVolumes(ids); err != nil {
		return diag.FromErr(err)
	}

	// Terraform state cleanup
	d.SetId("")
	log.WriteInfo("volumes deleted successfully")
	return nil
}

func ResourceAdminVolumeCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume create/update")

	serial := d.Get("serial").(int)

	params, _, err := buildCreateVolumeParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeIDs, err := recObj.ReconcileCreateAdminVolumes(params)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	setResourceIDFromVolumeIDs(d, volumeIDs)

	log.WriteInfo("volumes created successfully")
	return ResourceAdminVolumeRead(d)
}

func ResourceAdminVolumeUpdate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting volume update")

	serial := d.Get("serial").(int)

	// volume_id is required for update
	volumeID, ok := d.GetOk("volume_id")
	if !ok || volumeID.(int) <= 0 {
		return diag.Errorf("'volume_id' must be specified for update operations")
	}

	// Prevent illegal combinations
	if numVols, ok := d.GetOk("number_of_volumes"); ok && numVols.(int) > 0 {
		return diag.Errorf("'number_of_volumes' cannot be set during update — use 'volume_id' instead")
	}

	params, _, err := buildCreateVolumeParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update only one volume at a time
	err = recObj.ReconcileUpdateAdminVolume(volumeID.(int), params)
	if err != nil {
		return diag.FromErr(err)
	}

	// don't change resource ID to just the one volume that is updated, leave it as is
	// as other volumes may still exist in the resource
	log.WriteInfo(fmt.Sprintf("volume %d updated successfully", volumeID.(int)))
	return ResourceAdminVolumeRead(d)
}

// ------------------- Helpers -------------------
func getProvisionerManager(serial int) (provmanager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := provmodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	provObj, err := provimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get provisioner manager: %v", err)
		return nil, fmt.Errorf("failed to get provisioner manager: %w", err)
	}

	return provObj, nil
}

func getReconcilerManager(serial int) (recmanager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

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

	recObj, err := recimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get reconciler manager: %v", err)
		return nil, fmt.Errorf("failed to get reconciler manager: %w", err)
	}

	return recObj, nil
}

func convertOneVolumeInfoToSchema(v *gwymodel.VolumeInfoByID) map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	m := map[string]interface{}{
		"id":                           v.ID,
		"pool_id":                      v.PoolID,
		"total_capacity_mb":            v.TotalCapacity,
		"used_capacity_mb":             v.UsedCapacity,
		"free_capacity_mb":             v.FreeCapacity,
		"reserved_capacity_mb":         v.ReservedCapacity,
		"saving_setting":               v.SavingSetting,
		"capacity_saving_status":       v.CapacitySavingStatus,
		"number_of_connecting_servers": v.NumberOfConnectingServers,
		"number_of_snapshots":          v.NumberOfSnapshots,
		"volume_types":                 v.VolumeTypes,
	}

	// Optional fields
	if v.Nickname != nil {
		m["nickname"] = *v.Nickname
	}
	if v.PoolName != nil {
		m["pool_name"] = *v.PoolName
	}
	if v.IsDataReductionShareEnabled != nil {
		m["is_data_reduction_share_enabled"] = *v.IsDataReductionShareEnabled
	}
	if v.CompressionAcceleration != nil {
		m["compression_acceleration"] = *v.CompressionAcceleration
	}
	if v.CompressionAccelerationStatus != nil {
		m["compression_acceleration_status"] = *v.CompressionAccelerationStatus
	}
	if v.CapacitySavingProgress != nil {
		m["capacity_saving_progress"] = *v.CapacitySavingProgress
	}

	// Handle LUNs if present
	if len(v.LUNs) > 0 {
		luns := make([]map[string]interface{}, len(v.LUNs))
		for i, l := range v.LUNs {
			luns[i] = map[string]interface{}{
				"lun":       l.LUN,
				"server_id": l.ServerID,
				"port_id":   l.PortID,
			}
		}
		m["luns"] = luns
	}

	log.WriteDebug("Convert: %+v", m)
	return m
}

func convertVolumesInfoToSchema(volumes []gwymodel.VolumeInfoByID) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if len(volumes) == 0 {
		return nil
	}

	// Sort volumes by ID ascending
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].ID < volumes[j].ID
	})

	volumeList := make([]map[string]interface{}, 0, len(volumes))
	for i := range volumes {
		v := &volumes[i] // take address of actual element, not range copy
		m := convertOneVolumeInfoToSchema(v)
		volumeList = append(volumeList, m)
	}

	log.WriteDebug("Convert multiple: %+v", volumeList)
	return volumeList
}

func convertMultipleVolumeInfosListToSchema(volumes *gwymodel.VolumeInfoList) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Defensive check
	if volumes == nil || len(volumes.Data) == 0 {
		return nil
	}

	// Sort volumes by ID ascending
	sort.Slice(volumes.Data, func(i, j int) bool {
		return volumes.Data[i].ID < volumes.Data[j].ID
	})

	volumeList := make([]map[string]interface{}, len(volumes.Data))
	for i, v := range volumes.Data {
		m := map[string]interface{}{
			"id":                              v.ID,
			"pool_id":                         v.PoolID,
			"total_capacity_mb":               v.TotalCapacity,
			"used_capacity_mb":                v.UsedCapacity,
			"saving_setting":                  v.SavingSetting,
			"capacity_saving_status":          v.CapacitySavingStatus,
			"number_of_connecting_servers":    v.NumberOfConnectingServers,
			"number_of_snapshots":             v.NumberOfSnapshots,
			"volume_types":                    v.VolumeTypes,
			"is_data_reduction_share_enabled": false,
			"compression_acceleration":        false,
		}

		// Optional fields
		if v.Nickname != nil {
			m["nickname"] = *v.Nickname
		}
		if v.PoolName != nil {
			m["pool_name"] = *v.PoolName
		}
		if v.IsDataReductionShareEnabled != nil {
			m["is_data_reduction_share_enabled"] = *v.IsDataReductionShareEnabled
		}
		if v.CompressionAcceleration != nil {
			m["compression_acceleration"] = *v.CompressionAcceleration
		}

		volumeList[i] = m
	}

	return volumeList
}

func buildCreateVolumeParams(d *schema.ResourceData) (gwymodel.CreateVolumeParams, int, error) {
	var params gwymodel.CreateVolumeParams

	// Required: pool_id
	params.PoolID = d.Get("pool_id").(int)

	// Required: capacity
	capacityStr, ok := d.GetOk("capacity")
	if !ok {
		return params, 0, fmt.Errorf("capacity must be specified")
	}
	miB, err := utils.ParseCapacityToMiB(capacityStr.(string))
	if err != nil {
		return params, 0, fmt.Errorf("invalid capacity: %v", err)
	}
	params.Capacity = miB

	// Handle optional number_of_volumes
	desiredNum := 1 // Default to 1 if not provided
	if v, ok := d.GetOk("number_of_volumes"); ok {
		val := v.(int)
		if val > 0 {
			params.Number = &val
			desiredNum = val
		}
	} else {
		// Not set — ensure params.Number reflects the default
		params.Number = &desiredNum
	}

	// Optional: capacity_saving
	if v, ok := d.GetOk("capacity_saving"); ok {
		val := v.(string)
		params.SavingSetting = &val
	}

	// Optional: is_data_reduction_share_enabled
	if v, ok := d.GetOk("is_data_reduction_share_enabled"); ok {
		val := v.(bool)
		params.IsDataReductionShareEnabled = &val
	}

	// Nickname parameter (required)
	nicknameParam, ok := d.GetOk("nickname_param")
	if !ok || len(nicknameParam.([]interface{})) == 0 {
		return params, 0, fmt.Errorf("'nickname_param' must be specified with at least base_name")
	}

	m := nicknameParam.([]interface{})[0].(map[string]interface{})
	baseName, ok := m["base_name"].(string)
	if !ok || baseName == "" {
		return params, 0, fmt.Errorf("'base_name' in nickname_param is required")
	}

	nickname := gwymodel.VolumeNicknameParam{BaseName: baseName}

	if val, exists := m["start_number"]; exists {
		n := val.(int)
		nickname.StartNumber = &n
	}
	if val, exists := m["number_of_digits"]; exists {
		n := val.(int)
		nickname.NumberOfDigits = &n
	}

	// Validation: number_of_digits requires start_number
	if nickname.StartNumber == nil && nickname.NumberOfDigits != nil {
		return params, 0, fmt.Errorf("'number_of_digits' is specified but 'start_number' is not — please specify start_number when using number_of_digits")
	}

	params.NicknameParam = nickname

	return params, desiredNum, nil
}

// converts a slice of int IDs to a comma-separated string and sets it as the Terraform resource ID.
func setResourceIDFromVolumeIDs(d *schema.ResourceData, volumeIDs []int) {
	sort.Ints(volumeIDs) // ✅ Ensure consistent order
	idStrs := make([]string, len(volumeIDs))
	for i, id := range volumeIDs {
		idStrs[i] = strconv.Itoa(id)
	}
	d.SetId(strings.Join(idStrs, ","))
}
