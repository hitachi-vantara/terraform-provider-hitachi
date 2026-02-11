package terraform

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	// "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetVssbNode(d *schema.ResourceData) (*terraformmodel.StorageNodeVssb, error) {
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

	terraformStorageNode := terraformmodel.StorageNodeVssb{}

	if  d.Get("node_name") == nil {
		return &terraformStorageNode, nil
	}

	nodeName := d.Get("node_name").(string)
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_BEGIN), nodeName)

	reconStorageNode, err := reconObj.GetStorageNode(nodeName)
	if err != nil {
		log.WriteDebug("TFError| error getting GetStorageNode, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_NODE_FAILED), nodeName)
		return nil, err
	}

	// Converting reconciler to terraform
	err = copier.Copy(&terraformStorageNode, reconStorageNode)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NODE_END), nodeName)

	return &terraformStorageNode, nil
}

func GetVssbStorageNodes(d *schema.ResourceData) (*[]terraformmodel.StorageNodeVssb, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_BEGIN))

	reconStorageNodes, err := reconObj.GetStorageNodes()
	if err != nil {
		log.WriteDebug("TFError| error getting GetStorageNodes, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_STORAGE_NODES_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformStorageNodes := terraformmodel.StorageNodesVssb{}
	err = copier.Copy(&terraformStorageNodes, reconStorageNodes)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_STORAGE_NODES_END))

	return &terraformStorageNodes.Data, nil
}

func ConvertVssbStorageNodeToSchema(storageNode *terraformmodel.StorageNodeVssb) *map[string]interface{} {
	sp := map[string]interface{}{
		"id":                           storageNode.ID,
		"name":                         storageNode.Name,
		"bios_uuid":                    storageNode.BiosUuid,
		"fault_domain_id":              storageNode.FaultDomainID,
		"fault_domain_name":            storageNode.FaultDomainName,
		"cluster_role":                 storageNode.ClusterRole,
		"control_port_ipv4_address":    storageNode.ControlPortIpv4Address,
		"internode_port_ipv4_address":  storageNode.InternodePortIpv4Address,
		"software_version":             storageNode.SoftwareVersion,
		"model_name":                   storageNode.ModelName,
		"serial_number":                storageNode.SerialNumber,
		"memory":                       storageNode.Memory,
		"availability_zone_id":         storageNode.AvailabilityZoneID,
		"protection_domain_id":         storageNode.ProtectionDomainID,
		"drive_data_relocation_status": storageNode.DriveDataRelocationStatus,
		"status_summary":               storageNode.StatusSummary,
		"status":                       storageNode.Status,
	}

	item1 := []map[string]interface{}{}
	item1 = append(item1, map[string]interface{}{
		"capacity_of_drive": storageNode.InsufficientResourcesForRebuildCapacity.CapacityOfDrive,
		"number_of_drives":  storageNode.InsufficientResourcesForRebuildCapacity.NumberOfDrives,
	})
	sp["insufficient_resources_for_rebuild_capacity"] = item1

	// elements := []int{1,2}
	// setValue, _ := types.SetValueFrom(context.Background(), types.Int16, elements)
	// sp["insufficient_resources_for_rebuild_capacity"] = setValue

	// sp["insufficient_resources_for_rebuild_capacity"] = map[string]interface{}{
	// 	"capacity_of_drive": storageNode.InsufficientResourcesForRebuildCapacity.CapacityOfDrive,
	// 	"number_of_drives":  storageNode.InsufficientResourcesForRebuildCapacity.NumberOfDrives,
	// }

	item2 := []map[string]interface{}{}
	item2 = append(item2, map[string]interface{}{
		"number_of_drives": storageNode.InsufficientResourcesForRebuildCapacity.NumberOfDrives,
	})
	sp["rebuildable_resources"] = item2

	return &sp
}

// Add Storage Node
func CreateVssbStorageNode(d *schema.ResourceData) error {
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

	// FIXME: do we need to check if the node is already added, it is handled by the backend
	// storageNodeName := d.Get("node_name").(string)
	// if storageNodeName == "" {
	// 	log.WriteDebug("TFError| node_name cannot be empty")
	// 	return nil
	// }

	// see validateStorageNodeParameters for validate parameter combinations
	configurationFile := d.Get("configuration_file").(string)
	exportedConfigurationFile := d.Get("exported_configuration_file").(string)
	setupUserPassword := d.Get("setup_user_password").(string)
	expectedCloudProvider := d.Get("expected_cloud_provider").(string)
	vmConfigFileS3URI := d.Get("vm_configuration_file_s3_uri").(string)

	err = reconObj.AddStorageNode(
		configurationFile,
		exportedConfigurationFile,
		setupUserPassword,
		expectedCloudProvider,
		vmConfigFileS3URI,
	)
	if err != nil {
		log.WriteDebug("TFError| error in GetStoragenodes, err: %v", err)
		return err
	}

	return nil
}
