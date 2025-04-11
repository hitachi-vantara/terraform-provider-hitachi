package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbVolumes() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Volume:Obtains a list of volumes information.",
		ReadContext: dataSourceVssbVolumesRead,
		Schema:      schemaimpl.DataVolumeSchema,
	}
}

func dataSourceVssbVolumesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	volumes, err := impl.GetVssbVolumes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	volList := []map[string]interface{}{}
	for _, vol := range *volumes {
		eachVol := impl.ConvertVssbVolumesToSchema(&vol)
		log.WriteDebug("vol: %+v\n", *eachVol)
		volList = append(volList, *eachVol)
	}

	if err := d.Set("volumes", volList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.WriteInfo("all vssb volume read successfully")

	return nil

}

func DataSourceVssbVolumeNodes() *schema.Resource {
	return &schema.Resource{
		Description: "VOS Block Storage Volume:Obtains a list of volume information.",
		ReadContext: DataSourceVssbVolumeNodesRead,
		Schema:      schemaimpl.VolumeNodeSchema,
	}
}

func DataSourceVssbVolumeNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	volumeNode, err := impl.GetVssbVolumeNode(d)
	if err != nil {
		return diag.FromErr(err)
	}
	volume := impl.ConvertVssbVolumesToSchema(volumeNode)
	log.WriteDebug("vol: %+v\n", *volume)
	volumeList := []map[string]interface{}{
		*volume,
	}
	if err := d.Set("volume", volumeList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volumeNode.ID)
	log.WriteInfo("all vssb volume read successfully")

	return nil

}
