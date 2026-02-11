package terraform

type Pool struct {
	PoolID                   int    `json:"poolId"`
	PoolStatus               string `json:"poolStatus"`
	UsedCapacityRate         int    `json:"usedCapacityRate"`
	UsedPhysicalCapacityRate int    `json:"usedPhysicalCapacityRate"`
	SnapshotCount            int    `json:"snapshotCount"`
	PoolName                 string `json:"poolName"`
}

type Pools struct {
	Data []Pool `json:"data"`
}
