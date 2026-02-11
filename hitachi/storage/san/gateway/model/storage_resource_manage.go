package sanstorage

type LockResourcesReq struct {
	Parameters struct {
		WaitTime int `json:"waitTime"`
	} `json:"parameters"`
}
