package admin

type VolumeServerConnectionDetail struct {
	Id       string           `json:"id,omitempty"` // Composite ID: "{volumeId},{serverId}"
	VolumeId int              `json:"volumeId,omitempty"`
	ServerId int              `json:"serverId,omitempty"`
	Luns     []LunInformation `json:"luns,omitempty"`
}

type LunInformation struct {
	Lun    int    `json:"lun"`
	PortId string `json:"portId,omitempty"` // Port ID associated with the LUN
}

type VolumeServerConnectionsResponse struct {
	Data       []VolumeServerConnectionDetail `json:"data,omitempty"`
	Count      int                            `json:"count,omitempty"`      // Number of entries returned
	TotalCount int                            `json:"totalCount,omitempty"` // Total number of entries available
	HasNext    bool                           `json:"hasNext,omitempty"`    // Indicates if there are more results
}

/////////////////////////////
// PARAMS

type GetVolumeServerConnectionsParams struct {
	ServerId       *int    `json:"serverId,omitempty"`       // Server ID connected to the volume
	ServerNickname *string `json:"serverNickname,omitempty"` // Server nickname connected to the volume
	StartVolumeId  *int    `json:"startVolumeId,omitempty"`  // Starting volume ID for display
	Count          *int    `json:"count,omitempty"`          // Number of connection info entries to display
}

type AttachVolumeServerConnectionParam struct {
	VolumeIds []int `json:"volumeIds"` // List of volume IDs to attach
	ServerIds []int `json:"serverIds"` // List of server IDs to connect to
}
