package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	// "errors"
	// "context"
	// "fmt"
	"io/ioutil"
	// "strconv"
	// "time"
	// log "github.com/romana/rlog"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// sanprov "terraform-provider-hitachi/hitachi/storage/san/provisioner"
)

func GetSanSettingsFromFile(serialNumber string) (*sanmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	filePath := utils.SAN_SETTINGS_DIR + "/" + serialNumber
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	data := sanmodel.StorageSettingsAndInfo{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetVssbSettingsFromFile(vssbAddr string) (*vssbmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	filePath := utils.VSSB_SETTINGS_DIR + "/" + vssbAddr
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	data := vssbmodel.StorageSettingsAndInfo{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetStorageIdFromSerial(address, serial string) (string, error) {

	m, err := cache.GetInfraGwSerialToIdMap(address)
	if err != nil {
		return "", err
	}
	id := (*m)[serial]
	if id == "" {
		err_text := fmt.Sprintf("Serial Number %s does not exit", serial)
		err = errors.New(err_text)
		return "", err
	}
	return id, nil
}

func GetSerialFromStorageId(address, storageId string) (string, error) {

	m, err := cache.GetInfraGwIdToSerialMap(address)
	if err != nil {
		return "", err
	}
	serial := (*m)[storageId]
	if serial == "" {
		err_text := fmt.Sprintf("Storage Id %s does not exit", storageId)
		err = errors.New(err_text)
		return "", err
	}
	return serial, nil
}

func GetSerialString(d *schema.ResourceData) string {
	serial_number := -1
	serial := ""

	sid, okId := d.GetOk("serial")
	if okId {
		serial_number = sid.(int)
		serial = strconv.Itoa(serial_number)
	}
	return serial
}

func GetIscsiTargetId(address, storageId, port string, iscsi_id int) (string, error) {

	m, err := cache.GetInfraGwStorageIdToIscsiIdMap(address, storageId)
	if err != nil {
		return "", err
	}
	key := port + "_" + strconv.Itoa(iscsi_id)
	id := (*m)[key]
	if id == "" {
		err_text := fmt.Sprintf("Either port %s  or iscsi target number %d does not exit", port, iscsi_id)
		err = errors.New(err_text)
		return "", err
	}
	return id, nil
}
