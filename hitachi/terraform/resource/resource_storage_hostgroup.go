package terraform

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terrcommon "terraform-provider-hitachi/hitachi/terraform/common"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncHostGroupOperation = &sync.Mutex{}

func ResourceStorageHostGroup() *schema.Resource {
	return &schema.Resource{
		Description:   `VSP Storage Host Group: The following request creates a host group for the port. The host mode and the host mode option can also be specified at the same time when the host group is created.`,
		CreateContext: resourceStorageHostGroupCreate,
		ReadContext:   resourceStorageHostGroupRead,
		UpdateContext: resourceStorageHostGroupCreate,
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
	createID := fmt.Sprintf("%s,%d,%s", hostGroup.PortID, hostGroup.HostGroupNumber, hostGroup.HostGroupName)
	d.SetId(createID)
	log.WriteInfo("hg created successfully")

	return nil
}

func resourceStorageHostGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	hostgroup, err := impl.GetHostGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	hg := impl.ConvertHostGroupToSchema(hostgroup, serial)
	log.WriteDebug("hg: %+v\n", *hg)

	hgList := []map[string]interface{}{
		*hg,
	}

	if err := d.Set("hostgroup", hgList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s,%d,%s", hostgroup.PortID, hostgroup.HostGroupNumber, hostgroup.HostGroupName))
	log.WriteInfo("hg read successfully")

	return nil
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
	updatedID := fmt.Sprintf("%s,%d,%s", hostGroup.PortID, hostGroup.HostGroupNumber, hostGroup.HostGroupName)
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

	if d.Id() == "" {
		// create
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

		// Validate allowed name format
		pattern := `^[A-Za-z0-9.@_:-]{1,64}$`
		reg := regexp.MustCompile(pattern)

		hgName, ok := d.GetOk("hostgroup_name")
		// Check if the value matches the pattern
		if ok {
			if !reg.MatchString(hgName.(string)) {
				return fmt.Errorf("hostgroup_name must be 1–64 chars, alphanumeric or . @ _ : -, cannot start with '-'")
			}
			if strings.HasPrefix(hgName.(string), "-") {
				return fmt.Errorf("hostgroup_name cannot start with a hyphen (-)")
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
				return fmt.Errorf("%v", err.Error())
			}
		}
	} else {
		// update
		storedPortID, storedHgNum, storedHgName, err := terrcommon.ParseHostGroupFromID(d.Id())
		if err != nil {
			return err
		}

		// Read new values from config
		newPortID := d.Get("port_id").(string)

		//
		// PORT ID IMMUTABLE
		//
		if newPortID != storedPortID {
			return fmt.Errorf(
				"existing hostgroup does not match port_id provided (expected %s, got %s)",
				storedPortID, newPortID,
			)
		}

		//
		// HOSTGROUP NUMBER IMMUTABLE (if user provides it)
		//
		if v, ok := d.GetOkExists("hostgroup_number"); ok {
			newHgNum := v.(int)
			if newHgNum != 0 && newHgNum != storedHgNum {
				return fmt.Errorf(
					"existing hostgroup does not match hostgroup_number provided (expected %d, got %d)",
					storedHgNum, newHgNum,
				)
			}
		}

		//
		// HOSTGROUP NAME IMMUTABLE
		//
		if v, ok := d.GetOkExists("hostgroup_name"); ok {
			newHgName := v.(string)

			// Now check immutability
			if newHgName != storedHgName {
				return fmt.Errorf(
					"existing hostgroup does not match hostgroup_name provided (expected %s, got %s)",
					storedHgName, newHgName,
				)
			}
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
		hmode := impl.HostModeUserToRestConversion[strings.ToLower(userhmode)]
		if hmode == "" {
			err := fmt.Errorf("invalid hostmode specified %v", userhmode)
			return err

		}
	}

	// hgModeOpt, ok := d.GetOk("host_mode_options")

	//TODO:   hgModeOpt validation

	// Fix console output
	d.SetNewComputed("hostgroup")

	return validateLun(ctx, d, m)

}

func validateLun(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	lunSet, ok := d.GetOk("lun")
	if !ok {
		return nil
	}

	for _, item := range lunSet.(*schema.Set).List() {
		m := item.(map[string]interface{})

		// Retrieve actual values
		idVal := m["ldev_id"].(int)      // default = -1
		hexVal := m["ldev_id_hex"].(string) // default = ""

		hasID := idVal != -1   // -1 = not provided
		hasHex := hexVal != "" // empty string = not provided

		// Both provided → invalid
		if hasID && hasHex {
			return fmt.Errorf("only one of ldev_id or ldev_id_hex may be specified in each lun block")
		}

		// Neither provided → invalid
		if !hasID && !hasHex {
			return fmt.Errorf("one of ldev_id or ldev_id_hex must be specified in each lun block")
		}
	}

	return nil
}
