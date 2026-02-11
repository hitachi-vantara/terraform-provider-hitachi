package common

import (
	"errors"
	"reflect"
)

// ILogger - interface to define logging methods
type ILogger interface {
	WriteEnter(a ...interface{})
	WriteParam(format string, value interface{})

	WriteInfo(message interface{}, a ...interface{})
	WriteWarn(message interface{}, a ...interface{})
	WriteError(message interface{}, a ...interface{})

	WriteDebug(message string, a ...interface{})
	WriteExit()

	WriteAndReturnError(message interface{}, a ...interface{}) error
}

var loggerInstance ILogger = NewDefaultLogger()

func GetLogger() ILogger {
	return loggerInstance
}

func SetLogger(logger ILogger) {
	loggerInstance = logger
}

var messageCatalogs map[reflect.Type]*map[interface{}]string = make(map[reflect.Type]*map[interface{}]string)

func AddMessageCatalog(t reflect.Type, messageCatalog *map[interface{}]string) {
	if messageCatalogs == nil {
		return
	}
	messageCatalogs[t] = messageCatalog

}

func GetMessage(messageID interface{}) (string, error) {
	messageCatalogIDType := reflect.TypeOf(messageID)
	messageCatalog, isMessageCatalogPresent := messageCatalogs[messageCatalogIDType]
	if !isMessageCatalogPresent {
		return "", errors.New("message Catalog not found")
	}
	message := (*messageCatalog)[messageID]
	return message, nil
}
