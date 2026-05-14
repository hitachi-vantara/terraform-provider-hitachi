package sanstorage

import (
	"fmt"
	"strings"
	"sync"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

// --- Datasources ---

func (psm *sanStorageManager) ReconcileGetNvmSubsystem(subsystemID int) (*gwymodel.NvmSubsystem, error) {
	return psm.reconcileGetNvmSubsystemWithParams(subsystemID, true, true)
}

func (psm *sanStorageManager) reconcileGetNvmSubsystemWithParams(subsystemID int, includeHostNqns bool, includeNamespaces bool) (*gwymodel.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Fetching NVM subsystem: %d, includeHostNqns: %t, includeNamespaces: %t", subsystemID, includeHostNqns, includeNamespaces)

	// ports is automatically included
	input := reconmodel.NvmSubsystemGetMultipleInput{
		NvmSubsystemId:    &subsystemID,
		IncludeHostNqns:   includeHostNqns,
		IncludeNamespaces: includeNamespaces,
		All:               utils.Ptr(false),
	}

	resp, err := psm.ReconcileGetMultipleNvmSubsystems(input)
	if err != nil {
		log.WriteError("Failed to fetch NVM subsystem %d: %v", subsystemID, err)
		return nil, err
	}

	if resp == nil || len(resp.Data) == 0 {
		errStr := fmt.Errorf("NVM subsystem %d not found", subsystemID)
		log.WriteError(errStr.Error())
		return nil, errStr
	}

	return &resp.Data[0], nil
}

func (psm *sanStorageManager) ReconcileGetMultipleNvmSubsystems(input reconmodel.NvmSubsystemGetMultipleInput) (*gwymodel.NvmSubsystems, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Fetching single or multiple NVM subsystems with input: %+v", input)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	hasID := input.NvmSubsystemId != nil
	hasName := input.NvmSubsystemName != nil
	var baseSubsystems []gwymodel.NvmSubsystem

	// --- 1. IDENTITY RESOLUTION ---

	if hasID {
		// CASE: ID Only OR Both (ID + Name)
		log.WriteInfo("Fetching by ID: %d", *input.NvmSubsystemId)
		sub, err := provObj.GetNvmSubsystem(*input.NvmSubsystemId) // includes ports
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "404") || strings.Contains(strings.ToLower(err.Error()), "not found") {
				return &gwymodel.NvmSubsystems{Data: []gwymodel.NvmSubsystem{}}, nil
			}
			return nil, err
		}

		// Validate name if both were provided
		if hasName && !strings.EqualFold(sub.NvmSubsystemName, *input.NvmSubsystemName) {
			return nil, fmt.Errorf("identity mismatch: subsystem %d is named '%s' (expected '%s')",
				sub.NvmSubsystemId, sub.NvmSubsystemName, *input.NvmSubsystemName)
		}
		baseSubsystems = append(baseSubsystems, *sub)

		// check if there are namespaces
		input.IncludeNamespaces = false
		if len(baseSubsystems) > 0 && len(baseSubsystems[0].Namespaces) > 0 {
			input.IncludeNamespaces = true
		}
		input.IncludeNvmNqn = false // nqn included in get one

	} else if hasName {
		// CASE: Name Only
		log.WriteInfo("Searching for subsystem by name: %s", *input.NvmSubsystemName)
		resp, err := provObj.GetAllNvmSubsystems(gwymodel.GetNvmSubsystemsParams{})
		if err != nil {
			return nil, err
		}

		for _, sub := range resp.Data {
			if strings.EqualFold(sub.NvmSubsystemName, *input.NvmSubsystemName) {
				log.WriteDebug("Found nvme: %+v", sub)
				baseSubsystems = append(baseSubsystems, sub)
				break
			}
		}

		input.IncludeNvmNqn = true // nqn not in get all

	} else {
		// CASE: None (Both ID and Name are missing) -> Fetch All
		log.WriteInfo("No ID or Name provided. Fetching all NVM subsystems.")
		resp, err := provObj.GetAllNvmSubsystems(gwymodel.GetNvmSubsystemsParams{})
		if err != nil {
			return nil, err
		}
		baseSubsystems = resp.Data

		input.IncludeNvmNqn = true // nqn not in get all
	}

	// --- 2. ENRICHMENT ---

	if len(baseSubsystems) == 0 {
		return &gwymodel.NvmSubsystems{Data: baseSubsystems}, nil
	}

	// Map for concurrent enrichment
	subsystemMap := make(map[int]*gwymodel.NvmSubsystem)
	for i := range baseSubsystems {
		subsystemMap[baseSubsystems[i].NvmSubsystemId] = &baseSubsystems[i]
	}

	var tasks []func() error
	var mu sync.Mutex

	// Standard Tasks for NQNs and Ports
	if input.IncludeNvmNqn {
		tasks = append(tasks, psm.createEnrichmentTask(subsystemMap, "nqn", &mu))
	}

	if input.IncludePorts {
		tasks = append(tasks, psm.createEnrichmentTask(subsystemMap, "port", &mu))
	}

	if input.IncludeHostNqns {
		log.WriteDebug("Queueing Host NQN enrichment for %d subsystems", len(subsystemMap))
		for id := range subsystemMap {
			subID := id
			tasks = append(tasks, func() error {
				hostNqnData, err := provObj.GetAllHostNqns(subID)
				if err != nil {
					log.WriteError("Enrichment Task Error (HostNQN for ID %d): %v", subID, err)
					return nil
				}
				mu.Lock()
				if target, ok := subsystemMap[subID]; ok {
					target.HostNqns = hostNqnData.Data
				}
				mu.Unlock()
				return nil
			})
		}
	}

	// --- NAMESPACE & PATH ENRICHMENT (Per Subsystem ID) ---
	if input.IncludeNamespaces {
		for id := range subsystemMap {
			subID := id
			tasks = append(tasks, func() error {
				// A. Fetch Namespaces for this specific Subsystem ID
				nsResp, err := provObj.GetAllNamespaces(subID)
				if err != nil {
					log.WriteError("Failed to fetch namespaces for sub %d: %v", subID, err)
					return nil
				}

				// B. Fetch Paths for this specific Subsystem ID
				pathResp, err := provObj.GetNamespacePaths(gwymodel.GetNamespacePathsParams{
					NvmSubsystemId: subID,
				})
				if err != nil {
					log.WriteError("Failed to fetch paths for sub %d: %v", subID, err)
				}

				mu.Lock()
				defer mu.Unlock()

				target := subsystemMap[subID]
				target.Namespaces = nsResp.Data

				// C. Map Paths to the Namespaces using LdevId
				if pathResp != nil {
					// Create lookup map for this subsystem's namespaces
					nsMap := make(map[int]*gwymodel.Namespace)
					for i := range target.Namespaces {
						ns := &target.Namespaces[i]
						ns.Paths = []string{} // Initialize/Reset
						nsMap[ns.LdevId] = ns
					}

					for _, p := range pathResp.Data {
						// Match based on your API JSON: ldevId links Path to Namespace
						if ns, exists := nsMap[p.LdevId]; exists {
							ns.Paths = append(ns.Paths, p.HostNqn)
						}
					}
				}
				return nil
			})
		}
	}

	// --- 3. EXECUTE CONCURRENTLY ---
	if len(tasks) > 0 {
		_ = utils.RunConcurrentFuncs("NvmEnrichment", tasks)
		// Give the Hitachi Controller 500ms to release internal locks
		time.Sleep(500 * time.Millisecond)
	}

	return &gwymodel.NvmSubsystems{Data: baseSubsystems}, nil
}

// --- Resource Reconciler ---

func (psm *sanStorageManager) ReconcileNvmSubsystemApply(input reconmodel.NvmSubsystemReconcilerInput) (*gwymodel.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Starting reconciliation for NVM Subsystem with input: %+v", input)

	// Determine targetID: user-provided or auto-allocate
	targetID := -1
	if input.NvmSubsystemId != nil {
		targetID = *input.NvmSubsystemId
	}

	// Check existing subsystem
	var existing *gwymodel.NvmSubsystem
	if targetID != -1 {
		existing, _ = psm.ReconcileGetNvmSubsystem(targetID)
	}

	if existing != nil {
		log.WriteInfo("Subsystem %d exists. Validating updates...", targetID)

		if input.VirtualNvmSubsystemId != nil && existing.VirtualNvmSubsystemId != *input.VirtualNvmSubsystemId {
			return nil, fmt.Errorf("VirtualNvmSubsystemId is immutable: backend=%v input=%v",
				existing.VirtualNvmSubsystemId, *input.VirtualNvmSubsystemId)
		}

		time.Sleep(1 * time.Second)
		return psm.reconcileNvmSubsystemUpdate(targetID, input, existing)
	}

	// Create new subsystem
	if targetID == -1 {
		autoID, err := psm.getFirstAvailableNvmSubsystemId()
		if err != nil {
			return nil, err
		}
		targetID = autoID
	}

	return psm.reconcileNvmSubsystemCreate(targetID, input)
}

// --- Helpers ---

// createEnrichmentTask returns a function for concurrent execution that enriches
// base NVM subsystem objects with specific detail types (nqn, port, or namespace).
func (psm *sanStorageManager) createEnrichmentTask(subMap map[int]*gwymodel.NvmSubsystem, infoType string, mu *sync.Mutex) func() error {
	return func() error {
		log := commonlog.GetLogger()

		provObj, err := psm.getProvisionerManager()
		if err != nil {
			return err
		}

		// Fetch all subsystems with the specific enrichment info requested
		params := gwymodel.GetNvmSubsystemsParams{
			NvmSubsystemInfo: &infoType,
		}

		data, err := provObj.GetAllNvmSubsystems(params)
		if err != nil {
			log.WriteError("Enrichment task failed for type '%s': %v", infoType, err)
			return err
		}

		// Use the mutex to safely update the shared subsystem objects
		mu.Lock()
		defer mu.Unlock()

		for _, item := range data.Data {
			if target, ok := subMap[item.NvmSubsystemId]; ok {
				switch infoType {
				case "nqn":
					target.NvmSubsystemNqn = item.NvmSubsystemNqn
				case "port":
					target.PortIds = item.PortIds
				case "namespace":
					target.Namespaces = item.Namespaces
				}
			}
		}

		return nil
	}
}

func (psm *sanStorageManager) reconcileNvmSubsystemCreate(targetID int, input reconmodel.NvmSubsystemReconcilerInput) (*gwymodel.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Creating NVM Subsystem %d.", targetID)

	// Run Verification
	if err := psm.verifyNvmSubsystemCreateInput(input); err != nil {
		log.WriteError("Validation failed for subsystem creation: %v", err)
		return nil, err
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// Step 1: Create the base Subsystem
	var hmo *[]int
	if input.HostModeOptions != nil && len(*input.HostModeOptions) > 0 {
		hmo = input.HostModeOptions
	}

	request := gwymodel.CreateNvmSubsystemRequest{
		NvmSubsystemId:           targetID,
		VirtualNvmSubsystemId:    input.VirtualNvmSubsystemId,
		NvmSubsystemName:         input.NvmSubsystemName,
		HostMode:                 input.HostMode,
		HostModeOptions:          hmo,
		NamespaceSecuritySetting: mapSecurityInput(input.EnableNamespaceSecurity),
	}

	if _, err = provObj.CreateNvmSubsystem(request); err != nil {
		return nil, err
	}

	// Step 2: Orchestrate sub-resource creation
	// We pass targetID to all functions as it's the primary key
	if err := psm.createNvmSubsystemResources(targetID, input); err != nil {
		log.WriteError("Failed to populate sub-resources for subsystem %d: %v", targetID, err)
		return nil, err
	}

	// check if there are other inputs besides ports so we don't waste api calls
	includeHostNqns := false
	if input.HostNqns != nil {
		includeHostNqns = true
	}

	includeNamespaces := false
	if input.Namespaces != nil {
		includeNamespaces = true
	}

	return psm.reconcileGetNvmSubsystemWithParams(targetID, includeHostNqns, includeNamespaces)

}

func (psm *sanStorageManager) createNvmSubsystemResources(subId int, input reconmodel.NvmSubsystemReconcilerInput) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Add Ports ---
	if input.Ports != nil {
		for _, portID := range *input.Ports { // Dereference pointer to slice
			log.WriteInfo("Adding port %s to subsystem %d", portID, subId)
			req := gwymodel.AddNvmSubsystemPortRequest{
				NvmSubsystemId: subId,
				PortId:         portID,
			}
			if _, err := provObj.AddNvmSubsystemPort(req); err != nil {
				return err
			}
		}
	}

	// --- Register Host NQNs ---
	if input.HostNqns != nil {
		for _, host := range *input.HostNqns { // Dereference pointer to slice
			log.WriteInfo("Registering Host NQN %s to subsystem %d", host.Nqn, subId)
			req := gwymodel.RegisterHostNqnRequest{
				NvmSubsystemId: subId,
				HostNqn:        host.Nqn,
			}
			if _, err := provObj.RegisterHostNqn(req); err != nil {
				return err
			}

			if host.Nickname != "" {
				nickReq := gwymodel.SetHostNqnNicknameRequest{HostNqnNickname: host.Nickname}
				if _, err := provObj.SetHostNqnNickname(subId, host.Nqn, nickReq); err != nil {
					log.WriteWarn("Could not set nickname for Host NQN %s: %v", host.Nqn, err)
				}
			}
		}
	}

	// --- Create Namespaces ---
	nspacenum := 0 // auto namespace id
	if input.Namespaces != nil {
		for i, ns := range *input.Namespaces {
			// modify input.Namespaces, add namespaceId
			nspacenum = nspacenum + 1
			(*input.Namespaces)[i].NamespaceId = nspacenum

			log.WriteInfo("Creating Namespace with namespaceId %d for LDEV %d in subsystem %d", nspacenum, ns.LdevId, subId)
			req := gwymodel.CreateNamespaceRequest{
				NvmSubsystemId: subId,
				LdevId:         ns.LdevId,
				NamespaceId:    &nspacenum,
			}
			if _, err := provObj.CreateNamespace(req); err != nil {
				return err
			}

			if ns.Nickname != "" {
				nickReq := gwymodel.SetNamespaceNicknameRequest{NamespaceNickname: ns.Nickname}
				if _, err := provObj.SetNamespaceNickname(subId, nspacenum, nickReq); err != nil {
					log.WriteWarn("Could not set nickname for namespaceId %d LDEV %d: %v", nspacenum, ns.LdevId, err)
				}
			}
		}
	}

	// --- RECONCILE PATHS (MAPPINGS) ---
	// Note: We already checked input.Namespaces != nil above, but dereferencing
	// here ensures we iterate over the actual slice data.
	if input.Namespaces != nil {
		for _, ns := range *input.Namespaces {
			for _, hostNqn := range ns.Path {
				if hostNqn == "" {
					continue
				}

				// _, err := provObj.GetNamespacePathDetail(subId, hostNqn, ns.LdevId)
				// if err != nil {
				log.WriteInfo("Mapping NamespaceId %d LDEV %d to Host %s in NVM Subsystem %d", ns.NamespaceId, ns.LdevId, hostNqn, subId)

				pathReq := gwymodel.RegisterNamespacePathRequest{
					NvmSubsystemId: subId,
					HostNqn:        hostNqn,
					NamespaceId:    ns.NamespaceId,
				}

				if _, err := provObj.RegisterNamespacePath(pathReq); err != nil {
					log.WriteError("Failed to register path for LDEV %d and Host %s: %v", ns.LdevId, hostNqn, err)
					return err
				}
				// } else {
				// 	log.WriteDebug("Path already exists for LDEV %d and Host %s. Skipping.", ns.LdevId, hostNqn)
				// }
			}
		}
	}

	return nil
}

func (psm *sanStorageManager) reconcileNvmSubsystemUpdate(targetID int, input reconmodel.NvmSubsystemReconcilerInput, current *gwymodel.NvmSubsystem) (*gwymodel.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// Track if any modification occurs
	changed := false

	// --- TOP LEVEL ATTRIBUTE UPDATES ---

	if input.NvmSubsystemName != nil && *input.NvmSubsystemName != current.NvmSubsystemName {
		log.WriteInfo("Updating Name for ID %d", targetID)
		if _, err := provObj.UpdateNvmSubsystem(targetID, gwymodel.UpdateNvmSubsystemRequest{
			NvmSubsystemName: input.NvmSubsystemName,
		}); err != nil {
			return nil, err
		}
		changed = true
	}

	if input.HostMode != nil && *input.HostMode != current.HostMode {
		log.WriteInfo("Updating HostMode for ID %d", targetID)
		if _, err := provObj.UpdateNvmSubsystem(targetID, gwymodel.UpdateNvmSubsystemRequest{
			HostMode: input.HostMode,
		}); err != nil {
			return nil, err
		}
		changed = true
	}

	if input.HostModeOptions != nil && !utils.IntSliceEqual(*input.HostModeOptions, current.HostModeOptions) {
		log.WriteInfo("Updating Host Mode Options for ID %d: %v", targetID, *input.HostModeOptions)
		if _, err := provObj.UpdateNvmSubsystem(targetID, gwymodel.UpdateNvmSubsystemRequest{
			HostModeOptions: input.HostModeOptions,
		}); err != nil {
			return nil, err
		}
		changed = true
	}

	// Security
	desiredSec := mapSecurityInput(input.EnableNamespaceSecurity)
	if desiredSec != nil && *desiredSec != current.NamespaceSecuritySetting {
		log.WriteInfo("Updating Security for ID %d to %s", targetID, *desiredSec)
		if _, err := provObj.UpdateNvmSubsystem(targetID, gwymodel.UpdateNvmSubsystemRequest{
			NamespaceSecuritySetting: desiredSec,
		}); err != nil {
			return nil, err
		}
		changed = true
	}

	// --- SUB-RESOURCE RECONCILIATION ---
	// Note: updateNvmSubsystemResources should return a boolean 'subChanged'
	subChanged, err := psm.updateNvmSubsystemResources(targetID, input, current)
	if err != nil {
		log.WriteError("Failed to reconcile sub-resources for subsystem %d: %v", targetID, err)
		return nil, err
	}

	if subChanged {
		changed = true
	}

	// --- CONDITIONAL RE-FETCH ---
	if !changed {
		log.WriteInfo("No changes detected for subsystem %d. Skipping re-fetch.", targetID)
		return current, nil
	}

	log.WriteInfo("Changes applied to subsystem %d. Refreshing state from backend.", targetID)
	return psm.ReconcileGetNvmSubsystem(targetID)
}

func (psm *sanStorageManager) verifyNvmSubsystemCreateInput(input reconmodel.NvmSubsystemReconcilerInput) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	model, isSupported := checkIsHighEndSeries(psm.storageSetting.Serial)
	if !isSupported {
		return fmt.Errorf("NVM Subsystem creation is not supported on storage model: %s", model)
	}

	// Handle VirtualNvmSubsystemId logic
	if input.VirtualNvmSubsystemId != nil {
		_, isVsp5000 := checkIsVsp5000Series(psm.storageSetting.Serial)
		if isVsp5000 {
			log.WriteInfo("VirtualNvmSubsystemId %d specified for VSP 5000 series.", *input.VirtualNvmSubsystemId)
		} else {
			// IGNORE and LOG: Inform the user this is a 5000-series only feature
			log.WriteWarn("virtual_nvm_subsystem_id %d will be ignored. This attribute is only supported on VSP 5000 series (Current model: %s).",
				*input.VirtualNvmSubsystemId, model)
		}
	}

	return nil
}

func (psm *sanStorageManager) verifyNvmSubsystemUpdateInput(input reconmodel.NvmSubsystemReconcilerInput, current *gwymodel.NvmSubsystem) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	model, isSupported := checkIsHighEndSeries(psm.storageSetting.Serial)
	if !isSupported {
		return fmt.Errorf("NVM Subsystem updates are not supported on storage model: %s", model)
	}

	// Immutability Check: virtual_nvm_subsystem_id cannot be changed once created
	if input.VirtualNvmSubsystemId != nil {
		if *input.VirtualNvmSubsystemId != current.VirtualNvmSubsystemId {
			return fmt.Errorf("modification of virtual_nvm_subsystem_id is not supported and is only applicable to VSP 5000 series. Current: %d, Desired: %d",
				current.VirtualNvmSubsystemId, *input.VirtualNvmSubsystemId)
		}
	}

	return nil
}

func (psm *sanStorageManager) ReconcileNvmSubsystemDelete(id int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	log.WriteInfo("Deleting NVM Subsystem: %d", id)

	// ports are automatically included for get one
	includeHostNqns := true
	includeNamespaces := true

	current, err := psm.reconcileGetNvmSubsystemWithParams(id, includeHostNqns, includeNamespaces)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil // Already deleted
		}
		return err
	}

	// Clear Namespaces (and their Paths)
	log.WriteInfo("Cleaning up %d namespaces for subsystem %d", len(current.Namespaces), id)
	if err := psm.reconcileNvmNamespaces(id, []reconmodel.NamespaceInput{}, current.Namespaces); err != nil {
		return fmt.Errorf("failed to clean up namespaces during subsystem deletion: %w", err)
	}

	// Clear Host NQNs
	log.WriteInfo("Cleaning up %d Host NQNs for subsystem %d", len(current.HostNqns), id)
	if err := psm.reconcileNvmHostNqns(id, []reconmodel.HostNqnInput{}, current.HostNqns); err != nil {
		return fmt.Errorf("failed to clean up host NQNs during subsystem deletion: %w", err)
	}

	// Clear Ports
	log.WriteInfo("Cleaning up %d ports for subsystem %d", len(current.PortIds), id)
	if err := psm.reconcileNvmPorts(id, []string{}, current.PortIds); err != nil {
		return fmt.Errorf("failed to clean up ports during subsystem deletion: %w", err)
	}

	// Finally, Delete the Subsystem
	log.WriteInfo("Deleting NVM Subsystem: %d", id)
	_, err = provObj.DeleteNvmSubsystem(id)
	if err != nil {
		return fmt.Errorf("failed to delete NVM subsystem %d: %w", id, err)
	}

	return nil
}

func (psm *sanStorageManager) getFirstAvailableNvmSubsystemId() (int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return -1, err
	}

	log.WriteInfo("Fetching list of unused NVM Subsystem IDs from backend.")
	data, err := provObj.GetAllNvmSubsystems(gwymodel.GetNvmSubsystemsParams{NvmSubsystemOption: utils.Ptr("undefined")})
	if err != nil {
		log.WriteError("Failed to fetch available NVM subsystems: %v", err)
		return -1, err
	}

	// Check if the list of unused IDs is empty
	if len(data.Data) == 0 {
		err := fmt.Errorf("no available NVM Subsystem IDs found on the storage array")
		log.WriteError(err.Error())
		return -1, err
	}

	// Since the API returns unused IDs, we simply pick the first one in the list
	firstAvailable := data.Data[0].NvmSubsystemId
	log.WriteInfo("Selected first available NVM Subsystem ID: %d", firstAvailable)

	return firstAvailable, nil
}

func (psm *sanStorageManager) updateNvmSubsystemResources(subId int, input reconmodel.NvmSubsystemReconcilerInput, current *gwymodel.NvmSubsystem) (bool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	anyChanged := false

	// -----------------------
	// Ports
	// -----------------------
	if input.Ports != nil {
		desiredPorts := utils.RemoveDuplicateFromStringArray(*input.Ports)

		if !utils.StringSliceEqual(desiredPorts, current.PortIds) {
			log.WriteInfo("Reconciling Ports for NVM Subsystem %d (Desired: %v, Current: %v)",
				subId, desiredPorts, current.PortIds)

			if err := psm.reconcileNvmPorts(subId, desiredPorts, current.PortIds); err != nil {
				return false, err
			}
			anyChanged = true
		} else {
			log.WriteInfo("No Port changes detected for NVM Subsystem %d.", subId)
		}
	} else {
		log.WriteDebug("Ports attribute omitted in .tf; skipping.")
	}

	// -----------------------
	// Host NQNs
	// -----------------------
	if input.HostNqns != nil {
		desiredHostNqns := utils.DeduplicateByField(*input.HostNqns, func(n reconmodel.HostNqnInput) string {
			return n.Nqn
		})

		if psm.hasHostNqnChanges(desiredHostNqns, current.HostNqns) {
			log.WriteInfo("Reconciling Host NQNs for NVM Subsystem %d", subId)
			if err := psm.reconcileNvmHostNqns(subId, desiredHostNqns, current.HostNqns); err != nil {
				return false, err
			}
			anyChanged = true
		} else {
			log.WriteInfo("No Host NQN changes detected for NVM Subsystem %d.", subId)
		}
	} else {
		log.WriteDebug("HostNqns attribute omitted in .tf; skipping.")
	}

	// -----------------------
	// Namespaces
	// -----------------------
	if input.Namespaces != nil {
		desiredNamespaces := utils.DeduplicateByField(*input.Namespaces, func(ns reconmodel.NamespaceInput) int {
			return ns.LdevId
		})

		if psm.hasNamespaceChanges(desiredNamespaces, current.Namespaces) {
			log.WriteInfo("Reconciling Namespaces for NVM Subsystem %d", subId)
			if err := psm.reconcileNvmNamespaces(subId, desiredNamespaces, current.Namespaces); err != nil {
				return false, err
			}
			anyChanged = true
		} else {
			log.WriteInfo("No Namespace changes detected for NVM Subsystem %d.", subId)
		}
	} else {
		log.WriteDebug("Namespaces attribute omitted in .tf; skipping.")
	}

	return anyChanged, nil
}

func (psm *sanStorageManager) hasHostNqnChanges(desired []reconmodel.HostNqnInput, actual []gwymodel.HostNqnInfo) bool {
	if len(desired) != len(actual) {
		return true
	}
	// Create a map for quick lookup
	actualMap := make(map[string]gwymodel.HostNqnInfo)
	for _, a := range actual {
		actualMap[a.HostNqn] = a
	}

	for _, d := range desired {
		existing, found := actualMap[d.Nqn]
		if !found || (d.Nickname != existing.HostNqnNickname) {
			return true
		}
	}
	return false
}

func (psm *sanStorageManager) hasNamespaceChanges(desired []reconmodel.NamespaceInput, actual []gwymodel.Namespace) bool {
	// 1. If the count of namespaces differs, a change is definitely required
	if len(desired) != len(actual) {
		return true
	}

	// 2. Map actual namespaces by LdevId for O(1) lookup
	actualMap := make(map[int]gwymodel.Namespace)
	for _, a := range actual {
		actualMap[a.LdevId] = a
	}

	// 3. Compare each desired namespace against the existing one
	for _, d := range desired {
		existing, found := actualMap[d.LdevId]
		if !found {
			// LDEV ID is missing in the backend, needs creation
			return true
		}

		// Check Nickname (Note: using NamespaceNickname from your gateway struct)
		if d.Nickname != existing.NamespaceNickname {
			return true
		}

		// Check Paths/NQNs (Note: using existing.Paths from your gateway struct)
		// We use utils.StringSliceEqual to ignore order differences in NQN lists
		if !utils.StringSliceEqual(d.Path, existing.Paths) {
			return true
		}
	}

	return false
}

func (psm *sanStorageManager) reconcileNvmPorts(subId int, desired []string, actual []string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Categorize Ports for Logging ---
	desiredMap := make(map[string]bool)
	for _, p := range desired {
		desiredMap[p] = true
	}

	actualMap := make(map[string]bool)
	for _, p := range actual {
		actualMap[p] = true
	}

	var toAdd, toDelete, toRetain []string

	for p := range desiredMap {
		if actualMap[p] {
			toRetain = append(toRetain, p)
		} else {
			toAdd = append(toAdd, p)
		}
	}

	for p := range actualMap {
		if !desiredMap[p] {
			toDelete = append(toDelete, p)
		}
	}

	// Log the summary
	log.WriteInfo("Reconciling ports for subsystem %d: [Add: %v] [Delete: %v] [Retain: %v]",
		subId, toAdd, toDelete, toRetain)

	// --- Execution Logic ---

	// 1. Delete extra ports
	for _, p := range toDelete {
		log.WriteInfo("Removing port %s from subsystem %d", p, subId)
		if _, err := provObj.DeleteNvmSubsystemPort(subId, p); err != nil {
			return fmt.Errorf("failed to delete port %s: %w", p, err)
		}
	}

	// 2. Add missing ports
	for _, p := range toAdd {
		log.WriteInfo("Adding port %s to subsystem %d", p, subId)
		req := gwymodel.AddNvmSubsystemPortRequest{
			NvmSubsystemId: subId,
			PortId:         p,
		}
		if _, err := provObj.AddNvmSubsystemPort(req); err != nil {
			return fmt.Errorf("failed to add port %s: %w", p, err)
		}
	}

	return nil
}

func (psm *sanStorageManager) reconcileNvmHostNqns(subId int, desired []reconmodel.HostNqnInput, actual []gwymodel.HostNqnInfo) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Categorization ---
	actualMap := make(map[string]*gwymodel.HostNqnInfo)
	for i := range actual {
		actualMap[actual[i].HostNqn] = &actual[i]
	}

	var toAdd, toRemove, toUpdate, toRetain []string

	// Check for Add, Update, or Retain
	for _, d := range desired {
		existing, found := actualMap[d.Nqn]
		if !found {
			toAdd = append(toAdd, d.Nqn)
		} else if d.Nickname != "" && d.Nickname != existing.HostNqnNickname {
			toUpdate = append(toUpdate, d.Nqn)
		} else {
			toRetain = append(toRetain, d.Nqn)
		}
	}

	// Check for Remove
	desiredMap := make(map[string]bool)
	for _, d := range desired {
		desiredMap[d.Nqn] = true
	}
	for _, a := range actual {
		if !desiredMap[a.HostNqn] {
			toRemove = append(toRemove, a.HostNqn)
		}
	}

	// Summary Log
	log.WriteInfo("Reconciling Host NQNs for subsystem %d: [Add: %v] [Remove: %v] [Update Nickname: %v] [Retain: %v]",
		subId, toAdd, toRemove, toUpdate, toRetain)

	// --- Execution ---

	// 1. Remove extras
	for _, nqn := range toRemove {
		log.WriteInfo("Deleting Host NQN %s from subsystem %d", nqn, subId)
		if _, err := provObj.DeleteHostNqn(subId, nqn); err != nil {
			return fmt.Errorf("failed to delete host nqn %s: %w", nqn, err)
		}
	}

	// 2. Add new NQNs
	for _, nqn := range toAdd {
		log.WriteInfo("Registering Host NQN %s to subsystem %d", nqn, subId)
		req := gwymodel.RegisterHostNqnRequest{
			NvmSubsystemId: subId,
			HostNqn:        nqn,
		}
		if _, err := provObj.RegisterHostNqn(req); err != nil {
			return fmt.Errorf("failed to register host nqn %s: %w", nqn, err)
		}
	}

	// 3. Update Nicknames (for both newly added and existing if needed)
	// We iterate through 'desired' to ensure even newly added ones get their nickname set
	for _, d := range desired {
		existing, found := actualMap[d.Nqn]
		// Update if it's a known change (toUpdate) or if it's new (toAdd) and has a nickname
		isNewWithNickname := !found && d.Nickname != ""
		isExistingChange := found && d.Nickname != "" && d.Nickname != existing.HostNqnNickname

		if isNewWithNickname || isExistingChange {
			log.WriteInfo("Setting Nickname for Host NQN %s: %s", d.Nqn, d.Nickname)
			req := gwymodel.SetHostNqnNicknameRequest{HostNqnNickname: d.Nickname}
			if _, err := provObj.SetHostNqnNickname(subId, d.Nqn, req); err != nil {
				return fmt.Errorf("failed to set nickname for host nqn %s: %w", d.Nqn, err)
			}
		}
	}

	return nil
}

func (psm *sanStorageManager) reconcileNvmNamespaces(subId int, desired []reconmodel.NamespaceInput, actual []gwymodel.Namespace) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- 1. Categorization ---
	actualMap := make(map[int]gwymodel.Namespace)
	for _, a := range actual {
		actualMap[a.LdevId] = a
	}

	var toAdd, toDelete []int
	for _, d := range desired {
		if _, found := actualMap[d.LdevId]; !found {
			toAdd = append(toAdd, d.LdevId)
		}
	}

	desiredMap := make(map[int]bool)
	for _, d := range desired {
		desiredMap[d.LdevId] = true
	}
	for _, a := range actual {
		if !desiredMap[a.LdevId] {
			toDelete = append(toDelete, a.LdevId)
		}
	}

	log.WriteInfo("Reconciling Namespaces for subsystem %d: [Add: %v] [Delete: %v]", subId, toAdd, toDelete)

	// --- 2. Execution: Deletions (Paths First, then Namespace) ---
	for _, ldevID := range toDelete {
		existing := actualMap[ldevID]

		// A. REMOVE PATHS FIRST
		if len(existing.Paths) > 0 {
			log.WriteInfo("Clearing %d paths for Namespace %d (LDEV %d) before deletion", len(existing.Paths), existing.NamespaceId, ldevID)
			// Use an empty slice for desired paths to trigger a full unmap
			if err := psm.reconcileNamespacePaths(subId, existing.NamespaceId, []string{}, existing.Paths); err != nil {
				return fmt.Errorf("failed to clear paths for namespace %d before deletion: %w", existing.NamespaceId, err)
			}
		}

		// B. DELETE NAMESPACE
		log.WriteInfo("Deleting Namespace %d (LDEV %d) from subsystem %d", existing.NamespaceId, ldevID, subId)
		if _, err := provObj.DeleteNamespace(subId, existing.NamespaceId); err != nil {
			return fmt.Errorf("failed to delete namespace %d: %w", existing.NamespaceId, err)
		}
	}

	// --- 3. Execution: Creations ---
	for _, ldevID := range toAdd {
		log.WriteInfo("Creating new Namespace for LDEV %d in subsystem %d", ldevID, subId)
		if _, err := provObj.CreateNamespace(gwymodel.CreateNamespaceRequest{
			NvmSubsystemId: subId,
			LdevId:         ldevID,
		}); err != nil {
			return fmt.Errorf("failed to create namespace for LDEV %d: %w", ldevID, err)
		}
	}

	// --- 4. Execution: Updates & Path Management (For Retained/New) ---
	// If we changed the structure, refresh to get updated IDs/States
	var currentActual []gwymodel.Namespace
	if len(toAdd) > 0 || len(toDelete) > 0 {
		refresh, err := provObj.GetAllNamespaces(subId)
		if err != nil {
			return fmt.Errorf("failed to refresh namespaces for subsystem %d: %w", subId, err)
		}
		currentActual = refresh.Data
	} else {
		currentActual = actual
	}

	currentActualMap := make(map[int]gwymodel.Namespace)
	for _, a := range currentActual {
		currentActualMap[a.LdevId] = a
	}

	for _, nsInput := range desired {
		existing, found := currentActualMap[nsInput.LdevId]
		if !found {
			return fmt.Errorf("could not find Namespace for LDEV %d after synchronization", nsInput.LdevId)
		}

		// Nickname Update
		if nsInput.Nickname != "" && nsInput.Nickname != existing.NamespaceNickname {
			log.WriteInfo("Updating Nickname for Namespace %d (LDEV %d): %s", existing.NamespaceId, nsInput.LdevId, nsInput.Nickname)
			req := gwymodel.SetNamespaceNicknameRequest{NamespaceNickname: nsInput.Nickname}
			if _, err := provObj.SetNamespaceNickname(subId, existing.NamespaceId, req); err != nil {
				return fmt.Errorf("failed to set nickname for namespace %d: %w", existing.NamespaceId, err)
			}
		}

		// Path Reconcile
		uniquePaths := utils.RemoveDuplicateFromStringArray(nsInput.Path)
		if err := psm.reconcileNamespacePaths(subId, existing.NamespaceId, uniquePaths, existing.Paths); err != nil {
			return fmt.Errorf("failed to reconcile paths for namespace %d: %w", existing.NamespaceId, err)
		}
	}

	return nil
}

func (psm *sanStorageManager) reconcileNamespacePaths(subId int, namespaceId int, desiredPaths []string, currentPaths []string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// --- Categorization ---
	desiredMap := make(map[string]bool)
	for _, p := range desiredPaths {
		if p != "" {
			desiredMap[p] = true
		}
	}

	actualMap := make(map[string]bool)
	for _, p := range currentPaths {
		actualMap[p] = true
	}

	var toAdd, toDelete, toRetain []string

	for p := range desiredMap {
		if actualMap[p] {
			toRetain = append(toRetain, p)
		} else {
			toAdd = append(toAdd, p)
		}
	}

	for p := range actualMap {
		if !desiredMap[p] {
			toDelete = append(toDelete, p)
		}
	}

	// Summary Log - Critical for troubleshooting connectivity issues
	log.WriteInfo("Reconciling paths for Subsystem %d, Namespace %d: [Add: %v] [Delete: %v] [Retain: %v]",
		subId, namespaceId, toAdd, toDelete, toRetain)

	// --- Execution ---

	// 1. DELETE old paths (Unmapping)
	for _, hostNqn := range toDelete {
		log.WriteInfo("Unmapping Namespace %d from Host %s in subsystem %d", namespaceId, hostNqn, subId)
		if _, err := provObj.DeleteNamespacePath(subId, hostNqn, namespaceId); err != nil {
			log.WriteError("Failed to delete path for Namespace %d and Host %s: %v", namespaceId, hostNqn, err)
			return fmt.Errorf("failed to unmap path %s: %w", hostNqn, err)
		}
	}

	// 2. ADD new paths (Mapping)
	for _, hostNqn := range toAdd {
		log.WriteInfo("Mapping Namespace %d to Host %s in subsystem %d", namespaceId, hostNqn, subId)
		req := gwymodel.RegisterNamespacePathRequest{
			NvmSubsystemId: subId,
			HostNqn:        hostNqn,
			NamespaceId:    namespaceId,
		}
		if _, err := provObj.RegisterNamespacePath(req); err != nil {
			log.WriteError("Failed to register path for Namespace %d and Host %s: %v", namespaceId, hostNqn, err)
			return fmt.Errorf("failed to map path %s: %w", hostNqn, err)
		}
	}

	return nil
}

func mapSecurityInput(enable *bool) *string {
	if enable == nil {
		return nil
	}

	val := "Enable"
	if !*enable {
		val = "Disable"
	}

	return &val
}
