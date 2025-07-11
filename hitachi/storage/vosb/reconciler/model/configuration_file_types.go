package vssbstorage

type AddressSetting struct {
	Index                    int    `json:"index"`
	ControlPortIPv4Address   string `json:"controlPortIpv4Address"`
	InternodePortIPv4Address string `json:"internodePortIpv4Address"`
	ComputePortIPv4Address   string `json:"computePortIpv4Address"`
	ComputePortIPv6Address   string `json:"computePortIpv6Address,omitempty"`
}

type CreateConfigurationFileParam struct {
	ExportFileType        string           `json:"exportFileType,omitempty"`
	MachineImageID        string           `json:"machineImageId,omitempty"`
	NumberOfDrives        int              `json:"numberOfDrives,omitempty"`
	RecoverSingleDrive    bool             `json:"recoverSingleDrive,omitempty"`
	DriveID               string           `json:"driveId,omitempty"`
	RecoverSingleNode     bool             `json:"recoverSingleNode,omitempty"`
	NodeID                string           `json:"nodeId,omitempty"`
	AddressSetting        []AddressSetting `json:"addressSetting,omitempty"`
}
