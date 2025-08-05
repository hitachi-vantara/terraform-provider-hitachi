package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	"errors"
	"fmt"

	// "io/ioutil"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetVssbVolumes(d *schema.ResourceData) (*[]terraformmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
	computeNodeName := d.Get("compute_node_name").(string)
	reconStoragePools, err := reconObj.GetAllVolumes(computeNodeName)
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformVolumes := []terraformmodel.Volume{}
	err = copier.Copy(&terraformVolumes, reconStoragePools.Data)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))

	return &terraformVolumes, nil
}

func GetVssbVolumeNode(d *schema.ResourceData) (*terraformmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_BEGIN))
	var volumeName string = ""
	volName, ok := d.GetOk("volume_name")
	if !ok {
		volumeName = d.Get("name").(string)
	} else {
		volumeName = volName.(string)
	}
	if volumeName == "" {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, errors.New("either volume_name or name parameter is required")
	}
	reconStoragePools, err := reconObj.GetVolumeDetails(volumeName)
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_VOLUME_INFO_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformVolumeNodes := terraformmodel.Volume{}
	err = copier.Copy(&terraformVolumeNodes, reconStoragePools)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_VOLUME_INFO_END))

	return &terraformVolumeNodes, nil
}

func ConvertVssbVolumesToSchema(volumes *terraformmodel.Volume) *map[string]interface{} {
	vol := map[string]interface{}{
		"id":                           volumes.ID,
		"name":                         volumes.Name,
		"nick_name":                    volumes.NickName,
		"volume_number":                volumes.VolumeNumber,
		"pool_id":                      volumes.PoolId,
		"pool_name":                    volumes.PoolName,
		"total_capacity":               volumes.TotalCapacity,
		"used_capacity":                volumes.UsedCapacity,
		"number_of_connecting_servers": volumes.NumberOfConnectingServers,
		"number_of_snapshots":          volumes.NumberOfSnapshots,
		"protection_domain_id":         volumes.ProtectionDomainId,
		"full_allocated":               volumes.FullAllocated,
		"volume_type":                  volumes.VolumeType,
		"status_summary":               volumes.StatusSummary,
		"status":                       volumes.Status,
		"storage_controller_id":        volumes.StorageControllerId,
		"snapshot_attribute":           volumes.SnapshotAttribute,
		"snapshot_status":              volumes.SnapshotStatus,
		"saving_setting":               volumes.SavingSetting,
		"saving_mode":                  volumes.SavingMode,
		"data_reduction_status":        volumes.DataReductionStatus,
		"data_reduction_progress_rate": volumes.DataReductionProgressRate,
	}
	se := []map[string]interface{}{}
	data := map[string]interface{}{
		"system_data_capacity":                            volumes.SavingEffect.DataCapacity,
		"pre_capacity_data_reduction_without_system_data": volumes.SavingEffect.PreCapacity,
		"post_capacity_data_reduction":                    volumes.SavingEffect.PostCapacity,
	}
	se = append(se, data)
	vol["saving_effects"] = se
	// compute_nodes := ConvertVssbVolumeNodesToSchema(volumes.ComputeNodes)
	if volumes.ComputeNodes != nil {
		vol["compute_nodes"] = ConvertVssbVolumeComputeNodesToSchema(&volumes.ComputeNodes)
	}

	return &vol
}

func ConvertVssbVolumeComputeNodesToSchema(nodes *[]terraformmodel.VolumeNode) *[]map[string]interface{} {
	nodeSchema := []map[string]interface{}{}

	for _, node := range *nodes {
		Schema := map[string]interface{}{
			"id":             node.ID,
			"name":           node.Nickname,
			"os_type":        node.OsType,
			"total_capacity": node.TotalCapacity,
			"used_capacity":  node.UsedCapacity,
			"volume_count":   node.NumberOfVolumes,
			"lun":            node.Lun,
		}
		nodeSchema = append(nodeSchema, Schema)
	}

	return &nodeSchema
}

func CreateVolume(d *schema.ResourceData) (*terraformmodel.Volume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

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

	volumeSchema, err := CreateVolumeReqFromSchema(d)
	if err != nil {
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_BEGIN), volumeSchema.Name)
	terraformCreateVolume := reconcilermodel.CreateVolume{}
	err = copier.Copy(&terraformCreateVolume, volumeSchema)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from terraform to to reconcile structure, err: %v", err)
		return nil, err
	}

	reconVolumeResponse, err := reconObj.ReconcileVolume(&terraformCreateVolume)
	if err != nil {
		log.WriteDebug("TFError| error creating volume, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_VOLUME_FAILED), volumeSchema.Name)
		return nil, err
	}
	log.WriteDebug("Terraform Volume created/updated: %#v", *reconVolumeResponse)

	// Converting reconciler to terraform
	terraformVolumeRes := terraformmodel.Volume{}
	err = copier.Copy(&terraformVolumeRes, reconVolumeResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_VOLUME_END), volumeSchema.Name, terraformVolumeRes.ID)

	return &terraformVolumeRes, nil
}

func DeleteVolume(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vosb_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Username:       storageSetting.Username,
		Password:       storageSetting.Password,
		ClusterAddress: storageSetting.ClusterAddress,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}
	name, ok := d.GetOk("name")
	if !ok {
		return fmt.Errorf("name is the mandatory field")
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_BEGIN), name)

	id := d.State().ID
	err = reconObj.DeleteVolumeResource(&id)
	if err != nil {

		log.WriteDebug("TFError| error deleting volume, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_VOLUME_FAILED), name)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_VOLUME_END), name)

	return nil
}

func CreateVolumeReqFromSchema(d *schema.ResourceData) (*terraformmodel.CreateVolume, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateVolume{}

	name, ok := d.GetOk("name")
	if !ok {
		return nil, fmt.Errorf("name is the mandatory field")
	} else {
		createInput.Name = name.(string)
	}
	storagePool, ok := d.GetOk("storage_pool")
	if ok {
		createInput.PoolName = storagePool.(string)
	}

	capacity, ok := d.GetOk("capacity_gb")
	if ok {

		createInput.CapacityInGB = float32(capacity.(float64))
	}
	nick_name, ok := d.GetOk("nick_name")
	if ok {
		createInput.NickName = nick_name.(string)
	}
	computeNodes := d.Get("compute_nodes")
	computeNodeCheck := d.GetRawConfig().GetAttr("compute_nodes").IsNull()
	if !computeNodeCheck {
		nodes := []string{}
		for _, node := range computeNodes.([]interface{}) {
			nodes = append(nodes, node.(string))
		}
		createInput.ComputeNodes = nodes
	} else {
		createInput.ComputeNodes = nil
	}

	return &createInput, nil
}
