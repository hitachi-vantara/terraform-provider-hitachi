package terraform

import (
	"context"
	"errors"

	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraform "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraIscsiTarget() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage iSCSI Target:The following request obtains information about iSCSI Target.",
		ReadContext: DataSourceInfraIscsiTargetRead,
		Schema:      schemaimpl.DataInfraIscsiTargetSchema,
	}
}

func DataSourceInfraIscsiTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	iscsi_id := -1

	iscsi_name := d.Get("iscsi_name").(string)
	iid, okId := d.GetOk("iscsi_target_number")
	if okId {
		iscsi_id = iid.(int)
	}

	if iscsi_name != "" && iscsi_id != -1 {
		err := errors.New("both iscsi_name  and iscsi_target_number are not allowed. Either iscsi_name or iscsi_target_number or none of them can be specified")
		return diag.FromErr(err)
	}

	var response *[]terraform.InfraIscsiTargetInfo
	var err error
	list := []map[string]interface{}{}

	response, err = impl.GetInfraGwIscsiTargets(d)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, item := range *response {
		eachItem := impl.ConvertInfraIscsiTargetToSchema(&item)
		log.WriteDebug("it: %+v\n", *eachItem)
		list = append(list, *eachItem)
	}

	if err := d.Set("iscsi_target", list); err != nil {
		return diag.FromErr(err)
	}

	if iscsi_name == "" && iscsi_id == -1 {
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	} else {
		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			d.Set("iscsi_target_id", element.ResourceId)
			d.Set("iscsi_name", element.ISCSIName)
			d.Set("iscsi_target_number", element.ISCSIId)
			break
		}
	}
	log.WriteInfo("iscsiTargets read successfully")
	return nil

}
