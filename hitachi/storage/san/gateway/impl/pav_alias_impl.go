package sanstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetPavAliases returns PAV alias entries. When cuNumber is nil, returns all entries.
func (psm *sanStorageManager) GetPavAliases(cuNumber *int) (*[]sanmodel.PavAlias, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var resp sanmodel.PavAliasesResponse
	apiSuf := "objects/pav-ldevs"
	if cuNumber != nil {
		apiSuf = fmt.Sprintf("objects/pav-ldevs/%d", *cuNumber)
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &resp)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &resp.Data, nil
}

// AssignPavAlias assigns one or more alias LDEVs to a base LDEV.
func (psm *sanStorageManager) AssignPavAlias(baseLdevID int, aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/pav-ldevs/instance/actions/assign-alias/invoke"
	req := sanmodel.AssignPavAliasRequest{
		Parameters: sanmodel.AssignPavAliasParams{
			BaseLdevID:   baseLdevID,
			AliasLdevIDs: aliasLdevIDs,
		},
	}

	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, req)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}

// UnassignPavAlias unassigns one or more alias LDEVs.
func (psm *sanStorageManager) UnassignPavAlias(aliasLdevIDs []int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/pav-ldevs/instance/actions/unassign-alias/invoke"
	req := sanmodel.UnassignPavAliasRequest{
		Parameters: sanmodel.UnassignPavAliasParams{
			AliasLdevIDs: aliasLdevIDs,
		},
	}

	_, err := httpmethod.PostCall(psm.storageSetting, apiSuf, req)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return err
	}
	return nil
}
