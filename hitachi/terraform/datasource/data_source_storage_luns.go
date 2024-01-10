package terraform

import (
	"context"
	"time"

	// "fmt"
	"strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageLuns() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage Volume:It returns all luns information from given storage device.",
		ReadContext: DataSourceStorageLunsRead,
		Schema:      schemaimpl.DataLunsSchema,
	}
}

func DataSourceStorageLunsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	logicalUnits, err := impl.GetRangeOfLuns(d)
	if err != nil {
		return diag.FromErr(err)
	}

	lunList := []map[string]interface{}{}

	for _, lun := range *logicalUnits {
		eachLun := impl.ConvertLunToSchema(&lun, serial)
		log.WriteDebug("eachLun: %+v\n", &eachLun)
		lunList = append(lunList, *eachLun)
	}

	if err := d.Set("volumes", lunList); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	//d.SetId(strconv.Itoa(logicalUnits.LdevID))
	log.WriteInfo("range of luns read successfully")

	return nil
}
