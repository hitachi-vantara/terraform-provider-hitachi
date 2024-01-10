package infra_gw

type UcpSystemInfo struct {
	ResourceId           string                   `json:"resourceId"`
	Name                 string                   `json:"name"`
	ResourceState        string                   `json:"resourceState"`
	ComputeDevices       []ComputeDeviceInfo      `json:"computeDevices"`
	StorageDevices       []ShortStorageDeviceInfo `json:"storageDevices"`
	EthernetSwitches     []EthernetSwitchInfo     `json:"ethernetSwitches"`
	FibreChannelSwitches []FibreChannelSwitchInfo `json:"fibreChannelSwitches"`
	SerialNumber         string                   `json:"serialNumber"`
	GatewayAddress       string                   `json:"gatewayAddress"`
	Model                string                   `json:"model"`
	Vcenter              string                   `json:"vcenter"`
	Zone                 string                   `json:"zone"`
	VcenterResourceId    string                   `json:"vcenterResourceId"`
	Region               string                   `json:"region"`
	GeoInformation       GeoInformation           `json:"geoInformation"`
	WorkloadType         string                   `json:"workloadType"`
	ResultStatus         string                   `json:"resultStatus"`
	ResultMessage        string                   `json:"resultMessage"`
	PluginRegistered     bool                     `json:"pluginRegistered"`
	Linked               bool                     `json:"linked"`
}

type ComputeDeviceInfo struct {
	ResourceId         string `json:"resourceId"`
	BmcAddress         string `json:"bmcAddress"`
	BmcFirmwareVersion string `json:"bmcFirmwareVersion"`
	BiosVersion        string `json:"biosVersion"`
	ResourceState      string `json:"resourceState"`
	Model              string `json:"model"`
	Serial             string `json:"serial"`
	IsManagement       bool   `json:"isManagement"`
	HealthStatus       string `json:"healthStatus"`
	GatewayAddress     string `json:"gatewayAddress"`
}

type ShortStorageDeviceInfo struct {
	SerialNumber     string   `json:"serialNumber"`
	ResourceId       string   `json:"resourceId"`
	Address          string   `json:"address"`
	Model            string   `json:"model"`
	MicrocodeVersion string   `json:"microcodeVersion"`
	ResourceState    string   `json:"resourceState"`
	HealthState      string   `json:"healthState"`
	UcpSystems       []string `json:"ucpSystems"`
	SvpIp            string   `json:"svpIp"`
	GatewayAddress   string   `json:"gatewayAddress"`
}

type EthernetSwitchInfo struct {
	ResourceId      string `json:"resourceId"`
	Address         string `json:"address"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serialNumber"`
	FirmwareVersion string `json:"firmwareVersion"`
	ResourceState   string `json:"resourceState"`
	HealthStatus    string `json:"healthStatus"`
	GatewayAddress  string `json:"gatewayAddress"`
	IsManagement    bool   `json:"isManagement"`
}

type FibreChannelSwitchInfo struct {
	ResourceId      string `json:"resourceId"`
	Address         string `json:"address"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serialNumber"`
	FirmwareVersion string `json:"firmwareVersion"`
	ResourceState   string `json:"resourceState"`
	HealthState     string `json:"healthState"`
	SwitchName      string `json:"switchName"`
	GatewayAddress  string `json:"gatewayAddress"`
}

type GeoInformation struct {
	GeoLocation string `json:"geoLocation"`
	Country     string `json:"country"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Zipcode     string `json:"zipcode"`
}

type UcpSystem struct {
	Path    string        `json:"path"`
	Message string        `json:"message"`
	Data    UcpSystemInfo `json:"data"`
}

type UcpSystems struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    []UcpSystemInfo `json:"data"`
}

type CreateUcpSystemParam struct {
	// Required
	Name           string `json:"name"`
	SerialNumber   string `json:"serialNumber"`
	GatewayAddress string `json:"gatewayAddress"`
	Model          string `json:"model"`
	Region         string `json:"region"`
	Country        string `json:"country"`
	Zipcode        string `json:"zipcode"`

	// Optional
	ComputeDevices       []string `json:"computeDevices"`
	StorageDevices       []string `json:"storageDevices"`
	EthernetSwitches     []string `json:"ethernetSwitches"`
	FibreChannelSwitches []string `json:"fibreChannelSwitches"`
	Vcenter              string   `json:"vcenter"`
	Zone                 string   `json:"zone"`
	WorkloadType         string   `json:"workloadType"`
	ValidSerialNumber    bool     `json:"validSerialNumber"`
}
