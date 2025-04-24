package terraform

import (

	// "time"

	"context"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/jinzhu/copier"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageIscsiTargets() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage ISCSI Targets:The following request gets information about iSCSI targets of the ports.",
		ReadContext: DataSourceStorageIscsiTargetsRead,
		Schema:      schemaimpl.DataIscsiTargetsSchema,
	}
}

func DataSourceStorageIscsiTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	var iscsiTargets terraformmodel.IscsiTargets
	ports := d.Get("port_ids").([]interface{})
	if len(ports) > 0 {
		iscsiTargetsSource, err := impl.GetIscsiTargetsByPortIds(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = copier.Copy(&iscsiTargets, iscsiTargetsSource)

		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return diag.FromErr(err)
		}
	} else {
		iscsiTargetsSource, err := impl.GetAllIscsiTargets(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = copier.Copy(&iscsiTargets, iscsiTargetsSource)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return diag.FromErr(err)
		}

		if err := d.Set("port_ids", []string{}); err != nil {
			return diag.FromErr(err)
		}

	}
	itList := []map[string]interface{}{}
	for _, it := range iscsiTargets.IscsiTargets {
		eachIt := impl.ConvertSimpleIscsiTargetToSchema(&it, serial)
		log.WriteDebug("it: %+v\n", *eachIt)
		itList = append(itList, *eachIt)
	}

	if err := d.Set("total_iscsi_target_count", len(itList)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("iscsitargets", itList); err != nil {
		return diag.FromErr(err)
	}
	log.WriteInfo(iscsiTargets)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	// d.SetId(hostgroup.PortID + strconv.Itoa(hostgroup.HostGroupNumber))
	log.WriteInfo("all iscsi target read successfully")

	return nil
}
