package terraform

import (

	// "fmt"

	// "time"

	"context"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageChapUsers() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage ISCSI Targets:Using the specified port and iSCSI target, the following request gets the CHAP users information that is specified for the iSCSI target.",
		ReadContext: DataSourceStorageChapUsersRead,
		Schema:      schemaimpl.DataIscsiChapUsersSchema,
	}
}

func DataSourceStorageChapUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	chapUsers, err := impl.GetIscsiTargetChapUsers(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteDebug("chapUsers: %+v\n", chapUsers)

	itList := []map[string]interface{}{}
	for _, cu := range chapUsers.IscsiTargetChapUsers {
		eachIt := impl.ConvertIscsiTargetChapUserToSchema(&cu, serial)
		log.WriteDebug("it: %+v\n", *eachIt)
		itList = append(itList, *eachIt)
	}

	if err := d.Set("chap_users", itList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.WriteInfo("all iscsi target chap users read successfully")

	return nil
}
