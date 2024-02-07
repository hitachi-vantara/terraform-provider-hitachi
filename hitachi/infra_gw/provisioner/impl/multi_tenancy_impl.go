package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetUcpSystems gets all UCP Systems information
func (psm *infraGwManager) GetUserDetails(username string) (*bool, *string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, nil, err
	}

	users, err := gatewayObj.GetAllUsers()
	fmt.Println(users)
	return nil, nil, nil

}
