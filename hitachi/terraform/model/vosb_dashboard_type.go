package terraform

type Dashboard struct {
	HealthStatuses            []HealthStatus `json:"data"`
	NumberOfTotalVolumes      int            `json:"numberOfTotalVolumes"`
	NumberOfTotalServers      int            `json:"numberOfTotalServers"`
	NumberOfComputePorts      int            `json:"numberOfComputePorts"`
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
