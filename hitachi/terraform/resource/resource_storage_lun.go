package terraform

import (
	"context"
	"fmt"
	"strconv"

	// "time"
	// "errors"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

var syncLunOperation = &sync.Mutex{}

func ResourceStorageLun() *schema.Resource {
	return &schema.Resource{
		Description:   ":meta:subcategory:VSP Storage Volume:The following request creates a volume by using the specified parity groups or pools. Specify a parity group or pool id for creating a basic volume.",
		CreateContext: resourceStorageLunCreate,
		ReadContext:   resourceStorageLunRead,
		UpdateContext: resourceStorageLunUpdate,
		DeleteContext: resourceStorageLunDelete,
		Schema:        schemaimpl.ResourceLunSchema,
		CustomizeDiff: customDiffFunc(),
	}
}

func resourceStorageLunCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncLunOperation.Lock() //??
	defer syncLunOperation.Unlock()

	log.WriteInfo("starting lun create")

	serial := d.Get("serial").(int)

	logicalUnit, err := impl.CreateLun(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	lun := impl.ConvertLunToSchema(logicalUnit, serial)
	log.WriteDebug("lun: %+v\n", *lun)
	lunList := []map[string]interface{}{
		*lun,
	}
	if err := d.Set("volume", lunList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	//d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
	d.SetId(strconv.Itoa(logicalUnit.LdevID))
	log.WriteInfo("lun created successfully")

	return nil
}

func resourceStorageLunRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceStorageLunRead(ctx, d, m)
}

func resourceStorageLunUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting lun update")

	serial := d.Get("serial").(int)

	logicalUnit, err := impl.UpdateLun(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	terraformModelLun := terraformmodel.LogicalUnit{}
	err = copier.Copy(&terraformModelLun, logicalUnit)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil
	}

	lun := impl.ConvertLunToSchema(&terraformModelLun, serial)
	log.WriteDebug("lun: %+v\n", *lun)
	lunList := []map[string]interface{}{
		*lun,
	}
	if err := d.Set("volume", lunList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("ldev_id", logicalUnit.LdevID) // input may have an empty ldev_id
	d.SetId(strconv.Itoa(logicalUnit.LdevID))
	log.WriteInfo("lun updated successfully")

	return nil
}

func resourceStorageLunDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting lun delete")

	err := impl.DeleteLun(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("lun deleted successfully")
	return nil
}

func customDiffFunc() schema.CustomizeDiffFunc {
	return customdiff.All(
		validateSizeDiff(),
		// validateSerialDiff(),
		// validateLdevDiff(),
		// validatePoolDiff(),
		// validateParitygroupDiff(),
	)
}

func validateSizeDiff() schema.CustomizeDiffFunc {
	return customdiff.ValidateChange("size_gb", func(ctx context.Context, old, new, meta any) error {
		// size must only increase not decrease
		if new.(int) < old.(int) {
			return fmt.Errorf("new size_gb value must be greater than old value: %d", old.(int))
		}
		return nil
	})
}

func ValidateSanStorageVolumeDIff(d *schema.ResourceDiff) error {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	_, ok := d.GetOk("size_gb")
	if !ok {
		return fmt.Errorf("size_gb must be greater than 0 for create")
	}

	var pool_name = ""
	var paritygroup_id = ""
	pool_id, exists := d.GetOk("pool_id")
	okPO := exists || (pool_id.(int) == 0)

	pool_name = d.Get("pool_name").(string)
	paritygroup_id = d.Get("paritygroup_id").(string)
	log.WriteDebug("Pool ID=%v Pool Name=%v PG=%v\n", pool_id, pool_name, paritygroup_id)

	log.WriteDebug("ok=%v \n", ok)

	count := 0
	if okPO && pool_id != -1 {
		count++
	}
	if pool_name != "" {
		count++
	}
	if paritygroup_id != "" {
		count++
	}
	log.WriteDebug("count=%v\n", count)
	if count != 1 {
		return fmt.Errorf("either pool_id or pool_name or paritygroup_id is required to create volume")
	}

	// if pool_name != "" {
	// 	ppid, err := GetPoolIdFromPoolName(d, pool_name)
	// 	if err != nil {
	// 		return fmt.Errorf("could not find a pool with name %v", pool_name)
	// 	}
	// }
	return nil
}
