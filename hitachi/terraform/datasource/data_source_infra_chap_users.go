package terraform

import (
	"context"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	//terraform "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceInfraChapUsers() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage CHAP Users:The following request gets information about the chap users.",
		ReadContext: dataSourceInfraChapUsersRead,
		Schema:      schemaimpl.DataInfraIscsiChapUsersSchema,
	}
}

func dataSourceInfraChapUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	response, err := impl.GetInfraChapUsers(d)
	if err != nil {
		return diag.FromErr(err)
	}

	/*
		list := []map[string]interface{}{}
		for _, item := range *response {
			eachItem := impl.ConvertInfraHostGroupToSchema(&item)
			log.WriteDebug("it: %+v\n", *eachItem)
			list = append(list, *eachItem)
		}
	*/

	if err := d.Set("chap_users", response); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.WriteInfo("chap users read successfully")
	return nil

}
