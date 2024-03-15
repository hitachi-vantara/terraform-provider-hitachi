package terraform

import (
	"context"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//terraform "terraform-provider-hitachi/hitachi/terraform/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func DataSourceInfraHostGroups() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Host Group:The following request gets information about host groups of the ports.",
		ReadContext: dataSourceInfraHostGroupsRead,
		Schema:      schemaimpl.DataInfraHostGroupsSchema,
	}
}

func dataSourceInfraHostGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	storage_id, _, _, err := common.GetValidateStorageIDFromSerialResource(d, m)

	if err != nil {
		log.WriteDebug("Error in get storage ID %s", err)
		return diag.FromErr(err)
	}

	if storage_id != nil {
		response, mtResponse, err := impl.GetInfraHostGroupsByPortIds(d)
		if err != nil {
			return diag.FromErr(err)
		}
		list := []map[string]interface{}{}
		if mtResponse == nil {
			for _, item := range *response {
				eachItem := impl.ConvertInfraHostGroupToSchema(&item)
				log.WriteDebug("it: %+v\n", *eachItem)
				list = append(list, *eachItem)
			}
			if err := d.Set("hostgroups", list); err != nil {
				return diag.FromErr(err)
			}
		} else {
			for _, item := range *mtResponse {
				eachItem := impl.ConvertInfraMTHostGroupToSchema(&item)
				log.WriteDebug("it: %+v\n", *eachItem)
				list = append(list, *eachItem)
			}
			if err := d.Set("partner_hostgroups", list); err != nil {
				return diag.FromErr(err)
			}
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		d.Set("total_hostgroup_count", len(list))
	} else {
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

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	}

	log.WriteInfo("host groups read successfully")
	return nil

}
