package vssbstorage

import (
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
)

// GetAllStoragePools gets all storage pool details
func (psm *vssbStorageManager) GetAllStoragePools() (*vssbmodel.StoragePools, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_BEGIN))
	provStoragePools, err := provObj.GetAllStoragePools()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllComputeNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_POOLS_FAILED))
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePools := vssbmodel.StoragePools{}
	err = copier.Copy(&reconStoragePools, provStoragePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_POOLS_END))

	return &reconStoragePools, nil
}

// GetStoragePoolsByPoolNames gets storage pools by pool names
func (psm *vssbStorageManager) GetStoragePoolsByPoolNames(poolNames []string) (*vssbmodel.StoragePools, error) {
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

	names := strings.Join(poolNames, ",")

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_BEGIN), names)
	provStoragePools, err := provObj.GetStoragePoolsByPoolNames(poolNames)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllComputeNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_POOL_FAILED), names)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePools := vssbmodel.StoragePools{}
	err = copier.Copy(&reconStoragePools, provStoragePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_END), names)

	return &reconStoragePools, nil
}

func (psm *vssbStorageManager) GetStoragePoolByPoolName(poolName string) (*vssbmodel.StoragePool, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_BEGIN), poolName)
	provStoragePools, err := provObj.GetStoragePoolByPoolName(poolName)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllComputeNodes provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_STORAGE_POOL_FAILED), poolName)
		return nil, err
	}

	// Converting Prov to Reconciler
	reconStoragePools := vssbmodel.StoragePool{}
	err = copier.Copy(&reconStoragePools, provStoragePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from prov to reconciler structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_STORAGE_POOL_END), poolName)

	return &reconStoragePools, nil
}
