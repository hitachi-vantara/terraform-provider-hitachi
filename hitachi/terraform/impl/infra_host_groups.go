package terraform

import (
	"fmt"
	"strconv"
	"strings"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	//reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

var storage_serial_number int

var InfraHostModeUserToRestConversion = map[string]string{
	"aix":               "AIX",
	"dynix":             "DYNIX",
	"hi_ux":             "HI_UX",
	"hp":                "HP",
	"hp_xp":             "HP_XP",
	"linux":             "LINUX",
	"netware":           "NETWARE",
	"openvms":           "OPEN_VMS",
	"solaris":           "SOLARIS",
	"standard":          "STANDARD",
	"tru64":             "TRU64",
	"uvm":               "UVM",
	"vmware":            "VMWARE", //Deprecated
	"vmware extension":  "VMWARE_EXTENSION",
	"windows":           "WINDOWS", //Deprecated
	"windows extension": "WIN_EXTENSION",
}

var InfraHostModeRestToUserConversion = map[string]string{
	"AIX":              "AIX",
	"DYNIX":            "DYNIX",
	"HI_UX":            "HI_UX",
	"HP":               "HP",
	"HP_XP":            "HP_XP",
	"LINUX":            "Linux",
	"NETWARE":          "NetWare",
	"OPEN_VMS":         "OpenVMS",
	"SOLARIS":          "Solaris",
	"STANDARD":         "Standard",
	"TRU64":            "Tru64",
	"UVM":              "UVM",
	"VMWARE":           "VMware", //Deprecated
	"VMWARE_EXTENSION": "VMware Extension",
	"WIN":              "Windows", //Deprecated
	"WIN_EXTENSION":    "Windows Extension",
}

func GetInfraHostGroups(d *schema.ResourceData) (*[]terraformmodel.InfraHostGroupInfo, *[]terraformmodel.InfraMTHostGroupInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return nil, nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, nil, err
		}
		d.Set("storage_id", storageId)
	}

	if serial == "" {
		serial, err = common.GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	}
	d.Set("serial", storage_serial_number)

	port := d.Get("port_id").(string)
	hostgroup_name := d.Get("hostgroup_name").(string)
	hostgroup_id := -1
	hid, okId := d.GetOk("hostgroup_number")
	if okId {
		hostgroup_id = hid.(int)
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, nil, err
	}

	reconObj, err := reconimpl.NewEx(*storageSetting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_HOST_GROUPS_BEGIN), storageSetting.Address)
	if storageSetting.PartnerId == nil {
		reconResponse, err := reconObj.GetHostGroups(storageId, port)
		if err != nil {
			log.WriteDebug("TFError| error getting GetInfraHostGroups, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_HOST_GROUPS_FAILED), storageSetting.Address)
			return nil, nil, err
		}

		var result model.HostGroup
		if hostgroup_name != "" {
			found := false
			for _, hg := range reconResponse.Data {
				if hg.HostGroupName == hostgroup_name {
					result.Path = reconResponse.Path
					result.Message = reconResponse.Message
					result.Data = hg
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("hostgroup name  %s not found", hostgroup_name)
				log.WriteDebug("Hostgroup name  %s not found", hostgroup_name)
				return nil, nil, err
			}
		}
		if hostgroup_id != -1 {
			found := false
			for _, hg := range reconResponse.Data {
				if hg.HostGroupId == hostgroup_id {
					result.Path = reconResponse.Path
					result.Message = reconResponse.Message
					result.Data = hg
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("hostgroup number  %d not found", hostgroup_id)
				log.WriteDebug("Hostgroup number  %d not found", hostgroup_id)
				return nil, nil, err
			}
		}
		// Converting reconciler to terraform
		terraformResponse := terraformmodel.InfraHostGroups{}

		if hostgroup_name != "" || hostgroup_id != -1 {
			err = copier.Copy(&terraformResponse, &result)
		} else {
			err = copier.Copy(&terraformResponse, reconResponse)
		}

		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_HOST_GROUPS_END), storageSetting.Address)

		return &terraformResponse.Data, nil, nil
	}

	mtResponse, err := reconObj.GetMTStorageDevices()
	if err != nil {
		log.WriteDebug("TFError| error getting GetMTStorageDevices, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_DEVICES_FAILED), storageSetting.Address)
		return nil, nil, err
	}

	terraformMtResponse := terraformmodel.InfraMTHostGroups{}
	err = copier.Copy(&terraformMtResponse.Data, mtResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_END), storageSetting.Address)

	return nil, &terraformMtResponse.Data, nil
}

func GetInfraHostGroupsByPortIds(d *schema.ResourceData) (*[]terraformmodel.InfraHostGroupInfo, *[]terraformmodel.InfraMTHostGroupInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return nil, nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, nil, err
		}
		d.Set("storage_id", storageId)
	}

	if serial == "" {
		serial, err = common.GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, nil, err
		}
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, nil, err
	}

	reconObj, err := reconimpl.NewEx(*storageSetting)

	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, nil, err
	}

	portIdsMap := map[string]string{}
	ids, ok := d.GetOk("port_ids")
	if ok {
		pgIds := ids.([]interface{})
		for _, value := range pgIds {
			switch typedValue := value.(type) {
			case string:
				portIdsMap[typedValue] = typedValue
			}
		}
	}

	log.WriteDebug("TFDebug| port group filter will be apply %v size %v", portIdsMap, len(portIdsMap))

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_HOST_GROUPS_BEGIN), storageSetting.Address)

	if storageSetting.PartnerId == nil {
		reconResponse, err := reconObj.GetHostGroups(storageId, "")
		if err != nil {
			log.WriteDebug("TFError| error getting GetInfraHostGroupsByPortIds, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_PARITY_GROUPS_FAILED), storageSetting.Address)
			return nil, nil, err
		}

		var result model.HostGroups
		if len(portIdsMap) > 0 {
			result.Path = reconResponse.Path
			result.Message = reconResponse.Message
			for _, p := range reconResponse.Data {
				_, ok := portIdsMap[p.Port]
				if ok {
					result.Data = append(result.Data, p)
				}
			}
		}
		terraformResponse := terraformmodel.InfraHostGroups{}
		if len(portIdsMap) > 0 {
			err = copier.Copy(&terraformResponse, &result)
		} else {
			err = copier.Copy(&terraformResponse, reconResponse)
		}
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_HOST_GROUPS_END), storageSetting.Address)

		return &terraformResponse.Data, nil, nil
	}
	mtResponse, err := reconObj.GetHostGroupsByPartnerIdOrSubscriberID(storageId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetMTStorageDevices, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_HOST_GROUPS_FAILED), storageSetting.Address)
		return nil, nil, err
	}

	terraformMtResponse := terraformmodel.InfraMTHostGroups{}
	err = copier.Copy(&terraformMtResponse.Data, mtResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_HOST_GROUPS_END), storageSetting.Address)

	return nil, &terraformMtResponse.Data, nil
}

func CreateInfraHostGroup(d *schema.ResourceData) (*[]terraformmodel.InfraHostGroupInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, err
		}
		d.Set("storage_id", storageId)
	}

	if serial == "" {
		serial, err = common.GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	reconObj, err := reconimpl.NewEx(*storageSetting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraHostGroupRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_BEGIN), createInput.Port, createInput.HostGroupName)
	reconcilerCreateHostGroupRequest := model.CreateHostGroupParam{}
	err = copier.Copy(&reconcilerCreateHostGroupRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	hg, err := reconObj.ReconcileHostGroup(storageId, &reconcilerCreateHostGroupRequest)

	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_HOSTGROUP_FAILED), createInput.Port, createInput.HostGroupName)
		log.WriteDebug("TFError| error in Creating Hostgroup - ReconcileHostGroup , err: %v", err)
		return nil, err
	}

	terraformModelHostGroup := terraformmodel.InfraHostGroups{}
	err = copier.Copy(&terraformModelHostGroup, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_END), terraformModelHostGroup.Data[0].Port, terraformModelHostGroup.Data[0].HostGroupName)
	return &terraformModelHostGroup.Data, nil
}

func CreateInfraHostGroupRequestFromSchema(d *schema.ResourceData) (*terraformmodel.CreateInfraHostGroupParam, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateInfraHostGroupParam{}

	portId, ok := d.GetOk("port_id")
	if ok {
		pid := portId.(string)
		createInput.Port = pid
	}

	hgname, ok := d.GetOk("hostgroup_name")
	if ok {
		name := hgname.(string)
		createInput.HostGroupName = name
	}

	hostmode, ok := d.GetOk("host_mode")
	if ok {
		userhmode := hostmode.(string)
		hmode := InfraHostModeUserToRestConversion[strings.ToLower(userhmode)]
		if hmode == "" {
			err := fmt.Errorf("invalid hostmode specified %v", userhmode)
			return nil, err

		}
		createInput.HostMode = hmode
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
		luns := make([]int, len(vols))
		for _, ldev := range vols {
			w := ldev.(map[string]interface{})
			ldevValue := w["ldev_id"].(int)
			lunValue := w["lun"].(int)
			ilun := terraformmodel.Luns{
				LdevId: &ldevValue,
				Lun:    &lunValue,
			}
			ldevs = append(ldevs, ilun)
			luns = append(luns, lunValue)
		}
		//createInput.Ldevs = ldevs
		createInput.Luns = luns
	}

	wwns, ok := d.GetOk("wwn")
	if ok {
		aiws := wwns.(*schema.Set).List()
		inpws := []model.Wwn{}
		for _, wi := range aiws {
			w := wi.(map[string]interface{})
			iws := model.Wwn{
				Id:   w["host_wwn"].(string),
				Name: w["wwn_nickname"].(string),
			}
			inpws = append(inpws, iws)
		}
		createInput.Wwns = inpws
	}

	ucpSystem, ok := d.GetOk("system")
	if ok {
		us := ucpSystem.(string)
		createInput.UcpSystem = us
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func ConvertInfraMTHostGroupToSchema(pg *terraformmodel.InfraMTHostGroupInfo) *map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	sp := map[string]interface{}{
		"storage_serial_number": storage_serial_number,
		"resource_id":           pg.ResourceId,
		"type":                  pg.Type,
		"storage_id":            pg.StorageId,
		"device_id":             pg.DeviceId,
		"entitlement_status":    pg.EntitlementStatus,
		"partner_id":            pg.PartnerId,
		"subscriber_id":         pg.SubscriberId,
		"hostgroup_name":        pg.HostGroupInfo.HostGroupName,
		"host_group_number":     pg.HostGroupInfo.HostGroupId,
		"resource_group_id":     pg.HostGroupInfo.ResourceGroupId,
		"port_id":               pg.HostGroupInfo.Port,
		"host_mode":             InfraHostModeRestToUserConversion[pg.HostGroupInfo.HostMode],
	}
	return &sp
}
func ConvertInfraHostGroupToSchema(pg *terraformmodel.InfraHostGroupInfo) *map[string]interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	sp := map[string]interface{}{
		"storage_serial_number": storage_serial_number,
		"resource_id":           pg.ResourceId,
		"hostgroup_name":        pg.HostGroupName,
		"hostgroup_number":      pg.HostGroupId,
		"resource_group_id":     pg.ResourceGroupId,
		"port_id":               pg.Port,
		"host_mode":             InfraHostModeRestToUserConversion[pg.HostMode],
	}

	lupaths := []map[string]interface{}{}
	for _, lp := range pg.LunPaths {
		p := map[string]interface{}{
			"lun_id":  lp.LunId,
			"ldev_id": lp.LdevId,
		}
		lupaths = append(lupaths, p)
	}
	sp["lun_paths"] = lupaths

	wwns := []map[string]interface{}{}
	for _, wwn := range pg.Wwns {
		p := map[string]interface{}{
			"id":   wwn.Id,
			"name": wwn.Name,
		}
		wwns = append(wwns, p)
	}
	sp["wwns"] = wwns

	hModeOptions := []map[string]interface{}{}
	for _, hmo := range pg.HostModeOptions {
		p := map[string]interface{}{
			"host_mode_option":        hmo.HostModeOption,
			"host_mode_option_number": hmo.HostModeOptionNumber,
		}
		hModeOptions = append(hModeOptions, p)
	}
	sp["host_mode_options"] = hModeOptions

	return &sp
}

func UpdateInfraHostGroup(d *schema.ResourceData) (*[]terraformmodel.InfraHostGroupInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	err := common.ValidateSerialAndStorageId(serial, storageId)
	if err != nil {
		return nil, err
	}

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return nil, err
		}
		d.Set("storage_id", storageId)
	}

	if serial == "" {
		serial, err = common.GetSerialFromStorageId(address, storageId)
		if err != nil {
			return nil, err
		}
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	} else {
		storage_serial_number, err = strconv.Atoi(serial)
		if err != nil {
			return nil, err
		}
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	createInput, err := CreateInfraHostGroupRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_BEGIN), createInput.Port, createInput.HostGroupName)
	reconcilerCreateHostGroupRequest := model.CreateHostGroupParam{}
	err = copier.Copy(&reconcilerCreateHostGroupRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	hg, err := reconObj.ReconcileHostGroup(storageId, &reconcilerCreateHostGroupRequest)

	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_HOSTGROUP_FAILED), createInput.Port, createInput.HostGroupName)
		log.WriteDebug("TFError| error in Creating Hostgroup - ReconcileHostGroup , err: %v", err)
		return nil, err
	}

	terraformModelHostGroup := terraformmodel.InfraHostGroups{}
	err = copier.Copy(&terraformModelHostGroup, hg)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_HOSTGROUP_END), terraformModelHostGroup.Data[0].Port, terraformModelHostGroup.Data[0].HostGroupName)
	return &terraformModelHostGroup.Data, nil
}

func DeleteInfraHostGroup(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return err
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return err
	}

	setting := model.InfraGwSettings(*storageSetting)

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return err
	}
	storageId := d.Get("storage_id").(string)
	hgId := d.State().ID

	if storageSetting.PartnerId != nil {
		// Untag the Hostgroup from Subscriber first
		err = reconObj.DeleteMTHostGroup(storageId, hgId)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteStorageDevice, err: %v", err)
			return err
		}
	}

	err = reconObj.DeleteHostGroup(storageId, hgId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteStorageDevice, err: %v", err)
		return err
	}
	return nil
}
