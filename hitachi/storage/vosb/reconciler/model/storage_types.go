package vssbstorage

// StorageDeviceSettings .
type StorageDeviceSettings struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ClusterAddress string `json:"clusterAddress"`
}

type StorageSettingsAndInfo struct {
	Settings StorageDeviceSettings `json:"settings"`
	Info     *StorageVersionInfo   `json:"info"`
}

type StorageVersionInfo struct {
	ApiVersion  string `json:"apiVersion"`
	ProductName string `json:"productName"`
}

type Dashboard struct {
	HealthStatuses            []HealthStatus `json:"data"`
	NumberOfComputePorts      int            `json:"numberOfComputePorts"`
	NumberOfTotalVolumes      int            `json:"numberOfTotalVolumes"`
	NumberOfTotalServers      int            `json:"numberOfTotalServers"`
	NumberOfTotalStorageNodes int            `json:"numberOfTotalStorageNodes"`
	NumberOfStoragePools      int            `json:"numberOfStoragePools"`
	NumberOfFaultDomains      int            `json:"numberOfFaultDomains"`
	NumberOfDrives            int            `json:"numberOfDrives"`
	TotalPoolCapacityInMB     uint64         `json:"totalPoolCapacityInMB"`
	UsedPoolCapacityInMB      uint64         `json:"usedPoolCapacityInMB"`
	FreePoolCapacityInMB      uint64         `json:"freePoolCapacityInMB"`
	TotalEfficiency           int64          `json:"totalEfficiency"`
	EfficiencyDataReduction   int            `json:"efficiencyDataReduction"`
}

type HealthStatus struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type Drive struct {
	Id               string `json:"id"`
	WwId             string `json:"wwid"`
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
