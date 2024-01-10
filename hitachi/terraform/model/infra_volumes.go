package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraVolumeInfo struct {
	model.VolumeInfo
}

type InfraVolumes struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []InfraVolumeInfo `json:"data"`
}
