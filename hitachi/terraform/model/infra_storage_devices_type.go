package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraStorageDeviceInfo struct {
	model.StorageDeviceInfo
}

type InfraStorageDevices struct {
	Path    string                   `json:"path"`
	Message string                   `json:"message"`
	Data    []InfraStorageDeviceInfo `json:"data"`
}

type CreateInfraStorageDeviceParam struct {
	model.CreateStorageDeviceParam
}
