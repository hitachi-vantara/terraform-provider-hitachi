package terraform

import (
	"fmt"
	"strconv"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	provmanager "terraform-provider-hitachi/hitachi/storage/admin/provisioner"
	provimpl "terraform-provider-hitachi/hitachi/storage/admin/provisioner/impl"
	provmodel "terraform-provider-hitachi/hitachi/storage/admin/provisioner/model"
	reconcilerManager "terraform-provider-hitachi/hitachi/storage/admin/reconciler"
	reconcilerImpl "terraform-provider-hitachi/hitachi/storage/admin/reconciler/impl"
	reconcilerModel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Datasources -------------------

func DatasourceAdminServerHbasRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	v, ok := d.GetOk("server_id")
	if !ok {
		return diag.Errorf("server_id must be specified")
	}
	serverID := v.(int)

	// call provisioner directly
	provObj, err := getServerHbaProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverHBAs, err := provObj.GetServerHBAs(serverID)
	if err != nil {
		log.WriteDebug("failed to get server HBAs for server %v: %v", serverID, err)
		return diag.FromErr(fmt.Errorf("failed to get server HBAs for server %v: %v", serverID, err))
	}

	log.WriteDebug("server HBAs %+v", serverHBAs)
	if serverHBAs == nil {
		return nil
	}

	if err := d.Set("server_hba_info", convertMultipleServerHbasToSchema(serverHBAs)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set server_hba_info: %w", err))
	}
	if err := d.Set("server_hba_count", serverHBAs.Count); err != nil {
		return diag.FromErr(err)
	}

	// Set the resource ID so Terraform shows computed fields
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func DatasourceAdminServerHbaRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	v, ok := d.GetOk("server_id")
	if !ok {
		return diag.Errorf("server_id must be specified")
	}
	serverID := v.(int)

	v, ok = d.GetOk("initiator_name")
	if !ok {
		return diag.Errorf("initiator_name must be specified")
	}
	hbaWwn := v.(string)

	log.WriteDebug("Fetching server HBA for server ID %v and HBA WWN %v", serverID, hbaWwn)

	// call provisioner directly
	provObj, err := getServerHbaProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverHBA, err := provObj.GetServerHBAByWwn(serverID, hbaWwn)
	if err != nil {
		log.WriteDebug("failed to get server HBA for server %v, HBA WWN %v: %v", serverID, hbaWwn, err)
		return diag.FromErr(fmt.Errorf("failed to get server HBA for server %v, HBA WWN %v: %v", serverID, hbaWwn, err))
	}

	log.WriteInfo("server HBA info admin_server_hba %+v", serverHBA)
	if serverHBA == nil {
		return nil
	}

	if err := d.Set("server_hba_info", []map[string]interface{}{convertOneServerHbaToSchema(serverHBA)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set server_hba_info: %w", err))
	}

	// Set the resource ID so Terraform shows computed fields
	serverHbaStr := fmt.Sprintf("%d-%s", serverHBA.ServerID, serverHBA.HbaWwn)
	d.SetId(serverHbaStr)

	return nil
}

// ------------------- Helpers -------------------
func getServerHbaProvisionerManager(serial int) (provmanager.AdminStorageManager, error) {
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

func convertOneServerHbaToSchema(h *gwymodel.ServerHBA) map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Converting server HBA to schema: %+v", h)

	return map[string]interface{}{
		"server_id":  h.ServerID,
		"hba_wwn":    h.HbaWwn,
		"iscsi_name": h.IscsiName,
		"port_ids":   h.PortIds,
	}
}

func convertMultipleServerHbasToSchema(serverHBAs *gwymodel.ServerHBAList) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Defensive check
	if serverHBAs == nil || len(serverHBAs.Data) == 0 {
		return nil
	}

	hbaList := make([]map[string]interface{}, len(serverHBAs.Data))
	for i, h := range serverHBAs.Data {
		m := convertOneServerHbaToSchema(&h)
		hbaList[i] = m
	}

	return hbaList
}

// ------------------- Resource CRUD Operations -------------------

func ResourceAdminServerHbaCreate(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)

	// Build create parameters
	params, err := buildCreateAdminServerHbaParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	recObj, err := getReconcilerManagerForAdminServerHba(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverHBAList, err := recObj.ReconcileCreateAdminServerHBAs(serverID, params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID using one of the HBAs we're managing (from config), not from server response
	if len(params.HBAs) > 0 {
		firstConfigHBA := params.HBAs[0]
		var identifier string
		if firstConfigHBA.HbaWwn != "" {
			identifier = firstConfigHBA.HbaWwn
		} else {
			identifier = firstConfigHBA.IscsiName
		}
		d.SetId(fmt.Sprintf("%d-%d-%s", serial, serverID, identifier))
	} else {
		d.SetId(fmt.Sprintf("%d-%d", serial, serverID))
	}

	// Set output attributes
	if err := setServerHbaResourceAttributes(d, serverHBAList); err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("VSP One server HBAs created successfully for server ID: %d", serverID)
	return nil
}

func ResourceAdminServerHbaRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)

	// Get the current HBAs from the server to update computed fields
	provObj, err := getServerHbaProvisionerManager(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	serverHBAs, err := provObj.GetServerHBAs(serverID)
	if err != nil {
		log.WriteError("Failed to get current server HBAs for server %d: %v", serverID, err)
		// If we can't read the current state, don't fail but log the error
		// This ensures terraform destroy still works even if the server is gone
		log.WriteInfo("VSP One server HBA read completed with error (server may be deleted)")
		return nil
	}

	// Update the computed fields with current server state
	if err := setServerHbaResourceAttributes(d, serverHBAs); err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("VSP One server HBA read completed successfully")
	return nil
}

func ResourceAdminServerHbaDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	serverID := d.Get("server_id").(int)

	// Get HBAs from current state
	hbasRaw := d.Get("hbas").([]interface{})

	recObj, err := getReconcilerManagerForAdminServerHba(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Delete each HBA
	for _, hbaRaw := range hbasRaw {
		hbaMap := hbaRaw.(map[string]interface{})

		var initiatorName string
		if hbaWwn, ok := hbaMap["hba_wwn"].(string); ok && hbaWwn != "" {
			initiatorName = hbaWwn
		} else if iscsiName, ok := hbaMap["iscsi_name"].(string); ok && iscsiName != "" {
			initiatorName = iscsiName
		} else {
			log.WriteError("No valid initiator name found for HBA deletion")
			continue
		}

		_, err := recObj.ReconcileDeleteAdminServerHBA(serverID, initiatorName)
		if err != nil {
			log.WriteError("Failed to delete HBA %s: %v", initiatorName, err)
			// Continue with other HBAs even if one fails
		} else {
			log.WriteInfo("Successfully deleted HBA %s", initiatorName)
		}
	}

	// Clear the resource ID
	d.SetId("")
	log.WriteInfo("VSP One server HBAs deleted successfully for server ID: %d", serverID)
	return nil
}

// ------------------- Helpers -------------------

func getReconcilerManagerForAdminServerHba(serial int) (reconcilerManager.AdminStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilerModel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconcilerObj, err := reconcilerImpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get reconciler manager: %v", err)
		return nil, fmt.Errorf("failed to get reconciler manager: %w", err)
	}

	return reconcilerObj, nil
}

func buildCreateAdminServerHbaParams(d *schema.ResourceData) (gwymodel.CreateServerHBAParams, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	hbasRaw := d.Get("hbas").([]interface{})
	var hbas []gwymodel.ServerHBARequest

	for _, hbaRaw := range hbasRaw {
		hbaMap := hbaRaw.(map[string]interface{})

		hba := gwymodel.ServerHBARequest{}

		if hbaWwn, ok := hbaMap["hba_wwn"].(string); ok {
			hba.HbaWwn = hbaWwn
		}

		if iscsiName, ok := hbaMap["iscsi_name"].(string); ok {
			hba.IscsiName = iscsiName
		}

		// Validate that at least one of hba_wwn or iscsi_name is provided
		if hba.HbaWwn == "" && hba.IscsiName == "" {
			return gwymodel.CreateServerHBAParams{}, fmt.Errorf("either hba_wwn or iscsi_name must be provided for each HBA")
		}

		hbas = append(hbas, hba)
	}

	params := gwymodel.CreateServerHBAParams{
		HBAs: hbas,
	}

	log.WriteDebug("Built create params: %+v", params)
	return params, nil
}

func setServerHbaResourceAttributes(d *schema.ResourceData, serverHBAList *gwymodel.ServerHBAList) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Convert server HBA list to schema format
	hbaInfoList := convertMultipleServerHbasToSchema(serverHBAList)

	if err := d.Set("server_hba_info", hbaInfoList); err != nil {
		return fmt.Errorf("failed to set server_hba_info: %w", err)
	}

	if err := d.Set("server_hba_count", serverHBAList.Count); err != nil {
		return fmt.Errorf("failed to set server_hba_count: %w", err)
	}

	return nil
}
