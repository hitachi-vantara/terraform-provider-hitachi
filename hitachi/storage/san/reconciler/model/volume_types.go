package sanstorage

type LunRequest struct {
	LdevID *int    `json:"ldevId,omitempty"`
	Name   *string `json:"name,omitempty"`
	// Mainframe-specific create inputs
	Cylinder                           *int    `json:"cylinder,omitempty"`
	EmulationType                      *string `json:"emulationType,omitempty"`
	IsTseVolume                        *bool   `json:"isTseVolume,omitempty"`
	IsEseVolume                        *bool   `json:"isEseVolume,omitempty"`
	ClprID                             *int    `json:"clprId,omitempty"`
	MpBladeID                          *int    `json:"mpBladeId,omitempty"`
	Ssid                               *string `json:"ssid,omitempty"`
	BlockCapacity                      *int64  `json:"blockCapacity,omitempty"`
	PoolID                             *int    `json:"poolId"`
	ParityGroupID                      *string `json:"parityGroupId,omitempty"`
	ExternalParityGroupID              *string `json:"externalParityGroupId,omitempty"`
	ByteFormatCapacity                 string  `json:"byteFormatCapacity"`
	DataReductionMode                  *string `json:"dataReductionMode,omitempty"`
	IsDataReductionSharedVolumeEnabled *bool   `json:"isDataReductionSharedVolumeEnabled,omitempty"`
	IsCompressionAccelerationEnabled   *bool   `json:"isCompressionAccelerationEnabled,omitempty"`
}

type UpdateLunRequest struct {
	LdevID                           *int    `json:"ldevId,omitempty"`
	Name                             *string `json:"name,omitempty"`
	ByteFormatCapacity               *string `json:"byteFormatCapacity"`
	DataReductionMode                *string `json:"dataReductionMode,omitempty"`
	DataReductionProcessMode         *string `json:"dataReductionProcessMode,omitempty"`
	IsCompressionAccelerationEnabled *bool   `json:"isCompressionAccelerationEnabled,omitempty"`
	IsAluaEnabled                    *bool   `json:"isAluaEnabled,omitempty"`
}

type FormatLdevRequest struct {
	OperationType              *string `json:"operationType,omitempty"`
	IsDataReductionForceFormat *bool   `json:"isDataReductionForceFormat,omitempty"`
}
