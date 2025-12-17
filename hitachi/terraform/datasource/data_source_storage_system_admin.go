package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageSystemAdmin() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage System VSP One: It returns the storage device related information.",
		ReadContext: dataSourceStorageSystemAdminRead,
		Schema:      schemaimpl.StorageSystemAdminSchema,
	}
}

func dataSourceStorageSystemAdminRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ss, err := impl.GetStorageSystemAdmin(d)
	if err != nil {
		return diag.FromErr(err)
	}
	savingEffects := map[string]interface{}{
		"efficiency_data_reduction":                        ss.SavingEffects.EfficiencyDataReduction,
		"pre_capacity_data_reduction":                      ss.SavingEffects.PreCapacityDataReduction,
		"post_capacity_data_reduction":                     ss.SavingEffects.PostCapacityDataReduction,
		"efficiency_fmd_saving":                            ss.SavingEffects.EfficiencyFmdSaving,
		"pre_capacity_fmd_saving":                          ss.SavingEffects.PreCapacityFmdSaving,
		"post_capacity_fmd_saving":                         ss.SavingEffects.PostCapacityFmdSaving,
		"is_total_efficiency_support":                      ss.SavingEffects.IsTotalEfficiencySupport,
		"total_efficiency_status":                          ss.SavingEffects.TotalEfficiencyStatus,
		"data_reduction_without_system_data_status":        ss.SavingEffects.DataReductionWithoutSystemDataStatus,
		"software_saving_without_system_data_status":       ss.SavingEffects.SoftwareSavingWithoutSystemDataStatus,
		"total_efficiency":                                 ss.SavingEffects.TotalEfficiency,
		"data_reduction_without_system_data":               ss.SavingEffects.DataReductionWithoutSystemData,
		"pre_capacity_data_reduction_without_system_data":  ss.SavingEffects.PreCapacityDataReductionWithoutSystemData,
		"post_capacity_data_reduction_without_system_data": ss.SavingEffects.PostCapacityDataReductionWithoutSystemData,
		"software_saving_without_system_data":              ss.SavingEffects.SoftwareSavingWithoutSystemData,
		"calculation_start_time":                           ss.SavingEffects.CalculationStartTime,
		"calculation_end_time":                             ss.SavingEffects.CalculationEndTime,
	}

	ssadmin := map[string]interface{}{
		"model_name":                                ss.ModelName,
		"serial":                                    ss.Serial,
		"nickname":                                  ss.Nickname,
		"number_of_total_volumes":                   ss.NumberOfTotalVolumes,
		"number_of_free_drives":                     ss.NumberOfFreeDrives,
		"number_of_total_servers":                   ss.NumberOfTotalServers,
		"total_physical_capacity":                   ss.TotalPhysicalCapacity,
		"total_pool_capacity":                       ss.TotalPoolCapacity,
		"total_pool_physical_capacity":              ss.TotalPoolPhysicalCapacity,
		"used_pool_capacity":                        ss.UsedPoolCapacity,
		"free_pool_capacity":                        ss.FreePoolCapacity,
		"total_pool_capacity_with_ti_pool":          ss.TotalPoolCapacityWithTiPool,
		"total_pool_physical_capacity_with_ti_pool": ss.TotalPoolPhysicalCapacityWithTiPool,
		"used_pool_capacity_with_ti_pool":           ss.UsedPoolCapacityWithTiPool,
		"free_pool_capacity_with_ti_pool":           ss.FreePoolCapacityWithTiPool,
		"estimated_configurable_pool_capacity":      ss.EstimatedConfigurablePoolCapacity,
		"estimated_configurable_volume_capacity":    ss.EstimatedConfigurableVolumeCapacity,
		"gum_version":                               ss.GumVersion,
		"esm_os_version":                            ss.EsmOsVersion,
		"dkc_micro_version":                         ss.DkcMicroVersion,
		"warning_led_status":                        ss.WarningLedStatus,
		"esm_status":                                ss.EsmStatus,
		"ip_address_ipv4_service":                   ss.IpAddressIpv4Service,
		"ip_address_ipv4_ctl1":                      ss.IpAddressIpv4Ctl1,
		"ip_address_ipv4_ctl2":                      ss.IpAddressIpv4Ctl2,
		"ip_address_ipv6_service":                   ss.IpAddressIpv6Service,
		"ip_address_ipv6_ctl1":                      ss.IpAddressIpv6Ctl1,
		"ip_address_ipv6_ctl2":                      ss.IpAddressIpv6Ctl2,
		"saving_effects":                            []map[string]interface{}{savingEffects},
	}

	if !d.Get("with_estimated_configurable_capacities").(bool) {
		ssadmin["estimated_configurable_pool_capacity"] = nil
		ssadmin["estimated_configurable_volume_capacity"] = nil
	}

	log.WriteDebug("ssadmin: %+v\n", ssadmin)

	if err := d.Set("storage_system_admin", []map[string]interface{}{ssadmin}); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
