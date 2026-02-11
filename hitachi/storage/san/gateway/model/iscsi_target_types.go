package sanstorage

type IscsiTargetGwy struct {
	IscsiTargetID      string `json:"hostGroupId"`
	PortID             string `json:"portId"`
	IscsiTargetNumber  int    `json:"hostGroupNumber"`
	HostMode           string `json:"hostMode"`
	HostModeOptions    []int  `json:"hostModeOptions"`
	IscsiTargetName    string `json:"hostGroupName"`
	IscsiTargetNameIqn string `json:"iscsiName"`
	// IscsiTargetDirection string `json:"iscsiTargetDirection"`
	// AuthenticationMode   string `json:"authenticationMode"`
}

type IscsiTargets struct {
	IscsiTargets []IscsiTargetGwy `json:"data"`
}

type LunInfo struct {
	Lun    int `json:"lun"`
	LdevID int `json:"ldevId"`
}

// FIXME - check following json tags
type InitiatorInfo struct {
	InitiatorName     string `json:"initiatorName"`
	InitiatorNickName string `json:"initiatorNickName"`
}

type CreateIscsiTargetReq struct {
	PortID             string  `json:"portId,omitempty"`
	IscsiTargetName    string  `json:"hostGroupName,omitempty"`
	IscsiTargetNumber  *int    `json:"hostGroupNumber,omitempty"`
	IscsiTargetNameIqn *string `json:"iscsiName"`
	HostModeOptions    *[]int  `json:"hostModeOptions,omitempty"`
	//HP-UX, SOLARIS, AIX, WIN, LINUX/IRIX, TRU64, OVMS, NETWARE, VMWARE, VMWARE_EX, WIN_EX
	HostMode *string `json:"hostMode,omitempty"`
}

// SetIscsiNameReq
type SetIscsiNameReq struct {
	PortID             string `json:"portId"`
	IscsiTargetNameIqn string `json:"iscsiName"`
	IscsiTargetNumber  int    `json:"hostGroupNumber"`
}

// SetNicknameIscsiReq
type SetNicknameIscsiReq struct {
	IscsiNickname string `json:"iscsiNickname"`
}

type IscsiNameInformation struct {
	HostIscsiId        string `json:"hostIscsiId,omitempty"`
	PortID             string `json:"portId,omitempty"`
	IscsiTargetNumber  int    `json:"hostGroupNumber,omitempty"`
	IscsiTargetName    string `json:"hostGroupName,omitempty"`
	IscsiTargetNameIqn string `json:"iscsiName,omitempty"`
	IscsiNickname      string `json:"iscsiNickname,omitempty"`
}

type AllIscsiNameInformation struct {
	Data []IscsiNameInformation `json:"data"`
}

// IscsiTargetLuPaths
type IscsiTargetLuPaths struct {
	Data []IscsiTargetLuPath `json:"data"`
}

// IscsiTargetLuPath
type IscsiTargetLuPath struct {
	LunID             string `json:"lunId"`
	PortID            string `json:"portId"`
	IscsiTargetNumber int    `json:"hostGroupNumber"`
	HostMode          string `json:"hostMode"`
	Lun               int    `json:"lun"`
	LdevID            int    `json:"ldevId"`
	IsCommandDevice   bool   `json:"isCommandDevice"`
}

// SetIscsiHostModeAndOptions .
type SetIscsiHostModeAndOptions struct {
	HostMode        string `json:"hostMode"`
	HostModeOptions *[]int `json:"hostModeOptions,omitempty"`
}
