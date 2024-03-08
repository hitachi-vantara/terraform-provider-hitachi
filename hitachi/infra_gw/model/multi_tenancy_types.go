package infra_gw

const (
	AdminRole             = "ucpAdminRole"
	StorageAdminRole      = "UcpAdvisorStorageAdmin"
	DefaultSubscriberName = "TerraformSubscriber"
	DefaultPartnerName    = "TerraformPartner"
)

type Partner struct {
	PartnerID string `json:"partnerId"`
	Type      string `json:"type"`
	Time      int64  `json:"time"`
}

type Partners struct {
	Partners []Partner `json:"partners"`
}

type RegisterSubscriberReq struct {
	Name         string `json:"name"`
	PartnerID    string `json:"partnerId"`
	SubscriberID string `json:"subscriberId"`
	Description  string `json:"description"`
}

type RegisterPartnerReq struct {
	Name        string `json:"name"`
	PartnerID   string `json:"partnerId"`
	Description string `json:"description"`
}

type Subscriber struct {
	SubscriberId       string  `json:"subscriberId"`
	PartnerID          string  `json:"partnerId"`
	Type               string  `json:"type"`
	Time               int64   `json:"time"`
	Name               string  `json:"name"`
	SoftLimit          string  `json:"softLimit"`
	HardLimit          string  `json:"hardLimit"`
	QuotaLimit         string  `json:"quotaLimit"`
	UsedQuotaInPercent float64 `json:"usedQuotaInPercent"`
}

type Subscribers struct {
	Subscribers *[]Subscriber `json:"subscribers"`
}

type SubscriberDetails struct {
	DeviceId      string `json:"deviceId"`
	ResourceId    string `json:"resourceId"`
	SubscriberId  string `json:"subscriberId"`
	PartnerID     string `json:"partnerId"`
	Type          string `json:"type"`
	Time          int64  `json:"time"`
	ResourceValue string `json:"resourceValue"`
}

type MTDetails struct {
	PartnerId    *string `json:"partnerId"`
	SubscriberId *string `json:"subscriberId"`
}

type UpdateSubscriberReq struct {
	Name       string `json:"name"`
	SoftLimit  string `json:"softLimit"`
	HardLimit  string `json:"hardLimit"`
	QuotaLimit string `json:"quotaLimit"`
}
