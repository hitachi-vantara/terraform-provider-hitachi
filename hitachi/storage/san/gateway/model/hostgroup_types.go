package sanstorage

// HostGroupGwy .
type HostGroupGwy struct {
	HostGroupID     string `json:"hostGroupId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	HostMode        string `json:"hostMode"`
	HostModeOptions []int  `json:"hostModeOptions"`
}

// HostGroups .
type HostGroups struct {
	HostGroups []HostGroupGwy `json:"data"`
}

// HostModeAndOptions .
type HostModeAndOptions struct {
	HostModes       []HostModes       `json:"hostModes"`
	HostModeOptions []HostModeOptions `json:"hostModeOptions"`
}

// SetHostModeAndOptions .
type SetHostModeAndOptions struct {
	HostMode        string `json:"hostMode"`
	HostModeOptions *[]int `json:"hostModeOptions,omitempty"`
}

// HostModes .
type HostModes struct {
	HostModeID      int    `json:"hostModeId"`
	HostModeName    string `json:"hostModeName"`
	HostModeDisplay string `json:"hostModeDisplay"`
}

// HostModeOptions .
type HostModeOptions struct {
	HostModeOptionID          int    `json:"hostModeOptionId"`
	HostModeOptionDescription string `json:"hostModeOptionDescription"`
	Scope                     string `json:"scope"`
	RequiredHostModes         []int  `json:"requiredHostModes"`
}

// HostWwnDetails .
type HostWwnDetails struct {
	Data []HostWwnDetail `json:"data"`
}

// HostWwnDetail .
type HostWwnDetail struct {
	HostWwnID       string `json:"hostWwnId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	HostWwn         string `json:"hostWwn"`
	WwnNickname     string `json:"wwnNickname"`
}

// HostLuPaths .
type HostLuPaths struct {
	Data []HostLuPath `json:"data"`
}

// HostLuPath .
type HostLuPath struct {
	LunID           string `json:"lunId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostMode        string `json:"hostMode"`
	Lun             int    `json:"lun"`
	LdevID          int    `json:"ldevId"`
	IsCommandDevice bool   `json:"isCommandDevice"`
}

// CreateHostGroupReqGwy .
type CreateHostGroupReqGwy struct {
	PortID          *string `json:"portId,omitempty"`
	HostGroupName   *string `json:"hostGroupName,omitempty"`
	HostGroupNumber *int    `json:"hostGroupNumber,omitempty"`
	HostModeOptions []int   `json:"hostModeOptions,omitempty"`
	//HP-UX, SOLARIS, AIX, WIN, LINUX/IRIX, TRU64, OVMS, NETWARE, VMWARE, VMWARE_EX, WIN_EX
	HostMode *string `json:"hostMode,omitempty"`
}

// AddWwnToHgReqGwy .
type AddWwnToHgReqGwy struct {
	HostWwn         *string `json:"hostWwn,omitempty"`
	PortID          *string `json:"portId,omitempty"`
	HostGroupNumber *int    `json:"hostGroupNumber,omitempty"`
}

// AddLdevToHgReqGwy .
type AddLdevToHgReqGwy struct {
	//Specify this attribute when setting the LU paths for multiple ports at the same time. You can specify up to 6 port numbers.
	PortIds []string `json:"portIds,omitempty"`
	// Specify this attribute when setting the LU path for one port.
	PortID *string `json:"portId,omitempty"`
	//(Required) Host group number
	HostGroupNumber *int `json:"hostGroupNumber"`
	//An LDEV cannot be mapped to another LUN in the same host group
	LdevID *int `json:"ldevId"`
	//If this attribute is omitted, a value is automatically set.
	Lun *int `json:"lun,omitempty"`
}

// HostWwnNickName .
type HostWwnNickName struct {
	WwnNickname string `json:"wwnNickname"`
}
