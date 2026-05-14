package sanstorage

import (
	"fmt"
	"net/url"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/san/gateway/http"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
)

// GetNvmSubsystem gets information for a specific NVM subsystem by its ID
func (psm *sanStorageManager) GetNvmSubsystem(subsystemID int) (*sanmodel.NvmSubsystem, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var nvmSubsystem sanmodel.NvmSubsystem
	// API endpoint: v1/objects/nvm-subsystems/{object-ID}
	apiSuf := fmt.Sprintf("objects/nvm-subsystems/%d", subsystemID)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &nvmSubsystem)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &nvmSubsystem, nil
}

// GetAllNvmSubsystems gets all NVM subsystems.
// infoType can be: "basic", "nqn", "namespace", or "port". If empty, "basic" is assumed.
func (psm *sanStorageManager) GetAllNvmSubsystems(params sanmodel.GetNvmSubsystemsParams) (*sanmodel.NvmSubsystems, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var nvmSubsystems sanmodel.NvmSubsystems

	q := url.Values{}

	// (Optional) Type of information: basic, nqn, namespace, port
	if params.NvmSubsystemInfo != nil {
		q.Add("nvmSubsystemInfo", *params.NvmSubsystemInfo)
	}

	// (Optional) Condition: undefined
	if params.NvmSubsystemOption != nil {
		q.Add("nvmSubsystemOption", *params.NvmSubsystemOption)
	}

	log.WriteDebug("TFDebug| NvmSubsystem QueryParams:%+v", q)

	apiSuf := "objects/nvm-subsystems"
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &nvmSubsystems)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &nvmSubsystems, nil
}

// CreateNvmSubsystem creates a new NVM subsystem and returns a Job object
func (psm *sanStorageManager) CreateNvmSubsystem(request sanmodel.CreateNvmSubsystemRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/nvm-subsystems"

	log.WriteDebug("TFDebug| Creating NVM Subsystem ID: %d", request.NvmSubsystemId)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// UpdateNvmSubsystem updates settings for an existing NVM subsystem.
// Note: Only one attribute in the request struct should be non-nil.
func (psm *sanStorageManager) UpdateNvmSubsystem(subsystemID int, request sanmodel.UpdateNvmSubsystemRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/nvm-subsystems/%d", subsystemID)

	log.WriteDebug("TFDebug| Updating NVM Subsystem ID: %d", subsystemID)

	// Using your pattern where PatchCall returns (*string, error)
	resIds, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteNvmSubsystem deletes an NVM subsystem and returns the resource ID string
func (psm *sanStorageManager) DeleteNvmSubsystem(subsystemID int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/nvm-subsystems/%d", subsystemID)

	log.WriteDebug("TFDebug| Deleting NVM Subsystem ID: %d", subsystemID)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

/* --- NVM Subsystem Port Management --- */

// GetNvmSubsystemPort gets information for a specific port within an NVM subsystem
func (psm *sanStorageManager) GetNvmSubsystemPort(subsystemID int, portID string) (*sanmodel.NvmSubsystemPort, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var portDetail sanmodel.NvmSubsystemPort
	// Object ID pattern: nvmSubsystemId,portId
	apiSuf := fmt.Sprintf("objects/nvm-subsystem-ports/%d,%s", subsystemID, portID)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &portDetail)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &portDetail, nil
}

// GetNvmSubsystemPorts gets detailed port information for a specific NVM subsystem
func (psm *sanStorageManager) GetNvmSubsystemPorts(subsystemID int) (*sanmodel.NvmSubsystemPorts, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var nvmSubsystemPorts sanmodel.NvmSubsystemPorts
	q := url.Values{}
	q.Add("nvmSubsystemId", fmt.Sprintf("%d", subsystemID))
	q.Add("nvmSubsystemInfo", "port")

	apiSuf := fmt.Sprintf("objects/nvm-subsystems?%s", q.Encode())

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &nvmSubsystemPorts)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &nvmSubsystemPorts, nil
}

// AddNvmSubsystemPort maps a port to an NVM subsystem.
func (psm *sanStorageManager) AddNvmSubsystemPort(request sanmodel.AddNvmSubsystemPortRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/nvm-subsystem-ports"

	log.WriteDebug("TFDebug| Adding Port %s to NVM Subsystem ID: %d", request.PortId, request.NvmSubsystemId)

	// Using your pattern where PostCall returns (*string, error)
	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteNvmSubsystemPort deletes (unmaps) a port from an NVM subsystem
func (psm *sanStorageManager) DeleteNvmSubsystemPort(subsystemID int, portID string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Object ID pattern: nvmSubsystemId,portId
	apiSuf := fmt.Sprintf("objects/nvm-subsystem-ports/%d,%s", subsystemID, portID)

	log.WriteDebug("TFDebug| Deleting Port %s from NVM Subsystem ID: %d", portID, subsystemID)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

/* --- Host NQN Management --- */

// GetAllHostNqns gets host NQN information for a specific NVM subsystem.
func (psm *sanStorageManager) GetAllHostNqns(subsystemID int) (*sanmodel.HostNqns, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostNqns sanmodel.HostNqns
	q := url.Values{}
	q.Add("nvmSubsystemId", fmt.Sprintf("%d", subsystemID))

	apiSuf := "objects/host-nqns"
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hostNqns)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &hostNqns, nil
}

// GetHostNqn gets information about a specific host NQN.
func (psm *sanStorageManager) GetHostNqn(subsystemID int, hostNqn string) (*sanmodel.HostNqnInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var hostNqnInfo sanmodel.HostNqnInfo
	// Object ID: nvmSubsystemId,hostNqn
	apiSuf := fmt.Sprintf("objects/host-nqns/%d,%s", subsystemID, hostNqn)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &hostNqnInfo)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &hostNqnInfo, nil
}

// RegisterHostNqn registers a host NQN in an NVM subsystem.
func (psm *sanStorageManager) RegisterHostNqn(request sanmodel.RegisterHostNqnRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/host-nqns"

	log.WriteDebug("TFDebug| Registering Host NQN: %s in Subsystem: %d", request.HostNqn, request.NvmSubsystemId)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// SetHostNqnNickname sets or deletes the nickname for a host NQN.
func (psm *sanStorageManager) SetHostNqnNickname(subsystemID int, hostNqn string, request sanmodel.SetHostNqnNicknameRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Object ID: nvmSubsystemId,hostNqn
	apiSuf := fmt.Sprintf("objects/host-nqns/%d,%s", subsystemID, hostNqn)

	log.WriteDebug("TFDebug| Updating Host NQN Nickname. Subsystem: %d, NQN: %s", subsystemID, hostNqn)

	resIds, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteLoginHostNqn deletes the login information of the host NQN for a specified port.
func (psm *sanStorageManager) DeleteLoginHostNqn(portID string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// API Endpoint: objects/ports/{portId}/actions/delete-login-host-nqn/invoke
	apiSuf := fmt.Sprintf("objects/ports/%s/actions/delete-login-host-nqn/invoke", portID)

	log.WriteDebug("TFDebug| Deleting login information for Port: %s", portID)

	// Note: This is an 'invoke' action, usually sent via POST with no body
	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteHostNqn deletes a host NQN registered in an NVM subsystem.
func (psm *sanStorageManager) DeleteHostNqn(subsystemID int, hostNqn string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// Object ID: nvmSubsystemId,hostNqn
	apiSuf := fmt.Sprintf("objects/host-nqns/%d,%s", subsystemID, hostNqn)

	log.WriteDebug("TFDebug| Deleting Host NQN. Subsystem: %d, NQN: %s", subsystemID, hostNqn)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

/* --- Namespace Management --- */

// GetAllNamespaces gets namespace information for a specific NVM subsystem.
func (psm *sanStorageManager) GetAllNamespaces(subsystemID int) (*sanmodel.Namespaces, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var namespaces sanmodel.Namespaces
	q := url.Values{}
	q.Add("nvmSubsystemId", fmt.Sprintf("%d", subsystemID))

	apiSuf := "objects/namespaces"
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &namespaces)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &namespaces, nil
}

// GetNamespace gets information about a specific namespace.
func (psm *sanStorageManager) GetNamespace(subsystemID int, namespaceID int) (*sanmodel.Namespace, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var namespace sanmodel.Namespace
	// Object ID: nvmSubsystemId,namespaceId
	apiSuf := fmt.Sprintf("objects/namespaces/%d,%d", subsystemID, namespaceID)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &namespace)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &namespace, nil
}

// CreateNamespace creates a namespace in the NVM subsystem by using a specified LDEV number.
func (psm *sanStorageManager) CreateNamespace(request sanmodel.CreateNamespaceRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/namespaces"

	log.WriteDebug("TFDebug| Creating Namespace for Subsystem: %d, LdevId: %d", request.NvmSubsystemId, request.LdevId)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// SetNamespaceNickname sets or deletes the nickname for a namespace.
func (psm *sanStorageManager) SetNamespaceNickname(subsystemID int, namespaceID int, request sanmodel.SetNamespaceNicknameRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/namespaces/%d,%d", subsystemID, namespaceID)

	log.WriteDebug("TFDebug| Updating Nickname for Subsystem: %d, Namespace: %d", subsystemID, namespaceID)

	resIds, err := httpmethod.PatchCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteNamespace deletes the namespace in the NVM subsystem.
func (psm *sanStorageManager) DeleteNamespace(subsystemID int, namespaceID int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/namespaces/%d,%d", subsystemID, namespaceID)

	log.WriteDebug("TFDebug| Deleting Namespace. Subsystem: %d, Namespace: %d", subsystemID, namespaceID)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

/* --- Namespace Path (Mapping) Management --- */

// GetNamespacePaths gets host-namespace path information.
func (psm *sanStorageManager) GetNamespacePaths(params sanmodel.GetNamespacePathsParams) (*sanmodel.NamespacePathsResponse, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var response sanmodel.NamespacePathsResponse
	q := url.Values{}

	q.Add("nvmSubsystemId", fmt.Sprintf("%d", params.NvmSubsystemId))
	if params.NamespaceId != nil {
		q.Add("namespaceId", fmt.Sprintf("%d", *params.NamespaceId))
	}
	if params.HeadId != nil {
		q.Add("headId", *params.HeadId)
	}

	apiSuf := "objects/namespace-paths"
	if len(q) > 0 {
		apiSuf = fmt.Sprintf("%s?%s", apiSuf, q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &response)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &response, nil
}

// GetNamespacePathDetail gets specific host-namespace path information.
func (psm *sanStorageManager) GetNamespacePathDetail(subsystemId int, hostNqn string, namespaceId int) (*sanmodel.NamespacePath, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var pathDetail sanmodel.NamespacePath
	// Object ID: nvmSubsystemId,hostNqn,namespaceId
	apiSuf := fmt.Sprintf("objects/namespace-paths/%d,%s,%d", subsystemId, hostNqn, namespaceId)

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &pathDetail)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}

	return &pathDetail, nil
}

// RegisterNamespacePath sets the host-namespace path (mapping).
func (psm *sanStorageManager) RegisterNamespacePath(request sanmodel.RegisterNamespacePathRequest) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := "objects/namespace-paths"

	log.WriteDebug("TFDebug| Creating Namespace Path for Subsystem: %d, Namespace: %d", request.NvmSubsystemId, request.NamespaceId)

	resIds, err := httpmethod.PostCall(psm.storageSetting, apiSuf, request)
	if err != nil {
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}

// DeleteNamespacePath deletes a host-namespace path.
func (psm *sanStorageManager) DeleteNamespacePath(subsystemId int, hostNqn string, namespaceId int) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/namespace-paths/%d,%s,%d", subsystemId, hostNqn, namespaceId)

	log.WriteDebug("TFDebug| Deleting Namespace Path. Subsystem: %d, NQN: %s, Namespace: %d", subsystemId, hostNqn, namespaceId)

	resIds, err := httpmethod.DeleteCall(psm.storageSetting, apiSuf, nil)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return "", err
	}

	log.WriteDebug("TFDebug | Resource IDs = %v", *resIds)
	return *resIds, nil
}
