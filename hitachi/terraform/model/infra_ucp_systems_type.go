package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraUcpSystemInfo struct {
	model.UcpSystemInfo
}

type InfraUcpSystems struct {
	Path    string               `json:"path"`
	Message string               `json:"message"`
	Data    []InfraUcpSystemInfo `json:"data"`
}
