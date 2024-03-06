package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraStoragePortInfo struct {
	model.StoragePortInfo
}

type InfraMTStoragePortInfo struct {
	model.MTPortInfo
}

type InfraStoragePorts struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraStoragePortInfo `json:"data"`
}

type InfraMTStoragePorts struct {
	Path    string                   `json:"path"`
	Message string                   `json:"message"`
	Error   TFError                  `json:"error"`
	Data    []InfraMTStoragePortInfo `json:"data"`
}

type TFError struct {
	Message string `json:"message"`
}

type InfraGwStoragePort struct {
	Path    string               `json:"path"`
	Message string               `json:"message"`
	Data    InfraStoragePortInfo `json:"data"`
}
