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

type IscsiTargets struct {
	IscsiTargets []IscsiTarget `json:"data"`
}

type LunInfo struct {
	Lun    int `json:"lun"`
	LdevID int `json:"ldevId"`
}

type Initiator struct {
	IscsiTargetNameIqn string `json:"iscsiName,omitempty"`
	IscsiNickname      string `json:"iscsiNickname,omitempty"`
}

// CreateIscsiTargetReq .
type CreateIscsiTargetReq struct {
	PortID             string  `json:"portId,omitempty"`
	IscsiTargetName    string  `json:"hostGroupName,omitempty"`
	IscsiTargetNumber  *int    `json:"hostGroupNumber,omitempty"`
	IscsiTargetNameIqn *string `json:"iscsiName"`
	HostModeOptions    *[]int  `json:"hostModeOptions,omitempty"`
	//HP-UX, SOLARIS, AIX, WIN, LINUX/IRIX, TRU64, OVMS, NETWARE, VMWARE, VMWARE_EX, WIN_EX
	HostMode   *string     `json:"hostMode,omitempty"`
	Ldevs      []IscsiLuns `json:"ldevs,omitempty"`
	Initiators []Initiator `json:"initiators,omitempty"` // same as wwn
}

// IscsiLuns .
type IscsiLuns struct {
	LdevId *int `json:"ldevId"`
	Lun    *int `json:"lun,omitempty"`
}
