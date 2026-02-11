package admin

// BaseVolumeInfo contains the common fields returned by both GetVolumes and GetVolumeByID
type BaseVolumeInfo struct {
	ID                          int      `json:"id"`
	Nickname                    *string  `json:"nickname,omitempty"`
	PoolID                      int      `json:"poolId"`
	PoolName                    *string  `json:"poolName,omitempty"`
	TotalCapacity               int64    `json:"totalCapacity"`
	UsedCapacity                int64    `json:"usedCapacity"`
	SavingSetting               string   `json:"savingSetting"`
	IsDataReductionShareEnabled *bool    `json:"isDataReductionShareEnabled,omitempty"`
	CompressionAcceleration     *bool    `json:"compressionAcceleration,omitempty"`
	CapacitySavingStatus        string   `json:"capacitySavingStatus"`
	NumberOfConnectingServers   int      `json:"numberOfConnectingServers"`
	NumberOfSnapshots           int      `json:"numberOfSnapshots"`
	VolumeTypes                 []string `json:"volumeTypes"`
}

// VolumeInfoByID extends BaseVolumeInfo with additional fields returned by GetVolumeByID
type VolumeInfoByID struct {
	BaseVolumeInfo
	FreeCapacity                  int64   `json:"freeCapacity"`
	ReservedCapacity              int64   `json:"reservedCapacity"`
	CompressionAccelerationStatus *string `json:"compressionAccelerationStatus,omitempty"`
	CapacitySavingProgress        *int    `json:"capacitySavingProgress,omitempty"`
	LUNs                          []LUN   `json:"luns,omitempty"`
}

type LUN struct {
	LUN      int    `json:"lun"`
	ServerID int    `json:"serverId"`
	PortID   string `json:"portId"`
}

// List for GetVolumes
type VolumeInfoList struct {
	Data       []BaseVolumeInfo `json:"data"`
	Count      int              `json:"count"`
	TotalCount int              `json:"totalCount"`
	HasNext    bool             `json:"hasNext"`
}

/////////////////////////////
// PARAMS

// get volumes
type GetVolumeParams struct {
	PoolID           *int    `json:"poolId,omitempty"`
	PoolName         *string `json:"poolName,omitempty"`
	ServerID         *int    `json:"serverId,omitempty"`
	ServerNickname   *string `json:"serverNickname,omitempty"`
	Nickname         *string `json:"nickname,omitempty"`
	MinTotalCapacity *int64  `json:"minTotalCapacity,omitempty"`
	MaxTotalCapacity *int64  `json:"maxTotalCapacity,omitempty"`
	MinUsedCapacity  *int64  `json:"minUsedCapacity,omitempty"`
	MaxUsedCapacity  *int64  `json:"maxUsedCapacity,omitempty"`
	StartVolumeID    *int    `json:"startVolumeId,omitempty"`
	Count            *int    `json:"count,omitempty"`
}

// create volume
type CreateVolumeParams struct {
	Capacity                    int64               `json:"capacity"`
	Number                      *int                `json:"number,omitempty"` // default: 1
	NicknameParam               VolumeNicknameParam `json:"nicknameParam,omitempty"`
	SavingSetting               *string             `json:"savingSetting,omitempty"` // DEDUPLICATION_AND_COMPRESSION | COMPRESSION
	IsDataReductionShareEnabled *bool               `json:"isDataReductionShareEnabled,omitempty"`
	PoolID                      int                 `json:"poolId"`
	CompressionAcceleration     *bool               `json:"compressionAcceleration,omitempty"` // only for update operation
}

type VolumeNicknameParam struct {
	BaseName       string `json:"baseName"`
	StartNumber    *int   `json:"startNumber,omitempty"`
	NumberOfDigits *int   `json:"numberOfDigits,omitempty"`
}

type ExpandVolumeParams struct {
	Capacity int64 `json:"capacity"` // Incremental capacity (MiB)
}

type UpdateVolumeNicknameParams struct {
	Nickname string `json:"nickname"`
}

type UpdateVolumeReductionParams struct {
	SavingSetting           *string `json:"savingSetting,omitempty"` // e.g., "DEDUPLICATION_AND_COMPRESSION" or "COMPRESSION"
	CompressionAcceleration *bool   `json:"compressionAcceleration,omitempty"`
}
