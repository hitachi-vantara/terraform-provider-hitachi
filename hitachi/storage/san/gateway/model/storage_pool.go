package sanstorage

type Pool struct {
	PoolID     int    `json:"poolId"`
	PoolStatus string `json:"poolStatus"`
	PoolName   string `json:"poolName"`
}

type Pools struct {
	Data []Pool `json:"data"`
}
