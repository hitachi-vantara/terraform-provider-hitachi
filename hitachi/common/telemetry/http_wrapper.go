package telemetry

import (
	"reflect"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"time"
)

func WrapMethod(methodName string, method interface{}) interface{} {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	telemetryConsent := CheckTelemetryConsent()
	if !telemetryConsent {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return method
	}

	return reflect.MakeFunc(reflect.TypeOf(method), func(args []reflect.Value) []reflect.Value {
		log.WriteEnter()
		defer log.WriteExit()

		// If the method accepts variadic parameters, we handle it explicitly
		methodType := reflect.TypeOf(method)
		if methodType.In(methodType.NumIn()-1).Kind() == reflect.Slice {
			args = unpackVariadicArguments(args)
		}

		startTime := time.Now()

		// do the http call
		result := reflect.ValueOf(method).Call(args)

		endTime := time.Now()

		elapsedTime := endTime.Sub(startTime).Seconds()

		status, _ := checkHttpResult(result, methodName)

		firstArg := getHttpCallArgument(args, 0) // for storageSettings

		if status == "success" {
			thirdArg := getHttpCallArgument(args, 2) // for storageModel, apiversion
			UpdateTelemetryStats(status, elapsedTime, firstArg, thirdArg)
		} else {
			UpdateTelemetryStats(status, elapsedTime, firstArg, nil)
		}

		return result
	}).Interface()
}

func getHttpCallArgument(args []reflect.Value, i int) interface{} {
	if len(args) > 0 && args[0].Kind() == reflect.Struct {
		return args[i].Interface()
	}
	return nil
}

// Function to handle unpacking of variadic arguments (args...)
func unpackVariadicArguments(args []reflect.Value) []reflect.Value {
	// If there are no arguments, just return the original args
	if len(args) == 0 {
		return args
	}

	// Get the last argument, which should be a slice if it's variadic
	lastArg := args[len(args)-1]
	if lastArg.Kind() == reflect.Slice {
		// Access the slice containing the variadic arguments
		sliceValues := lastArg

		// Create a new slice of reflect.Value to hold the individual variadic arguments
		unpackedArgs := make([]reflect.Value, 0, sliceValues.Len())

		// Iterate over the slice and append each element as reflect.Value
		for i := 0; i < sliceValues.Len(); i++ {
			unpackedArgs = append(unpackedArgs, sliceValues.Index(i))
		}

		// Combine the original arguments (except the last one) and the unpacked variadic arguments
		args = append(args[:len(args)-1], unpackedArgs...)
	}

	// Return the updated args with the variadic arguments unpacked
	return args
}

func checkHttpResult(result []reflect.Value, methodName string) (string, int) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	resultCount := 0
	for i := 0; i < len(result); i++ {
		if !result[i].IsNil() {
			resultCount++
		}
	}

	log.WriteDebug("Method %s: Result length (counting non-nil): %d", methodName, resultCount)

	if methodName == "GetCall" {
		// GetCall returns only error which could be nil (nil are not counted in resultCount)
		// If result length is 0 or it's nil, treat it as success
		if resultCount == 0 || (len(result) > 0 && result[0].IsNil()) {
			log.WriteDebug("Method %s executed successfully (no data or no error)", methodName)
			return "success", resultCount
		} else {
			// If we have a result and it's not nil, consider failure
			err, ok := result[len(result)-1].Interface().(error)
			if ok && err != nil {
				log.WriteError("Error in method %s: %v", methodName, err)
				return "failure", resultCount
			}
			return "failure", resultCount
		}
	}

	if resultCount > 0 {
		// other http calls return (*data, error) (either: nil,err or somedata,nil)
		// Check if the last return value is an error
		err, ok := result[len(result)-1].Interface().(error)
		if ok && err != nil {
			log.WriteError("Error in method %s: %v", methodName, err)
			return "failure", resultCount
		}

		// If result[0] is nil, treat it as failure
		if result[0].IsNil() {
			log.WriteError("Method %s returned nil data", methodName)
			return "failure", resultCount
		}

		// Otherwise, consider it success
		log.WriteDebug("Method %s executed successfully", methodName)
		return "success", resultCount
	}

	// If there’s no data, consider failure
	log.WriteError("Method %s executed without returning data", methodName)
	return "failure", resultCount
}
