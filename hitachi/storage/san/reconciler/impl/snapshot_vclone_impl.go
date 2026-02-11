package sanstorage

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	provmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

var (
	CachePool           *provmodel.DynamicPool
	CachePvol           *gwymodel.LogicalUnit
	CacheSvol           *gwymodel.LogicalUnit
	snapshotGlobalCache *reconmodel.ReconcileSnapshotResult
	cacheLock           sync.Mutex
)

// --- Datasources ---

func (psm *sanStorageManager) ReconcileGetFamily(ldevID int) ([]gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Get the provisioner manager to access the hardware API
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// 2. Call the provisioner function you provided
	// This calls GET base-URL/v1/objects/snapshot-families
	log.WriteInfo("TIA| Fetching snapshot family for LDEV %d", ldevID)
	resp, err := provObj.GetSnapshotFamily(ldevID)
	if err != nil {
		return nil, log.WriteAndReturnError("failed to get snapshot family for LDEV %d: %v", ldevID, err)
	}

	// 3. Check if data exists in the response
	if resp == nil || len(resp.Data) == 0 {
		log.WriteInfo("TIA| No snapshot family members found for LDEV %d", ldevID)
		return []gwymodel.SnapshotFamily{}, nil
	}

	// 4. Return the array of family members
	// Your SnapshotFamilyListResponse already contains the Data slice
	// that matches your SnapshotFamily struct
	return resp.Data, nil
}

func (psm *sanStorageManager) ReconcileGetVirtualCloneParentVolumes() ([]int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Access the provisioner manager
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// 2. Call the hardware API wrapper
	// This maps to: GET base-URL/v1/objects/virtual-clone-parent-volumes
	log.WriteInfo("TIA| Fetching all Virtual Clone Parent LDEVs")
	resp, err := provObj.GetVirtualCloneParentVolumes()
	if err != nil {
		return nil, log.WriteAndReturnError("failed to get virtual clone parent volumes: %v", err)
	}

	// 3. Handle empty responses
	if resp == nil || len(resp.Data) == 0 {
		log.WriteInfo("TIA| No Virtual Clone Parent volumes found on this storage.")
		return []int{}, nil
	}

	// 4. Extract LdevIDs into an array of integers
	ldevIds := make([]int, 0, len(resp.Data))
	for _, item := range resp.Data {
		ldevIds = append(ldevIds, item.LdevID)
	}

	log.WriteInfo("TIA| Found %d parent volumes: %v", len(ldevIds), ldevIds)
	return ldevIds, nil
}

// --- Resource ---

// New main entry for Terraform Create/Update/Delete
func (psm *sanStorageManager) ReconcileSnapshotVclone(input reconmodel.SnapshotReconcilerInput) (*reconmodel.ReconcileSnapshotResult, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.Action == nil || *input.Action == "" {
		input.Action = utils.Ptr("read")
	}

	allresult := reconmodel.ReconcileSnapshotResult{}

	if *input.Action == "read" {
		return psm.ReconcileReadExistingSnapshotVclone(input)
	}

	if *input.Action != "vclone" && *input.Action != "vrestore" {
		// for plain snapshot
		snapshot, err := psm.ReconcileSnapshot(input)
		if err != nil {
			return nil, err
		}
		allresult.Snapshot = snapshot

		vcloneFamily, _ := psm.checkThenGetVclone(input, snapshot)
		allresult.VcloneFamily = vcloneFamily
	} else {
		// for vclone
		snapshot, vcloneFamily, err := psm.ReconcileVclone(input)
		if err != nil {
			return nil, err
		}
		allresult.Snapshot = snapshot
		allresult.VcloneFamily = vcloneFamily
	}

	log.WriteDebug("Snapshot: %+v", allresult.Snapshot)

	// Add Universal Metadata (HDP/HTI, DRS, etc.)
	svol := input.SvolLdevID
	if svol == nil || *svol == 0 {
		if allresult.Snapshot != nil && allresult.Snapshot.SvolLdevID != 0 {
			svol = &allresult.Snapshot.SvolLdevID
			log.WriteDebug("Snapshot svol=%v", *svol)
		} else if allresult.VcloneFamily != nil {
			svol = &allresult.VcloneFamily.LdevID
		}
	}
	if svol != nil {
		log.WriteDebug("Before getUniversalInfo: svol=%v", *svol)
	}
	uni, _ := psm.getUniversalInfo(input.PvolLdevID, svol, true)
	allresult.UniversalInfo = uni

	if allresult.Snapshot != nil {
		log.WriteDebug("ALLresult.Snapshot: %+v", *allresult.Snapshot)
	}
	if allresult.VcloneFamily != nil {
		log.WriteDebug("ALLresult.Vfamily: %+v", *allresult.VcloneFamily)
	}
	if allresult.UniversalInfo != nil {
		log.WriteDebug("ALLresult.UniversalInfo: %+v", *allresult.UniversalInfo)
	}

	return &allresult, nil
}

func (psm *sanStorageManager) ReconcileVclone(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	action := strings.ToLower(*input.Action)

	log.WriteInfo("Starting snapshot reconciliation for operation: %s", action)

	// pvol and mu are set in impl code
	log.WriteDebug("Input PvolLdevID: %v, MuNumber: %v", input.PvolLdevID, input.MuNumber)

	switch action {
	case "vclone":
		return psm.reconcileSnapshotVCloneCreateConvert(input)
	case "vrestore":
		return psm.reconcileSnapshotVCloneRestore(input)
	default:
		return nil, nil, log.WriteAndReturnError("unsupported snapshot operation: %s", *input.Action)
	}
}

// Terraform Read
func (psm *sanStorageManager) ReconcileReadExistingSnapshotVclone(input reconmodel.SnapshotReconcilerInput) (*reconmodel.ReconcileSnapshotResult, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	allresult := reconmodel.ReconcileSnapshotResult{}

	// 1. Attempt to get the traditional Snapshot pair (Uses MU Number)
	if input.MuNumber != nil {
		snap, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, true)
		if err != nil {
			if !isNotFoundError(err) {
				return nil, log.WriteAndReturnError("TFError| System error fetching snapshot for P-VOL %d: %v", *input.PvolLdevID, err)
			}
			log.WriteInfo("TFInfo| Snapshot MU %d not found for P-VOL %d. Checking for vClone...", *input.MuNumber, *input.PvolLdevID)
		} else {
			// Snapshot found - validate it matches our config (e.g. correct S-VOL)
			if vErr := psm.validateExistingSnapshotBackendFromInput(input, snap); vErr != nil {
				return nil, vErr
			}
			allresult.Snapshot = snap
			log.WriteInfo("TFInfo| Traditional snapshot pair found (P-VOL: %d, MU: %d)", *input.PvolLdevID, *input.MuNumber)
		}
	}

	vcloneFamily, _ := psm.checkThenGetVclone(input, allresult.Snapshot)
	allresult.VcloneFamily = vcloneFamily

	svol := input.SvolLdevID
	if svol == nil {
		if allresult.Snapshot != nil && allresult.Snapshot.SvolLdevID != 0 {
			svol = &allresult.Snapshot.SvolLdevID
		} else if allresult.VcloneFamily != nil {
			svol = &allresult.VcloneFamily.LdevID
		}
	}
	uni, _ := psm.getUniversalInfo(input.PvolLdevID, svol, false)
	allresult.UniversalInfo = uni

	// 3. Evaluation: If neither exists, the resource is truly gone
	if allresult.Snapshot == nil && allresult.VcloneFamily == nil {
		log.WriteInfo("TFInfo| No relationship found (Snapshot or vClone) for P-VOL %d. Proceeding as 'Not Found'.", *input.PvolLdevID)
		return &reconmodel.ReconcileSnapshotResult{}, nil
	}

	return &allresult, nil
}

// --- VClone Action Handlers ---

func (psm *sanStorageManager) reconcileSnapshotVCloneCreateConvert(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Validate Input and Hardware State
	// Returns (*Snapshot, *SnapshotFamily, error)
	snapshot, family, err := psm.validateTIAVcloneInput(input)
	if err != nil {
		return nil, nil, err
	}

	// IDEMPOTENCY CHECK: If already a vClone, return immediately
	if snapshot == nil && family != nil {
		log.WriteInfo("Idempotency| P-VOL %d (MU %d) is already a vClone. Skipping create/convert.",
			*input.PvolLdevID, *input.MuNumber)
		return nil, family, nil
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, err
	}

	pvolID := *input.PvolLdevID
	svolID := snapshot.SvolLdevID
	muNum := *input.MuNumber

	// 5. Action Logic based on current Snapshot Status
	switch snapshot.Status {
	case "PAIR":
		if snapshot.IsConsistencyGroup {
			return nil, nil, log.WriteAndReturnError("vClone 'create' is forbidden: P-VOL %d (MU %d) belongs to a Consistency Group.", pvolID, muNum)
		}
		log.WriteInfo("Status is PAIR. Performing vClone 'create' for P-VOL %d, MU %d", pvolID, muNum)
		_, err = provObj.CreateSnapshotVClone(pvolID, muNum)

	case "PSUS", "PFUS":
		log.WriteInfo("Status is %s. Performing vClone 'convert' for P-VOL %d, MU %d", snapshot.Status, pvolID, muNum)
		_, err = provObj.ConvertSnapshotVClone(pvolID, muNum)

	default:
		return nil, nil, log.WriteAndReturnError("Invalid status %s for vClone transition", snapshot.Status)
	}

	if err != nil {
		return nil, nil, err
	}

	// 6. It will return (nil, familyMember, nil) on success
	return psm.waitForSnapshotVCloneStatus(&pvolID, &svolID, &muNum)
}

func (psm *sanStorageManager) reconcileSnapshotVCloneRestore(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Initial State Check
	snapshot, family, err := psm.validateTIAVrestoreInput(input)
	if err != nil {
		return nil, nil, err
	}

	pvolID := *input.PvolLdevID
	muNum := *input.MuNumber

	// --- REPAIR / PAIR BACK LOGIC ---
	// If the snapshot record is missing but the family confirms a vClone exists
	if snapshot == nil && family != nil {
		log.WriteInfo("Restore| S-VOL %d is a standalone vClone. Pairing back to P-VOL %d...", family.SvolLdevID, pvolID)

		// This is mandatory for volumes with Compression/Deduplication enabled.
		forceCopy := true
		input.IsDataReductionForceCopy = &forceCopy

		// Now call the create logic to re-establish the relationship
		repairSnap, repairErr := psm.reconcileSnapshotCreate(input)
		if repairErr != nil {
			return nil, nil, log.WriteAndReturnError("Restore| Failed to pair back: %v", repairErr)
		}

		snapshot = repairSnap
		log.WriteInfo("Restore| Paired back successfully. Current status: %s", snapshot.Status)
	}

	if snapshot != nil && !isTIAdvancedPair(snapshot) {
		return nil, nil, log.WriteAndReturnError("vRestore requires TIA pairs (P-VOL %d, MU %d)", *input.PvolLdevID, *input.MuNumber)
	}

	var svolID int
	if snapshot != nil {
		svolID = snapshot.SvolLdevID
	} else if family != nil {
		svolID = family.SvolLdevID
	}

	// --- IDEMPOTENCY CHECK ---
	// If it's already in the final state (PSUS and Flags cleared), we are done
	if snapshot != nil && snapshot.Status == "PSUS" && !snapshot.IsVirtualCloneVolume {
		log.WriteInfo("Idempotency| Snapshot is already in PSUS (restored state). Skipping execution.")
		return snapshot, family, nil
	}

	// --- EXECUTION PHASE ---
	log.WriteInfo("Restore| Executing vRestore command for P-VOL %d, MU %d", pvolID, muNum)
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, err
	}

	_, err = provObj.RestoreSnapshotFromVClone(pvolID, muNum)
	if err != nil {
		return nil, nil, log.WriteAndReturnError("vRestore command failed: %v", err)
	}

	// --- FINAL WAIT ---
	// Polls for Status: PSUS, isVirtualCloneVolume: false
	return psm.waitForSnapshotVRestoreStatus(&pvolID, &svolID, &muNum)
}

func (psm *sanStorageManager) reconcileSnapshotTreeDelete(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.PvolLdevID == nil || input.MuNumber == nil {
		return nil, fmt.Errorf("PvolLdevID and Mirror Unit number required to delete a snapshot tree")
	}

	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		return nil, err
	}

	// Block TIA
	if isTIAdvancedPair(snapshot) {
		return nil, log.WriteAndReturnError("Snapshot Tree delete/cleanup is not supported for Thin Image Advanced (TIA).")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.DeleteSnapshotTreeRequest{
		Parameters: gwymodel.DeleteSnapshotTreeParams{
			LdevID: *input.PvolLdevID,
		},
	}

	log.WriteInfo("Deleting entire snapshot tree starting from Root LDEV: %d", *input.PvolLdevID)
	_, err = provObj.DeleteSnapshotTree(request)
	if err != nil {
		return nil, err
	}

	return nil, psm.waitForSnapshotDeletion(input.PvolLdevID, input.MuNumber)
}

func (psm *sanStorageManager) reconcileGarbageData(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.PvolLdevID == nil {
		return nil, fmt.Errorf("PvolLdevID is required for garbage data operations")
	}

	model, isVsp5000 := checkIsVsp5000Series(psm.storageSetting.Serial)
	if !isVsp5000 {
		errmsg := fmt.Sprintf("Garbage data deletion is only supported on VSP 5000 series (TI Standard). "+
			"Operation not supported on current model: %s", model)
		return nil, log.WriteAndReturnError("%v", errmsg)
	}

	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		return nil, err
	}

	// Block TIA
	if isTIAdvancedPair(snapshot) {
		return nil, log.WriteAndReturnError("Garbage data defrag: not supported for Thin Image Advanced (TIA) snapshots.")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	if input.DefragOperation == nil || *input.DefragOperation == "" {
		return nil, log.WriteAndReturnError(
			"DefragOperation is required for garbage data cleanup. " +
				"Please specify 'start' or 'stop' in the defrag_operation field.")
	}

	opType := strings.ToLower(*input.DefragOperation)

	request := gwymodel.DeleteGarbageDataRequest{
		Parameters: gwymodel.DeleteGarbageDataParams{
			LdevID:        *input.PvolLdevID,
			OperationType: opType,
		},
	}

	log.WriteInfo("Executing Garbage Data cleanup (%s) for Root LDEV: %d", opType, *input.PvolLdevID)
	_, err = provObj.DeleteGarbageData(request)
	if err != nil {
		return nil, fmt.Errorf("garbage data operation (%s) failed for LDEV %d: %w", opType, *input.PvolLdevID, err)
	}

	// get latest
	snapshot, err = psm.getSnapshot(input.PvolLdevID, input.MuNumber, true)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

// --- Validations ---

func (psm *sanStorageManager) validateTIAVcloneInput(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pvolID := *input.PvolLdevID
	muNum := *input.MuNumber

	// Hardware Model Check
	model, isBSeries := checkIsVspBSeries(psm.storageSetting.Serial)
	if !isBSeries {
		return nil, nil, log.WriteAndReturnError("vClone operations only supported on VSP B series. Model: %s", model)
	}

	// 1. Primary Lookup: Try to find a traditional Thin Image Snapshot Pair
	snapshot, err := psm.getSnapshot(&pvolID, &muNum, false)

	if snapshot != nil && !isTIAdvancedPair(snapshot) {
		return nil, nil, log.WriteAndReturnError("vClone requires TIA pairs (P-VOL %d, MU %d)", *input.PvolLdevID, *input.MuNumber)
	}

	// use old
	// 2. Secondary Lookup: Always check the Snapshot Family (Hardware Truth)
	// We fetch this regardless of whether the snapshot was found to handle transitions
	familyMember, fErr := psm.getVcloneFromFamily(&pvolID, &muNum, input.SvolLdevID, false)
	if fErr != nil {
		log.WriteWarn("Validation| Could not fetch family metadata for P-VOL %d: %v", pvolID, fErr)
	}

	if err != nil {
		// If not found in traditional table, check if it's already a promoted vClone
		if isNotFoundError(err) {
			log.WriteInfo("Validation| Snapshot pair not found. Checking Family metadata for P-VOL %d", pvolID)

			if familyMember != nil && (familyMember.IsVirtualCloneVolume || familyMember.IsVirtualCloneParentVolume) {
				log.WriteInfo("Validation| Verified relationship exists in Family table. Treating as existing vClone.")
				return nil, familyMember, nil
			}

			return nil, nil, log.WriteAndReturnError("snapshot relationship for P-VOL %d MU %d not found in traditional or family tables", pvolID, muNum)
		}
		// Return communication/system errors
		return nil, nil, log.WriteAndReturnError("error fetching snapshot metadata: %v", err)
	}

	// 3. Busy Check for Traditional Pairs
	if snapshot.PvolProcessingStatus == "P" {
		return nil, nil, log.WriteAndReturnError("snapshot %d,%d is busy (Processing Status: P). Cannot perform vClone action", pvolID, muNum)
	}

	log.WriteInfo("Validation| Successful for P-VOL %d, MU %d (Status: %s)", pvolID, muNum, snapshot.Status)

	return snapshot, familyMember, nil
}

func (psm *sanStorageManager) validateTIAVrestoreInput(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pvolID := *input.PvolLdevID
	muNum := *input.MuNumber

	// Hardware Model Check
	model, isBSeries := checkIsVspBSeries(psm.storageSetting.Serial)
	if !isBSeries {
		return nil, nil, log.WriteAndReturnError("vClone operations only supported on VSP B series. Model: %s", model)
	}

	// 1. Fetch Family metadata early (Used for "Pair it back" logic and final return)
	familyMember, fErr := psm.getVcloneFromFamily(input.PvolLdevID, input.MuNumber, input.SvolLdevID, false)
	if fErr != nil {
		log.WriteDebug("Validation| Could not fetch family metadata for P-VOL %d: %v", pvolID, fErr)
	}

	// 2. Attempt to fetch traditional snapshot metadata
	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)

	// 3. Handle Missing Traditional Pair
	if err != nil {
		if isNotFoundError(err) {
			log.WriteInfo("Traditional pair not found for %d,%d. Checking Family table for vClone state...", pvolID, muNum)

			if familyMember != nil && familyMember.IsVirtualCloneVolume {
				log.WriteInfo("Validation| Standalone vClone detected. Returning for re-pairing.")
				// RETURN NO ERROR HERE to allow reconcileSnapshotVCloneRestore to handle repair
				return nil, familyMember, nil
			}
			return nil, nil, log.WriteAndReturnError("vRestore failed: snapshot pair %d,%d not found in traditional or family tables.", pvolID, muNum)
		}
		return nil, nil, err
	}

	// 4. Strict 'Paired Back' Validation
	// Check traditional snapshot flags
	if !snapshot.IsVirtualCloneVolume || !snapshot.IsVirtualCloneParentVolume {
		return snapshot, familyMember, log.WriteAndReturnError(
			"vRestore failed: snapshot %d,%d is not in a valid 'paired back' vClone state. "+
				"Both isVirtualCloneVolume and isVirtualCloneParentVolume must be true.",
			pvolID, muNum,
		)
	}

	// 5. Status Check for the Paired Back relationship
	if snapshot.Status != "PAIR" && snapshot.Status != "PSUS" {
		return snapshot, familyMember, log.WriteAndReturnError("vRestore failed: paired-back snapshot %d,%d is in invalid status: %s", pvolID, muNum, snapshot.Status)
	}

	log.WriteInfo("Validation successful: Snapshot %d,%d is paired back and ready for vRestore.", pvolID, muNum)
	return snapshot, familyMember, nil
}

// --- Wait funcs ---

func (psm *sanStorageManager) waitForSnapshotVCloneStatus(pvolID *int, svolID *int, muNum *int) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// maxRetries := 15
	// delay := 10 * time.Second

	retryCfg := utils.GetRetryConfig()
	maxRetries := retryCfg.MaxRetries
	delay := time.Duration(retryCfg.Delay) * time.Second

	snapshotIsGone := false

	for i := 0; i < maxRetries; i++ {
		log.WriteInfo("Polling vClone status (Attempt %d/%d) for P-VOL %d", i+1, maxRetries, *pvolID)

		if !snapshotIsGone {
			// 1. Check if traditional pair exists
			_, err := psm.getSnapshot(pvolID, muNum, true)

			if err == nil {
				log.WriteDebug("Traditional pair %d:%d still active. Waiting...", *pvolID, *muNum)
				goto wait
			}

			if isNotFoundError(err) {
				log.WriteInfo("Traditional relationship is gone. Switching to Family metadata check.")
				snapshotIsGone = true
			} else {
				return nil, nil, err
			}
		}

		// 2. Check Family data (Looking for IsVirtualCloneVolume)
		if snapshotIsGone {
			log.WriteInfo("Checking Family metadata for S-VOL %d...", *svolID)

			familyMember, fErr := psm.getVcloneFromFamily(pvolID, muNum, svolID, true)

			if fErr == nil && familyMember != nil {
				// Since there is no status, we rely on the boolean flag and parent link
				if familyMember.IsVirtualCloneVolume && familyMember.ParentLdevID == *pvolID {
					log.WriteInfo("vClone confirmed: S-VOL %d is now a Virtual Clone of P-VOL %d", *svolID, *pvolID)
					return nil, familyMember, nil
				}
				log.WriteDebug("Family record found for %d, but IsVirtualCloneVolume is false.", *svolID)
			} else {
				log.WriteDebug("No Family record found yet for S-VOL %d.", *svolID)
			}
		}

	wait:
		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return nil, nil, fmt.Errorf("timeout: S-VOL %d never promoted to Virtual Clone after snapshot deletion", *svolID)
}

func (psm *sanStorageManager) waitForSnapshotVRestoreStatus(pvolID *int, svolID *int, muNum *int) (*gwymodel.Snapshot, *gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// maxRetries := 40
	// delay := 15 * time.Second

	retryCfg := utils.GetRetryConfig()
	maxRetries := retryCfg.MaxRetries + 10 // additional retries for vrestore
	delay := time.Duration(retryCfg.Delay) * time.Second

	for i := 0; i < maxRetries; i++ {
		log.WriteInfo("Polling vRestore status (Attempt %d/%d) for P-VOL %d (MU %d)", i+1, maxRetries, *pvolID, *muNum)

		// 1. Get current Snapshot data
		snap, err := psm.getSnapshot(pvolID, muNum, true)
		if err != nil {
			log.WriteDebug("Waiting for snapshot record to be accessible... (Error: %v)", err)
			goto wait
		}

		log.WriteInfo("Current Status: %s | isVClone: %t | isVCParent: %t",
			snap.Status, snap.IsVirtualCloneVolume, snap.IsVirtualCloneParentVolume)

		// 2. Logic Fix: vRestore result is PSUS with flags set to false
		// This indicates the restore is complete and the S-VOL is now a standard 'split' snapshot
		if snap.Status == "PSUS" && !snap.IsVirtualCloneVolume && !snap.IsVirtualCloneParentVolume {
			log.WriteInfo("vRestore success: Pair reached PSUS and vClone metadata has been cleared (restored to snapshot state).")

			// 3. Fetch Family metadata as a final check
			familyMember, fErr := psm.getVcloneFromFamily(pvolID, muNum, svolID, true)
			if fErr != nil {
				log.WriteWarn("vRestore| Could not fetch family metadata, but snapshot status is valid.")
			}

			return snap, familyMember, nil
		}

		// 4. Handle transitionary status
		if snap.Status == "RCPY" || snap.Status == "PFUS" {
			log.WriteDebug("vRestore in progress (Status: %s)...", snap.Status)
		} else if snap.Status == "PAIR" {
			log.WriteDebug("Pair still in PAIR status. Waiting for vRestore command to take effect...")
		}

	wait:
		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return nil, nil, fmt.Errorf("timeout: vRestore failed to reach PSUS status for P-VOL %d, MU %d", *pvolID, *muNum)
}

// --- Others ---

func (psm *sanStorageManager) getUniversalInfo(pvolID *int, svolID *int, forceRefresh bool) (*reconmodel.SnapshotUniversalInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Wrap IDs into the input struct the fetcher expects
	input := reconmodel.SnapshotReconcilerInput{
		PvolLdevID: pvolID,
		SvolLdevID: svolID,
	}

	// If we have a cached pool from the validation step, the fetcher will use it
	// unless forceRefresh is true.
	pool, pvol, svol, err := psm.fetchRequiredResources(input, forceRefresh)
	if err != nil {
		return nil, err
	}

	isDRS := slices.Contains(pvol.Attributes, "DRS")

	info := &reconmodel.SnapshotUniversalInfo{
		StorageSerial:  psm.storageSetting.Serial,
		PvolAttributes: pvol.Attributes,
		// isDRS:               isDRS,
		IsThinImageAdvanced: isDRS,
	}

	if pool != nil {
		info.SnapshotPoolType = pool.PoolType
	}

	if svol != nil {
		info.SvolAttributes = svol.Attributes
	}

	return info, nil
}

func (psm *sanStorageManager) fetchRequiredResources(input reconmodel.SnapshotReconcilerInput, forceRefreshLdev bool) (*provmodel.DynamicPool, *gwymodel.LogicalUnit, *gwymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, nil, err
	}

	// --- 1. Pool Metadata (Immutable for this command) ---
	poolID := -1
	if input.SnapshotPoolID != nil {
		poolID = *input.SnapshotPoolID
	}

	// Only fetch if CachePool is empty or the ID changed (unlikely in one run)
	if CachePool == nil || (poolID != -1 && CachePool.PoolID != poolID) {
		if poolID != -1 {
			p, err := provObj.GetDynamicPoolById(poolID)
			if err != nil {
				return nil, nil, nil, log.WriteAndReturnError("failed to fetch pool %d: %v", poolID, err)
			}
			CachePool = p
		}
	}

	// --- 2. P-VOL Metadata (Refreshable) ---
	pvolID := *input.PvolLdevID

	// Fetch if Cache is empty, ID is different, or we are forcing a refresh
	if forceRefreshLdev || CachePvol == nil || CachePvol.LdevID != pvolID {
		p, err := provObj.GetLun(pvolID)
		if err != nil {
			return nil, nil, nil, log.WriteAndReturnError("failed to fetch P-VOL %d: %v", pvolID, err)
		}
		if p == nil || p.EmulationType == "NOT DEFINED" {
			return nil, nil, nil, log.WriteAndReturnError("P-VOL %d not found", pvolID)
		}
		CachePvol = p
	}

	// --- 3. S-VOL Metadata (Refreshable) ---
	if input.SvolLdevID != nil && *input.SvolLdevID != 0 {
		svolID := *input.SvolLdevID

		if forceRefreshLdev || CacheSvol == nil || CacheSvol.LdevID != svolID {
			s, err := provObj.GetLun(svolID)
			if err != nil {
				return nil, nil, nil, log.WriteAndReturnError("failed to fetch S-VOL %d: %v", svolID, err)
			}
			CacheSvol = s
		}
	}

	return CachePool, CachePvol, CacheSvol, nil
}

func (psm *sanStorageManager) fetchRequiredResourcesConcurrent(input reconmodel.SnapshotReconcilerInput, forceRefreshLdev bool) (*provmodel.DynamicPool, *gwymodel.LogicalUnit, *gwymodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, nil, err
	}

	var wg sync.WaitGroup
	// Buffer the channel so goroutines don't block if we return early on error
	errChan := make(chan error, 3)

	// --- 1. Concurrent Pool Metadata ---
	poolID := -1
	if input.SnapshotPoolID != nil {
		poolID = *input.SnapshotPoolID
	}
	if poolID != -1 && (CachePool == nil || CachePool.PoolID != poolID) {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := provObj.GetDynamicPoolById(id)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch pool %d: %w", id, err)
				return
			}
			CachePool = p
		}(poolID)
	}

	// --- 2. Concurrent P-VOL Metadata ---
	pvolID := *input.PvolLdevID
	if forceRefreshLdev || CachePvol == nil || CachePvol.LdevID != pvolID {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := provObj.GetLun(id)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch P-VOL %d: %w", id, err)
				return
			}
			CachePvol = p
		}(pvolID)
	}

	// --- 3. Concurrent S-VOL Metadata ---
	if input.SvolLdevID != nil && *input.SvolLdevID != 0 {
		svolID := *input.SvolLdevID
		if forceRefreshLdev || CacheSvol == nil || CacheSvol.LdevID != svolID {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				s, err := provObj.GetLun(id)
				if err != nil {
					errChan <- fmt.Errorf("failed to fetch S-VOL %d: %w", id, err)
					return
				}
				CacheSvol = s
			}(svolID)
		}
	}

	// Wait for all requests to finish
	wg.Wait()
	close(errChan)

	// Check if any goroutine sent an error
	if err := <-errChan; err != nil {
		return nil, nil, nil, log.WriteAndReturnError("TFError| %v", err)
	}

	return CachePool, CachePvol, CacheSvol, nil
}

func (psm *sanStorageManager) getSnapshot(pvolID *int, mu *int, forceRefresh bool) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Validation
	if err := psm.validatePvolAndMu(pvolID, mu); err != nil {
		return nil, err
	}

	// 2. Check Global Cache
	// If cache exists and we aren't forcing a refresh, return the cached snapshot
	cacheLock.Lock()
	if !forceRefresh && snapshotGlobalCache != nil && snapshotGlobalCache.Snapshot != nil {
		log.WriteInfo("GlobalCache| Returning cached snapshot for P-VOL %d, MU %d", *pvolID, *mu)
		cachedSnap := snapshotGlobalCache.Snapshot
		cacheLock.Unlock()
		return cachedSnap, nil
	}
	cacheLock.Unlock()

	// 3. Backend Fetch
	log.WriteInfo("GlobalCache| Fetching fresh snapshot from backend (forceRefresh=%t)", forceRefresh)
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("TFError| failed to get provisioner manager: %v", err)
		return nil, err
	}

	snapshot, err := provObj.GetSnapshot(*pvolID, *mu)
	if err != nil {
		// We don't return an error here if it's a 404, but we let the caller handle it
		log.WriteError("TFError| failed to get snapshot for P-VOL %d, MU %d: %v", *pvolID, *mu, err)
		return nil, err
	}

	// 4. Update Global Cache
	// We update the Snapshot field specifically.
	// If the cache object doesn't exist yet, we initialize it.
	cacheLock.Lock()
	if snapshotGlobalCache == nil {
		snapshotGlobalCache = &reconmodel.ReconcileSnapshotResult{}
	}
	snapshotGlobalCache.Snapshot = snapshot
	cacheLock.Unlock()

	return snapshot, nil
}

func (psm *sanStorageManager) getVcloneFromFamily(pvolID *int, mu *int, svolID *int, forceRefresh bool) (*gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if pvolID == nil {
		return nil, nil
	}

	// 1. Check Global Cache
	cacheLock.Lock()
	if !forceRefresh && snapshotGlobalCache != nil && snapshotGlobalCache.VcloneFamily != nil {
		cf := snapshotGlobalCache.VcloneFamily

		// Ensure the cached member is still a Virtual Clone
		if cf.IsVirtualCloneVolume {
			match := false
			if svolID != nil && cf.LdevID == *svolID {
				match = true
			} else if mu != nil && cf.MuNumber == *mu {
				match = true
			}

			if match {
				log.WriteInfo("GlobalCache| Returning cached Virtual Clone for S-VOL %d", cf.LdevID)
				cacheLock.Unlock()
				return cf, nil
			}
		}
	}
	cacheLock.Unlock()

	// 2. Fetch from Backend
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	resp, err := provObj.GetSnapshotFamily(*pvolID)
	if err != nil || resp == nil {
		return nil, err
	}

	// 3. Filter for the specific member AND verify it's a vClone
	var foundMember *gwymodel.SnapshotFamily
	for i := range resp.Data {
		member := &resp.Data[i]

		// Condition 1: Must be a Virtual Clone
		if !member.IsVirtualCloneVolume {
			continue
		}

		if member.PrimaryOrSecondary != "" { // we wan't the svol, not pvol
			continue
		}

		// Condition 2: Must match the requested S-VOL ID or MU
		match := false
		if svolID != nil && member.LdevID == *svolID {
			match = true
		} else if member.ParentLdevID == *pvolID {
			match = true
		} else if mu != nil && member.MuNumber == *mu {
			match = true
		}

		if match {
			foundMember = member
			break
		}
	}

	// 4. Update Global Cache
	if foundMember != nil {
		cacheLock.Lock()
		if snapshotGlobalCache == nil {
			snapshotGlobalCache = &reconmodel.ReconcileSnapshotResult{}
		}
		snapshotGlobalCache.VcloneFamily = foundMember
		cacheLock.Unlock()
		log.WriteInfo("GlobalCache| Updated cache with verified Virtual Clone S-VOL %d", foundMember.LdevID)
	}

	return foundMember, nil
}

func (psm *sanStorageManager) checkThenGetVclone(input reconmodel.SnapshotReconcilerInput, snapshot *gwymodel.Snapshot) (*gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	model, isVspB := checkIsVspBSeries(psm.storageSetting.Serial)
	if !isVspB {
		log.WriteDebug("Get vclone family is only supported in B20/B85 series. Current model: %s", model)
		return nil, nil
	}

	tiaInput := false
	if snapshot == nil {
		// we do this to find out if the input is a tia if no snapshot
		pool, pvol, _, err := psm.fetchRequiredResources(input, false)
		if err != nil {
			return nil, err
		}
		tiaInput = isTIAInput(pvol, pool)
	}

	if (snapshot == nil && tiaInput) || (snapshot != nil && isTIAdvancedPair(snapshot)) {
		vcloneFamily, fErr := psm.getVcloneFromFamily(input.PvolLdevID, input.MuNumber, input.SvolLdevID, true)
		if fErr != nil {
			// just log, don't error out
			log.WriteDebug("TFDebug| Could not fetch family info: %v", fErr)
		} else if vcloneFamily != nil {
			// Since there is no status, we rely on the boolean flag and parent link
			if vcloneFamily.IsVirtualCloneVolume && vcloneFamily.ParentLdevID == *input.PvolLdevID {
				log.WriteInfo("vClone confirmed: S-VOL %d is a Virtual Clone of P-VOL %d", vcloneFamily.SvolLdevID, *input.PvolLdevID)
				return vcloneFamily, nil
			}
			log.WriteDebug("Family record found for %d, but IsVirtualCloneVolume is false.", input.SvolLdevID)
		} else {
			log.WriteDebug("No Family record found for S-VOL %d.", input.SvolLdevID)
		}
	}

	return nil, nil
}

func (psm *sanStorageManager) findValidTargetPort() (string, error) {
	// 1. Fetch ports with FIBRE type
	ports, err := psm.GetStoragePorts(nil, "FIBRE", "")
	if err != nil || ports == nil {
		return "", fmt.Errorf("failed to retrieve storage ports: %v", err)
	}

	// 2. Iterate and find the first one with "TAR" attribute
	for _, port := range *ports {
		for _, attr := range port.PortAttributes {
			if attr == "TAR" {
				return port.PortId, nil
			}
		}
	}

	return "", fmt.Errorf("no FIBRE ports with TAR attribute found for TIA repair")
}

func (psm *sanStorageManager) createRepairHostGroup(pvolID int) error {
	log := commonlog.GetLogger()

	portID, err := psm.findValidTargetPort()
	if err != nil {
		return err
	}

	hgName := fmt.Sprintf("TIA_REPAIR_%d", pvolID)
	myHostMode := "Standard"

	// Create the Luns slice using your model
	// We pass only the LdevID; omitting Lun allows the storage to auto-assign the next available LUN
	myLdevIds := []reconmodel.Luns{
		{
			LdevId: &pvolID,
		},
	}

	crReq := reconmodel.CreateHostGroupRequest{
		PortID:        &portID,
		HostGroupName: &hgName,
		HostMode:      &myHostMode,
		Ldevs:         myLdevIds,
	}

	log.WriteInfo("TIA| Mapping P-VOL %d to port %s (Auto-LUN)", pvolID, portID)

	_, err = psm.createHostGroup(&crReq)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}

func isVolVCloneParent(lun *gwymodel.LogicalUnit) bool {
	return hasVolAttribute(lun, "VCP")
}

func isVolVClone(lun *gwymodel.LogicalUnit) bool {
	return hasVolAttribute(lun, "VC")
}

func isTIACompatible(lun *gwymodel.LogicalUnit) bool {
	// TIA requires both HDP (Thin Provisioning) and DRS (Data Reduction)
	return isVolHDP(lun) && isVolDRS(lun)
}

func isVolHTI(lun *gwymodel.LogicalUnit) bool {
	return hasVolAttribute(lun, "HTI")
}
