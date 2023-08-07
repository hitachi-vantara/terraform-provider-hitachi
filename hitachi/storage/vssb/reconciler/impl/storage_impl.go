package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetStorageVersionInfo gets version information of vssb storage
func (psm *vssbStorageManager) GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	provVersionInfo, err := provObj.GetStorageVersionInfo()
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageVersionInfo provisioner call, err: %v", err)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconVersionInfo := vssbmodel.StorageVersionInfo{}
	err = copier.Copy(&reconVersionInfo, provVersionInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}

	return &reconVersionInfo, nil
}

// GetDashboardInfo gets the dashboard information of vssb storage
func (psm *vssbStorageManager) GetDashboardInfo() (*vssbmodel.Dashboard, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	provHealths, err := provObj.GetHealthStatuses()
	if err != nil {
		log.WriteDebug("TFError| error in GetHealthStatuses provisioner call, err: %v", err)
		return nil, err
	}

	provClusterInfo, err := provObj.GetStorageClusterInfo()
	if err != nil {
		log.WriteDebug("TFError| error in GetStorageClusterInfo provisioner call, err: %v", err)
		return nil, err
	}

	provDrivesInfo, err := provObj.GetDrivesInfo()
	if err != nil {
		log.WriteDebug("TFError| error in GetDrivesInfo provisioner call, err: %v", err)
		return nil, err
	}

	provPortsInfo, err := provObj.GetStoragePorts()
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragePorts provisioner call, err: %v", err)
		return nil, err
	}

	provStoragePoolsInfo, err := provObj.GetAllStoragePools()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllStoragePools provisioner call, err: %v", err)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconDashboardInfo, err := ConvertToDashBoard(provHealths, provClusterInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	reconDashboardInfo.NumberOfDrives = len(provDrivesInfo.Data)
	reconDashboardInfo.NumberOfComputePorts = len(provPortsInfo.Data)
	reconDashboardInfo.NumberOfStoragePools = len(provStoragePoolsInfo.Data)

	return reconDashboardInfo, nil
}

func ConvertToDashBoard(healths *provisonermodel.HealthStatuses, clusterInfo *provisonermodel.StorageClusterInfo) (*vssbmodel.Dashboard, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	data := *healths
	list := data["resources"]
	dashboardInfo := vssbmodel.Dashboard{}

	itList := []vssbmodel.HealthStatus{}
	for _, hs := range list {
		vhs := vssbmodel.HealthStatus{}
		vhs.Type = hs.Type
		vhs.Status = hs.Status
		itList = append(itList, vhs)
	}
	dashboardInfo.HealthStatuses = itList
	dashboardInfo.NumberOfTotalVolumes = clusterInfo.NumberOfTotalVolumes
	dashboardInfo.NumberOfTotalServers = clusterInfo.NumberOfTotalServers
	dashboardInfo.NumberOfTotalStorageNodes = clusterInfo.NumberOfReadyStorageNodes
	dashboardInfo.NumberOfFaultDomains = clusterInfo.NumberOfFaultDomains
	dashboardInfo.TotalPoolCapacityInMB = clusterInfo.TotalPoolCapacityInMB
	dashboardInfo.UsedPoolCapacityInMB = clusterInfo.UsedPoolCapacityInMB
	dashboardInfo.FreePoolCapacityInMB = clusterInfo.FreePoolCapacityInMB
	dashboardInfo.TotalEfficiency = clusterInfo.SavingEffects.TotalEfficiency
	dashboardInfo.EfficiencyDataReduction = clusterInfo.SavingEffects.EfficiencyDataReduction

	return &dashboardInfo, nil
}
