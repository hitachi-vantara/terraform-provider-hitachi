package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageSystem() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSP Storage System:It returns the storage device related information.",
		ReadContext: dataSourceStorageSystemRead,
		Schema:      schemaimpl.StorageSystemSchema,
	}
}

func dataSourceStorageSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ss, err := impl.GetStorageSystem(d)
	if err != nil {
		return diag.FromErr(err)
	}

	san := map[string]interface{}{
		"storage_device_id":      ss.StorageDeviceID,
		"storage_serial_number":  ss.SerialNumber,
		"storage_device_model":   ss.Model,
		"dkc_micro_code_version": ss.MicroVersion,
		"management_ip":          ss.MgmtIP,
		"svp_ip":                 ss.SvpIP,
		"controller1_ip":         ss.ControllerIP1,
		"controller2_ip":         ss.ControllerIP2,
		"total_capacity_in_mb":   ss.TotalCapacityInMB,
		"free_capacity_in_mb":    ss.FreeCapacityInMB,
		"used_capacity_in_mb":    ss.UsedCapacityInMB,
	}

	log.WriteDebug("san: %+v\n", san)
	sanList := []map[string]interface{}{
		san,
	}

	if err := d.Set("storage_system", sanList); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
