package sanstorage

type LunRequest struct {
	LdevID                             *int    `json:"ldevId,omitempty"`
	Name                               *string `json:"name,omitempty"`
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
