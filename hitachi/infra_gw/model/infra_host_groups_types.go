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
	Port            string `json:"port,omitempty"`
	HostGroupName   string `json:"hostGroupName,omitempty"`
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

type MTHostGroupInfo struct {
	ResourceId        string `json:"resourceId"`
	Type              string `json:"type"`
	StorageId         string `json:"storageId"`
	DeviceId          string `json:"deviceId"`
	EntitlementStatus string `json:"entitlementStatus"`
	PartnerId         string `json:"partnerId"`
	SubscriberId      string `json:"subscriberId"`

	HostGroupInfo struct {
		HostGroupName   string `json:"hostGroupName"`
		HostGroupId     int    `json:"hostGroupId"`
		ResourceGroupId int    `json:"resourceGroupId"`
		Port            string `json:"port"`
		HostMode        string `json:"hostMode"`
	} `json:"hostGroupInfo"`
}

type MTHostGroups struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Error   TFError           `json:"error"`
	Data    []MTHostGroupInfo `json:"data"`
}

type MTHostGroupDetailsInfo struct {
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
	SubscriberId      string `json:"subscriberId"`
	PartnerId         string `json:"partnerId"`
	EntitlementStatus string `json:"entitlementStatus"`
	StorageId         string `json:"storageId"`
	Time              int64  `json:"time"`
}

type MTHostGroupsDetails struct {
	Path    string                   `json:"path"`
	Message string                   `json:"message"`
	Error   TFError                  `json:"error"`
	Data    []MTHostGroupDetailsInfo `json:"data"`
}

type AddVolumesToHostGroupParam struct {
	LdevIds []int `json:"ldevIds"`
}

type DeleteVolumesToHostGroupParam AddVolumesToHostGroupParam
