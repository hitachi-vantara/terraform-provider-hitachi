package vssbstorage

// StorageDeviceSettings .
type StorageDeviceSettings struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ClusterAddress string `json:"clusterAddress"`
}

type StorageVersionInfo struct {
	ApiVersion  string `json:"apiVersion"`
	ProductName string `json:"productName"`
}

type HealthStatus struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	ProtectionDomainId string `json:"protectionDomainId"`
}

type HealthStatuses map[string][]HealthStatus

type FaultDomain struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	StatusSummary        string `json:"statusSummary"`
	Status               string `json:"status"`
	NumberOfStorageNodes int    `json:"numberOfStorageNodes"`
}

type FaultDomains struct {
	Data []FaultDomain `json:"data"`
}

type SavingEffectOfStorage struct {
	EfficiencyDataReduction                    int    `json:"efficiencyDataReduction"`
	PreCapacityDataReduction                   int64  `json:"preCapacityDataReduction"`
	PostCapacityDataReduction                  int64  `json:"postCapacityDataReduction"`
	TotalEfficiencyStatus                      string `json:"totalEfficiencyStatus"`
	DataReductionWithoutSystemDataStatus       string `json:"dataReductionWithoutSystemDataStatus"`
	TotalEfficiency                            int64  `json:"totalEfficiency"`
	DataReductionWithoutSystemData             int64  `json:"dataReductionWithoutSystemData"`
	PreCapacityDataReductionWithoutSystemData  int64  `json:"preCapacityDataReductionWithoutSystemData"`
	PostCapacityDataReductionWithoutSystemData int64  `json:"postCapacityDataReductionWithoutSystemData"`
	CalculationStartTime                       string `json:"calculationStartTime"`
	CalculationEndTime                         string `json:"calculationEndTime"`
}

type StorageClusterInfo struct {
	StorageDeviceId               string                `json:"storageDeviceId"`
	Id                            string                `json:"id"`
	ModelName                     string                `json:"modelName"`
	InternalId                    string                `json:"internalId"`
	Nickname                      string                `json:"nickname"`
	NumberOfTotalVolumes          int                   `json:"numberOfTotalVolumes"`
	NumberOfTotalServers          int                   `json:"numberOfTotalServers"`
	NumberOfTotalStorageNodes     int                   `json:"numberOfTotalStorageNodes"`
	NumberOfReadyStorageNodes     int                   `json:"numberOfReadyStorageNodes"`
	NumberOfFaultDomains          int                   `json:"numberOfFaultDomains"`
	TotalPoolRawCapacityInMB      int                   `json:"totalPoolRawCapacityInMB"`
	TotalPoolPhysicalCapacityInMB uint64                `json:"totalPoolPhysicalCapacityInMB"`
	TotalPoolCapacityInMB         uint64                `json:"totalPoolCapacityInMB"`
	UsedPoolCapacityInMB          uint64                `json:"usedPoolCapacityInMB"`
	FreePoolCapacityInMB          uint64                `json:"freePoolCapacityInMB"`
	SavingEffects                 SavingEffectOfStorage `json:"savingEffects"`
	SoftwareVersion               string                `json:"softwareVersion"`
	StatusSummary                 string                `json:"statusSummary"`
	Status                        string                `json:"status"`
	SystemRequirementsFileVersion int                   `json:"systemRequirementsFileVersion"`
	ServiceId                     string                `json:"serviceId"`
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
