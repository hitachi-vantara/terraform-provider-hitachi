package utils

import (
	"errors"
	"fmt"

	log "github.com/romana/rlog"
)

func GetFormatedError(err error, message string, value interface{}) error {
	errStr := "error is nil"
	if err != nil {
		errStr = err.Error()
	}
	newerr := errors.New(fmt.Sprintf("%s %+v ::: ERROR: %s", message, value, errStr))
	log.Error(newerr)
	return newerr
}