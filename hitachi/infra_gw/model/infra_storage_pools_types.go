package infra_gw

type StoragePoolInfo struct {
	ResourceId             string `json:"resourceId"`
	PoolId                 int    `json:"poolId"`
	LdevIds                []int  `json:"ldevIds"`
	Name                   string `json:"name"`
	DepletionThresholdRate int    `json:"depletionThresholdRate"`
	DpVolumes              []struct {
		LogicalUnitId int    `json:"logicalUnitId"`
		Size          string `json:"size"`
	} `json:"dpVolumes"`
	FreeCapacity                  int64  `json:"freeCapacity"`
	FreeCapacityInUnits           string `json:"freeCapacityInUnits"`
	ReplicationDepletionAlertRate int    `json:"replicationDepletionAlertRate"`
	ReplicationUsageRate          int    `json:"replicationUsageRate"`
	ResourceGroupId               int    `json:"resourceGroupId"`
	Status                        string `json:"status"`
	SubscriptionLimitRate         int    `json:"subscriptionLimitRate"`
	SubscriptionRate              int    `json:"subscriptionRate"`
	SubscriptionWarningRate       int    `json:"subscriptionWarningRate"`
	TotalCapacity                 int64  `json:"totalCapacity"`
	TotalCapacityInUnit           string `json:"totalCapacityInUnit"`
	Type                          string `json:"type"`
	VirtualVolumeCount            int    `json:"virtualVolumeCount"`
	WarningThresholdRate          int    `json:"warningThresholdRate"`
	DeduplicationEnabled          bool   `json:"deduplicationEnabled"`
}

type StoragePools struct {
	Path    string            `json:"path"`
	Message string            `json:"message"`
	Data    []StoragePoolInfo `json:"data"`
}

type StoragePool struct {
	Path    string          `json:"path"`
	Message string          `json:"message"`
	Data    StoragePoolInfo `json:"data"`
}
