package terraform

import (
	"context"
	// "fmt"
	"strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageLun() *schema.Resource {
	return &schema.Resource{
		Description: `VSP Storage Volume:It returns the Lun information such as capacity, ports, paritygroup, pool etc.`,
		ReadContext: DataSourceStorageLunRead,
		Schema:      schemaimpl.DataLunSchema,
	}
}

func DataSourceStorageLunRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	logicalUnit, err := impl.GetLun(d)
	if err != nil {
		return diag.FromErr(err)
	}

	lun := impl.ConvertLunToSchema(logicalUnit, serial)
	log.WriteDebug("lun: %+v\n", *lun)

	lunList := []map[string]interface{}{
		*lun,
	}
	if err := d.Set("volume", lunList); err != nil {
		return diag.FromErr(err)
	}

	// always run
	// d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.SetId(strconv.Itoa(logicalUnit.LdevID))
	log.WriteInfo("lun read successfully")

	return nil
}
