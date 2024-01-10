package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetVolumes gets volumes information
func (psm *infraGwManager) GetVolumes(id string) (*model.Volumes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response model.Volumes

	apiSuf := fmt.Sprintf("/storage/devices/%s/volumes", id)
	err := httpmethod.GetCall(psm.setting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &response, nil
}
