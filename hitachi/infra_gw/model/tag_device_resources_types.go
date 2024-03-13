package infra_gw

const (
	IscsiTargetPort = "Port"
)

type StorageResources struct {
	DeviceId      string `json:"deviceId"`
	ResourceId    string `json:"resourceId"`
	SubscriberId  string `json:"subscriberId"`
	Type          string `json:"type"`
	ResourceValue string `json:"resourceValue"`
	Time          int64  `json:"time"`
}

type StorageResourceResponse struct {
	Path    string           `json:"path"`
	Message string           `json:"message"`
	Data    StorageResources `json:"data"`
}
type AddStorageResourceRequest struct {
	ResourceId   string `json:"resourceId"`
	SubscriberId string `json:"subscriberId"`
	ResourceType string `json:"type"`
	PartnerId    string `json:"partnerId"`
}
