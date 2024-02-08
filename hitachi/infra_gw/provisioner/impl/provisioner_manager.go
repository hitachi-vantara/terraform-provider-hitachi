package infra_gw

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	spmanager "terraform-provider-hitachi/hitachi/infra_gw/provisioner"
)

// infraGwManager contain information for setting
type infraGwManager struct {
	setting model.InfraGwSettings
}

// A private function to construct an newInfraGWManagerEx
func newInfraGwManagerEx(setting model.InfraGwSettings) (*infraGwManager, error) {
	psm := &infraGwManager{
		setting: model.InfraGwSettings{
			Username: setting.Username,
			Password: setting.Password,
			Address:  setting.Address,
		},
	}
	return psm, nil
}

// NewEx returns a new Infra Gateway Provisioner
func NewEx(setting model.InfraGwSettings) (spmanager.InfraGwManager, error) {
	return newInfraGwManagerEx(setting)
}

func New(userName, password, address string) (spmanager.InfraGwManager, error) {

	psm := &infraGwManager{
		setting: model.InfraGwSettings{
			Username: userName,
			Password: password,
			Address:  address,
		},
	}

	
	status,partnerId, _ := psm.GetPartnerIdWithStatus(userName)

	if status {
		psm.setting.PartnerId = *partnerId
	}

	return psm, nil

}
