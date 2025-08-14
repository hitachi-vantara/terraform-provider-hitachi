package vssbstorage

import "time"

type StoragePool struct {
	ID                                       string          `json:"id,omitempty"`
	Name                                     string          `json:"name,omitempty"`
	ProtectionDomainId                       string          `json:"protectionDomainId,omitempty"`
	StatusSummary                            string          `json:"statusSummary,omitempty"`
	Status                                   string          `json:"status,omitempty"`
	TotalCapacity                            int             `json:"totalCapacity,omitempty"`
	TotalRawCapacity                         int             `json:"totalRawCapacity,omitempty"`
	UsedCapacity                             int             `json:"usedCapacity,omitempty"`
	FreeCapacity                             int             `json:"freeCapacity,omitempty"`
	TotalPhysicalCapacity                    int             `json:"totalPhysicalCapacity,omitempty"`
	MetaDataPhysicalCapacity                 int             `json:"metaDataPhysicalCapacity,omitempty"`
	ReservedPhysicalCapacity                 int             `json:"reservedPhysicalCapacity,omitempty"`
	UsablePhysicalCapacity                   int             `json:"usablePhysicalCapacity,omitempty"`
	BlockedPhysicalCapacity                  int             `json:"blockedPhysicalCapacity,omitempty"`
	CapacityManage                           CapacityManage  `json:"capacityManage,omitempty"`
	SavingEffect                             SpSavingEffects `json:"savingEffects,omitempty"`
	NumberOfVolumes                          int             `json:"numberOfVolumes,omitempty"`
	RedundantPolicy                          string          `json:"redundantPolicy,omitempty"`
	RedundantType                            string          `json:"redundantType,omitempty"`
	DataRedundancy                           int             `json:"dataRedundancy,omitempty"`
	StorageControllerCapacitiesGeneralStatus string          `json:"storageControllerCapacitiesGeneralStatus,omitempty"`
	TotalVolumeCapacity                      int             `json:"totalVolumeCapacity,omitempty"`
	ProvisionedVolumeCapacity                int             `json:"provisionedVolumeCapacity,omitempty"`
	OtherVolumeCapacity                      int             `json:"otherVolumeCapacity,omitempty"`
	TemporaryVolumeCapacity                  int             `json:"temporaryVolumeCapacity,omitempty"`
	RebuildCapacityPolicy                    string          `json:"rebuildCapacityPolicy,omitempty"`
	RebuildCapacityStatus                    string          `json:"rebuildCapacityStatus,omitempty"`
	RebuildCapacityResourceSetting           struct {
		NumberOfTolerableDriveFailures int `json:"numberOfTolerableDriveFailures,omitempty"`
	} `json:"rebuildCapacityResourceSetting,omitempty"`
	RebuildableResources struct {
		NumberOfDrives int `json:"numberOfDrives,omitempty"`
	} `json:"rebuildableResources,omitempty"`
}

type CapacityManage struct {
	UsedCapacityRate                    int `json:"usedCapacityRate,omitempty"`
	MaximumReserveRate                  int `json:"maximumReserveRate,omitempty"`
	ThresholdWarning                    int `json:"thresholdWarning,omitempty"`
	ThresholdDepletion                  int `json:"thresholdDepletion,omitempty"`
	ThresholdStorageControllerDepletion int `json:"thresholdStorageControllerDepletion,omitempty"`
}

type SpSavingEffects struct {
	EfficiencyDataReduction                    int       `json:"efficiencyDataReduction,omitempty"`
	PreCapacityDataReduction                   int       `json:"preCapacityDataReduction,omitempty"`
	PostCapacityDataReduction                  int       `json:"postCapacityDataReduction,omitempty"`
	TotalEfficiencyStatus                      string    `json:"totalEfficiencyStatus,omitempty"`
	DataReductionWithoutSystemDataStatus       string    `json:"dataReductionWithoutSystemDataStatus,omitempty"`
	TotalEfficiency                            int       `json:"totalEfficiency,omitempty"`
	DataReductionWithoutSystemData             int       `json:"dataReductionWithoutSystemData,omitempty"`
	PreCapacityDataReductionWithoutSystemData  int       `json:"preCapacityDataReductionWithoutSystemData,omitempty"`
	PostCapacityDataReductionWithoutSystemData int       `json:"postCapacityDataReductionWithoutSystemData,omitempty"`
	CalculationStartTime                       time.Time `json:"calculationStartTime,omitempty"`
	CalculationEndTime                         time.Time `json:"calculationEndTime,omitempty"`
}

type StoragePools struct {
	Data []StoragePool `json:"data"`
}

type ExpandStoragePoolReq struct {
	DriveIds []string `json:"driveIds"`
}
