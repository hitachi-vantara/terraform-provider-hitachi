package terraform

// StorageDeviceSettings .
type StorageDeviceSettings struct {
	Serial   int    `json:"serial"`
	Username string `json:"username"`
	Password string `json:"password"`
	MgmtIP   string `json:"mgmtIP"`
}

// StorageSystem .
type StorageSystem struct {
	StorageDeviceID   string `json:"storageDeviceId"`
	Model             string `json:"model"`
	SerialNumber      int    `json:"serialNumber"`
	MgmtIP            string `json:"mgmtIP"`
	SvpIP             string `json:"svpIp"`
	ControllerIP1     string `json:"controllerIP1"`
	ControllerIP2     string `json:"controllerIP2"`
	MicroVersion      string `json:"MicroVersion"`
	TotalCapacityInMB uint64 `json:"totalCapacityInMB"`
	FreeCapacityInMB  uint64 `json:"freeCapacityInMB"`
	UsedCapacityInMB  uint64 `json:"usedCapacityInMB"`
}

// StorageSettingsAndInfo
type StorageSettingsAndInfo struct {
	Settings StorageDeviceSettings `json:"settings"`
	Info     *StorageSystem        `json:"info"`
}

type StorageVersionInfo struct {
	ApiVersion  string `json:"apiVersion"`
	ProductName string `json:"productName"`
}

type AllStorageTypes struct {
	VspStorageSystem       []*StorageSystem
	VssbStorageVersionInfo []*StorageVersionInfo
}
