package sanstorage

// response of objects/storages/instance
type StorageSystemInfo struct {
	StorageDeviceID    string `json:"storageDeviceId"`
	Model              string `json:"model"`
	SerialNumber       int    `json:"serialNumber"`
	IP                 string `json:"ip"`
	Ctl1IP             string `json:"ctl1Ip"`
	Ctl2IP             string `json:"ctl2Ip"`
	DkcMicroVersion    string `json:"dkcMicroVersion"`
	CommunicationModes []struct {
		CommunicationMode string `json:"communicationMode"`
	} `json:"communicationModes"`
	IsSecure bool `json:"isSecure"`
}

// VSP 5000 only
// response of objects/storage-summaries/instance
// TODO: add paritygroup info
type StorageSystemSummary struct {
	Name                                      string `json:"name"`
	SvpMicroVersion                           string `json:"svpMicroVersion"`
	RmiServerVersion                          string `json:"rmiServerVersion"`
	NumOfDiskBoards                           int    `json:"numOfDiskBoards"`
	CacheMemoryCapacity                       int    `json:"cacheMemoryCapacity"`
	NumOfSpareDrives                          int    `json:"numOfSpareDrives"`
	TotalOpenVolumeCapacity                   int    `json:"totalOpenVolumeCapacity"`
	TotalOpenVolumeCapacityInKB               uint64 `json:"totalOpenVolumeCapacityInKB"`
	AllocatedOpenVolumeCapacity               int    `json:"allocatedOpenVolumeCapacity"`
	AllocatedOpenVolumeCapacityInKB           int    `json:"allocatedOpenVolumeCapacityInKB"`
	AllocatableOpenVolumeCapacity             int    `json:"allocatableOpenVolumeCapacity"`
	AllocatableOpenVolumeCapacityInKB         int    `json:"allocatableOpenVolumeCapacityInKB"`
	UnallocatedOpenVolumeCapacity             int    `json:"unallocatedOpenVolumeCapacity"`
	UnallocatedOpenVolumeCapacityInKB         uint64 `json:"unallocatedOpenVolumeCapacityInKB"`
	ReservedOpenVolumeCapacity                int    `json:"reservedOpenVolumeCapacity"`
	ReservedOpenVolumeCapacityInKB            uint64 `json:"reservedOpenVolumeCapacityInKB"`
	AllocatedOpenVolumePhysicalCapacity       int    `json:"allocatedOpenVolumePhysicalCapacity"`
	AllocatedOpenVolumePhysicalCapacityInKB   int    `json:"allocatedOpenVolumePhysicalCapacityInKB"`
	AllocatableOpenVolumePhysicalCapacity     int    `json:"allocatableOpenVolumePhysicalCapacity"`
	AllocatableOpenVolumePhysicalCapacityInKB int    `json:"allocatableOpenVolumePhysicalCapacityInKB"`
	ReservedOpenVolumePhysicalCapacity        int    `json:"reservedOpenVolumePhysicalCapacity"`
	ReservedOpenVolumePhysicalCapacityInKB    int    `json:"reservedOpenVolumePhysicalCapacityInKB"`
	AllocatedMainframeVolumeCapacity          int    `json:"allocatedMainframeVolumeCapacity"`
	AllocatedMainframeVolumeCapacityInKB      int    `json:"allocatedMainframeVolumeCapacityInKB"`
	ReservedMainframeVolumeCapacity           int    `json:"reservedMainframeVolumeCapacity"`
	ReservedMainframeVolumeCapacityInKB       int    `json:"reservedMainframeVolumeCapacityInKB"`
	TotalAllocatedVolumeCapacity              int    `json:"totalAllocatedVolumeCapacity"`
	TotalAllocatedVolumeCapacityInKB          int    `json:"totalAllocatedVolumeCapacityInKB"`
	TotalUnallocatedVolumeCapacity            int    `json:"totalUnallocatedVolumeCapacity"`
	TotalUnallocatedVolumeCapacityInKB        uint64 `json:"totalUnallocatedVolumeCapacityInKB"`
	TotalReservedVolumeCapacity               int    `json:"totalReservedVolumeCapacity"`
	TotalReservedVolumeCapacityInKB           uint64 `json:"totalReservedVolumeCapacityInKB"`
	TotalMainframeVolumeCapacity              int    `json:"totalMainframeVolumeCapacity"`
	TotalMainframeVolumeCapacityInKB          int    `json:"totalMainframeVolumeCapacityInKB"`
	TotalVolumeCapacity                       int    `json:"totalVolumeCapacity"`
	TotalVolumeCapacityInKB                   uint64 `json:"totalVolumeCapacityInKB"`
	NumOfOpenVolumes                          int    `json:"numOfOpenVolumes"`
	NumOfAllocatedOpenVolumes                 int    `json:"numOfAllocatedOpenVolumes"`
	NumOfAllocatableOpenVolumes               int    `json:"numOfAllocatableOpenVolumes"`
	NumOfReservedOpenVolumes                  int    `json:"numOfReservedOpenVolumes"`
}

type StorageCapacity struct {
	Internal struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"internal"`
	External struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"external"`
	Total struct {
		FreeSpace     uint64 `json:"freeSpace"`
		TotalCapacity uint64 `json:"totalCapacity"`
	} `json:"total"`
}
