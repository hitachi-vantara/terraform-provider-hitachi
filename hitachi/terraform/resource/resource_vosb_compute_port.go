package terraform

import (
	"context"
	"fmt"
	"strings"

	// "fmt"
	// "time"
	// "errors"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"

	//utils "terraform-provider-hitachi/hitachi/common/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"

	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"

	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncComputePortOperation = &sync.Mutex{}

func ResourceVssbStorageComputePort() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block iSCSI Target CHAP User:The following request sets the CHAP user.",
		CreateContext: resourceVssbComputrPortCreate,
		ReadContext:   resourceVssbComputrPortRead,
		UpdateContext: resourceVssbComputrPortUpdate,
		DeleteContext: resourceVssbComputrPortDelete,
		Schema:        schemaimpl.ResourceComputePortSchema,
		CustomizeDiff: VssbStorageComputePortCustomDiff,
	}
}

func resourceVssbComputrPortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("resource delete")

	_, err := impl.AllowChapUsersToAccessComputePort(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("chap user resource deleted successfully")
	return nil
}

func resourceVssbComputrPortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncComputePortOperation.Lock()
	defer syncComputePortOperation.Unlock()

	log.WriteInfo("starting associating chap users with port")
	portInfo, err := impl.AllowChapUsersToAccessComputePort(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	cp := impl.ConvertVssbIscsiPortAuthSettingsToSchema(portInfo)
	log.WriteDebug("Compute Port: %+v\n", *cp)
	cuList := []map[string]interface{}{
		*cp,
	}
	if err := d.Set("compute_port", cuList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(portInfo.Port.ID)
	return nil
}

func resourceVssbComputrPortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceVssbComputePortRead(ctx, d, m)
}

func resourceVssbComputrPortUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncComputePortOperation.Lock()
	defer syncComputePortOperation.Unlock()

	log.WriteInfo("starting associating chap users with port")
	portInfo, err := impl.AllowChapUsersToAccessComputePort(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//cp := impl.ConvertVssbPortDetailSettingsToSchema(portInfo)

	cp := impl.ConvertVssbIscsiPortAuthSettingsToSchema(portInfo)
	log.WriteDebug("Compute Port: %+v\n", *cp)
	cuList := []map[string]interface{}{
		*cp,
	}
	if err := d.Set("compute_port", cuList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(portInfo.Port.ID)
	log.WriteInfo("compute port updated successfully")
	return nil
}

func VssbStorageComputePortCustomDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Local Check
	name := d.Get("authentication_settings").(string)
	target_chap_users := d.Get("target_chap_users").([]interface{})
	// target_chap_users := d.Get("target_chap_users")

	// Check the length of the map.
	length := len(utils.ConvertInterfaceToSlice(target_chap_users))

	// target_chap_users = target_chap_users.([]interface{})
	if strings.ToLower(name) == "none" && length > 0 {
		return fmt.Errorf("if authentication_settings is set to 'None',target_chap_users value should be an empty list")
	}

	data := &schema.ResourceData{}
	res, err := impl.GetVssbComputePortByName(data, d.Get("vosb_address").(string), d.Get("name").(string))
	if err == nil {
		cp := impl.ConvertVssbIscsiPortAuthSettingsToSchema(res)
		log.WriteDebug("Compute Port: %+v\n", *cp)
		cuList := []map[string]interface{}{
			*cp,
		}
		if err := data.Set("compute_port", cuList); err != nil {
			data.SetId("")
			//return fmt.Errorf("could not set compute_port in the custom diff")
			return nil
		}
		/*
			if strings.ToLower(name) == "none" && length == 0 && len(res.ChapUsers.Data) > 0 {
				return fmt.Errorf("target_chap_users are still present in the VOSB storage when authentication_settings is set to 'None' and target_chap_users is to a empty list")
			}
		*/
		log.WriteInfo("compute port updated successfully")
	}

	return nil
}
