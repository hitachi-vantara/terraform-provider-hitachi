package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetParityGroups gets parity groups information
func (psm *infraGwManager) GetParityGroups(id string) (*model.ParityGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var parityGroups model.ParityGroups

	apiSuf := fmt.Sprintf("/storage/devices/%s/parityGroups", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &parityGroups)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &parityGroups, nil
}
