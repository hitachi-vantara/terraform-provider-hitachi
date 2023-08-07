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
	TargetChapUserName    string `json:"targetChapUserName"`
	TargetChapSecret      string `json:"targetChapSecret"`
	InitiatorChapUserName string `json:"initiatorChapUserName,omitempty"`
	InitiatorChapSecret   string `json:"initiatorChapSecret,omitempty"`
}

type PortAuthSettings struct {
	ID                  string `json:"id,omitempty"`
	AuthMode            string `json:"authMode,omitempty"` // values : "CHAP" , "CHAPComplyingWithInitiatorSetting" , "None"
	IsDiscoveryChapAuth bool   `json:"isDiscoveryChapAuth,omitempty"`
	IsMutualChapAuth    bool   `json:"isMutualChapAuth,omitempty"`
}

// ChapUserIdReq
type ChapUserIdReq struct {
	ChapUserId string `json:"chapUserId"`
}
