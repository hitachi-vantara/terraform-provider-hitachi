package sanstorage

// HostGroup .
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

// HostGroups .
type HostGroups struct {
	HostGroups []HostGroup `json:"data"`
}

// WwnDetail .
type WwnDetail struct {
	Wwn  string `json:"wwn"`
	Name string `json:"name"`
}

// LuPath .
type LuPath struct {
	Lun    int `json:"lun"`
	LdevID int `json:"ldevId"`
}

// CreateHostGroupRequest .
type CreateHostGroupRequest struct {
	PortID          *string `json:"portId,omitempty"`
	HostGroupName   *string `json:"hostGroupName,omitempty"`
	HostGroupNumber *int    `json:"hostGroupNumber,omitempty"`
	HostModeOptions []int   `json:"hostModeOptions,omitempty"`
	HostMode        *string `json:"hostMode,omitempty"`
	Ldevs           []Luns  `json:"ldevs,omitempty"`
	Wwns            []Wwn   `json:"wwns,omitempty"`
}

// Wwn .
type Wwn struct {
	Wwn  string `json:"wwn,omitempty"`
	Name string `json:"name,omitempty"`
}

// Luns .
type Luns struct {
	LdevId *int `json:"ldevId"`
	Lun    *int `json:"lun,omitempty"`
}
