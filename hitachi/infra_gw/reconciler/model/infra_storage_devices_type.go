package models

type CreateStorageDeviceParam struct {
	SerialNumber      string `json:"serialNumber,omitempty"`
	ManagementAddress string `json:"managementAddress,omitempty"`
	Username          string `json:"username,omitempty"`
	Password          string `json:"password,omitempty"`
	OutOfBand         bool   `json:"outOfBand,omitempty"`
	GatewayAddress    string `json:"gatewayAddress,omitempty"`
	UcpSystem         string `json:"ucpSystem,omitempty"`
}
