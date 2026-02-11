package terraform

// SupportedHostModes is the supported host modes response.
type SupportedHostModes struct {
	HostModes       []SupportedHostMode       `json:"hostModes"`
	HostModeOptions []SupportedHostModeOption `json:"hostModeOptions"`
}

// SupportedHostMode is a host mode item.
type SupportedHostMode struct {
	HostModeID      int    `json:"hostModeId"`
	HostModeName    string `json:"hostModeName"`
	HostModeDisplay string `json:"hostModeDisplay"`
}

// SupportedHostModeOption is a host mode option item.
type SupportedHostModeOption struct {
	HostModeOptionID          int    `json:"hostModeOptionId"`
	HostModeOptionDescription string `json:"hostModeOptionDescription"`
	Scope                     string `json:"scope"`
	RequiredHostModes         []int  `json:"requiredHostModes"`
}
