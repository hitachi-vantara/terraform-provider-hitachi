package utils

import (
	"errors"
	"fmt"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

func GetFormatedError(err error, message string, value interface{}) error {
	log := commonlog.GetLogger()

	errStr := "error is nil"
	if err != nil {
		errStr = err.Error()
	}
	newerr := errors.New(fmt.Sprintf("%s %+v ::: ERROR: %s", message, value, errStr))
	log.WriteError(newerr)
	return newerr
}