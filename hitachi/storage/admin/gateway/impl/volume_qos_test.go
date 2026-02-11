package admin

import (
	adminmodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
	"testing"
)

// go test -v -run TestGetVolumeQosAdminInfo
func xTestGetVolumeQosAdminInfo(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	volumeID := 101
	resp, err := psm.GetVolumeQosAdminInfo(volumeID)
	if err != nil {
		t.Errorf("Unexpected error in GetVolumeQosAdminInfo: %v", err)
		return
	}
	t.Logf("QoS Response: %+v", resp)
}

// go test -v -run TestStorageVolumeQosThresholdSettings
func xTestStorageVolumeQosThresholdSettings(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	volumeId := 101
	threshold := adminmodel.VolumeQosThreshold{
		IsUpperIopsEnabled:         true,
		UpperIops:                  5000,
		IsUpperTransferRateEnabled: true,
		UpperTransferRate:          200,
		IsLowerIopsEnabled:         true,
		LowerIops:                  1000,
		IsLowerTransferRateEnabled: true,
		LowerTransferRate:          50,
		IsResponsePriorityEnabled:  true,
		ResponsePriority:           3,
		TargetResponseTime:         10,
	}

	err = psm.SetVolumeQosAdminThreshold(volumeId, threshold)
	if err != nil {
		t.Errorf("Unexpected error in SetVolumeQosAdminThreshold: %v", err)
		return
	}
	t.Logf("Successfully updated QoS threshold for volume ID: %d", volumeId)
}

// go test -v -run TestStorageVolumeQosAlertSettings
func xTestStorageVolumeQosAlertSettings(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	volumeID := 101
	alert := adminmodel.VolumeQosAlertSetting{
		IsUpperAlertEnabled:        true,
		UpperAlertAllowableTime:    5,
		IsLowerAlertEnabled:        true,
		LowerAlertAllowableTime:    5,
		IsResponseAlertEnabled:     true,
		ResponseAlertAllowableTime: 5,
	}

	err = psm.SetVolumeQosAdminAlertSetting(volumeID, alert)
	if err != nil {
		t.Errorf("Unexpected error in SetVolumeQosAdminAlertSetting: %v", err)
		return
	}
	t.Logf("Successfully updated QoS alert settings for volume ID: %d", volumeID)
}
