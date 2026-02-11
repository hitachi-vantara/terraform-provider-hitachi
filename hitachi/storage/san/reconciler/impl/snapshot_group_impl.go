package sanstorage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

// --- Datasources ---

func (psm *sanStorageManager) ReconcileGetSnapshotGroup(snapshotGroupName string) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	detailType := "retention"
	params := gwymodel.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}

	resp, err := provObj.GetSnapshotGroup(snapshotGroupName, params)
	if err != nil {
		log.WriteError("TFError| failed to get snapshot group %s: %v", snapshotGroupName, err)
		return nil, err
	}

	return resp, nil
}

func (psm *sanStorageManager) ReconcileGetMultipleSnapshotGroups(includePairs bool) (*gwymodel.SnapshotGroupListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// 1. Get the base list of groups
	// If includePairs is false, we don't need any special DetailInfoType here
	log.WriteInfo("Fetching base snapshot group list (includePairs=%t)", includePairs)
	baseList, err := provObj.GetSnapshotGroups(gwymodel.GetSnapshotGroupsParams{})
	if err != nil {
		return nil, err
	}

	if baseList == nil || len(baseList.Data) == 0 {
		return &gwymodel.SnapshotGroupListResponse{Data: []gwymodel.SnapshotGroup{}}, nil
	}

	// 2. If we don't need pair details, return the base list immediately
	if !includePairs {
		log.WriteInfo("Returning base list only (skipping parallel detail calls)")
		return baseList, nil
	}

	// 3. If includePairs is true, proceed to fetch details in parallel
	log.WriteInfo("includePairs is true: performing parallel calls to fetch member data")

	var mu sync.Mutex
	var results []gwymodel.SnapshotGroup
	var firstErr error

	detailType := "retention"
	params := gwymodel.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}

	errs := utils.RunConcurrentOperations("ReconcileGetMultipleSnapshotGroups", baseList.Data,
		func(sg gwymodel.SnapshotGroup, idx int) error {
			resp, err := provObj.GetSnapshotGroup(sg.SnapshotGroupName, params)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return err
			}

			if resp != nil {
				mu.Lock()
				results = append(results, *resp)
				mu.Unlock()
			}
			return nil
		})

	for _, e := range errs {
		if e != nil && firstErr == nil {
			firstErr = e
		}
	}

	// Sort the detailed results
	sort.Slice(results, func(i, j int) bool {
		return results[i].SnapshotGroupName < results[j].SnapshotGroupName
	})

	return &gwymodel.SnapshotGroupListResponse{Data: results}, firstErr
}

// --- Main Reconciler Dispatcher ---

// ReconcileSnapshotGroupVFamily is the entry point that decides between standard and vFamily logic.
func (psm *sanStorageManager) ReconcileSnapshotGroupVFamily(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.Action == nil || *input.Action == "" {
		input.Action = utils.Ptr("read")
	}

	action := strings.ToLower(*input.Action)

	log.WriteInfo("Starting snapshot group reconciliation for operation: %s", action)

	if action == "read" {
		snap, vclones, err := psm.reconcileReadExistingSnapshotGroup(input)
		return snap, vclones, err
	}

	if action == "vclone" || action == "vrestore" {
		return psm.reconcileSnapshotGroupVFamilyDispatch(input)
	}

	// Default to standard reconciler
	res, err := psm.ReconcileSnapshotGroup(input)
	return res, nil, err
}

func (psm *sanStorageManager) reconcileSnapshotGroupVFamilyDispatch(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	model, isVspB := checkIsVspBSeries(psm.storageSetting.Serial)
	if !isVspB {
		return nil, nil, log.WriteAndReturnError("vclone or vrestore is only supported in B20/B85 series. Current model: %s", model)
	}

	action := strings.ToLower(*input.Action)

	switch action {
	case "vclone":
		group, family, err := psm.reconcileSnapshotGroupVClone(input)
		return group, family, err
	case "vrestore":
		group, family, err := psm.reconcileSnapshotGroupVRestore(input)
		return group, family, err
	default:
		return nil, nil, fmt.Errorf("unsupported vfamily action: %s", action)
	}
}

func (psm *sanStorageManager) ReconcileSnapshotGroup(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if input.SnapshotGroupName == nil || *input.SnapshotGroupName == "" {
		return nil, fmt.Errorf("snapshotGroupName is mandatory for snapshot group operations")
	}

	action := "read"
	if input.Action != nil {
		action = strings.ToLower(*input.Action)
	}

	log.WriteInfo("Starting snapshot group reconciliation for ID: %s, Action: %s", *input.SnapshotGroupName, action)

	switch action {
	case "split":
		return psm.reconcileSnapshotGroupSplit(input)
	case "resync":
		return psm.reconcileSnapshotGroupResync(input)
	case "restore":
		return psm.reconcileSnapshotGroupRestore(input)
	case "clone":
		return psm.reconcileSnapshotGroupClone(input)
	case "delete":
		return psm.reconcileSnapshotGroupDelete(input)
	case "update_retention":
		return psm.reconcileSetSnapshotGroupRetentionPeriod(input)
	default:
		return nil, fmt.Errorf("unsupported snapshot group action: %s", action)
	}
}

// --- Action Handlers ---

func (psm *sanStorageManager) reconcileReadExistingSnapshotGroup(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	groupName := *input.SnapshotGroupName
	var familyList []gwymodel.SnapshotFamily

	group, err := psm.ReconcileGetSnapshotGroup(groupName)
	if err != nil {
		if isNotFoundError(err) {
			log.WriteInfo("TFInfo| Snapshot Group %s not found.", groupName)
			return nil, nil, nil
		}
		return nil, nil, log.WriteAndReturnError("TFError| System error fetching snapshot group %s: %v", groupName, err)
	}

	model, isVspB := checkIsVspBSeries(psm.storageSetting.Serial)
	if !isVspB {
		log.WriteDebug("Get vclone family is only supported in B20/B85 series. Current model: %s", model)
		return group, nil, nil
	}

	// Iterate through members to collect Family Metadata
	for _, snap := range group.Snapshots {
		// Only attempt family lookup for Thin Image Advanced pairs
		if isTIAdvancedPair(&snap) {
			log.WriteDebug("TFDebug| Fetching family metadata for member: SVOL %d", snap.SvolLdevID)

			// Replicating single snapshot pattern: Get vclone family info
			vFamily, fErr := psm.getVcloneFromFamily(&snap.PvolLdevID, &snap.MuNumber, &snap.SvolLdevID, true)
			if fErr != nil {
				log.WriteDebug("TFDebug| Could not fetch family info for SVOL %d: %v", snap.SvolLdevID, fErr)
				continue
			}

			if vFamily != nil {
				familyList = append(familyList, *vFamily)
			}
		}
	}

	log.WriteInfo("TFInfo| Reconciled Snapshot Group %s with %d family records found.", groupName, len(familyList))
	return group, familyList, nil
}

// reconcileSnapshotgGroupSplit handles splitting an existing snapshot pairs in a snapshot group
// Note: For Thin Image Advanced pairs, if both of the following conditions are met, operations in units of snapshot groups might fail.
// If the following conditions are met, perform operations on one pair at a time by using the API request for splitting snapshot.
//
//	The snapshot group was not created in CTG mode.
//	The snapshot group contains two or more pairs that have the same volume as the P-VOL.
func (psm *sanStorageManager) reconcileSnapshotGroupSplit(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	request := gwymodel.SplitSnapshotRequest{
		Parameters: gwymodel.SplitSnapshotParams{
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.SplitSnapshotGroup(snapshotGroupName, request)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PSUS", "PFUS"}
	tiaTargetStatuses := []string{"PSUS", "PFUS", "PFUL"}
	return psm.waitForSnapshotGroupStatus(snapshotGroupName, tiTargetStatuses, tiaTargetStatuses)
}

// reconcileSnapshotGroupResync handles resynchronizing an existing snapshot pairs in a snapshot group
// Note: If the snapshot group includes at least one pair for which a snapshot data retention period (retentionPeriod) is set,
// you will not be able to resynchronize pairs in units of snapshot groups.
// In this case, wait until the snapshot data retention periods for all pairs end, and then perform the operation.
// Alternatively, individually resynchronize pairs for which no snapshot data retention period is set.
func (psm *sanStorageManager) reconcileSnapshotGroupResync(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	err := psm.validateSnapshotGroupRetention(snapshotGroupName)
	if err != nil {
		return nil, err
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.ResyncSnapshotRequest{
		Parameters: gwymodel.ResyncSnapshotParams{
			AutoSplit:       input.AutoSplit,
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	_, err = provObj.ResyncSnapshotGroup(snapshotGroupName, request)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PAIR"}
	tiaTargetStatuses := []string{"PAIR"}
	if input.AutoSplit != nil && *input.AutoSplit == true {
		tiTargetStatuses = []string{"PSUS", "PFUS", "PFUL"}
		tiaTargetStatuses = []string{"PSUS", "PFUS", "PFUL"}
	}

	return psm.waitForSnapshotGroupStatus(snapshotGroupName, tiTargetStatuses, tiaTargetStatuses)
}

// reconcileSnapshotGroupRestore handles restoring P-VOL data from an S-VOL
// For Thin Image Advanced pairs, autoSplit attribute is ignored even if it is specified.
func (psm *sanStorageManager) reconcileSnapshotGroupRestore(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.RestoreSnapshotRequest{
		Parameters: gwymodel.RestoreSnapshotParams{
			AutoSplit: input.AutoSplit,
		},
	}

	_, err = provObj.RestoreSnapshotGroup(snapshotGroupName, request)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PAIR"}
	tiaTargetStatuses := []string{"PAIR", "PSUS"} // TIA does not go to PAIR

	// If it's TIA, AutoSplit is ignored
	if input.AutoSplit != nil && *input.AutoSplit == true {
		tiTargetStatuses = []string{"PSUS"}
		// Note: We don't change tiaTargetStatuses here because the TIA ignores the flag
		// Tested it, it's not ignored
	}

	return psm.waitForSnapshotGroupStatus(snapshotGroupName, tiTargetStatuses, tiaTargetStatuses)
}

// reconcileSnapshotClone handles cloning a snapshot to a new volume (in a snapshot group)
// Only for TI Standard pairs. Not for Thin Image Advanced pairs
func (psm *sanStorageManager) reconcileSnapshotGroupClone(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.CloneSnapshotRequest{
		Parameters: gwymodel.CloneSnapshotParams{
			CopySpeed: input.CopySpeed,
		},
	}

	_, err = provObj.CloneSnapshotGroup(snapshotGroupName, request)
	if err != nil {
		return nil, err
	}

	// TODO: No need to wait for all pairs, maybe just check that none are in PAIR/COPY?
	// at this point, TI snapshot is gone
	err = psm.waitForSnapshotGroupDeletion(snapshotGroupName)
	return nil, err
}

// reconcileSnapshotGroupDelete
// Note: If the snapshot group includes at least one pair for which a snapshot data retention period (retentionPeriod) is set,
// you will not be able to delete pairs in units of snapshot groups.
// In this case, wait until the snapshot data retention periods for all pairs end, and then perform the operation.
// Alternatively, individually delete pairs for which no snapshot data retention period is set.
func (psm *sanStorageManager) reconcileSnapshotGroupDelete(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	err := psm.validateSnapshotGroupRetention(snapshotGroupName)
	if err != nil {
		return nil, err
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	_, err = provObj.DeleteSnapshotGroup(snapshotGroupName)
	if err != nil {
		return nil, err
	}

	// TODO: check this
	err = psm.waitForSnapshotGroupDeletion(snapshotGroupName)
	return nil, err
}

// reconcileSetSnapshotGroupRetentionPeriod only for TIA pairs
func (psm *sanStorageManager) reconcileSetSnapshotGroupRetentionPeriod(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	snapshotGroupName := *input.SnapshotGroupName

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	request := gwymodel.SetSnapshotRetentionPeriodRequest{
		Parameters: gwymodel.SetSnapshotRetentionPeriodParams{
			RetentionPeriod: input.RetentionPeriod,
		},
	}

	_, err = provObj.SetSnapshotGroupRetentionPeriod(snapshotGroupName, request)
	if err != nil {
		return nil, err
	}

	tiTargetStatuses := []string{"PSUS", "PFUS"}
	tiaTargetStatuses := []string{"PSUS", "PFUS", "PFUL"}
	return psm.waitForSnapshotGroupStatus(snapshotGroupName, tiTargetStatuses, tiaTargetStatuses)
}

// --- Helper Functions ---

func (psm *sanStorageManager) validateSnapshotGroupRetention(groupName string) error {
	log := commonlog.GetLogger()
	provObj, _ := psm.getProvisionerManager()

	params := gwymodel.GetSnapshotGroupsParams{DetailInfoType: utils.Ptr("retention")}

	group, err := provObj.GetSnapshotGroup(groupName, params)
	if err != nil {
		return fmt.Errorf("failed to retrieve snapshot group %s for retention check: %w", groupName, err)
	}

	log.WriteInfo("RetentionPreVerify| Checking %d snapshots in group: %s", len(group.Snapshots), groupName)

	for _, snap := range group.Snapshots {
		if err := checkTIARetention(&snap); err != nil {
			log.WriteError("RetentionPreVerify| Locked Member Found in group %s: P-VOL %d",
				groupName, snap.PvolLdevID)
			return fmt.Errorf("Snapshot group '%s' member check failed: %w", groupName, err)
		}
	}

	log.WriteInfo("RetentionPreVerify| All members in group %s passed retention check.", groupName)
	return nil
}

func (psm *sanStorageManager) waitForSnapshotGroupStatus(snapshotGroupName string, tiTargetStatuses []string, tiaTargetStatuses []string) (*gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, err
	}

	// 1. Use dynamic retry configuration
	retryCfg := utils.GetRetryConfig()

	detailType := "retention"
	params := gwymodel.GetSnapshotGroupsParams{
		DetailInfoType: &detailType,
	}

	// 2. Main Polling Loop
	for attempt := 1; attempt <= retryCfg.MaxRetries; attempt++ {
		group, err := provObj.GetSnapshotGroup(snapshotGroupName, params)

		// Handle case where group is missing or API fails
		if err != nil || group == nil {
			log.WriteDebug("Attempt %d/%d: Snapshot group %s not found or error: %v",
				attempt, retryCfg.MaxRetries, snapshotGroupName, err)

			if attempt < retryCfg.MaxRetries {
				time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
				continue
			}
			return nil, fmt.Errorf("snapshot group %s not found after %d attempts: %w", snapshotGroupName, retryCfg.MaxRetries, err)
		}

		// If the group is empty, there is nothing to wait for
		if len(group.Snapshots) == 0 {
			log.WriteWarn("Snapshot group %s has no member pairs to track", snapshotGroupName)
			return group, nil
		}

		// 3. Check statuses of all members (Local in-memory check)
		allMatch := true
		for _, pair := range group.Snapshots {
			var currentTargets []string
			pairTypeStr := "TI Standard"

			if isTIAdvancedPair(&pair) {
				currentTargets = tiaTargetStatuses
				pairTypeStr = "TIA"
			} else {
				currentTargets = tiTargetStatuses
			}

			match := false
			for _, target := range currentTargets {
				if pair.Status == target {
					match = true
					break
				}
			}

			if !match {
				allMatch = false
				log.WriteDebug("Group %s: Member P-VOL %d (%s) is in status %s; waiting for %v",
					snapshotGroupName, pair.PvolLdevID, pairTypeStr, pair.Status, currentTargets)
				break
			}
		}

		if allMatch {
			log.WriteInfo("All members of group %s reached their respective target statuses", snapshotGroupName)
			return group, nil
		}

		// 4. Dynamic Delay
		if attempt < retryCfg.MaxRetries {
			time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
		}
	}

	return nil, fmt.Errorf("timeout waiting for snapshot group %s. TI targets: %v, TIA targets: %v",
		snapshotGroupName, tiTargetStatuses, tiaTargetStatuses)
}

func (psm *sanStorageManager) waitForSnapshotGroupDeletion(snapshotGroupName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return fmt.Errorf("failed to get provisioner manager: %w", err)
	}

	// 1. Use your dynamic retry configuration
	retryCfg := utils.GetRetryConfig()
	params := gwymodel.GetSnapshotGroupsParams{}

	log.WriteInfo("WaitDelete| Monitoring deletion of group %s (Max Retries: %d)",
		snapshotGroupName, retryCfg.MaxRetries)

	// 2. Polling loop
	for attempt := 1; attempt <= retryCfg.MaxRetries; attempt++ {
		_, err := provObj.GetSnapshotGroup(snapshotGroupName, params)

		if err != nil {
			// Success condition: The group is no longer found
			if isNotFoundError(err) {
				log.WriteInfo("WaitDelete| Snapshot group %s successfully deleted.", snapshotGroupName)
				return nil
			}
			// Actual API error
			return fmt.Errorf("error while polling for deletion: %w", err)
		}

		log.WriteDebug("WaitDelete| Group %s still exists, retrying... (%d/%d)",
			snapshotGroupName, attempt, retryCfg.MaxRetries)

		// 3. Use dynamic delay
		if attempt < retryCfg.MaxRetries {
			time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
		}
	}

	return fmt.Errorf("timeout: snapshot group %s still exists after %d attempts",
		snapshotGroupName, retryCfg.MaxRetries)
}

//////

func (psm *sanStorageManager) reconcileSnapshotGroupVClone(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	groupName := *input.SnapshotGroupName
	log.WriteInfo("VClone| Starting reconciliation for group: %s", groupName)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, err
	}

	// 1. Pre-verification
	doCreate, doConvert, originalGroup, err := psm.preVerifySnapshotGroupVClone(groupName)
	if err != nil {
		log.WriteError("VClone| Pre-verification failed for group %s: %v", groupName, err)
		return nil, nil, err
	}

	// 2. Execution Phase
	if doCreate {
		log.WriteInfo("VClone| API: Calling CreateSnapshotGroupVClone (Promote) for group: %s", groupName)
		if _, err = provObj.CreateSnapshotGroupVClone(groupName); err != nil {
			return nil, nil, log.WriteAndReturnError("VClone| Create API failed: %v", err)
		}
	} else if doConvert {
		log.WriteInfo("VClone| API: Calling ConvertSnapshotGroupVClone (Finalize) for group: %s", groupName)
		if _, err = provObj.ConvertSnapshotGroupVClone(groupName); err != nil {
			return nil, nil, log.WriteAndReturnError("VClone| Convert API failed: %v", err)
		}
	}

	// 3. Wait and Gather
	log.WriteInfo("VClone| Entering Wait/Gather phase for %d members in group %s", len(originalGroup.Snapshots), groupName)
	latestGroup, vclones, err := psm.waitForVClonePromotion(groupName, originalGroup.Snapshots)
	if err != nil {
		log.WriteWarn("VClone| Wait/Gather completed with errors for group %s: %v", groupName, err)
		return originalGroup, nil, err
	}

	log.WriteInfo("VClone| Successfully reconciled group %s. Found %d verified vClones.", groupName, len(vclones))
	return latestGroup, vclones, nil
}

func (psm *sanStorageManager) waitForVClonePromotion(groupName string, originalSnaps []gwymodel.Snapshot) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Fetch dynamic retry configuration
	retryCfg := utils.GetRetryConfig() // used individually per api

	var mu sync.Mutex
	var verifiedFamily []gwymodel.SnapshotFamily

	// 2. Parallel polling for promotion status
	_ = utils.RunConcurrentOperations("WaitVClonePromotion", originalSnaps,
		func(snap gwymodel.Snapshot, idx int) error {
			log.WriteInfo("VCloneWait| Checking status for SVOL LDEV: %d", snap.SvolLdevID)

			var member *gwymodel.SnapshotFamily
			var err error
			found := false

			// Individual retry loop for each SVOL
			for attempt := 1; attempt <= retryCfg.MaxRetries; attempt++ {
				member, err = psm.getVcloneFromFamily(&snap.PvolLdevID, &snap.MuNumber, &snap.SvolLdevID, true)

				// Condition: Successfully retrieved and flag is true
				if err == nil && member != nil && member.IsVirtualCloneVolume {
					log.WriteInfo("VCloneWait| SVOL %d verified as Virtual Clone (Attempt %d/%d)",
						snap.SvolLdevID, attempt, retryCfg.MaxRetries)
					found = true
					break
				}

				log.WriteDebug("VCloneWait| SVOL %d still transitioning, retrying... (%d/%d)",
					snap.SvolLdevID, attempt, retryCfg.MaxRetries)

				if attempt < retryCfg.MaxRetries {
					time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
				}
			}

			if found && member != nil {
				mu.Lock()
				verifiedFamily = append(verifiedFamily, *member)
				mu.Unlock()
			} else {
				log.WriteWarn("VCloneWait| SVOL %d could not be verified in Family Table after promotion", snap.SvolLdevID)
			}

			return nil // Continue checking other members even if one times out
		})

	// 3. Fetch final state of the group
	provObj, _ := psm.getProvisionerManager()
	params := gwymodel.GetSnapshotGroupsParams{DetailInfoType: utils.Ptr("retention")}
	latestGroup, err := provObj.GetSnapshotGroup(groupName, params)
	if err != nil {
		log.WriteInfo("VCloneWait| Group %s no longer accessible (expected if converted to standalone Clones)", groupName)
		return nil, verifiedFamily, nil
	}

	return latestGroup, verifiedFamily, nil
}

func (psm *sanStorageManager) preVerifySnapshotGroupVClone(groupName string) (bool, bool, *gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	provObj, _ := psm.getProvisionerManager()
	params := gwymodel.GetSnapshotGroupsParams{DetailInfoType: utils.Ptr("retention")}

	group, err := provObj.GetSnapshotGroup(groupName, params)
	if err != nil {
		return false, false, nil, err
	}

	pairCount := 0
	psusCount := 0
	totalMembers := len(group.Snapshots)

	for _, snap := range group.Snapshots {
		status := strings.ToUpper(snap.Status)
		if status == "PAIR" {
			pairCount++
		} else if status == "PSUS" || status == "SSUS" {
			psusCount++
		}
	}

	// Block mixed states for consistency
	if pairCount > 0 && psusCount > 0 {
		return false, false, group, fmt.Errorf("vClone promotion aborted: group %s is in a mixed state (%d PAIR, %d PSUS)", groupName, pairCount, psusCount)
	}

	doCreate := false
	doConvert := false

	if pairCount == totalMembers && totalMembers > 0 {
		log.WriteInfo("VClonePreVerify| All members are PAIR. Using operationType: 'create'")
		doCreate = true
		doConvert = false // Based on your doc, 'create' is sufficient for PAIR
	} else if psusCount == totalMembers && totalMembers > 0 {
		log.WriteInfo("VClonePreVerify| All members are PSUS. Using operationType: 'convert'")
		doCreate = false
		doConvert = true
	} else {
		// If they are already vClones, both will be false and the wait loop will verify them
		log.WriteInfo("VClonePreVerify| Members are neither PAIR nor PSUS. Checking if already promoted...")
	}

	return doCreate, doConvert, group, nil
}

func (psm *sanStorageManager) reconcileSnapshotGroupVRestore(input reconmodel.SnapshotGroupReconcilerInput) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	groupName := *input.SnapshotGroupName
	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, err
	}

	// 1. Pre-verification: Check if any members actually need restoring
	_, originalGroup, err := psm.preVerifySnapshotGroupVRestore(groupName)
	if err != nil {
		return nil, nil, err
	}

	// 2. Execution Phase: Call Pair-Back/Revert API
	log.WriteInfo("VRestore| Calling RestoreSnapshotGroupFromVClone for group: %s", groupName)
	_, err = provObj.RestoreSnapshotGroupFromVClone(groupName)
	if err != nil {
		return originalGroup, nil, log.WriteAndReturnError("VRestore| API failed for group %s: %v", groupName, err)
	}

	// 3. Wait and Gather Phase: Monitor transition back to PSUS
	latestGroup, vclones, err := psm.waitForVRestoreCompletion(groupName, originalGroup.Snapshots)
	if err != nil {
		log.WriteError("VRestore| Wait phase encountered errors for group %s: %v", groupName, err)
		return latestGroup, vclones, err
	}

	log.WriteInfo("VRestore| Group %s successfully restored.", groupName)
	return latestGroup, vclones, nil
}

func (psm *sanStorageManager) preVerifySnapshotGroupVRestore(groupName string) (bool, *gwymodel.SnapshotGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provObj, _ := psm.getProvisionerManager()

	// 1. Initial Group Retrieval
	params := gwymodel.GetSnapshotGroupsParams{DetailInfoType: utils.Ptr("retention")}
	group, err := provObj.GetSnapshotGroup(groupName, params)
	if err != nil {
		return false, nil, err
	}

	// 2. Synchronous retention check (cheap/local data check)
	for _, snap := range group.Snapshots {
		if err := checkTIARetention(&snap); err != nil {
			return false, group, fmt.Errorf("vRestore group check failed: %w", err)
		}
	}

	var mu sync.Mutex
	var firstValidationError error

	// 3. Concurrent Metadata Retrieval & Validation
	// This handles the expensive getVcloneFromFamily API calls in parallel
	_ = utils.RunConcurrentOperations("PreVerifyVRestoreMembers", group.Snapshots,
		func(snap gwymodel.Snapshot, idx int) error {
			// Early exit if another goroutine already found an error
			mu.Lock()
			if firstValidationError != nil {
				mu.Unlock()
				return nil
			}
			mu.Unlock()

			member, err := psm.getVcloneFromFamily(&snap.PvolLdevID, &snap.MuNumber, &snap.SvolLdevID, true)

			mu.Lock()
			defer mu.Unlock()

			// Check if error occurred or member missing
			if err != nil || member == nil {
				firstValidationError = fmt.Errorf("vRestore failed: could not retrieve family metadata for member SVOL %d", snap.SvolLdevID)
				return firstValidationError
			}

			// Apply the strict vClone validation logic
			if !member.IsVirtualCloneVolume || !member.IsVirtualCloneParentVolume {
				firstValidationError = fmt.Errorf(
					"vRestore failed: snapshot (P-VOL: %d, MU: %d, S-VOL: %d) in group '%s' is not in a valid 'paired back' vClone state. "+
						"Both isVirtualCloneVolume and isVirtualCloneParentVolume must be true.",
					snap.PvolLdevID, snap.MuNumber, snap.SvolLdevID, groupName,
				)
				return firstValidationError
			}

			log.WriteDebug("VRestorePreVerify| Member %d passed strict validation.", snap.SvolLdevID)
			return nil
		})

	// If any concurrent check set an error, return it now
	if firstValidationError != nil {
		return false, group, firstValidationError
	}

	log.WriteInfo("VRestorePreVerify| Group %s is valid for Virtual Restore.", groupName)
	return true, group, nil
}

func (psm *sanStorageManager) waitForVRestoreCompletion(groupName string, originalSnaps []gwymodel.Snapshot) (*gwymodel.SnapshotGroup, []gwymodel.SnapshotFamily, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// 1. Fetch retry configuration from your environment-aware utility
	retryCfg := utils.GetRetryConfig() // used individually per api

	var mu sync.Mutex
	var remainingVclones []gwymodel.SnapshotFamily

	// 2. Parallel polling using your concurrency utility
	_ = utils.RunConcurrentOperations("WaitVRestoreCompletion", originalSnaps,
		func(snap gwymodel.Snapshot, idx int) error {
			log.WriteInfo("VRestoreWait| Monitoring SVOL %d for reversion...", snap.SvolLdevID)

			success := false
			// Use the dynamic MaxRetries from config
			for attempt := 1; attempt <= retryCfg.MaxRetries; attempt++ {
				member, err := psm.getVcloneFromFamily(&snap.PvolLdevID, &snap.MuNumber, &snap.SvolLdevID, true)

				// Success check: If error occurs (volume gone), member is nil, or flag is false
				if err != nil || member == nil || !member.IsVirtualCloneVolume {
					log.WriteInfo("VRestoreWait| SVOL %d: Successfully reverted (Attempt %d/%d)",
						snap.SvolLdevID, attempt, retryCfg.MaxRetries)
					success = true
					break
				}

				log.WriteDebug("VRestoreWait| SVOL %d: Still in vClone state (Attempt %d/%d)",
					snap.SvolLdevID, attempt, retryCfg.MaxRetries)

				// Use the dynamic Delay from config
				if attempt < retryCfg.MaxRetries {
					time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
				}
			}

			if !success {
				log.WriteWarn("VRestoreWait| SVOL %d: Failed to revert within %d attempts",
					snap.SvolLdevID, retryCfg.MaxRetries)

				// Final check to capture the state of the failed SVOL
				if m, _ := psm.getVcloneFromFamily(&snap.PvolLdevID, &snap.MuNumber, &snap.SvolLdevID, true); m != nil {
					mu.Lock()
					remainingVclones = append(remainingVclones, *m)
					mu.Unlock()
				}
			}
			return nil
		})

	// 3. Refresh final group status
	provObj, _ := psm.getProvisionerManager()
	latestGroup, err := provObj.GetSnapshotGroup(groupName, gwymodel.GetSnapshotGroupsParams{DetailInfoType: utils.Ptr("retention")})

	return latestGroup, remainingVclones, err
}
