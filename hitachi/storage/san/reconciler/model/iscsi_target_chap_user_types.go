package sanstorage

type IscsiTargetChapUser struct {
	ChapUserID      string `json:"chapUserId"`
	PortID          string `json:"portId"`
	HostGroupNumber int    `json:"hostGroupNumber"`
	HostGroupName   string `json:"hostGroupName"`
	ChapUserName    string `json:"chapUserName"`
	WayOfChapUser   string `json:"wayOfChapUser"`
}

type IscsiTargetChapUsers struct {
	IscsiTargetChapUsers []IscsiTargetChapUser `json:"data"`
}

// ChapUserRequest
type ChapUserRequest struct {
	ChapUserName      string `json:"chapUserName"`
	PortID            string `json:"portId"`
	IscsiTargetNumber int    `json:"iscsiTargetNumber"`
	WayOfChapUser     string `json:"wayOfChapUser"`
	ChapUserSecret    string `json:"chapUserstring"`
}
