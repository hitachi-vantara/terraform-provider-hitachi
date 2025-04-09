package vssbstorage

import (
	"fmt"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	httpmethod "terraform-provider-hitachi/hitachi/storage/vssb/gateway/http"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/gateway/model"
)

// ChangeUserPassword is used to change the password of a user
func (psm *vssbStorageManager) ChangeUserPassword(userId string, reqBody *vssbmodel.ChangeUserPasswordReq) (*vssbmodel.User, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var userInfo vssbmodel.User

	apiSuf := fmt.Sprintf("objects/users/%s/password", userId)
	sensitiveLogging := false
	err := httpmethod.PatchCallSync(psm.storageSetting, apiSuf, reqBody, &userInfo, sensitiveLogging)
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in %s API call, err: %v", apiSuf, err)
		return nil, err
	}
	return &userInfo, nil
}
