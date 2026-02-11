package sanstorage

import (
	"fmt"
	"reflect"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/san/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/reconciler/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/jinzhu/copier"
)

// ReconcileIscsiTarget will reconcile and call Create/Update iScsi
func (psm *sanStorageManager) ReconcileIscsiTarget(createInput *sanmodel.CreateIscsiTargetReq) (*sanmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var reconcilerIscsi *sanmodel.IscsiTarget = &sanmodel.IscsiTarget{}
	// 1) If IscsiTarget Exisit - Update 2) IscsiTarget Not Exist - Create New
	if createInput.IscsiTargetNumber != nil {
		// Get IscsiTarget
		iscsiTarget, err := psm.GetIscsiTarget(createInput.PortID, *createInput.IscsiTargetNumber)
		if err != nil {
			log.WriteDebug("TFError| error in GetIscsiTarget provisioner call, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), createInput.PortID, *createInput.IscsiTargetNumber)
			return nil, err
		}
		// If IscsiTarget not exist we will get "-" from Rest API
		if iscsiTarget.IscsiTargetName == "-" {
			reconcilerIscsi, err = psm.CreateIscsiTarget(createInput)
			if err != nil {
				log.WriteDebug("TFError| error in CreateIscsiTarget call, err: %v", err)
				return reconcilerIscsi, err
			}

		} else {
			// IscsiTarget already exist
			reconcilerIscsi, err = psm.updateIscsiTarget(iscsiTarget, createInput)
			if err != nil {
				log.WriteDebug("TFError| error in updateIscsiTarget call, err: %v", err)
				return reconcilerIscsi, err
			}
		}
	} else {
		// IscsiTarget number not given so new Iscsi will be create
		var err error = nil
		reconcilerIscsi, err = psm.CreateIscsiTarget(createInput)
		if err != nil {
			log.WriteDebug("TFError| error in CreateIscsiTarget call, err: %v", err)
			return reconcilerIscsi, err
		}
	}
	return reconcilerIscsi, nil
}

// CreateIscsiTarget will create the Iscsi Target
func (psm *sanStorageManager) CreateIscsiTarget(createInput *sanmodel.CreateIscsiTargetReq) (*sanmodel.IscsiTarget, error) {
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

	provCreateInput := provisonermodel.CreateIscsiTargetReq{}
	err = copier.Copy(&provCreateInput, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_BEGIN), createInput.PortID, createInput.IscsiTargetNumber)
	provIscsiTarget, err := provObj.CreateIscsiTarget(provCreateInput)
	if err != nil {
		log.WriteDebug("TFError| error in CreateIscsiTarget call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_FAILED), createInput.PortID, createInput.IscsiTargetNumber)
		return nil, err
	}

	// Converting  Reconciler to Provisioner
	reconcilerIscsiTarget := sanmodel.IscsiTarget{}
	err = copier.Copy(&reconcilerIscsiTarget, provIscsiTarget)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_END), createInput.PortID, createInput.IscsiTargetNumber)
	return &reconcilerIscsiTarget, nil
}

// updateIscsiTarget will update existing Iscsi Target
func (psm *sanStorageManager) updateIscsiTarget(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) (*sanmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// ISCSI TARGET HOSTGROUP MODE
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_ISCSITARGET_BEGIN), createInput.PortID, *createInput.IscsiTargetNumber)
	err := psm.reconcileIscsiTargetHostModeAndHostModeOptions(existingIscsiTarget, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileIscsiTargetHostModeAndHostModeOptions, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_MODE_OPTION_ISCSITARGET_FAILED), createInput.PortID, *createInput.IscsiTargetNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_ISCSITARGET_END), createInput.PortID, *createInput.IscsiTargetNumber)

	// ISCSI TARGET  INITIATOR
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_INITIATOR_ISCSITARGET_BEGIN), createInput.PortID, *createInput.IscsiTargetNumber)
	err = psm.reconcileIscsiInitiators(existingIscsiTarget, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileIscsiInitiators, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_INITIATOR_ISCSITARGET_FAILED), createInput.PortID, *createInput.IscsiTargetNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_INITIATOR_ISCSITARGET_END), createInput.PortID, *createInput.IscsiTargetNumber)

	// ISCSI TARGET LUN
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_LUN_ISCSITARGET_BEGIN), createInput.PortID, *createInput.IscsiTargetNumber)
	err = psm.reconcileLunPaths(existingIscsiTarget, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in reconcileLunPaths, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_LUN_ISCSITARGET_FAILED), createInput.PortID, *createInput.IscsiTargetNumber)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_LUN_ISCSITARGET_END), createInput.PortID, *createInput.IscsiTargetNumber)

	// After Update - Get IscsiTarget information
	hostGroup, err := psm.GetIscsiTarget(createInput.PortID, *createInput.IscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetIscsiTarget provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), createInput.PortID, *createInput.IscsiTargetNumber)
		return nil, err
	}
	return hostGroup, nil
}

// reconcileIscsiTargetHostModeAndHostModeOptions will update hostmode and hostmode option for existing Iscsi Target
func (psm *sanStorageManager) reconcileIscsiTargetHostModeAndHostModeOptions(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	finalHostMode := ""
	finalHostModeOptions := createInput.HostModeOptions
	// If it is same then no need to update
	if createInput.HostMode != nil {
		if (existingIscsiTarget.HostMode == *createInput.HostMode) && (reflect.DeepEqual(existingIscsiTarget.HostModeOptions, *createInput.HostModeOptions)) {
			log.WriteDebug("TFDebug| No Need to Update - HostMode and Options are same")
			return nil
		}
	}
	// If Hostmode not given in TF file but existing in available hostgroup then need to set default = LINUX/IRIX
	if (existingIscsiTarget.HostMode != "") && (createInput.HostMode == nil) {
		finalHostMode = "LINUX/IRIX"
	} else {
		finalHostMode = *createInput.HostMode
	}
	// If HostMode Option is removed from TF file but available in existing Hostgroup then need to remove/reset hostgroup mode
	if (len(existingIscsiTarget.HostModeOptions) > 0) && (createInput.HostModeOptions == nil) {
		finalHostModeOptions = &([]int{-1})
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
	provisionerHostMode := provisonermodel.SetIscsiHostModeAndOptions{
		HostMode:        finalHostMode,
		HostModeOptions: finalHostModeOptions,
	}
	err = provObj.SetIscsiHostGroupModeAndOptions(createInput.PortID, *createInput.IscsiTargetNumber, provisionerHostMode)
	if err != nil {
		log.WriteDebug("TFError| error in SetIscsiHostGroupModeAndOptions  call, err: %v", err)
		return err
	}

	return nil
}

// reconcileIscsiInitiators will update hostmode and hostmode option for existing Iscsi Target
func (psm *sanStorageManager) reconcileIscsiInitiators(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) error {
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
	if psm.isIscsiNickNameAlreadyExist(existingIscsiTarget, createInput) {
		msg := "iSCSI nickname is already used in the same port."
		return fmt.Errorf(msg)
	}
	AddInitiator, UpdateInitiator, DeleteInitiator := psm.getUpdatedIscsiInitiatorInformation(existingIscsiTarget, createInput)

	// Delete Initiator
	for _, delete := range DeleteInitiator {
		err := provObj.DeleteIscsiNameFromIscsiTarget(createInput.PortID, *createInput.IscsiTargetNumber, delete.IscsiTargetNameIqn)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteIscsiNameFromIscsiTarget call, err: %v", err)
			return err
		}
	}
	// Update Initiator
	for _, update := range UpdateInitiator {
		req := provisonermodel.SetNicknameIscsiReq{
			IscsiNickname: update.IscsiNickname,
		}
		err := provObj.SetNicknameForIscsiName(createInput.PortID, *createInput.IscsiTargetNumber, update.IscsiTargetNameIqn, req)
		if err != nil {
			log.WriteDebug("TFError| error in SetNicknameForIscsiName call, err: %v", err)
			return err
		}
	}

	// Add Initiator
	for _, add := range AddInitiator {
		// Getting Provisioner host mode structure
		if add.IscsiTargetNameIqn != "" {
			provisionerIqn := provisonermodel.SetIscsiNameReq{
				PortID:             createInput.PortID,
				IscsiTargetNameIqn: add.IscsiTargetNameIqn,
				IscsiTargetNumber:  *createInput.IscsiTargetNumber,
			}
			err := provObj.SetIscsiNameForIscsiTarget(provisionerIqn)
			if err != nil {
				log.WriteDebug("TFError| error in SetIscsiNameForIscsiTarget call, err: %v", err)
				return err
			}
		}
		if add.IscsiNickname != "" {
			req := provisonermodel.SetNicknameIscsiReq{
				IscsiNickname: add.IscsiNickname,
			}
			err := provObj.SetNicknameForIscsiName(createInput.PortID, *createInput.IscsiTargetNumber, add.IscsiTargetNameIqn, req)
			if err != nil {
				log.WriteDebug("TFError| error in SetNicknameForIscsiName call, err: %v", err)
				return err
			}
		}
	}
	return nil
}

// getUpdatedIscsiInitiatorInformation collecte Update,  Add and deleted Wwn array
func (psm *sanStorageManager) getUpdatedIscsiInitiatorInformation(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) (AddInitiator, UpdateInitiator, DeleteInitiator []sanmodel.Initiator) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	/*
	 1) Find NOT available in existingIscsiTarget.IscsiTargetNameIqn and available in createInput.IscsiTargetNameIqn - ADD New
	 2) Find available in existingIscsiTarget.IscsiTargetNameIqn and NOT available in createInput.IscsiTargetNameIqn - DELETE
	 3) Find IscsiTargetNameIqn matching in existingIscsiTarget.IscsiTargetNameIqn and createInput.IscsiTargetNameIqn but nickname not matching - Update NickName
	*/
	// Find updated
	for _, current := range createInput.Initiators {
		for _, existing := range existingIscsiTarget.Initiators {
			if (current.IscsiTargetNameIqn == existing.IscsiTargetNameIqn) && (current.IscsiNickname != existing.IscsiNickname) {
				UpdateInitiator = append(UpdateInitiator, current)
			}
		}
	}

	// Find New wwn
	isNewInitiator := true
	for _, current := range createInput.Initiators {
		for _, existing := range existingIscsiTarget.Initiators {
			if current.IscsiTargetNameIqn == existing.IscsiTargetNameIqn {
				isNewInitiator = false
				break
			}
		}
		if isNewInitiator {
			AddInitiator = append(AddInitiator, current)
		}
		isNewInitiator = true
	}

	// Find Deleted Initiator
	isDeletedInitiator := true
	for _, existing := range existingIscsiTarget.Initiators {
		for _, current := range createInput.Initiators {
			if existing.IscsiTargetNameIqn == current.IscsiTargetNameIqn {
				isDeletedInitiator = false
				break
			}
		}
		if isDeletedInitiator {
			initiatorinfo := sanmodel.Initiator{IscsiTargetNameIqn: existing.IscsiTargetNameIqn, IscsiNickname: existing.IscsiNickname}
			DeleteInitiator = append(DeleteInitiator, initiatorinfo)
		}
		isDeletedInitiator = true
	}

	return AddInitiator, UpdateInitiator, DeleteInitiator
}

// isIscsiNickNameAlreadyExist check if nickname already exist
func (psm *sanStorageManager) isIscsiNickNameAlreadyExist(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) bool {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Check if same wwn nickname available in TF file
	iCount := 0
	for _, c1wwn := range createInput.Initiators {
		for _, c2wwn := range createInput.Initiators {
			if c1wwn.IscsiNickname == c2wwn.IscsiNickname {
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

// reconcileLunPaths will update Lun Paths
func (psm *sanStorageManager) reconcileLunPaths(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if psm.isLunIdExistInIscsiTarget(existingIscsiTarget, createInput) {
		// If Lun id already present Puma gives Error: An error occurred in the storage system. (message = Another LDEV is already mapped to LUN.)
		msg := "Another LDEV is already mapped to LUN."
		return fmt.Errorf(msg)
	}
	// Get Added and Delete LU array
	AddLU, DeleteLU := psm.getIscsiUpdateLUPathInformation(existingIscsiTarget, createInput)

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
			err := provObj.RemoveLdevFromHG(createInput.PortID, *createInput.IscsiTargetNumber, *delete.Lun)
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
			PortID:          &createInput.PortID,
			HostGroupNumber: createInput.IscsiTargetNumber,
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

// isLunIdExistInIscsiTarget is used to check Lun Id exist in Iscsi target
func (psm *sanStorageManager) isLunIdExistInIscsiTarget(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) bool {

	// Check TF file Lun Id is already exist in HostGroup
	for _, cLdev := range createInput.Ldevs {
		for _, eLdev := range existingIscsiTarget.LuPaths {
			if cLdev.Lun != nil {
				if (*cLdev.Lun == eLdev.Lun) && (*cLdev.LdevId != eLdev.LdevID) {
					return true
				}
			}
		}
	}
	return false
}

// getIscsiUpdateLUPathInformation used to get Added and Deleted LU Path information
func (psm *sanStorageManager) getIscsiUpdateLUPathInformation(existingIscsiTarget *sanmodel.IscsiTarget, createInput *sanmodel.CreateIscsiTargetReq) (AddLu, DeleteLu []sanmodel.Luns) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	isNewLu := true
	for _, clu := range createInput.Ldevs {
		for _, elu := range existingIscsiTarget.LuPaths {
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
	for _, elu := range existingIscsiTarget.LuPaths {
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

// GetIscsiTarget
func (psm *sanStorageManager) GetIscsiTarget(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_BEGIN), portID, iscsiTargetNumber)
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
	provIscsiTarget, err := provObj.GetIscsiTarget(portID, iscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| error in GetIscsiTarget provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), portID, iscsiTargetNumber)
		return nil, err
	}
	// Converting Prov to Reconciler
	reconIscsiTarget := sanmodel.IscsiTarget{}
	err = copier.Copy(&reconIscsiTarget, provIscsiTarget)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_END), portID, iscsiTargetNumber)
	return &reconIscsiTarget, nil
}

// GetIscsiTargetsByPortIds
func (psm *sanStorageManager) GetIscsiTargetsByPortIds(portIds []string) (*sanmodel.IscsiTargets, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), objStorage.Serial)

	// get iscsi ports
	provStoragePorts, err := provObj.GetStoragePorts()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), objStorage.Serial)
		return nil, err
	}

	// filter only iscsi ports
	iscsiPorts := []string{}
	for _, p := range *provStoragePorts {
		if p.PortType == "ISCSI" {
			iscsiPorts = append(iscsiPorts, p.PortId)
		}
	}

	log.WriteDebug("All ISCSI ports: %+v", iscsiPorts)

	portIdsIscsi := []string{}
	if len(portIds) > 0 {
		// check if iscsi port
		for _, p := range portIds {
			if containsStringIgnoreCase(iscsiPorts, p) {
				portIdsIscsi = append(portIdsIscsi, p)
			} else {
				log.WriteDebug("TFError| portId %v is not an ISCSI port", p)
			}
		}
	} else {
		portIdsIscsi = iscsiPorts
	}

	log.WriteDebug("Get Info for ISCSI ports: %+v", portIdsIscsi)

	provIscsiTargets := []provisonermodel.IscsiTargetGwy{}
	for _, p := range portIdsIscsi {
		log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), objStorage.Serial)
		pIscsis, err := provObj.GetIscsiTargetsByPortId(p)
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_ISCSITARGET_FAILED), objStorage.Serial)
			return nil, err
		}
		provIscsiTargets = append(provIscsiTargets, pIscsis.IscsiTargets...)
	}

	// Converting Prov to Reconciler
	reconIscsiTargets := []sanmodel.IscsiTarget{}
	err = copier.Copy(&reconIscsiTargets, provIscsiTargets)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_END), objStorage.Serial)

	return &sanmodel.IscsiTargets{
		IscsiTargets: reconIscsiTargets,
	}, nil
}

// GetAllIscsiTargets
func (psm *sanStorageManager) GetAllIscsiTargets() (*sanmodel.IscsiTargets, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), objStorage.Serial)
	provIscsiTargets, err := provObj.GetAllIscsiTargets()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllIscsiTargets provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_ISCSITARGET_FAILED), objStorage.Serial)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconIscsiTargets := sanmodel.IscsiTargets{}
	err = copier.Copy(&reconIscsiTargets, provIscsiTargets)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_END), objStorage.Serial)

	return &reconIscsiTargets, nil
}

// DeleteIscsiTarget
func (psm *sanStorageManager) DeleteIscsiTarget(portID string, iscsiTargetNumber int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_BEGIN), portID, iscsiTargetNumber)
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
	err = provObj.DeleteIscsiTarget(portID, iscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteIscsiTarget provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_FAILED), portID, iscsiTargetNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_END), portID, iscsiTargetNumber)

	return nil
}

func containsStringIgnoreCase(slice []string, str string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}
