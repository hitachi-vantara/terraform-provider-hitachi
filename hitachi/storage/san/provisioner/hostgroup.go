package sanstorage

import (
	// "encoding/json"
	// "errors"
	// "context"
	"fmt"
	// "io/ioutil"
	// "strconv"
	// "time"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// "terraform-provider-hitachi/hitachi/common/utils"
	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sangateway "terraform-provider-hitachi/hitachi/storage/san/gateway"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetHostGroup(storageSetting sanmodel.StorageDeviceSettings, portID string, hostGroupNumber int) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	////////// in parallel
	var wga sync.WaitGroup
	var outMap sync.Map
	var errMap sync.Map

	n := 1
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		hg, err := sangateway.GetHostGroup(storageSetting, portID, hostGroupNumber)
		errMap.Store(n, err)
		outMap.Store(n, hg)
		return
	}(n)

	n = 2
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		hgwwns, err := sangateway.GetHostGroupWwns(storageSetting, portID, hostGroupNumber)
		errMap.Store(n, err)
		outMap.Store(n, hgwwns)
		return
	}(n)

	n = 3
	wga.Add(1)
	go func(n int) {
		defer wga.Done()
		hglupaths, err := sangateway.GetHostGroupLuPaths(storageSetting, portID, hostGroupNumber)
		errMap.Store(n, err)
		outMap.Store(n, hglupaths)
		return
	}(n)

	log.WriteDebug("Waiting for GetHostGroups goroutines to end ...")
	wga.Wait()
	log.WriteDebug("GetHostGroups goroutines ended")

	// check for errors
	for nth := 1; nth <= 3; nth++ {
		serr, _ := errMap.Load(nth)
		if serr != nil {
			log.WriteDebug("serr: %+v", serr.(error))
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

	hg := ihg.(*sanmodel.HostGroupGwy)
	hgwwns := ihgwwns.(*[]sanmodel.HostWwnDetail)
	hglupaths := ihglupaths.(*[]sanmodel.HostLuPath)

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
	for _, hw := range *hgwwns {
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
	for _, lp := range *hglupaths {
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

	log.WriteDebug("hg=%+v", phg)
	return &phg, nil
}

func GetHostGroupModeAndOptions(storageSetting sanmodel.StorageDeviceSettings) (*sanmodel.HostModeAndOptions, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	hgmo, err := sangateway.GetHostGroupModeAndOptions(storageSetting)
	if err != nil {
		return nil, err
	}
	log.WriteDebug("hgmo=%+v", hgmo)
	return hgmo, nil
}

func CreateHostGroup(storageSetting sanmodel.StorageDeviceSettings, hgBody sanmodel.CreateHostGroupRequest) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	crReq := sanmodel.CreateHostGroupReqGwy{
		PortID:          hgBody.PortID,
		HostGroupName:   hgBody.HostGroupName,
		HostGroupNumber: hgBody.HostGroupNumber,
		HostModeOptions: hgBody.HostModeOptions,
		HostMode:        hgBody.HostMode,
	}

	log.WriteDebug("crReq: %+v", crReq)

	// check if already existing
	_, phgNum, err := sangateway.CreateHostGroup(storageSetting, crReq)
	if err != nil {
		log.WriteDebug("ERROR: %+v", err)
		return nil, err
	}

	// these can be in parallel
	if len(hgBody.Wwns) > 0 {
		// check if wwn is already used, before create
		for _, wwn := range hgBody.Wwns {
			updReq := sanmodel.AddWwnToHgReqGwy {
				HostWwn: &wwn.Wwn,
				PortID: hgBody.PortID,
				HostGroupNumber: hgBody.HostGroupNumber,
			}
			err := sangateway.AddWwnToHG(storageSetting, updReq)
			if err != nil {
				log.WriteDebug("ERROR: %+v", err)
				return nil, err
			}
			// update nickname
		}
	}

	if len(hgBody.Ldevs) > 0 {
		// check if ldev is already used, before create
		for _, ldev := range hgBody.Ldevs {
			updReq := sanmodel.AddLdevToHgReqGwy {
				PortID: hgBody.PortID,
				HostGroupNumber: hgBody.HostGroupNumber,
				LdevID: &ldev,
				// no hg lun
			}
			err := sangateway.AddLdevToHG(storageSetting, updReq)
			if err != nil {
				log.WriteDebug("ERROR: %+v", err)
				return nil, err
			}
		}
	}

	hg, err := GetHostGroup(storageSetting, *hgBody.PortID, *phgNum)
	if err != nil {
		return nil, err
	}

	return hg, nil
}
