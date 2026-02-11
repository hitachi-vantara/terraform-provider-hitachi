package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetStoragePorts used to get storage ports information with optional filters
func (psm *sanStorageManager) GetStoragePorts(detailInfoTypes []string, portType, portAttributes string) (*[]sanmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var storagePorts sanmodel.StoragePorts
	apiSuf := "objects/ports"

	// Build query parameters
	queryParams := []string{}

	// Add portType parameter
	if portType != "" {
		queryParams = append(queryParams, "portType="+portType)
	}

	// Add portAttributes parameter
	if portAttributes != "" {
		queryParams = append(queryParams, "portAttributes="+portAttributes)
	}

	// Add detailInfoType parameters
	if len(detailInfoTypes) > 0 {
		var validTypes []string
		for _, dit := range detailInfoTypes {
			if dit != "" {
				validTypes = append(validTypes, dit)
			}
		}

		if len(validTypes) > 0 {
			// Prioritize portMode first since API only processes the first parameter
			var prioritizedTypes []string
			var otherTypes []string

			for _, dit := range validTypes {
				if dit == "portMode" {
					prioritizedTypes = append(prioritizedTypes, dit)
				} else {
					otherTypes = append(otherTypes, dit)
				}
			}

			// Combine prioritized types first, then others
			allTypes := append(prioritizedTypes, otherTypes...)

			queryParams = append(queryParams, "detailInfoType="+allTypes[0])
			// Add additional detailInfoType parameters
			for _, dit := range allTypes[1:] {
				queryParams = append(queryParams, "detailInfoType="+dit)
			}
		}
	}

	// Build final API suffix with query parameters
	if len(queryParams) > 0 {
		apiSuf += "?"
		for i, param := range queryParams {
			if i > 0 {
				apiSuf += "&"
			}
			apiSuf += param
		}
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &storagePorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in objects/ports API call, err: %v", err)
		return nil, err
	}

	return &storagePorts.Data, nil
}
