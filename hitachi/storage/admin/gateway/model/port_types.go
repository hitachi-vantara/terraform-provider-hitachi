package admin

// ------------------- Port Information Types -------------------

// IPv4Information represents IPv4 network configuration
type IPv4Information struct {
	Address        string `json:"address"`
	SubnetMask     string `json:"subnetMask"`
	DefaultGateway string `json:"defaultGateway"`
}

// IPv6Information represents IPv6 network configuration
type IPv6Information struct {
	Linklocal              string `json:"linklocal"`
	LinklocalAddress       string `json:"linklocalAddress"`
	LinklocalAddressStatus string `json:"linklocalAddressStatus"`
	Global                 string `json:"global"`
	GlobalAddress          string `json:"globalAddress"`
	GlobalAddressStatus    string `json:"globalAddressStatus"`
	DefaultGateway         string `json:"defaultGateway"`
}

// FCInformation represents Fibre Channel port information
type FCInformation struct {
	AlPa                string `json:"alPa"`
	FabricSwitchSetting bool   `json:"fabricSwitchSetting"`
	ConnectionType      string `json:"connectionType"`
	SfpDataTransferRate string `json:"sfpDataTransferRate"`
	PortMode            string `json:"portMode"`
}

// ISCSIInformation represents iSCSI port information
type ISCSIInformation struct {
	VlanUse             bool             `json:"vlanUse"`
	VlanId              int              `json:"vlanId"`
	IpMode              string           `json:"ipMode"`
	Ipv4Information     *IPv4Information `json:"ipv4Information,omitempty"`
	Ipv6Information     *IPv6Information `json:"ipv6Information,omitempty"`
	IsIpv6Updating      bool             `json:"isIpv6Updating"`
	SelectiveAck        bool             `json:"selectiveAck"`
	DelayedAck          bool             `json:"delayedAck"`
	MtuSize             string           `json:"mtuSize"`
	LinkMtuSize         string           `json:"linkMtuSize"`
	VirtualPortEnabled  bool             `json:"virtualPortEnabled"`
	TcpPort             int              `json:"tcpPort"`
	WindowSize          string           `json:"windowSize"`
	KeepAliveTimer      int              `json:"keepAliveTimer"`
	IsnsServerMode      bool             `json:"isnsServerMode"`
	IsnsServerIpAddress string           `json:"isnsServerIpAddress"`
	IsnsServerPort      int              `json:"isnsServerPort"`
}

// NVMeTCPInformation represents NVMe-TCP port information
type NVMeTCPInformation struct {
	VlanUse            bool             `json:"vlanUse"`
	VlanId             int              `json:"vlanId"`
	IpMode             string           `json:"ipMode"`
	Ipv4Information    *IPv4Information `json:"ipv4Information,omitempty"`
	Ipv6Information    *IPv6Information `json:"ipv6Information,omitempty"`
	IsIpv6Updating     bool             `json:"isIpv6Updating"`
	SelectiveAck       bool             `json:"selectiveAck"`
	DelayedAck         bool             `json:"delayedAck"`
	MtuSize            string           `json:"mtuSize"`
	LinkMtuSize        string           `json:"linkMtuSize"`
	VirtualPortEnabled bool             `json:"virtualPortEnabled"`
	TcpPort            int              `json:"tcpPort"`
	DiscoveryTcpPort   int              `json:"discoveryTcpPort"`
	WindowSize         string           `json:"windowSize"`
}

// PortInfo represents individual port information
type PortInfo struct {
	ID                 string              `json:"id"`
	Protocol           string              `json:"protocol"`
	PortWwn            string              `json:"portWwn"`
	PortIscsiName      string              `json:"portIscsiName"`
	PortSpeed          string              `json:"portSpeed"`
	ActualPortSpeed    string              `json:"actualPortSpeed"`
	PortSecurity       bool                `json:"portSecurity"`
	FcInformation      *FCInformation      `json:"fcInformation,omitempty"`
	IscsiInformation   *ISCSIInformation   `json:"iscsiInformation,omitempty"`
	NvmeTcpInformation *NVMeTCPInformation `json:"nvmeTcpInformation,omitempty"`
}

// PortInfoList represents the response for multiple ports (API response format)
type PortInfoList struct {
	Data  []PortInfo `json:"data"`
	Count int        `json:"count"`
}

// GetPortParams represents query parameters for getting ports
type GetPortParams struct {
	Protocol *string `json:"protocol,omitempty"`
}

// UpdatePortParams represents parameters for updating a port
type UpdatePortParams struct {
	PortSpeed          *string             `json:"portSpeed,omitempty"`
	PortSecurity       *bool               `json:"portSecurity,omitempty"`
	FcInformationParam      *FCInformationParam      `json:"fcInformation,omitempty"`
	IscsiInformationParam   *IscsiInformationParam   `json:"iscsiInformation,omitempty"`
	NvmeTcpInformationParam *NvmeTcpInformationParam `json:"nvmeTcpInformation,omitempty"`
}

type FCInformationParam struct {
    AlPa                *string `json:"alPa,omitempty"`                // hex 01–EF
    FabricSwitchSetting *bool   `json:"fabricSwitchSetting,omitempty"`
    ConnectionType      *string `json:"connectionType,omitempty"`
}

type IscsiInformationParam struct {
    VlanUse          *bool            `json:"vlanUse,omitempty"`
    AddVlanID        *int             `json:"addVlanId,omitempty"`
    DeleteVlanID     *int             `json:"deleteVlanId,omitempty"`
    IPMode           *string          `json:"ipMode,omitempty"`
    IPv4Info         *IPv4Information `json:"ipv4Information,omitempty"`
    IPv6Info         *IPv6Information `json:"ipv6Information,omitempty"`
    TCPPort          *int             `json:"tcpPort,omitempty"`
    SelectiveAck     *bool            `json:"selectiveAck,omitempty"`
    DelayedAck       *bool            `json:"delayedAck,omitempty"`
    WindowSize       *string          `json:"windowSize,omitempty"`
    MTUSize          *string          `json:"mtuSize,omitempty"`
    KeepAliveTimer   *int             `json:"keepAliveTimer,omitempty"`
    IsnsServerMode   *bool            `json:"isnsServerMode,omitempty"`
    IsnsServerIP     *string          `json:"isnsServerIpAddress,omitempty"`
    IsnsServerPort   *int             `json:"isnsServerPort,omitempty"`
}

type NvmeTcpInformationParam struct {
    VlanUse          *bool            `json:"vlanUse,omitempty"`
    AddVlanID        *int             `json:"addVlanId,omitempty"`
    DeleteVlanID     *int             `json:"deleteVlanId,omitempty"`
    IPMode           *string          `json:"ipMode,omitempty"`
    IPv4Info         *IPv4Information `json:"ipv4Information,omitempty"`
    IPv6Info         *IPv6Information `json:"ipv6Information,omitempty"`
    TCPPort          *int             `json:"tcpPort,omitempty"`
    DiscoveryTCPPort *int             `json:"discoveryTcpPort,omitempty"`
    SelectiveAck     *bool            `json:"selectiveAck,omitempty"`
    DelayedAck       *bool            `json:"delayedAck,omitempty"`
    WindowSize       *string          `json:"windowSize,omitempty"`
    MTUSize          *string          `json:"mtuSize,omitempty"`
}
