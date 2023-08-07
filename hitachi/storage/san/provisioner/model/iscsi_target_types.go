package sanstorage

type IscsiTarget struct {
	IscsiTargetID      string      `json:"hostGroupId"`
	PortID             string      `json:"portId"`
	IscsiTargetNumber  int         `json:"hostGroupNumber"`
	IscsiTargetName    string      `json:"hostGroupName"`
	IscsiTargetNameIqn string      `json:"iscsiName"`
	HostMode           string      `json:"hostMode"`
	HostModeOptions    []int       `json:"hostModeOptions"`
	Initiators         []Initiator `json:"initiators"`
	LuPaths            []LuPath    `json:"luPaths"`
	Ldevs              []int       `json:"ldevs"`
	ItLuns             []int       `json:"hgLuns"`
}

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

type Initiator struct {
	IscsiTargetNameIqn string `json:"iscsiName"`
	IscsiNickname      string `json:"iscsiNickname"`
}

type CreateIscsiTargetReq struct {
	PortID             string       `json:"portId,omitempty"`
	IscsiTargetName    string       `json:"hostGroupName,omitempty"`
	Initiators         *[]Initiator `json:"initiators,omitempty"`
	IscsiTargetNumber  *int         `json:"hostGroupNumber,omitempty"`
	IscsiTargetNameIqn *string      `json:"iscsiName"`
	HostModeOptions    *[]int       `json:"hostModeOptions,omitempty"`
	//HP-UX, SOLARIS, AIX, WIN, LINUX/IRIX, TRU64, OVMS, NETWARE, VMWARE, VMWARE_EX, WIN_EX
	HostMode *string      `json:"hostMode,omitempty"`
	Ldevs    *[]IscsiLuns `json:"ldevs,omitempty"`
}

// IscsiLuns .
type IscsiLuns struct {
	LdevId *int `json:"ldevId"`
	Lun    *int `json:"lun,omitempty"`
}

type IscsiNameInformation struct {
	HostIscsiId        string `json:"hostIscsiId,omitempty"`
	PortID             string `json:"portId,omitempty"`
	IscsiTargetNumber  int    `json:"hostGroupNumber,omitempty"`
	IscsiTargetName    string `json:"hostGroupName,omitempty"`
	IscsiTargetNameIqn string `json:"iscsiName,omitempty"`
	IscsiNickname      string `json:"iscsiNickname,omitempty"`
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

// SetIscsiNameReq .
type SetIscsiNameReq struct {
	PortID             string `json:"portId"`
	IscsiTargetNameIqn string `json:"iscsiName"`
	IscsiTargetNumber  int    `json:"hostGroupNumber"`
}

// SetNicknameIscsiReq .
type SetNicknameIscsiReq struct {
	IscsiNickname string `json:"iscsiNickname"`
}
