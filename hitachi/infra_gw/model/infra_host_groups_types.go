package infra_gw

type HostGroupInfo struct {
	ResourceId      string `json:"resourceId"`
	HostGroupName   string `json:"hostGroupName"`
	HostGroupId     int    `json:"hostGroupId"`
	ResourceGroupId int    `json:"resourceGroupId"`
	Port            string `json:"port"`
	LunPaths        []struct {
		LunId  int `json:"lunId"`
		LdevId int `json:"ldevId"`
	} `json:"lunPaths"`
	HostMode string `json:"hostMode"`
	Wwns     []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"wwns"`
	HostModeOptions []struct {
		HostModeOption       string `json:"hostModeOption"`
		HostModeOptionNumber int    `json:"hostModeOptionNumber"`
	} `json:"hostModeOptions"`
}

type HostGroups struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    []HostGroupInfo `json:"data"`
}

type HostGroup struct {
	Path    string        `json:"path"`
	Message string        `json:"message"`
	Data    HostGroupInfo `json:"data"`
}

type CreateHostGroupParam struct {
	Port          string `json:"port,omitempty"`
	HostGroupName string `json:"hostGroupName,omitempty"`
	//HostGroupNumber int    `json:"hostGroupNumber,omitempty"`
	HostModeOptions []int  `json:"hostModeOptions,omitempty"`
	HostMode        string `json:"hostMode,omitempty"`
	Luns            []int  `json:"ldevs,omitempty"`
	Wwns            []Wwn  `json:"wwns,omitempty"`
	UcpSystem       string `json:"ucpSystem,omitempty"`
}

// Luns .
type Luns struct {
	LdevId *int `json:"ldevId"`
	Lun    *int `json:"lun,omitempty"`
}

type Wwn struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PatcheHostGroupParam struct {
	HostMode        string `json:"hostMode,omitempty"`
	HostModeOptions []int  `json:"hostModeOptions,omitempty"`
}
