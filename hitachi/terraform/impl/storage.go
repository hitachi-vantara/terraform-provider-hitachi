package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"
	"strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	infragwmodel "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimplinfragw "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	reconimplvssb "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/impl"
	reconcilermodelvssb "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func RegisterStorageSystem(d *schema.ResourceData) (*terraformmodel.AllStorageTypes, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ssList := []*terraformmodel.StorageSystem{}
	var err error
	ss_items := d.Get("san_storage_system").([]interface{})
	if len(ss_items) > 0 {
		ssList, err = GetSanStorageSystem(ss_items)
		if err != nil {
			log.WriteDebug("TFError| error in GetSanStorageSystem, err: %v", err)
			return nil, err
		}
	}

	vssbList := []*terraformmodel.StorageVersionInfo{}
	ss_vssb_items := d.Get("hitachi_vss_block_provider").([]interface{})
	if len(ss_vssb_items) > 0 {
		vssbList, err = GetVssbStorageSystem(ss_vssb_items)
		if err != nil {
			log.WriteDebug("TFError| error in GetVssbStorageSystem, err: %v", err)
			return nil, err
		}
	}

	infra_gw_list := []*terraformmodel.InfraGwSettings{}
	ss_ingra_gw_items := d.Get("hitachi_infrastructure_gateway_provider").([]interface{})
	if len(ss_ingra_gw_items) > 0 {
		infra_gw_list, err = GetInfraGwSystem(ss_ingra_gw_items)
		if err != nil {
			log.WriteDebug("TFError| error in GetVssbStorageSystem, err: %v", err)
			return nil, err
		}
	}

	allStorageTypes := terraformmodel.AllStorageTypes{}
	if ssList != nil {
		allStorageTypes.VspStorageSystem = append(allStorageTypes.VspStorageSystem, ssList...)
	}
	if vssbList != nil {
		allStorageTypes.VssbStorageVersionInfo = append(allStorageTypes.VssbStorageVersionInfo, vssbList...)
	}
	if infra_gw_list != nil {
		allStorageTypes.InfraGwInfo = append(allStorageTypes.InfraGwInfo, infra_gw_list...)
	}
	return &allStorageTypes, nil
}

func GetSanStorageSystem(ssItems []interface{}) (ssList []*terraformmodel.StorageSystem, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	for _, item := range ssItems {
		i := item.(map[string]interface{})

		serial := i["serial"].(int)
		mgmtIP := i["management_ip"].(string)
		usernameEncoded := i["username"].(string)
		passwordEncoded := i["password"].(string)

		storageSetting := reconcilermodel.StorageDeviceSettings{
			Serial:   serial,
			Username: usernameEncoded,
			Password: passwordEncoded,
			MgmtIP:   mgmtIP,
		}

		reconObj, err := reconimpl.NewEx(storageSetting)
		if err != nil {
			log.WriteDebug("TFError| error in NewEx, err: %v", err)
			return nil, err
		}

		storageSystem, err := reconObj.GetStorageSystem()
		if err != nil {
			log.WriteDebug("TFError| error getting storage system, err: %v", err)
			return nil, err
		}

		settingAndInfo := reconcilermodel.StorageSettingsAndInfo{
			Settings: storageSetting,
			Info:     storageSystem,
		}

		// save this to a cache
		cache.WriteToSanCache(strconv.Itoa(serial), settingAndInfo)

		terraformStorageSystem := terraformmodel.StorageSystem{}
		err = copier.Copy(&terraformStorageSystem, storageSystem)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

		ssList = append(ssList, &terraformStorageSystem)
	}
	return
}

func GetInfraGwSystem(ssVssbItems []interface{}) (ssList []*terraformmodel.InfraGwSettings, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	for _, item := range ssVssbItems {
		i := item.(map[string]interface{})

		address := i["address"].(string)
		username := i["username"].(string)
		password := i["password"].(string)

		setting := infragwmodel.InfraGwSettings{
			Username: username,
			Password: password,
			Address:  address,
		}

		reconObj, err := reconimplinfragw.NewEx(setting)
		if err != nil {
			log.WriteDebug("TFError| error in NewEx, err: %v", err)
			return nil, err
		}

		mtDetails, err := reconObj.GetPartnerAndSubscriberId(username)

		if err == nil {
			setting.PartnerId = &mtDetails.PartnerId
			setting.SubscriberId = &mtDetails.SubscriberId
		} else {
			log.WriteDebug("TFError| error in GetPartnerAndSubscriberId, err: %v", err)
			return nil, err
		}

		storageDevices, err := reconObj.GetStorageDevices()
		if err != nil {
			log.WriteDebug("TFError| error in NewEx, err: %v", err)
			return nil, err
		}

		m := make(map[string]string)
		mm := make(map[string]string)
		for _, x := range storageDevices.Data {
			m[x.SerialNumber] = x.ResourceId
			mm[x.ResourceId] = x.SerialNumber
		}

		if err != nil {
			log.WriteDebug("TFError| error getting storage system, err: %v", err)
			return nil, err
		}

		iscsiMap := make(map[string]map[string]string)
		for _, val := range m {
			iscsiTargets, err := reconObj.GetIscsiTargets(val, "")
			if err != nil {
				log.WriteDebug("TFError| error getting iscsi targets, err: %v", err)
				return nil, err
			}
			m2 := make(map[string]string)
			for _, x := range iscsiTargets.Data {
				key := x.PortId + "_" + strconv.Itoa(x.ISCSIId)
				m2[key] = x.ResourceId
			}
			iscsiMap[val] = m2
		}

		settingAndInfo := infragwmodel.InfraGwStorageSettingsAndInfo{
			Settings:          setting,
			SerialToStorageId: m,
			StorageIdToSerial: mm,
			IscsiTargetIdMap:  iscsiMap,
		}

		log.WriteDebug("TFDebug| Infra GW Info: %v", settingAndInfo)

		// save this to a cache
		cache.SetCurrentAddress(address)
		cache.WriteToInfraGwCache(address, settingAndInfo)

		terraformInfo := terraformmodel.InfraGwSettings{}
		err = copier.Copy(&terraformInfo, setting)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

		ssList = append(ssList, &terraformInfo)
	}
	return
}

func GetVssbStorageSystem(ssVssbItems []interface{}) (ssList []*terraformmodel.StorageVersionInfo, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	for _, item := range ssVssbItems {
		i := item.(map[string]interface{})

		mgmtIP := i["vss_block_address"].(string)
		usernameEncoded := i["username"].(string)
		passwordEncoded := i["password"].(string)

		storageSetting := reconcilermodelvssb.StorageDeviceSettings{
			Username:       usernameEncoded,
			Password:       passwordEncoded,
			ClusterAddress: mgmtIP,
		}

		reconObj, err := reconimplvssb.NewEx(storageSetting)
		if err != nil {
			log.WriteDebug("TFError| error in NewEx, err: %v", err)
			return nil, err
		}

		versionInfo, err := reconObj.GetStorageVersionInfo()
		if err != nil {
			log.WriteDebug("TFError| error getting storage system, err: %v", err)
			return nil, err
		}

		settingAndInfo := reconcilermodelvssb.StorageSettingsAndInfo{
			Settings: storageSetting,
			Info:     versionInfo,
		}

		// save this to a cache
		cache.WriteToVssbCache(mgmtIP, settingAndInfo)

		terraformVersionInfo := terraformmodel.StorageVersionInfo{}
		err = copier.Copy(&terraformVersionInfo, versionInfo)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

		ssList = append(ssList, &terraformVersionInfo)
	}
	return
}

func GetStorageSystem(d *schema.ResourceData) (*terraformmodel.StorageSystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_BEGIN), setting.MgmtIP)
	reconStorageSystem, err := reconObj.GetStorageSystem()
	if err != nil {
		log.WriteDebug("TFError| error getting storage system, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_FAILED), setting.MgmtIP)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStorageSystem := terraformmodel.StorageSystem{}
	err = copier.Copy(&terraformStorageSystem, reconStorageSystem)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_END), setting.MgmtIP)

	return &terraformStorageSystem, nil
}
