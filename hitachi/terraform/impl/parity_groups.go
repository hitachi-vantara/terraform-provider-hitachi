package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	"fmt"
	"strings"

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

	driveTypeName := ""
	if v, ok := d.GetOk("drive_type_name"); ok {
		driveTypeName = canonicalizeParityGroupDriveTypeName(v.(string))
		log.WriteDebug("TFDebug| GetParityGroups: drive_type_name filter: %s", driveTypeName)
	}

	var clprID *int
	if v, ok := d.GetOkExists("clpr_id"); ok {
		id := v.(int)
		clprID = &id
		log.WriteDebug("TFDebug| GetParityGroups: clpr_id filter: %d", *clprID)
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

	// Extract include_detail_info and include_cache_info parameters
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, add FMC
		detailInfoTypes = append(detailInfoTypes, "FMC")
		log.WriteDebug("TFDebug| GetParityGroups: include_detail_info=true, added FMC")
	}
	if includeCacheInfo, ok := d.GetOk("include_cache_info"); ok && includeCacheInfo.(bool) {
		// When includeCacheInfo is true, add class
		detailInfoTypes = append(detailInfoTypes, "class")
		log.WriteDebug("TFDebug| GetParityGroups: include_cache_info=true, added class")
	}

	detailInfoType := ""
	if len(detailInfoTypes) > 0 {
		detailInfoType = strings.Join(detailInfoTypes, ",")
		log.WriteDebug("TFDebug| GetParityGroups: final detailInfoType: %s", detailInfoType)
	} else {
		log.WriteDebug("TFDebug| GetParityGroups: no detail info parameters set")
	}
	var parityGroups *[]reconcilermodel.ParityGroup
	if detailInfoType != "" {
		parityGroups, err = reconObj.GetParityGroups(detailInfoType, driveTypeName, clprID, parityGroupIds)
	} else {
		parityGroups, err = reconObj.GetParityGroups("", driveTypeName, clprID, parityGroupIds)
	}
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

	if driveTypeName != "" || clprID != nil {
		filtered := make([]terraformmodel.ParityGroup, 0, len(terraformModelParityGroups))
		for _, pg := range terraformModelParityGroups {
			if driveTypeName != "" && !strings.EqualFold(pg.DriveTypeName, driveTypeName) {
				continue
			}
			if clprID != nil && pg.ClprId != *clprID {
				continue
			}
			filtered = append(filtered, pg)
		}
		terraformModelParityGroups = filtered
		log.WriteDebug("TFDebug| GetParityGroups: filtered parity groups count: %d", len(terraformModelParityGroups))
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

	// Extract include_detail_info and include_cache_info parameters
	var detailInfoTypes []string
	if includeDetailInfo, ok := d.GetOk("include_detail_info"); ok && includeDetailInfo.(bool) {
		// When includeDetailInfo is true, add FMC
		detailInfoTypes = append(detailInfoTypes, "FMC")
		log.WriteDebug("TFDebug| GetParityGroup: include_detail_info=true, added FMC")
	}
	if includeCacheInfo, ok := d.GetOk("include_cache_info"); ok && includeCacheInfo.(bool) {
		// When includeCacheInfo is true, add class
		detailInfoTypes = append(detailInfoTypes, "class")
		log.WriteDebug("TFDebug| GetParityGroup: include_cache_info=true, added class")
	}

	detailInfoType := ""
	if len(detailInfoTypes) > 0 {
		detailInfoType = strings.Join(detailInfoTypes, ",")
		log.WriteDebug("TFDebug| GetParityGroup: final detailInfoType: %s", detailInfoType)
	} else {
		log.WriteDebug("TFDebug| GetParityGroup: no detail info parameters set")
	}

	var parityGroup *reconcilermodel.ParityGroup

	if detailInfoType != "" {
		// Use the multi-parity group method with filtering
		parityGroups, err := reconObj.GetParityGroups(detailInfoType, "", nil, []string{pgid})
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), setting.Serial)
			return nil, err
		}
		if len(*parityGroups) == 0 {
			log.WriteError("Parity group not found: %s", pgid)
			return nil, fmt.Errorf("parity group not found: %s", pgid)
		}
		parityGroup = &(*parityGroups)[0]
	} else {
		parityGroup, err = reconObj.GetParityGroup(pgid)
		if err != nil {
			log.WriteError(mc.GetMessage(mc.ERR_GET_PARITY_GROUP_FAILED), setting.Serial)
			return nil, err
		}
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

func canonicalizeParityGroupDriveTypeName(input string) string {
	allowed := []string{"SAS", "SCM", "SSD", "SSD(FMC)", "SSD(MLC)", "SSD(QLC)", "SSD(RI)"}
	for _, v := range allowed {
		if strings.EqualFold(v, input) {
			return v
		}
	}
	return input
}

func ConvertParityGroupToSchema(parityGroup *terraformmodel.ParityGroup, serial int) *map[string]interface{} {
	// Convert spaces array
	spacesList := []map[string]interface{}{}
	for _, space := range parityGroup.Spaces {
		spaceMap := map[string]interface{}{
			"partition_number": space.PartitionNumber,
			"ldev_id":          space.LdevId,
			"status":           space.Status,
			"lba_location":     space.LbaLocation,
			"lba_size":         space.LbaSize,
		}
		spacesList = append(spacesList, spaceMap)
	}

	parity := map[string]interface{}{
		"storage_serial_number":              serial,
		"parity_group_id":                    parityGroup.ParityGroupId,
		"group_type":                         parityGroup.GroupType,
		"num_of_ldevs":                       parityGroup.NumberOfLdevs,
		"used_capacity_rate":                 parityGroup.UsedCapacityRate,
		"available_volume_capacity":          parityGroup.AvailableVolumeCapacity,
		"raid_level":                         parityGroup.RaidLevel,
		"raid_type":                          parityGroup.RaidType,
		"clpr_id":                            parityGroup.ClprId,
		"drive_type":                         parityGroup.DriveType,
		"drive_type_name":                    parityGroup.DriveTypeName,
		"is_copy_back_mode_enabled":          parityGroup.IsCopyBackModeEnabled,
		"is_encryption_enabled":              parityGroup.IsEncryptionEnabled,
		"total_capacity":                     parityGroup.TotalCapacity,
		"physical_capacity":                  parityGroup.PhysicalCapacity,
		"available_physical_capacity":        parityGroup.AvailablePhysicalCapacity,
		"is_accelerated_compression_enabled": parityGroup.IsAcceleratedCompressionEnabled,
		"spaces":                             spacesList,
		"emulation_type":                     parityGroup.EmulationType,
		"available_volume_capacity_in_kb":    parityGroup.AvailableVolumeCapacityInKB,
	}

	return &parity
}
