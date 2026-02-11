package admin

import "time"

// AdminPool represents a pool in the admin system
type AdminPool struct {
	ID                           int                     `json:"id"`
	Name                         string                  `json:"name"`
	Status                       string                  `json:"status"`
	Encryption                   string                  `json:"encryption"`
	TotalCapacity                int64                   `json:"totalCapacity"`
	EffectiveCapacity            int64                   `json:"effectiveCapacity"`
	UsedCapacity                 int64                   `json:"usedCapacity"`
	FreeCapacity                 int64                   `json:"freeCapacity"`
	CapacityManage               AdminPoolCapacityManage `json:"capacityManage"`
	SavingEffects                AdminPoolSavingEffects  `json:"savingEffects"`
	ConfigStatus                 []string                `json:"configStatus"`
	NumberOfVolumes              int                     `json:"numberOfVolumes"`
	NumberOfTiers                int                     `json:"numberOfTiers"`
	NumberOfDriveTypes           int                     `json:"numberOfDriveTypes"`
	Tiers                        interface{}             `json:"tiers"`
	Drives                       []AdminPoolDrive        `json:"drives"`
	SubscriptionLimit            interface{}             `json:"subscriptionLimit"`
	ContainsCapacitySavingVolume bool                    `json:"containsCapacitySavingVolume"`
}

// AdminPoolCapacityManage represents capacity management information
type AdminPoolCapacityManage struct {
	UsedCapacityRate   int `json:"usedCapacityRate"`
	ThresholdWarning   int `json:"thresholdWarning"`
	ThresholdDepletion int `json:"thresholdDepletion"`
}

// AdminPoolSavingEffects represents saving effects information
type AdminPoolSavingEffects struct {
	EfficiencyDataReduction               int        `json:"efficiencyDataReduction"`
	EfficiencyFmdSaving                   int        `json:"efficiencyFmdSaving"`
	PreCapacityFmdSaving                  int64      `json:"preCapacityFmdSaving"`
	PostCapacityFmdSaving                 int64      `json:"postCapacityFmdSaving"`
	IsTotalEfficiencySupport              bool       `json:"isTotalEfficiencySupport"`
	TotalEfficiencyStatus                 string     `json:"totalEfficiencyStatus"`
	DataReductionWithoutSystemDataStatus  string     `json:"dataReductionWithoutSystemDataStatus"`
	SoftwareSavingWithoutSystemDataStatus string     `json:"softwareSavingWithoutSystemDataStatus"`
	TotalEfficiency                       int64      `json:"totalEfficiency"`
	DataReductionWithoutSystemData        int64      `json:"dataReductionWithoutSystemData"`
	SoftwareSavingWithoutSystemData       int64      `json:"softwareSavingWithoutSystemData"`
	CalculationStartTime                  *time.Time `json:"calculationStartTime"`
	CalculationEndTime                    *time.Time `json:"calculationEndTime"`
}

// AdminPoolDrive represents drive information
type AdminPoolDrive struct {
	DriveType            string   `json:"driveType"`
	DriveInterface       string   `json:"driveInterface"`
	DriveRpm             string   `json:"driveRpm"`
	DriveCapacity        int      `json:"driveCapacity"`
	DisplayDriveCapacity string   `json:"displayDriveCapacity"`
	TotalCapacity        int64    `json:"totalCapacity"`
	NumberOfDrives       int      `json:"numberOfDrives"`
	Locations            []string `json:"locations"`
	RaidLevel            string   `json:"raidLevel"`
	ParityGroupType      string   `json:"parityGroupType"`
}

// AdminPoolListResponse represents the response for listing pools
type AdminPoolListResponse struct {
	Data  []AdminPool `json:"data"`
	Count int         `json:"count"`
}

// AdminPoolListParams represents parameters for listing pools
type AdminPoolListParams struct {
	Name         *string `json:"name,omitempty"`
	Status       *string `json:"status,omitempty"`
	ConfigStatus *string `json:"configStatus,omitempty"`
}

// CreateAdminPoolParams represents parameters for creating a pool
type CreateAdminPoolParams struct {
	Name                string                 `json:"name"`
	IsEncryptionEnabled bool                   `json:"isEncryptionEnabled"`
	Drives              []CreateAdminPoolDrive `json:"drives"`
}

// CreateAdminPoolDrive represents drive parameters for pool creation
type CreateAdminPoolDrive struct {
	DriveTypeCode   string `json:"driveTypeCode"`
	DataDriveCount  int    `json:"dataDriveCount"`
	RaidLevel       string `json:"raidLevel"`
	ParityGroupType string `json:"parityGroupType"`
}

// UpdateAdminPoolParams represents parameters for updating a pool
type UpdateAdminPoolParams struct {
	Name               string `json:"name,omitempty"`
	ThresholdWarning   int    `json:"thresholdWarning,omitempty"`
	ThresholdDepletion int    `json:"thresholdDepletion,omitempty"`
}

// ExpandAdminPoolParams represents parameters for expanding a pool
type ExpandAdminPoolParams struct {
	AdditionalDrives []ExpandAdminPoolDrive `json:"additionalDrives"`
}

// ExpandAdminPoolDrive represents drive parameters for pool expansion
type ExpandAdminPoolDrive struct {
	DriveTypeCode   string `json:"driveTypeCode"`
	DataDriveCount  int    `json:"dataDriveCount"`
	RaidLevel       string `json:"raidLevel"`
	ParityGroupType string `json:"parityGroupType"`
}
