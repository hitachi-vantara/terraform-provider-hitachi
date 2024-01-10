package terraform

import (
	"errors"
	"strconv"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	common "terraform-provider-hitachi/hitachi/terraform/common"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetInfraGwIscsiTargets(d *schema.ResourceData) (*[]terraformmodel.InfraIscsiTargetInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	if serial == "" && storageId == "" {
		err := errors.New("both serial and storage_id can't be empty. Please specify one")
		return nil, err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
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
	port := d.Get("port_id").(string)

	iscsi_name := d.Get("iscsi_name").(string)
	iscsi_id := -1
	iid, okId := d.GetOk("iscsi_target_number")
	if okId {
		iscsi_id = iid.(int)
	}

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetIscsiTargets(storageId, port)
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraGwIscsiTargets, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GW_GET_ISCSI_TARGETS_FAILED), setting.Address)
		return nil, err
	}

	var result model.IscsiTarget
	if iscsi_name != "" {
		for _, iscsi := range reconResponse.Data {
			if iscsi.ISCSIName == iscsi_name {
				result.Path = reconResponse.Path
				result.Message = reconResponse.Message
				result.Data = iscsi
				break
			}
		}
	}
	if iscsi_id != -1 {
		for _, iscsi := range reconResponse.Data {
			if iscsi.ISCSIId == iscsi_id {
				result.Path = reconResponse.Path
				result.Message = reconResponse.Message
				result.Data = iscsi
				break
			}
		}
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraIscsiTargets{}

	if iscsi_name != "" || iscsi_id != -1 {
		err = copier.Copy(&terraformResponse, &result)
	} else {
		err = copier.Copy(&terraformResponse, reconResponse)
	}
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_END), setting.Address)

	return &terraformResponse.Data, nil
}

func GetInfraIscsiTargetsById(d *schema.ResourceData) (*[]terraformmodel.InfraIscsiTargetInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	if serial == "" && storageId == "" {
		err := errors.New("both serial and storage_id can't be empty. Please specify one")
		return nil, err
	}

	if serial != "" && storageId != "" {
		err := errors.New("both serial and storage_id are not allowed. Either serial or storage_id can be specified")
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

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetIscsiTargets(storageId, "")
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraIscsiTargetsById, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GW_GET_ISCSI_TARGETS_FAILED), setting.Address)
		return nil, err
	}

	var result model.IscsiTargets
	if len(portIdsMap) > 0 {
		result.Path = reconResponse.Path
		result.Message = reconResponse.Message
		for _, p := range reconResponse.Data {
			_, ok := portIdsMap[p.PortId]
			if ok {
				result.Data = append(result.Data, p)
			}
		}
	}

	terraformResponse := terraformmodel.InfraIscsiTargets{}
	if len(portIdsMap) > 0 {
		err = copier.Copy(&terraformResponse, &result)
	} else {
		err = copier.Copy(&terraformResponse, reconResponse)
	}
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_END), setting.Address)

	return &terraformResponse.Data, nil
}

func GetInfraGwIscsiTarget(d *schema.ResourceData, id string) (*[]terraformmodel.InfraIscsiTargetInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}
	storage_id := d.Get("storage_id").(string)

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	setting := model.InfraGwSettings{
		Username: storageSetting.Username,
		Password: storageSetting.Password,
		Address:  storageSetting.Address,
	}

	reconObj, err := reconimpl.NewEx(setting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_BEGIN), setting.Address)
	reconResponse, err := reconObj.GetIscsiTarget(storage_id, id)
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraGwIscsiTarget, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GW_GET_ISCSI_TARGETS_FAILED), setting.Address)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraIscsiTargets{}
	err = copier.Copy(&terraformResponse, reconResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GW_GET_ISCSI_TARGETS_END), setting.Address)

	return &terraformResponse.Data, nil
}

func ConvertInfraIscsiTargetToSchema(pg *terraformmodel.InfraIscsiTargetInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"storage_serial_number": storage_serial_number,
		"resource_id":           pg.ResourceId,
		"port_id":               pg.PortId,
		"resource_group_id":     pg.ResourceGroupId,
		"target_user":           pg.TargetUser,
		"iqn":                   pg.Iqn,
		"iqn_initiators":        pg.IqnInitiators,
		"chap_users":            pg.ChapUsers,
		"iscsi_name":            pg.ISCSIName,
		"iscsi_id":              pg.ISCSIId,
	}

	hostMode := []map[string]interface{}{}
	hm := map[string]interface{}{
		"host_common_settings": pg.HostMode.HostCommonSettings,
		"host_middle_ware":     pg.HostMode.HostMiddleWare,
		"host_mode":            pg.HostMode.HostMode,
		"host_platform_option": pg.HostMode.HostPlatformOption,
		"is_df":                pg.HostMode.IsDF,
		"is_raid":              pg.HostMode.IsRAID,
		"raid_host_mode_char":  pg.HostMode.RaidHostModeChar,
	}

	hModeOptions := []map[string]interface{}{}
	for _, hmo := range pg.HostMode.HostModeOptions {
		p := map[string]interface{}{
			"df_option":          hmo.DfOption,
			"is_ams_legal":       hmo.IsAMSLegal,
			"is_df":              hmo.IsDF,
			"is_hus_legal":       hmo.IsHUSLegal,
			"is_raid":            hmo.IsRAID,
			"raid_option":        hmo.RaidOption,
			"raid_option_number": hmo.RaidOptionNumber,
		}
		hModeOptions = append(hModeOptions, p)
	}
	hm["host_mode_options"] = hModeOptions

	hostMode = append(hostMode, hm)
	sp["host_mode"] = hostMode

	logicalUnits := []map[string]interface{}{}
	for _, lu := range pg.LogicalUnits {
		p := map[string]interface{}{
			"host_lun_id":                lu.HostLunId,
			"logical_unit_id":            lu.LogicalUnitId,
			"logical_unit_id_hex_format": lu.LogicalUnitIdHexFormat,
		}
		logicalUnits = append(logicalUnits, p)
	}
	sp["logical_units"] = logicalUnits

	authParams := []map[string]interface{}{}
	a := map[string]interface{}{
		"is_chap_enabled":     pg.AuthParam.IsChapEnabled,
		"is_chap_required":    pg.AuthParam.IsChapRequired,
		"is_mutual_auth":      pg.AuthParam.IsMutualAuth,
		"authentication_mode": pg.AuthParam.AuthenticationMode,
	}
	authParams = append(authParams, a)
	sp["auth_params"] = authParams

	return &sp
}
