package terraform

import (
	"context"

	// "fmt"

	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageChapUser() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage ISCSI Targets:Using the specified port and iSCSI target, the following request gets the CHAP user information that is specified for the iSCSI target.",
		ReadContext: DataSourceStorageChapUserRead,
		Schema:      schemaimpl.DataIscsiChapUserSchema,
	}
}

func DataSourceStorageChapUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Resource Data = %v", d)

	serial := d.Get("serial").(int)

	chapUser, err := impl.GetIscsiTargetChapUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	cu := impl.ConvertIscsiTargetChapUserToSchema(chapUser, serial)
	cuList := []map[string]interface{}{
		*cu,
	}

	if err := d.Set("chap_user", cuList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(chapUser.ChapUserID)

	log.WriteInfo("Chap User read successfully")

	return nil
}
