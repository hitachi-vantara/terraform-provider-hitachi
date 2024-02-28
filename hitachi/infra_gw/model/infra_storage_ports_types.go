package infra_gw

type StoragePortInfo struct {
	PortId            string `json:"portId"`
	Type              string `json:"type"`
	Speed             string `json:"speed"`
	ResourceGroupId   int    `json:"resourceGroupId"`
	Wwn               string `json:"wwn"`
	ResourceId        string `json:"resourceId"`
	Attribute         string `json:"attribute"`
	ConnectionType    string `json:"connectionType"`
	FabricOn          bool   `json:"fabricOn"`
	Mode              string `json:"mode"`
	IsSecurityEnabled bool   `json:"isSecurityEnabled"`
}

type StoragePorts struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []StoragePortInfo `json:"data"`
}

type StoragePort struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    StoragePortInfo `json:"data"`
}

type MTHurPairInfo struct {
	ResourceId               string `json:"resourceId"`
	ConsistencyGroupId       int    `json:"consistencyGroupId"`
	CopyRate                 int    `json:"copyRate"`
	FenceLevel               string `json:"fenceLevel"`
	MirrorUnitId             int    `json:"mirrorUnitId"`
	PairName                 string `json:"pairName"`
	PrimaryVolumeId          int    `json:"primaryVolumeId"`
	PrimaryVolumeStorageId   int    `json:"primaryVolumeStorageId"`
	SecondaryVolumeId        int    `json:"secondaryVolumeId"`
	SecondaryVolumeStorageId int    `json:"secondaryVolumeStorageId"`
	Status                   string `json:"status"`
	SvolAccessMode           string `json:"svolAccessMode"`
	Type                     string `json:"type"`
}

type MTSnapshotPairInfo struct {
	ConsistencyGroupId       int    `json:"consistencyGroupId"`
	CopyPaceTrackSize        string `json:"copyPaceTrackSize"`
	CopyRate                 int    `json:"copyRate"`
	MirrorUnitId             int    `json:"mirrorUnitId"`
	PrimaryHexVolumeId       string `json:"primaryHexVolumeId"`
	PrimaryVolumeId          int    `json:"primaryVolumeId"`
	StorageSerialNumber      int    `json:"storageSerialNumber"`
	SecondaryHexVolumeId     string `json:"secondaryHexVolumeId"`
	SecondaryVolumeId        int    `json:"secondaryVolumeId"`
	SecondaryVolumeStorageId int    `json:"secondaryVolumeStorageId"`
	Status                   string `json:"status"`
	SvolAccessMode           string `json:"svolAccessMode"`
	Type                     string `json:"type"`
}

type MTShadowImageInfo struct {
	ConsistencyGroupId       int    `json:"consistencyGroupId"`
	CopyPaceTrackSize        string `json:"copyPaceTrackSize"`
	CopyRate                 int    `json:"copyRate"`
	MirrorUnitId             int    `json:"mirrorUnitId"`
	PairName                 string `json:"pairName"`
	PrimaryHexVolumeId       string `json:"primaryHexVolumeId"`
	PrimaryVolumeId          int    `json:"primaryVolumeId"`
	PrimaryVolumeStorageId   int    `json:"primaryVolumeStorageId"`
	SecondaryHexVolumeId     string `json:"secondaryHexVolumeId"`
	SecondaryVolumeId        int    `json:"secondaryVolumeId"`
	SecondaryVolumeStorageId int    `json:"secondaryVolumeStorageId"`
	Status                   string `json:"status"`
	SvolAccessMode           string `json:"svolAccessMode"`
	Type                     string `json:"type"`
}

type MTIscsiTargetInfo struct {
	IscsiName       string `json:"iscsiName"`
	IscsiId         int    `json:"iscsiId"`
	ResourceGroupId int    `json:"resourceGroupId"`
	HostMode        string `json:"hostMode"`
	PortId          string `json:"portId"`
}

type MTStorageVolumeInfo struct {
	LdevId        int    `json:"ldevId"`
	PoolId        int    `json:"poolId"`
	TotalCapacity int    `json:"totalCapacity"`
	UsedCapacity  int    `json:"usedCapacity"`
	PoolName      string `json:"poolName"`
}

type MTHostGroupInfo struct {
	HostGroupName   string `json:"hostGroupName"`
	HostGroupId     int    `json:"hostGroupId"`
	ResourceGroupId int    `json:"resourceGroupId"`
	Port            string `json:"port"`
	HostMode        string `json:"hostMode"`
}

type DpVolume struct {
	LogicalUnitId int    `json:"logicalUnitId"`
	Size          string `json:"size"`
}

type MTStoragePoolInfo struct {
	PoolId                        int        `json:"poolId"`
	LdevIds                       []int      `json:"ldevIds"`
	PoolName                      string     `json:"poolName"`
	PoolType                      string     `json:"poolType"`
	DepletionThresholdRate        int        `json:"depletionThresholdRate"`
	DpVolumes                     []DpVolume `json:"dpVolumes"`
	FreeCapacity                  int        `json:"freeCapacity"`
	FreeCapacityInUnits           string     `json:"freeCapacityInUnits"`
	ReplicationDataReleasedRate   int        `json:"replicationDataReleasedRate"`
	ReplicationDepletionAlertRate int        `json:"replicationDepletionAlertRate"`
	ReplicationUsageRate          int        `json:"replicationUsageRate"`
	ResourceGroupId               int        `json:"resourceGroupId"`
	Status                        string     `json:"status"`
	SubscriptionLimitRate         int        `json:"subscriptionLimitRate"`
	SubscriptionRate              int        `json:"subscriptionRate"`
	SubscriptionWarningRate       int        `json:"subscriptionWarningRate"`
	TotalCapacity                 int        `json:"totalCapacity"`
	TotalCapacityInUnit           string     `json:"totalCapacityInUnit"`
	UtilizationRate               int        `json:"utilizationRate"`
	VirtualVolumeCount            int        `json:"virtualVolumeCount"`
	WarningThresholdRate          int        `json:"warningThresholdRate"`
	IsDeduplicationEnabled        bool       `json:"isDeduplicationEnabled"`
}

type MTPortInfo struct {
	PortType          string `json:"portType"`
	PortId            string `json:"portId"`
	Speed             string `json:"speed"`
	ResourceGroupId   int    `json:"resourceGroupId"`
	IsSecurityEnabled bool   `json:"isSecurityEnabled"`
	Wwn               string `json:"wwn"`
	Attribute         string `json:"attribute"`
	ConnectionType    string `json:"connectionType"`
	FabricOn          bool   `json:"fabricOn"`
	Mode              string `json:"mode"`
}

type MTStoragePortInfo struct {
	ResourceId        string              `json:"resourceId"`
	Type              string              `json:"type"`
	StorageId         string              `json:"storageId"`
	DeviceId          string              `json:"deviceId"`
	EntitlementStatus string              `json:"entitlementStatus"`
	PartnerId         string              `json:"partnerId"`
	SubscriberId      string              `json:"subscriberId"`
	PortInfo          MTPortInfo          `json:"portInfo"`
	StoragePoolInfo   MTStoragePoolInfo   `json:"storagePoolInfo"`
	HostGroupInfo     MTHostGroupInfo     `json:"hostGroupInfo"`
	StorageVolumeInfo MTStorageVolumeInfo `json:"storageVolumeInfo"`
	IscsiTargetInfo   MTIscsiTargetInfo   `json:"iscsiTargetInfo"`
	ShadowImageInfo   MTShadowImageInfo   `json:"shadowImageInfo"`
	SnapshotPairInfo  MTSnapshotPairInfo  `json:"snapshotPairInfo"`
	HurPairInfo       MTHurPairInfo       `json:"hurPairInfo"`
}
