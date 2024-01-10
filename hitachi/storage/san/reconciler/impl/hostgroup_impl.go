package sanstorage

import (
	"fmt"
	"reflect"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetHostGroup Get host information
func (psm *sanStorageManager) GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), portID, hostGroupNumber)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	provHostGroup, err := provObj.GetHostGroup(portID, hostGroupNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), portID, hostGroupNumber)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconcilerHostGroup := sanmodel.HostGroup{}
	err = copier.Copy(&reconcilerHostGroup, provHostGroup)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), portID, hostGroupNumber)
	return &reconcilerHostGroup, nil
}

// GetAllHostGroup Gets all host information
func (psm *sanStorageManager) GetAllHostGroups() (*sanmodel.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_BEGIN), objStorage.Serial)
	provHostGroups, err := provObj.GetAllHostGroups()
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconcilerHostGroups := sanmodel.HostGroups{}
	err = copier.Copy(&reconcilerHostGroups, provHostGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_END), objStorage.Serial)

	return &reconcilerHostGroups, nil
}

// GetHostGroupsByPortIds Gets all host information
func (psm *sanStorageManager) GetHostGroupsByPortIds(portIds []string) (*sanmodel.HostGroups, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_BEGIN), objStorage.Serial)
	provHostGroups, err := provObj.GetAllHostGroups()
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), objStorage.Serial)
		return nil, err
	}
	provHostGroupsFilter := provisonermodel.HostGroups{}

	for _, host := range provHostGroups.HostGroups {
		for _, id := range portIds {
			if host.PortID == id {
				provHostGroupsFilter.HostGroups = append(provHostGroupsFilter.HostGroups, host)
			}
		}
	}
	// Converting Prov to Reconciler
	reconcilerHostGroups := sanmodel.HostGroups{}
	err = copier.Copy(&reconcilerHostGroups, provHostGroupsFilter)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_END), objStorage.Serial)

	return &reconcilerHostGroups, nil
}

// ReconcileHostGroup will reconcile and call Create/Update hostgroup
func (psm *sanStorageManager) ReconcileHostGroup(createInput *sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var reconcilerHg *sanmodel.HostGroup = &sanmodel.HostGroup{}
	// 1) If Hostroup Exisit - Update 2) Hostgroup Not Exist - Create New
	if createInput.HostGroupNumber != nil {
		// Get Hostgroup
		hostGroup, err := psm.GetHostGroup(*createInput.PortID, *createInput.HostGroupNumber)
		if err != nil {
			log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), *createInput.PortID, *createInput.HostGroupNumber)
			return nil, err
		}
		// If Hostgroup not exist we will get "-" from Rest API
		if hostGroup.HostGroupName == "-" {
			// Hostgroup not exist - create new
			reconcilerHg, err = psm.createHostGroup(createInput)
			if err != nil {
				log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
				return reconcilerHg, err
			}
		} else {
			// Hostgroup already exist
			reconcilerHg, err = psm.updateHostgroup(hostGroup, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in updateHostgroup call, err: %v", err)
				return reconcilerHg, err
			}
		}
	} else {
		// Hostgroup number not given so new hostgroup will be create
		var err error = nil
		reconcilerHg, err = psm.createHostGroup(createInput)
		if err != nil {
			log.WriteDebug("TFError| error in createHostGroup call, err: %v", err)
			return reconcilerHg, err
		}
	}

	return reconcilerHg, nil
}

// createHostGroup will create new hostgroup
func (psm *sanStorageManager) createHostGroup(createInput *sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_BEGIN), createInput.PortID, createInput.HostGroupNumber)
	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}
	// Converting  Reconciler to Provisioner
	provisionerCreate := provisonermodel.CreateHostGroupRequest{}
	err = copier.Copy(&provisionerCreate, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to provisioner structure, err: %v", err)
		return nil, err
	}
	provHostGroup, err := provObj.CreateHostGroup(provisionerCreate)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroup  call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_HOSTGROUP_FAILED), createInput.PortID, createInput.HostGroupNumber)
		return nil, err
	}
	// Converting  Reconciler to Provisioner
	reconcilerHg := sanmodel.HostGroup{}
	err = copier.Copy(&reconcilerHg, provHostGroup)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_END), createInput.PortID, createInput.HostGroupNumber)
	return &reconcilerHg, nil
}

// updateHostgroup will update existing hostgroup
func (psm *sanStorageManager) updateHostgroup(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// HOSTGROUP MODE
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_HOSTGROUP_BEGIN), createInput.PortID, createInput.HostGroupNumber)
	err := psm.reconcileHostGroupMode(existingHostgroup, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileHostGroupMode, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_MODE_OPTION_HOSTGROUP_FAILED), createInput.PortID, createInput.HostGroupNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_HOSTGROUP_END), createInput.PortID, createInput.HostGroupNumber)

	// HOSTGROUP WWN
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_WWN_HOSTGROUP_BEGIN), createInput.PortID, createInput.HostGroupNumber)
	err = psm.reconcileHostGroupWwns(existingHostgroup, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileHostGroupWwns, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_WWN_HOSTGROUP_FAILED), createInput.PortID, createInput.HostGroupNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_WWN_HOSTGROUP_END), createInput.PortID, createInput.HostGroupNumber)

	// HOSTGROUP LUN
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_LUN_HOSTGROUP_BEGIN), createInput.PortID, createInput.HostGroupNumber)
	err = psm.reconcileHostGroupLuns(existingHostgroup, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileHostGroupWwns, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_LUN_HOSTGROUP_FAILED), createInput.PortID, createInput.HostGroupNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_LUN_HOSTGROUP_END), createInput.PortID, createInput.HostGroupNumber)

	// After Update - Get Hostgroup
	hostGroup, err := psm.GetHostGroup(*createInput.PortID, *createInput.HostGroupNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), *createInput.PortID, *createInput.HostGroupNumber)
		return nil, err
	}
	return hostGroup, nil
}

// reconcileHostGroupMode will update hostmode and hostmode option for existing hostgroup
func (psm *sanStorageManager) reconcileHostGroupMode(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	finalHostMode := ""
	finalHostModeOptions := createInput.HostModeOptions
	// If it is same then no need to update
	if createInput.HostMode != nil {
		if (existingHostgroup.HostMode == *createInput.HostMode) && (reflect.DeepEqual(existingHostgroup.HostModeOptions, createInput.HostModeOptions)) {
			log.WriteDebug("TFDebug| No Need to Update - HostMode and Options are same")
			return nil
		}
	}
	// If Hostmode not given in TF file but existing in available hostgroup then need to set default = LINUX/IRIX
	if (existingHostgroup.HostMode != "") && (createInput.HostMode == nil) {
		finalHostMode = "LINUX/IRIX"
	} else {
		finalHostMode = *createInput.HostMode
	}
	// If HostMode Option is removed from TF file but available in existing Hostgroup then need to remove/reset hostgroup mode
	if (len(existingHostgroup.HostModeOptions) > 0) && (createInput.HostModeOptions == nil) {
		finalHostModeOptions = []int{-1}
	}

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	// Getting Provisioner host mode structure
	provisionerHostMode := provisonermodel.SetHostModeAndOptions{
		HostMode:        finalHostMode,
		HostModeOptions: &finalHostModeOptions,
	}
	err = provObj.SetHostGroupModeAndOptions(*createInput.PortID, *createInput.HostGroupNumber, provisionerHostMode)
	if err != nil {
		log.WriteDebug("TFError| error in SetHostGroupModeAndOptions  call, err: %v", err)
		return err
	}
	return nil
}

// reconcileHostGroupWwns will update WWN details for existing hostgroup
func (psm *sanStorageManager) reconcileHostGroupWwns(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}
	if psm.isNickNameAlreadyExist(existingHostgroup, createInput) {
		//"An error occurred in the storage system. (message = WWN nickname is already used in the same port.)"
		msg := "WWN nickname is already used in the same port."
		return fmt.Errorf(msg)
	}
	AddWwn, UpdateWwn, DeleteWwn := psm.getUpdatedWwnInformation(existingHostgroup, createInput)

	// Delete WWN
	for _, delete := range DeleteWwn {
		err := provObj.DeleteWwn(*createInput.PortID, *createInput.HostGroupNumber, delete.Wwn)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteWwn call, err: %v", err)
			return err
		}
	}
	// Update WWN
	for _, update := range UpdateWwn {
		err := provObj.SetHostWwnNickName(*createInput.PortID, *createInput.HostGroupNumber, update.Wwn, update.Name)
		if err != nil {
			log.WriteDebug("TFError| error in SetHostWwnNickName call, err: %v", err)
			return err
		}
	}
	// Add WWN
	for _, add := range AddWwn {
		// Getting Provisioner host mode structure
		provisionerWwn := provisonermodel.AddWwnToHg{
			HostWwn:         &add.Wwn,
			PortID:          createInput.PortID,
			HostGroupNumber: createInput.HostGroupNumber,
		}
		if add.Wwn != "" {
			err := provObj.AddWwnToHG(provisionerWwn)
			if err != nil {
				log.WriteDebug("TFError| error in AddWwnToHG call, err: %v", err)
				return err
			}
		}
		if add.Name != "" {
			err := provObj.SetHostWwnNickName(*createInput.PortID, *createInput.HostGroupNumber, add.Wwn, add.Name)
			if err != nil {
				log.WriteDebug("TFError| error in SetHostWwnNickName call, err: %v", err)
				return err
			}
		}
	}

	return nil
}

// isNickNameAlreadyExist check if nickname already exist
func (psm *sanStorageManager) isNickNameAlreadyExist(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) bool {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Check if same wwn nickname available in TF file
	iCount := 0
	for _, c1wwn := range createInput.Wwns {
		for _, c2wwn := range createInput.Wwns {
			if c1wwn.Name == c2wwn.Name {
				iCount++
			}
		}
		if iCount > 1 {
			return true
		}
		iCount = 0
	}

	return false
}

// getUpdatedWwnInformation collecte Update,  Add and deleted Wwn array
func (psm *sanStorageManager) getUpdatedWwnInformation(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) (AddWwn, UpdateWwn, DeleteWwn []sanmodel.Wwn) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	/*
	 1) Find not available in existingHostgroup.wwn and available in createInput.wwn - ADD New
	 2) Find available in existingHostgroup.wwn and not available in createInput.wwn - DELETE
	 3) Find WWN matching in existingHostgroup.wwn and createInput.wwn but nickname not matching - Update NickName
	*/
	// Find updated
	for _, cwwn := range createInput.Wwns {
		for _, ewwn := range existingHostgroup.WwnDetails {
			if (cwwn.Wwn == ewwn.Wwn) && (cwwn.Name != ewwn.Name) {
				UpdateWwn = append(UpdateWwn, cwwn)
			}
		}
	}

	// Find New wwn
	isNewWwn := true
	for _, cwwn := range createInput.Wwns {
		for _, ewwn := range existingHostgroup.WwnDetails {
			if cwwn.Wwn == ewwn.Wwn {
				isNewWwn = false
				break
			}
		}
		if isNewWwn {
			AddWwn = append(AddWwn, cwwn)
		}
		isNewWwn = true
	}

	// Find Deleted wwn
	isDeletedWwn := true
	for _, ewwn := range existingHostgroup.WwnDetails {
		for _, cwwn := range createInput.Wwns {
			if ewwn.Wwn == cwwn.Wwn {
				isDeletedWwn = false
				break
			}
		}
		if isDeletedWwn {
			wwninfo := sanmodel.Wwn{Wwn: ewwn.Wwn, Name: ewwn.Name}
			DeleteWwn = append(DeleteWwn, wwninfo)
		}
		isDeletedWwn = true
	}

	return AddWwn, UpdateWwn, DeleteWwn
}

// reconcileHostGroupLuns will update Hostgroup Lun details for existing hostgroup
func (psm *sanStorageManager) reconcileHostGroupLuns(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if psm.isLunIdExistInHostGroup(existingHostgroup, createInput) {
		// If Lun id already present Puma gives Error: An error occurred in the storage system. (message = Another LDEV is already mapped to LUN.)
		msg := "Another LDEV is already mapped to LUN."
		return fmt.Errorf(msg)
	}
	// Get Added and Delete LU array
	AddLU, DeleteLU := psm.getUpdateLUPathInformation(existingHostgroup, createInput)

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}
	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	// DELETE LU PATH
	for _, delete := range DeleteLU {
		if delete.Lun != nil {
			err := provObj.RemoveLdevFromHG(*createInput.PortID, *createInput.HostGroupNumber, *delete.Lun)
			if err != nil {
				log.WriteDebug("TFError| error in RemoveLdevFromHG call, err: %v", err)
				return err
			}
		}
	}

	// ADD LU PATH
	for _, add := range AddLU {
		// Getting Provisioner host mode structure
		provisionerLu := provisonermodel.AddLdevToHg{
			PortID:          createInput.PortID,
			HostGroupNumber: createInput.HostGroupNumber,
			LdevID:          add.LdevId,
			Lun:             add.Lun,
		}
		err := provObj.AddLdevToHG(provisionerLu)
		if err != nil {
			log.WriteDebug("TFError| error in AddLdevToHG call, err: %v", err)
			return err
		}
	}

	return nil
}

// TODO : IN FUTURE USE IF REQUIRE - COMMENT CODE
/*
// isLdevIdExistInHostGroup is used to check Ldev Id exist in Hostgroup
func (psm *sanStorageManager) isLdevIdExistInHostGroup(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) bool {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	// Check TF file LdevID is already exist in HostGroup
	for _, cLdev := range createInput.Ldevs {
		for _, eLdev := range existingHostgroup.LuPaths {
			if cLdev.LdevId != nil {
				if cLdev.Lun == nil {
					if *cLdev.LdevId == eLdev.LdevID {
						return true
					}
				} else {
					// IF in TF file and Existing Hostgroup have not same data then only need to insert.
					if (*cLdev.LdevId == eLdev.LdevID) && (*cLdev.Lun != eLdev.Lun) {
						return true
					}
				}
			}
		}
	}
	return false
}
*/

// isLunIdExistInHostGroup is used to check Lun Id exist in hostgroup
func (psm *sanStorageManager) isLunIdExistInHostGroup(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) bool {
	// Check TF file Lun Id is already exist in HostGroup
	for _, cLdev := range createInput.Ldevs {
		for _, eLdev := range existingHostgroup.LuPaths {
			if cLdev.Lun != nil {
				if (*cLdev.Lun == eLdev.Lun) && (*cLdev.LdevId != eLdev.LdevID) {
					return true
				}
			}
		}
	}
	return false
}

// getUpdateLUPathInformation used to get Added and Deleted LU Path information
func (psm *sanStorageManager) getUpdateLUPathInformation(existingHostgroup *sanmodel.HostGroup, createInput *sanmodel.CreateHostGroupRequest) (AddLu, DeleteLu []sanmodel.Luns) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	isNewLu := true
	for _, clu := range createInput.Ldevs {
		for _, elu := range existingHostgroup.LuPaths {
			if (clu.LdevId != nil) && (clu.Lun != nil) {
				if (*clu.LdevId == elu.LdevID) && (*clu.Lun == elu.Lun) {
					isNewLu = false
					break
				}
			}
		}
		if isNewLu {
			if clu.Lun == nil {
				AddLu = append(AddLu, sanmodel.Luns{LdevId: clu.LdevId, Lun: nil})
			} else {
				AddLu = append(AddLu, sanmodel.Luns{LdevId: clu.LdevId, Lun: clu.Lun})
			}
		}
		isNewLu = true
	}

	// Find Deleted LU
	isDeletedLu := true
	for _, elu := range existingHostgroup.LuPaths {
		for _, clu := range createInput.Ldevs {
			if (clu.LdevId != nil) && (clu.Lun != nil) {
				if (elu.LdevID == *clu.LdevId) && (elu.Lun == *clu.Lun) {
					isDeletedLu = false
					break
				}
			}
		}
		if isDeletedLu {
			ldev := elu.LdevID
			lun := elu.Lun
			DeleteLu = append(DeleteLu, sanmodel.Luns{LdevId: &ldev, Lun: &lun})
		}
		isDeletedLu = true
	}

	return AddLu, DeleteLu
}

// DeleteHostGroup will delete the hostgroup
func (psm *sanStorageManager) DeleteHostGroup(portID string, hostGroupNumber int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_BEGIN), portID, hostGroupNumber)
	err = provObj.DeleteHostGroup(portID, hostGroupNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetHostGroup provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_HOSTGROUP_FAILED), portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_END), portID, hostGroupNumber)

	return nil
}
