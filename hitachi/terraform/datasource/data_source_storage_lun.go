package terraform

import (
	"context"
	"fmt"
	"strconv"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageLun() *schema.Resource {
	return &schema.Resource{
		Description:   `VSP Storage Volume: It returns the Lun information such as capacity, ports, paritygroup, pool etc.`,
		ReadContext:   DataSourceStorageLunRead,
		Schema:        schemaimpl.DataLunSchema,
	}
}

func DataSourceStorageLunRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	err := datasourceStorageLunValidation(d)
	if err != nil {
		return diag.FromErr(err)
	}

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

func datasourceStorageLunValidation(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// -----------------------------------------------------
	// Validate LDEV ID / HEX
	// -----------------------------------------------------

	// Retrieve both values once
	ldevIDValue, ldevIDProvided := d.GetOk("ldev_id")
	ldevHexValue, ldevHexProvided := d.GetOk("ldev_id_hex")

	log.WriteDebug("ldevIDProvided: %v, ldevHexProvided: %v", ldevIDProvided, ldevHexProvided)	
	log.WriteDebug("ldevIDValue: %v, ldevHexValue: %v", ldevIDValue, ldevHexValue)
	// Only one allowed
	if ldevIDProvided && ldevHexProvided {
		return fmt.Errorf("only one of ldev_id or ldev_id_hex may be specified, not both")
	}
	if !ldevIDProvided && !ldevHexProvided {
		return fmt.Errorf("one of ldev_id or ldev_id_hex is required")
	}

	var newLdevID *int

	// --- Handle ldev_id ---
	if ldevIDProvided {
		id, ok := ldevIDValue.(int)
		if !ok {
			return fmt.Errorf("invalid type for ldev_id")
		}
		newLdevID = &id
	}

	// --- Handle ldev_id_hex ---
	if newLdevID == nil && ldevHexProvided {
		hexStr, ok := ldevHexValue.(string)
		if !ok {
			return fmt.Errorf("invalid type for ldev_id_hex")
		}

		parsed, err := utils.HexStringToInt(hexStr)
		if err != nil {
			return fmt.Errorf("failed to parse ldev_id_hex: %v", err)
		}
		newLdevID = &parsed
	}

	return nil
}
