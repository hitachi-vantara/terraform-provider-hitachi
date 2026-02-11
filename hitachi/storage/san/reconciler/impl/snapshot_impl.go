package sanstorage

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	provmanager "terraform-provider-hitachi/hitachi/storage/san/provisioner"
	provimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

// --- Datasources ---

func (psm *sanStorageManager) ReconcileGetSnapshot(pvolID *int, mu *int) (*gwymodel.Snapshot, error) {
	return psm.getSnapshot(pvolID, mu, false)
}

func (psm *sanStorageManager) ReconcileGetMultipleSnapshots(input reconmodel.SnapshotGetMultipleInput) (*gwymodel.SnapshotListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Starting GetMultipleSnapshots")

	// Validate Input
	if err := psm.validateSnapshotGetMultipleInput(input); err != nil {
		return nil, err
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("TFError| failed to get provisioner manager: %v", err)
		return nil, err
	}

	params := gwymodel.GetSnapshotsParams{
		SnapshotGroupName: input.SnapshotGroupName,
		PvolLdevID:        input.PvolLdevID,
		SvolLdevID:        input.SvolLdevID,
		MuNumber:          input.MuNumber,
	}

	snapshotList, err := provObj.GetSnapshots(params)
	if err != nil {
		log.WriteError("TFError| failed to get snapshots with params %+v: %v", params, err)
		return nil, err
	}

	log.WriteInfo("Successfully retrieved %d snapshots", len(snapshotList.Data))
	return snapshotList, nil
}

func (psm *sanStorageManager) ReconcileGetMultipleSnapshotsRange(input reconmodel.SnapshotGetMultipleRangeInput) (*gwymodel.SnapshotListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Determine Range
	startID := 0
	if input.StartPvolLdevID != nil {
		startID = *input.StartPvolLdevID
	}
	endID := 65535 // Default to max LDEV ID
	if input.EndPvolLdevID != nil {
		endID = *input.EndPvolLdevID
	}

	// 2. Fetch all snapshot data using the high-performance Group pattern
	log.WriteInfo("Fetching all Snapshot Groups to identify pairs in range %d to %d", startID, endID)

	groupResp, err := psm.ReconcileGetMultipleSnapshotGroups(true)
	if err != nil {
		log.WriteError("TFError| failed to fetch snapshot groups for indexing: %v", err)
		return nil, err
	}

	if groupResp == nil || len(groupResp.Data) == 0 {
		return &gwymodel.SnapshotListResponse{Data: []gwymodel.Snapshot{}}, nil
	}

	// 3. Filter the aggregated snapshots by the requested P-VOL range
	var results []gwymodel.Snapshot
	seen := make(map[string]bool) // Prevent duplicates if a pair appears in multiple views

	for _, group := range groupResp.Data {
		for _, snap := range group.Snapshots {
			// Check if this snapshot falls within the user's requested LDEV range
			if snap.PvolLdevID >= startID && snap.PvolLdevID <= endID {

				// Create a unique key for PvolID + MuNumber
				key := fmt.Sprintf("%d-%d", snap.PvolLdevID, snap.MuNumber)
				if !seen[key] {
					results = append(results, snap)
					seen[key] = true
				}
			}
		}
	}

	// 4. Sort results for Terraform consistency
	sort.Slice(results, func(i, j int) bool {
		if results[i].PvolLdevID == results[j].PvolLdevID {
			return results[i].MuNumber < results[j].MuNumber
		}
		return results[i].PvolLdevID < results[j].PvolLdevID
	})

	log.WriteInfo("Filtered %d snapshots from group index within range %d-%d", len(results), startID, endID)
	return &gwymodel.SnapshotListResponse{Data: results}, nil
}

// --- Resource ---

// Action		Primary Direction		Target of Change		Data Impact
// Split		P to S (Stop sync)		S-VOL					Freezes S-VOL at point-in-time.
// Resync		P to S (Start sync)		S-VOL					Overwrites S-VOL with latest P-VOL data.
// Restore		S to P					P-VOL					Overwrites Production (P-VOL) with Snapshot (S-VOL) data.
// Clone		P to New				New Volume				Creates a physically independent disk.

func (psm *sanStorageManager) ReconcileSnapshot(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.Action == nil || *input.Action == "" {
		input.Action = utils.Ptr("read")
	}

	action := strings.ToLower(*input.Action)

	log.WriteInfo("Starting snapshot reconciliation for operation: %s", action)

	// pvol and mu are set in impl code
	log.WriteDebug("Input PvolLdevID: %v, MuNumber: %v", input.PvolLdevID, input.MuNumber)

	switch action {
	case "create":
		return psm.reconcileSnapshotCreateOrReadExisting(input)
	case "split":
		return psm.reconcileSnapshotSplit(input)
	case "resync":
		return psm.reconcileSnapshotResync(input)
	case "restore":
		return psm.reconcileSnapshotRestore(input)
	case "clone":
		return psm.reconcileSnapshotClone(input)
	case "assign_svol":
		return psm.reconcileSnapshotAssign(input)
	case "unassign_svol":
		return psm.reconcileSnapshotUnassign(input)
	case "update_retention_period":
		return psm.reconcileSetSnapshotRetentionPeriod(input)
	case "delete":
		return psm.ReconcileSnapshotDelete(input)
	case "defrag":
		return psm.reconcileGarbageData(input)
	case "deletetree":
		return psm.reconcileSnapshotTreeDelete(input)
	default:
		return nil, log.WriteAndReturnError("unsupported snapshot operation: %s", *input.Action)
	}
}

// --- Action Handlers ---

func (psm *sanStorageManager) reconcileSnapshotCreateOrReadExisting(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Read existing state
	result, err := psm.ReconcileReadExistingSnapshotVclone(input)
	if err != nil {
		return nil, err
	}

	if result != nil {
		// Case A: Traditional Snapshot exists
		if result.Snapshot != nil {
			backend := result.Snapshot
			log.WriteInfo("TFInfo| Snapshot pair exists (P-VOL %d, MU %d).", backend.PvolLdevID, backend.MuNumber)

			// A.1 Check S-VOL Assignment Drift
			inputSvol := 0
			if input.SvolLdevID != nil {
				inputSvol = *input.SvolLdevID
			}

			if inputSvol != 0 && backend.SvolLdevID == 0 {
				log.WriteInfo("TFInfo| S-VOL Assignment: Input %d, Backend 0. Running Assign.", inputSvol)
				return psm.reconcileSnapshotAssign(input)
			}

			if inputSvol == 0 && backend.SvolLdevID != 0 {
				log.WriteInfo("TFInfo| S-VOL Unassignment: Input 0, Backend %d. Running Unassign.", backend.SvolLdevID)
				return psm.reconcileSnapshotUnassign(input)
			}

			// A.2 Check Retention Period Drift with Documentation Constraints
			if input.RetentionPeriod != nil && *input.RetentionPeriod != backend.RetentionPeriod {
				newVal := *input.RetentionPeriod
				currVal := backend.RetentionPeriod

				// Hitachi Doc: "You cannot specify a value smaller than the current value."
				if newVal < currVal {
					return nil, log.WriteAndReturnError(
						"Cannot decrease RetentionPeriod for TIA pairs. (Requested: %d, Current: %d)",
						newVal, currVal,
					)
				} else {
					log.WriteInfo("TFInfo| Retention Update: Increasing from %d to %d hours.", currVal, newVal)
					return psm.reconcileSetSnapshotRetentionPeriod(input)
				}
			}

			return backend, nil
		}

		// Case B: vClone exists
		if result.VcloneFamily != nil && result.VcloneFamily.IsVirtualCloneVolume {
			log.WriteInfo("TFInfo| vClone detected. Re-creating pair for restore capability.")
		}
	}

	// 2. Fall through: No pair exists, or it's a vClone that needs a pair re-created
	return psm.reconcileSnapshotCreate(input)
}

func (psm *sanStorageManager) reconcileSnapshotCreate(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	err := psm.validateSnapshotCreateInput(input)
	if err != nil {
		return nil, err
	}

	request := gwymodel.CreateSnapshotParams{
		PvolLdevID:               *input.PvolLdevID,
		SnapshotPoolID:           *input.SnapshotPoolID,
		SnapshotGroupName:        *input.SnapshotGroupName,
		MuNumber:                 input.MuNumber,
		IsClone:                  input.IsClone,
		SvolLdevID:               input.SvolLdevID,
		IsConsistencyGroup:       input.IsConsistencyGroup,
		AutoSplit:                input.AutoSplit,
		CanCascade:               input.CanCascade,
		ClonesAutomation:         input.ClonesAutomation,
		CopySpeed:                input.CopySpeed,
		IsDataReductionForceCopy: input.IsDataReductionForceCopy,
		RetentionPeriod:          input.RetentionPeriod,
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	resID, err := provObj.CreateSnapshot(request)
	if err != nil {
		return nil, err
	}

	// Use the helper to get pvol and mu from snapshot id
	pvolID, muID, err := parseSnapshotResID(resID)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PAIR"}
	tiaTargetStatuses := []string{"PAIR"}
	snapshot, err := psm.waitForSnapshotStatus(pvolID, muID, tiTargetStatuses, tiaTargetStatuses)
	return snapshot, err
}

// reconcileSnapshotSplit handles splitting an existing snapshot pair
// Not for a Thin Image pair that belongs to a snapshot group in the consistency group mode (CTG mode).
func (psm *sanStorageManager) reconcileSnapshotSplit(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateTIOperationInput(input)
	if err != nil {
		return nil, err
	}

	if input.RetentionPeriod != nil && !isTIAdvancedPair(snapshot) {
		log.WriteInfo("RetentionPeriod parameter is ignored for Thin Image Standard pairs during split operation")
	}

	if snapshot.IsClone {
		return nil, log.WriteAndReturnError("cannot split: snapshot %d:%d was created as a Clone; cannot split", *input.PvolLdevID, *input.MuNumber)
	}

	// --- Status Validation ---
	currentStatus := snapshot.Status
	if currentStatus == "PSUS" || currentStatus == "PFUS" || currentStatus == "PFUL" {
		log.WriteInfo("already split: pair %d,%d is already in a split status (PSUS/PFUS/PFUL) (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
		return snapshot, err
	}
	if currentStatus != "PAIR" {
		return nil, log.WriteAndReturnError("cannot split: pair %d,%d must be in PAIR status (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
	}

	request := gwymodel.SplitSnapshotRequest{
		Parameters: gwymodel.SplitSnapshotParams{
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.SplitSnapshot(*input.PvolLdevID, *input.MuNumber, request)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PSUS", "PFUS"}
	tiaTargetStatuses := []string{"PSUS", "PFUS", "PFUL"}
	snapshot, err = psm.waitForSnapshotStatus(input.PvolLdevID, input.MuNumber, tiTargetStatuses, tiaTargetStatuses)
	return snapshot, err
}

// reconcileSnapshotResync handles resynchronizing an existing snapshot pair
// Not for a Thin Image pair that belongs to a snapshot group in the consistency group mode (CTG mode).
// When the pair is resynchronized, all snapshot data will be deleted.
// You can store new snapshot data by specifying the setting auto-split the resynchronized pair.
func (psm *sanStorageManager) reconcileSnapshotResync(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateTIOperationInput(input)
	if err != nil {
		return nil, err
	}

	if input.RetentionPeriod != nil && !isTIAdvancedPair(snapshot) {
		log.WriteInfo("RetentionPeriod parameter is ignored for Thin Image Standard pairs during resync operation")
	}

	if input.RetentionPeriod != nil && (input.AutoSplit == nil || *input.AutoSplit == false) {
		return nil, log.WriteAndReturnError("RetentionPeriod can only be specified if the autoSplit attribute is set to true for resync operation of TIA pair")
	}

	err = checkTIARetention(snapshot)
	if err != nil {
		return nil, log.WriteAndReturnError(err.Error())
	}

	// --- Status Validation ---
	currentStatus := snapshot.Status
	if currentStatus == "PAIR" || currentStatus == "COPY" {
		log.WriteInfo("already synced: pair %d,%d is already in a synced status (PAIR/COPY) (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
		return snapshot, err
	}
	if currentStatus != "PSUS" && currentStatus != "PFUS" && currentStatus != "PFUL" {
		return nil, log.WriteAndReturnError("cannot resync: pair %d,%d must be in a split status (PSUS/PFUS/PFUL) (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
	}

	params := gwymodel.ResyncSnapshotRequest{
		Parameters: gwymodel.ResyncSnapshotParams{
			AutoSplit:       input.AutoSplit,
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.ResyncSnapshot(*input.PvolLdevID, *input.MuNumber, params)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PAIR"}
	tiaTargetStatuses := []string{"PAIR"}
	if input.AutoSplit != nil && *input.AutoSplit == true {
		tiTargetStatuses = []string{"PSUS", "PFUS"}
		tiaTargetStatuses = []string{"PSUS", "PFUS", "PFUL"}
	}

	snapshot, err = psm.waitForSnapshotStatus(input.PvolLdevID, input.MuNumber, tiTargetStatuses, tiaTargetStatuses)
	return snapshot, err
}

// reconcileSnapshotRestore handles restoring P-VOL data from an S-VOL
// Not for a Thin Image pair that belongs to a snapshot group in the consistency group mode (CTG mode).
// When the pair is restored, the data of the snapshot specified for the primary volume is overwritten.
// For Thin Image Advanced pairs, autoSplit attribute is ignored even if it is specified.
func (psm *sanStorageManager) reconcileSnapshotRestore(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateTIOperationInput(input)
	if err != nil {
		return nil, err
	}

	isTIA := isTIAdvancedPair(snapshot)
	if isTIA {
		log.WriteInfo("autoSplit parameter is ignored for Thin Image Advanced pairs during restore operation")
	}

	// --- Status Validation ---
	currentStatus := snapshot.Status
	if currentStatus == "PAIR" || currentStatus == "COPY" {
		log.WriteInfo("already synced: pair %d,%d is already in a synced status (PAIR/COPY) (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
		return snapshot, err
	}
	if currentStatus != "PSUS" && currentStatus != "PFUS" && currentStatus != "PFUL" {
		return nil, log.WriteAndReturnError("cannot restore: pair %d,%d must be in a split status (PSUS/PFUS/PFUL) (current: %s)",
			*input.PvolLdevID, *input.MuNumber, currentStatus)
	}

	params := gwymodel.RestoreSnapshotRequest{
		Parameters: gwymodel.RestoreSnapshotParams{
			AutoSplit: input.AutoSplit,
		},
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.RestoreSnapshot(*input.PvolLdevID, *input.MuNumber, params)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PAIR"}
	tiaTargetStatuses := []string{"PAIR", "PSUS"} // TIA does not go to PAIR

	// If it's NOT TIA, AutoSplit actually works.
	// If it's TIA, AutoSplit is ignored
	if !isTIA && input.AutoSplit != nil && *input.AutoSplit == true {
		tiTargetStatuses = []string{"PSUS", "PFUS"}
		// Note: We don't change tiaTargetStatuses here because the B20/TIA ignores the flag
	}

	snapshot, err = psm.waitForSnapshotStatus(input.PvolLdevID, input.MuNumber, tiTargetStatuses, tiaTargetStatuses)
	return snapshot, err
}

// reconcileSnapshotClone handles cloning a snapshot to a new volume
// Only for TI Standard pairs. Not for Thin Image Advanced pairs
func (psm *sanStorageManager) reconcileSnapshotClone(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateTIOperationInput(input)
	if err != nil {
		return nil, err
	}

	if isTIAdvancedPair(snapshot) {
		return nil, log.WriteAndReturnError("Clone operation can only be performed on Thin Image Standard pair, not TIA pair (P-VOL %d, MU %d)",
			*input.PvolLdevID, *input.MuNumber)
	}

	if !snapshot.IsClone {
		return nil, log.WriteAndReturnError("snapshot %d:%d was not created as a Clone; cannot clone", *input.PvolLdevID, *input.MuNumber)
	}

	// Should not be split already
	// // --- Status Validation ---
	// currentStatus := snapshot.Status
	// if currentStatus != "PSUS" && currentStatus != "PFUS" {
	// 	return nil, log.WriteAndReturnError("cannot clone: pair %d,%d must be split (PSUS/PFUS) to create a physical clone (current: %s)",
	// 		*input.PvolLdevID, *input.MuNumber, currentStatus)
	// }

	params := gwymodel.CloneSnapshotRequest{
		Parameters: gwymodel.CloneSnapshotParams{
			CopySpeed: input.CopySpeed,
		},
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.CloneSnapshot(*input.PvolLdevID, *input.MuNumber, params)
	if err != nil {
		return nil, err
	}

	// at this point, TI snapshot is gone
	err = psm.waitForSnapshotDeletion(input.PvolLdevID, input.MuNumber)
	return nil, err
}

func (psm *sanStorageManager) reconcileSnapshotAssign(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateAssignVolumeInput(input)
	if err != nil {
		return nil, err
	}

	currentStatus := snapshot.Status
	if currentStatus == "PAIR" {
		log.WriteInfo("Assign: S-VOL is in PAIR status. Splitting snapshot P-VOL %d MU %d before assign.", *input.PvolLdevID, *input.MuNumber)
		snapshot, err = psm.reconcileSnapshotSplit(input)
		if err != nil {
			return nil, log.WriteAndReturnError("Assign failed: Splitting snapshot P-VOL %d MU %d before assign failed. %+v", *input.PvolLdevID, *input.MuNumber, err)
		}
	}

	if snapshot.Status != "PSUS" {
		return nil, log.WriteAndReturnError("Assign failed: snapshot %d,%d must be in PSUS status (current: %s)",
			*input.PvolLdevID, *input.MuNumber, snapshot.Status)
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	params := gwymodel.AssignSnapshotVolumeRequest{
		Parameters: gwymodel.AssignSnapshotVolumeParams{
			SvolLdevID: *input.SvolLdevID,
		},
	}

	log.WriteInfo("Assign| Executing assign-volume: P-VOL %d, MU %d -> S-VOL %d",
		*input.PvolLdevID, *input.MuNumber, *input.SvolLdevID)

	_, err = provObj.AssignSnapshotVolume(*input.PvolLdevID, *input.MuNumber, params)
	if err != nil {
		return nil, log.WriteAndReturnError("Assign command failed: %v", err)
	}

	return psm.waitForSnapshotStatus(input.PvolLdevID, input.MuNumber, nil, []string{"PSUS"})
}

// reconcileSnapshotUnassign handles unassigning a volume from a snapshot
func (psm *sanStorageManager) reconcileSnapshotUnassign(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateUnassignVolumeInput(input)
	if err != nil {
		return nil, err
	}

	if snapshot.SvolLdevID == 0 {
		log.WriteInfo("Unassign: S-VOL is already unassigned for P-VOL %d MU %d. Skipping API call.",
			*input.PvolLdevID, *input.MuNumber)
		return snapshot, nil
	}

	saveSvol := snapshot.SvolLdevID

	currentStatus := snapshot.Status
	if snapshot.Status != "PAIR" && snapshot.Status != "PSUS" && snapshot.Status != "PSUE" {
		return nil, log.WriteAndReturnError("Unassign failed: snapshot %d,%d must be in PAIR, PSUS or PSUE status (current: %s)",
			*input.PvolLdevID, *input.MuNumber, snapshot.Status)
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	log.WriteInfo("Unassign| Executing Unassign API for P-VOL %d, MU %d (Detaching S-VOL %d)",
		*input.PvolLdevID, *input.MuNumber, snapshot.SvolLdevID)

	_, err = provObj.UnassignSnapshotVolume(*input.PvolLdevID, *input.MuNumber)
	if err != nil {
		return nil, log.WriteAndReturnError("Unassign failed for P-VOL %d, MU %d: %v",
			*input.PvolLdevID, *input.MuNumber, err)
	}

	snapshot, err = psm.waitForSnapshotStatus(input.PvolLdevID, input.MuNumber, nil, []string{currentStatus})
	if err != nil {
		return nil, err
	}

	// automatically delete svol
	log.WriteInfo("Unassign: Automatically deleting unassigned S-VOl %d", saveSvol)
	err = psm.DeleteLun(saveSvol)
	if err != nil {
		return nil, log.WriteAndReturnError("Unassign failed: deleting S-VOL %d failed after unassign: %+v", saveSvol, err)
	}

	return snapshot, nil
}

// ReconcileSnapshotDelete handles removing the snapshot pair
func (psm *sanStorageManager) ReconcileSnapshotDelete(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if err := psm.validatePvolAndMu(input.PvolLdevID, input.MuNumber); err != nil {
		return nil, err
	}

	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		if isNotFoundError(err) {
			log.WriteInfo("Snapshot %d,%d already gone", *input.PvolLdevID, *input.MuNumber)
			return nil, nil
		} else {
			return nil, log.WriteAndReturnError(err.Error())
		}
	}

	err = checkTIARetention(snapshot)
	if err != nil {
		return nil, log.WriteAndReturnError(err.Error())
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.DeleteSnapshot(*input.PvolLdevID, *input.MuNumber)
	if err != nil {
		return nil, err
	}

	err = psm.waitForSnapshotDeletion(input.PvolLdevID, input.MuNumber)
	return nil, err
}

// reconcileSetSnapshotRetentionPeriod only for TIA pairs
func (psm *sanStorageManager) reconcileSetSnapshotRetentionPeriod(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshot, err := psm.validateTIOperationInput(input)
	if err != nil {
		return nil, err
	}

	if !isTIAdvancedPair(snapshot) {
		return nil, log.WriteAndReturnError("Update retention period can only be performed on TIA pair (P-VOL %d, MU %d)",
			*input.PvolLdevID, *input.MuNumber)
	}

	// User has to make sure the pair is already split because we can't resync right away because of the retention period if originally paired.
	if snapshot.Status != "PSUS" {
		return nil, log.WriteAndReturnError("cannot update retention period: pair %d,%d must be in split (PSUS) status (current: %s)",
			*input.PvolLdevID, *input.MuNumber, snapshot.Status)
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.SetSnapshotRetentionPeriodRequest{
		Parameters: gwymodel.SetSnapshotRetentionPeriodParams{
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	_, err = provObj.SetSnapshotRetentionPeriod(*input.PvolLdevID, *input.MuNumber, request)
	if err != nil {
		return nil, err
	}

	return nil, err
}

// --- Validations ---

func (psm *sanStorageManager) validatePvolAndMu(pvolID *int, mu *int) error {
	log := commonlog.GetLogger()

	if pvolID == nil || mu == nil {
		return log.WriteAndReturnError("P-VOL id and mirror unit id cannot be nil")
	}

	// Based on documentation: LDEV ID must be > 0 and MU must be >= 0
	if *pvolID <= 0 || *mu < 0 {
		return log.WriteAndReturnError("invalid identification: P-VOL id must be > 0 (got %d) and mirror unit id must be >= 0 (got %d)", *pvolID, *mu)
	}

	return nil
}

func (psm *sanStorageManager) validateSnapshotGetMultipleInput(input reconmodel.SnapshotGetMultipleInput) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Pattern Enforcement
	hasPvol := input.PvolLdevID != nil
	hasSvol := input.SvolLdevID != nil
	hasGroup := input.SnapshotGroupName != nil
	hasMu := input.MuNumber != nil

	// Check against the 4 valid combinations from the documentation
	isValidCombination := false

	switch {
	// Combination 1: P-VOL LDEV + Snapshot Group Name
	case hasPvol && hasGroup && !hasSvol && !hasMu:
		isValidCombination = true

	// Combination 2: P-VOL LDEV + MU Number
	case hasPvol && hasMu && !hasSvol && !hasGroup:
		isValidCombination = true

	// Combination 3: Only P-VOL LDEV
	case hasPvol && !hasGroup && !hasSvol && !hasMu:
		isValidCombination = true

	// Combination 4: Only S-VOL LDEV
	case hasSvol && !hasPvol && !hasGroup && !hasMu:
		isValidCombination = true
	}

	if !isValidCombination {
		log.WriteError("TFError| Input: Pvol:%t, Svol:%t, Group:%t, MU:%t)", hasPvol, hasSvol, hasGroup, hasMu)
		return log.WriteAndReturnError("invalid parameter combination. Allowed patterns: [Pvol+Group], [Pvol+MU], [Only Pvol], or [Only Svol]")
	}

	// Individual Value Validation
	if hasGroup && (len(*input.SnapshotGroupName) < 1 || len(*input.SnapshotGroupName) > 32) {
		return log.WriteAndReturnError("snapshotGroupName must be 1-32 characters (received: %d)", len(*input.SnapshotGroupName))
	}
	if hasPvol && *input.PvolLdevID < 0 {
		return log.WriteAndReturnError("P-VOL id must be 0 or greater (received: %d)", *input.PvolLdevID)
	}
	if hasSvol && *input.SvolLdevID < 0 {
		return log.WriteAndReturnError("S-VOL id must be 0 or greater (received: %d)", *input.SvolLdevID)
	}
	if hasMu && (*input.MuNumber < 0) {
		return log.WriteAndReturnError("mirror unit id must be between 0 and 3 (received: %d)", *input.MuNumber)
	}

	log.WriteInfo("Validation successful for standard Get request. Pvol: %v, Svol: %v, Group: %v, MU: %v",
		input.PvolLdevID, input.SvolLdevID, input.SnapshotGroupName, input.MuNumber)
	return nil
}

// --- Validate Create Inputs ---

func (psm *sanStorageManager) validateSnapshotCreateInput(input reconmodel.SnapshotReconcilerInput) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if err := psm.validateGlobalRequiredFields(input); err != nil {
		return err
	}

	needRefresh := false // read has done it
	pool, pvol, svol, err := psm.fetchRequiredResources(input, needRefresh)
	if err != nil {
		return err
	}

	tiaInput := isTIAInput(pvol, pool)

	if tiaInput {
		if err := psm.validateTIAInputs(input, pvol, svol, pool); err != nil {
			return err
		}
	} else {
		if err := psm.validateTIStandardInputs(input, pvol, svol); err != nil {
			return err
		}
	}

	// 3. Final logic-based cross-checks (Clones, Speeds, etc.)
	return psm.validateLogicConstraints(input, tiaInput)
}

// --- Sub-function: Global Required Fields ---
func (psm *sanStorageManager) validateGlobalRequiredFields(input reconmodel.SnapshotReconcilerInput) error {
	log := commonlog.GetLogger()

	if input.SnapshotGroupName == nil || *input.SnapshotGroupName == "" {
		return log.WriteAndReturnError("snapshot_group_name is required (1-32 characters)")
	}
	if input.SnapshotPoolID == nil {
		return log.WriteAndReturnError("snapshot_pool_id is required")
	}
	if input.PvolLdevID == nil {
		return log.WriteAndReturnError("P-VOL id is required")
	}
	return nil
}

// --- Sub-function: TIA Specific Rules ---
func (psm *sanStorageManager) validateTIAInputs(input reconmodel.SnapshotReconcilerInput, pvol *gwymodel.LogicalUnit, svol *gwymodel.LogicalUnit, pool *provmodel.DynamicPool) error {
	log := commonlog.GetLogger()

	log.WriteInfo("Validating Thin Image Advanced specific rules")

	// Documentation: TIA pairs must use HDP pools.
	if pool.PoolType != "HDP" {
		return log.WriteAndReturnError("Thin Image Advanced pairs require an HDP pool (current pool type: %s)", pool.PoolType)
	}

	// Rule: snapshot Pool must match P-VOL Pool for TIA
	if pool.PoolID != pvol.PoolID {
		return log.WriteAndReturnError("Snapshot pool ID (%d) must match P-VOL pool ID (%d) for Thin Image Advanced", pool.PoolID, pvol.PoolID)
	}

	// Rule: isClone must be false for TIA
	if input.IsClone != nil && *input.IsClone {
		return log.WriteAndReturnError("Thin Image Advanced pairs do not support isClone: true")
	}

	// Rule: canCascade must be true for TIA
	if input.CanCascade == nil || !*input.CanCascade {
		return log.WriteAndReturnError("canCascade must be true for Thin Image Advanced")
	}

	// Rule: isDataReductionForceCopy must be true for TIA
	if input.IsDataReductionForceCopy == nil || !*input.IsDataReductionForceCopy {
		return log.WriteAndReturnError("isDataReductionForceCopy must be true for Thin Image Advanced")
	}

	// Rule: Capacity Saving must be ENABLED or CONVERTING (mapped via DataReductionMode)
	mode := strings.ToLower(pvol.DataReductionMode)
	if mode == "" || mode == "disabled" || mode == "none" {
		return log.WriteAndReturnError("P-VOL %d must have capacity saving (compression/deduplication) enabled for Thin Image Advanced (current mode: %s)", pvol.LdevID, pvol.DataReductionMode)
	}

	// Additional Check: Ensure it is a DRS volume
	if !isVolDRS(pvol) {
		return log.WriteAndReturnError("P-VOL %d must have Data Reduction Shared (DRS) enabled for TIA", pvol.LdevID)
	}

	// --- S-VOL Specific Checks for TIA ---
	if svol != nil {
		// Rule: S-VOL Pool must match P-VOL Pool for TIA
		if svol.PoolID != pvol.PoolID {
			return log.WriteAndReturnError("S-VOL pool ID (%d) must match P-VOL pool ID (%d) for Thin Image Advanced", svol.PoolID, pvol.PoolID)
		}

		// Rule: S-VOL must have capacity saving enabled
		svolMode := strings.ToLower(svol.DataReductionMode)
		if svolMode == "" || svolMode == "disabled" || svolMode == "none" {
			return log.WriteAndReturnError("S-VOL %d must have capacity saving enabled for TIA", svol.LdevID)
		}

		// Rule: S-VOL must have Data Reduction Shared (DRS) enabled
		if !isVolDRS(svol) {
			return log.WriteAndReturnError("secondary volume (S-VOL) %d must have Data Reduction Shared (DRS) enabled for TIA", svol.LdevID)
		}
	}

	return nil
}

// --- Sub-function: TI Standard Specific Rules ---
func (psm *sanStorageManager) validateTIStandardInputs(input reconmodel.SnapshotReconcilerInput, pvol *gwymodel.LogicalUnit, svol *gwymodel.LogicalUnit) error {
	log := commonlog.GetLogger()

	log.WriteInfo("Validating Thin Image Standard specific rules")

	// got: Program product is not installed, "SSB2": "9010", "SSB1": "2E21"
	model, isVspB := checkIsVspBSeries(psm.storageSetting.Serial)
	if isVspB {
		errmsg := fmt.Sprintf("TI standard is not supported on VSP B series. Current model: %s", model)
		return log.WriteAndReturnError("%v", errmsg)
	}

	// Standard doesn't have the strict "Always True" DRS requirements
	// but we check if ForceCopy is missing when capacity saving happens to be on
	hasCapSaving := pvol.DataReductionStatus != "DISABLED" && pvol.DataReductionStatus != ""
	if hasCapSaving && (input.IsDataReductionForceCopy == nil || !*input.IsDataReductionForceCopy) {
		return log.WriteAndReturnError("isDataReductionForceCopy must be true when P-VOL capacity saving is enabled")
	}

	// 2. S-VOL Specific Validations for Standard TI
	if svol != nil {
		isClone := input.IsClone != nil && *input.IsClone
		canCascade := input.CanCascade != nil && *input.CanCascade

		// Documentation: If isClone or canCascade is true, S-VOL must be a DP volume (HDP)
		// Usually, we check this by ensuring the S-VOL belongs to an HDP pool type
		if isClone || canCascade {
			// No need to call pool api
			// // Fetch S-VOL pool to verify it's a DP volume
			// provObj, _ := psm.getProvisionerManager()
			// svolPool, err := provObj.GetDynamicPoolById(svol.PoolID)

			// if err == nil && svolPool != nil && svolPool.PoolType != "HDP" {
			// 	return log.WriteAndReturnError("S-VOL %d must be a DP volume (HDP pool) when isClone or canCascade is specified", svol.LdevID)
			// }

			// just check svol attributes
			if !isVolHDP(svol) {
				return log.WriteAndReturnError("S-VOL %d must be a DP volume (HDP pool) when isClone or canCascade is specified", svol.LdevID)
			}
		}
	}

	return nil
}

// --- Sub-function: Logic & Dependencies (Clones/Retention) ---
func (psm *sanStorageManager) validateLogicConstraints(input reconmodel.SnapshotReconcilerInput, isTIA bool) error {
	log := commonlog.GetLogger()

	isClone := input.IsClone != nil && *input.IsClone
	autoSplit := input.AutoSplit != nil && *input.AutoSplit
	clonesAuto := input.ClonesAutomation != nil && *input.ClonesAutomation

	// Rule: Mutual Exclusivity
	if isClone && autoSplit {
		return log.WriteAndReturnError("cannot specify both isClone and autoSplit")
	}

	// Rule: Clone Requirements
	if isClone {
		if input.SvolLdevID == nil {
			return log.WriteAndReturnError("svolLdevId is required when isClone is true")
		}
		if input.CanCascade == nil || !*input.CanCascade {
			return log.WriteAndReturnError("canCascade must be true when isClone is true")
		}
	}

	// FIX for FAIL: "clonesAutomation without isClone"
	if clonesAuto && !isClone {
		return log.WriteAndReturnError("clonesAutomation can only be specified when isClone is true")
	}

	// Rule: Retention Period
	if input.RetentionPeriod != nil {
		if isTIA {
			if !autoSplit {
				return log.WriteAndReturnError("retentionPeriod requires autoSplit to be true")
			}
			// Doc: range check 1 to 12288
			if *input.RetentionPeriod < 1 || *input.RetentionPeriod > 12288 {
				return log.WriteAndReturnError("retentionPeriod must be between 1 and 12288")
			}
		}
	}

	// Rule: Copy Speed
	if input.CopySpeed != nil {
		// Rule: Requires both isClone AND ClonesAutomation
		if !isClone || !clonesAuto {
			return log.WriteAndReturnError("copySpeed requires both isClone and clonesAutomation")
		}
		speed := strings.ToLower(*input.CopySpeed)
		validSpeeds := []string{"slower", "medium", "faster"}
		if !slices.Contains(validSpeeds, speed) {
			return log.WriteAndReturnError("invalid copySpeed: %s", speed)
		}
	}

	// FIX for FAIL: "negative muNumber"
	if input.MuNumber != nil {
		if *input.MuNumber < 0 || *input.MuNumber > 1023 {
			return log.WriteAndReturnError("muNumber must be between 0 and 1023")
		}
	}

	return nil
}

// --- Validate Existing Snapshot ---
func (psm *sanStorageManager) validateExistingSnapshotBackendFromInput(input reconmodel.SnapshotReconcilerInput, backend *gwymodel.Snapshot) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var mismatches []string

	// TIA is typically indicated by IsRedirectOnWrite (ROW) in the backend
	snapshotType := "TI Standard"
	if backend.IsRedirectOnWrite {
		snapshotType = "TI Advanced (TIA)"
	}
	log.WriteInfo("Verifying input against backend for %s (P-VOL %d, MU %d)", snapshotType, backend.PvolLdevID, backend.MuNumber)

	if input.SnapshotGroupName != nil && *input.SnapshotGroupName != backend.SnapshotGroupName {
		mismatches = append(mismatches, fmt.Sprintf("SnapshotGroupName (input: %s, backend: %s)", *input.SnapshotGroupName, backend.SnapshotGroupName))
	}

	if input.SnapshotPoolID != nil && *input.SnapshotPoolID != backend.SnapshotPoolID {
		mismatches = append(mismatches, fmt.Sprintf("SnapshotPoolID (input: %d, backend: %d)", *input.SnapshotPoolID, backend.SnapshotPoolID))
	}

	// not for attach/unattach svol during create
	if *input.Action == "read" {
		// S-VOL Comparison
		if input.SvolLdevID != nil {
			if *input.SvolLdevID != backend.SvolLdevID {
				mismatches = append(mismatches, fmt.Sprintf("S-VOL (input: %d, backend: %d)", *input.SvolLdevID, backend.SvolLdevID))
			}
		}
	}

	if input.IsConsistencyGroup != nil && *input.IsConsistencyGroup != backend.IsConsistencyGroup {
		mismatches = append(mismatches, fmt.Sprintf("IsConsistencyGroup (input: %t, backend: %t)", *input.IsConsistencyGroup, backend.IsConsistencyGroup))
	}

	if input.IsClone != nil && *input.IsClone != backend.IsClone {
		mismatches = append(mismatches, fmt.Sprintf("IsClone (input: %t, backend: %t)", *input.IsClone, backend.IsClone))
	}

	if input.CanCascade != nil && *input.CanCascade != backend.CanCascade {
		mismatches = append(mismatches, fmt.Sprintf("CanCascade (input: %t, backend: %t)", *input.CanCascade, backend.CanCascade))
	}

	// If user inputs a TIA-only attribute or logic suggests TIA, but backend is ROW=false
	if backend.IsRedirectOnWrite {
		// TIA should NOT be a clone
		if backend.IsClone {
			mismatches = append(mismatches, "Type Mismatch: Backend is RedirectOnWrite (TIA) but also marked as Clone")
		}
	}

	if len(mismatches) > 0 {
		errorDetail := strings.Join(mismatches, "; ")
		return log.WriteAndReturnError("Snapshot configuration mismatch: %s", errorDetail)
	}

	log.WriteInfo("Backend verification successful for %s: %s", snapshotType, backend.SnapshotID)
	return nil
}

// --- Validate TI Operation (like split, resync, ...) Inputs ---
func (psm *sanStorageManager) validateTIOperationInput(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	action := *input.Action

	// --- Basic ID Requirements ---
	if input.PvolLdevID == nil || input.MuNumber == nil {
		return nil, log.WriteAndReturnError("P-VOL id and mirror unit id are required for operation input: %s", action)
	}

	// Fetch the current snapshot/pair metadata
	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		return nil, log.WriteAndReturnError("failed to fetch snapshot metadata for P-VOL %d, MU %d: %v",
			*input.PvolLdevID, *input.MuNumber, err)
	}

	// check if floating snapshot
	if snapshot.SvolLdevID == 0 && (*input.Action == "resync" || *input.Action == "restore") {
		return nil, log.WriteAndReturnError("invalid operation: cannot perform action '%s' on a floating snapshot (P-VOL: %d, MU: %d). An S-VOL must be assigned first",
			*input.Action, *input.PvolLdevID, *input.MuNumber)
	}

	// just log if TIA or TI standard
	isTIAdvancedPair(snapshot)

	// REQUIREMENT: Individual pair operations are restricted for CTG-enabled groups
	// If IsConsistencyGroup is true, you must perform actions on the whole group, not the LDEV.
	if snapshot.IsConsistencyGroup {
		return nil, log.WriteAndReturnError("operation '%s' not allowed: snapshot %d,%d belongs to a group in CTG mode. Operations must be performed on the snapshot group level",
			action, *input.PvolLdevID, *input.MuNumber)
	}

	// REQUIREMENT: Check for "Busy" status (Processing Status: P)
	if snapshot.PvolProcessingStatus == "P" {
		return nil, log.WriteAndReturnError("snapshot %d,%d is currently busy (Processing Status: P). Wait and try again",
			*input.PvolLdevID, *input.MuNumber)
	}

	log.WriteInfo("Validation successful for operation input '%s' on P-VOL %d, MU %d", action, *input.PvolLdevID, *input.MuNumber)
	return snapshot, nil
}

func (psm *sanStorageManager) validateAssignVolumeInput(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pvolID := *input.PvolLdevID
	muNum := *input.MuNumber
	svolID := *input.SvolLdevID

	if input.SvolLdevID == nil {
		return nil, log.WriteAndReturnError("S-VOL id is required for assignment")
	}

	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		return nil, log.WriteAndReturnError("Assign failed: snapshot data for %d,%d not found", pvolID, muNum)
	}

	if snapshot.SvolLdevID != 0 {
		return nil, log.WriteAndReturnError(
			"Assign failed: snapshot %d,%d already has S-VOL %d assigned.",
			pvolID, muNum, snapshot.SvolLdevID,
		)
	}

	needRefresh := false // read has done it
	_, pvol, svol, err := psm.fetchRequiredResources(input, needRefresh)
	if err != nil {
		return nil, err
	}

	if svol == nil || svol.EmulationType == "NOT DEFINED" {
		return nil, log.WriteAndReturnError("Assign failed: target S-VOL %d not found", svolID)
	}

	if svol.Status != "NML" {
		return nil, fmt.Errorf("S-VOL %d is in status %s; it must be NML to be assigned", svolID, svol.Status)
	}

	// Check Capacity
	if svol.ByteFormatCapacity != pvol.ByteFormatCapacity {
		return nil, fmt.Errorf("capacity mismatch: P-VOL %d is %v, S-VOL %d is %v",
			pvolID, pvol.ByteFormatCapacity, svolID, svol.ByteFormatCapacity)
	}

	// TIA / DRS Validation
	if snapshot.IsRedirectOnWrite {
		log.WriteInfo("Assign| TIA detected. Validating pool and DRS for S-VOL %d", svolID)

		if snapshot.SnapshotPoolID != svol.PoolID {
			return nil, log.WriteAndReturnError(
				"Assign failed: Pool mismatch. S-VOL %d is in Pool %d, but Snapshot requires Pool %d",
				svolID, svol.PoolID, snapshot.SnapshotPoolID,
			)
		}

		if !isVolDRS(svol) {
			return nil, log.WriteAndReturnError(
				"Assign failed: S-VOL %d lacks the 'DRS' attribute required for Thin Image Advanced.",
				svolID,
			)
		}
	}

	log.WriteInfo("Validation successful for assigning S-VOL %d to Snapshot %d,%d", svolID, pvolID, muNum)
	return snapshot, nil
}

func (psm *sanStorageManager) validateUnassignVolumeInput(input reconmodel.SnapshotReconcilerInput) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	pvolID := *input.PvolLdevID
	muNum := *input.MuNumber

	snapshot, err := psm.getSnapshot(input.PvolLdevID, input.MuNumber, false)
	if err != nil {
		return nil, log.WriteAndReturnError("Unassign| Snapshot %d,%d not found; likely already purged.", pvolID, muNum)
	}

	return snapshot, nil
}

// checkTIARetention verifies if a snapshot is TIA and if it is currently locked.
func checkTIARetention(snapshot *gwymodel.Snapshot) error {
	if isTIAdvancedPair(snapshot) {
		if snapshot.RetentionPeriod > 0 {
			return fmt.Errorf(
				"Operation blocked: Thin Image Advanced (TIA) snapshot has an active retention lock "+
					"(%d hours remaining); storage will not allow modification or deletion until the lock expires",
				snapshot.RetentionPeriod,
			)
		}
	}
	return nil
}

// --- Wait funcs ---

func (psm *sanStorageManager) waitForSnapshotStatus(pvolID *int, mu *int, tiTargetStatuses []string, tiaTargetStatuses []string) (*gwymodel.Snapshot, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var lastSnap *gwymodel.Snapshot

	if err := psm.validatePvolAndMu(pvolID, mu); err != nil {
		return nil, err
	}

	// // Configure timing: 5 mins total (30 retries * 10s)
	// maxRetries := 30
	// delay := 10 * time.Second

	retryCfg := utils.GetRetryConfig()
	maxRetries := retryCfg.MaxRetries
	delay := time.Duration(retryCfg.Delay) * time.Second

	for i := 0; i < maxRetries; i++ {
		snap, err := psm.getSnapshot(pvolID, mu, true)
		if err != nil {
			log.WriteError("TFError| failed to fetch snapshot metadata for P-VOL %d, MU %d: %v", *pvolID, *mu, err)
			return nil, err
		}

		lastSnap = snap

		if snap != nil {
			var currentTargets []string
			var isTIA bool

			// 1. Identify Snapshot Architecture
			if snap.IsRedirectOnWrite {
				isTIA = true
				currentTargets = tiaTargetStatuses
				log.WriteDebug("Polling| Detected Thin Image Advanced (TIA) snapshot")
			} else {
				isTIA = false
				currentTargets = tiTargetStatuses
				log.WriteDebug("Polling| Detected Thin Image (TI) snapshot")
			}

			// 2. Handle Nil Target List
			// If the user passed nil for this specific type, we exit early and return current state
			if currentTargets == nil {
				typeStr := "TI Standard"
				if isTIA {
					typeStr = "TIA"
				}
				log.WriteInfo("Polling| No target statuses provided for %s. Returning current snapshot state immediately.", typeStr)
				return snap, nil
			}

			// 3. Check against Target Statuses
			for _, target := range currentTargets {
				if snap.Status == target {
					log.WriteInfo("Polling| Snapshot %d,%d reached target status: %s", *pvolID, *mu, target)
					return snap, nil
				}
			}

			log.WriteDebug("Polling| Snapshot %d,%d status: %s. Target list: %v (Attempt %d/%d)",
				*pvolID, *mu, snap.Status, currentTargets, i+1, maxRetries)
		} else {
			log.WriteWarn("Polling| Received nil snapshot object from gateway for P-VOL %d, MU %d", *pvolID, *mu)
		}

		time.Sleep(delay)
	}

	// Final target list for error reporting
	finalTargets := tiTargetStatuses
	if lastSnap != nil && lastSnap.IsRedirectOnWrite {
		finalTargets = tiaTargetStatuses
	}

	return lastSnap, log.WriteAndReturnError("timeout: snapshot %d,%d did not reach %v within 5 minutes. Current status: %s",
		*pvolID, *mu, finalTargets, lastSnap.Status)
}

func (psm *sanStorageManager) waitForSnapshotDeletion(pvolID *int, mu *int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if err := psm.validatePvolAndMu(pvolID, mu); err != nil {
		return err
	}

	// maxRetries := 20
	// delay := 5 * time.Second
	retryCfg := utils.GetRetryConfig()
	maxRetries := retryCfg.MaxRetries
	delay := time.Duration(retryCfg.Delay) * time.Second

	for i := 0; i < maxRetries; i++ {
		_, err := psm.getSnapshot(pvolID, mu, true)
		if err != nil {
			if isNotFoundError(err) {
				log.WriteInfo("Snapshot %d,%d successfully deleted", *pvolID, *mu)
				return nil
			}
		}
		time.Sleep(delay)
	}
	return log.WriteAndReturnError("timeout waiting for snapshot %d,%d deletion", *pvolID, *mu)
}

// --- Others ---

func isTIAdvancedPair(snapshot *gwymodel.Snapshot) bool {
	log := commonlog.GetLogger()

	if snapshot.IsRedirectOnWrite {
		log.WriteInfo("Snapshot %d,%d is a TIA (Redirect-on-Write) pair", snapshot.PvolLdevID, snapshot.MuNumber)
		return true
	} else {
		log.WriteInfo("Snapshot %d,%d is a standard TI pair", snapshot.PvolLdevID, snapshot.MuNumber)
		return false
	}
}

func hasVolAttribute(lun *gwymodel.LogicalUnit, attr string) bool {
	log := commonlog.GetLogger()

	found := slices.ContainsFunc(lun.Attributes, func(s string) bool {
		exist := strings.EqualFold(s, attr)
		if exist {
			log.WriteDebug("LUN %d has attribute %s", lun.LdevID, attr)
		}
		return exist
	})

	return found
}

func isVolDRS(lun *gwymodel.LogicalUnit) bool {
	return hasVolAttribute(lun, "DRS")
}

func isVolHDP(lun *gwymodel.LogicalUnit) bool {
	return hasVolAttribute(lun, "HDP")
}

func isTIAInput(pvol *gwymodel.LogicalUnit, pool *provmodel.DynamicPool) bool {
	log := commonlog.GetLogger()

	if pvol != nil && pool != nil {
		return pool.PoolType == "HDP" && isVolDRS(pvol)
	}
	if pvol != nil && pool == nil {
		log.WriteDebug("No pool input given")
		return isVolDRS(pvol)
	}
	log.WriteDebug("No pool and pvol input given")
	return false
}

func checkIsHighEndSeries(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)

	// Returns true if it's a 5000 series, B20, B85, or E series.
	isSupported := strings.Contains(m, "VSP 5") ||
		strings.Contains(m, "VSP ONE B2") ||
		strings.Contains(m, "VSP ONE B85") ||
		strings.Contains(m, "VSP E")

	return model, isSupported
}

func checkIsVspESeries(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)
	return model, strings.Contains(m, "VSP E")
}

func checkIsVsp5000Series(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)
	return model, strings.Contains(m, "VSP 5")
}

func checkIsVspBSeries(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)
	return model, strings.Contains(m, "VSP ONE B2") ||
		strings.Contains(m, "VSP ONE B85")
}

func checkIsVspB20Series(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)
	return model, strings.Contains(m, "VSP ONE B2")
}

func checkIsVspB85(storageSerialNumber int) (string, bool) {
	model := getStorageModel(storageSerialNumber)
	m := strings.ToUpper(model)
	return model, strings.Contains(m, "VSP ONE B85")
}

func parseSnapshotResID(resID string) (*int, *int, error) {
	log := commonlog.GetLogger()

	// 1. Split the string by comma
	parts := strings.Split(resID, ",")
	if len(parts) != 2 {
		return nil, nil, log.WriteAndReturnError("invalid state id format: expected 'pvolId,muId', got '%s'", resID)
	}

	// 2. Parse PvolLdevID
	pvolID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, nil, log.WriteAndReturnError("failed to parse pvolId from state Id '%s': %v", resID, err)
	}

	// 3. Parse MuNumber
	muID, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, nil, log.WriteAndReturnError("failed to parse mirror unit id from state Id '%s': %v", resID, err)
	}

	log.WriteDebug("Successfully parsed state Id %s into PvolID: %d, MuID: %d", resID, pvolID, muID)
	return &pvolID, &muID, nil
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "does not exist") ||
		strings.Contains(errStr, "404")
}

func getStorageModel(storageSerialNumber int) string {
	storageSettingsAndInfo, _ := cache.ReadFromSanCache(strconv.Itoa(storageSerialNumber))
	model := ""
	if storageSettingsAndInfo.Info != nil && storageSettingsAndInfo.Info.Model != "" {
		model = storageSettingsAndInfo.Info.Model
	}
	return model
}

func (psm *sanStorageManager) getProvisionerManager() (provmanager.SanStorageManager, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	setting := provmodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provimpl.NewEx(setting)
	if err != nil {
		log.WriteError("failed to get provisioner manager: %v", err)
		return nil, log.WriteAndReturnError("failed to get provisioner manager: %w", err)
	}

	return provObj, nil
}
