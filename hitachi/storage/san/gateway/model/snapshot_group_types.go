package sanstorage

// GetSnapshotGroupsParams defines the query parameters for filtering snapshot groups.
type GetSnapshotGroupsParams struct {
	// snapshotGroupName: (Optional) Get info about pairs in a specific group.
	SnapshotGroupName *string `json:"snapshotGroupName,omitempty"`
	// detailInfoType: (Optional) Set to "pair" to get nested pair details for all groups.
	DetailInfoType *string `json:"detailInfoType,omitempty"`
}

// SnapshotGroup represents a group container or a group entry in a list.
type SnapshotGroup struct {
	SnapshotGroupName string     `json:"snapshotGroupName"`
	SnapshotGroupID   string     `json:"snapshotGroupId,omitempty"`
	Snapshots         []Snapshot `json:"snapshots,omitempty"` // Populated if detailInfoType=pair
}

// SnapshotGroupListResponse is the top-level container for snapshot group data.
type SnapshotGroupListResponse struct {
	Data []SnapshotGroup `json:"data"`
}

type DeleteSnapshotTreeRequest struct {
	Parameters DeleteSnapshotTreeParams `json:"parameters"`
}

type DeleteSnapshotTreeParams struct {
	// LdevID is the LDEV number of the root volume of the snapshot tree.
	LdevID int `json:"ldevId"`
}

type DeleteGarbageDataRequest struct {
	Parameters DeleteGarbageDataParams `json:"parameters"`
}

type DeleteGarbageDataParams struct {
	LdevID        int    `json:"ldevId"`
	OperationType string `json:"operationType"` // start or stop
}
