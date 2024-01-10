package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetDynamicPools used to fetch all Parity group and with filtering also
func (psm *sanStorageManager) GetParityGroups(parityGroupIds ...[]string) (*[]sanmodel.ParityGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteDebug("TFDebug| Storage Serial:%d, ManagementIP:%s\n", psm.storageSetting.Serial, psm.storageSetting.MgmtIP)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), objStorage.Serial)
	parityGroups, err := gatewayObj.GetParityGroups()
	if err != nil {
		log.WriteDebug("TFError| error in GetParityGroups gateway call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), objStorage.Serial)
		return nil, err
	}

	provParityGroups := []sanmodel.ParityGroup{}
	err = copier.Copy(&provParityGroups, parityGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	filteredParityGroups := []sanmodel.ParityGroup{}
	// If Parity Group ID Filter available
	if len(parityGroupIds) > 0 {
		if parityGroupIds[0] != nil {
			for _, parity := range provParityGroups {
				for _, id := range parityGroupIds[0] {
					if parity.ParityGroupId == id {
						filteredParityGroups = append(filteredParityGroups, parity)
					}
				}
			}
		} else {
			filteredParityGroups = provParityGroups
		}

	} else {
		// If no filtering
		filteredParityGroups = provParityGroups
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), objStorage.Serial)
	return &filteredParityGroups, nil
}
