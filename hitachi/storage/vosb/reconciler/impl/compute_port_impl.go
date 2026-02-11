package vssbstorage

import (
	"fmt"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	provisonerimpl "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/impl"
	provisonermodel "terraform-provider-hitachi/hitachi/storage/vosb/provisioner/model"
	mc "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/message-catalog"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"

	"github.com/jinzhu/copier"
)

// AllowChapUsersToAccessComputePort allows chap users to access compute port of vssb storage
func (psm *vssbStorageManager) AllowChapUsersToAccessComputePort(portId string, authMode string, inputChapUsers []string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_AUTH_SETTINGS_BEGIN), portId)
	existingPortAuthSettings, err := provObj.GetPortAuthSettings(portId)
	if err != nil {
		log.WriteDebug("TFError| error in GetPortAuthSettings provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_AUTH_SETTINGS_FAILED), portId)
		return err
	}

	if len(inputChapUsers) == 0 || strings.Compare(authMode, "None") == 0 {
		err = psm.DeleteAllChapUsersFromComputePort(portId)
		if err != nil {
			log.WriteDebug("TFError| error in DeleteAllChapUsersFromComputePort call, err: %v", err)
			return err
		}

		err = psm.UpdatePortAuthSettings(portId, existingPortAuthSettings, authMode)
		if err != nil {
			log.WriteDebug("TFError| error in UpdatePortAuthSettings  call, err: %v", err)
			return err
		}

		return nil
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_CHAPUSERS_BEGIN))
	provChapUsers, err := provObj.GetAllChapUsers()
	log.WriteDebug("TFError| Prov chapUsersInfo: %v", provChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in GetAllChapUsers provisioner call, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_CHAPUSERS_FAILED))
		return err
	}

	chapUserIds, err := CheckAndGetIdsForInputChapUsers(provChapUsers, inputChapUsers)
	if err != nil {
		log.WriteDebug("TFError| All input target chap users  are not present, err: %v", err)
		return err
	}

	err = psm.DeleteAllChapUsersFromComputePort(portId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteAllChapUsersFromComputePort call, err: %v", err)
		return err
	}
	log.WriteDebug("TFError| DeleteAllChapUsersFromComputePort  call successful")

	err = psm.UpdatePortAuthSettings(portId, existingPortAuthSettings, authMode)
	if err != nil {
		log.WriteDebug("TFError| error in UpdatePortAuthSettings  call, err: %v", err)
		return err
	}

	log.WriteDebug("TFError| UpdatePortAuthSettings  call successful")

	err = psm.AddChapUsersToComputePort(portId, chapUserIds)
	if err != nil {
		log.WriteDebug("TFError| error in UpdatePortAuthSettings  call, err: %v", err)
		return err
	}

	log.WriteDebug("TFError| AddChapUsersToComputePort  call successful")

	return nil
}

// CheckAndGetIdsForInputChapUsers checks if all input chap users are prsent in the system, if success return their corresponding Ids
func CheckAndGetIdsForInputChapUsers(sysChapUsers *provisonermodel.ChapUsers, inputChapUsers []string) ([]string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	m := make(map[string]provisonermodel.ChapUser)
	for _, cu := range sysChapUsers.Data {
		m[cu.TargetChapUserName] = cu
	}

	var sb strings.Builder
	sb.WriteString("")
	allPresent := true

	chapUsersIdList := []string{}

	for _, icu := range inputChapUsers {
		if v, found := m[icu]; found {
			chapUsersIdList = append(chapUsersIdList, v.ID)
		} else {
			sb.WriteString(icu + " ")
			allPresent = false
		}
	}

	log.WriteDebug("TFError| chap user Id list %v", chapUsersIdList)

	if allPresent {
		return chapUsersIdList, nil
	} else {
		sb.WriteString("chap users are not present in the system")
		return nil, fmt.Errorf(sb.String())
	}

}

// DeleteAllChapUsersFromComputePort removes chap users from a specific compute port of vssb storage
func (psm *vssbStorageManager) DeleteAllChapUsersFromComputePort(portId string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = provObj.DeleteAllChapUsersFromComputePort(portId)
	if err != nil {
		log.WriteDebug("TFError| error in DeleteAllChapUsersFromComputePort provisioner call, err: %v", err)
		return err
	}

	return nil
}

func (psm *vssbStorageManager) UpdatePortAuthSettings(portId string, existingPortAuthSettings *provisonermodel.PortAuthSettings, authMode string) error {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	var pas = provisonermodel.PortAuthSettings{}
	pas.AuthMode = authMode
	if strings.Compare(authMode, "None") == 0 {
		pas.IsDiscoveryChapAuth = false
		pas.IsMutualChapAuth = false
	} else if strings.Compare(authMode, "CHAP") == 0 {
		pas.IsDiscoveryChapAuth = existingPortAuthSettings.IsDiscoveryChapAuth
		pas.IsMutualChapAuth = true
	} else {
		pas.IsDiscoveryChapAuth = existingPortAuthSettings.IsDiscoveryChapAuth
		pas.IsMutualChapAuth = existingPortAuthSettings.IsMutualChapAuth
	}

	err = provObj.UpdatePortAuthSettings(portId, &pas)
	if err != nil {
		log.WriteDebug("TFError| error in UpdatePortAuthSettings provisioner call, err: %v", err)
		return err
	}

	log.WriteInfo("UpdatePortAuthSettings Successful")
	return nil
}

// AddChapUsersToComputePort allows chap users to access a specific compute port of vssb storage
func (psm *vssbStorageManager) AddChapUsersToComputePort(portId string, chapUserIds []string) error {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return err
	}

	err = provObj.AddChapUsersToComputePort(portId, chapUserIds)
	if err != nil {
		log.WriteDebug("TFError| error in AddChapUsersToComputePort provisioner call, err: %v", err)
		return err
	}

	return nil
}

func (psm *vssbStorageManager) GetPortInfoByID(portId string) (*vssbmodel.PortDetailSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portId)

	reconPort, err := provObj.GetPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPort, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, err
	}

	reconPortAuthSetting, err := provObj.GetPortAuthSettings(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPort, err: %v", err)
		return nil, err
	}

	reconChapUsers, err := provObj.GetChapUsersAllowedToAccessPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetChapUsersAllowedToAccessPort, err: %v", err)
		return nil, err
	}

	// Converting reconciler to terraform
	reconPortDetails := vssbmodel.PortDetailSettings{}
	err = copier.Copy(&reconPortDetails.Port, reconPort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	err = copier.Copy(&reconPortDetails.AuthSettings, reconPortAuthSetting)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	err = copier.Copy(&reconPortDetails.ChapUsers, reconChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	return &reconPortDetails, nil

}

func (psm *vssbStorageManager) GetIscsiPortAuthInfo(portId string) (*vssbmodel.PortDetailSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := provisonermodel.StorageDeviceSettings{
		Username:       psm.storageSetting.Username,
		Password:       psm.storageSetting.Password,
		ClusterAddress: psm.storageSetting.ClusterAddress,
	}

	provObj, err := provisonerimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_PORT_BEGIN), portId)

	reconPort, err := provObj.GetPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPort, err: %v", err)
		log.WriteError(mc.GetMessage(mc.ERR_GET_PORT_FAILED), portId)
		return nil, err
	}

	reconPortAuthSetting, err := provObj.GetPortAuthSettings(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetPort, err: %v", err)
		return nil, err
	}

	reconChapUsers, err := provObj.GetChapUsersAllowedToAccessPort(portId)
	if err != nil {
		log.WriteDebug("TFError| error getting GetChapUsersAllowedToAccessPort, err: %v", err)
		return nil, err
	}

	// Converting reconciler to terraform
	reconPortDetails := vssbmodel.PortDetailSettings{}
	err = copier.Copy(&reconPortDetails.Port, reconPort)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	err = copier.Copy(&reconPortDetails.AuthSettings, reconPortAuthSetting)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	err = copier.Copy(&reconPortDetails.ChapUsers, reconChapUsers)
	if err != nil {
		log.WriteDebug("TFError| error in Copy from reconciler to terraform structure, err: %v", err)
		return nil, err
	}

	return &reconPortDetails, nil

}
