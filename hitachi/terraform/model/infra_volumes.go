package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraVolumeInfo struct {
	model.VolumeInfo
}

type MtInfraVolumeInfo struct {
	model.MTVolumeInfo
}


type InfraVolumes struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []InfraVolumeInfo `json:"data"`
}


type MTInfraVolumes struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []MtInfraVolumeInfo `json:"data"`
}

type InfraVolumeTypes struct {
	Name                         string `json:"name,omitempty"`
	PoolID                       *int    `json:"poolId,omitempty"`
	ParityGroupId                string `json:"parityGroupId,omitempty"`
	Capacity                     string `json:"capacity,omitempty"`
	ResourceGroupId              *int    `json:"resourceGroupId,omitempty"`
	LunId                        *int    `json:"lunId,omitempty"`
	System                       string `json:"ucpSystem,omitempty"`
	DeduplicationCompressionMode string `json:"deduplicationCompressionMode,omitempty"`
}