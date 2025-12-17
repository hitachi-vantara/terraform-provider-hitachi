package admin

type IscsiTargetInfoByPort struct {
	PortID          string  `json:"portId"`                    // Assigned port
	TargetIscsiName *string `json:"targetIscsiName,omitempty"` // (nullable) Target port iSCSI name (for iSCSI)
}

type IscsiTargetInfoList struct {
	Data  []IscsiTargetInfoByPort `json:"data"`  // List of target ports
	Count int                     `json:"count"` // Number of items stored in data
}

/////////////////////////////
// PARAMS

type RenameIscsiTargetNameParam struct {
	TargetIscsiName string `json:"targetIscsiName"`
}

// For Server
type HostGroupForAddToServerParam struct {
	PortID        string  `json:"portId"`
	HostGroupID   *int    `json:"hostGroupId,omitempty"`
	HostGroupName *string `json:"hostGroupName,omitempty"`
}

type AddHostGroupsToServerParam struct {
	HostGroups []HostGroupForAddToServerParam `json:"hostGroups"` // Host Group list
}
