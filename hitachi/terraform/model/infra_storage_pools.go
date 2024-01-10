package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraStoragePoolInfo struct {
	model.StoragePoolInfo
}

type InfraStoragePools struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraStoragePoolInfo `json:"data"`
}

type InfraStoragePool struct {
	Path    string               `json:"path"`
	Message string               `json:"message"`
	Data    InfraStoragePoolInfo `json:"data"`
}
