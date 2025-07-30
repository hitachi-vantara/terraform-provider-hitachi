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

func (psm *sanStorageManager) GetHostGroup(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), portID, hostGroupNumber)
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
		hg, getErr := gatewayObj.GetHostGroup(portID, hostGroupNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetHostGroup, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provHostGroup := sanmodel.HostGroupGwy{}
		copyErr := copier.Copy(&provHostGroup, hg)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}

		errMap.Store(n, nil)
		outMap.Store(n, provHostGroup)
		return
	}(n)

	n = 2
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		hgwwns, getErr := gatewayObj.GetHostGroupWwns(portID, hostGroupNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetHostGroupWwns, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provHostGroupWwnDetail := []sanmodel.HostWwnDetail{}
		copyErr := copier.Copy(&provHostGroupWwnDetail, hgwwns)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy wwns, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}
		errMap.Store(n, nil)
		outMap.Store(n, provHostGroupWwnDetail)
		return
	}(n)

	n = 3
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		hglupaths, getErr := gatewayObj.GetHostGroupLuPaths(portID, hostGroupNumber)
		if getErr != nil {
			log.WriteDebug("TFError| error in GetHostGroupLuPaths, err: %v", getErr)
			errMap.Store(n, getErr)
			return
		}

		provHostGroupLupaths := []sanmodel.HostLuPath{}
		copyErr := copier.Copy(&provHostGroupLupaths, hglupaths)
		if copyErr != nil {
			log.WriteDebug("TFError| error in Copy LUNs, err: %v", copyErr)
			errMap.Store(n, copyErr)
			return
		}

		errMap.Store(n, nil)
		outMap.Store(n, provHostGroupLupaths)
		return
	}(n)

	log.WriteDebug("TFDebug|Waiting for GetHostGroups goroutines to end ...")
	wga.Wait()
	log.WriteDebug("TFDebug|GetHostGroups goroutines ended")

	// check for errors
	for nth := 1; nth <= 3; nth++ {
		serr, _ := errMap.Load(nth)
		if serr != nil {
			log.WriteDebug("TFError|serr: %+v", serr.(error))
			log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), portID, hostGroupNumber)
			return nil, serr.(error)
		}
	}

	ihg, ok := outMap.Load(1)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[1]")
	}
	ihgwwns, ok := outMap.Load(2)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[2]")
	}
	ihglupaths, ok := outMap.Load(3)
	if !ok {
		return nil, fmt.Errorf("cannot load sync map outMap[3]")
	}

	hg := ihg.(sanmodel.HostGroupGwy)
	hgwwns := ihgwwns.([]sanmodel.HostWwnDetail)
	hglupaths := ihglupaths.([]sanmodel.HostLuPath)

	//////////

	phg := sanmodel.HostGroup{}
	phg.PortID = hg.PortID
	phg.HostGroupNumber = hg.HostGroupNumber
	phg.HostGroupName = hg.HostGroupName
	phg.HostMode = hg.HostMode
	phg.HostModeOptions = hg.HostModeOptions

	wwnDetails := []sanmodel.WwnDetail{}
	wwns := []string{}
	// ? do we need to convert wwn to uppercase ?
	for _, hw := range hgwwns {
		name := hw.WwnNickname
		if hw.WwnNickname == "-" {
			name = ""
		}
		wd := sanmodel.WwnDetail{
			Wwn:  hw.HostWwn,
			Name: name,
		}
		wwnDetails = append(wwnDetails, wd)
		wwns = append(wwns, hw.HostWwn)
	}
	phg.WwnDetails = wwnDetails
	phg.Wwns = wwns

	lupaths := []sanmodel.LuPath{}
	ldevs := []int{}
	hgluns := []int{}
	for _, lp := range hglupaths {
		lp := sanmodel.LuPath{
			Lun:    lp.Lun,
			LdevID: lp.LdevID,
		}
		lupaths = append(lupaths, lp)
		ldevs = append(ldevs, lp.LdevID)
		hgluns = append(hgluns, lp.Lun)
	}
	phg.LuPaths = lupaths
	phg.Ldevs = ldevs
	phg.HgLuns = hgluns

	log.WriteDebug("TFDebug| hostgroup structure=%+v", phg)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), portID, hostGroupNumber)
	return &phg, nil
}

// GetAllHostGroups .
func (psm *sanStorageManager) GetAllHostGroups() (*sanmodel.HostGroups, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_BEGIN), objStorage.Serial)
	hostGroups, err := gatewayObj.GetAllHostGroups()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllHostGroups err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), objStorage.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_END), objStorage.Serial)

	provHostGroups := sanmodel.HostGroups{}
	err = copier.Copy(&provHostGroups, hostGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provHostGroups, nil
}

// CreateHostGroup .
func (psm *sanStorageManager) CreateHostGroup(hgBody sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_BEGIN), hgBody.PortID, hgBody.HostGroupNumber)
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
	objHostgroup := sangatewaymodel.CreateHostGroupReqGwy{
		PortID:          hgBody.PortID,
		HostGroupName:   hgBody.HostGroupName,
		HostGroupNumber: hgBody.HostGroupNumber,
		HostModeOptions: hgBody.HostModeOptions,
		HostMode:        hgBody.HostMode,
	}
	log.WriteDebug("TFDebug| CreateHostGroupReqGwy: %+v", objHostgroup)
	if objHostgroup.HostGroupName != nil {
		log.WriteDebug("TFDebug| HostGroupName: %v", *objHostgroup.HostGroupName)
	}
	if objHostgroup.PortID != nil {
		log.WriteDebug("TFDebug| PortID: %v", *objHostgroup.PortID)
	}

	// since GetHostGroup is by port and HostGroupNumber
	// have to get all first then 
	// use HostGroupName and portID to check if hostgroup already exists
	doCreate := true
	hostGroups, err := gatewayObj.GetAllHostGroups()
	phgNum := hgBody.HostGroupNumber
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllHostGroups err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), objStorage.Serial)
		return nil, err
	}
	for _, hg := range hostGroups.HostGroups {
		log.WriteDebug("TFDebug| HostGroupName: %v, PortID: %v, HostGroupNumber: %v", hg.HostGroupName, hg.PortID, hg.HostGroupNumber)
		if objHostgroup.HostGroupName != nil && objHostgroup.PortID != nil &&
			hg.HostGroupName == *objHostgroup.HostGroupName && hg.PortID == *objHostgroup.PortID {
				hostGroupNumber := hg.HostGroupNumber
				phgNum = &hostGroupNumber
				doCreate = false
				break
		}
	}

	if doCreate {
		_, phgNum, err = gatewayObj.CreateHostGroup(objHostgroup)
		if err != nil || phgNum == nil {
			log.WriteError(mc.GetMessage(mc.ERR_CREATE_HOSTGROUP_FAILED), hgBody.PortID, hgBody.HostGroupNumber)
			log.WriteDebug("TFError| failed to call CreateHostGroup err: %+v", err)
			return nil, err
		}
	}
	log.WriteDebug("TFDebug| phgNum: %v", *phgNum)
	log.WriteDebug("TFDebug| doCreate: %+v", doCreate)
	
	// If user not pass then it will create automatically
	if hgBody.HostGroupNumber == nil {
		hgBody.HostGroupNumber = phgNum
	}
	if len(hgBody.Wwns) > 0 {
		err := psm.AddWwnForNewHostgroup(hgBody)
		if err != nil {
			log.WriteDebug("TFError| failed to call AddWwnForNewHostgroup err: %+v", err)
			return nil, err
		}
	}

	if len(hgBody.Ldevs) > 0 {
		err := psm.AddLdevForNewHostgroup(hgBody)
		if err != nil {
			log.WriteDebug("TFError| failed to call AddLdevForNewHostgroup err: %+v", err)
			return nil, err
		}
	}
	hg, err := psm.GetHostGroup(*hgBody.PortID, *phgNum)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetHostGroup err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), *hgBody.PortID, *phgNum)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_END), hgBody.PortID, hgBody.HostGroupNumber)
	return hg, nil
}

// AddWwnForNewHostgroup will add WWN and Nickname for new hostgroup
func (psm *sanStorageManager) AddWwnForNewHostgroup(hgBody sanmodel.CreateHostGroupRequest) error {
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

	if len(hgBody.Wwns) > 0 {
		// check if wwn is already used, before create
		for _, wwn := range hgBody.Wwns {
			updReq := sangatewaymodel.AddWwnToHgReqGwy{
				HostWwn:         &wwn.Wwn,
				PortID:          hgBody.PortID,
				HostGroupNumber: hgBody.HostGroupNumber,
			}
			// Add WWN
			err := gatewayObj.AddWwnToHG(updReq)
			if err != nil {
				log.WriteDebug("TFError| failed to call AddWwnToHG err: %+v", err)
				log.WriteError(mc.GetMessage(mc.ERR_ADD_WWN_TO_HOSTGROUP_FAILED), wwn.Wwn)
				return err
			}
			// Set NickName for WWN
			err = gatewayObj.SetHostWwnNickName(*hgBody.PortID, *hgBody.HostGroupNumber, wwn.Wwn, wwn.Name)
			if err != nil {
				log.WriteDebug("TFError| failed to call SetHostWwnNickName err: %+v", err)
				log.WriteError(mc.GetMessage(mc.ERR_SET_NICKNAME_TO_WWN_FAILED), wwn.Name, wwn.Wwn)
				return err
			}
		}
	}
	return nil
}

// AddLdevForNewHostgroup add Ldev for new hostgroup
func (psm *sanStorageManager) AddLdevForNewHostgroup(hgBody sanmodel.CreateHostGroupRequest) error {
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

	if len(hgBody.Ldevs) > 0 {
		// check if ldev is already used, before create
		for _, ldev := range hgBody.Ldevs {
			updReq := sangatewaymodel.AddLdevToHgReqGwy{
				PortID:          hgBody.PortID,
				HostGroupNumber: hgBody.HostGroupNumber,
				LdevID:          ldev.LdevId,
				Lun:             ldev.Lun,
				// no hg lun
			}
			if *ldev.LdevId != -1 {
				err := gatewayObj.AddLdevToHG(updReq)
				if err != nil {
					log.WriteDebug("TFError| failed to call AddLdevToHG err: %+v", err)
					log.WriteError(mc.GetMessage(mc.ERR_ADD_LDEV_TO_HOSTGROUP_FAILED), updReq.LdevID)
					return err
				}
			} else {
				log.WriteDebug("TFError| LdevID is -1 meaning ldev_id was not provided in the input")
			}
		}
	}
	return nil
}

// DeleteHostGroup will delete the hostgroup
func (psm *sanStorageManager) DeleteHostGroup(portID string, hostGroupNumber int) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_BEGIN), portID, hostGroupNumber)
	err = gatewayObj.DeleteHostGroup(portID, hostGroupNumber)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteHostGroup err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_HOSTGROUP_FAILED), portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_END), portID, hostGroupNumber)

	return nil
}

// SetHostGroupModeAndOptions .
func (psm *sanStorageManager) SetHostGroupModeAndOptions(portID string, hostGroupNumber int, reqBody sanmodel.SetHostModeAndOptions) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_HOSTGROUP_BEGIN), portID, hostGroupNumber)
	modeRequest := sangatewaymodel.SetHostModeAndOptions{
		HostMode:        reqBody.HostMode,
		HostModeOptions: reqBody.HostModeOptions,
	}
	err = gatewayObj.SetHostGroupModeAndOptions(portID, hostGroupNumber, modeRequest)
	if err != nil {
		log.WriteDebug("TFError| failed to call SetHostGroupModeAndOptions err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_MODE_OPTION_HOSTGROUP_FAILED), portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_MODE_OPTION_HOSTGROUP_END), portID, hostGroupNumber)

	return nil
}

// AddWwnToHGv will add wwn to hostgroup
func (psm *sanStorageManager) AddWwnToHG(reqBody sanmodel.AddWwnToHg) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_WWN_TO_HOSTGROUP_BEGIN), reqBody.HostWwn)
	wwnReq := sangatewaymodel.AddWwnToHgReqGwy{
		HostWwn:         reqBody.HostWwn,
		PortID:          reqBody.PortID,
		HostGroupNumber: reqBody.HostGroupNumber,
	}
	err = gatewayObj.AddWwnToHG(wwnReq)
	if err != nil {
		log.WriteDebug("TFError| failed to call AddWwnToHG err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_ADD_WWN_TO_HOSTGROUP_FAILED), wwnReq.HostWwn)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_WWN_TO_HOSTGROUP_END), reqBody.HostWwn)

	return nil
}

// SetHostWwnNickName will set nick name for wwn
func (psm *sanStorageManager) SetHostWwnNickName(portID string, hostGroupNumber int, hostWwn string, wwnNickname string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_SET_NICKNAME_HOSTGROUP_BEGIN), wwnNickname, hostWwn)
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

	err = gatewayObj.SetHostWwnNickName(portID, hostGroupNumber, hostWwn, wwnNickname)
	if err != nil {
		log.WriteDebug("TFError| failed to call SetHostWwnNickName err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_SET_NICKNAME_TO_WWN_FAILED), wwnNickname, hostWwn)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_NICKNAME_HOSTGROUP_END), wwnNickname, hostWwn)

	return nil
}

// DeleteWwn will delete wwn
func (psm *sanStorageManager) DeleteWwn(portID string, hostGroupNumber int, wwn string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_WWN_BEGIN), wwn, portID, hostGroupNumber)
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

	err = gatewayObj.DeleteWwn(portID, hostGroupNumber, wwn)
	if err != nil {
		log.WriteDebug("TFError| failed to call DeleteWwn err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_WWN_FAILED), wwn, portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_WWN_END), wwn, portID, hostGroupNumber)

	return nil
}

// RemoveLdevFromHG will remove LU Path from Hostgroup
func (psm *sanStorageManager) RemoveLdevFromHG(portID string, hostGroupNumber int, lunID int) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_REMOVE_LUN_HOSTGROUP_BEGIN), lunID, portID, hostGroupNumber)
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

	err = gatewayObj.RemoveLdevFromHG(portID, hostGroupNumber, lunID)
	if err != nil {
		log.WriteDebug("TFError| failed to call RemoveLdevFromHG err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_REMOVE_LUN_HOSTGROUP_FAILED), lunID, portID, hostGroupNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.ERR_REMOVE_LUN_HOSTGROUP_FAILED), lunID, portID, hostGroupNumber)

	return nil
}

// AddLdevToHG will add LuPath to Host group
func (psm *sanStorageManager) AddLdevToHG(reqBody sanmodel.AddLdevToHg) error {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_LDEV_TO_HOSTGROUP_BEGIN), reqBody.LdevID)
	luReq := sangatewaymodel.AddLdevToHgReqGwy{
		PortID:          reqBody.PortID,
		HostGroupNumber: reqBody.HostGroupNumber,
		LdevID:          reqBody.LdevID,
		Lun:             reqBody.Lun,
	}
	if *reqBody.LdevID != -1 {
		err = gatewayObj.AddLdevToHG(luReq)
		if err != nil {
			log.WriteDebug("TFError| failed to call AddLdevToHG err: %+v", err)
			log.WriteError(mc.GetMessage(mc.ERR_ADD_LDEV_TO_HOSTGROUP_FAILED), reqBody.LdevID)
			return err
		}
	} else {
		log.WriteDebug("TFError| LdevID is -1 meaning ldev_id was not provided in the input")
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_LDEV_TO_HOSTGROUP_END), reqBody.LdevID)

	return nil
}
