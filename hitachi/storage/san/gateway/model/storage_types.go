package sanstorage

// StorageDeviceSettings is details for Storage settings
type StorageDeviceSettings struct {
	Serial                  int    `json:"serial"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	MgmtIP                  string `json:"mgmtIP"`
	TerraformResourceMethod string `json:"terraformResourceMethod"`
}

// StorageSystemInfo is used for Storage System Informaiton
type StorageSystemInfo struct {
	StorageDeviceID                    string `json:"storageDeviceId"`
	Model                              string `json:"model"`
	SerialNumber                       int    `json:"serialNumber"`
	IP                                 string `json:"ip"`
	Ctl1IP                             string `json:"ctl1Ip"`
	Ctl2IP                             string `json:"ctl2Ip"`
	DkcMicroVersion                    string `json:"dkcMicroVersion"`
	DetailDkcMicroVersion              string `json:"detailDkcMicroVersion,omitempty"`
	IsCompressionAccelerationAvailable bool   `json:"isCompressionAccelerationAvailable,omitempty"`
	CommunicationModes                 []struct {
		CommunicationMode string `json:"communicationMode"`
	} `json:"communicationModes"`
	IsSecure bool `json:"isSecure,omitempty"`
}

// StorageCapacity .
type StorageCapacity struct {
	Internal struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"internal"`
	External struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"external"`
	Total struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"total"`
}
