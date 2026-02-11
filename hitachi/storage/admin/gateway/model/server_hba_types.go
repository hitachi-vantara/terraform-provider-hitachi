package admin

// ServerHBA represents individual server HBA information
type ServerHBA struct {
	ServerID  int      `json:"serverId"`
	HbaWwn    string   `json:"hbaWwn"`
	IscsiName string   `json:"iscsiName"`
	PortIds   []string `json:"portIds"`
}

// ServerHBAList represents the response for multiple server HBAs (API response format)
type ServerHBAList struct {
	Data  []ServerHBA `json:"data"`
	Count int         `json:"count"`
}

// GetServerHBAParams represents query parameters for getting server HBAs
type GetServerHBAParams struct {
	ServerID int `json:"serverId"`
}

// CreateServerHBAParams represents parameters for creating/adding server HBAs
type CreateServerHBAParams struct {
	HBAs []ServerHBARequest `json:"hbas"`
}

// ServerHBARequest represents individual HBA data for create requests
type ServerHBARequest struct {
	HbaWwn    string `json:"hbaWwn,omitempty"`
	IscsiName string `json:"iscsiName,omitempty"`
}

// DeleteServerHBAParams represents parameters for deleting server HBAs (via URL)
type DeleteServerHBAParams struct {
	ServerID      int    `json:"serverId"`
	InitiatorName string `json:"initiatorName"`
}
