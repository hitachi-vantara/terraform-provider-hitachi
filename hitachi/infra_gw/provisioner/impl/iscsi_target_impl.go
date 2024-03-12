package infra_gw

import (
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	mc "terraform-provider-hitachi/hitachi/infra_gw/provisioner/message-catalog"
)

// GetIscsiTargets gets IscsiTargets information
func (psm *infraGwManager) GetIscsiTargets(id string, port string) (*model.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_ISCSI_TARGETS_BEGIN))
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.INFO_GET_INFRA_ISCSI_TARGETS_END))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err

	}

	if psm.setting.PartnerId != nil {
		return gatewayObj.GetMTIscsiTargets(id)
	}

	return gatewayObj.GetIscsiTargets(id, port)
}

// GetIscsiTarget gets IscsiTarget information
func (psm *infraGwManager) GetIscsiTarget(id string, iscsiTargetId string) (*model.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_INFRA_ISCSI_TARGET_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_INFRA_ISCSI_TARGET_FAILED))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	if psm.setting.PartnerId != nil {
		return gatewayObj.GetMTIscsiTarget(id, iscsiTargetId)
	}

	return gatewayObj.GetIscsiTarget(id, iscsiTargetId)
}

// CreateIscsiTarget .
func (psm *infraGwManager) CreateIscsiTarget(storageId string, reqBody model.CreateIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_INFRA_ISCSI_TARGET_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_INFRA_ISCSI_TARGET_FAILED))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	if psm.setting.PartnerId != nil {
		// Check if the resource is attached to the subscriber or not
		err := psm.TagResourceToSubscriber(storageId, reqBody.Port)
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_INFRA_ISCSI_TARGET_FAILED))
			log.WriteDebug("TFError| error in TagResourceToSubscriber call, err: %v", err)
			return nil, err
		}
		return gatewayObj.CreateMTIscsiTarget(storageId, reqBody)
	}

	return gatewayObj.CreateIscsiTarget(storageId, reqBody)
}

// UpdateIscsiTarget .
func (psm *infraGwManager) UpdateIscsiTarget(storageId, hostGroupId string, reqBody model.CreateIscsiTargetParam) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_INFRA_ISCSI_TARGET_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_INFRA_ISCSI_TARGET_FAILED))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	updateReqBody := model.UpdateHostModeParam{
		HostMode:        reqBody.HostMode,
		HostModeOptions: reqBody.HostModeOptions,
	}

	return gatewayObj.UpdateHostMode(storageId, hostGroupId, updateReqBody)
}

func (psm *infraGwManager) TagResourceToSubscriber(storageId, resourceId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_INFRA_ISCSI_TARGET_BEGIN))

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_INFRA_ISCSI_TARGET_FAILED))

		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	_, err = gatewayObj.GetStorageResource(storageId, model.IscsiTargetPort, resourceId)

	if err != nil && strings.Contains(err.Error(), "Resource not found") {
		reqPayload := model.AddStorageResourceRequest{
			ResourceType: model.IscsiTargetPort,
			ResourceId:   resourceId,
			SubscriberId: *psm.setting.SubscriberId,
			PartnerId:    *psm.setting.PartnerId,
		}
		log.WriteDebug("TFError| error in GetStorageResource call, err: %v", err)
		_, err = gatewayObj.AddStorageResource(storageId, &reqPayload)
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_UPDATE_INFRA_ISCSI_TARGET_FAILED))

			log.WriteDebug("TFError| error in AddStorageResource call, err: %v", err)
			return err
		}
	} else if err != nil {

		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_INFRA_ISCSI_TARGET_FAILED))

		log.WriteDebug("TFError| error in AddStorageResource call, err: %v", err)
		return err

	}
	return nil
}
