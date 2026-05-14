package terraform

import (
	"context"
	"fmt"
	"strings"

	// "fmt"

	// "time"
	// "errors"
	"sync"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncHComputeNodeOperation = &sync.Mutex{}

func ResourceVssbStorageComputeNode() *schema.Resource {
	return &schema.Resource{
		Description: "VSP One SDS Block Compute Node: Registers the information of the compute node.",
		Importer: &schema.ResourceImporter{
			StateContext: importVosbComputeNode,
		},
		CreateContext: resourceVssbStorageComputeNodeCreate,
		ReadContext:   resourceVssbStorageComputeNodeRead,
		UpdateContext: resourceVssbStorageComputeNodeUpdate,
		DeleteContext: resourceVssbStorageComputeNodeDelete,
		Schema:        schemaimpl.ResourceVssbStorageComputeNodeSchema,
		CustomizeDiff: resourceComputeNodeCustomDiff,
	}
}

func validateVssbComputeNodeConnections(d *schema.ResourceData) error {
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
		return err
	}

	connection, ok := d.GetOk("iscsi_connection")
	if !ok {
		return nil
	}

	iscsiConn := connection.(*schema.Set).List()
	storagePorts, err := reconObj.GetStoragePorts()
	if err != nil {
		return err
	}

	for _, conn := range iscsiConn {
		v := conn.(map[string]interface{})
		iqnName := v["iscsi_initiator"].(string)
		if iqnName != "" && !utils.IsIqn(iqnName) {
			return fmt.Errorf("iscsi_initiator %s is invalid", iqnName)
		}

		portNames, _ := v["port_names"].([]interface{})
		for _, value := range portNames {
			portName, ok := value.(string)
			if !ok {
				continue
			}
			portFound := false
			if storagePorts.Data != nil {
				for _, port := range storagePorts.Data {
					if port.Nickname == portName {
						portFound = true
						break
					}
				}
			}
			if !portFound {
				return fmt.Errorf("port name %s is invalid", portName)
			}
		}
	}

	return nil
}

func resourceComputeNodeCustomDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Only do deterministic schema validations here (no API calls).
	if connection, ok := d.GetOk("iscsi_connection"); ok {
		iscsiConn := connection.(*schema.Set).List()
		for _, conn := range iscsiConn {
			v := conn.(map[string]interface{})
			iqnName := v["iscsi_initiator"].(string)
			if iqnName != "" && !utils.IsIqn(iqnName) {
				log.WriteDebug("TFDebug | iqnName: %s", iqnName)
				return fmt.Errorf("iscsi_initiator %s is invalid", iqnName)
			}
		}
	}

	// fix for ResourceVssbStorageComputeNodeSchema 'compute_nodes' not updated in console output
	if err := d.SetNewComputed("compute_nodes"); err != nil {
		return err
	}

	return nil
}

func resourceVssbStorageComputeNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting compute node resource delete")

	err := impl.DeleteVssbComputeNodeResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.WriteInfo("compute node resource deleted successfully")
	return nil
}

func resourceVssbStorageComputeNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncHComputeNodeOperation.Lock()
	defer syncHComputeNodeOperation.Unlock()

	vssbAddr, _ := d.Get("vosb_address").(string)
	if strings.TrimSpace(vssbAddr) == "" {
		return diag.FromErr(fmt.Errorf("vosb_address is required to create a compute node"))
	}

	if strings.TrimSpace(d.Get("compute_node_name").(string)) == "" {
		return diag.FromErr(fmt.Errorf("compute_node_name is required to create a compute node"))
	}

	log.WriteInfo("starting compute node creation")
	if err := validateVssbComputeNodeConnections(d); err != nil {
		return diag.FromErr(err)
	}
	computeNode, err := impl.CreateVssbComputeNode(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cpn := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)
	log.WriteDebug("cpn: %+v\n", *cpn)
	cpnList := []map[string]interface{}{
		*cpn,
	}
	if err := d.Set("compute_nodes", cpnList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("compute_node_name", computeNode.Node.Nickname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("os_type", computeNode.Node.OsType); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(computeNode.Node.ID)
	log.WriteInfo("compute node created successfully")
	return nil
}

func resourceVssbStorageComputeNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbComputeNodesRead(ctx, d, m)
}

func resourceVssbStorageComputeNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncHComputeNodeOperation.Lock()
	defer syncHComputeNodeOperation.Unlock()

	log.WriteInfo("starting compute node update")
	vssbAddr, _ := d.Get("vosb_address").(string)
	if strings.TrimSpace(vssbAddr) == "" {
		return diag.FromErr(fmt.Errorf("vosb_address is required to update a compute node"))
	}
	if strings.TrimSpace(d.Get("compute_node_name").(string)) == "" {
		return diag.FromErr(fmt.Errorf("compute_node_name is required to update a compute node"))
	}
	if err := validateVssbComputeNodeConnections(d); err != nil {
		return diag.FromErr(err)
	}
	computeNode, err := impl.UpdateVssbComputeNode(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cpn := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)
	log.WriteDebug("cpn: %+v\n", *cpn)
	cpnList := []map[string]interface{}{
		*cpn,
	}
	if err := d.Set("compute_nodes", cpnList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("compute_node_name", computeNode.Node.Nickname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("os_type", computeNode.Node.OsType); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(computeNode.Node.ID)
	log.WriteInfo("compute node updated successfully")
	return nil
}
