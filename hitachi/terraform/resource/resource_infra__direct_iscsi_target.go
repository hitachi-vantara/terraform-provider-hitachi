package terraform

import (
	"context"
	"fmt"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	common "terraform-provider-hitachi/hitachi/terraform/common"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"

	//resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncInfraIscsiTargetOperation = &sync.Mutex{}

func ResourceInfraIscsiTarget() *schema.Resource {
	return &schema.Resource{
		Description:   `:meta:subcategory:VSP Storage Host Group:The following request creates a host group for the port. The host mode and the host mode option can also be specified at the same time when the host group is created.`,
		CreateContext: resourceInfraIscsiTargetCreate,

		ReadContext:   resourceInfraIscsiTargetRead,
		UpdateContext: resourceInfraIscsiTargetUpdate,
		DeleteContext: resourceInfraIscsiTargetDelete,
		Schema:        schemaimpl.ResourceInfraIscsiTargetSchema,
		//CustomizeDiff: resourceMyResourceCustomDiffInfraIscsiTarget,
	}
}

func resourceInfraIscsiTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	syncInfraIscsiTargetOperation.Lock() //??
	defer syncInfraIscsiTargetOperation.Unlock()

	log.WriteInfo("starting Infra Iscsi Target create")
	storage_id, _, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)
	if storage_id != nil {

		response, err := impl.CreateInfraIscsiTarget(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		list := []map[string]interface{}{}
		for _, item := range *response {
			eachItem := impl.ConvertInfraIscsiTargetToSchema(&item)
			log.WriteDebug("it: %+v\n", *eachItem)
			list = append(list, *eachItem)
		}

		if err := d.Set("iscsitarget", list); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			/*
				d.Set("hostgroup_name", element.HostGroupName)
				d.Set("hostgroup_number", element.HostGroupId)
				d.Set("port", element.Port)
			*/
			break
		}
		log.WriteInfo("Iscsi Target created successfully")
	} else {
		serial := d.Get("serial").(int)

		iscsiTarget, err := impl.CreateIscsiTarget(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		it := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
		log.WriteDebug("iscsi target: %+v\n", *it)
		itList := []map[string]interface{}{
			*it,
		}
		if err := d.Set("iscsitarget", itList); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.Set("iscsi_target_number", iscsiTarget.IscsiTargetNumber)
		createID := fmt.Sprintf("%s%d", iscsiTarget.PortID, iscsiTarget.IscsiTargetNumber)
		d.SetId(createID)
		log.WriteInfo("iscsi target created successfully")
	}

	return nil
}

func resourceInfraIscsiTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	storage_id, _, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {
		return datasourceimpl.DataSourceInfraIscsiTargetRead(ctx, d, m)
	} else {
		return datasourceimpl.DataSourceStorageIscsiTargetRead(ctx, d, m)
	}

}

func resourceInfraIscsiTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting Infra Iscsi Target update")
	storage_id, _, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {

		response, err := impl.UpdateInfraIscsiTarget(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		list := []map[string]interface{}{}
		for _, item := range *response {
			eachItem := impl.ConvertInfraIscsiTargetToSchema(&item)
			log.WriteDebug("it: %+v\n", *eachItem)
			list = append(list, *eachItem)
		}

		if err := d.Set("iscsitarget", list); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		for _, item := range *response {
			element := &item
			d.SetId(element.ResourceId)
			/*
				d.Set("hostgroup_name", element.HostGroupName)
				d.Set("hostgroup_number", element.HostGroupId)
				d.Set("port", element.Port)
			*/
			break
		}

		log.WriteInfo("Infra Iscsi Target updated successfully")

	} else {
		serial := d.Get("serial").(int)

		iscsiTarget, err := impl.UpdateIscsiTarget(d)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		hg := impl.ConvertIscsiTargetToSchema(iscsiTarget, serial)
		log.WriteDebug("iscsiTarget: %+v\n", *hg)
		hgList := []map[string]interface{}{
			*hg,
		}
		if err := d.Set("iscsitarget", hgList); err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}

		d.Set("iscsi_target_number", iscsiTarget.IscsiTargetNumber)
		updatedID := fmt.Sprintf("%s%d", iscsiTarget.PortID, iscsiTarget.IscsiTargetNumber)
		d.SetId(updatedID)
		log.WriteInfo("iscsi target updated successfully")
	}

	return nil
}

func resourceInfraIscsiTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo("starting Infra Iscsi Target delete")
	storage_id, _, _, _ := common.GetValidateStorageIDFromSerialResource(d, m)

	if storage_id != nil {

		err := impl.DeleteIscsiTarget(d)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		log.WriteInfo("Infra Iscsi Target deleted successfully")
		return nil

	} else {

		err := impl.DeleteIscsiTarget(d)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		log.WriteInfo("iscsi target deleted successfully")
		return nil
	}

}

/*
func resourceMyResourceCustomDiffInfraHostGroup(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetInfraSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
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
*/
