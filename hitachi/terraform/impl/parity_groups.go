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

	// "terraform-provider-hitachi/hitachi/common/utils"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetParityGroups(d *schema.ResourceData) (*[]terraformmodel.ParityGroup, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), setting.Serial)
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	var parityGroupIds []string
	ids, ok := d.GetOk("parity_group_ids")
	if ok {
		pgIds := ids.([]interface{})
		parityGroupArray := make([]string, len(pgIds))
		for index, value := range pgIds {
			switch typedValue := value.(type) {
			case string:
				parityGroupArray[index] = typedValue
			}
		}
		parityGroupIds = parityGroupArray
		log.WriteDebug("TFDebug| parity group filter will be apply")
	}

	parityGroups, err := reconObj.GetParityGroups(parityGroupIds)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelParityGroups := []terraformmodel.ParityGroup{}
	err = copier.Copy(&terraformModelParityGroups, parityGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), setting.Serial)

	return &terraformModelParityGroups, nil
}

func GetParityGroup(d *schema.ResourceData) (*terraformmodel.ParityGroup, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_BEGIN), setting.Serial)
	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	pgid := d.Get("parity_group_id").(string)

	parityGroup, err := reconObj.GetParityGroup(pgid)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelParityGroup := terraformmodel.ParityGroup{}
	err = copier.Copy(&terraformModelParityGroup, parityGroup)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PARITY_GROUP_END), setting.Serial)

	return &terraformModelParityGroup, nil
}

func ConvertParityGroupToSchema(parityGroup *terraformmodel.ParityGroup, serial int) *map[string]interface{} {
	parity := map[string]interface{}{
		"storage_serial_number":              serial,
		"parity_group_id":                    parityGroup.ParityGroupId,
		"num_of_ldevs":                       parityGroup.NumberOfLdevs,
		"used_capacity_rate":                 parityGroup.UsedCapacityRate,
		"available_volume_capacity":          parityGroup.AvailableVolumeCapacity,
		"raid_level":                         parityGroup.RaidLevel,
		"raid_type":                          parityGroup.RaidType,
		"clpr_id":                            parityGroup.ClprId,
		"drive_type":                         parityGroup.DriveType,
		"drive_type_name":                    parityGroup.DriveTypeName,
		"total_capacity":                     parityGroup.TotalCapacity,
		"physical_capacity":                  parityGroup.PhysicalCapacity,
		"available_physical_capacity":        parityGroup.AvailablePhysicalCapacity,
		"is_accelerated_compression_enabled": parityGroup.IsAcceleratedCompressionEnabled,
		"available_volume_capacity_in_kb":    parityGroup.AvailableVolumeCapacityInKB,
	}

	return &parity
}
