package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetParityGroups is used to fetch all parity group details with optional detailInfoType and also we can filter by ids
func (psm *sanStorageManager) GetParityGroups(detailInfoType string, driveTypeName string, clprId *int, parityGroupIds ...[]string) (*[]sanmodel.ParityGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), objStorage.Serial)
	var provParityGroups *[]provisonermodel.ParityGroup
	if detailInfoType != "" {
		provParityGroups, err = provObj.GetParityGroups(detailInfoType, driveTypeName, clprId, parityGroupIds...)
	} else {
		provParityGroups, err = provObj.GetParityGroups("", driveTypeName, clprId, parityGroupIds...)
	}
	if err != nil {
		log.WriteDebug("TFError| error in GetParityGroups provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcileParityGroups := []sanmodel.ParityGroup{}
	err = copier.Copy(&reconcileParityGroups, provParityGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), objStorage.Serial)
	return &reconcileParityGroups, nil
}

func (psm *sanStorageManager) GetParityGroup(parityGroupId string) (*sanmodel.ParityGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), objStorage.Serial)
	provParityGroup, err := provObj.GetParityGroup(parityGroupId)
	if err != nil {
		log.WriteDebug("TFError| error in GetParityGroups provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcileParityGroup := sanmodel.ParityGroup{}
	err = copier.Copy(&reconcileParityGroup, provParityGroup)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), objStorage.Serial)
	return &reconcileParityGroup, nil
}
