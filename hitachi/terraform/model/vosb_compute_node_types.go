package terraform

// ComputeNode .
type Server struct {
	ID              string `json:"id,omitempty"`
	Nickname        string `json:"nickname,omitempty"`
	OsType          string `json:"osType,omitempty"`
	TotalCapacity   int    `json:"totalCapacity,omitempty"`
	UsedCapacity    int    `json:"usedCapacity,omitempty"`
	NumberOfPaths   int    `json:"numberOfPaths,omitempty"`
	VpsId           string `json:"vpsId,omitempty"`
	VpsName         string `json:"vpsName,omitempty"`
	NumberOfVolumes int    `json:"numberOfVolumes,omitempty"`
	Paths           []Path `json:"paths,omitempty"`
}

type ComputeNodeWithPathDetails struct {
	Node         Server           `json:"paths,omitempty"`
	ComputePaths ComputeNodePaths `json:"computepaths,omitempty"`
}

// Path .
type Path struct {
	HbaName string   `json:"hbaName,omitempty"`
	PortIds []string `json:"portIds,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// ComputeNodes .
type ComputeNodes struct {
	Data []Server `json:"data"`
}

// ComputeNodePath information
type ComputeNodePath struct {
	ID                string `json:"id,omitempty"`
	ServerId          string `json:"serverId,omitempty"`
	IScsiInitiatorIqn string `json:"hbaName,omitempty"`
	IScsiInitiatorId  string `json:"hbaId,omitempty"`
	PortId            string `json:"portId,omitempty"`
	PortName          string `json:"portName,omitempty"`
	PortNickname      string `json:"portNickname,omitempty"`
}

// ComputeNodePaths
type ComputeNodePaths struct {
	Data []ComputeNodePath `json:"data"`
}

// ComputeResource .
type ComputeResource struct {
	ID               string           `json:"id,omitempty"`
	Name             string           `json:"name"`
	OsType           string           `json:"ostype"`
	IscsiConnections []IscsiConnector `json:"iscsiconnection,omitempty"`
	FcConnections    []FcConnector    `json:"fcconnections,omitempty"`
}

// FcConnector
type FcConnector struct {
	HostWWN string `json:"hostwwn,omitempty"`
}

// IscsiConnector .
type IscsiConnector struct {
	IscsiInitiator string   `json:"iscsiinitiator,omitempty"`
	PortNames      []string `json:"portnames,omitempty"`
}

// ComputeResourceOutput .
type ComputeResourceOutput struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	OsType          string           `json:"ostype"`
	IscsiConnection []IscsiConnector `json:"iscsiconnection,omitempty"`
}
