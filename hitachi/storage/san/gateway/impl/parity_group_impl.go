package sanstorage

import (
	"fmt"
	"net/url"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetParityGroups
func (psm *sanStorageManager) GetParityGroups(detailInfoType string, driveTypeName string, clprId *int) (*[]sanmodel.ParityGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var parityGroups sanmodel.ParityGroups
	apiSuf := "objects/parity-groups"

	params := url.Values{}
	if detailInfoType != "" {
		params.Set("detailInfoType", detailInfoType)
	}
	if driveTypeName != "" {
		params.Set("driveTypeName", driveTypeName)
	}
	if clprId != nil {
		params.Set("clprId", strconv.Itoa(*clprId))
	}
	if len(params) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, params.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &parityGroups)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &parityGroups.Data, nil
}

// GetParityGroup
func (psm *sanStorageManager) GetParityGroup(pgid string) (*sanmodel.ParityGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var parityGroup sanmodel.ParityGroup
	apiSuf := fmt.Sprintf("objects/parity-groups/%s", pgid)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &parityGroup)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &parityGroup, nil
}
