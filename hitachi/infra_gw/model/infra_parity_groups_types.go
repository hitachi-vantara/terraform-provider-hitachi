package infra_gw

type ParityGroupInfo struct {
	ResourceId               string `json:"resourceId"`
	ParityGroupId            string `json:"parityGroupId"`
	FreeCapacity             int64  `json:"freeCapacity"`
	ResourceGroupId          int    `json:"resourceGroupId"`
	TotalCapacity            int64  `json:"totalCapacity"`
	LdevIds                  []int  `json:"ldevIds"`
	RaidLevel                string `json:"raidLevel"`
	DriveType                string `json:"driveType"`
	CopybackMode             bool   `json:"copybackMode"`
	Status                   string `json:"status"`
	IsPoolArrayGroup         bool   `json:"isPoolArrayGroup"`
	IsAcceleratedCompression bool   `json:"isAcceleratedCompression"`
	IsEncryptionEnabled      bool   `json:"isEncryptionEnabled"`
}

type ParityGroups struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []ParityGroupInfo `json:"data"`
}

type ParityGroup struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    ParityGroupInfo `json:"data"`
}
