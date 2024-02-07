package infra_gw

type Partner struct {
	PartnerID string `json:"partnerId"`
	Type      string `json:"type"`
	Time      int64  `json:"time"`
}

type Partners struct {
	Partners []*Partner `json:"partners"`
}

type RegisterSubscriberReq struct {
	Name         string `json:"name"`
	PartnerID    int64  `json:"partnerId"`
	SubscriberID int64  `json:"subscriberId"`
	Description  string `json:"description"`
}

type RegisterPartnerReq struct {
	Name        string `json:"name"`
	PartnerID   int64  `json:"partnerId"`
	Description string `json:"description"`
}

type Subscriber struct {
	SubscriberId string `json:"subscriberId"`
	PartnerID    int64  `json:"partnerId"`
	Type         string `json:"type"`
	Time         int64  `json:"time"`
}

type Subscribers struct {
	Subscribers *[]Subscriber `json:"subscribers"`
}

type SubscriberDetails struct {
	DeviceId      string `json:"deviceId"`
	ResourceId    string `json:"resourceId"`
	SubscriberId  string `json:"subscriberId"`
	PartnerID     int64  `json:"partnerId"`
	Type          string `json:"type"`
	Time          int64  `json:"time"`
	ResourceValue string `json:"resourceValue"`
}
