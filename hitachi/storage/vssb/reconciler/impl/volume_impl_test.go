package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vssb/reconciler/model"
	"testing"
)

func newVolumeTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "admin",    //"YWRtaW4=",
		Password:       "vssb-789", //"dnNzYi03ODk=",
		ClusterAddress: "10.76.47.55",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetAllVolumes
func xTestGetAllVolumes(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	//computeNodeName := ""
	computeNodeName := "esxi-151"
	resp, err := psm.GetAllVolumes(computeNodeName)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllVolumes
func xTestGetVolumeDetails(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	//computeNodeName := ""
	volumeName := "Mongonode3_vol3"
	resp, err := psm.GetVolumeDetails(volumeName)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateVolume
func xTestCreateVolume(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	name := "test_volume_name"
	poolName := "SP01"
	var capacity float32 = 0.50
	Nickname := "test_volume_nick_name"
	computNodes := []string{}

	volumeI := vssbmodel.CreateVolume{
		Name:         &name,
		PoolName:     &poolName,
		CapacityInGB: &capacity,
		NickName:     &Nickname,
		ComputeNodes: computNodes,
	}
	resp, err := psm.ReconcileVolume(&volumeI)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestDeleteVolume
func xTestDeleteVolume(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	id := "dc062f96-2844-4902-a69e-cbc94bbe5fb6"
	err = psm.DeleteVolumeResource(&id)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
}
