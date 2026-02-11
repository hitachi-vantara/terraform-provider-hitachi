package terraform

import (
	"context"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func DataSourceStorageSupportedHostModes() *schema.Resource {
	return &schema.Resource{
		Description: "VSP Supported Host Modes: The following request obtains information about supported host modes and host mode options.",
		ReadContext: DataSourceStorageSupportedHostModesRead,
		Schema:      schemaimpl.DataSupportedHostModesSchema,
	}
}

func DataSourceStorageSupportedHostModesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var out terraformmodel.SupportedHostModes

	source, err := impl.GetSupportedHostModes(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := copier.Copy(&out, source); err != nil {
		log.WriteDebug("TFError| error in Copy from impl to terraform structure, err: %v", err)
		return diag.FromErr(err)
	}

	hostModes := make([]map[string]interface{}, 0, len(out.HostModes))
	for _, hm := range out.HostModes {
		hostModes = append(hostModes, map[string]interface{}{
			"host_mode_id":      hm.HostModeID,
			"host_mode_name":    hm.HostModeName,
			"host_mode_display": hm.HostModeDisplay,
		})
	}

	hostModeOptions := make([]map[string]interface{}, 0, len(out.HostModeOptions))
	for _, hmo := range out.HostModeOptions {
		hostModeOptions = append(hostModeOptions, map[string]interface{}{
			"host_mode_option_id":          hmo.HostModeOptionID,
			"host_mode_option_description": hmo.HostModeOptionDescription,
			"scope":                        hmo.Scope,
			"required_host_modes":          hmo.RequiredHostModes,
		})
	}

	if err := d.Set("host_modes", hostModes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_mode_options", hostModeOptions); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}
