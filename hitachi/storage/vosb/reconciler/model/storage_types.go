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
