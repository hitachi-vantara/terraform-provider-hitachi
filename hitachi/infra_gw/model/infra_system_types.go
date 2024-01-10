package infra_gw

type InfraGwSettings struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type InfraGwStorageSettingsAndInfo struct {
	Settings          InfraGwSettings   `json:"settings"`
	SerialToStorageId map[string]string `json:"serialToStorageId"`
	StorageIdToSerial map[string]string `json:"storageIdToSerial"`
	IscsiTargetIdMap  map[string]map[string]string
}
