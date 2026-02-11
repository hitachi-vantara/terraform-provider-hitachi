package admin

// AdminServer represents a server in the admin system
type AdminServer struct {
	ID                     int    `json:"id"`
	Nickname               string `json:"nickname"`
	Protocol               string `json:"protocol"`
	OsType                 string `json:"osType"`
	TotalCapacity          int64  `json:"totalCapacity"`
	UsedCapacity           int64  `json:"usedCapacity"`
	NumberOfPaths          int    `json:"numberOfPaths"`
	IsInconsistent         bool   `json:"isInconsistent"`
	ModificationInProgress bool   `json:"modificationInProgress"`
	IsReserved             bool   `json:"isReserved"`
	HasUnalignedOsTypes    bool   `json:"hasUnalignedOsTypes"`
}

// AdminServerListResponse represents the response for listing servers
type AdminServerListResponse struct {
	Data  []AdminServer `json:"data"`
	Count int           `json:"count"`
}

// AdminServerPath represents server path information
type AdminServerPath struct {
	HbaWwn    string   `json:"hbaWwn"`
	IscsiName string   `json:"iscsiName"`
	PortIds   []string `json:"portIds"`
}

// AdminServerInfo represents detailed server information
type AdminServerInfo struct {
	ID                        int               `json:"id"`
	Nickname                  string            `json:"nickname"`
	Protocol                  string            `json:"protocol"`
	OsType                    string            `json:"osType"`
	OsTypeOptions             []int             `json:"osTypeOptions"`
	TotalCapacity             int64             `json:"totalCapacity"`
	UsedCapacity              int64             `json:"usedCapacity"`
	NumberOfVolumes           int               `json:"numberOfVolumes"`
	NumberOfPaths             int               `json:"numberOfPaths"`
	Paths                     []AdminServerPath `json:"paths"`
	IsInconsistent            bool              `json:"isInconsistent"`
	IsReserved                bool              `json:"isReserved"`
	HasNonFullmeshLuPaths     bool              `json:"hasNonFullmeshLuPaths"`
	HasUnalignedOsTypes       bool              `json:"hasUnalignedOsTypes"`
	HasUnalignedOsTypeOptions bool              `json:"hasUnalignedOsTypeOptions"`
}

// AdminServerListParams represents parameters for listing servers
type AdminServerListParams struct {
	Nickname  *string `json:"nickname,omitempty"`
	HbaWwn    *string `json:"hbaWwn,omitempty"`
	IscsiName *string `json:"iscsiName,omitempty"`
}

// CreateAdminServerParams represents parameters for creating a server
type CreateAdminServerParams struct {
	ServerNickname string `json:"serverNickname"`
	Protocol       string `json:"protocol,omitempty"`
	OsType         string `json:"osType,omitempty"`
	OsTypeOptions  []int  `json:"osTypeOptions,omitempty"`
	IsReserved     bool   `json:"isReserved"`
}

// UpdateAdminServerParams represents parameters for updating a server
type UpdateAdminServerParams struct {
	Nickname      string `json:"nickname,omitempty"`
	OsType        string `json:"osType,omitempty"`
	OsTypeOptions []int  `json:"osTypeOptions,omitempty"`
}

// DeleteAdminServerParams represents parameters for deleting a server
type DeleteAdminServerParams struct {
	KeepLunConfig bool `json:"keepLunConfig"`
}

// SetAdminServerPathParams represents parameters for setting server path
type SetAdminServerPathParams struct {
	HbaWwn    string   `json:"hbaWwn,omitempty"`
	IscsiName string   `json:"iscsiName,omitempty"`
	PortIds   []string `json:"portIds"`
}

// DeleteAdminServerPathParams represents parameters for deleting server path
type DeleteAdminServerPathParams struct {
	HbaWwn    string `json:"hbaWwn,omitempty"`
	IscsiName string `json:"iscsiName,omitempty"`
	PortId    string `json:"portId"`
}

// AdminServerPathParams represents parameters for getting server path
type AdminServerPathParams struct {
	ServerID  int    `json:"serverId"`
	HbaWwn    string `json:"hbaWwn,omitempty"`
	IscsiName string `json:"iscsiName,omitempty"`
	PortId    string `json:"portId"`
}

// AdminServerPathInfo represents server path information
type AdminServerPathInfo struct {
	ID        string `json:"id"`
	ServerID  int    `json:"serverId"`
	HbaWwn    string `json:"hbaWwn"`
	IscsiName string `json:"iscsiName"`
	PortId    string `json:"portId"`
}
