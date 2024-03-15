package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraHostGroupInfo struct {
	model.HostGroupInfo
}

type InfraMTHostGroupInfo struct {
	model.MTHostGroupInfo
}

type InfraHostGroups struct {
	Path    string               `json:"path"`
	Message string               `json:"message"`
	Data    []InfraHostGroupInfo `json:"data"`
}

type InfraMTHostGroups struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraMTHostGroupInfo `json:"data"`
}

type CreateInfraHostGroupParam struct {
	model.CreateHostGroupParam
}
