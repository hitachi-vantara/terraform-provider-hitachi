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
	ChapUsers         []string `json:"chapUsers"`
	ISCSIName         string   `json:"iSCSIName"`
	ISCSIId           int      `json:"iSCSIId"`
	SubscriberId      string   `json:"subscriberId"`
	PartnerId         string   `json:"partnerId"`
	StorageId         string   `json:"storageId"`
	Time              int64    `json:"time"`
	EntitlementStatus string   `json:"entitlementStatus"`
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

type CreateIscsiTargetParam struct {
	IscsiName          string     `json:"iscsiName,omitempty"`
	HostMode           string     `json:"hostMode,omitempty"`
	HostModeOptions    []int      `json:"hostModeOptions,omitempty"`
	Port               string     `json:"port,omitempty"`
	IqnInitiators      []string   `json:"iqnInitiators,omitempty"`
	Luns               []int      `json:"ldevs,omitempty"`
	Iqn                string     `json:"iqn,omitempty"`
	AuthenticationMode string     `json:"authenticationMode,omitempty"`
	IsMutualAuth       bool       `json:"isMutualAuth"`
	ChapUsers          []ChapUser `json:"chapUsers,omitempty"`
	UcpSystem          string     `json:"ucpSystem,omitempty"`
}

type ChapUser struct {
	ChapSecret string `json:"chapSecret,omitempty"`
	ChapUserId string `json:"chapUserId,omitempty"`
}

type AddVolumesToIscsiTargetParam struct {
	LdevIds []int `json:"ldevIds,omitempty"`
}

type RemoveVolumesFromIscsiTargetParam AddVolumesToIscsiTargetParam

type AddIqnInitiatorsToIscsiTargetParam struct {
	IqnInitiators []string `json:"iqnInitiators,omitempty"`
}

type RemoveIqnInitiatorsFromIscsiTargetParam AddIqnInitiatorsToIscsiTargetParam
type UpdateTargetIqnInIscsiTargetParam AddIqnInitiatorsToIscsiTargetParam

type UpdateHostModeParam struct {
	HostMode        string `json:"hostMode,omitempty"`
	HostModeOptions []int  `json:"hostModeOptions,omitempty"`
}

type SetChapUserParam ChapUser
type UpdateChapUserParam ChapUser
