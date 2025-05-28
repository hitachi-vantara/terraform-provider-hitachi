package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vosb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetStorageVersionInfo gets version information of vssb storage
func (psm *vssbStorageManager) GetStorageVersionInfo() (*vssbmodel.StorageVersionInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	versionInfo, err := gatewayObj.GetStorageVersionInfo()
	if err != nil {
		return nil, err
	}

	provVersionInfo := vssbmodel.StorageVersionInfo{}
	err = copier.Copy(&provVersionInfo, versionInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provVersionInfo, nil
}

// GetHealthStatuses Obtains the health status.
func (psm *vssbStorageManager) GetHealthStatuses() (*vssbmodel.HealthStatuses, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	healthStatuses, err := gatewayObj.GetHealthStatuses()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetHealthStatuses, err: %+v", err)
		return nil, err
	}

	provHealthStatuses := vssbmodel.HealthStatuses{}
	err = copier.Copy(&provHealthStatuses, healthStatuses)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provHealthStatuses, nil
}

// GetStorageClusterInfo Obtains the storage cluster information.
func (psm *vssbStorageManager) GetStorageClusterInfo() (*vssbmodel.StorageClusterInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	storageClusterInfo, err := gatewayObj.GetStorageClusterInfo()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStorageClusterInfo, err: %+v", err)
		return nil, err
	}

	provStorageClusterInfo := vssbmodel.StorageClusterInfo{}
	err = copier.Copy(&provStorageClusterInfo, storageClusterInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provStorageClusterInfo, nil
}

// GetDrivesInfo Obtains a list of drive information.
func (psm *vssbStorageManager) GetDrivesInfo() (*vssbmodel.Drives, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := vssbgatewaymodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	drivesInfo, err := gatewayObj.GetDrivesInfo()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStorageClusterInfo, err: %+v", err)
		return nil, err
	}

	provDrivesInfo := vssbmodel.Drives{}
	err = copier.Copy(&provDrivesInfo, drivesInfo)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}

	return &provDrivesInfo, nil
}
