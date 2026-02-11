package sanstorage

// HostModeAndOptions represents supported host modes and host mode options.
type HostModeAndOptions struct {
	HostModes       []HostModes       `json:"hostModes"`
	HostModeOptions []HostModeOptions `json:"hostModeOptions"`
}

// HostModes represents a supported host mode.
type HostModes struct {
	HostModeID      int    `json:"hostModeId"`
	HostModeName    string `json:"hostModeName"`
	HostModeDisplay string `json:"hostModeDisplay"`
}

// HostModeOptions represents a supported host mode option.
type HostModeOptions struct {
	HostModeOptionID          int    `json:"hostModeOptionId"`
	HostModeOptionDescription string `json:"hostModeOptionDescription"`
	Scope                     string `json:"scope"`
	RequiredHostModes         []int  `json:"requiredHostModes"`
}
