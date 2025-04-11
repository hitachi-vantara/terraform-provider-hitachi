package terraform

import (
	"encoding/json"
	// "errors"
	// "context"
	// "fmt"
	"io/ioutil"
	// "strconv"
	// "time"
	// log "github.com/romana/rlog"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	// sanprov "terraform-provider-hitachi/hitachi/storage/san/provisioner"
)

// file not used

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

	filePath := utils.VOSB_SETTINGS_DIR + "/" + vssbAddr
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
