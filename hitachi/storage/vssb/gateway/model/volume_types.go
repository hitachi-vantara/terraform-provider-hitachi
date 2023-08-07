package vssbstorage

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
}

type SavingEffects struct {
	DataCapacity int `json:"systemDataCapacity,omitempty"`
	PreCapacity  int `json:"preCapacityDataReductionWithoutSystemData,omitempty"`
	PostCapacity int `json:"postCapacityDataReduction,omitempty"`
}

type Volumes struct {
	Data []Volume `json:"data"`
}

// CreateLunRequestGwy .
type CreateVolumeRequestGwy struct {
	Capacity      *int32        `json:"capacity"`
	PoolID        *string       `json:"poolId"`
	NickNameParam NickNameParam `json:"nicknameParam,omitempty"`
	NameParam     NameParam     `json:"nameParam,omitempty"`
}

type NickNameParam struct {
	BaseName *string `json:"baseName"`
}

type NameParam struct {
	BaseName *string `json:"baseName"`
}

type AddVolumeToComputeNodeReq struct {
	VolumeID string `json:"volumeId"`
	ServerID string `json:"serverId"`
}

type UpdateVolumeNickNameReq struct {
	NickName string `json:"nickname"`
}

type UpdateVolumeSizeReq struct {
	AdditionalCapacity *int32 `json:"additionalCapacity"`
}
// {
// 	"capacity": 52,
// 	"number": 1,
// 	"nameParam": {
// 	  "baseName": "volume_test_1"
// 	},
// 	"nicknameParam": {
// 	  "baseName": "volume_nick_name"
// 	},
// 	"savingSetting": "Disabled",
// 	"poolId": "c89c1cbc-e181-4db3-ae44-eb7c6e496192",
// 	"fullAllocated": false
//   }
