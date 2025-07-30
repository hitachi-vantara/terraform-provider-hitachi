package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	"fmt"
	"sort"
	"strconv"
	"strings"

	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// "terraform-provider-hitachi/hitachi/common/utils"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetIscsiTarget(d *schema.ResourceData) (*terraformmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	portID := d.Get("port_id").(string)
	itNum := d.Get("iscsi_target_number").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_BEGIN), portID, itNum)
	it, err := reconObj.GetIscsiTarget(portID, itNum)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED), portID, itNum)
		newerr := fmt.Errorf(mc.GetMessage(mc.ERR_GET_ISCSITARGET_FAILED) + "\n%s", portID, itNum, err.Error())
		return nil, newerr
	}

	terraformModelIscsiTarget := terraformmodel.IscsiTarget{}
	err = copier.Copy(&terraformModelIscsiTarget, it)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_END), portID, itNum)

	return &terraformModelIscsiTarget, nil
}

func GetAllIscsiTargets(d *schema.ResourceData) (*terraformmodel.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), setting.Serial)
	iscsiTargets, err := reconObj.GetAllIscsiTargets()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_ISCSITARGET_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelIscsiTargets := terraformmodel.IscsiTargets{}
	err = copier.Copy(&terraformModelIscsiTargets, iscsiTargets)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_END), setting.Serial)

	return &terraformModelIscsiTargets, nil
}

func GetIscsiTargetsByPortIds(d *schema.ResourceData) (*terraformmodel.IscsiTargets, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	port_ids := d.Get("port_ids").([]interface{})

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}
	port_id_recon := []string{}
	for _, id := range port_ids {
		port_id_recon = append(port_id_recon, id.(string))
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_BEGIN), setting.Serial)
	iscsiTargets, err := reconObj.GetIscsiTargetsByPortIds(port_id_recon)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_ISCSITARGET_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelIscsiTargets := terraformmodel.IscsiTargets{}
	err = copier.Copy(&terraformModelIscsiTargets, iscsiTargets)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_ISCSITARGET_END), setting.Serial)

	return &terraformModelIscsiTargets, nil
}

func CreateIscsiTarget(d *schema.ResourceData) (*terraformmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateIscsiTargetRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_BEGIN), createInput.PortID, createInput.IscsiTargetNumber)
	reconcilerCreateIscsiTargetRequest := reconcilermodel.CreateIscsiTargetReq{}
	err = copier.Copy(&reconcilerCreateIscsiTargetRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	it, err := reconObj.ReconcileIscsiTarget(&reconcilerCreateIscsiTargetRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_FAILED), createInput.PortID, createInput.IscsiTargetNumber)
		log.WriteDebug("TFError| error in Creating IscsiTarget - ReconcileIscsiTarget , err: %v", err)
		return nil, err
	}

	terraformModelIscsiTarget := terraformmodel.IscsiTarget{}
	err = copier.Copy(&terraformModelIscsiTarget, it)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_END), terraformModelIscsiTarget.PortID, terraformModelIscsiTarget.IscsiTargetNumber)
	return &terraformModelIscsiTarget, nil
}

func ConvertIscsiTargetToSchema(iscsiTarget *terraformmodel.IscsiTarget, serial int) *map[string]interface{} {
	hg := map[string]interface{}{
		"storage_serial_number": serial,
		"iscsi_target_number":   iscsiTarget.IscsiTargetNumber,
		"port_id":               iscsiTarget.PortID,
		"iscsi_target_id":       iscsiTarget.IscsiTargetID,
		"iscsi_target_alias":    iscsiTarget.IscsiTargetName,
		"iscsi_target_name":     iscsiTarget.IscsiTargetNameIqn,
		"host_mode":             HostModeRestToUserConversion[iscsiTarget.HostMode],
		"host_mode_options":     iscsiTarget.HostModeOptions,
	}

	initiatorsDetail := []map[string]interface{}{}
	for _, hw := range iscsiTarget.Initiators {
		w := map[string]interface{}{
			"initiator_name":     hw.IscsiTargetNameIqn,
			"initiator_nickname": hw.IscsiNickname,
		}
		initiatorsDetail = append(initiatorsDetail, w)
	}
	hg["initiators"] = initiatorsDetail

	lupaths := []map[string]interface{}{}
	hgluns := []int{}
	ldevs := []int{}
	for _, lp := range iscsiTarget.LuPaths {
		p := map[string]interface{}{
			"lun_id":  lp.Lun,
			"ldev_id": lp.LdevID,
		}
		lupaths = append(lupaths, p)
		hgluns = append(hgluns, lp.Lun)
		ldevs = append(ldevs, lp.LdevID)
	}
	hg["lun_paths"] = lupaths
	sort.Ints(hgluns)
	sort.Ints(ldevs)
	hg["luns"] = hgluns
	hg["ldevs"] = ldevs

	return &hg
}

func ConvertSimpleIscsiTargetToSchema(iscsiTarget *terraformmodel.IscsiTarget, serial int) *map[string]interface{} {
	hg := map[string]interface{}{
		"storage_serial_number": serial,
		"iscsi_target_number":   iscsiTarget.IscsiTargetNumber,
		"port_id":               iscsiTarget.PortID,
		"iscsi_target_id":       iscsiTarget.IscsiTargetID,
		"iscsi_target_alias":    iscsiTarget.IscsiTargetName,
		"host_mode":             HostModeRestToUserConversion[iscsiTarget.HostMode],
	}

	return &hg
}

// UpdateIscsiTarget used to update modified data of hostgroup
func UpdateIscsiTarget(d *schema.ResourceData) (*terraformmodel.IscsiTarget, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}

	setting := reconcilermodel.StorageDeviceSettings{
		Serial:   storageSetting.Serial,
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		MgmtIP:   storageSetting.MgmtIP,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	updateInput, err := CreateIscsiTargetRequestFromSchema(d)
	if err != nil {
		log.WriteDebug("TFError| error in CreateIscsiTargetRequestFromSchema, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_ISCSITARGET_BEGIN), updateInput.PortID, updateInput.IscsiTargetNumber)
	reconcilerUpdateIscsiTargetRequest := reconcilermodel.CreateIscsiTargetReq{}
	err = copier.Copy(&reconcilerUpdateIscsiTargetRequest, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	hg, err := reconObj.ReconcileIscsiTarget(&reconcilerUpdateIscsiTargetRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_ISCSITARGET_FAILED), updateInput.PortID, *updateInput.IscsiTargetNumber)
		log.WriteDebug("TFError| error in Updating Hostgroup - ReconcileHostGroup , err: %v", err)
		return nil, err
	}

	terraformModelIscsiTarget := terraformmodel.IscsiTarget{}
	err = copier.Copy(&terraformModelIscsiTarget, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_ISCSITARGET_END), updateInput.PortID, *updateInput.IscsiTargetNumber)
	return &terraformModelIscsiTarget, nil
}

func CreateIscsiTargetRequestFromSchema(d *schema.ResourceData) (*terraformmodel.CreateIscsiTargetReq, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateIscsiTargetReq{}

	portId, ok := d.GetOk("port_id")
	if ok {
		pid := portId.(string)
		createInput.PortID = pid
	}

	hgnum, ok := d.GetOk("iscsi_target_number")
	if ok {
		hid := hgnum.(int)
		createInput.IscsiTargetNumber = &hid
	}

	hgname, ok := d.GetOk("iscsi_target_alias")
	if ok {
		name := hgname.(string)
		createInput.IscsiTargetName = name
	}

	itNameIqn, ok := d.GetOk("iscsi_target_name")
	if ok {
		name := itNameIqn.(string)
		createInput.IscsiTargetNameIqn = &name
	}

	hostmode, ok := d.GetOk("host_mode")
	if ok {
		userhmode := hostmode.(string)
		hmode := HostModeUserToRestConversion[strings.ToLower(userhmode)]
		createInput.HostMode = &hmode
	}
	hostoptions, ok := d.GetOk("host_mode_options")
	if ok {
		hopt := hostoptions.([]interface{})
		hostModeOptions := make([]int, len(hopt))
		for index, value := range hopt {
			switch typedValue := value.(type) {
			case int:
				hostModeOptions[index] = typedValue
			}
		}
		createInput.HostModeOptions = &hostModeOptions
	}

	volumes, ok := d.GetOk("lun")
	if ok {
		vols := volumes.(*schema.Set).List()
		ldevs := []terraformmodel.IscsiLuns{}
		for _, ldev := range vols {
			w := ldev.(map[string]interface{})
			ldevValue := w["ldev_id"].(int)
			lunValue := w["lun_id"].(int)
			ilun := terraformmodel.IscsiLuns{
				LdevId: &ldevValue,
				Lun:    &lunValue,
			}
			ldevs = append(ldevs, ilun)
		}
		createInput.Ldevs = ldevs
	}

	initiators, ok := d.GetOk("initiator")
	if ok {
		aiws := initiators.(*schema.Set).List()
		inpws := []terraformmodel.Initiator{}
		for _, wi := range aiws {
			w := wi.(map[string]interface{})
			iws := terraformmodel.Initiator{
				IscsiTargetNameIqn: w["initiator_name"].(string),
				IscsiNickname:      w["initiator_nickname"].(string),
			}
			inpws = append(inpws, iws)
		}
		createInput.Initiators = inpws
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func DeleteIscsiTarget(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	portId, ok := d.GetOk("port_id")
	log.WriteDebug("portId: %+v", portId)

	if !ok {
		iscsitarget, ok := d.GetOk("iscsitarget")
		if !ok {
			return fmt.Errorf("no iscsitarget data in resource")
		}
		log.WriteDebug("iscsitarget: %+v", iscsitarget.([]map[string]interface{})[0])
		portId, ok = iscsitarget.([]map[string]interface{})[0]["portId"]
		if !ok {
			return fmt.Errorf("found no portId in iscsitarget")
		}
		log.WriteDebug("iscsitarget portId: %+v", portId)
	}
	portID := portId.(string)

	itNum, ok := d.GetOk("iscsi_target_number")
	if !ok {
		iscsitarget, ok := d.GetOk("iscsitarget")
		if !ok {
			return fmt.Errorf("no iscsitarget data in resource")
		}
		log.WriteDebug("iscsitarget: %+v", iscsitarget.([]map[string]interface{})[0])
		itNum, ok = iscsitarget.([]map[string]interface{})[0]["iscsi_target_number"]
		if !ok {
			return fmt.Errorf("found no itNum in iscsitarget")
		}
		log.WriteDebug("iscsitarget itNum: %+v", itNum)
	}
	itNumber := itNum.(int)

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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_BEGIN), portID, itNumber)
	err = reconObj.DeleteIscsiTarget(portID, itNumber)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_FAILED), portID, itNumber)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_END), portID, itNumber)

	return nil
}
