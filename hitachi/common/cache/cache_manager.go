package common

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	adminmodel "terraform-provider-hitachi/hitachi/storage/admin/reconciler/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
)

var SanSettingsMap map[string]sanmodel.StorageSettingsAndInfo
var VssbSettingsMap map[string]vssbmodel.StorageSettingsAndInfo
var AdminSettings map[string]adminmodel.StorageSettingsAndInfo

func init() {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("memory cache init function has been called")

	SanSettingsMap = make(map[string]sanmodel.StorageSettingsAndInfo)
	VssbSettingsMap = make(map[string]vssbmodel.StorageSettingsAndInfo)
	AdminSettings = make(map[string]adminmodel.StorageSettingsAndInfo)
}

func WriteToSanCache(key string, data sanmodel.StorageSettingsAndInfo) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	//Uncomment when required
	// log.WriteDebug("key: %+v  data: %v type %v", key, data, reflect.TypeOf(data))

	SanSettingsMap[key] = data

}

func WriteToVssbCache(key string, data vssbmodel.StorageSettingsAndInfo) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	//Uncomment when required
	//log.WriteDebug("key: %+v  data: %v type %v", key, data, reflect.TypeOf(data))
	VssbSettingsMap[key] = data
}

func WriteToAdminCache(key string, data adminmodel.StorageSettingsAndInfo) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	//Uncomment when required
	//log.WriteDebug("key: %+v  data: %v type %v", key, data, reflect.TypeOf(data))
	AdminSettings[key] = data
}

func ReadFromSanCache(key string) (sanmodel.StorageSettingsAndInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("ReadFromCache key: %+v ", key)

	value, ok := SanSettingsMap[key]
	if ok {
		return value, nil
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("key: %+v  data: %v", key, value)
	return value, fmt.Errorf("storage serial number %v not found", key)
}

func ReadFromVssbCache(key string) (vssbmodel.StorageSettingsAndInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("ReadFromCache key: %+v ", key)

	value, ok := VssbSettingsMap[key]
	if ok {
		return value, nil
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("key: %+v  data: %v", key, value)
	return value, fmt.Errorf("vosb address %v not found", key)
}

func ReadFromAdminCache(key string) (adminmodel.StorageSettingsAndInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("ReadFromCache key: %+v ", key)

	value, ok := AdminSettings[key]
	if ok {
		return value, nil
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("key: %+v  data: %v", key, value)
	return value, fmt.Errorf("storage serial number %v not found", key)
}

func GetSanSettingsFromCache(serialNumber string) (*sanmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromSanCache(serialNumber)

	//Uncomment if required
	// log.WriteDebug("ssi: %+v type: %v", data, reflect.TypeOf(data))

	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetVssbSettingsFromCache(vssbAddr string) (*vssbmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromVssbCache(vssbAddr)

	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetAdminSettingsFromCache(serialNumber string) (*adminmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromAdminCache(serialNumber)

	//Uncomment if required
	// log.WriteDebug("ssi: %+v type: %v", data, reflect.TypeOf(data))

	if err != nil {
		log.WriteError(err)
		return nil, err
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}
