package sanstorage

// LogicalUnit .
type LogicalUnit struct {
	// these are returned from gateway
	LdevID              int    `json:"ldevId"`
	VirtualLdevID       int64  `json:"virtualLdevId,omitempty"`
	VirtualSerialNumber string `json:"virtualSerialNumber,omitempty"`
	VirtualModel        string `json:"virtualModel,omitempty"`
	ClprID              int    `json:"clprId"`
	EmulationType       string `json:"emulationType"`
	ByteFormatCapacity  string `json:"byteFormatCapacity"`
	BlockCapacity       uint64 `json:"blockCapacity"`
	Cylinder            int    `json:"cylinder,omitempty"`
	NumOfPorts          int    `json:"numOfPorts"`
	Ports               []struct {
		PortID          string `json:"portId"`
		HostGroupNumber int    `json:"hostGroupNumber"`
		HostGroupName   string `json:"hostGroupName"`
		Lun             int    `json:"lun"`
	} `json:"ports"`
	Attributes                       []string `json:"attributes"`
	Label                            string   `json:"label"`
	Status                           string   `json:"status"`
	MpBladeID                        int      `json:"mpBladeId"`
	Ssid                             string   `json:"ssid"`
	PoolID                           int      `json:"poolId"`
	NumOfUsedBlock                   uint64   `json:"numOfUsedBlock"`
	IsFullAllocationEnabled          bool     `json:"isFullAllocationEnabled"`
	ResourceGroupID                  int      `json:"resourceGroupId"`
	DataReductionStatus              string   `json:"dataReductionStatus"`
	DataReductionMode                string   `json:"dataReductionMode"`
	DataReductionProcessMode         string   `json:"dataReductionProcessMode"`
	DataReductionProgressRate        int      `json:"dataReductionProgressRate"`
	IsAluaEnabled                    bool     `json:"isAluaEnabled"`
	NaaID                            string   `json:"naaId"`
	IsCompressionAccelerationEnabled bool     `json:"isCompressionAccelerationEnabled"`
	CompressionAccelerationStatus    string   `json:"compressionAccelerationStatus"`

	// RAID-related
	RaidLevel               string   `json:"raidLevel,omitempty"`
	RaidType                string   `json:"raidType,omitempty"`
	NumOfParityGroups       int      `json:"numOfParityGroups,omitempty"`
	ParityGroupIds          []string `json:"parityGroupIds,omitempty"`
	DriveType               string   `json:"driveType,omitempty"`
	DriveByteFormatCapacity string   `json:"driveByteFormatCapacity,omitempty"`
	DriveBlockCapacity      int64    `json:"driveBlockCapacity,omitempty"`

	// ⚡️ Newly added fields
	ComposingPoolId int `json:"composingPoolId,omitempty"`

	SnapshotPoolId int `json:"snapshotPoolId,omitempty"`

	ExternalVendorId       string `json:"externalVendorId,omitempty"`
	ExternalProductId      string `json:"externalProductId,omitempty"`
	ExternalVolumeId       string `json:"externalVolumeId,omitempty"`
	ExternalVolumeIdString string `json:"externalVolumeIdString,omitempty"`

	NumOfExternalPorts int `json:"numOfExternalPorts,omitempty"`
	ExternalPorts      []struct {
		PortID          string `json:"portId"`
		HostGroupNumber int    `json:"hostGroupNumber"`
		Lun             int    `json:"lun"`
		Wwn             string `json:"wwn"`
	} `json:"externalPorts,omitempty"`

	QuorumDiskId              int    `json:"quorumDiskId,omitempty"`
	QuorumStorageSerialNumber string `json:"quorumStorageSerialNumber,omitempty"`
	QuorumStorageTypeId       string `json:"quorumStorageTypeId,omitempty"`

	NamespaceID    string `json:"namespaceId,omitempty"`
	NvmSubsystemId string `json:"nvmSubsystemId,omitempty"`

	IsRelocationEnabled bool   `json:"isRelocationEnabled,omitempty"`
	TierLevel           string `json:"tierLevel,omitempty"`

	UsedCapacityPerTierLevel1 int64 `json:"usedCapacityPerTierLevel1,omitempty"`
	UsedCapacityPerTierLevel2 int64 `json:"usedCapacityPerTierLevel2,omitempty"`
	UsedCapacityPerTierLevel3 int64 `json:"usedCapacityPerTierLevel3,omitempty"`

	TierLevelForNewPageAllocation string `json:"tierLevelForNewPageAllocation,omitempty"`

	OperationType                  string `json:"operationType,omitempty"`
	PreparingOperationProgressRate int    `json:"preparingOperationProgressRate,omitempty"`

	ParentLdevId int `json:"parentLdevId,omitempty"`

	// below will be populated by provisioner
	TotalCapacityInMB uint64 `json:"totalCapacityInMB"`
	FreeCapacityInMB  uint64 `json:"freeCapacityInMB"`
	UsedCapacityInMB  uint64 `json:"usedCapacityInMB"`
}

// LogicalUnits .
type LogicalUnits struct {
	ListOfLun []LogicalUnit `json:"data,omitempty"`
}

// CreateLunRequestGwy .
type CreateLunRequestGwy struct {
	Label  *string `json:"label,omitempty"`
	LdevID *int    `json:"ldevId,omitempty"`
	// Mainframe-specific fields
	Cylinder                           *int    `json:"cylinder,omitempty"`
	EmulationType                      *string `json:"emulationType,omitempty"`
	IsTseVolume                        *bool   `json:"isTseVolume,omitempty"`
	IsEseVolume                        *bool   `json:"isEseVolume,omitempty"`
	ClprID                             *int    `json:"clprId,omitempty"`
	MpBladeID                          *int    `json:"mpBladeId,omitempty"`
	Ssid                               *string `json:"ssid,omitempty"`
	BlockCapacity                      *int64  `json:"blockCapacity,omitempty"`
	PoolID                             *int    `json:"poolId,omitempty"`
	ParityGroupID                      *string `json:"parityGroupId,omitempty"`
	ExternalParityGroupID              *string `json:"externalParityGroupId,omitempty"`
	ByteFormatCapacity                 string  `json:"byteFormatCapacity,omitempty"`
	DataReductionMode                  *string `json:"dataReductionMode,omitempty"`
	IsDataReductionSharedVolumeEnabled *bool   `json:"isDataReductionSharedVolumeEnabled,omitempty"`
	IsCompressionAccelerationEnabled   *bool   `json:"isCompressionAccelerationEnabled,omitempty"`
}

// UpdateLunRequestGwy .
type UpdateLunRequestGwy struct {
	Label                            *string `json:"label,omitempty"`
	DataReductionMode                *string `json:"dataReductionMode,omitempty"`
	DataReductionProcessMode         *string `json:"dataReductionProcessMode,omitempty"`
	IsCompressionAccelerationEnabled *bool   `json:"isCompressionAccelerationEnabled,omitempty"`
	IsAluaEnabled                    *bool   `json:"isAluaEnabled,omitempty"`
}

// ExpandLunRequestGwy .
type ExpandLunRequestGwy struct {
	Parameters ExpandLunParameters `json:"parameters,omitempty"`
}

// ExpandLunParameters .
type ExpandLunParameters struct {
	AdditionalByteFormatCapacity string `json:"additionalByteFormatCapacity,omitempty"`
	// AdditionalBlockCapacity      int    `json:"additionalBlockCapacity,omitempty"`
	// EnhancedExpansion            bool   `json:"enhancedExpansion,omitempty"`
}

// SetEseVolumeRequestGwy is the request body for the Set ESE action
type SetEseVolumeRequestGwy struct {
	Parameters SetEseVolumeParameters `json:"parameters"`
}

type SetEseVolumeParameters struct {
	IsEseVolume bool `json:"isEseVolume"`
}

// FormatLdevRequestGwy is the request body for invoking the LDEV format action
type FormatLdevRequestGwy struct {
	Parameters FormatLdevParameters `json:"parameters,omitempty"`
}

type FormatLdevParameters struct {
	IsDataReductionForceFormat bool   `json:"isDataReductionForceFormat,omitempty"`
	OperationType              string `json:"operationType,omitempty"` // FMT or QFMT
}

// ChangeStatusRequestGwy is the request body for invoking the change-status action
type ChangeStatusRequestGwy struct {
	Parameters ChangeStatusParameters `json:"parameters,omitempty"`
}

type ChangeStatusParameters struct {
	Status string `json:"status,omitempty"`
}
