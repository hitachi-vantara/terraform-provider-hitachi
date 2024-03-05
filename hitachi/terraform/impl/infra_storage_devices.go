package terraform

import (
	"strconv"
	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	common "terraform-provider-hitachi/hitachi/terraform/common"

	// mc "terraform-provider-hitachi/hitachi/messagecatalog"

	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"

	model "terraform-provider-hitachi/hitachi/infra_gw/model"
	reconimpl "terraform-provider-hitachi/hitachi/infra_gw/reconciler/impl"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

func GetInfraStorageDevices(d *schema.ResourceData) (*[]terraformmodel.InfraStorageDeviceInfo, *terraformmodel.InfraMTStorageDevices, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, nil, err
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

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_BEGIN), storageSetting.Address)

	if storageSetting.PartnerId == nil {
		reconResponse, err := reconObj.GetStorageDevices()
		if err != nil {
			log.WriteDebug("TFError| error getting GetInfraStorageDevices, err: %v", err)
			log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_DEVICES_FAILED), storageSetting.Address)
			return nil, nil, err
		}

		// Converting reconciler to terraform
		terraformResponse := terraformmodel.InfraStorageDevices{}
		err = copier.Copy(&terraformResponse, reconResponse)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, nil, err
		}
		log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_END), storageSetting.Address)

		log.WriteDebug("all: %+v\n", terraformResponse)
		log.WriteDebug("data: %+v\n", terraformResponse.Data)

		return &terraformResponse.Data, nil, nil
	}
	mtResponse, err := reconObj.GetMTStorageDevices()
	if err != nil {
		log.WriteDebug("TFError| error getting GetMTStorageDevices, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_DEVICES_FAILED), storageSetting.Address)
		return nil, nil, err
	}

	terraformMtResponse := terraformmodel.InfraMTStorageDevices{}
	err = copier.Copy(&terraformMtResponse, mtResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_END), storageSetting.Address)

	return nil, &terraformMtResponse, nil

}

func GetInfraStorageDevice(d *schema.ResourceData, serial string) (*[]terraformmodel.InfraStorageDeviceInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	id, err := common.GetStorageIdFromSerial(address, serial)
	if err != nil {
		return nil, err
	}
	d.Set("storage_id", id)

	storageSetting, err := cache.GetInfraSettingsFromCache(address)
	if err != nil {
		return nil, err
	}

	reconObj, err := reconimpl.NewEx(*storageSetting)
	if err != nil {
		log.WriteDebug("TFError| error in terraform NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_BEGIN), storageSetting.Address)
	reconResponse, err := reconObj.GetStorageDevice(id)
	if err != nil {
		log.WriteDebug("TFError| error getting GetInfraGwStorageDevice, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_INFRA_GET_STORAGE_DEVICES_FAILED), storageSetting.Address)
		return nil, err
	}

	// Converting reconciler to terraform
	terraformResponse := terraformmodel.InfraStorageDevices{}
	err = copier.Copy(&terraformResponse, reconResponse)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_INFRA_GET_STORAGE_DEVICES_END), storageSetting.Address)

	return &terraformResponse.Data, nil
}

func CreateInfraStorageDevice(d *schema.ResourceData) (*[]terraformmodel.InfraStorageDeviceInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return nil, err
	}

	storageId, _ := common.GetStorageIdFromSerial(address, serial)

	storage_serial_number, err = strconv.Atoi(serial)
	if err != nil {
		return nil, err
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

	createInput, err := CreateInfraStorageDeviceRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_INFRA_STORAGE_DEVICE_BEGIN), createInput.SerialNumber, createInput.ManagementAddress)
	reconcilerCreateStorageDeviceRequest := model.CreateStorageDeviceParam{}
	err = copier.Copy(&reconcilerCreateStorageDeviceRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	sd, err := reconObj.ReconcileStorageDevice(storageId, &reconcilerCreateStorageDeviceRequest)

	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_INFRA_STORAGE_DEVICE_FAILED), createInput.SerialNumber, createInput.ManagementAddress)
		log.WriteDebug("TFError| error in  ReconcileStorageDevice , err: %v", err)
		return nil, err
	}

	terraformModelResponse := terraformmodel.InfraStorageDevices{}
	err = copier.Copy(&terraformModelResponse, sd)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_INFRA_STORAGE_DEVICE_END), createInput.SerialNumber, createInput.ManagementAddress)
	return &terraformModelResponse.Data, nil
}

func DeleteInfraStorageDevice(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

	address, err := cache.GetCurrentAddress()
	if err != nil {
		return err
	}

	if storageId == "" {
		storageId, err = common.GetStorageIdFromSerial(address, serial)
		if err != nil {
			return err
		}
		d.Set("storage_id", storageId)
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

	//storageId = d.State().ID

	err = reconObj.DeleteStorageDevice(storageId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteStorageDevice, err: %v", err)
		return err
	}
	return nil
}

func UpdateInfraStorageDevice(d *schema.ResourceData) (*[]terraformmodel.InfraStorageDeviceInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := common.GetSerialString(d)
	storageId := d.Get("storage_id").(string)

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

	createInput, err := CreateInfraStorageDeviceRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_INFRA_STORAGE_DEVICE_BEGIN), createInput.SerialNumber, createInput.ManagementAddress)
	reconcilerCreateStorageDeviceRequest := model.CreateStorageDeviceParam{}
	err = copier.Copy(&reconcilerCreateStorageDeviceRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	sd, err := reconObj.ReconcileStorageDevice(storageId, &reconcilerCreateStorageDeviceRequest)

	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_INFRA_STORAGE_DEVICE_FAILED), createInput.SerialNumber, createInput.ManagementAddress)
		log.WriteDebug("TFError| error in  ReconcileStorageDevice , err: %v", err)
		return nil, err
	}

	terraformModelResponse := terraformmodel.InfraStorageDevices{}
	err = copier.Copy(&terraformModelResponse, sd)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_INFRA_STORAGE_DEVICE_END), createInput.SerialNumber, createInput.ManagementAddress)
	return &terraformModelResponse.Data, nil

}

func CreateInfraStorageDeviceRequestFromSchema(d *schema.ResourceData) (*terraformmodel.CreateInfraStorageDeviceParam, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateInfraStorageDeviceParam{}

	createInput.SerialNumber = common.GetSerialString(d)

	managementAddress, ok := d.GetOk("management_address")
	if ok {
		ma := managementAddress.(string)
		createInput.ManagementAddress = ma
	}

	username, ok := d.GetOk("username")
	if ok {
		un := username.(string)
		createInput.Username = un
	}

	password, ok := d.GetOk("password")
	if ok {
		pw := password.(string)
		createInput.Password = pw
	}

	gwAddress, ok := d.GetOk("gateway_address")
	if ok {
		gwa := gwAddress.(string)
		createInput.GatewayAddress = gwa
	} else {
		createInput.GatewayAddress = ""
	}

	outOfBand, ok := d.GetOk("out_of_band")
	if ok {
		oob := outOfBand.(bool)
		createInput.OutOfBand = oob
	}

	system, ok := d.GetOk("systems")
	if ok {
		createInput.UcpSystem = system.(string)
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func ConvertPartnersInfraStorageDeviceToSchema(pg *terraformmodel.InfraMTStorageDevice) *map[string]interface{} {
	sp := map[string]interface{}{
		"storage_serial_number": pg.Storage.SerialNumber,
		"resource_id":           pg.Storage.ResourceId,
		"ucp_systems":           pg.Storage.UcpSystems,
		"storage_id":            pg.StorageId,
		"status":                pg.Status,
	}

	if pg.PartnerId != "" {
		sp["partner_id"] = pg.PartnerId
	}

	if pg.SubscriberId != "" {
		sp["subscriber_id"] = pg.SubscriberId
	}
	return &sp
}

func ConvertInfraStorageDeviceToSchema(pg *terraformmodel.InfraStorageDeviceInfo) *map[string]interface{} {
	sp := map[string]interface{}{
		"resource_id":          pg.ResourceId,
		"serial_number":        pg.SerialNumber,
		"management_address":   pg.ManagementAddress,
		"controller_address":   pg.ControllerAddress,
		"username":             pg.Username,
		"systems":              pg.UcpSystems,
		"device_type":          pg.DeviceType,
		"model":                pg.Model,
		"microcode_version":    pg.MicrocodeVersion,
		"total_capacity_in_mb": pg.TotalCapacityInMb,
		"free_capacity_in_mb":  pg.FreeCapacityInMb,
		"total_capacity":       pg.TotalCapacity,
		"free_capacity":        pg.FreeCapacity,
		"resource_state":       pg.ResourceState,
		"tags":                 pg.Tags,
		"operational_status":   pg.OperationalStatus,

		"health_status":                          pg.HealthStatus,
		"free_gad_consistency_group_id":          pg.FreeGadConsistencyGroupId,
		"free_local_clone_consistency_group_id":  pg.FreeLocalCloneConsistencyGroupId,
		"free_remote_clone_consistency_group_id": pg.FreeRemoteCloneConsistencyGroupId,
	}

	ses := []map[string]interface{}{}

	p := map[string]interface{}{
		"compression_ratio": pg.StorageEfficiencyStat.CompressionRatio,
		"start_time":        pg.StorageEfficiencyStat.StartTime,
		"end_time":          pg.StorageEfficiencyStat.EndTime,
		"provisioning_rate": pg.StorageEfficiencyStat.ProvisioningRate,
		"snapshot_ratio":    pg.StorageEfficiencyStat.SnapshotRatio,
		"total_ratio":       pg.StorageEfficiencyStat.TotalRatio,
	}

	aces := []map[string]interface{}{}
	a := map[string]interface{}{
		"compression_ratio": pg.StorageEfficiencyStat.AccelCompEfficiencyStat.CompressionRatio,
		"dedupe_ratio":      pg.StorageEfficiencyStat.AccelCompEfficiencyStat.DedupeRatio,
		"reclaim_ratio":     pg.StorageEfficiencyStat.AccelCompEfficiencyStat.ReclaimRatio,
		"total_ratio":       pg.StorageEfficiencyStat.AccelCompEfficiencyStat.TotalRatio,
	}
	aces = append(aces, a)
	p["accel_comp_efficiency_stat"] = aces

	dces := []map[string]interface{}{}
	d := map[string]interface{}{
		"compression_ratio": pg.StorageEfficiencyStat.DedupeCompEfficiencyStat.CompressionRatio,
		"dedupe_ratio":      pg.StorageEfficiencyStat.DedupeCompEfficiencyStat.DedupeRatio,
		"reclaim_ratio":     pg.StorageEfficiencyStat.DedupeCompEfficiencyStat.ReclaimRatio,
		"total_ratio":       pg.StorageEfficiencyStat.DedupeCompEfficiencyStat.TotalRatio,
	}
	dces = append(dces, d)
	p["dedupe_comp_efficiency_stat"] = dces

	ses = append(ses, p)
	sp["storage_efficiency_stat"] = ses

	syslogConfig := []map[string]interface{}{}

	s := map[string]interface{}{
		"detailed": pg.SyslogConfig.Detailed,
	}
	syslogConfig = append(syslogConfig, s)

	syslogServers := []map[string]interface{}{}
	for _, sls := range pg.SyslogConfig.SyslogServers {
		p := map[string]interface{}{
			"id":                    sls.Id,
			"syslog_server_address": sls.SyslogServerAddress,
			"syslog_server_port":    sls.SyslogServerPort,
		}
		syslogServers = append(syslogServers, p)
	}

	ss := map[string]interface{}{
		"syslog_servers": syslogServers,
	}
	syslogConfig = append(syslogConfig, ss)
	sp["syslog_config"] = syslogConfig

	sdl := []map[string]interface{}{}
	for _, l := range pg.StorageDeviceLicenses {
		p := map[string]interface{}{
			"is_enabled":   l.IsEnabled,
			"is_installed": l.IsInstalled,
			"type":         l.Type,
			"name":         l.Name,
		}
		sdl = append(sdl, p)
	}
	sp["storage_device_licenses"] = sdl

	deviceLimits := []map[string]interface{}{}

	m := map[string]interface{}{
		"health_description": pg.DeviceLimits.HealthDescription,
		"is_unified":         pg.DeviceLimits.IsUnified,
	}
	deviceLimits = append(deviceLimits, m)

	egnr := []map[string]interface{}{}

	m1 := map[string]interface{}{
		"is_valid":  pg.DeviceLimits.ExternalGroupNumberRange.IsValid,
		"max_value": pg.DeviceLimits.ExternalGroupNumberRange.MaxValue,
		"min_value": pg.DeviceLimits.ExternalGroupNumberRange.MinValue,
	}
	egnr = append(egnr, m1)
	m["external_group_number_range"] = egnr

	egsnr := []map[string]interface{}{}
	m2 := map[string]interface{}{
		"is_valid":  pg.DeviceLimits.ExternalGroupSubNumberRange.IsValid,
		"max_value": pg.DeviceLimits.ExternalGroupSubNumberRange.MaxValue,
		"min_value": pg.DeviceLimits.ExternalGroupSubNumberRange.MinValue,
	}
	egsnr = append(egsnr, m2)
	m["external_group_number_range"] = egsnr

	pgnr := []map[string]interface{}{}
	m3 := map[string]interface{}{
		"is_valid":  pg.DeviceLimits.ParityGroupNumberRange.IsValid,
		"max_value": pg.DeviceLimits.ParityGroupNumberRange.MaxValue,
		"min_value": pg.DeviceLimits.ParityGroupNumberRange.MinValue,
	}
	pgnr = append(pgnr, m3)
	m["parity_group_number_range"] = pgnr

	pgsnr := []map[string]interface{}{}
	m4 := map[string]interface{}{
		"is_valid":  pg.DeviceLimits.ParityGroupSubNumberRange.IsValid,
		"max_value": pg.DeviceLimits.ParityGroupSubNumberRange.MaxValue,
		"min_value": pg.DeviceLimits.ParityGroupSubNumberRange.MinValue,
	}
	pgsnr = append(pgsnr, m4)
	m["parity_group_number_range"] = pgsnr

	deviceLimits = append(deviceLimits, m)
	sp["device_limits"] = deviceLimits

	return &sp
}
