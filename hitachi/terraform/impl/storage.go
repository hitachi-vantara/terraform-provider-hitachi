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

	reconimpladmin "terraform-provider-hitachi/hitachi/storage/admin/reconciler/impl"
	reconcilermodeladmin "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	reconimplvssb "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodelvssb "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
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
	ss_vosb_items := d.Get("hitachi_vosb_provider").([]interface{})
	if len(ss_vosb_items) > 0 {
		vssbList, err = GetVssbStorageSystem(ss_vosb_items)
		if err != nil {
			log.WriteDebug("TFError| error in GetVssbStorageSystem, err: %v", err)
			return nil, err
		}
	}

	ssAdminList := []*terraformmodel.StorageSystemAdmin{}
	ss_admin_items := d.Get("hitachi_vsp_one_provider").([]interface{})
	if len(ss_admin_items) > 0 {
		ssAdminList, err = GetAdminStorageSystem(ss_admin_items)
		if err != nil {
			log.WriteDebug("TFError| error in GetStorageSystemAdmin, err: %v", err)
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
	if ssAdminList != nil {
		allStorageTypes.AdminStorageSystem = append(allStorageTypes.AdminStorageSystem, ssAdminList...)
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

		storageSystem, err := reconObj.GetStorageSystemInfo()
		if err != nil {
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

func GetVssbStorageSystem(ssVssbItems []interface{}) (ssList []*terraformmodel.StorageVersionInfo, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	for _, item := range ssVssbItems {
		i := item.(map[string]interface{})

		mgmtIP := i["vosb_address"].(string)
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
			// log.WriteDebug("TFError| error getting storage system, err: %v", err)
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

func GetAdminStorageSystem(ssItems []interface{}) (ssList []*terraformmodel.StorageSystemAdmin, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	for _, item := range ssItems {
		i := item.(map[string]interface{})

		serial := i["serial"].(int)
		mgmtIP := i["management_ip"].(string)
		usernameEncoded := i["username"].(string)
		passwordEncoded := i["password"].(string)

		storageSetting := reconcilermodeladmin.StorageDeviceSettings{
			Serial:   serial,
			Username: usernameEncoded,
			Password: passwordEncoded,
			MgmtIP:   mgmtIP,
		}

		reconObj, err := reconimpladmin.NewEx(storageSetting)
		if err != nil {
			log.WriteDebug("TFError| error in NewEx, err: %v", err)
			return nil, err
		}

		storageSystem, err := reconObj.GetStorageAdminInfo(true)
		if err != nil {
			return nil, err
		}

		settingAndInfo := reconcilermodeladmin.StorageSettingsAndInfo{
			Settings: storageSetting,
			Info:     storageSystem,
		}

		// save this to a cache
		cache.WriteToAdminCache(strconv.Itoa(serial), settingAndInfo)

		terraformStorageSystem := terraformmodel.StorageSystemAdmin{}
		err = copier.Copy(&terraformStorageSystem, storageSystem)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

		ssList = append(ssList, &terraformStorageSystem)
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

func GetStorageSystemAdmin(d *schema.ResourceData) (*terraformmodel.StorageSystemAdmin, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	configurable_capacities := d.Get("with_estimated_configurable_capacities").(bool)
	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetAdminSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodeladmin.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpladmin.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_ADMIN_BEGIN), setting.MgmtIP)
	reconStorageSystem, err := reconObj.GetStorageAdminInfo(configurable_capacities)
	if err != nil {
		log.WriteDebug("TFError| error getting storage system, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_SYSTEM_ADMIN_FAILED), setting.MgmtIP)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStorageSystem := terraformmodel.StorageSystemAdmin{}
	err = copier.Copy(&terraformStorageSystem, reconStorageSystem)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_SYSTEM_ADMIN_END), setting.MgmtIP)

	return &terraformStorageSystem, nil
}
