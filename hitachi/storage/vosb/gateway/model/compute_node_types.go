package vssbstorage

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
	Lun             int    `json:"lun,omitempty"`
	Paths           []Path `json:"paths,omitempty"`
}

type Path struct {
	HbaName  string   `json:"hbaName"`
	PortIds  []string `json:"portIds"`
	Protocol string   `json:"protocol"`
}

type Servers struct {
	Data []Server `json:"data"`
}

// ComputeNodeInformation used to edit information of compute node
type ComputeNodeInformation struct {
	Nickname string `json:"nickname"`
	OsType   string `json:"osType"`
}

// ComputeNodeCreateReq used to register information of compute node
type ComputeNodeCreateReq struct {
	ServerNickname string `json:"serverNickname"`
	OsType         string `json:"osType"`
}

// RegisterInitiator used to register protocol and Iscsi name
type RegisterInitiator struct {
	Protocol  string `json:"protocol"`
	IscsiName string `json:"iscsiName"`
}

// RegisterHba used to register protocol and HbaWwn
type RegisterHba struct {
	Protocol    string `json:"protocol"`
	HbaWwn      string `json:"hbaWwn"`
	IsTargetAny bool   `json:"isTargetAny"`
}

// Initiator is used for single initiator
type Initiator struct {
	ID        string   `json:"id"`
	ServerID  string   `json:"serverId"`
	Name      string   `json:"name"`
	IscsiName string   `json:"iscsiName"`
	Protocol  string   `json:"protocol"`
	PortIDs   []string `json:"portIds"`
}

// Initiators used to get all initiators
type Initiators struct {
	Data []Initiator `json:"data"`
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

// ComputeNodePathReq request body for getAll/delete/add apis
type ComputeNodePathReq struct {
	HbaId  string `json:"hbaId"`
	PortId string `json:"portId"`
}

// VolumeServerConnectionInfo
type VolumeServerConnectionInfo struct {
	ID       string `json:"id"`
	ServerId string `json:"serverId"`
	VolumeId string `json:"volumeId"`
	Lun      int    `json:"lun"`
}

// VolumeServerConnectionsInfo
type VolumeServerConnectionsInfo struct {
	Data []VolumeServerConnectionInfo `json:"data"`
}

// SetPathVolumeServerReq
type SetPathVolumeServerReq struct {
	ServerId string `json:"serverId"`
	VolumeId string `json:"volumeId"`
}

// ReleaseMultiConVolumeServerReq
type ReleaseMultiConVolumeServerReq struct {
	ServerIds []string `json:"serverIds"`
	VolumeIds []string `json:"volumeIds"`
}
