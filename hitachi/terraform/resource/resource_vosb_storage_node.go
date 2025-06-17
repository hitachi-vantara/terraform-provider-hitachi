package terraform

import (
	"context"
	"strconv"
	"time"

	// "fmt"

	// "fmt"

	// "time"
	// "errors"
	"sync"

	// cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	// utils "terraform-provider-hitachi/hitachi/common/utils"
	// reconimpl "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/impl"
	// reconcilermodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncHStorageNodeOperation = &sync.Mutex{}

func ResourceVssbStorageNode() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block Storage Node:Registers the information of the storage node.",
		CreateContext: resourceVssbStorageNodeCreate,
		ReadContext:   resourceVssbStorageNodeRead,
		UpdateContext: resourceVssbStorageNodeUpdate,
		DeleteContext: resourceVssbStorageNodeDelete,
		Schema:        schemaimpl.ResourceVssbStorageNodeSchema,
		// CustomizeDiff: resourceStorageNodeCustomDiff,
	}
}


func resourceVssbStorageNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting storage node resource delete")

	// err := impl.DeleteVssbStorageNodeResource(d)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId("")
	log.WriteInfo("storage node resource deleted successfully")
	return nil
}

func resourceVssbStorageNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncHStorageNodeOperation.Lock()
	defer syncHStorageNodeOperation.Unlock()

	log.WriteInfo("starting storage node creation")

	err := impl.CreateVssbStorageNode(d)
	if err != nil {
		return diag.FromErr(err)
	}

	setOutput(d)
	log.WriteInfo("storage node created successfully")
	return nil
}

func setOutput(d *schema.ResourceData) {
		storageNodes, err := impl.GetVssbStorageNodes(d)
		if err != nil {
			return 
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storageNodes {
			eachSp := impl.ConvertVssbStorageNodeToSchema(&sp)
			// log.WriteDebug("storage node: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("storage_nodes", spList); err != nil {
			return 
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))	
}


func resourceVssbStorageNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbStorageNodesRead(ctx, d, m)
}

func resourceVssbStorageNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncHStorageNodeOperation.Lock()
	defer syncHStorageNodeOperation.Unlock()

	log.WriteInfo("starting storage node update")
	// computeNode, err := impl.UpdateVssbStorageNode(d)
	// if err != nil {
	// 	d.SetId("")
	// 	return diag.FromErr(err)
	// }

	// cpn := impl.ConvertVssbStorageNodeWithPathDetailsToSchema(computeNode)
	// log.WriteDebug("cpn: %+v\n", *cpn)
	// cpnList := []map[string]interface{}{
	// 	*cpn,
	// }
	// if err := d.Set("compute_nodes", cpnList); err != nil {
	// 	d.SetId("")
	// 	return diag.FromErr(err)
	// }

	// d.Set("name", computeNode.Node.Nickname)
	// d.Set("os_type", computeNode.Node.OsType)
	// d.SetId(computeNode.Node.ID)
	// log.WriteInfo("storage node updated successfully")
	return nil
}
