package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	// "time"

	"fmt"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetVssbComputeNodes(d *schema.ResourceData) (*[]terraformmodel.Server, error) {

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_BEGIN))
	reconServers, err := reconObj.GetAllComputeNodes()
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_SERVERS_FAILED))

		return nil, err
	}

	// Converting reconciler to terraform
	terraformComputeNodes := []terraformmodel.Server{}
	err = copier.Copy(&terraformComputeNodes, reconServers.Data)

	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_SERVERS_END))

	return &terraformComputeNodes, nil
}

func GetVssbComputeNode(d *schema.ResourceData, id string) (*terraformmodel.ComputeNodeWithPathDetails, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_BEGIN))
	reconServer, err := reconObj.GetComputeNode(id)
	if err != nil {
		log.WriteDebug("TFError| error getting GetAllStoragePools, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_SERVER_FAILED))
		return nil, err
	}

	// Converting reconciler to terraform
	terraformComputeNode := terraformmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&terraformComputeNode, reconServer)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	err = copier.Copy(&terraformComputeNode.ComputePaths, &reconServer.ComputePaths)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_SERVER_END))

	return &terraformComputeNode, nil
}

func DeleteVssbComputeNodeResource(d *schema.ResourceData) error {
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

	serverId := d.State().ID
	if serverId == "" {
		return fmt.Errorf("id field is required to delete the compute node resource.")
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_BEGIN), serverId)
	err = reconObj.DeleteComputeNodeResource(serverId)
	if err != nil {
		log.WriteDebug("TFError| error getting DeleteComputeNodeResource, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_SERVER_FAILED), serverId)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_SERVER_END), serverId)

	return nil
}

func ConvertVssbComputeNodeToSchema(computeNode *terraformmodel.Server) *map[string]interface{} {

	compNode := map[string]interface{}{
		"id":                computeNode.ID,
		"nickname":          computeNode.Nickname,
		"number_of_volumes": computeNode.NumberOfVolumes,
		"os_type":           computeNode.OsType,
		"total_capacity":    computeNode.TotalCapacity,
		"used_capacity":     computeNode.UsedCapacity,
		"number_of_paths":   computeNode.NumberOfPaths,
	}
	pa := []map[string]interface{}{}

	for _, item := range computeNode.Paths {
		data := map[string]interface{}{
			"hba_name": item.HbaName,
			"port_ids": item.PortIds,
		}
		pa = append(pa, data)
	}
	compNode["paths"] = pa

	return &compNode

}

func ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode *terraformmodel.ComputeNodeWithPathDetails) *map[string]interface{} {

	compNode := map[string]interface{}{
		"id":                computeNode.Node.ID,
		"nickname":          computeNode.Node.Nickname,
		"number_of_volumes": computeNode.Node.NumberOfVolumes,
		"os_type":           computeNode.Node.OsType,
		"total_capacity":    computeNode.Node.TotalCapacity,
		"used_capacity":     computeNode.Node.UsedCapacity,
		"number_of_paths":   computeNode.Node.NumberOfPaths,
	}
	pa := []map[string]interface{}{}

	for _, item := range computeNode.Node.Paths {
		data := map[string]interface{}{
			"hba_name": item.HbaName,
			"port_ids": item.PortIds,
			"protocol": item.Protocol,
		}
		pa = append(pa, data)
	}
	compNode["paths"] = pa

	pa2 := []map[string]interface{}{}

	for _, item := range computeNode.ComputePaths.Data {
		data := map[string]interface{}{
			"port_id":         item.PortId,
			"target_port_identifier": item.PortName,
			"port_name":       item.PortNickname,
		}
		pa2 = append(pa2, data)
	}
	compNode["port_details"] = pa2

	return &compNode

}

func CreateVssbComputeNode(d *schema.ResourceData) (*terraformmodel.ComputeNodeWithPathDetails, error) {
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
	createInput, err := CreateVssbComputeNodeFromSchema(d)
	if err != nil {
		return nil, err
	}

	protocolsToCheck := map[string]bool{}
	if len(createInput.FcConnections) > 0 {
		protocolsToCheck["FC"] = true
	}
	if len(createInput.IscsiConnections) > 0 {
		protocolsToCheck["iSCSI"] = true
	}
	if len(protocolsToCheck) > 0 {
		storagePorts, err := reconObj.GetStoragePorts()
		if err != nil {
			log.WriteDebug("TFError| error in GetStoragePorts, err: %v", err)
			return nil, err
		}
		availableProtocol := map[string]bool{}
		for _, port := range storagePorts.Data {
			availableProtocol[port.Protocol] = true
		}
		for requiredProtocol, _ := range protocolsToCheck {
			if !availableProtocol[requiredProtocol] {
				err := fmt.Errorf("no ports with protocol %s found — check your storage configuration", requiredProtocol)
				log.WriteDebug("TFError| error in Creating ComputeNode - ReconcileComputeNode , err: %v", err)
				return nil, err
			}
		}
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_COMPUTE_NODE_BEGIN), createInput.Name)
	reconcilerCreateComputeNode := reconcilermodel.ComputeResource{}
	err = copier.Copy(&reconcilerCreateComputeNode, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	compute, err := reconObj.ReconcileComputeNode(&reconcilerCreateComputeNode)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_COMPUTE_NODE_FAILED), createInput.Name)
		log.WriteDebug("TFError| error in Creating ComputeNode - ReconcileComputeNode , err: %v", err)
		return nil, err
	}

	terraformModelCompute := terraformmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&terraformModelCompute, compute)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_COMPUTE_NODE_END), createInput.Name)
	return &terraformModelCompute, nil
}

func UpdateVssbComputeNode(d *schema.ResourceData) (*terraformmodel.ComputeNodeWithPathDetails, error) {
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
	updateInput, err := CreateVssbComputeNodeFromSchema(d)
	if err != nil {
		return nil, err
	}

	protocolsToCheck := map[string]bool{}
	if len(updateInput.FcConnections) > 0 {
		protocolsToCheck["FC"] = true
	}
	if len(updateInput.IscsiConnections) > 0 {
		protocolsToCheck["iSCSI"] = true
	}
	if len(protocolsToCheck) > 0 {
		storagePorts, err := reconObj.GetStoragePorts()
		if err != nil {
			log.WriteDebug("TFError| error in GetStoragePorts, err: %v", err)
			return nil, err
		}
		availableProtocol := map[string]bool{}
		for _, port := range storagePorts.Data {
			availableProtocol[port.Protocol] = true
		}
		for requiredProtocol, _ := range protocolsToCheck {
			if !availableProtocol[requiredProtocol] {
				err := fmt.Errorf("no ports with protocol %s found — check your storage configuration", requiredProtocol)
				log.WriteDebug("TFError| error in Creating ComputeNode - ReconcileComputeNode , err: %v", err)
				return nil, err
			}
		}
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_COMPUTE_NODE_BEGIN), updateInput.Name)
	reconcilerCreateComputeNode := reconcilermodel.ComputeResource{}
	err = copier.Copy(&reconcilerCreateComputeNode, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	// Add ID for Compute Node
	serverId := d.State().ID
	if serverId == "" {
		return nil, fmt.Errorf("unable to find compute node id, update operation failed")
	}
	reconcilerCreateComputeNode.ID = serverId

	compute, err := reconObj.ReconcileComputeNode(&reconcilerCreateComputeNode)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_COMPUTE_NODE_FAILED), updateInput.Name)
		log.WriteDebug("TFError| error in Updating ComputeNode - ReconcileComputeNode , err: %v", err)
		return nil, err
	}

	terraformModelCompute := terraformmodel.ComputeNodeWithPathDetails{}
	err = copier.Copy(&terraformModelCompute, compute)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_COMPUTE_NODE_END), updateInput.Name)
	return &terraformModelCompute, nil
}

func CreateVssbComputeNodeFromSchema(d *schema.ResourceData) (*terraformmodel.ComputeResource, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.ComputeResource{}
	name, ok := d.GetOk("compute_node_name")
	if ok {
		pid := name.(string)
		createInput.Name = pid
	}
	osType, ok := d.GetOk("os_type")
	if ok {
		pid := osType.(string)
		createInput.OsType = pid
	}

	connection, ok := d.GetOk("iscsi_connection")
	if ok {
		iscsiConn := connection.(*schema.Set).List()
		iscOutput := []terraformmodel.IscsiConnector{}
		for _, conn := range iscsiConn {
			v := conn.(map[string]interface{})
			port := v["port_names"].([]interface{})
			ports := make([]string, len(port))
			for index, value := range port {
				switch typedValue := value.(type) {
				case string:
					ports[index] = typedValue
				}
			}
			isc := terraformmodel.IscsiConnector{
				IscsiInitiator: v["iscsi_initiator"].(string),
				PortNames:      ports,
			}
			iscOutput = append(iscOutput, isc)
		}
		createInput.IscsiConnections = iscOutput
	}

	fcConnection, ok := d.GetOk("fc_connection")
	if ok {
		fcConn := fcConnection.(*schema.Set).List()
		fccOutput := []terraformmodel.FcConnector{}
		for _, conn := range fcConn {
			v := conn.(map[string]interface{})
			isc := terraformmodel.FcConnector{
				HostWWN: v["host_wwn"].(string),
			}
			fccOutput = append(fccOutput, isc)
		}
		createInput.FcConnections = fccOutput
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}
