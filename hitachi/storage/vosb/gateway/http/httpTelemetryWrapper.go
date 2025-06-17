package vssbstorage

import (
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/telemetry"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	"time"
)

func GetCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, output interface{}) error {
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
	if apiSuf == "objects/storages/instance" || apiSuf == "configuration/version" {
		outputForModelOrVersion = output
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, outputForModelOrVersion)
	return err
}

func PostCallForm(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return postCallForm(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := postCallForm(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}

func PostCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

func PatchCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return patchCall(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := patchCall(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}

func PatchCallSync(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}, output interface{}) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return patchCallSync(storageSetting, apiSuf, reqBody, output)
	}

	startTime := time.Now()
	err := patchCallSync(storageSetting, apiSuf, reqBody, output)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return err
}

func DeleteCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
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

func DownloadFileCall(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, toFilePath string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return downloadFileCall(storageSetting, apiSuf, toFilePath)
	}

	startTime := time.Now()
	filePath, err := downloadFileCall(storageSetting, apiSuf, toFilePath)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return filePath, err
}


func PostCallFormExt(storageSetting vssbmodel.StorageDeviceSettings, apiSuf string, reqBody interface{}) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if !telemetry.CheckTelemetryConsent() {
		log.WriteInfo("Telemetry consent not given. Skipping telemetry tracking.")
		return postCallForm(storageSetting, apiSuf, reqBody)
	}

	startTime := time.Now()
	result, err := postCallForm(storageSetting, apiSuf, reqBody)
	elapsedTime := time.Since(startTime).Seconds()

	status := "failure"
	if err == nil && result != nil {
		status = "success"
	}

	telemetry.UpdateTelemetryStats(status, elapsedTime, storageSetting, nil)
	return result, err
}