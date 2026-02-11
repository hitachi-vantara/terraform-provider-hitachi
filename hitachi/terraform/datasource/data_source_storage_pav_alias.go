package terraform

import (
	"context"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStoragePavAlias() *schema.Resource {
	return &schema.Resource{
		Description: "VSP PAV Alias: The following request obtains information about PAV aliases.",
		ReadContext: DataSourceStoragePavAliasRead,
		Schema:      schemaimpl.DataPavAliasSchema,
	}
}

func DataSourceStoragePavAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	list, err := impl.GetPavAliases(d)
	if err != nil {
		return diag.FromErr(err)
	}

	items := impl.PavAliasItemsFromList(list)

	if err := d.Set("pav_aliases", items); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}
