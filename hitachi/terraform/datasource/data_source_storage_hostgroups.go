package terraform

import (
	"context"
	"time"

	// "fmt"
	"strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func DataSourceStorageHostGroups() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Host Group: The following request gets information about host groups of the ports.",
		ReadContext: DataSourceStorageHostGroupsRead,
		Schema:      schemaimpl.DataHostGroupsSchema,
	}
}

func DataSourceStorageHostGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	var hostgroups terraformmodel.HostGroups
	ports := d.Get("port_ids").([]interface{})
	if len(ports) > 0 {
		hostgroupsSource, err := impl.GetHostGroupsByPortIds(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = copier.Copy(&hostgroups, hostgroupsSource)

		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return diag.FromErr(err)
		}
	} else {
		hostgroupsSource, err := impl.GetAllHostGroups(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = copier.Copy(&hostgroups, hostgroupsSource)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return diag.FromErr(err)
		}

	}

	hgList := []map[string]interface{}{}
	for _, hg := range hostgroups.HostGroups {
		eachHg := impl.ConvertSimpleHostGroupToSchema(&hg, serial)
		log.WriteDebug("hg: %+v\n", *eachHg)
		hgList = append(hgList, *eachHg)
	}

	if err := d.Set("hostgroups", hgList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("total_hostgroup_count", len(hgList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	// d.SetId(hostgroup.PortID + strconv.Itoa(hostgroup.HostGroupNumber))
	log.WriteInfo("all hg read successfully")

	return nil
}
