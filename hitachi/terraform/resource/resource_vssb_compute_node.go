package terraform

import (
	"context"
	"fmt"

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
	reconimpl "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncHComputeNodeOperation = &sync.Mutex{}

func ResourceVssbStorageComputeNode() *schema.Resource {
	return &schema.Resource{
		Description:   ":meta:subcategory:VSS Block Compute Node:Registers the information of the compute node.",
		CreateContext: resourceVssbStorageComputeNodeCreate,
		ReadContext:   resourceVssbStorageComputeNodeRead,
		UpdateContext: resourceVssbStorageComputeNodeUpdate,
		DeleteContext: resourceVssbStorageComputeNodeDelete,
		Schema:        schemaimpl.ResourceVssbStorageComputeNodeSchema,
		CustomizeDiff: resourceComputeNodeCustomDiff,
	}
}

func resourceComputeNodeCustomDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	vssbAddr := d.Get("vss_block_address").(string)

	storageSetting, err := cache.GetVssbSettingsFromCache(vssbAddr)
	if err != nil {
		return err
	}

	// Local Check
	name, ok := d.GetOk("target_chap_user_name")
	if !ok {
		return fmt.Errorf("name is required")
	}
	

	// REST API Check
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
	connection, ok := d.GetOk("iscsi_connection")
	if ok {
		iscsiConn := connection.(*schema.Set).List()
		// Get All Storage Ports
		storagePorts, err := reconObj.GetStoragePorts()
		if err != nil {
			log.WriteDebug("TFError| error in GetStoragePorts, err: %v", err)
			return err
		}

		for _, conn := range iscsiConn {
			v := conn.(map[string]interface{})
			// Check IQN value
			iqnName := v["iscsi_initiator"].(string)
			if (iqnName != "") && (!utils.IsIqn(iqnName)) {
				log.WriteDebug("TFDebug | iqnName: %s", iqnName)
				return fmt.Errorf("iscsi_initiator %s is invalid", iqnName)
			}

			/// TODO - FIX ME - Array is not working for Plan
			portNames := v["port_names"].([]interface{})
			for _, value := range portNames {
				switch typedValue := value.(type) {
				case string:
					{
						// TODO-FIXME - Code execution is not coming inside
						// Check if Port Name Exist
						portFound := false
						if storagePorts.Data != nil {
							for _, port := range storagePorts.Data {
								if port.Nickname == typedValue {
									portFound = true
									break
								}
							}
							// If Input Port is invalid
							if !portFound {
								return fmt.Errorf("port name %s is invalid", typedValue)
							}
						}
					} // Case End
				} // Switch End
			}
		}
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

	d.SetId("")
	log.WriteInfo("compute node resource deleted successfully")
	return nil
}

func resourceVssbStorageComputeNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncHComputeNodeOperation.Lock()
	defer syncHComputeNodeOperation.Unlock()

	log.WriteInfo("starting compute node creation")
	computeNode, err := impl.CreateVssbComputeNode(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cpn := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)
	log.WriteDebug("cpn: %+v\n", *cpn)
	cpnList := []map[string]interface{}{
		*cpn,
	}
	if err := d.Set("compute_nodes", cpnList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("name", computeNode.Node.Nickname)
	d.Set("os_type", computeNode.Node.OsType)
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
	computeNode, err := impl.UpdateVssbComputeNode(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cpn := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)
	log.WriteDebug("cpn: %+v\n", *cpn)
	cpnList := []map[string]interface{}{
		*cpn,
	}
	if err := d.Set("compute_nodes", cpnList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("name", computeNode.Node.Nickname)
	d.Set("os_type", computeNode.Node.OsType)
	d.SetId(computeNode.Node.ID)
	log.WriteInfo("compute node updated successfully")
	return nil
}
