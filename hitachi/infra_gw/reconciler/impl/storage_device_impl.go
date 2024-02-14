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

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStorageDevices()
}

// GetStorageDevice gets storage device information
func (psm *infraGwManager) GetStorageDevice(id string) (*model.StorageDevice, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	return provObj.GetStorageDevice(id)
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
	log.WriteDebug("TFError| Return value of AddStorageDevice : %v", storageId)

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

func (psm *infraGwManager) findUcpSystemBySerial(serial string) (*model.UcpSystem, error) {
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

// createUcpSystem .
func (psm *infraGwManager) createUcpSystem(reqBody *model.CreateStorageDeviceParam) (*model.UcpSystem, error) {
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

	ucpSerialNumber := "Logical-UCP-" + reqBody.SerialNumber

	ucp, err := psm.findUcpSystemBySerial(ucpSerialNumber)

	if err != nil {
		// ucp system does not exist, create first

		body := model.CreateUcpSystemParam{
			Name:           "ucp-system-" + reqBody.SerialNumber,
			SerialNumber:   "Logical-UCP-" + reqBody.SerialNumber,
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
			reconcilerUcpSystem, err := psm.createUcpSystem(createInput)
			if err != nil {
				log.WriteDebug("TFError| error in createUcpSystem call, err: %v", err)
				return nil, err
			}
			createInput.UcpSystem = reconcilerUcpSystem.Data.Name
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
