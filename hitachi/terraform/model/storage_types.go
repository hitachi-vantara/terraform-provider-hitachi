package terraform

// StorageDeviceSettings .
type StorageDeviceSettings struct {
	Serial   int    `json:"serial"`
	Username string `json:"username"`
	Password string `json:"password"`
	MgmtIP   string `json:"mgmtIP"`
}

// StorageSystem .
type StorageSystem struct {
	StorageDeviceID                    string              `json:"storageDeviceId"`
	Model                              string              `json:"model"`
	SerialNumber                       int                 `json:"serialNumber"`
	MgmtIP                             string              `json:"mgmtIP"`
	IP                                 string              `json:"ip"`
	ControllerIP1                      string              `json:"controllerIP1"`
	ControllerIP2                      string              `json:"controllerIP2"`
	MicroVersion                       string              `json:"MicroVersion"`
	TotalCapacityInMB                  uint64              `json:"totalCapacityInMB"`
	FreeCapacityInMB                   uint64              `json:"freeCapacityInMB"`
	UsedCapacityInMB                   uint64              `json:"usedCapacityInMB"`
	IsCompressionAccelerationAvailable bool                `json:"isCompressionAccelerationAvailable,omitempty"`
	DetailDkcMicroVersion              string              `json:"detailDkcMicroVersion,omitempty"`
	CommunicationModes                 []CommunicationMode `json:"communicationModes,omitempty"`
	IsSecure                           bool                `json:"isSecure,omitempty"`
}

// CommunicationMode for storage system communication modes
type CommunicationMode struct {
	CommunicationMode string `json:"communicationMode"`
}

// SavingEffectsAdmin
type SavingEffectsAdmin struct {
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

// StorageSystemAdmin
type StorageSystemAdmin struct {
	ModelName                           string              `json:"modelName,omitempty"`
	Serial                              string              `json:"serial,omitempty"`
	Nickname                            string              `json:"nickname,omitempty"`
	NumberOfTotalVolumes                int                 `json:"numberOfTotalVolumes,omitempty"`
	NumberOfFreeDrives                  int                 `json:"numberOfFreeDrives,omitempty"`
	NumberOfTotalServers                int                 `json:"numberOfTotalServers,omitempty"`
	TotalPhysicalCapacity               int                 `json:"totalPhysicalCapacity,omitempty"`
	TotalPoolCapacity                   int                 `json:"totalPoolCapacity,omitempty"`
	TotalPoolPhysicalCapacity           int                 `json:"totalPoolPhysicalCapacity,omitempty"`
	UsedPoolCapacity                    int                 `json:"usedPoolCapacity,omitempty"`
	FreePoolCapacity                    int                 `json:"freePoolCapacity,omitempty"`
	TotalPoolCapacityWithTiPool         int                 `json:"totalPoolCapacityWithTiPool,omitempty"`
	TotalPoolPhysicalCapacityWithTiPool int                 `json:"totalPoolPhysicalCapacityWithTiPool,omitempty"`
	UsedPoolCapacityWithTiPool          int                 `json:"usedPoolCapacityWithTiPool,omitempty"`
	FreePoolCapacityWithTiPool          int                 `json:"freePoolCapacityWithTiPool,omitempty"`
	EstimatedConfigurablePoolCapacity   int                 `json:"estimatedConfigurablePoolCapacity,omitempty"`
	EstimatedConfigurableVolumeCapacity int                 `json:"estimatedConfigurableVolumeCapacity,omitempty"`
	SavingEffects                       *SavingEffectsAdmin `json:"savingEffects,omitempty"`
	GumVersion                          string              `json:"gumVersion,omitempty"`
	EsmOsVersion                        string              `json:"esmOsVersion,omitempty"`
	DkcMicroVersion                     string              `json:"dkcMicroVersion,omitempty"`
	WarningLedStatus                    string              `json:"warningLedStatus,omitempty"`
	EsmStatus                           string              `json:"esmStatus,omitempty"`
	IpAddressIpv4Service                string              `json:"ipAddressIpv4Service,omitempty"`
	IpAddressIpv4Ctl1                   string              `json:"ipAddressIpv4Ctl1,omitempty"`
	IpAddressIpv4Ctl2                   string              `json:"ipAddressIpv4Ctl2,omitempty"`
	IpAddressIpv6Service                string              `json:"ipAddressIpv6Service,omitempty"`
	IpAddressIpv6Ctl1                   string              `json:"ipAddressIpv6Ctl1,omitempty"`
	IpAddressIpv6Ctl2                   string              `json:"ipAddressIpv6Ctl2,omitempty"`
}

type StorageSettingsAndInfo struct {
	Settings StorageDeviceSettings `json:"settings"`
	Info     *StorageSystem        `json:"info"`
}

type StorageVersionInfo struct {
	ApiVersion  string `json:"apiVersion"`
	ProductName string `json:"productName"`
}

type AllStorageTypes struct {
	VspStorageSystem       []*StorageSystem
	VssbStorageVersionInfo []*StorageVersionInfo
	AdminStorageSystem     []*StorageSystemAdmin
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
