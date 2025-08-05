package terraform

import (
	"context"
	"fmt"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbComputeNodes() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Compute Node:Obtains a list of compute node information.",
		ReadContext: DataSourceVssbComputeNodesRead,
		Schema:      schemaimpl.DataComputeNodeSchema,
	}
}

func DataSourceVssbComputeNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	computeNodes, err := impl.GetVssbComputeNodes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	compNodeName, ok := d.Get("compute_node_name").(string)

	log.WriteDebug("compute Node name : %v, ok : %v", compNodeName, ok)

	if compNodeName == "" {
		// if compute_node_name is not present in the datasource we display all compute nodes

		compNodeList := []map[string]interface{}{}
		for _, comNode := range *computeNodes {

			computeNode, err := impl.GetVssbComputeNode(d, comNode.ID)
			if err != nil {
				return diag.FromErr(err)
			}

			eachNode := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)
			log.WriteDebug("compute node: %+v\n", *eachNode)
			compNodeList = append(compNodeList, *eachNode)
		}

		if err := d.Set("compute_nodes", compNodeList); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("compute_node_name", ""); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("all vssb compute node read successfully")

	} else {
		// if compute_node_name is present find the id of the node

		var comNodeId string
		for _, comNode := range *computeNodes {
			if comNode.Nickname == compNodeName {
				comNodeId = comNode.ID
				break
			}
		}
		log.WriteDebug("compute Node ID  : %v", comNodeId)

		if comNodeId != "" {
			computeNode, err := impl.GetVssbComputeNode(d, comNodeId)
			if err != nil {
				return diag.FromErr(err)
			}

			terraNode := impl.ConvertVssbComputeNodeWithPathDetailsToSchema(computeNode)

			cnList := []map[string]interface{}{
				*terraNode,
			}

			if err := d.Set("compute_nodes", cnList); err != nil {
				return diag.FromErr(err)
			}

			d.SetId(comNodeId)

			log.WriteInfo("Compute Node read successfully")

		} else {
			err := fmt.Errorf("compute node %v does not exit", compNodeName)
			return diag.FromErr(err)
		}
	}
	return nil

}
