package admin

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	recmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"
)

// ReconcileReadVolumeServerConnections fetches multiple volume-server connections in parallel.
func (psm *adminStorageManager) ReconcileReadVolumeServerConnections(pairs []recmodel.VolumeServerPair) ([]gwymodel.VolumeServerConnectionDetail, []recmodel.VolumeServerPair, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting reconciliation read for volume-server connections")
	log.WriteDebug("Pairs to reconcile: %+v", pairs)

	if len(pairs) == 0 {
		return nil, nil, fmt.Errorf("no volume-server connection pairs provided")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return nil, nil, err
	}

	var mu sync.Mutex
	var results []gwymodel.VolumeServerConnectionDetail
	var existing []recmodel.VolumeServerPair
	var firstErr error

	errs := utils.RunConcurrentOperations(
		"Read", pairs,
		func(pair recmodel.VolumeServerPair, idx int) error {

			conn, err := provObj.GetOneVolumeServerConnection(pair.VolumeID, pair.ServerID)
			if err != nil {
				if IsNotFoundError(err) || strings.Contains(err.Error(), "not found") {
					log.WriteWarn("[Read #%d] Connection not found: volume=%d server=%d",
						idx+1, pair.VolumeID, pair.ServerID)
					return nil
				}

				log.WriteError("[Read #%d] Error fetching volume=%d server=%d: %v",
					idx+1, pair.VolumeID, pair.ServerID, err)

				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("volume=%d server=%d: %w", pair.VolumeID, pair.ServerID, err)
				}
				mu.Unlock()

				return err
			}

			if conn != nil {
				mu.Lock()
				results = append(results, *conn)
				existing = append(existing, pair)
				mu.Unlock()
			}
			return nil
		},
	)

	for _, e := range errs {
		if e != nil && firstErr == nil {
			firstErr = e
		}
	}

	if len(existing) == 0 {
		log.WriteWarn("No existing connections found among provided pairs")
		if firstErr == nil {
			firstErr = fmt.Errorf("no existing volume-server connections found")
		}
		return nil, nil, firstErr
	}

	// Stable order
	sort.Slice(results, func(i, j int) bool {
		if results[i].ServerId == results[j].ServerId {
			return results[i].VolumeId < results[j].VolumeId
		}
		return results[i].ServerId < results[j].ServerId
	})
	sort.Slice(existing, func(i, j int) bool {
		if existing[i].ServerID == existing[j].ServerID {
			return existing[i].VolumeID < existing[j].VolumeID
		}
		return existing[i].ServerID < existing[j].ServerID
	})

	log.WriteInfo("Reconciliation read completed successfully: %d connections", len(existing))
	return results, existing, firstErr
}

// ReconcileDeleteVolumeServerConnections deletes multiple volume-server connections concurrently.
func (psm *adminStorageManager) ReconcileDeleteVolumeServerConnections(pairs []recmodel.VolumeServerPair) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting reconciliation delete for volume-server connections")

	if len(pairs) == 0 {
		return fmt.Errorf("no volume-server connection pairs provided")
	}

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	errs := utils.RunConcurrentOperations(
		"Delete",
		pairs,
		func(pair recmodel.VolumeServerPair, idx int) error {

			log.WriteDebug("[Delete #%d] Detaching volume=%d server=%d",
				idx+1, pair.VolumeID, pair.ServerID)

			err := provObj.DetachVolumeFromServer(pair.VolumeID, pair.ServerID)
			if err != nil {
				if IsNotFoundError(err) || strings.Contains(err.Error(), "not found") {
					log.WriteWarn("[Delete #%d] volume=%d server=%d not found or already detached",
						idx+1, pair.VolumeID, pair.ServerID)
					return nil
				}

				log.WriteError("[Delete #%d] Failed volume=%d server=%d: %v",
					idx+1, pair.VolumeID, pair.ServerID, err)

				return fmt.Errorf("volume=%d server=%d: %v",
					pair.VolumeID, pair.ServerID, err)
			}

			return nil
		},
	)

	// Collect errors
	var allErrs []string
	for _, e := range errs {
		if e != nil {
			allErrs = append(allErrs, e.Error())
		}
	}

	if len(allErrs) > 0 {
		errMsg := fmt.Sprintf("one or more detaches failed: %s", strings.Join(allErrs, "; "))
		log.WriteError(errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	log.WriteInfo("All volume-server connections deleted successfully")
	return nil
}

// ReconcileUpdateVolumeServerConnections performs reconciliation by calling
// a single Create (bulk attach) for all desired pairs, then detaching obsolete ones
// individually.
func (psm *adminStorageManager) ReconcileUpdateVolumeServerConnections(existingPairs, desiredPairs []recmodel.VolumeServerPair) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting reconciliation update for volume-server connections (Create-first)")

	if len(desiredPairs) == 0 {
		return fmt.Errorf("no desired pairs provided")
	}

	log.WriteDebug("Existing pairs: %+v", existingPairs)
	log.WriteDebug("Desired pairs: %+v", desiredPairs)

	// Compute pairs to detach (in existing but not in desired)
	toDetach := diffVolumeServerPairs(existingPairs, desiredPairs)
	log.WriteInfo("Pairs to detach after Create: %+v", toDetach)

	// Compute pairs to add (in desired but not in existing)
	toAdd := diffVolumeServerPairs(desiredPairs, existingPairs)
	log.WriteInfo("Pairs to add (new connections): %+v", toAdd)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		return err
	}

	// 🔹 STEP 1: Call Create only if there are new pairs to add
	if len(toAdd) > 0 {
		log.WriteInfo("New volume-server connections detected — calling Create")

		// Collect unique volumeIds and serverIds from toAdd only
		volSet := make(map[int]struct{})
		srvSet := make(map[int]struct{})
		for _, pair := range toAdd {
			volSet[pair.VolumeID] = struct{}{}
			srvSet[pair.ServerID] = struct{}{}
		}

		volumeIds := make([]int, 0, len(volSet))
		serverIds := make([]int, 0, len(srvSet))
		for v := range volSet {
			volumeIds = append(volumeIds, v)
		}
		for s := range srvSet {
			serverIds = append(serverIds, s)
		}

		params := &gwymodel.AttachVolumeServerConnectionParam{
			VolumeIds: volumeIds,
			ServerIds: serverIds,
		}

		if _, err := provObj.AttachVolumeToServers(*params); err != nil {
			log.WriteError("Bulk Create (AttachVolumeToServers) failed: %v", err)
			return fmt.Errorf("bulk create failed: %v", err)
		}

		log.WriteInfo("Bulk Create completed successfully")
	} else {
		log.WriteInfo("No new volume-server connections to create — skipping Create step")
	}

	// 🔹 STEP 2: Detach obsolete connections one at a time
	if len(toDetach) > 0 {
		log.WriteInfo("Detaching obsolete connections...")

		errs := utils.RunConcurrentOperations(
			"Update-Detach",
			toDetach,
			func(pair recmodel.VolumeServerPair, idx int) error {

				log.WriteDebug("[Detach #%d] Detaching volume=%d server=%d",
					idx+1, pair.VolumeID, pair.ServerID)

				err := provObj.DetachVolumeFromServer(pair.VolumeID, pair.ServerID)
				if err != nil {

					if IsNotFoundError(err) || strings.Contains(err.Error(), "not found") {
						log.WriteWarn("[Detach #%d] volume=%d server=%d not found or already detached",
							idx+1, pair.VolumeID, pair.ServerID)
						return nil
					}

					log.WriteError("[Detach #%d] Failed volume=%d server=%d: %v",
						idx+1, pair.VolumeID, pair.ServerID, err)

					// return formatted error but still allow other ops to continue
					return fmt.Errorf("volume=%d server=%d: %v",
						pair.VolumeID, pair.ServerID, err)
				}

				return nil
			},
		)

		// Aggregate errors into a single message
		var allErrs []string
		for _, e := range errs {
			if e != nil {
				allErrs = append(allErrs, e.Error())
			}
		}

		if len(allErrs) > 0 {
			errMsg := fmt.Sprintf("one or more detaches failed: %s", strings.Join(allErrs, "; "))
			log.WriteError(errMsg)
			return fmt.Errorf("%s", errMsg)
		}

		log.WriteInfo("All obsolete detaches completed successfully")
	} else {
		log.WriteInfo("No obsolete connections to detach")
	}

	log.WriteInfo("Reconciliation update completed successfully (Create-if-needed strategy)")
	return nil
}

// diffVolumeServerPairs returns items in 'a' that are not in 'b'.
func diffVolumeServerPairs(a, b []recmodel.VolumeServerPair) []recmodel.VolumeServerPair {
	m := make(map[string]bool, len(b))
	for _, pair := range b {
		key := fmt.Sprintf("%d,%d", pair.VolumeID, pair.ServerID)
		m[key] = true
	}

	diff := make([]recmodel.VolumeServerPair, 0)
	for _, pair := range a {
		key := fmt.Sprintf("%d,%d", pair.VolumeID, pair.ServerID)
		if !m[key] {
			diff = append(diff, pair)
		}
	}
	return diff
}
