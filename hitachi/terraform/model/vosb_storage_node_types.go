package terraform

type StorageNodeVssb struct {
	ID                 string `json:"id"`
	BiosUuid           string `json:"biosUuid"`
	ProtectionDomainID string `json:"protectionDomainId"`
	FaultDomainID      string `json:"faultDomainId"`
	FaultDomainName    string `json:"faultDomainName"`
	Name               string `json:"name"`
	ClusterRole        string `json:"clusterRole"`
	// storageNodeAttributes                  list `json:"storageNodeAttributes"`
	StatusSummary             string `json:"statusSummary"`
	Status                    string `json:"status"`
	DriveDataRelocationStatus string `json:"driveDataRelocationStatus"`
	ControlPortIpv4Address    string `json:"controlPortIpv4Address"`
	InternodePortIpv4Address  string `json:"internodePortIpv4Address"`
	SoftwareVersion           string `json:"softwareVersion"`
	ModelName                 string `json:"modelName"`
	SerialNumber              string `json:"serialNumber"`
	Memory                    int    `json:"memory"`
	AvailabilityZoneID        string `json:"availabilityZoneId"`

	InsufficientResourcesForRebuildCapacity struct {
		CapacityOfDrive int `json:"capacityOfDrive"`
		NumberOfDrives  int `json:"numberOfDrives"`
	} `json:"insufficientResourcesForRebuildCapacity"`

	RebuildableResources struct {
		NumberOfDrives int `json:"numberOfDrives"`
	} `json:"rebuildableResources"`
}

type StorageNodesVssb struct {
	Data []StorageNodeVssb `json:"data"`
}
