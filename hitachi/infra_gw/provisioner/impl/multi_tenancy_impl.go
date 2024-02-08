package infra_gw

import (
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/infra_gw/gateway/impl"
	model "terraform-provider-hitachi/hitachi/infra_gw/model"
)

// GetUserAdminRoleStatus gets all UCP Systems information
func (psm *infraGwManager) GetUserAdminRoleStatus(username string) (*bool, *string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	status := false
	user, err := psm.GetUserDetailsByName(username)

	if err != nil {
		log.WriteDebug("TFError| error in GetAllUsers gateway call, err: %v", err)
		return nil, nil, err
	}

	for _, role := range user.Roles {
		if strings.Contains(role.Name, model.AdminRole) || strings.Contains(role.Name, model.StorageAdminRole) {
			log.WriteInfo("Found user role with %s", role.Name)
			status = true
			return &status, &user.Id, nil
		}
	}

	log.WriteInfo("User doest not have any admin role")

	return &status, nil, nil

}

func (psm *infraGwManager) GetUserDetailsByName(username string) (*model.User, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return nil, nil
	}

	users, err := gatewayObj.GetAllUsers()

	if err != nil {
		log.WriteDebug("TFError| error in GetAllUsers gateway call, err: %v", err)
		return nil, err
	}
	for _, user := range users.Data.Users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, nil

}

func (psm *infraGwManager) GetPartnerIdWithStatus(username string) (bool, *string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	objStorage := model.InfraGwSettings{
		Username: psm.setting.Username,
		Password: psm.setting.Password,
		Address:  psm.setting.Address,
	}

	gatewayObj, err := gatewayimpl.NewEx(objStorage)
	if err != nil {
		log.WriteDebug("TFError| error in NewEx call, err: %v", err)
		return false, nil, err
	}

	adminStatus, _, err := psm.GetUserAdminRoleStatus(username)
	if err != nil {
		log.WriteDebug("TFError| error in GetUserAdminRoleStatus call, err: %v", err)
		return *adminStatus, nil, err
	}

	if !*adminStatus {
		return *adminStatus, nil, nil
	}

	partners, err := gatewayObj.GetAllPartners()
	if err != nil {
		log.WriteDebug("TFError| error in GetAllUsers gateway call, err: %v", err)
		return false, nil, err
	}

	if partners != nil || len(*partners) > 0 {
		return *adminStatus, &(*partners)[0].PartnerID, nil
	}
	return false, nil, nil
}
