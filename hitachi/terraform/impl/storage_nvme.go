package terraform

import (
	"fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	recmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- NVM Subsystem Datasources -------------------

func DatasourceVspNvmSubsystemRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// The spec requires an ID for a detailed "Get One" lookup
	subsystemID, ok := d.GetOk("nvm_subsystem_id")
	if !ok {
		return diag.Errorf("nvm_subsystem_id is required for this datasource")
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Call the specific "GetOne" reconciler method
	subsystem, err := reconObj.ReconcileGetNvmSubsystem(subsystemID.(int))
	if err != nil {
		log.WriteDebug("failed to get NVM subsystem %d: %v", subsystemID, err)
		return diag.FromErr(fmt.Errorf("failed to get NVM subsystem: %v", err))
	}

	if err := d.Set("nvm_subsystems", convertNvmSubsystemsToSchema([]gwymodel.NvmSubsystem{*subsystem})); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystems: %w", err))
	}

	if err := d.Set("nvm_subsystem_count", 1); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystem_count: %w", err))
	}

	// Set a unique ID for the datasource session
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func DatasourceVspNvmSubsystemsRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Map Schema inputs to Reconciler Input Model
	input := recmodel.NvmSubsystemGetMultipleInput{}

	// Basic Filters - Map using the new schema keys
	if v, ok := d.GetOk("nvm_subsystem_id"); ok {
		val := v.(int)
		input.NvmSubsystemId = &val
	}

	if v, ok := d.GetOk("nvm_subsystem_name"); ok {
		val := v.(string)
		input.NvmSubsystemName = &val
	}

	// Logic for the 'All' field:
	// If no specific filters are applied, we assume the user wants all subsystems.
	isAll := (input.NvmSubsystemId == nil && input.NvmSubsystemName == nil)
	input.All = &isAll

	// Inclusion Flags
	// Using d.Get() is correct here because these have Defaults in the schema
	input.IncludePorts = d.Get("include_ports").(bool)
	input.IncludeHostNqns = d.Get("include_host_nqns").(bool)
	input.IncludeNamespaces = d.Get("include_namespaces").(bool)

	// Call Reconciler
	resp, err := reconObj.ReconcileGetMultipleNvmSubsystems(input)
	if err != nil {
		log.WriteDebug("failed to get multiple NVM subsystems: %v", err)
		return diag.FromErr(fmt.Errorf("failed to get NVM subsystems: %v", err))
	}

	// Map Response data to Schema
	if err := d.Set("nvm_subsystems", convertNvmSubsystemsToSchema(resp.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystems: %w", err))
	}

	if err := d.Set("nvm_subsystem_count", len(resp.Data)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystem_count: %w", err))
	}

	// Set unique ID for the datasource session
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

// ------------------- NVM Subsystem Resource -------------------

// ------------------- Resource Read -------------------
func ResourceVspNvmSubsystemRead(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	stateID := d.Id()
	if stateID == "" {
		log.WriteInfo("No State ID found; resource not yet tracked.")
		return nil
	}

	id, err := strconv.Atoi(stateID)
	if err != nil {
		return diag.Errorf("failed to parse resource ID from state: %s", stateID)
	}

	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteInfo("Refreshing NVM Subsystem %d from storage", id)
	subsystem, err := reconObj.ReconcileGetNvmSubsystem(id)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to refresh NVM subsystem %d: %w", id, err))
	}
	if subsystem == nil {
		log.WriteInfo("Reconciler returned nil for NVM Subsystem %d. Removing from state.", id)
		d.SetId("")
		return nil
	}

	// Update input fields in state
	if err := setNvmSubsystemInputFields(d, *subsystem); err != nil {
		return diag.FromErr(err)
	}

	// Update computed list output
	if err := d.Set("nvm_subsystems", convertNvmSubsystemsToSchema([]gwymodel.NvmSubsystem{*subsystem})); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystems: %w", err))
	}

	return nil
}

// ------------------- Resource Apply / Create / Update -------------------
func ResourceVspNvmSubsystemApply(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	input, err := getNvmSubsystemReconcilerInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := reconObj.ReconcileNvmSubsystemApply(input)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("reconciliation returned nil result for serial %d", serial))
	}

	// Set resource ID (primary identity)
	d.SetId(strconv.Itoa(result.NvmSubsystemId))

	// Set input fields in state
	if err := setNvmSubsystemInputFields(d, *result); err != nil {
		return diag.FromErr(err)
	}

	// Set output list
	if err := d.Set("nvm_subsystems", convertNvmSubsystemsToSchema([]gwymodel.NvmSubsystem{*result})); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nvm_subsystems: %w", err))
	}

	return nil
}

func ResourceVspNvmSubsystemDelete(d *schema.ResourceData) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	reconObj, err := getReconcilerManagerSan(serial)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use the Resource ID from state
	id := d.Id()
	if id == "" {
		log.WriteWarn("Resource ID is empty during deletion")
		return diag.Errorf("invalid resource ID for deletion")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return diag.Errorf("failed to parse resource ID '%s' as integer: %v", id, err)
	}

	err = reconObj.ReconcileNvmSubsystemDelete(idInt)
	if err != nil {
		// Even on error, clear state to avoid dangling resource
		d.SetId("")
		log.WriteError("Failed to delete NVM subsystem %d: %v", idInt, err)
		return diag.FromErr(err)
	}

	// Clear state completely
	d.SetId("")
	_ = d.Set("nvm_subsystem_id", nil)
	_ = d.Set("virtual_nvm_subsystem_id", nil)
	_ = d.Set("nvm_subsystem_name", nil)
	_ = d.Set("host_mode", nil)
	_ = d.Set("host_mode_options", nil)
	_ = d.Set("enable_namespace_security", nil)
	_ = d.Set("ports", nil)
	_ = d.Set("host_nqns", nil)
	_ = d.Set("namespaces", nil)
	_ = d.Set("nvm_subsystems", nil)

	log.WriteInfo("Successfully deleted NVM subsystem %d", idInt)
	return nil
}

// ------------------- Helpers & Converters -------------------

func convertNvmSubsystemsToSchema(subsystems []gwymodel.NvmSubsystem) []map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	result := make([]map[string]interface{}, 0, len(subsystems))
	for _, s := range subsystems {
		m := map[string]interface{}{
			"nvm_subsystem_id":           s.NvmSubsystemId,
			"virtual_nvm_subsystem_id":   s.VirtualNvmSubsystemId,
			"nvm_subsystem_name":         s.NvmSubsystemName,
			"resource_group_id":          s.ResourceGroupId,
			"namespace_security_setting": s.NamespaceSecuritySetting,
			"t10_pi_mode":                s.T10piMode,
			"host_mode":                  s.HostMode,
			"host_mode_options":          s.HostModeOptions,
			"nvm_subsystem_nqn":          s.NvmSubsystemNqn,
			"port_ids":                   s.PortIds,
		}

		// --- Convert Host NQNs ---
		hostNqnList := make([]map[string]interface{}, 0, len(s.HostNqns))
		for _, hn := range s.HostNqns {
			hostNqnList = append(hostNqnList, map[string]interface{}{
				"host_nqn_id":       hn.HostNqnId,
				"host_nqn":          hn.HostNqn,
				"host_nqn_nickname": hn.HostNqnNickname,
			})
		}
		m["host_nqns"] = hostNqnList

		// --- Convert Namespaces ---
		nsList := make([]map[string]interface{}, 0, len(s.Namespaces))
		for _, ns := range s.Namespaces {
			nsMap := map[string]interface{}{
				"namespace_id":       ns.NamespaceId,
				"ldev_id":            ns.LdevId,
				"ldev_id_hex":        utils.IntToHexString(ns.LdevId),
				"namespace_nickname": ns.NamespaceNickname,
			}

			// Map paths (even if nil)
			if ns.Paths != nil {
				nsMap["paths"] = ns.Paths
			} else {
				nsMap["paths"] = []string{}
			}

			nsList = append(nsList, nsMap)
		}
		m["namespaces"] = nsList

		result = append(result, m)
	}
	return result
}

func setNvmSubsystemInputFields(d *schema.ResourceData, s gwymodel.NvmSubsystem) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Identity / Primary ID
	if err := d.Set("nvm_subsystem_id", s.NvmSubsystemId); err != nil {
		return fmt.Errorf("failed to set nvm_subsystem_id: %w", err)
	}

	// Configurable fields (input fields)
	if err := d.Set("virtual_nvm_subsystem_id", s.VirtualNvmSubsystemId); err != nil {
		return fmt.Errorf("failed to set virtual_nvm_subsystem_id: %w", err)
	}
	if err := d.Set("nvm_subsystem_name", s.NvmSubsystemName); err != nil {
		return fmt.Errorf("failed to set nvm_subsystem_name: %w", err)
	}
	if err := d.Set("host_mode", s.HostMode); err != nil {
		return fmt.Errorf("failed to set host_mode: %w", err)
	}
	if err := d.Set("host_mode_options", s.HostModeOptions); err != nil {
		return fmt.Errorf("failed to set host_mode_options: %w", err)
	}
	if err := d.Set("enable_namespace_security", s.NamespaceSecuritySetting == "Enable"); err != nil {
		return fmt.Errorf("failed to set enable_namespace_security: %w", err)
	}
	if err := d.Set("ports", s.PortIds); err != nil {
		return fmt.Errorf("failed to set ports: %w", err)
	}

	// Host NQNs
	hostNqns := make([]map[string]interface{}, 0, len(s.HostNqns))
	for _, hn := range s.HostNqns {
		hostNqns = append(hostNqns, map[string]interface{}{
			"nqn":      hn.HostNqn,
			"nickname": hn.HostNqnNickname,
		})
	}
	if err := d.Set("host_nqns", hostNqns); err != nil {
		return fmt.Errorf("failed to set host_nqns: %w", err)
	}

	// Namespaces
	namespaces := make([]map[string]interface{}, 0, len(s.Namespaces))
	for _, ns := range s.Namespaces {
		nsMap := map[string]interface{}{
			"ldev_id":     ns.LdevId,
			"ldev_id_hex": utils.IntToHexString(ns.LdevId),
			"nickname":    ns.NamespaceNickname,
		}
		if ns.Paths != nil {
			nsMap["path"] = ns.Paths
		} else {
			nsMap["path"] = []string{}
		}
		namespaces = append(namespaces, nsMap)
	}
	if err := d.Set("namespaces", namespaces); err != nil {
		return fmt.Errorf("failed to set namespaces: %w", err)
	}

	return nil
}

func getNvmSubsystemReconcilerInput(d *schema.ResourceData) (recmodel.NvmSubsystemReconcilerInput, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	input := recmodel.NvmSubsystemReconcilerInput{
		VirtualNvmSubsystemId: terrcommon.GetIntPointer(d, "virtual_nvm_subsystem_id"),
		NvmSubsystemName:      terrcommon.GetStringPointer(d, "nvm_subsystem_name"),
		HostMode:              terrcommon.GetStringPointer(d, "host_mode"),
	}

	// Check if the user actually provided an ID
	if v, ok := d.GetOkExists("nvm_subsystem_id"); ok {
		subId := v.(int)
		input.NvmSubsystemId = &subId
	} else {
		// If not provided, it stays nil.
		input.NvmSubsystemId = nil
	}

	// --- Boolean Handling (Avoid default false if missing) ---
	if _, ok := d.GetOkExists("enable_namespace_security"); ok {
		val := d.Get("enable_namespace_security").(bool)
		input.EnableNamespaceSecurity = &val
	}

	// --- Host Mode Options (TypeList) ---
	if v, ok := d.GetOkExists("host_mode_options"); ok {
		raw := v.([]interface{})
		options := make([]int, len(raw))
		for i, val := range raw {
			options[i] = val.(int)
		}
		input.HostModeOptions = &options
	}

	// --- Ports ---
	// Check if "ports" is actually in the HCL file
	if !d.GetRawConfig().GetAttr("ports").IsNull() {
		v, ok := d.GetOkExists("ports")
		if ok {
			var list []interface{}
			if s, ok := v.(*schema.Set); ok {
				list = s.List()
			}

			ports := make([]string, 0, len(list))
			for _, val := range list {
				ports = append(ports, val.(string))
			}
			input.Ports = &ports
		}
	}

	// --- Host NQNs ---
	if !d.GetRawConfig().GetAttr("host_nqns").IsNull() {
		v, ok := d.GetOkExists("host_nqns")
		if ok {
			var list []interface{}
			if s, ok := v.(*schema.Set); ok {
				list = s.List()
			}

			hostNqns := make([]recmodel.HostNqnInput, len(list))
			for i, item := range list {
				m := item.(map[string]interface{})
				hostNqns[i] = recmodel.HostNqnInput{
					Nqn:      m["nqn"].(string),
					Nickname: m["nickname"].(string),
				}
			}
			input.HostNqns = &hostNqns
		}
	}

	// --- Namespaces ---
	if !d.GetRawConfig().GetAttr("namespaces").IsNull() {
		v, ok := d.GetOkExists("namespaces")
		if ok {
			var list []interface{}
			if s, ok := v.(*schema.Set); ok {
				list = s.List()
			}

			namespaces := make([]recmodel.NamespaceInput, len(list))
			for i, item := range list {
				m := item.(map[string]interface{})

				// --- Support for ldev_id or ldev_id_hex ---
				finalLdev, err := terrcommon.ExtractLdevFromMap(m, "ldev_id", "ldev_id_hex")
				if err != nil {
					return input, err
				}

				// Ensure at least one was provided if the block exists
				var ldevID int
				if finalLdev != nil {
					ldevID = *finalLdev
				} else {
					return input, fmt.Errorf("either ldev_id or ldev_id_hex must be provided for namespace %d", i)
				}

				paths := []string{}
				if pVal, ok := m["path"]; ok && pVal != nil {
					if pSet, ok := pVal.(*schema.Set); ok {
						for _, pi := range pSet.List() {
							paths = append(paths, pi.(string))
						}
					}
				}

				namespaces[i] = recmodel.NamespaceInput{
					LdevId:   ldevID,
					Nickname: m["nickname"].(string),
					Path:     paths,
				}
			}
			input.Namespaces = &namespaces
		}
	}

	return input, nil
}

func mapNamespace(item map[string]interface{}) recmodel.NamespaceInput {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var paths []string
	if p, ok := item["path"]; ok && p != nil {
		rawPaths := p.([]interface{})
		paths = make([]string, len(rawPaths))
		for i, val := range rawPaths {
			paths[i] = val.(string)
		}
	}

	return recmodel.NamespaceInput{
		LdevId:   item["ldev_id"].(int),
		Nickname: item["nickname"].(string),
		Path:     paths,
	}
}
