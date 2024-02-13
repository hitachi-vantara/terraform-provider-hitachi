package infra_gw

type VolumeInfo struct {
	ResourceId                   string `json:"resourceId"`
	DeduplicationCompressionMode string `json:"deduplicationCompressionMode"`
	EmulationType                string `json:"emulationType"`
	FormatOrShredRate            int    `json:"formatOrShredRate"`
	LdevId                       int    `json:"ldevId"`
	Name                         string `json:"name"`
	ParityGroupId                string `json:"parityGroupId"`
	PoolId                       int    `json:"poolId"`
	ResourceGroupId              int    `json:"resourceGroupId"`
	Status                       string `json:"status"`
	TotalCapacity                int64  `json:"totalCapacity"`
	UsedCapacity                 int64  `json:"usedCapacity"`
	VirtualStorageDeviceId       string `json:"virtualStorageDeviceId"`
	StripeSize                   int    `json:"stripeSize"`
	Type                         string `json:"type"`
	PathCount                    int    `json:"pathCount"`
	ProvisionType                string `json:"provisionType"`
	IsCommandDevice              bool   `json:"isCommandDevice"`
	LogicalUnitIdHexFormat       string `json:"logicalUnitIdHexFormat"`
	VirtualLogicalUnitId         int    `json:"virtualLogicalUnitId"`
	NaaId                        string `json:"naaId"`
	DedupCompressionProgress     int    `json:"dedupCompressionProgress"`
	DedupCompressionStatus       string `json:"dedupCompressionStatus"`
	IsALUA                       bool   `json:"isALUA"`
	IsDynamicPoolVolume          bool   `json:"isDynamicPoolVolume"`
	IsJournalPoolVolume          bool   `json:"isJournalPoolVolume"`
	IsPoolVolume                 bool   `json:"isPoolVolume"`
	PoolName                     string `json:"poolName"`
	QuorumDiskId                 int    `json:"quorumDiskId"`
	IsInGadPair                  bool   `json:"isInGadPair"`
	IsVVol                       bool   `json:"isVVol"`
}

type Volumes struct {
	Path    string       `json:"path"`
	Message string       `json:"message"`
	Data    []VolumeInfo `json:"data"`
}

type Volume struct {
	Path    string     `json:"path"`
	Message string     `json:"message"`
	Data    VolumeInfo `json:"data"`
}

type CreateVolumeParams struct {
	Name                         string `json:"name,omitempty"`
	PoolID                       int    `json:"poolId,omitempty"`
	ParityGroupId                string `json:"parityGroupId,omitempty"`
	Capacity                     string `json:"capacity,omitempty"`
	ResourceGroupId              int    `json:"resourceGroupId,omitempty"`
	LunId                        int    `json:"lunId,omitempty"`
	System                       string `json:"ucpSystem,omitempty"`
	DeduplicationCompressionMode string `json:"deduplicationCompressionMode,omitempty"`
}

type UpdateVolumeParams struct {
	Name                         string `json:"lunName,omitempty"`
	Capacity                     string `json:"capacity,omitempty"`
	DeduplicationCompressionMode string `json:"deduplicationCompressionMode,omitempty"`
}
