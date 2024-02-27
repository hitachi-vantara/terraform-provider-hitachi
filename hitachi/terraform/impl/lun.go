package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	"fmt"
	"strings"

	// "io/ioutil"
	"strconv"
	// "time"

	cache "terraform-provider-hitachi/hitachi/common/cache"
	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// "terraform-provider-hitachi/hitachi/common/utils"

	reconimpl "terraform-provider-hitachi/hitachi/storage/san/reconciler/impl"
	reconcilermodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	mc "terraform-provider-hitachi/hitachi/terraform/message-catalog"
	terraformmodel "terraform-provider-hitachi/hitachi/terraform/model"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jinzhu/copier"
)

// GetLun gets a lun
func GetLun(d *schema.ResourceData) (*terraformmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	// check if this is getting executed from "data"
	lunID := d.Get("ldev_id").(int)

	// check if this is getting executed from "resource"
	_, lunOk := d.GetOk("ldev_id")
	if !lunOk {
		lunFromState := d.State().ID
		log.WriteDebug("TFDebug| lunFromState from state: %s", lunFromState)
		if lunFromState != "" {
			lun, err := strconv.Atoi(lunFromState)
			if err != nil {
				log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
				return nil, err
			}
			lunID = lun
		}
	}

	log.WriteDebug("TFDebug| lunID: %v", lunID)

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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_BEGIN), lunID)
	lun, err := reconObj.GetLun(lunID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_LUN_FAILED), lunID)
		return nil, err
	}

	terraformModelLun := terraformmodel.LogicalUnit{}
	err = copier.Copy(&terraformModelLun, lun)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_END), lunID)

	return &terraformModelLun, nil
}

// GetRangeOfLuns gets the desired luns based on range specified
func GetRangeOfLuns(d *schema.ResourceData) (*[]terraformmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	startLdevID := d.Get("start_ldev_id").(int)
	if startLdevID < 0 {
		return nil, fmt.Errorf("start_ldev_id must be greater than or equal to 0")
	}

	endLdevID := d.Get("end_ldev_id").(int)
	if endLdevID < 0 {
		return nil, fmt.Errorf("end_ldev_id must be greater than or equal to 0")
	}

	if endLdevID < startLdevID {
		return nil, fmt.Errorf("end_ldev_id must be greater than or equal to start_ldev_id")
	}

	isUndefindLdev := d.Get("undefined_ldev").(bool)

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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_BEGIN), startLdevID, endLdevID)
	luns, err := reconObj.GetRangeOfLuns(startLdevID, endLdevID, isUndefindLdev)
	if err != nil {
		log.WriteDebug("TFError| error in GetRangeOfLuns terraform call, err: %v", err)
		return nil, err
	}

	terraformModelLuns := []terraformmodel.LogicalUnit{}
	err = copier.Copy(&terraformModelLuns, luns)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_LUN_RANGE_END), startLdevID, endLdevID)

	return &terraformModelLuns, nil
}

// CreateLun creates a lun
func CreateLun(d *schema.ResourceData) (*terraformmodel.LogicalUnit, error) {
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

	pldevID := CheckSchemaIfLunGet(d)
	if pldevID != nil {
		lun, err := reconObj.GetLun(*pldevID)
		if err != nil {
			return nil, err
		}
		if lun.ByteFormatCapacity == "" {
			// does not exist, or in the process of being deleted
			return nil, fmt.Errorf("volume does not exist")
		}
		terraformModelLun := terraformmodel.LogicalUnit{}
		err = copier.Copy(&terraformModelLun, lun)
		if err != nil {
			log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
			return nil, err
		}

		return &terraformModelLun, nil
	}

	createInput, err := CreateLunRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	reconcilerCreateLunRequest := reconcilermodel.LunRequest{}
	err = copier.Copy(&reconcilerCreateLunRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	lun, err := reconObj.SetLun(&reconcilerCreateLunRequest)
	if err != nil {
		log.WriteDebug("TFError| error in SetLun, err: %v", err)
		return nil, err
	}

	terraformModelLun := terraformmodel.LogicalUnit{}
	err = copier.Copy(&terraformModelLun, lun)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	return &terraformModelLun, nil
}

func CreateLunRequestFromSchema(d *schema.ResourceData) (*terraformmodel.CreateLunRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.CreateLunRequest{}

	size_gb, _ := d.GetOk("size_gb")

	createInput.CapacityInGB = uint64(size_gb.(int))

	ldevId, _ := d.GetOk("ldev_id")
	if ldevId.(int) >= 0 {
		lid := ldevId.(int)
		createInput.LdevID = &lid
	}

	name, ok := d.GetOk("name")
	if ok {
		label := name.(string)
		createInput.Name = &label
	}

	var pool_name = ""
	var paritygroup_id = ""
	pool_id, _ := d.GetOk("pool_id")

	pool_name = d.Get("pool_name").(string)
	paritygroup_id = d.Get("paritygroup_id").(string)
	log.WriteDebug("Pool ID=%v Pool Name=%v PG=%v\n", pool_id, pool_name, paritygroup_id)

	log.WriteDebug("ok=%v \n", ok)

	if pool_id.(int) >= 0 {
		pool_id_int := pool_id.(int)
		createInput.PoolID = &pool_id_int
	} else if pool_name != "" {
		ppid, err := GetPoolIdFromPoolName(d, pool_name)
		createInput.PoolID = ppid
		if err != nil {
			return nil, fmt.Errorf("could not find a pool with name %v", pool_name)
		}
	} else if paritygroup_id != "" {
		createInput.ParityGroupID = &paritygroup_id
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func GetPoolIdFromPoolName(d *schema.ResourceData, poolName string) (*int, error) {
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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	//log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevID, setting.Serial)

	pools, err := reconObj.GetPools()
	if err != nil {
		//log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevID, setting.Serial)
		return nil, err
	}
	//log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevID, setting.Serial)

	poolId := -1

	for _, pool := range *pools {
		if strings.EqualFold(pool.PoolName, poolName) {
			poolId = pool.PoolID
			break

		}
	}

	return &poolId, nil
}

func CheckSchemaIfLunGet(d *schema.ResourceData) *int {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	fields := []string{
		"size_gb",
		"name",
		//"dedup_mode",
		"pool_id",
		"pool_name",
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

// DeleteLun deletes a lun
func DeleteLun(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	ldevID := d.Get("ldev_id")
	log.WriteDebug("ldevID: %+v", ldevID)
	lunID := 0
	isLdevIdSetFromState := false
	if ldevID.(int) <= 0 {
		lunFromState := d.State().ID
		if lunFromState != "" {
			lun, err := strconv.Atoi(lunFromState)
			if err != nil {
				log.WriteDebug("TFError| error while converting string to int lunID, err: %v", err)
				return err
			}
			lunID = lun
			isLdevIdSetFromState = true
		} else {
			volume, ok := d.GetOk("volume")
			if !ok {
				return fmt.Errorf("no volume data in resource")
			}
			log.WriteDebug("volume: %+v", volume.([]map[string]interface{})[0])
			ldevID, ok = volume.([]map[string]interface{})[0]["ldev_id"]
			if !ok {
				return fmt.Errorf("found no ldev_id in info")
			}
			log.WriteDebug("volume ldevID: %+v", ldevID)
		}
	}

	if !isLdevIdSetFromState {
		lunID = ldevID.(int)
	}

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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_BEGIN), ldevID, setting.Serial)

	err = reconObj.DeleteLun(lunID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LUN_FAILED), ldevID, setting.Serial)
		return err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LUN_END), ldevID, setting.Serial)

	return nil
}

func ConvertLunToSchema(logicalUnit *terraformmodel.LogicalUnit, serial int) *map[string]interface{} {
	lun := map[string]interface{}{
		"storage_serial_number":      serial,
		"ldev_id":                    logicalUnit.LdevID,
		"clpr_id":                    logicalUnit.ClprID,
		"emulation_type":             logicalUnit.EmulationType,
		"num_ports":                  logicalUnit.NumOfPorts,
		"attributes":                 logicalUnit.Attributes,
		"label":                      logicalUnit.Label,
		"status":                     logicalUnit.Status,
		"mpblade_id":                 logicalUnit.MpBladeID,
		"ss_id":                      logicalUnit.Ssid,
		"pool_id":                    logicalUnit.PoolID,
		"parity_group_id":            logicalUnit.ParityGroupId,
		"is_full_allocation_enabled": logicalUnit.IsFullAllocationEnabled,
		"resource_group_id":          logicalUnit.ResourceGroupID,
		//"data_reduction_mode":        logicalUnit.DataReductionMode,
		"is_alua_enabled":      logicalUnit.IsAluaEnabled,
		"naa_id":               logicalUnit.NaaID,
		"total_capacity_in_mb": logicalUnit.TotalCapacityInMB,
		"free_capacity_in_mb":  logicalUnit.FreeCapacityInMB,
		"used_capacity_in_mb":  logicalUnit.UsedCapacityInMB,
	}

	ports := []map[string]interface{}{}
	for _, pin := range logicalUnit.Ports {
		p := map[string]interface{}{
			"port_id":        pin.PortID,
			"hostgroup_id":   pin.HostGroupNumber,
			"hostgroup_name": pin.HostGroupName,
			"lun_id":         pin.Lun,
		}
		ports = append(ports, p)
	}
	lun["ports"] = ports

	return &lun
}

// UpdateLun updates a lun
func UpdateLun(d *schema.ResourceData) (*terraformmodel.LogicalUnit, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	updateInput, err := UpdateLunRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

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
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	reconcilerUpdateLunRequest := reconcilermodel.UpdateLunRequest{}
	err = copier.Copy(&reconcilerUpdateLunRequest, updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_BEGIN), reconcilerUpdateLunRequest.LdevID, setting.Serial)
	lun, err := reconObj.UpdateLun(&reconcilerUpdateLunRequest)
	if err != nil {
		log.WriteDebug("TFError| error in UpdateLun reconciler call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_LUN_FAILED), reconcilerUpdateLunRequest.LdevID, setting.Serial)
		return nil, err
	}

	terraformModelLun := terraformmodel.LogicalUnit{}
	err = copier.Copy(&terraformModelLun, lun)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_LUN_END), reconcilerUpdateLunRequest.LdevID, setting.Serial)

	return &terraformModelLun, nil
}

func UpdateLunRequestFromSchema(d *schema.ResourceData) (*terraformmodel.UpdateLunRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Input Res: %+v", d)
	log.WriteDebug("Input volume: %+v", d.Get("volume"))
	log.WriteDebug("Input State: %+v", d.Get("state"))
	log.WriteDebug("Input Diff: %+v", d.Get("diff"))

	updateInput := terraformmodel.UpdateLunRequest{}

	if d.HasChange("size_gb") {
		old, new := d.GetChange("size_gb")
		expandSize := new.(int) - old.(int)
		size_gb := uint64(expandSize)
		updateInput.CapacityInGB = &size_gb
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateInput.Name = &name
	}

	// Remove dedup from this version
	/*
		if d.HasChange("dedup_mode") {
			dedup := d.Get("dedup_mode").(string)
			if dedup == "" {
				dedup = "disabled"
			}
			updateInput.DataReductionMode = &dedup
		}
	*/

	pldevID, err := getLdevIdFromSchema(d)
	if err != nil {
		return nil, err
	}
	updateInput.LdevID = pldevID

	log.WriteDebug("updateInput: %+v", updateInput)
	return &updateInput, nil
}

func getLdevIdFromSchema(d *schema.ResourceData) (*int, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	ldevID := d.Get("ldev_id")
	ldevIDInt := ldevID.(int)
	log.WriteDebug("spec input ldevID: %+v", ldevID)
	if ldevIDInt <= 0 {
		volume, ok := d.GetOk("volume")
		if !ok {
			return nil, fmt.Errorf("no info data in resource")
		}
		log.WriteDebug("volume: %+v", volume.([]interface{})[0])
		ldevID, ok = volume.([]interface{})[0].(map[string]interface{})["ldev_id"]
		if !ok {
			return nil, fmt.Errorf("found no ldev_id in info")
		}
		log.WriteDebug("volume ldevID: %+v", ldevID)
	}
	ldev := ldevID.(int)
	return &ldev, nil
}
