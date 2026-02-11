package terraform

import (

	// "time"

	"context"
	"strconv"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/jinzhu/copier"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceStorageParityGroup() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Storage Parity Group: The following request obtains information about a parity group.",
		ReadContext: DataSourceStorageParityGroupRead,
		Schema:      schemaimpl.DataParityGroupSchema,
	}
}

func DataSourceStorageParityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	var parityGroup terraformmodel.ParityGroup

	parityGroupSource, err := impl.GetParityGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = copier.Copy(&parityGroup, parityGroupSource)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return diag.FromErr(err)
	}

	pgList := []map[string]interface{}{}
	pg := impl.ConvertParityGroupToSchema(&parityGroup, serial)
	pgList = append(pgList, *pg)
	log.WriteDebug("pgList: %+v\n", pgList)

	if err := d.Set("parity_group", pgList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.WriteInfo("parity group read successfully")

	return nil
}
