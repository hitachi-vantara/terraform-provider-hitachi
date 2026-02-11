package terraform

import (
	"fmt"
	"strconv"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
)

// StopAllVolumeFormat invokes reconciler to call the appliance-wide stop-format action.
// Accepts a `schema.ResourceData` (so the resource can pass serial); returns error on failure.
func StopAllVolumeFormat(d interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// support both *schema.ResourceData and map[string]interface{} for easier testing
	var serial int

	switch rd := d.(type) {
	case interface{ Get(string) interface{} }:
		serial = rd.Get("serial").(int)
	default:
		return fmt.Errorf("unsupported input to StopAllVolumeFormat")
	}

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		return err
	}

	if err := reconObj.StopAllVolumeFormat(); err != nil {
		return err
	}

	return nil
}
