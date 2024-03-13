package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/infra_gw/gateway/http"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetStorageResource gets storage resource information
func (psm *infraGwManager) GetStorageResource(storageId, resourceType, resourceId string) (*model.StorageResourceResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	var StorageResources model.StorageResourceResponse

	apiSuf := fmt.Sprintf("/storage/%s/resource/%s?type=%s", storageId, resourceId, resourceType)
	err := httpmethod.GetCall(psm.setting, apiSuf, nil, &StorageResources)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	log.WriteDebug("MT Storage Devices %v", StorageResources)
	return &StorageResources, nil
}

// AddStorageResource adds storage resource information
func (psm *infraGwManager) AddStorageResource(storageId string, reqData *model.AddStorageResourceRequest) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	apiSuf := fmt.Sprintf("/storage/%s/resource/", storageId)
	resourceId, err := httpmethod.PostCall(psm.setting, apiSuf, &reqData, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	log.WriteDebug("MT Storage resource id %v", resourceId)
	return resourceId, nil
}

func (psm *infraGwManager) RemoveStorageResource(storageId, resourceId, subscriberId, Type string) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	psm.setting.V3API = true

	apiSuf := fmt.Sprintf("/storage/%s/resource/%s?type=%s&subscriberid=%s", storageId, resourceId, Type, subscriberId)
	id, err := httpmethod.DeleteCall(psm.setting, apiSuf, nil, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return id, nil
}
