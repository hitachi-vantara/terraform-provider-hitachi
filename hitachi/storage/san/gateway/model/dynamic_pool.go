package sanstorage

type DynamicPool struct {
	PoolID                                         int    `json:"poolId"`
	PoolStatus                                     string `json:"poolStatus"`
	UsedCapacityRate                               int    `json:"usedCapacityRate"`
	UsedPhysicalCapacityRate                       int    `json:"usedPhysicalCapacityRate"`
	SnapshotCount                                  int    `json:"snapshotCount"`
	PoolName                                       string `json:"poolName"`
	AvailableVolumeCapacity                        int    `json:"availableVolumeCapacity"`
	AvailablePhysicalVolumeCapacity                int    `json:"availablePhysicalVolumeCapacity"`
	TotalPoolCapacity                              int    `json:"totalPoolCapacity"`
	TotalPhysicalCapacity                          int    `json:"totalPhysicalCapacity"`
	NumOfLdevs                                     int    `json:"numOfLdevs"`
	FirstLdevID                                    int    `json:"firstLdevId"`
	WarningThreshold                               int    `json:"warningThreshold"`
	DepletionThreshold                             int    `json:"depletionThreshold"`
	VirtualVolumeCapacityRate                      int    `json:"virtualVolumeCapacityRate"`
	IsMainframe                                    bool   `json:"isMainframe"`
	IsShrinking                                    bool   `json:"isShrinking"`
	LocatedVolumeCount                             int    `json:"locatedVolumeCount"`
	TotalLocatedCapacity                           int    `json:"totalLocatedCapacity"`
	BlockingMode                                   string `json:"blockingMode"`
	TotalReservedCapacity                          int    `json:"totalReservedCapacity"`
	ReservedVolumeCount                            int    `json:"reservedVolumeCount"`
	PoolType                                       string `json:"poolType"`
	DuplicationNumber                              int    `json:"duplicationNumber"`
	EffectiveCapacity                              int    `json:"effectiveCapacity"`
	DataReductionAccelerateCompCapacity            int    `json:"dataReductionAccelerateCompCapacity"`
	DataReductionCapacity                          int    `json:"dataReductionCapacity"`
	DataReductionBeforeCapacity                    int    `json:"dataReductionBeforeCapacity"`
	DataReductionAccelerateCompRate                int    `json:"dataReductionAccelerateCompRate"`
	DuplicationRate                                int    `json:"duplicationRate"`
	CompressionRate                                int    `json:"compressionRate"`
	DataReductionRate                              int    `json:"dataReductionRate"`
	DataReductionAccelerateCompIncludingSystemData struct {
		IsReductionCapacityAvailable bool `json:"isReductionCapacityAvailable"`
		ReductionCapacity            int  `json:"reductionCapacity"`
		IsReductionRateAvailable     bool `json:"isReductionRateAvailable"`
		ReductionRate                int  `json:"reductionRate"`
	} `json:"dataReductionAccelerateCompIncludingSystemData"`
	DataReductionIncludingSystemData struct {
		IsReductionCapacityAvailable bool `json:"isReductionCapacityAvailable"`
		ReductionCapacity            int  `json:"reductionCapacity"`
		IsReductionRateAvailable     bool `json:"isReductionRateAvailable"`
		ReductionRate                int  `json:"reductionRate"`
	} `json:"dataReductionIncludingSystemData"`
	SnapshotUsedCapacity int    `json:"snapshotUsedCapacity"`
	SuspendSnapshot      bool   `json:"suspendSnapshot"`
	FormattedCapacity    int    `json:"formattedCapacity"`
	AutoAddPoolVol       string `json:"autoAddPoolVol"`
	Efficiency           struct {
		IsCalculated         bool   `json:"isCalculated"`
		TotalRatio           string `json:"totalRatio"`
		CompressionRatio     string `json:"compressionRatio"`
		SnapshotRatio        string `json:"snapshotRatio"`
		ProvisioningRate     string `json:"provisioningRate"`
		CalculationStartTime string `json:"calculationStartTime"`
		CalculationEndTime   string `json:"calculationEndTime"`
		DedupAndCompression  struct {
			TotalRatio       string `json:"totalRatio"`
			CompressionRatio string `json:"compressionRatio"`
			DedupeRatio      string `json:"dedupeRatio"`
			ReclaimRatio     string `json:"reclaimRatio"`
		} `json:"dedupAndCompression"`
		AcceleratedCompression struct {
			TotalRatio       string `json:"totalRatio"`
			CompressionRatio string `json:"compressionRatio"`
			ReclaimRatio     string `json:"reclaimRatio"`
		} `json:"acceleratedCompression"`
	} `json:"efficiency"`
	CapacitiesExcludingSystemData struct {
		UsedVirtualVolumeCapacity int `json:"usedVirtualVolumeCapacity"`
		CompressedCapacity        int `json:"compressedCapacity"`
		DedupedCapacity           int `json:"dedupedCapacity"`
		ReclaimedCapacity         int `json:"reclaimedCapacity"`
		SystemDataCapacity        int `json:"systemDataCapacity"`
		PreUsedCapacity           int `json:"preUsedCapacity"`
		PreCompressedCapacity     int `json:"preCompressedCapacity"`
		PreDedupredCapacity       int `json:"preDedupredCapacity"`
	} `json:"capacitiesExcludingSystemData"`
}

type DynamicPools struct {
	Data []DynamicPool `json:"data"`
}
