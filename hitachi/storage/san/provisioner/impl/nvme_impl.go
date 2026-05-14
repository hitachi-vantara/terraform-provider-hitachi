package sanstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	gatewayimpl "terraform-provider-hitachi/hitachi/storage/san/gateway/impl"
	model "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	mc "terraform-provider-hitachi/hitachi/storage/san/provisioner/message-catalog"
)

/* --- Helper Method --- */

func (psm *sanStorageManager) getDeviceSettings() model.StorageDeviceSettings {
	return model.StorageDeviceSettings{
		Serial:   psm.storageSetting.Serial,
		Username: psm.storageSetting.Username,
		Password: psm.storageSetting.Password,
		MgmtIP:   psm.storageSetting.MgmtIP,
	}
}

func (psm *sanStorageManager) GetAllNvmSubsystems(params model.GetNvmSubsystemsParams) (*model.NvmSubsystems, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEMS_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	subsystems, err := gatewayObj.GetAllNvmSubsystems(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NVM_SUBSYSTEMS_FAILED), psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEMS_END), psm.storageSetting.Serial)
	return subsystems, nil
}

func (psm *sanStorageManager) GetNvmSubsystem(subsystemID int) (*model.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	subsystem, err := gatewayObj.GetNvmSubsystem(subsystemID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NVM_SUBSYSTEM_FAILED), subsystemID, psm.storageSetting.Serial)
		return nil, err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_END), subsystemID, psm.storageSetting.Serial)
	return subsystem, nil
}

func (psm *sanStorageManager) CreateNvmSubsystem(request model.CreateNvmSubsystemRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_NVM_SUBSYSTEM_BEGIN), psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.CreateNvmSubsystem(request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_NVM_SUBSYSTEM_FAILED), psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_NVM_SUBSYSTEM_END), psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) UpdateNvmSubsystem(subsystemID int, request model.UpdateNvmSubsystemRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_NVM_SUBSYSTEM_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.UpdateNvmSubsystem(subsystemID, request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_UPDATE_NVM_SUBSYSTEM_FAILED), subsystemID, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_UPDATE_NVM_SUBSYSTEM_END), subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteNvmSubsystem(subsystemID int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NVM_SUBSYSTEM_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteNvmSubsystem(subsystemID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_NVM_SUBSYSTEM_FAILED), subsystemID, psm.storageSetting.Serial)
		return "", err
	}

	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NVM_SUBSYSTEM_END), subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

/* --- NVM Subsystem Port Management --- */

func (psm *sanStorageManager) GetNvmSubsystemPort(subsystemID int, portID string) (*model.NvmSubsystemPort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_PORT_BEGIN), portID, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	port, err := gatewayObj.GetNvmSubsystemPort(subsystemID, portID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NVM_SUBSYSTEM_PORT_FAILED), portID, subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_PORT_END), portID, subsystemID, psm.storageSetting.Serial)
	return port, nil
}

func (psm *sanStorageManager) GetNvmSubsystemPorts(subsystemID int) (*model.NvmSubsystemPorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_PORTS_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	ports, err := gatewayObj.GetNvmSubsystemPorts(subsystemID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NVM_SUBSYSTEM_PORTS_FAILED), subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NVM_SUBSYSTEM_PORTS_END), subsystemID, psm.storageSetting.Serial)
	return ports, nil
}

func (psm *sanStorageManager) AddNvmSubsystemPort(request model.AddNvmSubsystemPortRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NVM_SUBSYSTEM_PORT_BEGIN), request.PortId, request.NvmSubsystemId, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.AddNvmSubsystemPort(request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_ADD_NVM_SUBSYSTEM_PORT_FAILED), request.PortId, request.NvmSubsystemId, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_ADD_NVM_SUBSYSTEM_PORT_END), request.PortId, request.NvmSubsystemId, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteNvmSubsystemPort(subsystemID int, portID string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NVM_SUBSYSTEM_PORT_BEGIN), portID, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteNvmSubsystemPort(subsystemID, portID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_NVM_SUBSYSTEM_PORT_FAILED), portID, subsystemID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NVM_SUBSYSTEM_PORT_END), portID, subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

/* --- Host NQN Management --- */

func (psm *sanStorageManager) GetAllHostNqns(subsystemID int) (*model.HostNqns, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOST_NQNS_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	nqns, err := gatewayObj.GetAllHostNqns(subsystemID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_HOST_NQNS_FAILED), subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_HOST_NQNS_END), subsystemID, psm.storageSetting.Serial)
	return nqns, nil
}

func (psm *sanStorageManager) GetHostNqn(subsystemID int, hostNqn string) (*model.HostNqnInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOST_NQN_BEGIN), hostNqn, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	nqnInfo, err := gatewayObj.GetHostNqn(subsystemID, hostNqn)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_HOST_NQN_FAILED), hostNqn, subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_HOST_NQN_END), hostNqn, subsystemID, psm.storageSetting.Serial)
	return nqnInfo, nil
}

func (psm *sanStorageManager) RegisterHostNqn(request model.RegisterHostNqnRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_REGISTER_HOST_NQN_BEGIN), request.HostNqn, request.NvmSubsystemId, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.RegisterHostNqn(request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_REGISTER_HOST_NQN_FAILED), request.HostNqn, request.NvmSubsystemId, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REGISTER_HOST_NQN_END), request.HostNqn, request.NvmSubsystemId, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) SetHostNqnNickname(subsystemID int, hostNqn string, request model.SetHostNqnNicknameRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_HOST_NQN_NICKNAME_BEGIN), hostNqn, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.SetHostNqnNickname(subsystemID, hostNqn, request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_SET_HOST_NQN_NICKNAME_FAILED), hostNqn, subsystemID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_HOST_NQN_NICKNAME_END), hostNqn, subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteLoginHostNqn(portID string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LOGIN_HOST_NQN_BEGIN), portID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteLoginHostNqn(portID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_LOGIN_HOST_NQN_FAILED), portID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_LOGIN_HOST_NQN_END), portID, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteHostNqn(subsystemID int, hostNqn string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOST_NQN_BEGIN), hostNqn, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteHostNqn(subsystemID, hostNqn)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_HOST_NQN_FAILED), hostNqn, subsystemID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_HOST_NQN_END), hostNqn, subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

/* --- Namespace Management --- */

func (psm *sanStorageManager) GetNamespace(subsystemID int, namespaceID int) (*model.Namespace, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_BEGIN), namespaceID, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	ns, err := gatewayObj.GetNamespace(subsystemID, namespaceID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NAMESPACE_FAILED), namespaceID, subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_END), namespaceID, subsystemID, psm.storageSetting.Serial)
	return ns, nil
}

func (psm *sanStorageManager) GetAllNamespaces(subsystemID int) (*model.Namespaces, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_NAMESPACES_BEGIN), subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	namespaces, err := gatewayObj.GetAllNamespaces(subsystemID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_ALL_NAMESPACES_FAILED), subsystemID, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_ALL_NAMESPACES_END), subsystemID, psm.storageSetting.Serial)
	return namespaces, nil
}

func (psm *sanStorageManager) CreateNamespace(request model.CreateNamespaceRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_NAMESPACE_BEGIN), request.NvmSubsystemId, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.CreateNamespace(request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_CREATE_NAMESPACE_FAILED), request.NvmSubsystemId, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_CREATE_NAMESPACE_END), request.NvmSubsystemId, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) SetNamespaceNickname(subsystemID int, namespaceID int, request model.SetNamespaceNicknameRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_NAMESPACE_NICKNAME_BEGIN), namespaceID, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.SetNamespaceNickname(subsystemID, namespaceID, request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_SET_NAMESPACE_NICKNAME_FAILED), namespaceID, subsystemID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_SET_NAMESPACE_NICKNAME_END), namespaceID, subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteNamespace(subsystemID int, namespaceID int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NAMESPACE_BEGIN), namespaceID, subsystemID, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteNamespace(subsystemID, namespaceID)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_NAMESPACE_FAILED), namespaceID, subsystemID, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NAMESPACE_END), namespaceID, subsystemID, psm.storageSetting.Serial)
	return resIDs, nil
}

/* --- Namespace Path (Mapping) Management --- */

func (psm *sanStorageManager) GetNamespacePaths(params model.GetNamespacePathsParams) (*model.NamespacePathsResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_PATHS_BEGIN), params.NvmSubsystemId, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	paths, err := gatewayObj.GetNamespacePaths(params)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NAMESPACE_PATHS_FAILED), params.NvmSubsystemId, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_PATHS_END), params.NvmSubsystemId, psm.storageSetting.Serial)
	return paths, nil
}

func (psm *sanStorageManager) GetNamespacePathDetail(subsystemId int, hostNqn string, namespaceId int) (*model.NamespacePath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_PATH_DETAIL_BEGIN), namespaceId, hostNqn, subsystemId, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return nil, err
	}

	pathDetail, err := gatewayObj.GetNamespacePathDetail(subsystemId, hostNqn, namespaceId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_GET_NAMESPACE_PATH_DETAIL_FAILED), namespaceId, hostNqn, subsystemId, psm.storageSetting.Serial)
		return nil, err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_GET_NAMESPACE_PATH_DETAIL_END), namespaceId, hostNqn, subsystemId, psm.storageSetting.Serial)
	return pathDetail, nil
}

func (psm *sanStorageManager) RegisterNamespacePath(request model.RegisterNamespacePathRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_REGISTER_NAMESPACE_PATH_BEGIN), request.NamespaceId, request.HostNqn, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.RegisterNamespacePath(request)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_REGISTER_NAMESPACE_PATH_FAILED), request.NamespaceId, request.HostNqn, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_REGISTER_NAMESPACE_PATH_END), request.NamespaceId, request.HostNqn, psm.storageSetting.Serial)
	return resIDs, nil
}

func (psm *sanStorageManager) DeleteNamespacePath(subsystemId int, hostNqn string, namespaceId int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NAMESPACE_PATH_BEGIN), namespaceId, hostNqn, psm.storageSetting.Serial)

	gatewayObj, err := gatewayimpl.NewEx(psm.getDeviceSettings())
	if err != nil {
		return "", err
	}

	resIDs, err := gatewayObj.DeleteNamespacePath(subsystemId, hostNqn, namespaceId)
	if err != nil {
		log.WriteError(mc.GetMessage(mc.ERR_DELETE_NAMESPACE_PATH_FAILED), namespaceId, hostNqn, psm.storageSetting.Serial)
		return "", err
	}
	log.WriteInfo(mc.GetMessage(mc.INFO_DELETE_NAMESPACE_PATH_END), namespaceId, hostNqn, psm.storageSetting.Serial)
	return resIDs, nil
}
