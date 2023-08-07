package terraform

import (
	"context"
	"strconv"
	"time"

	// "fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbDashboard() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSS Block Dashboard:Obtains the information about Dashboard Information.",
		ReadContext: DataSourceVssbDashboardRead,
		Schema:      schemaimpl.DataVssbDashboardSchema,
	}
}

func DataSourceVssbDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	dashboard, err := impl.GetDashboardInfo(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.WriteDebug("Dashboard: %+v\n", dashboard)

	db := impl.ConvertVssbDashboardToSchema(dashboard)
	itList := []map[string]interface{}{
		*db,
	}

	if err := d.Set("dashboard_info", itList); err != nil {
		log.WriteDebug("err: %v\n", err)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.WriteInfo("Dashboard read successfully")

	return nil
}
