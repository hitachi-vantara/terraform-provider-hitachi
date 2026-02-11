package terraform

import (
	"context"
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVspPavLdev() *schema.Resource {
	return &schema.Resource{
		Description:   "VSP PAV LDEV: The following request assigns/unassigns alias LDEVs for a base LDEV.",
		CreateContext: resourceVspPavLdevCreate,
		ReadContext:   resourceVspPavLdevRead,
		UpdateContext: resourceVspPavLdevUpdate,
		DeleteContext: resourceVspPavLdevDelete,
		Schema:        schemaimpl.ResourceVspPavLdevSchema,
	}
}

func resourceVspPavLdevCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if err := impl.AssignPavAlias(d); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	serial := d.Get("serial").(int)
	baseLdevID := d.Get("base_ldev_id").(int)
	d.SetId(fmt.Sprintf("%d:%d", serial, baseLdevID))

	return resourceVspPavLdevRead(ctx, d, m)
}

func resourceVspPavLdevRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	list, err := impl.GetPavAliases(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if list == nil {
		d.SetId("")
		return nil
	}

	remaining := make(map[int]struct{})
	for _, v := range d.Get("alias_ldev_ids").([]interface{}) {
		remaining[v.(int)] = struct{}{}
	}
	for _, it := range *list {
		delete(remaining, it.LdevID)
	}
	if len(remaining) != 0 {
		d.SetId("")
		return nil
	}

	items := impl.PavAliasItemsFromList(list)

	if err := d.Set("pav_aliases", items); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVspPavLdevDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if err := impl.UnassignPavAlias(d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceVspPavLdevUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	return diag.Errorf("update is not supported for hitachi_vsp_pav_ldev (PAV alias assignment). Recreate the resource to change 'alias_ldev_ids' or 'base_ldev_id'.")
}
