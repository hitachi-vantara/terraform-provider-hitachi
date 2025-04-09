package terraform

// Node .
type Node struct {
	ID              string `json:"id,omitempty"`
	Nickname        string `json:"nickName,omitempty"`
	TotalCapacity   int    `json:"totalCapacity,omitempty"`
	UsedCapacity    int    `json:"usedCapacity,omitempty"`
	OsType          string `json:"osType,omitempty"`
	NumberOfVolumes int    `json:"numberOfVolumes,omitempty"`
}

// Nodes .
type Nodes struct {
	Data []Node `json:"data"`
}

type VolumeNode struct {
	ID              string `json:"id,omitempty"`
	Nickname        string `json:"nickName,omitempty"`
	TotalCapacity   int    `json:"totalCapacity,omitempty"`
	UsedCapacity    int    `json:"usedCapacity,omitempty"`
	OsType          string `json:"osType,omitempty"`
	NumberOfVolumes int    `json:"numberOfVolumes,omitempty"`
	Lun             int    `json:"lun,omitempty"`
}

type VolumeNodes struct {
	Data []VolumeNode `json:"data"`
}