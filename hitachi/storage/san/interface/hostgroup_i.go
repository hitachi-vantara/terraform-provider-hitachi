package sanstorage

import (
	// "encoding/json"
	// "errors"
	// "context"
	// "fmt"
	// "io/ioutil"
	// "strconv"
	// "time"
	// "sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	// "terraform-provider-hitachi/hitachi/common/utils"
	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	// sangateway "terraform-provider-hitachi/hitachi/storage/san/gateway"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func (ssm *sanstorageManager) GetHostGroupI(portID string, hostGroupNumber int) (*sanmodel.HostGroup, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	return nil, nil

}
