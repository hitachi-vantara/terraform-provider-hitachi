package terraform

import (
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

type InfraIscsiTargetInfo struct {
	model.IscsiTargetInfo
}

type InfraIscsiTargets struct {
	Path    string                 `json:"path"`
	Message string                 `json:"message"`
	Data    []InfraIscsiTargetInfo `json:"data"`
}
