package terraform

import (
	"context"
	"errors"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//terraform "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraHostGroup() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Host Group:The following request gets information about host group of the port.",
		ReadContext: DataSourceInfraHostGroupRead,
		Schema:      schemaimpl.DataInfraHostGroupSchema,
	}
}

func DataSourceInfraHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	hostgroup_name := d.Get("hostgroup_name").(string)

	hostgroup_id := -1
	hid, okId := d.GetOk("hostgroup_number")
	if okId {
		hostgroup_id = hid.(int)
	}

	if hostgroup_name != "" && hostgroup_id != -1 {
		err := errors.New("both hostgroup_name  and hostgroup_number are not allowed. Either hostgroup_name or hostgroup_number or none of them can be specified")
		return diag.FromErr(err)
	}

	storage_id, _, _, err := common.GetValidateStorageIDFromSerialResource(d, m)

	if err != nil {
		log.WriteDebug("Error in get storage ID %s", err)
		return diag.FromErr(err)
	}

	if storage_id != nil {
		response, mtResponse, err := impl.GetInfraHostGroups(d)
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

			if err := d.Set("hostgroup", list); err != nil {
				return diag.FromErr(err)
			}

			if hostgroup_name == "" && hostgroup_id == -1 {
				d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
			} else {
				for _, item := range *response {
					element := &item
					d.SetId(element.ResourceId)
					d.Set("hostgroup_name", element.HostGroupName)
					d.Set("hostgroup_number", element.HostGroupId)
					d.Set("port", element.Port)
					break
				}
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
	} else {
		serial := d.Get("serial").(int)

		hostgroup, err := impl.GetHostGroup(d)
		if err != nil {
			return diag.FromErr(err)
		}

		hg := impl.ConvertHostGroupToSchema(hostgroup, serial)
		log.WriteDebug("hg: %+v\n", *hg)

		hgList := []map[string]interface{}{
			*hg,
		}

		if err := d.Set("hostgroup", hgList); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(hostgroup.PortID + strconv.Itoa(hostgroup.HostGroupNumber))
	}
	log.WriteInfo("host groups read successfully")
	return nil

}
