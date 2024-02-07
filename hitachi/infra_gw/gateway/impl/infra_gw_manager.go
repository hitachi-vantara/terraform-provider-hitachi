package infra_gw

import (
	spmanager "terraform-provider-hitachi/hitachi/infra_gw/gateway"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
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
			V3API: false,
		},
	}
	return psm, nil
}

// NewEx returns a new Infra Gateway Provisioner
func NewEx(setting model.InfraGwSettings) (spmanager.InfraGwManager, error) {
	return newInfraGwManagerEx(setting)
}
