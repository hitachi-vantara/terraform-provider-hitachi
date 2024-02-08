package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetAllPartners gets partners information
func (psm *infraGwManager) GetAllPartners() (*[]model.Partner, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var partners []model.Partner
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = "/partners"

	err := httpmethod.GetCall(psm.setting, apiSuf, &partners)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &partners, nil
}

// GetPartner gets partner information
func (psm *infraGwManager) GetPartner(partnerId string) (*model.Partner, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var partner model.Partner
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = fmt.Sprintf("partner/%s", partnerId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &partner)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &partner, nil
}

// GetAllSubscribers gets all subscribers information
func (psm *infraGwManager) GetAllSubscribers(partnerId string) (*model.Subscribers, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var subscribers model.Subscribers
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = fmt.Sprintf("partner/%s/subscribers", partnerId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &subscribers)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &subscribers, nil
}

// GetSubscriber gets  subscriber information
func (psm *infraGwManager) GetSubscriber(partnerId string, subscriberId string) (*model.Subscriber, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var subscriber model.Subscriber
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = fmt.Sprintf("partner/%s/subscriber/%s", partnerId, subscriberId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &subscriber)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &subscriber, nil
}

// GetSubscriberResources gets  subscriber resources information
func (psm *infraGwManager) GetSubscriberResources(partnerId string, subscriberId string) (*model.SubscriberDetails, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var subscriber model.SubscriberDetails
	var apiSuf string
	psm.setting.V3API = true

	apiSuf = fmt.Sprintf("partner/%s/subscriber/%s/resources", partnerId, subscriberId)

	err := httpmethod.GetCall(psm.setting, apiSuf, &subscriber)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &subscriber, nil
}

// RegisterSubscriber registers a subscriber
func (psm *infraGwManager) RegisterSubscriber(reqBody *model.RegisterSubscriberReq) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	apiSuf := "/register/subscriber"

	resp, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resp, nil
}

// RegisterPartner registers a new partner
func (psm *infraGwManager) RegisterPartner(reqBody *model.RegisterPartnerReq) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	apiSuf := "/register/partner"

	resp, err := httpmethod.PostCall(psm.setting, apiSuf, reqBody)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resp, nil
}

// UnRegisterSubscriber removes a subscriber from the subscription
func (psm *infraGwManager) UnRegisterSubscriber(subscriberId string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	apiSuf := fmt.Sprintf("unregister/subscriber/%s", subscriberId)

	resp, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return resp, nil
}
