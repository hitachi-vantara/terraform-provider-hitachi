package sanstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetDynamicPools
func (psm *sanStorageManager) GetDynamicPools(isMainframe *bool, poolType string, detailInfoType ...string) (*[]sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var dynamicPools sanmodel.DynamicPools
	apiSuf := "objects/pools"

	// Build query parameters
	queryParams := []string{}

	// Add isMainframe filter if provided
	if isMainframe != nil {
		queryParams = append(queryParams, fmt.Sprintf("isMainframe=%t", *isMainframe))
	}

	// Add poolType filter if provided
	if poolType != "" {
		queryParams = append(queryParams, fmt.Sprintf("poolType=%s", poolType))
	}

	// Add detailInfoType as query parameter if provided
	if len(detailInfoType) > 0 && detailInfoType[0] != "" {
		queryParams = append(queryParams, fmt.Sprintf("detailInfoType=%s", detailInfoType[0]))
	}

	// Construct final API URL with query parameters
	if len(queryParams) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, strings.Join(queryParams, "&"))
	}

	log.WriteDebug("TFDebug| GetDynamicPools: API URL constructed: %s", apiSuf)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &dynamicPools)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &dynamicPools.Data, nil
}

// GetDynamicPoolById
func (psm *sanStorageManager) GetDynamicPoolById(poolId int) (*sanmodel.DynamicPool, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var dynamicPool sanmodel.DynamicPool
	apiSuf := fmt.Sprintf("objects/pools/%d", poolId)

	log.WriteDebug("TFDebug| GetDynamicPoolById: API URL: %s", apiSuf)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &dynamicPool)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &dynamicPool, nil
}
