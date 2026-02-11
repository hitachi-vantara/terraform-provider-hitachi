package sanstorage

type HostGroupGwy struct {
	HostGroupID     string `json:"hostGroupId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	HostMode        string `json:"hostMode"`
	HostModeOptions []int  `json:"hostModeOptions"`
}

type HostModeAndOptions struct {
	HostModes       []HostModes       `json:"hostModes"`
	HostModeOptions []HostModeOptions `json:"hostModeOptions"`
}
type HostModes struct {
	HostModeID      int    `json:"hostModeId"`
	HostModeName    string `json:"hostModeName"`
	HostModeDisplay string `json:"hostModeDisplay"`
}
type HostModeOptions struct {
	HostModeOptionID          int    `json:"hostModeOptionId"`
	HostModeOptionDescription string `json:"hostModeOptionDescription"`
}

type HostWwnDetails struct {
	Data []HostWwnDetail `json:"data"`
}
type HostWwnDetail struct {
	HostWwnID       string `json:"hostWwnId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	HostWwn         string `json:"hostWwn"`
	WwnNickname     string `json:"wwnNickname"`
}

type HostLuPaths struct {
	Data []HostLuPath `json:"data"`
}
type HostLuPath struct {
	LunID           string `json:"lunId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostMode        string `json:"hostMode"`
	Lun             int    `json:"lun"`
	LdevID          int    `json:"ldevId"`
	IsCommandDevice bool   `json:"isCommandDevice"`
}

type CreateHostGroupReqGwy struct {
	PortID          *string `json:"portId,omitempty"`
	HostGroupName   *string `json:"hostGroupName,omitempty"`
	HostGroupNumber *int `json:"hostGroupNumber,omitempty"`
	HostModeOptions []int   `json:"hostModeOptions,omitempty"`
	//HP-UX, SOLARIS, AIX, WIN, LINUX/IRIX, TRU64, OVMS, NETWARE, VMWARE, VMWARE_EX, WIN_EX
	HostMode        *string `json:"hostMode,omitempty"`
}

type AddWwnToHgReqGwy struct {
	HostWwn         *string `json:"hostWwn,omitempty"`
	PortID          *string `json:"portId,omitempty"`
	HostGroupNumber *int    `json:"hostGroupNumber,omitempty"`
}

type AddLdevToHgReqGwy struct {
	PortIds         []string `json:"portIds,omitempty"`
	PortID          *string      `json:"portId,omitempty"`
	HostGroupNumber *int      `json:"hostGroupNumber,omitempty"`
	LdevID          *int      `json:"ldevId,omitempty"`
	Lun             *int      `json:"lun,omitempty"`
}