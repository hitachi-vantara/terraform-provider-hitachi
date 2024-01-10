package vssbstorage

type PortDetailSettings struct {
	Port         StoragePort      `json:"storagePort"`
	AuthSettings PortAuthSettings `json:"authSettings"`
	ChapUsers    ChapUsers        `json:"chapUsers"`
}

type IscsiPortAuthSettings struct {
	AuthSettings PortAuthSettings `json:"authSettings"`
	ChapUsers    ChapUsers        `json:"chapUsers"`
}
