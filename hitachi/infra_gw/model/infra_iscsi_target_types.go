package infra_gw

type IscsiTargetInfo struct {
	ResourceId string `json:"resourceId"`
	PortId     string `json:"portId"`
	HostMode   struct {
		HostCommonSettings string `json:"hostCommonSettings"`
		HostMiddleWare     string `json:"hostMiddleWare"`
		HostMode           string `json:"hostMode"`
		HostModeOptions    []struct {
			DfOption         string `json:"dfOption"`
			IsAMSLegal       bool   `json:"isAMSLegal"`
			IsDF             bool   `json:"isDF"`
			IsHUSLegal       bool   `json:"isHUSLegal"`
			IsRAID           bool   `json:"isRAID"`
			RaidOption       string `json:"raidOption"`
			RaidOptionNumber int    `json:"raidOptionNumber"`
		} `json:"hostModeOptions"`
		HostPlatformOption string `json:"hostPlatformOption"`
		IsDF               bool   `json:"isDF"`
		IsRAID             bool   `json:"isRAID"`
		RaidHostModeChar   string `json:"raidHostModeChar"`
	} `json:"hostMode"`
	ResourceGroupId int      `json:"resourceGroupId"`
	TargetUser      string   `json:"targetUser"`
	Iqn             string   `json:"iqn"`
	IqnInitiators   []string `json:"iqnInitiators"`
	LogicalUnits    []struct {
		HostLunId              int    `json:"hostLunId"`
		LogicalUnitId          int    `json:"logicalUnitId"`
		LogicalUnitIdHexFormat string `json:"logicalUnitIdHexFormat"`
	} `json:"logicalUnits"`
	AuthParam struct {
		IsChapEnabled      bool   `json:"isChapEnabled"`
		IsChapRequired     bool   `json:"isChapRequired"`
		IsMutualAuth       bool   `json:"isMutualAuth"`
		AuthenticationMode string `json:"authenticationMode"`
	} `json:"authParam"`
	ChapUsers []string `json:"chapUsers"`
	ISCSIName string   `json:"iSCSIName"`
	ISCSIId   int      `json:"iSCSIId"`
}

type IscsiTargets struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []IscsiTargetInfo `json:"data"`
}

type IscsiTarget struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    IscsiTargetInfo `json:"data"`
}
