package vssbstorage

type AddressSetting struct {
	Index                    int    `json:"index"`
	ControlPortIPv4Address   string `json:"control_port_ipv4_address"`
	InternodePortIPv4Address string `json:"internode_port_ipv4_address"`
	ComputePortIPv4Address   string `json:"compute_port_ipv4_address"`
	ComputePortIPv6Address   string `json:"compute_port_ipv6_address,omitempty"`
}

type CreateConfigurationFileParam struct {
	ExportFileType        string           `json:"export_file_type,omitempty"`
	MachineImageID        string           `json:"machine_image_id,omitempty"`
	NumberOfDrives        int              `json:"number_of_drives,omitempty"`
	RecoverSingleDrive    bool             `json:"recover_single_drive,omitempty"`
	DriveID               string           `json:"drive_id,omitempty"`
	RecoverSingleNode     bool             `json:"recover_single_node,omitempty"`
	NodeID                string           `json:"node_id,omitempty"`
	AddressSetting        []AddressSetting `json:"address_setting,omitempty"`
}
