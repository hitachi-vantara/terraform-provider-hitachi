package sanstorage

type HostGroup struct {
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	HostMode        string `json:"hostMode"`
	HostModeOptions []int  `json:"hostModeOptions"`

	WwnDetails []WwnDetail `json:"wwnDetails"`
	LuPaths    []LuPath    `json:"luPaths"`

	Wwns   []string `json:"wwns"`
	Ldevs  []int    `json:"ldevs"`
	HgLuns []int    `json:"hgLuns"`
}

type WwnDetail struct {
	Wwn  string `json:"wwn"`
	Name string `json:"name"`
}

type LuPath struct {
	Lun    int `json:"lun"`
	LdevID int `json:"ldevId"`
}

type CreateHostGroupRequest struct {
	PortID          *string `json:"portId,omitempty"`
	HostGroupName   *string `json:"hostGroupName,omitempty"`
	HostGroupNumber *int    `json:"hostGroupNumber,omitempty"`
	HostModeOptions []int   `json:"hostModeOptions,omitempty"`
	HostMode        *string `json:"hostMode,omitempty"`
	Ldevs []int `json:"ldevs,omitempty"`
	Wwns  []Wwn `json:"wwns,omitempty"`
}

type Wwn struct {
	Wwn  string `json:"wwn,omitempty"`
	Name string `json:"name,omitempty"`
}
