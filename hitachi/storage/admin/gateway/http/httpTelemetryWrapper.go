package admin

import (
	"strings"
	"time"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/telemetry"
	model "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

func GetCall(storageSetting model.StorageDeviceSettings, apiSuf string, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return getCall(storageSetting, apiSuf, output)
	}

	startTime := time.Now()
	err := getCall(storageSetting, apiSuf, output)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil {
		status = "success"
	}

	// only for get call
	outputForModelOrVersion := interface{}(nil)
	if apiSuf == "objects/storages/instance" || apiSuf == "configuration/version" || strings.HasPrefix(apiSuf, "objects/storage"){
		outputForModelOrVersion = output
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, outputForModelOrVersion)
	return err
}

func PostCall(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return postCall(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := postCall(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}

func PatchCall(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return patchCallAsync(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := patchCallAsync(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}

func PatchCallSync(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return patchCallSync(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := patchCallSync(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}

func DeleteCall(storageSetting model.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return deleteCall(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := deleteCall(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}
