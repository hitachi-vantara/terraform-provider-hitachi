package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraStoragePortInfo struct {
	model.StoragePortInfo
}

type InfraStoragePorts struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraStoragePortInfo `json:"data"`
}

type InfraGwStoragePort struct {
	Path    string               `json:"path"`
	Message string               `json:"message"`
	Data    InfraStoragePortInfo `json:"data"`
}
