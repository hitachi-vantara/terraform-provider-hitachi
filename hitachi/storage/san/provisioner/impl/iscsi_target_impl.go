package sanstorage

import (
	"fmt"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	sangatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/provisioner/model"

	"github.com/jinzhu/copier"
)

func (psm *sanStorageManager) GetIscsiTarget(portID string, iscsiTargetNumber int) (*sanmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_BEGIN), portID, iscsiTargetNumber)
	////////// in parallel
	var wga sync.WaitGroup
	var outMap sync.Map
	var errMap sync.Map

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	n := 1
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		it, getErr := gatewayObj.GetIscsiTarget(portID, iscsiTargetNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetIscsiTarget, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provIscsiTarget := sanmodel.IscsiTargetGwy{}
		copyErr := copier.Copy(&provIscsiTarget, it)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}

		errMap.Store(n, nil)
		outMap.Store(n, provIscsiTarget)
		return
	}(n)

	n = 2
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		itNameInfo, getErr := gatewayObj.GetIscsiNameInformation(portID, iscsiTargetNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetIscsiNameInformation, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provIscsiNameInfo := []sanmodel.IscsiNameInformation{}
		copyErr := copier.Copy(&provIscsiNameInfo, itNameInfo)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}

		errMap.Store(n, nil)
		outMap.Store(n, provIscsiNameInfo)
		return
	}(n)

	n = 3
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		itLuPaths, getErr := gatewayObj.GetIscsiTargetGroupLuPaths(portID, iscsiTargetNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetIscsiTargetGroupLuPaths, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provIscsiTargetLupaths := []sanmodel.IscsiTargetLuPath{}
		copyErr := copier.Copy(&provIscsiTargetLupaths, itLuPaths)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}
		
		errMap.Store(n, nil)
		outMap.Store(n, provIscsiTargetLupaths)
		return
	}(n)

	log.WriteDebug("TFDebug|Waiting for GetIscsiTarget goroutines to end ...")
	wga.Wait()
	log.WriteDebug("TFDebug|GetIscsiTarget goroutines ended")

	// check for errors
	for nth := 1; nth <= 3; nth++ {
		serr, _ := errMap.Load(nth)
		if serr != nil {
			log.WriteDebug("TFError|serr: %+v", serr.(error))
			log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), portID, iscsiTargetNumber)
			return nil, serr.(error)
		}
	}

	ihg, ok := outMap.Load(1)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[1]")
	}
	iitNameInfo, ok := outMap.Load(2)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[2]")
	}
	itLuPaths, ok := outMap.Load(3)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[3]")
	}

	hg := ihg.(sanmodel.IscsiTargetGwy)
	itNames := iitNameInfo.([]sanmodel.IscsiNameInformation)
	luPaths := itLuPaths.([]sanmodel.IscsiTargetLuPath)

	phg := sanmodel.IscsiTarget{}
	phg.IscsiTargetID = hg.IscsiTargetID
	phg.PortID = hg.PortID
	phg.IscsiTargetNumber = hg.IscsiTargetNumber
	phg.IscsiTargetName = hg.IscsiTargetName
	phg.HostMode = hg.HostMode
	phg.HostModeOptions = hg.HostModeOptions
	phg.IscsiTargetNameIqn = hg.IscsiTargetNameIqn

	iscsiNameDetails := []sanmodel.Initiator{}

	for _, hw := range itNames {
		name := hw.IscsiNickname
		if hw.IscsiNickname == "-" {
			name = ""
		}
		wd := sanmodel.Initiator{
			IscsiTargetNameIqn: hw.IscsiTargetNameIqn,
			IscsiNickname:      name,
		}
		iscsiNameDetails = append(iscsiNameDetails, wd)
	}
	phg.Initiators = iscsiNameDetails

	iscsiLuPaths := []sanmodel.LuPath{}
	ldevs := []int{}
	itLuns := []int{}
	for _, lp := range luPaths {
		lp := sanmodel.LuPath{
			Lun:    lp.Lun,
			LdevID: lp.LdevID,
		}
		iscsiLuPaths = append(iscsiLuPaths, lp)
		ldevs = append(ldevs, lp.LdevID)
		itLuns = append(itLuns, lp.Lun)
	}
	phg.LuPaths = iscsiLuPaths
	phg.Ldevs = ldevs
	phg.ItLuns = itLuns

	log.WriteDebug("TFDebug| iscsiTarget structure=%+v", phg)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_END), portID, iscsiTargetNumber)
	return &phg, nil
}

// GetAllIscsiTargets .
func (psm *sanStorageManager) GetAllIscsiTargets() (*sanmodel.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), objStorage.Serial)
	iscsiTargets, err := gatewayObj.GetAllIscsiTargets()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllIscsiTargets err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_ISCSITARGET_FAILED), objStorage.Serial)
		return nil, err
	}

	provIscsiTargets := sanmodel.IscsiTargets{}
	err = copier.Copy(&provIscsiTargets, iscsiTargets)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_END), objStorage.Serial)

	return &provIscsiTargets, nil
}

// SetHostGroupModeAndOptions .
func (psm *sanStorageManager) SetIscsiHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetIscsiHostModeAndOptions) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_ISCSITARGET_BEGIN), portID, hostGroupNumber)
	modeRequest := sangatewaymodel.SetIscsiHostModeAndOptions{
		HostMode:        reqBody.HostMode,
		HostModeOptions: reqBody.HostModeOptions,
	}
	err = gatewayObj.SetIScsiTargetHostModeAndHostModeOptions(portID, hostGroupNumber, modeRequest)
	if err != nil {
		log.WriteDebug("TFError| failed to call SetIScsiTargetHostModeAndHostModeOptions err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_MODE_OPTION_ISCSITARGET_FAILED), portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_ISCSITARGET_END), portID, hostGroupNumber)

	return nil
}

// CreateIscsiTarget
func (psm *sanStorageManager) CreateIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) (*sanmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	reqIscsiTarget := sangatewaymodel.CreateIscsiTargetReq{
		PortID:             reqBody.PortID,
		IscsiTargetName:    reqBody.IscsiTargetName,
		IscsiTargetNumber:  reqBody.IscsiTargetNumber,
		IscsiTargetNameIqn: reqBody.IscsiTargetNameIqn,
		HostModeOptions:    reqBody.HostModeOptions,
		HostMode:           reqBody.HostMode,
	}

	if reqIscsiTarget.HostMode == nil {
		defaultHostMode := "LINUX/IRIX"
		reqIscsiTarget.HostMode = &defaultHostMode
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_BEGIN), reqBody.PortID, reqBody.IscsiTargetNumber)
	_, pItNum, err := gatewayObj.CreateIscsiTarget(reqIscsiTarget)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_FAILED), reqBody.PortID, reqBody.IscsiTargetNumber)
		log.WriteDebug("TFError| failed to call CreateIscsiTarget err: %+v", err)
		return nil, err
	}

	if reqBody.IscsiTargetNumber == nil {
		reqBody.IscsiTargetNumber = pItNum
	}

	if len(*reqBody.Ldevs) > 0 {
		err := psm.AddLdevToIscsiTarget(reqBody)
		if err != nil {
			log.WriteDebug("TFError| failed to call AddLdevToIscsiTarget err: %+v", err)
			return nil, err
		}
	}

	if len(*reqBody.Initiators) > 0 {
		err := psm.AddInitiatorsToIscsiTarget(reqBody)
		if err != nil {
			log.WriteDebug("TFError| failed to call AddInitiatorsToIscsiTarget err: %+v", err)
			return nil, err
		}
	}

	iscsiTarget, err := psm.GetIscsiTarget(reqBody.PortID, *reqBody.IscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetIscsiTarget err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), reqBody.PortID, *pItNum)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_END), reqBody.PortID, reqBody.IscsiTargetNumber)

	return iscsiTarget, nil
}

// AddLdevToIscsiTarget
func (psm *sanStorageManager) AddLdevToIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if len(*reqBody.Ldevs) > 0 {
		// check if ldev is already used, before create
		for _, ldev := range *reqBody.Ldevs {
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_LDEV_TO_ISCSITARGET_BEGIN), &ldev.LdevId)
			updReq := sangatewaymodel.AddLdevToHgReqGwy{
				PortID:          &reqBody.PortID,
				HostGroupNumber: reqBody.IscsiTargetNumber,
				LdevID:          ldev.LdevId,
				Lun:             ldev.Lun,
			}
			err := gatewayObj.AddLdevToHG(updReq)
			if err != nil {
				log.WriteDebug("TFError| failed to call AddLdevToHG err: %+v", err)
				log.WriteError(mc.GetMessage(mc.ERR_ADD_LDEV_TO_ISCSITARGET_FAILED), &ldev.LdevId)
				return err
			}
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_LDEV_TO_ISCSITARGET_END), &ldev.LdevId)
		}
	}
	return nil
}

// AddInitiatorsToIscsiTarget
func (psm *sanStorageManager) AddInitiatorsToIscsiTarget(reqBody sanmodel.CreateIscsiTargetReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	if len(*reqBody.Initiators) > 0 {
		for _, initiator := range *reqBody.Initiators {
			// First, set iscsi name to iscsitarget
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_IQN_NAME_TO_ISCSITARGET_BEGIN), initiator.IscsiTargetNameIqn)
			updReq := sangatewaymodel.SetIscsiNameReq{
				IscsiTargetNameIqn: initiator.IscsiTargetNameIqn,
				PortID:             reqBody.PortID,
				IscsiTargetNumber:  *reqBody.IscsiTargetNumber,
			}
			err := gatewayObj.SetIscsiNameForIscsiTarget(updReq)
			if err != nil {
				log.WriteDebug("TFError| failed to call SetIscsiNameForIscsiTarget err: %+v", err)
				log.WriteError(mc.GetMessage(mc.ERR_ADD_IQN_NAME_TO_ISCSITARGET_FAILED), initiator.IscsiTargetNameIqn)
				return err
			}
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_IQN_NAME_TO_ISCSITARGET_END), initiator.IscsiTargetNameIqn)

			// Second, set nickname to iscsi name
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NICKNAME_TO_ISCSITARGET_BEGIN), initiator.IscsiNickname, initiator.IscsiTargetNameIqn)
			nameReq := sangatewaymodel.SetNicknameIscsiReq{
				IscsiNickname: initiator.IscsiNickname,
			}
			err = gatewayObj.SetNicknameForIscsiName(reqBody.PortID, *reqBody.IscsiTargetNumber, initiator.IscsiTargetNameIqn, nameReq)
			if err != nil {
				log.WriteDebug("TFError| failed to call SetNicknameForIscsiName err: %+v", err)
				log.WriteError(mc.GetMessage(mc.ERR_ADD_NICKNAME_TO_ISCSITARGET_FAILED), initiator.IscsiNickname, initiator.IscsiTargetNameIqn)
				return err
			}
			log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NICKNAME_TO_ISCSITARGET_END), initiator.IscsiNickname, initiator.IscsiTargetNameIqn)

		}
	}
	return nil
}

// SetIscsiNameForIscsiTarget .
func (psm *sanStorageManager) SetIscsiNameForIscsiTarget(reqBody sanmodel.SetIscsiNameReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_IQN_NAME_TO_ISCSITARGET_BEGIN), reqBody.IscsiTargetNameIqn)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	req := sangatewaymodel.SetIscsiNameReq{
		PortID:             reqBody.PortID,
		IscsiTargetNameIqn: reqBody.IscsiTargetNameIqn,
		IscsiTargetNumber:  reqBody.IscsiTargetNumber,
	}
	err = gatewayObj.SetIscsiNameForIscsiTarget(req)
	if err != nil {
		log.WriteDebug("TFError| failed to call SetIscsiNameForIscsiTarget err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_IQN_NAME_TO_ISCSITARGET_FAILED), reqBody.IscsiTargetNameIqn)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_IQN_NAME_TO_ISCSITARGET_END), reqBody.IscsiTargetNameIqn)
	return nil
}

// SetNicknameForIscsiName .
func (psm *sanStorageManager) SetNicknameForIscsiName(portID string, iscsiTargetNumber int, iscsiName string, reqBody sanmodel.SetNicknameIscsiReq) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NICKNAME_TO_ISCSITARGET_BEGIN), reqBody.IscsiNickname, iscsiName)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	req := sangatewaymodel.SetNicknameIscsiReq{
		IscsiNickname: reqBody.IscsiNickname,
	}
	err = gatewayObj.SetNicknameForIscsiName(portID, iscsiTargetNumber, iscsiName, req)
	if err != nil {
		log.WriteDebug("TFError| failed to call SetNicknameForIscsiName err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_NICKNAME_TO_ISCSITARGET_FAILED), reqBody.IscsiNickname, iscsiName)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NICKNAME_TO_ISCSITARGET_END), reqBody.IscsiNickname, iscsiName)
	return nil
}

// DeleteIscsiNameFromIscsiTarget .
func (psm *sanStorageManager) DeleteIscsiNameFromIscsiTarget(portID string, iscsiTargetNumber int, iscsiName string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_IQN_NAME_BEGIN), iscsiName, portID, iscsiTargetNumber)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = gatewayObj.DeleteIscsiNameFromIscsiTarget(portID, iscsiTargetNumber, iscsiName)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteIscsiNameFromIscsiTarget err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_IQN_NAME_FAILED), iscsiName, portID, iscsiTargetNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_IQN_NAME_END), iscsiName, portID, iscsiTargetNumber)
	return nil
}

// DeleteIscsiTarget
func (psm *sanStorageManager) DeleteIscsiTarget(portID string, iscsiTargetNumber int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := sangatewaymodel.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_BEGIN), portID, iscsiTargetNumber)
	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = gatewayObj.DeleteIscsiTarget(portID, iscsiTargetNumber)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteIscsiTarget, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_FAILED), portID, iscsiTargetNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_END), portID, iscsiTargetNumber)

	return nil
}
