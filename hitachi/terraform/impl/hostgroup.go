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
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

var HostModeUserToRestConversion = map[string]string{
	"standard":          "LINUX/IRIX",
	"vmware":            "VMWARE", //Deprecated
	"hp":                "HP-UX",
	"openvms":           "OVMS",
	"tru64":             "TRU64",
	"solaris":           "SOLARIS",
	"netware":           "NETWARE",
	"windows":           "WIN", //Deprecated
	"aix":               "AIX",
	"vmware extension":  "VMWARE_EX",
	"windows extension": "WIN_EX",
}

var HostModeRestToUserConversion = map[string]string{
	"LINUX/IRIX": "Standard",
	"VMWARE":     "VMware", //Deprecated
	"HP-UX":      "HP",
	"OVMS":       "OpenVMS",
	"TRU64":      "Tru64",
	"SOLARIS":    "Solaris",
	"NETWARE":    "NetWare",
	"WIN":        "Windows", //Deprecated
	"AIX":        "AIX",
	"VMWARE_EX":  "VMware Extension",
	"WIN_EX":     "Windows Extension",
}

func GetHostGroup(d *schema.ResourceData) (*terraformmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	portID := d.Get("port_id").(string)
	hgNum := d.Get("hostgroup_number").(int)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_BEGIN), portID, hgNum)
	hg, err := reconObj.GetHostGroup(portID, hgNum)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_HOSTGROUP_FAILED), portID, hgNum)
		return nil, err
	}

	terraformModelHostGroup := terraformmodel.HostGroup{}
	err = copier.Copy(&terraformModelHostGroup, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOSTGROUP_END), portID, hgNum)

	return &terraformModelHostGroup, nil
}

func GetAllHostGroups(d *schema.ResourceData) (*terraformmodel.HostGroups, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_BEGIN), setting.Serial)
	hostGroups, err := reconObj.GetAllHostGroups()
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelHostGroups := terraformmodel.HostGroups{}
	err = copier.Copy(&terraformModelHostGroups, hostGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_END), setting.Serial)

	return &terraformModelHostGroups, nil
}

func GetHostGroupsByPortIds(d *schema.ResourceData) (*terraformmodel.HostGroups, error) {
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
	portIdRecon := []string{}
	for _, id := range port_ids {
		portIdRecon = append(portIdRecon, id.(string))
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_BEGIN), setting.Serial)
	hostGroups, err := reconObj.GetHostGroupsByPortIds(portIdRecon)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOSTGROUP_FAILED), setting.Serial)
		return nil, err
	}

	terraformModelHostGroups := terraformmodel.HostGroups{}
	err = copier.Copy(&terraformModelHostGroups, hostGroups)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOSTGROUP_END), setting.Serial)

	return &terraformModelHostGroups, nil
}

func CreateHostGroup(d *schema.ResourceData) (*terraformmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	storageSetting, err := cache.GetSanSettingsFromCache(strconv.Itoa(serial))
	if err != nil {
		return nil, err
	}
	// pportID, phgnum := CheckSchemaIfHostGroupGet(d)
	// if pportID != nil && phgnum != nil {
	// 	hg, err := sanprov.GetHostGroup(*storageSetting, *pportID, *phgnum)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if hg.ByteFormatCapacity == "" {
	// 		// does not exist, or in the process of being deleted
	// 		return nil, fmt.Errorf("Volume does not exist")
	// 	}
	// 	return hg, nil
	// }

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

	createInput, err := CreateHostGroupRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_BEGIN), createInput.PortID, createInput.HostGroupNumber)
	reconcilerCreateHostGroupRequest := reconcilermodel.CreateHostGroupRequest{}
	err = copier.Copy(&reconcilerCreateHostGroupRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	hg, err := reconObj.ReconcileHostGroup(&reconcilerCreateHostGroupRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_HOSTGROUP_FAILED), createInput.PortID, createInput.HostGroupNumber)
		log.WriteDebug("TFError| error in Creating Hostgroup - ReconcileHostGroup , err: %v", err)
		return nil, err
	}

	terraformModelHostGroup := terraformmodel.HostGroup{}
	err = copier.Copy(&terraformModelHostGroup, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_END), terraformModelHostGroup.PortID, terraformModelHostGroup.HostGroupNumber)
	return &terraformModelHostGroup, nil
}

func CreateHostGroupRequestFromSchema(d *schema.ResourceData) (*terraformmodel.CreateHostGroupRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateHostGroupRequest{}

	portId, ok := d.GetOk("port_id")
	if ok {
		pid := portId.(string)
		createInput.PortID = &pid
	}

	hgnum, ok := d.GetOk("hostgroup_number")
	if ok {
		hid := hgnum.(int)
		createInput.HostGroupNumber = &hid
	}

	hgname, ok := d.GetOk("hostgroup_name")
	if ok {
		name := hgname.(string)
		createInput.HostGroupName = &name
	}

	hostmode, ok := d.GetOk("host_mode")
	if ok {
		userhmode := hostmode.(string)
		hmode := HostModeUserToRestConversion[strings.ToLower(userhmode)]
		if hmode == "" {
			err := fmt.Errorf("invalid hostmode specified %v", userhmode)
			return nil, err

		}
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
		createInput.HostModeOptions = hostModeOptions
	}

	volumes, ok := d.GetOk("lun")
	if ok {
		vols := volumes.(*schema.Set).List()
		ldevs := []terraformmodel.Luns{}
		for _, ldev := range vols {
			w := ldev.(map[string]interface{})
			ldevValue := w["ldev_id"].(int)
			lunValue := w["lun"].(int)
			ilun := terraformmodel.Luns{
				LdevId: &ldevValue,
				Lun:    &lunValue,
			}
			ldevs = append(ldevs, ilun)
		}
		createInput.Ldevs = ldevs
	}

	wwns, ok := d.GetOk("wwn")
	if ok {
		aiws := wwns.(*schema.Set).List()
		inpws := []terraformmodel.Wwn{}
		for _, wi := range aiws {
			w := wi.(map[string]interface{})
			iws := terraformmodel.Wwn{
				Wwn:  w["host_wwn"].(string),
				Name: w["wwn_nickname"].(string),
			}
			inpws = append(inpws, iws)
		}
		createInput.Wwns = inpws
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func CheckSchemaIfHostGroupGet(d *schema.ResourceData) *int {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	fields := []string{
		"size_gb",
		"name",
		//"dedup_mode",
		"pool_id",
		"paritygroup_id",
	}

	for _, f := range fields {
		if _, ok := d.GetOk(f); ok {
			return nil
		}
	}

	ldevId, ok := d.GetOk("ldev_id")
	if ok {
		lid := ldevId.(int)
		return &lid
	}
	return nil
}

func DeleteHostGroup(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	portId, ok := d.GetOk("port_id")
	log.WriteDebug("portId: %+v", portId)

	if !ok {
		hostgroup, ok := d.GetOk("hostgroup")
		if !ok {
			return fmt.Errorf("no hostgroup data in resource")
		}
		log.WriteDebug("hostgroup: %+v", hostgroup.([]map[string]interface{})[0])
		portId, ok = hostgroup.([]map[string]interface{})[0]["portId"]
		if !ok {
			return fmt.Errorf("found no portId in hostgroup")
		}
		log.WriteDebug("hostgroup ldevID: %+v", portId)
	}
	portID := portId.(string)

	hgNum, ok := d.GetOk("hostgroup_number")
	log.WriteDebug("hostgroup_number: %+v", hgNum)

	if !ok {
		hostgroup, ok := d.GetOk("hostgroup")
		if !ok {
			return fmt.Errorf("no hostgroup data in resource")
		}
		log.WriteDebug("hostgroup: %+v", hostgroup.([]map[string]interface{})[0])
		hgNum, ok = hostgroup.([]map[string]interface{})[0]["hostgroup_number"]
		if !ok {
			return fmt.Errorf("found no hgNum in hostgroup")
		}
		log.WriteDebug("hostgroup hgNum: %+v", hgNum)
	}
	hgNumber := hgNum.(int)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_BEGIN), portID, hgNumber)
	err = reconObj.DeleteHostGroup(portID, hgNumber)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_HOSTGROUP_FAILED), portID, hgNumber)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOSTGROUP_END), portID, hgNumber)

	return nil
}

func ConvertSimpleHostGroupToSchema(hostgroup *terraformmodel.HostGroup, serial int) *map[string]interface{} {
	hg := map[string]interface{}{
		"storage_serial_number": serial,
		"hostgroup_number":      hostgroup.HostGroupNumber,
		"port_id":               hostgroup.PortID,
		"hostgroup_name":        hostgroup.HostGroupName,
		"host_mode":             HostModeRestToUserConversion[hostgroup.HostMode],
	}

	return &hg
}

func ConvertHostGroupToSchema(hostgroup *terraformmodel.HostGroup, serial int) *map[string]interface{} {
	hg := map[string]interface{}{
		"storage_serial_number": serial,
		"hostgroup_number":      hostgroup.HostGroupNumber,
		"port_id":               hostgroup.PortID,
		"hostgroup_name":        hostgroup.HostGroupName,
		"host_mode":             HostModeRestToUserConversion[hostgroup.HostMode],
		"host_mode_options":     hostgroup.HostModeOptions,
	}

	wwnsDetail := []map[string]interface{}{}
	wwns := []string{}
	for _, hw := range hostgroup.WwnDetails {
		w := map[string]interface{}{
			"wwn":  hw.Wwn,
			"name": hw.Name,
		}
		wwnsDetail = append(wwnsDetail, w)
		wwns = append(wwns, hw.Wwn)
	}
	hg["wwns_detail"] = wwnsDetail
	sort.Strings(wwns)
	hg["wwns"] = wwns

	lupaths := []map[string]interface{}{}
	hgluns := []int{}
	ldevs := []int{}
	for _, lp := range hostgroup.LuPaths {
		p := map[string]interface{}{
			"hg_lun_id": lp.Lun,
			"ldev_id":   lp.LdevID,
		}
		lupaths = append(lupaths, p)
		hgluns = append(hgluns, lp.Lun)
		ldevs = append(ldevs, lp.LdevID)
	}
	hg["lun_paths"] = lupaths
	sort.Ints(hgluns)
	sort.Ints(ldevs)
	hg["hg_luns"] = hgluns
	hg["ldevs"] = ldevs

	return &hg
}

// UpdateHostGroup used to update modified data of hostgroup
func UpdateHostGroup(d *schema.ResourceData) (*terraformmodel.HostGroup, error) {
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

	updateInput, err := CreateHostGroupRequestFromSchema(d)
	if err != nil {
		log.WriteDebug("TFError| error in CreateHostGroupRequestFromSchema, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_HOSTGROUP_BEGIN), updateInput.PortID, updateInput.HostGroupNumber)
	reconcilerUpdateHostGroupRequest := reconcilermodel.CreateHostGroupRequest{}
	err = copier.Copy(&reconcilerUpdateHostGroupRequest, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	hg, err := reconObj.ReconcileHostGroup(&reconcilerUpdateHostGroupRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_HOSTGROUP_FAILED), updateInput.PortID, updateInput.HostGroupNumber)
		log.WriteDebug("TFError| error in Updating Hostgroup - ReconcileHostGroup , err: %v", err)
		return nil, err
	}

	terraformModelHostGroup := terraformmodel.HostGroup{}
	err = copier.Copy(&terraformModelHostGroup, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_HOSTGROUP_END), updateInput.PortID, updateInput.HostGroupNumber)
	return &terraformModelHostGroup, nil
}
