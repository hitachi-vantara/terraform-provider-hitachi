package sanstorage

import (
	gwymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

type SnapshotReconcilerInput struct {
	// StorageModel *string // e.g., "VSP G Series", "VSP F Series", "VSP E Series"
	Action *string // "read", "create", "split", "resync", "restore", "delete"

	// --- Required Fields ---
	SnapshotGroupName *string `json:"snapshotGroupName"` // Name of the snapshot group (1-32 chars)
	SnapshotPoolID    *int    `json:"snapshotPoolId"`    // ID of the Thin Image or HDP pool
	PvolLdevID        *int    `json:"pvolLdevId"`        // LDEV number of the Primary Volume

	// --- Optional Fields ---
	MuNumber           *int  `json:"muNumber,omitempty"`           // MU number of the P-VOL
	SvolLdevID         *int  `json:"svolLdevId,omitempty"`         // Required if IsClone is true
	IsConsistencyGroup *bool `json:"isConsistencyGroup,omitempty"` // Default: false
	AutoSplit          *bool `json:"autoSplit,omitempty"`          // Default: false
	CanCascade         *bool `json:"canCascade,omitempty"`         // Default: same as IsClone

	// --- Clone ---
	IsClone          *bool `json:"isClone,omitempty"`          // Default: false
	ClonesAutomation *bool `json:"clonesAutomation,omitempty"` // Default: false
	// CopySpeed: slower, medium, faster. Default: medium
	CopySpeed *string `json:"copySpeed,omitempty"`

	IsDataReductionForceCopy *bool   `json:"isDataReductionForceCopy,omitempty"` // Default: false
	RetentionPeriod          *int    `json:"retentionPeriod,omitempty"`          // Retention period in hours
	DefragOperation          *string `json:"defragOperation,omitempty"`
}

// added extra info struct
type SnapshotUniversalInfo struct {
	StorageSerial       int      `json:"storageSerial"`
	IsThinImageAdvanced bool     `json:"isThinImageAdvanced"`
	SnapshotPoolType    string   `json:"snapshotPoolType"` // "HDP" or "HTI"
	PvolAttributes      []string `json:"pvolAttributes"`   // e.g., "VCP", "CVS"
	SvolAttributes      []string `json:"svolAttributes"`   // e.g., "VC", "HDP"
}

type ReconcileSnapshotResult struct {
	Snapshot      *gwymodel.Snapshot
	VcloneFamily  *gwymodel.SnapshotFamily
	UniversalInfo *SnapshotUniversalInfo
}

type SnapshotGetMultipleInput struct {
	SnapshotGroupName *string `json:"snapshotGroupName"`    // Name of the snapshot group (1-32 chars)
	PvolLdevID        *int    `json:"pvolLdevId"`           // LDEV number of the Primary Volume
	SvolLdevID        *int    `json:"svolLdevId,omitempty"` // Required if IsClone is true
	MuNumber          *int    `json:"muNumber,omitempty"`   // MU number of the P-VOL
}

// This API can be used when the storage system is the VSP 5000 series.
type SnapshotGetMultipleRangeInput struct {
	StartPvolLdevID *int `json:"startPvolLdevId,omitempty"` // Default: 0 if omitted.
	EndPvolLdevID   *int `json:"endPvolLdevId,omitempty"`   // Default: Maximum LDEV number if omitted.
}
