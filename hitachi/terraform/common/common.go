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
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	reconcilermodel "terraform-provider-hitachi/hitachi/infra_gw/model"

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
		err_text := fmt.Sprintf("Serial Number %s does not exit on %s", serial, address)
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
		err_text := fmt.Sprintf("Storage Id %s does not exit on %s", storageId, address)
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

func GetSerialStringFromDiff(d *schema.ResourceDiff) string {
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

func ValidateSerialAndStorageId(serial, storageId string) error {
	if serial == "" && storageId == "" {
		err := errors.New("both serial and storage_id can't be empty. Please specify one")
		return err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
		return err
	}
	return nil
}

func GbToMbString(gb int) string {
	mb := gb * 1024
	return fmt.Sprintf("%dMB", mb)
}

func GetValidateStorageIDFromSerial(d *schema.ResourceDiff) (*string, error) {

	serial := GetSerialStringFromDiff(d)
	storageId := d.Get("storage_id").(string)

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
		return nil, err
	} else if serial == "" && storageId == "" {
		err := errors.New("either serial or storage_id can't be empty. Please specify one")
		return nil, err
	}
	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	if storageId == "" {
		storageId, err = GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, err
		}
	}
	return &storageId, nil
}

func GetValidateStorageIDFromSerialResource(d *schema.ResourceData, m interface{}) (*string, *string, error) {
	if m != nil {

		providerLists := m.(*terraformmodel.AllStorageTypes)

		if len(providerLists.InfraGwInfo) == 0 {
			return nil, nil, nil
		}
	}

	serial := GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return nil, nil, err
	}

	var storage_serial_number int
	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, nil, err
	}
	if storageId == "" {
		storageId, err = GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, nil, err
		}
	}

	if serial == "" {
		serial, err = GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	}
	d.Set("serial", storage_serial_number)

	return &storageId, &address, nil
}

func GetInfraGatewaySettings(d *schema.ResourceData, m interface{}) (*string, *reconcilermodel.InfraGwSettings, error) {

	storage_id, address, err := GetValidateStorageIDFromSerialResource(d, m)

	if err != nil {
		return nil, nil, err
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(*address)
	if err != nil {
		return nil, nil, err
	}
	setting := reconcilermodel.InfraGwSettings(*storageSetting)

	if setting.PartnerId != nil {
		subId, ok := d.GetOk("subscriber_id")
		if ok {
			subIdw := subId.(string)
			setting.SubscriberId = &subIdw
		}
	}

	return storage_id, &setting, nil
}
