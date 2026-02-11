package sanstorage

type SnapshotGroupReconcilerInput struct {
	Action *string // "read", "create", "split" and so on

	SnapshotGroupName *string `json:"snapshotGroupName"`   // Name of the snapshot group (1-32 chars)
	AutoSplit         *bool   `json:"autoSplit,omitempty"` // Default: false
	CopySpeed         *string `json:"copySpeed,omitempty"`
	RetentionPeriod   *int    `json:"retentionPeriod,omitempty"` // Retention period in hours
}
