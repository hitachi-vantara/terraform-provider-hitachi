package sanstorage

type ISCSI struct {
	HostGroupID          string `json:"hostGroupId"`
	PortID               string `json:"portId"`
	HostGroupNumber      int    `json:"hostGroupNumber"`
	HostGroupName        string `json:"hostGroupName"`
	IscsiName            string `json:"iscsiName"`
	AuthenticationMode   string `json:"authenticationMode"`
	IscsiTargetDirection string `json:"iscsiTargetDirection"`
	HostMode             string `json:"hostMode"`
	HostModeOptions      []int  `json:"hostModeOptions"`
}