package admin

import (
	"fmt"
	"net/url"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/admin/gateway/http"
	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// GetPorts retrieves all ports based on query parameters
func (psm *adminStorageManager) GetPorts(params gwymodel.GetPortParams) (*gwymodel.PortInfoList, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var portInfoList gwymodel.PortInfoList

	// Build query string dynamically based on non-nil fields
	q := url.Values{}
	if params.Protocol != nil {
		q.Add("protocol", *params.Protocol)
	}

	log.WriteDebug("TFDebug| QueryParams:%+v", q)

	// apiSuf := fmt.Sprintf("objects/ports?%s", q.Encode())
	apiSuf := fmt.Sprintf("objects/ports")

	if len(q) > 0 {
		apiSuf = fmt.Sprintf("objects/ports?%s", q.Encode())
	}

	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &portInfoList)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	return &portInfoList, nil
}

// GetPortByID retrieves a specific port by its ID
func (psm *adminStorageManager) GetPortByID(portID string) (*gwymodel.PortInfo, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var portInfo gwymodel.PortInfo

	apiSuf := fmt.Sprintf("objects/ports/%s", portID)
	err := httpmethod.GetCall(psm.storageSetting, apiSuf, &portInfo)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to get port with id %s: %w", portID, err))
		return nil, err
	}

	return &portInfo, nil
}

// UpdatePort updates port configuration by its ID
func (psm *adminStorageManager) UpdatePort(portID string, params gwymodel.UpdatePortParams) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	apiSuf := fmt.Sprintf("objects/ports/%s", portID)
	_, err := httpmethod.PatchCallSync(psm.storageSetting, apiSuf, params)
	if err != nil {
		log.WriteError(fmt.Errorf("failed to update port with id %s: %w", portID, err))
		return err
	}

	return nil
}
