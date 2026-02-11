package sanstorage

type SnapshotFamily struct {
	LdevID                     int    `json:"ldevId"`
	SnapshotGroupName          string `json:"snapshotGroupName,omitempty"`
	PrimaryOrSecondary         string `json:"primaryOrSecondary,omitempty"`
	Status                     string `json:"status,omitempty"`
	PvolLdevID                 int    `json:"pvolLdevId,omitempty"`
	MuNumber                   int    `json:"muNumber,omitempty"`
	SvolLdevID                 int    `json:"svolLdevId,omitempty"`
	PoolID                     int    `json:"poolId,omitempty"`
	IsVirtualCloneVolume       bool   `json:"isVirtualCloneVolume"`
	IsVirtualCloneParentVolume bool   `json:"isVirtualCloneParentVolume"`
	SplitTime                  string `json:"splitTime,omitempty"`
	ParentLdevID               int    `json:"parentLdevId"`
	SnapshotGroupID            string `json:"snapshotGroupId,omitempty"`
	SnapshotID                 string `json:"snapshotId,omitempty"`
}

type SnapshotFamilyListResponse struct {
	Data []SnapshotFamily `json:"data"`
}

type VirtualCloneRequest struct {
	Parameters VirtualCloneParams `json:"parameters"`
}

type VirtualCloneParams struct {
	OperationType string `json:"operationType"` // "create", "convert", or "restore"
}

type VirtualCloneParentVolume struct {
	LdevID int `json:"ldevId"`
}

type VirtualCloneParentVolumeList struct {
	Data []VirtualCloneParentVolume `json:"data"`
}
