package sanstorage

type PavAliasesResponse struct {
	Data []PavAlias `json:"data"`
}

type PavAlias struct {
	CuNumber      int    `json:"cuNumber"`
	LdevID        int    `json:"ldevId"`
	PavAttribute  string `json:"pavAttribute"`
	CBaseVolumeID *int   `json:"cBaseVolumeId,omitempty"`
	SBaseVolumeID *int   `json:"sBaseVolumeId,omitempty"`
}

type AssignPavAliasRequest struct {
	Parameters AssignPavAliasParams `json:"parameters"`
}
type AssignPavAliasParams struct {
	BaseLdevID   int   `json:"baseLdevId"`
	AliasLdevIDs []int `json:"aliasLdevIds"`
}

type UnassignPavAliasRequest struct {
	Parameters UnassignPavAliasParams `json:"parameters"`
}

type UnassignPavAliasParams struct {
	AliasLdevIDs []int `json:"aliasLdevIds"`
}
