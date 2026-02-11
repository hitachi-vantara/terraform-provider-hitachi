package admin

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func (psm *adminStorageManager) ReconcileCreateAdminPool(params gwymodel.CreateAdminPoolParams) (int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool creation reconciliation for pool %s", params.Name)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Create Pool Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return 0, err
	}

	err = provObj.CreateAdminPool(params)
	if err != nil {
		log.WriteError("Failed to create admin pool %s: %v", params.Name, err)
		return 0, fmt.Errorf("admin pool creation failed: %w", err)
	}

	log.WriteInfo("Admin pool %s created successfully", params.Name)

	// For pool creation, we need to find the created pool since the API doesn't return the pool ID
	// We'll search by name since pool names should be unique
	listParams := gwymodel.AdminPoolListParams{
		Name: &params.Name,
	}

	poolList, err := provObj.GetAdminPoolList(listParams)
	if err != nil {
		log.WriteError("Failed to get created pool info: %v", err)
		return 0, fmt.Errorf("failed to get created pool info: %w", err)
	}

	if poolList == nil || len(poolList.Data) == 0 {
		return 0, fmt.Errorf("created pool %s not found", params.Name)
	}

	// Return the pool ID of the first matching pool (should be the one we just created)
	createdPoolID := poolList.Data[0].ID
	log.WriteInfo("Successfully retrieved created pool ID: %d", createdPoolID)

	return createdPoolID, nil
}

func (psm *adminStorageManager) ReconcileReadAdminPool(poolID int) (*gwymodel.AdminPool, bool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool read reconciliation for ID %d", poolID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, false, err
	}

	result, err := provObj.GetAdminPoolInfo(poolID)
	if err != nil {
		if IsNotFoundError(err) {
			log.WriteWarn("Admin pool %d not found or already deleted", poolID)
			return nil, false, err
		}
		log.WriteError("Failed to get admin pool info for ID %d: %v", poolID, err)
		return nil, false, fmt.Errorf("failed to get admin pool info: %w", err)
	}

	if result == nil {
		log.WriteWarn("Admin pool %d returned nil", poolID)
		return nil, false, fmt.Errorf("admin pool %d not found", poolID)
	}

	// Poll for capacity consistency with configurable timeout and interval
	pollTimeout := getPollTimeout()
	pollInterval := getPollInterval()
	startTime := time.Now()
	timedOut := false

	var previousResult *gwymodel.AdminPool

	for time.Since(startTime) < pollTimeout {
		// If we have a previous result, check if capacity values are consistent between old and new
		if previousResult != nil && isCapacityConsistent(previousResult, result) {
			log.WriteInfo("Capacity values are consistent between old and new results for pool %d", poolID)
			break
		}

		if previousResult == nil {
			log.WriteInfo("First capacity reading for pool %d, polling again in %v to compare", poolID, pollInterval)
		} else {
			log.WriteInfo("Capacity values not yet consistent between readings for pool %d, polling again in %v", poolID, pollInterval)
		}

		// Store current result as previous for next iteration
		previousResult = result

		time.Sleep(pollInterval)

		// Re-fetch pool info
		result, err = provObj.GetAdminPoolInfo(poolID)
		if err != nil {
			log.WriteError("Failed to re-fetch admin pool info for ID %d during polling: %v", poolID, err)
			break // Don't fail, just stop polling
		}

		if result == nil {
			log.WriteWarn("Admin pool %d returned nil during polling", poolID)
			break // Don't fail, just stop polling
		}
	}

	if time.Since(startTime) >= pollTimeout {
		timedOut = true
		log.WriteWarn("Capacity consistency polling timed out after %v for pool %d, proceeding anyway", pollTimeout, poolID)
	}

	log.WriteInfo("Admin pool read reconciliation completed successfully for ID %d", poolID)
	return result, timedOut, nil
}

// isCapacityConsistent checks if capacity values are consistent between old and new pool results
func isCapacityConsistent(oldPool, newPool *gwymodel.AdminPool) bool {

	totalCapacityMatches := oldPool.TotalCapacity == newPool.TotalCapacity
	effectiveCapacityMatches := oldPool.EffectiveCapacity == newPool.EffectiveCapacity
	freeCapacityMatches := oldPool.FreeCapacity == newPool.FreeCapacity

	// All capacity values must match exactly between old and new results
	return totalCapacityMatches && effectiveCapacityMatches && freeCapacityMatches
}

// getPollTimeout returns the polling timeout from environment variable or default
func getPollTimeout() time.Duration {
	if timeoutStr := os.Getenv("POOL_POLL_TIMEOUT_SEC"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			return time.Duration(timeout) * time.Second
		}
	}
	return 60 * time.Minute // default 60 minutes
}

// getPollInterval returns the polling interval from environment variable or default
func getPollInterval() time.Duration {
	if intervalStr := os.Getenv("POOL_POLL_INTERVAL_SEC"); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			return time.Duration(interval) * time.Second
		}
	}
	return 30 * time.Second // default 30 seconds
}

func (psm *adminStorageManager) ReconcileUpdateAdminPool(poolID int, params gwymodel.UpdateAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool update reconciliation for ID %d", poolID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Update Pool Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.UpdateAdminPool(poolID, params)
	if err != nil {
		log.WriteError("Failed to update admin pool ID %d: %v", poolID, err)
		return fmt.Errorf("admin pool update failed: %w", err)
	}

	log.WriteInfo("Admin pool %d updated successfully", poolID)
	return nil
}

func (psm *adminStorageManager) ReconcileExpandAdminPool(poolID int, params gwymodel.ExpandAdminPoolParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool expansion reconciliation for ID %d", poolID)

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("Expand Pool Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.ExpandAdminPool(poolID, params)
	if err != nil {
		log.WriteError("Failed to expand admin pool ID %d: %v", poolID, err)
		return fmt.Errorf("admin pool expansion failed: %w", err)
	}

	log.WriteInfo("Admin pool %d expanded successfully", poolID)
	return nil
}

func (psm *adminStorageManager) ReconcileDeleteAdminPool(poolID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool deletion reconciliation for ID %d", poolID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return err
	}

	err = provObj.DeleteAdminPool(poolID)
	if err != nil {
		log.WriteError("Failed to delete admin pool ID %d: %v", poolID, err)
		return fmt.Errorf("admin pool deletion failed: %w", err)
	}

	log.WriteInfo("Admin pool deleted successfully for ID %d", poolID)
	return nil
}

func (psm *adminStorageManager) GetAdminPoolList(params gwymodel.AdminPoolListParams) (*gwymodel.AdminPoolListResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool list retrieval")

	// Debug log the parameters
	b, _ := json.MarshalIndent(params, "", "  ")
	log.WriteDebug("List Pool Params: %v", string(b))

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminPoolList(params)
	if err != nil {
		log.WriteError("Failed to get admin pool list: %v", err)
		return nil, fmt.Errorf("admin pool list retrieval failed: %w", err)
	}

	if result == nil {
		log.WriteWarn("No admin pool list returned")
		return nil, fmt.Errorf("no admin pool list returned")
	}

	log.WriteInfo("Admin pool list retrieved successfully with %d pools", len(result.Data))
	return result, nil
}

func (psm *adminStorageManager) GetAdminPoolInfo(poolID int) (*gwymodel.AdminPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("Starting admin pool info retrieval for ID %d", poolID)

	provObj, err := psm.getProvisionerManager()
	if err != nil {
		log.WriteError("Failed to get provisioner manager: %v", err)
		return nil, err
	}

	result, err := provObj.GetAdminPoolInfo(poolID)
	if err != nil {
		log.WriteError("Failed to get admin pool info for ID %d: %v", poolID, err)
		return nil, fmt.Errorf("admin pool info retrieval failed: %w", err)
	}

	if result == nil {
		log.WriteWarn("No admin pool info returned for ID %d", poolID)
		return nil, fmt.Errorf("no admin pool info returned for ID %d", poolID)
	}

	log.WriteInfo("Admin pool info retrieved successfully for ID %d", poolID)
	return result, nil
}
