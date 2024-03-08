package infra_gw

type StorageDeviceInfo struct {
	ResourceId        string `json:"resourceId"`
	SerialNumber      string `json:"serialNumber"`
	ManagementAddress string `json:"managementAddress"`
	ControllerAddress string `json:"controllerAddress"`
	Username          string `json:"username"`
	DeviceType        string `json:"deviceType"`
	Model             string `json:"model"`
	MicrocodeVersion  string `json:"microcodeVersion"`
	TotalCapacityInMb int    `json:"totalCapacityInMb"`
	FreeCapacityInMb  int    `json:"freeCapacityInMb"`
	TotalCapacity     string `json:"totalCapacity"`
	FreeCapacity      string `json:"freeCapacity"`
	OperationalStatus string `json:"operationalStatus"`

	ResourceState string   `json:"resourceState"`
	UcpSystems    []string `json:"ucpSystems"`
	Tags          []string `json:"tags"`
	HealthStatus  string   `json:"healthStatus"`

	FreeGadConsistencyGroupId         int `json:"freeGadConsistencyGroupId"`
	FreeLocalCloneConsistencyGroupId  int `json:"freeLocalCloneConsistencyGroupId"`
	FreeRemoteCloneConsistencyGroupId int `json:"freeRemoteCloneConsistencyGroupId"`

	StorageEfficiencyStat struct {
		AccelCompEfficiencyStat struct {
			CompressionRatio string `json:"compressionRatio"`
			DedupeRatio      string `json:"dedupeRatio"`
			ReclaimRatio     string `json:"reclaimRatio"`
			TotalRatio       string `json:"totalRatio"`
		} `json:"accelCompEfficiencyStat"`
		DedupeCompEfficiencyStat struct {
			CompressionRatio string `json:"compressionRatio"`
			DedupeRatio      string `json:"dedupeRatio"`
			ReclaimRatio     string `json:"reclaimRatio"`
			TotalRatio       string `json:"totalRatio"`
		} `json:"dedupeCompEfficiencyStat"`
		CompressionRatio string `json:"compressionRatio"`
		StartTime        string `json:"startTime"`
		EndTime          string `json:"endTime"`
		ProvisioningRate string `json:"provisioningRate"`
		SnapshotRatio    string `json:"snapshotRatio"`
		TotalRatio       string `json:"totalRatio"`
	} `json:"storageEfficiencyStat"`
	SyslogConfig struct {
		SyslogServers []struct {
			Id                  int    `json:"id"`
			SyslogServerAddress string `json:"syslogServerAddress"`
			SyslogServerPort    string `json:"syslogServerPort"`
		} `json:"syslogServers"`
		Detailed bool `json:"detailed"`
	} `json:"syslogConfig"`
	StorageDeviceLicenses []struct {
		IsEnabled   bool   `json:"isEnabled"`
		IsInstalled bool   `json:"isInstalled"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"storageDeviceLicenses"`
	DeviceLimits struct {
		ExternalGroupNumberRange struct {
			IsValid  bool `json:"isValid"`
			MaxValue int  `json:"maxValue"`
			MinValue int  `json:"minValue"`
		} `json:"externalGroupNumberRange"`
		ExternalGroupSubNumberRange struct {
			IsValid  bool `json:"isValid"`
			MaxValue int  `json:"maxValue"`
			MinValue int  `json:"minValue"`
		} `json:"dedupeCompEfficiencyStat"`
		ParityGroupNumberRange struct {
			IsValid  bool `json:"isValid"`
			MaxValue int  `json:"maxValue"`
			MinValue int  `json:"minValue"`
		} `json:"ParityGroupNumberRange"`
		ParityGroupSubNumberRange struct {
			IsValid  bool `json:"isValid"`
			MaxValue int  `json:"maxValue"`
			MinValue int  `json:"minValue"`
		} `json:"parityGroupSubNumberRange"`
		HealthDescription string `json:"healthDescription"`
		IsUnified         bool   `json:"isUnified"`
	} `json:"deviceLimits"`
}

type StorageDevices struct {
	Path    string              `json:"path"`
	Message string              `json:"message"`
	Data    []StorageDeviceInfo `json:"data"`
}

type StorageDevice struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    StorageDeviceInfo `json:"data"`
}

type CreateStorageDeviceParam struct {
	SerialNumber      string `json:"serialNumber,omitempty"`
	ManagementAddress string `json:"managementAddress,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	OutOfBand         bool   `json:"outOfBand,omitempty"`
	GatewayAddress    string `json:"gatewayAddress,omitempty"`
	UcpSystem         string `json:"ucpSystem,omitempty"`
}

type PatchStorageDeviceParam struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	OutOfBand bool   `json:"outOfBand,omitempty"`
}

type StorageDeviceToPartnerReq struct {
	PartnerId  string `json:"partnerId"`
	ResourceId string `json:"resourceId"`
}

const (
	DefaultSystemName         = "terraform-ucp-system"
	DefaultSystemSerialNumber = "Logical-UCP-95054"
)

type MTStorageDevice struct {
	Storage struct {
		ResourceId   string   `json:"resourceId"`
		SerialNumber string   `json:"serialNumber"`
		UcpSystems   []string `json:"ucpSystems"`
	} `json:"storage"`
	Status       string `json:"status"`
	PartnerId    string `json:"partnerId"`
	SubscriberId string `json:"subscriberId"`
	StorageId    string `json:"storageId"`
}

type MTStorageDevices []MTStorageDevice

type CreateMTStorageDeviceParam struct {
	ResourceId string `json:"resourceId"`
	PartnerId  string `json:"partnerId"`
}
