package terraform

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"

	"fmt"
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

var chapUserTypeConversion = map[string]string{
	"initiator": "INI",
	"target":    "TAR",
}

var chapUserTypeToSchema = map[string]string{
	"INI": "initiator",
	"TAR": "target",
}

func GetIscsiTargetChapUser(d *schema.ResourceData) (*terraformmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)
	portID := d.Get("port_id").(string)
	itNum := d.Get("iscsi_target_number").(int)
	cuType := d.Get("chap_user_type").(string)
	chapUserName := d.Get("chap_user_name").(string)

	chapUserType := chapUserTypeConversion[strings.ToLower(cuType)]
	if chapUserType == "" {
		err := fmt.Errorf("invalid chap user type specified %v", chapUserType)
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
		log.WriteDebug("TFError| error in NewEx, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_BEGIN), portID, itNum, chapUserName, chapUserType)
	itcu, err := reconObj.GetChapUser(portID, itNum, chapUserName, chapUserType)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSER_FAILED), portID, itNum, chapUserName, chapUserType)
		return nil, err
	}

	terraformModelIscsiTargetChapUser := terraformmodel.IscsiTargetChapUser{}
	err = copier.Copy(&terraformModelIscsiTargetChapUser, itcu)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSER_END), portID, itNum, chapUserName, chapUserType)

	return &terraformModelIscsiTargetChapUser, nil
}

func GetIscsiTargetChapUsers(d *schema.ResourceData) (*terraformmodel.IscsiTargetChapUsers, error) {
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

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_BEGIN), portID, itNum)

	iscsiChapUsers, err := reconObj.GetChapUsers(portID, itNum)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ISCSITARGET_CHAPUSERS_FAILED), portID, itNum)
		return nil, err
	}

	terraformModelIscsiTargetChapUsers := terraformmodel.IscsiTargetChapUsers{}
	err = copier.Copy(&terraformModelIscsiTargetChapUsers, iscsiChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ISCSITARGET_CHAPUSERS_END), portID, itNum)

	return &terraformModelIscsiTargetChapUsers, nil
}

func CreateIscsiTargetChapUser(d *schema.ResourceData) (*terraformmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Resource Data: %v", d)

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

	createInput, err := CreateIscsiTargetChapUserRequestFromSchema(d)
	if err != nil {
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_CHAPUSER_BEGIN), createInput.PortID, createInput.IscsiTargetNumber, createInput.ChapUserName, createInput.WayOfChapUser)
	reconcilerCreateIscsiChapUserRequest := reconcilermodel.ChapUserRequest{}
	err = copier.Copy(&reconcilerCreateIscsiChapUserRequest, createInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	itcu, err := reconObj.ReconcileChapUser(&reconcilerCreateIscsiChapUserRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_ISCSITARGET_CHAPUSER_FAILED), createInput.PortID, createInput.IscsiTargetNumber, createInput.ChapUserName, createInput.WayOfChapUser)
		log.WriteDebug("TFError| error in Creating IscsiTarget Chap User - ReconcileIscsiTargetChapUser, err: %v", err)
		return nil, err
	}

	terraformModelIscsiTargetChapUser := terraformmodel.IscsiTargetChapUser{}
	err = copier.Copy(&terraformModelIscsiTargetChapUser, itcu)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_ISCSITARGET_CHAPUSER_END), terraformModelIscsiTargetChapUser.PortID, terraformModelIscsiTargetChapUser.HostGroupNumber, terraformModelIscsiTargetChapUser.ChapUserName, terraformModelIscsiTargetChapUser.WayOfChapUser)
	return &terraformModelIscsiTargetChapUser, nil
}

func ConvertIscsiTargetChapUserToSchema(iscsiTargetChapUser *terraformmodel.IscsiTargetChapUser, serial int) *map[string]interface{} {

	itcu := map[string]interface{}{
		"storage_serial_number": serial,
		"iscsi_target_number":   iscsiTargetChapUser.HostGroupNumber,
		"port_id":               iscsiTargetChapUser.PortID,
		"chap_user_type":        chapUserTypeToSchema[iscsiTargetChapUser.WayOfChapUser],
		"chap_user_name":        iscsiTargetChapUser.ChapUserName,
		"chap_user_id":          iscsiTargetChapUser.ChapUserID,
	}

	return &itcu
}

// UpdateIscsiTargetChapUser used to update modified data of Chap User
func UpdateIscsiTargetChapUser(d *schema.ResourceData) (*terraformmodel.IscsiTargetChapUser, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteDebug("Input Res: %+v", d)
	log.WriteDebug("Input Info: %+v", d.Get("chap_user"))
	log.WriteDebug("Input State: %+v", d.Get("state"))
	log.WriteDebug("Input Diff: %+v", d.Get("diff"))
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

	updateInput, err := CreateIscsiTargetChapUserRequestFromSchema(d)
	if err != nil {
		log.WriteDebug("TFError| error in CreateIscsiTargetChapUserRequestFromSchema, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSITARGET_CHAPUSER_BEGIN), updateInput.PortID, updateInput.IscsiTargetNumber, updateInput.ChapUserName, updateInput.WayOfChapUser)
	reconcilerUpdateIscsiChapUserRequest := reconcilermodel.ChapUserRequest{}
	err = copier.Copy(&reconcilerUpdateIscsiChapUserRequest, &updateInput)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	itcu, err := reconObj.ReconcileChapUser(&reconcilerUpdateIscsiChapUserRequest)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CHANGE_ISCSITARGET_CHAPUSER_FAILED), updateInput.PortID, updateInput.IscsiTargetNumber, updateInput.ChapUserName, updateInput.WayOfChapUser)
		log.WriteDebug("TFError| error in Updating Chap User - ReconcileChap User , err: %v", err)
		return nil, err
	}

	terraformModelIscsiTargetChapUser := terraformmodel.IscsiTargetChapUser{}
	err = copier.Copy(&terraformModelIscsiTargetChapUser, itcu)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CHANGE_ISCSITARGET_CHAPUSER_END), updateInput.PortID, updateInput.IscsiTargetNumber, updateInput.ChapUserName, updateInput.WayOfChapUser)
	return &terraformModelIscsiTargetChapUser, nil
}

func CreateIscsiTargetChapUserRequestFromSchema(d *schema.ResourceData) (*terraformmodel.ChapUserRequest, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	createInput := terraformmodel.ChapUserRequest{}

	cuname, ok := d.GetOk("chap_user_name")
	if ok {
		name := cuname.(string)
		createInput.ChapUserName = name
	}

	portId, ok := d.GetOk("port_id")
	if ok {
		pid := portId.(string)
		createInput.PortID = pid
	}

	hgnum, ok := d.GetOk("iscsi_target_number")
	if ok {
		hid := hgnum.(int)
		createInput.IscsiTargetNumber = hid
	}

	wayofcu, ok := d.GetOk("chap_user_type")
	if ok {
		woc := wayofcu.(string)
		createInput.WayOfChapUser = chapUserTypeConversion[strings.ToLower(woc)]
		if createInput.WayOfChapUser == "" {
			err := fmt.Errorf("invalid chap user type specified %v", woc)
			return nil, err

		}

	}

	cupass, ok := d.GetOk("chap_user_password")
	if ok {
		pass := cupass.(string)
		createInput.ChapUserSecret = pass
	}

	log.WriteDebug("createInput: %+v", createInput)
	return &createInput, nil
}

func DeleteIscsiTargetChapUser(d *schema.ResourceData) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	serial := d.Get("serial").(int)

	portId, ok := d.GetOk("port_id")
	log.WriteDebug("portId: %+v", portId)

	if !ok {
		chap_user, ok := d.GetOk("chap_user")
		if !ok {
			return fmt.Errorf("no chap_user data in resource")
		}
		log.WriteDebug("chap_user: %+v", chap_user.([]map[string]interface{})[0])
		portId, ok = chap_user.([]map[string]interface{})[0]["portId"]
		if !ok {
			return fmt.Errorf("found no portId in chap_user")
		}
		log.WriteDebug("chap_user portId: %+v", portId)
	}
	portID := portId.(string)

	itNum, ok := d.GetOk("iscsi_target_number")
	if !ok {
		chap_user, ok := d.GetOk("chap_user")
		if !ok {
			return fmt.Errorf("no chap_user data in resource")
		}
		log.WriteDebug("chap_user: %+v", chap_user.([]map[string]interface{})[0])
		itNum, ok = chap_user.([]map[string]interface{})[0]["iscsi_target_number"]
		if !ok {
			return fmt.Errorf("found no itNum in chap_user")
		}
		log.WriteDebug("chap_user itNum: %+v", itNum)
	}
	itNumber := itNum.(int)

	chapUserName, ok := d.GetOk("chap_user_name")
	log.WriteDebug("chapUserName: %+v", chapUserName)

	if !ok {
		chap_user, ok := d.GetOk("chap_user")
		if !ok {
			return fmt.Errorf("no chap_user data in resource")
		}
		log.WriteDebug("chap_user: %+v", chap_user.([]map[string]interface{})[0])
		chapUserName, ok = chap_user.([]map[string]interface{})[0]["chap_user_name"]
		if !ok {
			return fmt.Errorf("found no chapUserName in chap_user")
		}
		log.WriteDebug("chap_user chapUserName: %+v", chapUserName)
	}
	cuname := chapUserName.(string)

	wayOfChapUser, ok := d.GetOk("chap_user_type")
	log.WriteDebug("chapUserName: %+v", wayOfChapUser)

	if !ok {
		chap_user, ok := d.GetOk("chap_user")
		if !ok {
			return fmt.Errorf("no chap_user data in resource")
		}
		log.WriteDebug("chap_user: %+v", chap_user.([]map[string]interface{})[0])
		wayOfChapUser, ok = chap_user.([]map[string]interface{})[0]["chap_user_type"]
		if !ok {
			return fmt.Errorf("found no wayOfChapUser in chap_user")
		}
		log.WriteDebug("chap_user chapUserName: %+v", wayOfChapUser)
	}
	woc_file := wayOfChapUser.(string)
	woc := chapUserTypeConversion[strings.ToLower(woc_file)]

	if woc == "" {
		err := fmt.Errorf("invalid chap user type specified %v", woc_file)
		return err

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

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_BEGIN), portID, itNumber, cuname, woc)
	err = reconObj.DeleteChapUser(portID, itNumber, cuname, woc)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_ISCSITARGET_CHAPUSER_FAILED), portID, itNumber, cuname, woc)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_ISCSITARGET_CHAPUSER_END), portID, itNumber, cuname, woc)

	return nil
}
