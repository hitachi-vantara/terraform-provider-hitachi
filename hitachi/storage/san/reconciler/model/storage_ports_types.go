package sanstorage

type StoragePort struct {
	PortId             string   `json:"portId"`
	PortType           string   `json:"portType"`
	PortAttributes     []string `json:"portAttributes"`
	PortSpeed          string   `json:"portSpeed"`
	LoopId             string   `json:"loopId"`
	FabricMode         bool     `json:"fabricMode"`
	PortConnection     string   `json:"portConnection"`
	LunSecuritySetting bool     `json:"lunSecuritySetting"`
	Wwn                string   `json:"wwn"`
	PortMode           string   `json:"portMode"`
}
