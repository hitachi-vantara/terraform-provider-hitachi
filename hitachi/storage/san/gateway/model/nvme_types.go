package sanstorage

// NvmSubsystem represents the detailed structure returned by the Hitachi API
type NvmSubsystem struct {
	NvmSubsystemId           int           `json:"nvmSubsystemId"`
	VirtualNvmSubsystemId    int           `json:"virtualNvmSubsystemId,omitempty"`
	NvmSubsystemName         string        `json:"nvmSubsystemName,omitempty"`
	ResourceGroupId          int           `json:"resourceGroupId,omitempty"`
	NamespaceSecuritySetting string        `json:"namespaceSecuritySetting,omitempty"`
	T10piMode                string        `json:"t10piMode,omitempty"`
	HostMode                 string        `json:"hostMode,omitempty"`
	HostModeOptions          []int         `json:"hostModeOptions,omitempty"`
	NvmSubsystemNqn          string        `json:"nvmSubsystemNqn,omitempty"`
	HostNqns                 []HostNqnInfo `json:"hostNqns,omitempty"`
	Namespaces               []Namespace   `json:"namespaces,omitempty"`
	PortIds                  []string      `json:"portIds,omitempty"`
}

// NvmSubsystems represents the list response wrapper
type NvmSubsystems struct {
	Data []NvmSubsystem `json:"data"`
}

type HostNqnInfo struct {
	HostNqnId       string `json:"hostNqnId"`
	HostNqn         string `json:"hostNqn"`
	NvmSubsystemId  int    `json:"nvmSubsystemId"`
	HostNqnNickname string `json:"hostNqnNickname"`
}

type HostNqns struct {
	Data []HostNqnInfo `json:"data"`
}

// GetNvmSubsystemsParams for query parameters
type GetNvmSubsystemsParams struct {
	NvmSubsystemInfo   *string `json:"nvmSubsystemInfo,omitempty"`
	NvmSubsystemOption *string `json:"nvmSubsystemOption,omitempty"`
}

// CreateNvmSubsystemRequest used for the POST body
type CreateNvmSubsystemRequest struct {
	NvmSubsystemId           int     `json:"nvmSubsystemId"`
	VirtualNvmSubsystemId    *int    `json:"virtualNvmSubsystemId,omitempty"`
	NvmSubsystemName         *string `json:"nvmSubsystemName,omitempty"`
	HostMode                 *string `json:"hostMode,omitempty"`
	HostModeOptions          *[]int  `json:"hostModeOptions,omitempty"`
	NamespaceSecuritySetting *string `json:"namespaceSecuritySetting,omitempty"`
}

// AddNvmSubsystemPortRequest for POST /objects/nvm-subsystem-ports
type AddNvmSubsystemPortRequest struct {
	NvmSubsystemId int    `json:"nvmSubsystemId"`
	PortId         string `json:"portId"`
}

// NvmSubsystemPort for GET /objects/nvm-subsystem-ports/{object-ID}
type NvmSubsystemPort struct {
	NvmSubsystemId     int    `json:"nvmSubsystemId"`
	PortId             string `json:"portId"`
	NvmSubsystemPortId string `json:"nvmSubsystemPortId,omitempty"` // e.g., "Fibre Channel" or "iSCSI"
}

type NvmSubsystemPorts struct {
	Data []NvmSubsystemPort `json:"data"`
}

// UpdateNvmSubsystemRequest used for the PATCH body
type UpdateNvmSubsystemRequest struct {
	NvmSubsystemName         *string `json:"nvmSubsystemName,omitempty"`
	HostMode                 *string `json:"hostMode,omitempty"`
	HostModeOptions          *[]int  `json:"hostModeOptions,omitempty"`
	NamespaceSecuritySetting *string `json:"namespaceSecuritySetting,omitempty"`
}

type RegisterHostNqnRequest struct {
	NvmSubsystemId int    `json:"nvmSubsystemId"`
	HostNqn        string `json:"hostNqn"`
}

// NamespaceNicknameRequest for PATCH /objects/namespaces/object-ID
type NamespaceNicknameRequest struct {
	NamespaceNickname string `json:"namespaceNickname"`
}

// NamespacePath represents the host-namespace path information
type NamespacePath struct {
	NamespacePathId string `json:"namespacePathId"`
	NvmSubsystemId  int    `json:"nvmSubsystemId"`
	HostNqn         string `json:"hostNqn"`
	NamespaceId     int    `json:"namespaceId"`
	LdevId          int    `json:"ldevId"`
}

// NamespacePathsResponse for GET /objects/namespace-paths
type NamespacePathsResponse struct {
	Data    []NamespacePath `json:"data"`
	HasNext bool            `json:"hasNext"`
	NextId  string          `json:"nextId,omitempty"`
}

// GetNamespacePathsParams for GET query parameters
type GetNamespacePathsParams struct {
	NvmSubsystemId int     `json:"nvmSubsystemId"`
	NamespaceId    *int    `json:"namespaceId,omitempty"`
	HeadId         *string `json:"headId,omitempty"`
}

// RegisterNamespacePathRequest for POST /objects/namespace-paths
type RegisterNamespacePathRequest struct {
	NvmSubsystemId int    `json:"nvmSubsystemId"`
	HostNqn        string `json:"hostNqn"`
	NamespaceId    int    `json:"namespaceId"`
}

// SetHostNqnNicknameRequest for PATCH /objects/host-nqns/object-ID
type SetHostNqnNicknameRequest struct {
	HostNqnNickname string `json:"hostNqnNickname"`
}

type Namespace struct {
	NamespaceObjectId  string   `json:"namespaceObjectId"`
	NamespaceId        int      `json:"namespaceId"`
	NamespaceNickname  string   `json:"namespaceNickname"`
	NvmSubsystemId     int      `json:"nvmSubsystemId"`
	NvmSubsystemName   string   `json:"nvmSubsystemName"`
	LdevId             int      `json:"ldevId"`
	ByteFormatCapacity string   `json:"byteFormatCapacity"`
	BlockCapacity      int64    `json:"blockCapacity"`
	Paths              []string `json:"paths,omitempty"` // added to capture the host-namespace paths, not part of the original API response
}

type Namespaces struct {
	Data []Namespace `json:"data"`
}

type CreateNamespaceRequest struct {
	NvmSubsystemId int  `json:"nvmSubsystemId"`
	NamespaceId    *int `json:"namespaceId,omitempty"`
	LdevId         int  `json:"ldevId"`
}

type SetNamespaceNicknameRequest struct {
	NamespaceNickname string `json:"namespaceNickname"`
}
