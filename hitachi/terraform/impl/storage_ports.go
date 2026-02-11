package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	"strconv"
	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetStoragePorts(d *schema.ResourceData) (*[]terraformmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// Extract include_detail_info parameter
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, pass all available detail types
		detailInfoTypes = []string{"logins", "portMode", "loginHostNqn"}
		log.WriteDebug("TFDebug| GetStoragePorts: include_detail_info=true, using detailInfoTypes: %v", detailInfoTypes)
	} else {
		log.WriteDebug("TFDebug| GetStoragePorts: include_detail_info=false or not set, no detailInfoType")
	}

	// Extract portType and portAttributes parameters
	var portType, portAttributes string
	if pt, ok := d.GetOk("port_type"); ok {
		portType = pt.(string)
		log.WriteDebug("TFDebug| GetStoragePorts: portType=%s", portType)
	}
	if pa, ok := d.GetOk("port_attributes"); ok {
		portAttributes = pa.(string)
		log.WriteDebug("TFDebug| GetStoragePorts: portAttributes=%s", portAttributes)
	}

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
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_BEGIN), setting.Serial)
	reconStoragePorts, err := reconObj.GetStoragePorts(detailInfoTypes, portType, portAttributes)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePorts, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_PORTS_FAILED), setting.Serial)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStoragePorts := []terraformmodel.StoragePort{}
	err = copier.Copy(&terraformStoragePorts, reconStoragePorts)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_END), setting.Serial)

	return &terraformStoragePorts, nil
}

func GetStoragePortByPortId(d *schema.ResourceData) (*terraformmodel.StoragePort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	portID := d.Get("port_id").(string)

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
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_BEGIN), portID, setting.Serial)
	reconStoragePort, err := reconObj.GetStoragePortByPortId(portID)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStoragePortByPortId, err: %v", err)
		log.WriteError(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_BEGIN), portID, setting.Serial)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStoragePort := terraformmodel.StoragePort{}
	err = copier.Copy(&terraformStoragePort, reconStoragePort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_PORTS_PORTID_BEGIN), portID, setting.Serial)

	return &terraformStoragePort, nil
}

func ConvertStoragePortToSchema(storagePort *terraformmodel.StoragePort) *map[string]interface{} {
	sp := map[string]interface{}{
		"port_id":              storagePort.PortId,
		"port_type":            storagePort.PortType,
		"port_attributes":      storagePort.PortAttributes,
		"port_speed":           storagePort.PortSpeed,
		"loop_id":              storagePort.LoopId,
		"fabric_mode":          storagePort.FabricMode,
		"port_connection":      storagePort.PortConnection,
		"lun_security_setting": storagePort.LunSecuritySetting,
		"wwn":                  storagePort.Wwn,
		"port_mode":            storagePort.PortMode,
	}

	return &sp
}
