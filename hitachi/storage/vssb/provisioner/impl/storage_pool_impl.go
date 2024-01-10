package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/vssb/gateway/impl"
	vssbgatewaymodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/message-catelog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/provisioner/model"

	"github.com/jinzhu/copier"
)

// GetAllStoragePools gets all storage pool details
func (psm *vssbStorageManager) GetAllStoragePools() (*vssbmodel.StoragePools, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_BEGIN))
	storagePools, err := gatewayObj.GetAllStoragePools()
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllStoragePools, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_POOLS_FAILED))
		return nil, err
	}

	provStoragePools := vssbmodel.StoragePools{}
	err = copier.Copy(&provStoragePools, storagePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_END))

	return &provStoragePools, nil
}

// GetStoragePoolsByPoolNames gets storage pools by pool names
func (psm *vssbStorageManager) GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error) {
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

	names := strings.Join(poolNames, ",")

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_BEGIN), names)
	storagePools, err := gatewayObj.GetStoragePoolsByPoolNames(poolNames)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetAllStoragePools, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_POOL_FAILED), names)
		return nil, err
	}

	provStoragePools := vssbmodel.StoragePools{}
	err = copier.Copy(&provStoragePools, storagePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_END), names)

	return &provStoragePools, nil
}

// GetStoragePoolByPoolName gets storage pool by pool name
func (psm *vssbStorageManager) GetStoragePoolByPoolName(poolName string) (*vssbmodel.StoragePool, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_BEGIN), poolName)
	poolsName := []string{poolName}
	storagePool, err := gatewayObj.GetStoragePoolsByPoolNames(poolsName)
	if err != nil {
		log.WriteDebug("TFError| failed to call GetStoragePoolByPoolName, err: %+v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_POOL_FAILED), poolName)
		return nil, err
	}
	provStoragePool := vssbmodel.StoragePool{}
	var found bool = false
	for _, v := range storagePool.Data {
		if v.Name == poolName {
			found = true
			err = copier.Copy(&provStoragePool, v)
			if err != nil {
				log.WriteDebug("TFError| error in Copy from gateway to provisioner structure, err: %v", err)
				return nil, err
			}

		}
	}
	if !found {
		log.WriteDebug("TFError| no strorage pool found for the given pool name %s", poolName)
		return nil, fmt.Errorf("no strorage pool found for the given pool name %s", poolName)
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_END), poolName)

	return &provStoragePool, nil
}
