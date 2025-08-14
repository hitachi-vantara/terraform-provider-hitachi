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
		Description:   `VSP Storage Volume: The following request creates a volume by using the specified parity groups or pools. Specify a parity group or pool id for creating a basic volume.`,
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

func validateSerialDiff() schema.CustomizeDiffFunc {
	return customdiff.ValidateChange("serial", func(ctx context.Context, old, new, meta any) error {
		if new.(int) != old.(int) {
			return fmt.Errorf("serial should not change: old value: %d", old.(int))
		}
		return nil
	})
}

func validateLdevDiff() schema.CustomizeDiffFunc {
	return customdiff.ValidateChange("ldev_id", func(ctx context.Context, old, new, meta any) error {
		if new.(int) != old.(int) {
			return fmt.Errorf("ldev_id should not change: old value: %d", old.(int))
		}
		return nil
	})
}

func validatePoolDiff() schema.CustomizeDiffFunc {
	return customdiff.ValidateChange("pool_id", func(ctx context.Context, old, new, meta any) error {
		if new.(int) != old.(int) {
			return fmt.Errorf("pool_id should not change: old value: %d", old.(int))
		}
		return nil
	})
}

func validateParitygroupDiff() schema.CustomizeDiffFunc {
	return customdiff.ValidateChange("paritygroup_id", func(ctx context.Context, old, new, meta any) error {
		if new.(string) != old.(string) {
			return fmt.Errorf("paritygroup_id should not change: old value: %s", old.(string))
		}
		return nil
	})
}
