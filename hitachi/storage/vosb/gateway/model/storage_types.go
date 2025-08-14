package vssbstorage

type StorageDeviceSettings struct {
	Username                string `json:"username"`
	Password                string `json:"password"`
	ClusterAddress          string `json:"clusterAddress"`
	TerraformResourceMethod string `json:"terraformResourceMethod"`
	ContentType 			string `json:"contentType"`
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
	TotalPoolRawCapacityInMB      uint64                `json:"totalPoolRawCapacity"`
	TotalPoolPhysicalCapacityInMB uint64                `json:"totalPoolPhysicalCapacity"`
	TotalPoolCapacityInMB         uint64                `json:"totalPoolCapacity"`
	UsedPoolCapacityInMB          uint64                `json:"usedPoolCapacity"`
	FreePoolCapacityInMB          uint64                `json:"freePoolCapacity"`
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

/*
type CapacityManage struct {

}

type Pool struct {
	Id                            string `json:"id"`
	Name                          string `json:"name"`
	ProtectionDomainId            string `json:"protectionDomainId"`
	StatusSummary                 string `json:"statusSummary"`
	Status                        string `json:"status"`
	TotalCapacityInMB             uint64 `json:"totalCapacity"`
	TotalRawCapacityInMB          uint64 `json:"totalRawCapacity"`
	UsedCapacityInMB              uint64 `json:"usedCapacity"`
	FreeCapacityInMB              uint64 `json:"freeCapacity"`
	TotalPhysicalCapacityInMB     uint64 `json:"totalPhysicalCapacity"`
	MetaDataPhysicalCapacityInMB  uint64 `json:"metaDataPhysicalCapacity"`
	ReservedPhysicalCapacityInMB  uint64 `json:"reservedPhysicalCapacity"`
	UsablePhysicalCapacityInMB    uint64 `json:"usablePhysicalCapacity"`
	BlockedPhysicalCapacityInMB   uint64 `json:"blockedPhysicalCapacity"`
	TotalVolumeCapacityInMB       uint64 `json:"totalVolumeCapacity"`
	ProvisionedVolumeCapacityInMB uint64 `json:"provisionedVolumeCapacity"`
	OtherVolumeCapacityInMB       uint64 `json:"otherVolumeCapacity"`
	TemporaryVolumeCapacityInMB   uint64 `json:"temporaryVolumeCapacity"`
}
*/
