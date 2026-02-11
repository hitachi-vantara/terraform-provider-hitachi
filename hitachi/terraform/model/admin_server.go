package terraform

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

type AdminServerListResponse struct {
	Data  []AdminServer `json:"data"`
	Count int           `json:"count"`
}

type AdminServerPath struct {
	HbaWwn    string   `json:"hbaWwn"`
	IscsiName string   `json:"iscsiName"`
	PortIds   []string `json:"portIds"`
}

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
