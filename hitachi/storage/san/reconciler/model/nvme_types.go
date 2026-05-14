package sanstorage

// NvmSubsystemReconcilerInput defines the schema for reconciling an NVM Subsystem
type NvmSubsystemReconcilerInput struct {
	NvmSubsystemId          *int    `json:"nvmSubsystemId"`
	VirtualNvmSubsystemId   *int    `json:"virtualNvmSubsystemId,omitempty"`
	NvmSubsystemName        *string `json:"nvmSubsystemName,omitempty"`
	HostMode                *string `json:"hostMode,omitempty"`
	HostModeOptions         *[]int  `json:"hostModeOptions,omitempty"`
	EnableNamespaceSecurity *bool   `json:"enableNamespaceSecurity,omitempty"`

	Ports      *[]string         `json:"ports,omitempty"`
	HostNqns   *[]HostNqnInput   `json:"hostNqns,omitempty"`
	Namespaces *[]NamespaceInput `json:"namespaces,omitempty"`
}

type HostNqnInput struct {
	Nqn      string `json:"nqn"`
	Nickname string `json:"nickname,omitempty"`
}

type NamespaceInput struct {
	NamespaceId int      `json:"namespaceId"` // auto set
	LdevId      int      `json:"ldevId"`
	Nickname    string   `json:"nickname,omitempty"`
	Path        []string `json:"path,omitempty"`
}

// NvmSubsystemGetMultipleInput defines the filtering and inclusion
// parameters for the NVM Subsystem Data Source.
type NvmSubsystemGetMultipleInput struct {
	// Filters
	All              *bool   `json:"all,omitempty"`
	NvmSubsystemId   *int    `json:"nvmSubsystemId,omitempty"`
	NvmSubsystemName *string `json:"nvmSubsystemName,omitempty"`

	// Inclusion Flags (maps to query parameters)
	IncludePorts      bool `json:"includePorts"`
	IncludeHostNqns   bool `json:"includeHostNqns"`
	IncludeNamespaces bool `json:"includeNamespaces"`
	IncludeNvmNqn     bool `json:"includeNvmNqn"` // internal
}
