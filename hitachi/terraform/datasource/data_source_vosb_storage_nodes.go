package terraform

import (
	"context"
	"errors"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbStorageNodes() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Node:Obtains a list of storage nodes information.",
		ReadContext: DataSourceVssbStorageNodesRead,
		Schema:      schemaimpl.DataVssbStorageNodeSchema,
	}
}

func DataSourceVssbStorageNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	node_name := ""
	if  d.Get("node_name") != nil {
		node_name = d.Get("node_name").(string)
	}

	if node_name == "" {

		storageNodes, err := impl.GetVssbStorageNodes(d)
		if err != nil {
			return diag.FromErr(err)
		}

		spList := []map[string]interface{}{}
		for _, sp := range *storageNodes {
			eachSp := impl.ConvertVssbStorageNodeToSchema(&sp)
			log.WriteDebug("storage node: %+v\n", *eachSp)
			spList = append(spList, *eachSp)
		}

		if err := d.Set("nodes", spList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("all vssb storage nodes read successfully")

		return nil
	} else {
		if !utils.IsValidName(node_name) {
			err := errors.New("node name can not exceed 255 characters")
			return diag.FromErr(err)
		}

		node, err := impl.GetVssbNode(d)

		if err != nil {
			return diag.FromErr(err)
		}

		pList := []map[string]interface{}{
			*impl.ConvertVssbStorageNodeToSchema(node),
		}
		log.WriteDebug("pList: %+v\n", pList)
		if err := d.Set("nodes", pList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(node.ID)
		log.WriteInfo("node read successfully")
		return nil
	}

}
