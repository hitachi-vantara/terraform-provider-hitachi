package infra_gw

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	provisonerimpl "terraform-provider-hitachi/hitachi/infra_gw/provisioner/impl"
)

// GetStorageDevices gets storage devices information
func (psm *infraGwManager) GetStorageDevices() (*model.StorageDevices, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStorageDevices()
}

// GetMTStorageDevices gets storage devices information
func (psm *infraGwManager) GetMTStorageDevices() (*[]model.MTStorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetMTStorageDevices()
}

// GetStorageDevice gets storage device information
func (psm *infraGwManager) GetStorageDevice(id string) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStorageDevice(id)
}

// GetMTStorageDevice gets storage device information
func (psm *infraGwManager) GetMTStorageDevice(id string) (*model.MTStorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetMTStorageDevice(id)
}

func (psm *infraGwManager) addStorageDevice(createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	storageId, err := provObj.AddStorageDevice(*createInput)
	if err != nil {
		log.WriteDebug("TFError| error in AddStorageDevice call, err: %v", err)
		return nil, err
	}
	log.WriteDebug("Return value of AddStorageDevice : %s", *storageId)

	if psm.setting.PartnerId != nil {

		reqBody := model.CreateMTStorageDeviceParam{
			ResourceId: *storageId,
			PartnerId:  *psm.setting.PartnerId,
		}
		id, err := provObj.AddMTStorageDevice(reqBody)
		if err != nil {
			log.WriteDebug("TFError| error in AddStorageDevice call, err: %v", err)
			return nil, err
		}
		log.WriteDebug("Return value of AddStorageDevice : %s", *id)
	}

	return psm.GetStorageDevice(*storageId)
}

func (psm *infraGwManager) updateStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	updateRequest := model.PatchStorageDeviceParam{
		Username:  createInput.Username,
		Password:  createInput.Password,
		OutOfBand: createInput.OutOfBand,
	}
	sId, err := provObj.UpdateStorageDevice(storageId, updateRequest)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateStorageDevice call, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| Return value of UpdateStorageDevice : %v", storageId)

	return psm.GetStorageDevice(*sId)
}

func (psm *infraGwManager) FindUcpSystemBySerial(serial string) (*model.UcpSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	ucpSystems, err := provObj.GetUcpSystems()
	if err != nil {
		e2 := fmt.Errorf("failed to get Ucp Systems, error code: %v", err)
		log.WriteDebug("TFError| error in GetUcpSystems call, err: %v", err)
		return nil, e2
	}

	var result model.UcpSystem
	for _, ucp := range ucpSystems.Data {
		if ucp.SerialNumber == serial {
			result.Path = ucpSystems.Path
			result.Message = ucpSystems.Message
			result.Data = ucp
			return &result, nil
		}
	}

	e2 := fmt.Errorf("UCP System with Serial Number %s does not exist", serial)
	return nil, e2

}

func (psm *infraGwManager) FindUcpSystemByName(name string) (*bool, *model.UcpSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, nil, err
	}
	ucpSystems, err := provObj.GetUcpSystems()
	if err != nil {
		e2 := fmt.Errorf("failed to get Ucp Systems, error code: %v", err)
		log.WriteDebug("TFError| error in GetUcpSystems call, err: %v", err)
		return nil, nil, e2
	}

	var result model.UcpSystem
	var found = false
	for _, ucp := range ucpSystems.Data {
		if ucp.Name == name {
			result.Path = ucpSystems.Path
			result.Message = ucpSystems.Message
			result.Data = ucp
			found = true
			return &found, &result, nil
		}
	}

	return &found, nil, nil

}

func (psm *infraGwManager) FindStorageSystemByNameAndSerial(name string, serialNumber string) (*model.UcpSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	ucpSystems, err := provObj.GetUcpSystems()
	if err != nil {
		e2 := fmt.Errorf("failed to get Ucp Systems, error code: %v", err)
		log.WriteDebug("TFError| error in GetUcpSystems call, err: %v", err)
		return nil, e2
	}

	var result model.UcpSystem
	for _, ucp := range ucpSystems.Data {
		// if ucp.Name == name && ucp.SerialNumber == serialNumber {
		// 	result.Path = ucpSystems.Path
		// 	result.Message = ucpSystems.Message
		// 	result.Data = ucp
		// 	return &result, nil
		// } else if ucp.Name != name && ucp.SerialNumber == serialNumber {
		// 	e2 := fmt.Errorf("storage with Serial number %s on-boarded to different system named %s", serialNumber, ucp.Name)
		// 	return nil, e2
		// }
		if ucp.Name == name {
			for _, v := range ucp.StorageDevices {
				if v.SerialNumber == serialNumber {
					result.Path = ucpSystems.Path
					result.Message = ucpSystems.Message
					result.Data = ucp
					return &result, nil
				}

			}
		}
	}

	e2 := fmt.Errorf("storage serial number is not on-boarded %s", serialNumber)
	return nil, e2

}

// createUcpSystem .
func (psm *infraGwManager) GetOrCreateDefaultUcpSystem(reqBody *model.CreateStorageDeviceParam) (*model.UcpSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	// First check if the ucp system we are about to create already exists
	if reqBody.GatewayAddress == "" {
		reqBody.GatewayAddress = objStorage.Address
	}
	found, ucp, err := psm.FindUcpSystemByName(model.DefaultSystemName)
	if err != nil {
		log.WriteDebug("TFError| error in FindUcpSystemByName call, err: %v", err)
		return nil, err
	}

	if !(*found) {
		// ucp system does not exist, create first

		body := model.CreateUcpSystemParam{
			Name:           model.DefaultSystemName,
			SerialNumber:   model.DefaultSystemSerialNumber,
			GatewayAddress: reqBody.GatewayAddress,
			Model:          "Logical UCP",
			Region:         "AMERICA",
			Country:        "United States",
			Zipcode:        "95054",
		}
		id, err := provObj.CreateUcpSystem(body)
		if err != nil {
			e2 := fmt.Errorf("failed to create Ucp System, error code: %v", err)
			log.WriteDebug("TFError| error in CreateUcpSystem call, err: %v", err)
			return nil, e2
		}

		return psm.GetUcpSystemById(*id)
	} else {
		return ucp, nil
	}

}

func (psm *infraGwManager) ReconcileStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ucpSystem := createInput.UcpSystem

	if storageId == "" {
		// The storage is not present in the ucp systems, we have to add the storage to a ucp system

		if ucpSystem == "" {
			// The user did not provide any ucp_system information, so we will create one and onboard the storage to that system
			reconcilerUcpSystem, err := psm.GetOrCreateDefaultUcpSystem(createInput)
			if err != nil {
				log.WriteDebug("TFError| error in createUcpSystem call, err: %v", err)
				return nil, err
			}
			createInput.UcpSystem = reconcilerUcpSystem.Data.Name
		} else {
			_, existSystem, err := psm.FindUcpSystemByName(ucpSystem)
			if err != nil {
				log.WriteDebug("TFError| error in FindUcpSystemByName call, err: %v", err)
				return nil, err
			}
			createInput.UcpSystem = existSystem.Data.Name
		}
		reconcilerSd, err := psm.addStorageDevice(createInput)
		if err != nil {
			log.WriteDebug("TFError| error in addStorageDevice call, err: %v", err)
			return nil, err
		}

		return reconcilerSd, nil

	} else {
		// The storage id is present, so this is an update
		reconcilerSd, err := psm.updateStorageDevice(storageId, createInput)
		if err != nil {
			log.WriteDebug("TFError| error in updateStorageDevice call, err: %v", err)
			return nil, err
		}
		return reconcilerSd, nil

	}
}

// func (psm *infraGwManager) ReconcileStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error) {
// 	log := commonlog.GetLogger()
// 	log.WriteEnter()
// 	defer log.WriteExit()

// 	ucpSystem := createInput.UcpSystem

// 	if storageId == "" {
// 		// The storage is not present in the ucp systems, we have to add the storage to a ucp system

// 		if ucpSystem == "" {
// 			// The user did not provide any ucp_system information, so we will create one and onboard the storage to that system
// 			reconcilerUcpSystem, err := psm.GetOrCreateDefaultUcpSystem(createInput)
// 			if err != nil {
// 				log.WriteDebug("TFError| error in createUcpSystem call, err: %v", err)
// 				return nil, err
// 			}
// 			createInput.UcpSystem = reconcilerUcpSystem.Data.Name
// 		} else {
// 			_, existSystem, err := psm.FindUcpSystemByName(ucpSystem)
// 			if err != nil {
// 				log.WriteDebug("TFError| error in FindUcpSystemByName call, err: %v", err)
// 				return nil, err
// 			}
// 			createInput.UcpSystem = existSystem.Data.Name
// 		}
// 		reconcilerSd, err := psm.addStorageDevice(createInput)
// 		if err != nil {
// 			log.WriteDebug("TFError| error in addStorageDevice call, err: %v", err)
// 			return nil, err
// 		}

// 		return reconcilerSd, nil

// 	} else {
// 		// The storage id is present, so this is an update
// 		reconcilerSd, err := psm.updateStorageDevice(storageId, createInput)
// 		if err != nil {
// 			log.WriteDebug("TFError| error in updateStorageDevice call, err: %v", err)
// 			return nil, err
// 		}
// 		return reconcilerSd, nil

// 	}
// }

func (psm *infraGwManager) DeleteStorageDevice(storageId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	provSetting := model.InfraGwSettings(psm.setting)

	provObj, err := provisonerimpl.NewEx(provSetting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	sd, err := psm.GetStorageDevice(storageId)
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageDevice call, err: %v", err)
		return err
	}

	if len(sd.Data.UcpSystems) > 0 {
		// The storage is part of a UCP system, first remove it from the ucp system
		ucpId := sd.Data.UcpSystems[0]
		err = provObj.DeleteStorageDeviceFromUcp(ucpId, storageId)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteStorageDeviceFromUcp call, err: %v", err)
			return err
		}

	}

	// Now remove the storage from UCP inventory
	err = provObj.DeleteStorageDevice(storageId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteVolume call, err: %v", err)
		return err
	}

	return nil
}
