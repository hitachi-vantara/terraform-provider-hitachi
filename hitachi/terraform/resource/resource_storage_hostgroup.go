package terraform

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	// "fmt"

	// "time"
	// "errors"
	"sync"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	terraform "terraform-provider-hitachi/hitachi/terraform/impl"

	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncHostGroupOperation = &sync.Mutex{}

func ResourceStorageHostGroup() *schema.Resource {
	return &schema.Resource{
		Description:   `VSP Storage Host Group: The following request creates a host group for the port. The host mode and the host mode option can also be specified at the same time when the host group is created.`,
		CreateContext: resourceStorageHostGroupCreate,

		ReadContext:   resourceStorageHostGroupRead,
		UpdateContext: resourceStorageHostGroupUpdate,
		DeleteContext: resourceStorageHostGroupDelete,
		Schema:        schemaimpl.ResourceHostGroupSchema,
		CustomizeDiff: resourceMyResourceCustomDiffHostGroup,
	}
}

func resourceStorageHostGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncHostGroupOperation.Lock() //??
	defer syncHostGroupOperation.Unlock()

	log.WriteInfo("starting hg create")

	serial := d.Get("serial").(int)

	hostGroup, err := impl.CreateHostGroup(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hg := impl.ConvertHostGroupToSchema(hostGroup, serial)
	log.WriteDebug("hg: %+v\n", *hg)
	hgList := []map[string]interface{}{
		*hg,
	}
	if err := d.Set("hostgroup", hgList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("hostgroup_number", hostGroup.HostGroupNumber)
	createID := fmt.Sprintf("%s%d", hostGroup.PortID, hostGroup.HostGroupNumber)
	d.SetId(createID)
	log.WriteInfo("hg created successfully")

	return nil
}

func resourceStorageHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return datasourceimpl.DataSourceStorageHostGroupRead(ctx, d, m)
}

func resourceStorageHostGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting hg update")

	serial := d.Get("serial").(int)

	hostGroup, err := impl.UpdateHostGroup(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	hg := impl.ConvertHostGroupToSchema(hostGroup, serial)
	log.WriteDebug("hg: %+v\n", *hg)
	hgList := []map[string]interface{}{
		*hg,
	}
	if err := d.Set("hostgroup", hgList); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("hostgroup_number", hostGroup.HostGroupNumber)
	updatedID := fmt.Sprintf("%s%d", hostGroup.PortID, hostGroup.HostGroupNumber)
	d.SetId(updatedID)
	log.WriteInfo("hg updated successfully")

	return nil
}

func resourceStorageHostGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting hg delete")

	err := impl.DeleteHostGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.WriteInfo("hg deleted successfully")
	return nil
}

func resourceMyResourceCustomDiffHostGroup(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in Reconciler NewEx, err: %v", err)
		return err
	}
	// Define the regular expression pattern
	pattern := "^[a-zA-Z0-9_-]{1,64}$"
	reg := regexp.MustCompile(pattern)

	hgName, ok := d.GetOk("hostgroup_name")
	// Check if the value matches the pattern
	if ok {
		if !reg.MatchString(hgName.(string)) {
			return fmt.Errorf("hostgroup_name Value is alphanumeric and can only accpet '-' and '_' and within the range of 1-64 characters")
		}
	}
	// vlidate hostgroup_number ranges from 0 to 255
	hg_number, ok := d.GetOk("hostgroup_number")
	if ok {
		hgNumberInt := hg_number.(int)
		if hgNumberInt < 0 || hgNumberInt > 255 {
			return fmt.Errorf("hostgroup_number Value should be between 0 and 255")
		}
	}
	portId, ok := d.GetOk("port_id")
	if ok {
		_, err := reconObj.GetStoragePortByPortId(portId.(string))
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	/*
		//validate hostmodes from given regex
		validHostModes := []string{"HP-UX", "SOLARIS", "AIX", "WIN", "LINUX/IRIX", "TRU64", "OVMS", "NETWARE", "VMWARE", "VMWARE_EX", "WIN_EX"}
		pattern = fmt.Sprintf(`\b(?:%s)\b`, strings.Join(validHostModes, "|"))
		hgMode, ok := d.GetOk("host_mode")
		if ok {

			match, _ := regexp.MatchString(pattern, hgMode.(string))

			if !match {
				return fmt.Errorf("hostmode Value should be with in :  %s", strings.Join(validHostModes, ","))
			}
		}
	*/
	hostmode, ok := d.GetOk("host_mode")
	if ok {
		userhmode := hostmode.(string)
		hmode := terraform.HostModeUserToRestConversion[strings.ToLower(userhmode)]
		if hmode == "" {
			err := fmt.Errorf("invalid hostmode specified %v", userhmode)
			return err

		}
	}

	// hgModeOpt, ok := d.GetOk("host_mode_options")

	//TODO:   hgModeOpt validation
	return nil
}
