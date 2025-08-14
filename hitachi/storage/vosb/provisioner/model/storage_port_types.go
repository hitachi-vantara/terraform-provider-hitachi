package vssbstorage

type StoragePort struct {
	ID                  string `json:"id"`
	Protocol            string `json:"protocol"`
	Type                string `json:"type"`
	Nickname            string `json:"nickname"`
	Name                string `json:"name"`
	ConfiguredPortSpeed string `json:"configuredPortSpeed"`
	PortSpeed           string `json:"portSpeed"`
	PortSpeedDuplex     string `json:"portSpeedDuplex"`
	ProtectionDomainID  string `json:"protectionDomainId"`
	StorageNodeID       string `json:"storageNodeId"`
	InterfaceName       string `json:"interfaceName"`
	StatusSummary       string `json:"statusSummary"`
	Status              string `json:"status"`
	FcInformation       struct {
		ConnectionType      string `json:"connectionType"`
		SfpDataTransferRate string `json:"sfpDataTransferRate"`
		PhysicalWwn         string `json:"physicalWwn"`
	} `json:"fcInformation"`
	IscsiInformation struct {
		IPMode          string `json:"ipMode"`
		Ipv4Information struct {
			Address        string `json:"address"`
			SubnetMask     string `json:"subnetMask"`
			DefaultGateway string `json:"defaultGateway"`
		} `json:"ipv4Information"`
		Ipv6Information struct {
			LinklocalAddressMode string `json:"linklocalAddressMode"`
			LinklocalAddress     string `json:"linklocalAddress"`
			GlobalAddressMode    string `json:"globalAddressMode"`
			GlobalAddress1       string `json:"globalAddress1"`
			SubnetPrefixLength1  int    `json:"subnetPrefixLength1"`
			DefaultGateway       string `json:"defaultGateway"`
		} `json:"ipv6Information"`
		DelayedAck          bool   `json:"delayedAck"`
		MtuSize             int    `json:"mtuSize"`
		MacAddress          string `json:"macAddress"`
		IsIsnsClientEnabled bool   `json:"isIsnsClientEnabled"`
		IsnsServers         []struct {
			Index      int    `json:"index"`
			ServerName string `json:"serverName"`
			Port       int    `json:"port"`
		} `json:"isnsServers"`
	} `json:"iscsiInformation"`
}

type StoragePorts struct {
	Data []StoragePort `json:"data"`
}
