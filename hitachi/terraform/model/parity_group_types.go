package terraform

// ParityGroup is group of parity
type ParityGroup struct {
	ParityGroupId                   string `json:"parityGroupId"`
	NumberOfLdevs                   int    `json:"numOfLdevs"`
	UsedCapacityRate                int    `json:"usedCapacityRate"`
	AvailableVolumeCapacity         int    `json:"availableVolumeCapacity"`
	RaidLevel                       string `json:"raidLevel"`
	RaidType                        string `json:"raidType"`
	ClprId                          int    `json:"clprId"`
	DriveType                       string `json:"driveType"`
	DriveTypeName                   string `json:"driveTypeName"`
	TotalCapacity                   int    `json:"totalCapacity"`
	PhysicalCapacity                int    `json:"physicalCapacity"`
	AvailablePhysicalCapacity       int    `json:"availablePhysicalCapacity"`
	IsAcceleratedCompressionEnabled bool   `json:"isAcceleratedCompressionEnabled"`
	AvailableVolumeCapacityInKB     int    `json:"availableVolumeCapacityInKB"`
}

// ParityGroups is collection of all parity groups
type ParityGroups struct {
	Data []ParityGroup `json:"data"`
}
