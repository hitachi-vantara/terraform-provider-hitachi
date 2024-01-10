package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraParityGroupInfo struct {
	model.ParityGroupInfo
}

type InfraParityGroups struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraParityGroupInfo `json:"data"`
}
