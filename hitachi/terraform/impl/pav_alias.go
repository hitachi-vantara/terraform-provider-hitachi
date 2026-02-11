package terraform

import (
	"strconv"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewaymodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetPavAliases(d *schema.ResourceData) (*[]gatewaymodel.PavAlias, error) {
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

	baseLdevIDRaw, hasBaseLdev := d.GetOk("base_ldev_id")
	aliasLdevIDs := extractIntList(d, "alias_ldev_ids")

	var cuNumberPtr *int
	if v, ok := d.GetOk("cu_number"); ok {
		cu := v.(int)
		cuNumberPtr = &cu
	}

	// If we need to filter by base/alias LDEVs, fetch all entries.
	if hasBaseLdev || len(aliasLdevIDs) > 0 {
		cuNumberPtr = nil
	}

	reconList, err := reconObj.GetPavAliases(cuNumberPtr)
	if err != nil {
		return nil, err
	}
	if !hasBaseLdev && len(aliasLdevIDs) == 0 {
		return reconList, nil
	}

	if hasBaseLdev {
		baseLdevID := baseLdevIDRaw.(int)
		aliasSet := make(map[int]struct{}, len(aliasLdevIDs))
		for _, id := range aliasLdevIDs {
			aliasSet[id] = struct{}{}
		}

		result := make([]gatewaymodel.PavAlias, 0)
		for _, it := range *reconList {
			if it.PavAttribute != "ALIAS" {
				continue
			}
			baseMatch := (it.CBaseVolumeID != nil && *it.CBaseVolumeID == baseLdevID) || (it.SBaseVolumeID != nil && *it.SBaseVolumeID == baseLdevID)
			if !baseMatch {
				continue
			}
			if len(aliasSet) > 0 {
				if _, ok := aliasSet[it.LdevID]; !ok {
					continue
				}
			}
			result = append(result, it)
		}
		return &result, nil
	}

	byLdev := make(map[int]gatewaymodel.PavAlias, len(*reconList))
	for _, it := range *reconList {
		byLdev[it.LdevID] = it
	}

	result := make([]gatewaymodel.PavAlias, 0, len(aliasLdevIDs))
	for _, id := range aliasLdevIDs {
		if it, ok := byLdev[id]; ok {
			result = append(result, it)
		}
	}
	return &result, nil
}

func AssignPavAlias(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
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
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}

	baseLdevID := d.Get("base_ldev_id").(int)
	aliasLdevIDs := extractIntList(d, "alias_ldev_ids")
	return reconObj.AssignPavAlias(baseLdevID, aliasLdevIDs)
}

func UnassignPavAlias(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
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
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return err
	}

	aliasLdevIDs := extractIntList(d, "alias_ldev_ids")
	return reconObj.UnassignPavAlias(aliasLdevIDs)
}

func PavAliasItemsFromList(list *[]gatewaymodel.PavAlias) []map[string]interface{} {
	if list == nil {
		return nil
	}
	items := make([]map[string]interface{}, 0, len(*list))
	for _, it := range *list {
		m := map[string]interface{}{
			"cu_number":     it.CuNumber,
			"ldev_id":       it.LdevID,
			"pav_attribute": it.PavAttribute,
		}
		if it.CBaseVolumeID != nil {
			m["c_base_volume_id"] = *it.CBaseVolumeID
		}
		if it.SBaseVolumeID != nil {
			m["s_base_volume_id"] = *it.SBaseVolumeID
		}
		items = append(items, m)
	}
	return items
}

func extractIntList(d *schema.ResourceData, key string) []int {
	raw := d.Get(key).([]interface{})
	out := make([]int, 0, len(raw))
	for _, v := range raw {
		out = append(out, v.(int))
	}
	return out
}
