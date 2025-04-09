package terraform

import (
	// "encoding/json"
	// "context"
	//"fmt"
	// "io/ioutil"
	// "time"

	"math"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetDashboardInfo(d *schema.ResourceData) (*terraformmodel.Dashboard, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DASHBOARD_BEGIN))

	reconDashboard, err := reconObj.GetDashboardInfo()

	if err != nil {
		log.WriteDebug("TFError| error getting GetDashboardInfo, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_DASHBOARD_FAILED))

		return nil, err
	}

	// Converting reconciler to terraform
	terraformDashboard := terraformmodel.Dashboard{}
	err = copier.Copy(&terraformDashboard, reconDashboard)

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_DASHBOARD_END))

	return &terraformDashboard, nil
}

func convertMbToGb(mb uint64) float64 {
	gb := math.Round(float64(mb) / float64(1024))
	return gb
}

func ConvertVssbDashboardToSchema(dashboard *terraformmodel.Dashboard) *map[string]interface{} {

	db := map[string]interface{}{
		"fault_domain_count": dashboard.NumberOfFaultDomains,
		"volume_count":       dashboard.NumberOfTotalVolumes,
		"compute_node_count": dashboard.NumberOfTotalServers,
		"compute_port_count": dashboard.NumberOfComputePorts,
		"drive_count":        dashboard.NumberOfDrives,
		"storage_node_count": dashboard.NumberOfTotalStorageNodes,
		"storage_pool_count": dashboard.NumberOfStoragePools,
		"total_capacity_mb":  dashboard.TotalPoolCapacityInMB,
		"used_capacity_mb":   dashboard.UsedPoolCapacityInMB,
		"free_capacity_mb":   dashboard.FreePoolCapacityInMB,
		"total_efficiency":   dashboard.TotalEfficiency,
		"data_reduction":     dashboard.EfficiencyDataReduction,
		"total_capacity_gb":  convertMbToGb(dashboard.TotalPoolCapacityInMB),
		"used_capacity_gb":   convertMbToGb(dashboard.UsedPoolCapacityInMB),
		"free_capacity_gb":   convertMbToGb(dashboard.FreePoolCapacityInMB),
	}
	hs := []map[string]interface{}{}
	for _, item := range dashboard.HealthStatuses {
		data := map[string]interface{}{
			"type":   item.Type,
			"status": item.Status,
		}
		hs = append(hs, data)
	}
	db["health_status"] = hs

	return &db
}
