package infra_gw

type StoragePortInfo struct {
	PortId            string `json:"portId"`
	Type              string `json:"type"`
	Speed             string `json:"speed"`
	ResourceGroupId   int    `json:"resourceGroupId"`
	Wwn               string `json:"wwn"`
	ResourceId        string `json:"resourceId"`
	Attribute         string `json:"attribute"`
	ConnectionType    string `json:"connectionType"`
	FabricOn          bool   `json:"fabricOn"`
	Mode              string `json:"mode"`
	IsSecurityEnabled bool   `json:"isSecurityEnabled"`
}

type StoragePorts struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []StoragePortInfo `json:"data"`
}

type StoragePort struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    StoragePortInfo `json:"data"`
}
