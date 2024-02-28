package common

import (
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	infragwmodel "terraform-provider-hitachi/hitachi/infra_gw/model"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"
)

var SanSettingsMap map[string]sanmodel.StorageSettingsAndInfo
var VssbSettingsMap map[string]vssbmodel.StorageSettingsAndInfo
var InfraGwSettingsMap map[string]infragwmodel.InfraGwStorageSettingsAndInfo

var currentAddress map[string]string

func SetCurrentAddress(a string) {
	currentAddress["current_address"] = a
}

func GetCurrentAddress() (string, error) {
	value, ok := currentAddress["current_address"]
	if ok {
		return value, nil
	}
	return value, fmt.Errorf("current_address not found in the cache")
}

func init() {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("cache init function has been called")

	SanSettingsMap = make(map[string]sanmodel.StorageSettingsAndInfo)
	VssbSettingsMap = make(map[string]vssbmodel.StorageSettingsAndInfo)
	InfraGwSettingsMap = make(map[string]infragwmodel.InfraGwStorageSettingsAndInfo)

	currentAddress = make(map[string]string)
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

func WriteToInfraGwCache(key string, data infragwmodel.InfraGwStorageSettingsAndInfo) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	//Uncomment when required
	//log.WriteDebug("key: %+v  data: %v type %v", key, data, reflect.TypeOf(data))
	InfraGwSettingsMap[key] = data
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
	return value, fmt.Errorf("%s not found in the cache", key)
}

func ReadFromInfraGwCache(key string) (infragwmodel.InfraGwStorageSettingsAndInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("ReadFromCache key: %+v ", key)

	value, ok := InfraGwSettingsMap[key]
	if ok {
		return value, nil
	}
	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("key: %+v  data: %v", key, value)
	return value, fmt.Errorf("%s not found in the cache", key)
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
	return value, fmt.Errorf("%s not found in the cache", key)
}

func GetSanSettingsFromCache(serialNumber string) (*sanmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromSanCache(serialNumber)

	if serialNumber == "0" {
		return nil, fmt.Errorf("valid serial number is not provided for direct connection")
	}
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

func GetInfraSettingsFromCache(infraGwAddr string) (*infragwmodel.InfraGwSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromInfraGwCache(infraGwAddr)

	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetInfraGwSerialToIdMap(infraGwAddr string) (*map[string]string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromInfraGwCache(infraGwAddr)

	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	return &data.SerialToStorageId, nil
}

func GetInfraGwIdToSerialMap(infraGwAddr string) (*map[string]string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromInfraGwCache(infraGwAddr)

	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	log.WriteDebug("data: %+v", data)
	log.WriteDebug("StorageIdToSerialMap: %+v", data.StorageIdToSerial)

	return &data.StorageIdToSerial, nil
}

func GetInfraGwStorageIdToIscsiIdMap(infraGwAddr, storageId string) (*map[string]string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data, err := ReadFromInfraGwCache(infraGwAddr)

	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	// DO NOT UMCOMMENT THIS LINE, it prints username/password in the log file
	//log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	IscsiTargetIdMap := data.IscsiTargetIdMap
	m := IscsiTargetIdMap[storageId]
	return &m, nil
}
