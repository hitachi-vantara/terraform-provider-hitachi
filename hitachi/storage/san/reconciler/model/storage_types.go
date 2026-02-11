package sanstorage

// StorageDeviceSettings .
type StorageDeviceSettings struct {
	Serial   int    `json:"serial"`
	Username string `json:"username"`
	Password string `json:"password"`
	MgmtIP   string `json:"mgmtIP"`
}

// StorageSystem .
type StorageSystem struct {
	StorageDeviceID                    string `json:"storageDeviceId"`
	Model                              string `json:"model"`
	SerialNumber                       int    `json:"serialNumber"`
	MgmtIP                             string `json:"mgmtIP"`
	IP                                 string `json:"ip"`
	ControllerIP1                      string `json:"controllerIP1"`
	ControllerIP2                      string `json:"controllerIP2"`
	MicroVersion                       string `json:"MicroVersion"`
	DetailDkcMicroVersion              string `json:"detailDkcMicroVersion,omitempty"`
	IsCompressionAccelerationAvailable bool   `json:"isCompressionAccelerationAvailable,omitempty"`
	IsSecure                           bool   `json:"isSecure,omitempty"`
	TotalCapacityInMB                  uint64 `json:"totalCapacityInMB"`
	FreeCapacityInMB                   uint64 `json:"freeCapacityInMB"`
	UsedCapacityInMB                   uint64 `json:"usedCapacityInMB"`
}

// StorageSettingsAndInfo
type StorageSettingsAndInfo struct {
	Settings StorageDeviceSettings `json:"settings"`
	Info     *StorageSystem        `json:"info"`
}
