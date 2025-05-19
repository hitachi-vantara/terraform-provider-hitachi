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

type Drive struct {
	Id               string `json:"id"`
	WwwId            string `json:"wwwid"`
	StatusSummary    string `json:"statusSummary"`
	Status           string `json:"status"`
	TypeCode         string `json:"typeCode"`
	SerialNumber     string `json:"serialNumber"`
	StorageNodeId    string `json:"storageNodeId"`
	DeviceFileName   string `json:"deviceFileName"`
	VendorName       string `json:"vendorName"`
	FirmwareRevision string `json:"firmwareRevision"`
	LocatorLedStatus string `json:"locatorLedStatus"`
	DriveType        string `json:"driveType"`
	DriveCapacity    int    `json:"driveCapacity"`
}

type Drives struct {
	Data []Drive `json:"data"`
}
