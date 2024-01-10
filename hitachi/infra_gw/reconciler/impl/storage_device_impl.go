package infra_gw

import (
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

func (psm *infraGwManager) addStorageDevice(storageId string, createInput *model.CreateStorageDeviceParam) (*model.StorageDevice, error) {
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
	t, err := provObj.AddStorageDevice(storageId, *createInput)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup call, err: %v", err)
		return nil, err
	}
	log.WriteDebug("TFError| Return value of AddStorageDevice : %v", t)

	return psm.GetStorageDevice(storageId)
	/*
		if t.Data.Status == "Success" {
			return psm.GetStorageDevice(storageId)
		} else {

			port := reqBody.Port
			hostGroupName := reqBody.HostGroupName
			return nil, fmt.Errorf("Create hostgroup on storage devide %s with port %s and hostgroup name %s failed", storageId, port, hostGroupName), false
		}
	*/
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
		} else {
			// The user provided the ucp_system information so we will try to add the storage to the provided ucp_system
		}
		reconcilerSd, err := psm.addStorageDevice(storageId, createInput)
		if err != nil {
			log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
			return reconcilerSd, err
		}
		return reconcilerSd, err
	} else {

		// The storage id is present, so this is an update
		/*
			updateRequest := model.UpdateStorageDeviceParam{
				Username:  createInput.Username,
				Password:  createInput.Password,
				OutOfBand: createInput.OutOfBand,
			}
		*/
		return nil, nil

	}

}
