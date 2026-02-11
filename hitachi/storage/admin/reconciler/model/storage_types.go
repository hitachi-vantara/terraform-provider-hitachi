package admin

type StorageDeviceSettings struct {
	Serial                  int    `json:"serial"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	MgmtIP                  string `json:"mgmtIP"`
	TerraformResourceMethod string `json:"terraformResourceMethod"`
}

// SavingEffectsAdmin
type SavingEffects struct {
	EfficiencyDataReduction                    int    `json:"efficiencyDataReduction,omitempty"`
	PreCapacityDataReduction                   int    `json:"preCapacityDataReduction,omitempty"`
	PostCapacityDataReduction                  int    `json:"postCapacityDataReduction,omitempty"`
	EfficiencyFmdSaving                        int    `json:"efficiencyFmdSaving,omitempty"`
	PreCapacityFmdSaving                       int    `json:"preCapacityFmdSaving,omitempty"`
	PostCapacityFmdSaving                      int    `json:"postCapacityFmdSaving,omitempty"`
	IsTotalEfficiencySupport                   bool   `json:"isTotalEfficiencySupport,omitempty"`
	TotalEfficiencyStatus                      string `json:"totalEfficiencyStatus,omitempty"`
	DataReductionWithoutSystemDataStatus       string `json:"dataReductionWithoutSystemDataStatus,omitempty"`
	SoftwareSavingWithoutSystemDataStatus      string `json:"softwareSavingWithoutSystemDataStatus,omitempty"`
	TotalEfficiency                            int    `json:"totalEfficiency,omitempty"`
	DataReductionWithoutSystemData             int    `json:"dataReductionWithoutSystemData,omitempty"`
	PreCapacityDataReductionWithoutSystemData  int    `json:"preCapacityDataReductionWithoutSystemData,omitempty"`
	PostCapacityDataReductionWithoutSystemData int    `json:"postCapacityDataReductionWithoutSystemData,omitempty"`
	SoftwareSavingWithoutSystemData            int    `json:"softwareSavingWithoutSystemData,omitempty"`
	CalculationStartTime                       string `json:"calculationStartTime,omitempty"`
	CalculationEndTime                         string `json:"calculationEndTime,omitempty"`
}

// StorageAdminInfo represents the admin information
type StorageAdminInfo struct {
	ModelName                           string        `json:"modelName,omitempty"`
	Serial                              string        `json:"serial,omitempty"`
	Nickname                            string        `json:"nickname,omitempty"`
	NumberOfTotalVolumes                int           `json:"numberOfTotalVolumes,omitempty"`
	NumberOfFreeDrives                  int           `json:"numberOfFreeDrives,omitempty"`
	NumberOfTotalServers                int           `json:"numberOfTotalServers,omitempty"`
	TotalPhysicalCapacity               int64         `json:"totalPhysicalCapacity,omitempty"`
	TotalPoolCapacity                   int64         `json:"totalPoolCapacity,omitempty"`
	TotalPoolPhysicalCapacity           int64         `json:"totalPoolPhysicalCapacity,omitempty"`
	UsedPoolCapacity                    int64         `json:"usedPoolCapacity,omitempty"`
	FreePoolCapacity                    int64         `json:"freePoolCapacity,omitempty"`
	TotalPoolCapacityWithTiPool         int64         `json:"totalPoolCapacityWithTiPool,omitempty"`
	TotalPoolPhysicalCapacityWithTiPool int64         `json:"totalPoolPhysicalCapacityWithTiPool,omitempty"`
	UsedPoolCapacityWithTiPool          int64         `json:"usedPoolCapacityWithTiPool,omitempty"`
	FreePoolCapacityWithTiPool          int64         `json:"freePoolCapacityWithTiPool,omitempty"`
	EstimatedConfigurablePoolCapacity   int64         `json:"estimatedConfigurablePoolCapacity,omitempty"`
	EstimatedConfigurableVolumeCapacity int64         `json:"estimatedConfigurableVolumeCapacity,omitempty"`
	SavingEffects                       SavingEffects `json:"savingEffects,omitempty"`
	GumVersion                          string        `json:"gumVersion,omitempty"`
	EsmOsVersion                        string        `json:"esmOsVersion,omitempty"`
	DkcMicroVersion                     string        `json:"dkcMicroVersion,omitempty"`
	WarningLedStatus                    string        `json:"warningLedStatus,omitempty"`
	EsmStatus                           string        `json:"esmStatus,omitempty"`
	IpAddressIpv4Service                string        `json:"ipAddressIpv4Service,omitempty"`
	IpAddressIpv4Ctl1                   string        `json:"ipAddressIpv4Ctl1,omitempty"`
	IpAddressIpv4Ctl2                   string        `json:"ipAddressIpv4Ctl2,omitempty"`
	IpAddressIpv6Service                string        `json:"ipAddressIpv6Service,omitempty"`
	IpAddressIpv6Ctl1                   string        `json:"ipAddressIpv6Ctl1,omitempty"`
	IpAddressIpv6Ctl2                   string        `json:"ipAddressIpv6Ctl2,omitempty"`
}

// StorageSettingsAndInfo
type StorageSettingsAndInfo struct {
	Settings StorageDeviceSettings `json:"settings"`
	Info     *StorageAdminInfo     `json:"info"`
}
