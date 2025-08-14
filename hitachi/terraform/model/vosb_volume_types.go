package terraform

// Volume .
type Volume struct {
	SavingEffect              SavingEffects `json:"savingEffects,omitempty"`
	ID                        string        `json:"id,omitempty"`
	Name                      string        `json:"name,omitempty"`
	NickName                  string        `json:"nickname,omitempty"`
	VolumeNumber              int           `json:"volumeNumber,omitempty"`
	PoolId                    string        `json:"poolId,omitempty"`
	PoolName                  string        `json:"poolName,omitempty"`
	TotalCapacity             int           `json:"totalCapacity,omitempty"`
	UsedCapacity              int           `json:"usedCapacity,omitempty"`
	NumberOfConnectingServers int           `json:"numberOfConnectingServers,omitempty"`
	NumberOfSnapshots         int           `json:"numberOfSnapshots,omitempty"`
	ProtectionDomainId        string        `json:"protectionDomainId,omitempty"`
	FullAllocated             bool          `json:"fullAllocated,omitempty"`
	VolumeType                string        `json:"volumeType,omitempty"`
	StatusSummary             string        `json:"statusSummary,omitempty"`
	Status                    string        `json:"status,omitempty"`
	StorageControllerId       string        `json:"storageControllerId,omitempty"`
	SnapshotAttribute         string        `json:"snapshotAttribute,omitempty"`
	SnapshotStatus            string        `json:"snapshotStatus,omitempty"`
	SavingSetting             string        `json:"savingSetting,omitempty"`
	SavingMode                string        `json:"savingMode,omitempty"`
	DataReductionStatus       string        `json:"dataReductionStatus,omitempty"`
	DataReductionProgressRate int           `json:"dataReductionProgressRate,omitempty"`
	ComputeNodes              []VolumeNode  `json:"computeNodes,omitempty"`
}

// SavingEffects .
type SavingEffects struct {
	DataCapacity int `json:"systemDataCapacity,omitempty"`
	PreCapacity  int `json:"preCapacityDataReductionWithoutSystemData,omitempty"`
	PostCapacity int `json:"postCapacityDataReduction,omitempty"`
}

// Volumes .
type Volumes struct {
	Data []Volume `json:"data"`
}

type CreateVolume struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	PoolName     string   `json:"poolName,omitempty"`
	CapacityInGB float32  `json:"capacityInGB,omitempty"`
	NickName     string   `json:"nickname,omitempty"`
	ComputeNodes []string `json:"computeNodes,omitempty"`
}
