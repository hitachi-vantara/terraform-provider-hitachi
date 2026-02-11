package sanstorage

// Snapshot represents a single Thin Image pair returned from the API.
type Snapshot struct {
	PvolLdevID                 int     `json:"pvolLdevId"`           // LDEV number of P-VOL
	MuNumber                   int     `json:"muNumber"`             // MU number of the P-VOL
	SvolLdevID                 int     `json:"svolLdevId,omitempty"` // LDEV number of S-VOL
	SnapshotPoolID             int     `json:"snapshotPoolId"`       // ID of the pool where snapshot data is created
	SnapshotID                 string  `json:"snapshotId"`           // Format: "pvolLdevId,muNumber"
	SnapshotGroupName          string  `json:"snapshotGroupName"`
	PrimaryOrSecondary         string  `json:"primaryOrSecondary"`             // Attribute of the LDEV (P-VOL or S-VOL)
	Status                     string  `json:"status"`                         // Pair status
	IsRedirectOnWrite          bool    `json:"isRedirectOnWrite"`              // True if the pair is Thin Image Advanced.
	IsConsistencyGroup         bool    `json:"isConsistencyGroup"`             // True if created in CTG mode.
	IsWrittenInSvol            bool    `json:"isWrittenInSvol"`                // True if data was written to the S-VOL when status was PSUS/PFUS.
	IsClone                    bool    `json:"isClone"`                        // True if the pair has the clone attribute.
	CanCascade                 bool    `json:"canCascade"`                     // True if the pair can be a cascaded pair.
	SnapshotDataReadOnly       bool    `json:"snapshotDataReadOnly"`           // True if the snapshot data has the read-only attribute.
	ConcordanceRate            *int    `json:"concordanceRate,omitempty"`      // Concordance rate, optional.
	ProgressRate               *int    `json:"progressRate,omitempty"`         // Progress of the processing, optional.
	SplitTime                  *string `json:"splitTime,omitempty"`            // Time when snapshot data was created.
	PvolProcessingStatus       string  `json:"pvolProcessingStatus,omitempty"` // Processing status of the P-VOL pair (E/N).
	SvolProcessingStatus       string  `json:"svolProcessingStatus,omitempty"` // Processing status of the S-VOL pair (E/N).
	RetentionPeriod            int     `json:"retentionPeriod"`                // in unit of hours if data retention is set in TIA
	IsVirtualCloneVolume       bool    `json:"isVirtualCloneVolume"`           // for VSP One B20.
	IsVirtualCloneParentVolume bool    `json:"isVirtualCloneParentVolume"`     // for VSP One B20.
}

// SnapshotListResponse represents the overall API response structure for a list of Thin Image pairs.
type SnapshotListResponse struct {
	Data []Snapshot `json:"data"`
}

// SnapshotAll represents a single Thin Image pair returned when requesting all pairs.
type SnapshotAll struct {
	SnapshotReplicationID string  `json:"snapshotReplicationId"`        // Format: "pvolLdevId,muNumber"
	SnapshotGroupName     string  `json:"snapshotGroupName,omitempty"`  // Name of the snapshot group (omitted if not in a group)
	PvolLdevID            int     `json:"pvolLdevId"`                   // LDEV number of P-VOL
	MuNumber              int     `json:"muNumber"`                     // MU number of the P-VOL (Used in the ID)
	SnapshotPoolID        int     `json:"snapshotPoolId"`               // ID of the pool where snapshot data was created
	SvolLdevID            int     `json:"svolLdevId,omitempty"`         // LDEV number of S-VOL (omitted if no S-VOL exists)
	ConsistencyGroupID    int     `json:"consistencyGroupId,omitempty"` // Consistency group ID (omitted if no consistency group exists)
	Status                string  `json:"status"`                       // Pair status (e.g., SMPP, COPY, PAIR, PSUS, PFUL, PSUE, PFUS, RCPY, PSUP, CPYP, OTHER)
	ConcordanceRate       *int    `json:"concordanceRate,omitempty"`    // Concordance rate for pairs (omitted if isRedirectOnWrite is true)
	IsRedirectOnWrite     bool    `json:"isRedirectOnWrite"`            // True if the pair is Thin Image Advanced.
	IsClone               bool    `json:"isClone"`                      // True if the pair has the clone attribute.
	CanCascade            bool    `json:"canCascade"`                   // True if the pair can be a cascaded pair.
	SplitTime             *string `json:"splitTime,omitempty"`
}

// SnapshotAllListResponse represents the overall API response structure for a list of all Thin Image pairs.
type SnapshotAllListResponse struct {
	Data []SnapshotAll `json:"data"`
}

//// Params

// GetSnapshotsParams holds the query parameters used to filter the list of Thin Image pairs.
type GetSnapshotsParams struct {
	SnapshotGroupName *string `url:"snapshotGroupName,omitempty"`
	PvolLdevID        *int    `url:"pvolLdevId,omitempty"`
	SvolLdevID        *int    `url:"svolLdevId,omitempty"`
	MuNumber          *int    `url:"muNumber,omitempty"`
	DetailInfoType    *string `url:"detailInfoType,omitempty"`
}

// SnapshotID represents the two required identifiers for a specific Thin Image pair.
type SnapshotID struct {
	PvolLdevID int `json:"pvolLdevId"`
	MuNumber   int `json:"muNumber"`
}

// GetSnapshotReplicationsRangeParams contains the optional query parameters for specifying
type GetSnapshotReplicationsRangeParams struct {
	StartPvolLdevID *int `url:"startPvolLdevId,omitempty"` // Default: 0 if omitted.
	EndPvolLdevID   *int `url:"endPvolLdevId,omitempty"`   // Default: Maximum LDEV number if omitted.
}

// CreateSnapshotParams defines the attributes for creating a single Thin Image pair.
type CreateSnapshotParams struct {
	// --- Required Fields ---
	SnapshotGroupName string `json:"snapshotGroupName"` // Name of the snapshot group (1-32 chars)
	SnapshotPoolID    int    `json:"snapshotPoolId"`    // ID of the Thin Image or HDP pool
	PvolLdevID        int    `json:"pvolLdevId"`        // LDEV number of the Primary Volume

	// --- Optional Fields ---
	MuNumber           *int  `json:"muNumber,omitempty"`           // MU number of the P-VOL
	SvolLdevID         *int  `json:"svolLdevId,omitempty"`         // Required if IsClone is true
	IsConsistencyGroup *bool `json:"isConsistencyGroup,omitempty"` // Default: false
	AutoSplit          *bool `json:"autoSplit,omitempty"`          // Default: false
	CanCascade         *bool `json:"canCascade,omitempty"`         // Default: same as IsClone
	IsClone            *bool `json:"isClone,omitempty"`            // Default: false
	ClonesAutomation   *bool `json:"clonesAutomation,omitempty"`   // Default: false

	// CopySpeed: slower, medium, faster. Default: medium
	CopySpeed                *string `json:"copySpeed,omitempty"`
	IsDataReductionForceCopy *bool   `json:"isDataReductionForceCopy,omitempty"` // Default: false

	RetentionPeriod *int `json:"retentionPeriod"` // in unit of hours if data retention is set in TIA
}

// SplitSnapshotRequest defines the body for the split action.
type SplitSnapshotRequest struct {
	Parameters SplitSnapshotParams `json:"parameters"`
}

// SplitSnapshotParams defines the optional retention period for TIA snapshots.
type SplitSnapshotParams struct {
	// for TIA
	RetentionPeriod *int `json:"retentionPeriod,omitempty"`
}

// ResyncSnapshotRequest defines the body for the resync action.
type ResyncSnapshotRequest struct {
	Parameters ResyncSnapshotParams `json:"parameters"`
}

// ResyncSnapshotParams defines the specific toggle for auto-splitting after resync.
type ResyncSnapshotParams struct {
	// AutoSplit: true to split the pair and store snapshot data after resync. Default: false.
	AutoSplit *bool `json:"autoSplit,omitempty"`
	// for TIA
	RetentionPeriod *int `json:"retentionPeriod,omitempty"`
}

// RestoreSnapshotRequest defines the body for the restore action.
type RestoreSnapshotRequest struct {
	Parameters RestoreSnapshotParams `json:"parameters"`
}

// RestoreSnapshotParams defines the toggle for auto-splitting after restore.
type RestoreSnapshotParams struct {
	// AutoSplit: If true, the pair is split and snapshot data is stored after restore.
	// Ignored for Thin Image Advanced pairs. Default: false.
	AutoSplit *bool `json:"autoSplit,omitempty"`
}

// CloneSnapshotRequest defines the body for the clone action.
type CloneSnapshotRequest struct {
	Parameters CloneSnapshotParams `json:"parameters"`
}

// CloneSnapshotParams defines the optional settings for the cloning process.
type CloneSnapshotParams struct {
	// CopySpeed: slower, medium, faster. Default: medium.
	// This item is not case sensitive.
	CopySpeed *string `json:"copySpeed,omitempty"`
}

// AssignSnapshotVolumeRequest defines the body for the assign-volume action.
type AssignSnapshotVolumeRequest struct {
	Parameters AssignSnapshotVolumeParams `json:"parameters"`
}

// AssignSnapshotVolumeParams defines the parameters for the secondary volume assignment.
type AssignSnapshotVolumeParams struct {
	// SvolLdevID is the required LDEV number of the S-VOL to be assigned.
	// This must be a virtual volume for Thin Image created beforehand.
	SvolLdevID int `json:"svolLdevId"`
}

// SetSnapshotRetentionPeriodRequest defines the body for the split action.
type SetSnapshotRetentionPeriodRequest struct {
	Parameters SetSnapshotRetentionPeriodParams `json:"parameters"`
}

// SetSnapshotRetentionPeriodParams defines the optional retention period for TIA snapshots.
type SetSnapshotRetentionPeriodParams struct {
	// for TIA
	RetentionPeriod *int `json:"retentionPeriod,omitempty"`
}
