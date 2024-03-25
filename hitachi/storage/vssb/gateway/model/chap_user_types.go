package vssbstorage

type ChapUser struct {
	ID                    string `json:"id,omitempty"`
	TargetChapUserName    string `json:"targetChapUserName,omitempty"`
	InitiatorChapUserName string `json:"initiatorChapUserName,omitempty"`
}

type ChapUsers struct {
	Data []ChapUser `json:"data"`
}

type ChapUserReq struct {
	ID                    string `json:"id,omitempty"`
	TargetChapUserName    string `json:"targetChapUserName,omitempty"`
	TargetChapSecret      string `json:"targetChapSecret,omitempty"`
	InitiatorChapUserName string `json:"initiatorChapUserName,omitempty"`
	InitiatorChapSecret   string `json:"initiatorChapSecret,omitempty"`
}

type PortAuthSettings struct {
	AuthMode            string `json:"authMode"` // values : "CHAP" , "CHAPComplyingWithInitiatorSetting" , "None"
	IsDiscoveryChapAuth bool   `json:"isDiscoveryChapAuth"`
	IsMutualChapAuth    bool   `json:"isMutualChapAuth"`
}

// ChapUserIdReq
type ChapUserIdReq struct {
	ChapUserId string `json:"chapUserId"`
}
