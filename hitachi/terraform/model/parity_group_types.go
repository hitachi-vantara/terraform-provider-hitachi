package terraform

// Space represents a partition within a parity group
type Space struct {
	PartitionNumber int    `json:"partitionNumber"`
	LdevId          int    `json:"ldevId,omitempty"`
	Status          string `json:"status"`
	LbaLocation     string `json:"lbaLocation"`
	LbaSize         string `json:"lbaSize"`
}

// ParityGroup is group of parity
type ParityGroup struct {
	ParityGroupId                   string  `json:"parityGroupId"`
	GroupType                       string  `json:"groupType"`
	NumberOfLdevs                   int     `json:"numOfLdevs"`
	UsedCapacityRate                int     `json:"usedCapacityRate"`
	AvailableVolumeCapacity         int     `json:"availableVolumeCapacity"`
	RaidLevel                       string  `json:"raidLevel"`
	RaidType                        string  `json:"raidType"`
	ClprId                          int     `json:"clprId"`
	DriveType                       string  `json:"driveType"`
	DriveTypeName                   string  `json:"driveTypeName"`
	IsCopyBackModeEnabled           bool    `json:"isCopyBackModeEnabled"`
	IsEncryptionEnabled             bool    `json:"isEncryptionEnabled"`
	TotalCapacity                   int     `json:"totalCapacity"`
	PhysicalCapacity                int     `json:"physicalCapacity"`
	AvailablePhysicalCapacity       int     `json:"availablePhysicalCapacity"`
	IsAcceleratedCompressionEnabled bool    `json:"isAcceleratedCompressionEnabled"`
	Spaces                          []Space `json:"spaces,omitempty"`
	EmulationType                   string  `json:"emulationType"`
	AvailableVolumeCapacityInKB     int     `json:"availableVolumeCapacityInKB"`
}

// ParityGroups is collection of all parity groups
type ParityGroups struct {
	Data []ParityGroup `json:"data"`
}
