package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetAllPartners gets partners information
func (psm *infraGwManager) GetAllUsers() (*model.Partners, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var partners model.Partners

	apiSuf := "/rbac/users?onlyUcpUsers=true"

	err := httpmethod.GetCall(psm.setting, apiSuf, &partners)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &partners, nil
}

func (psm *infraGwManager) GetUsersRoles(userId string) (*model.RoleWithDetails, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var useRoles model.RoleWithDetails
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = fmt.Sprintf("/rbac/users/%s/roles", userId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &useRoles)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &useRoles, nil
}
