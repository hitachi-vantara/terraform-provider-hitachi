package vssbstorage

import (
	"fmt"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/gateway/model"
	"testing"
)

func newVolumeTestManager() (*vssbStorageManager, error) {

	objStorage := vssbmodel.StorageDeviceSettings{
		Username:       "admin",
		Password:       "vssb-789",
		ClusterAddress: "10.76.47.55",
	}
	psm, err := newVssbStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newVssbStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestCreateVolume
func xTestCreateVolume(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	var capacity int32 = 55
	poolID := "c89c1cbc-e181-4db3-ae44-eb7c6e496192"
	name := "test-volume_new_case"
	nickname := vssbmodel.NameParam{
		BaseName: &name,
	}
	name2 := vssbmodel.NickNameParam{
		BaseName: &name,
	}
	reqBody := vssbmodel.CreateVolumeRequestGwy{
		Capacity:      &capacity,
		PoolID:        &poolID,
		NameParam:     nickname,
		NickNameParam: name2,
	}
	resp, err := psm.CreateVolume(&reqBody)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetAllVolumes
func xTestGetAllVolumes(t *testing.T) {
	psm, err := newVolumeTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetAllVolumes()
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
	volName := "0937f399-6ab2-4c44-af16-3260e069f007"

	err = psm.DeleteVolume(&volName)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", err)
}

/*
VOLUME OPUTPUT:
-------------
{"data":[
{
"savingEffects":
{
"systemDataCapacity":0,
"preCapacityDataReductionWithoutSystemData":0,
"postCapacityDataReduction":0
},
"id":"246b4850-8e8b-4c01-8a94-9faf8f470c75",
"name":"Mongonode3_vol3",
"nickname":"Mongonode3_vol3",
"volumeNumber":12,
"poolId":"c89c1cbc-e181-4db3-ae44-eb7c6e496192",
"poolName":"SP01",
"totalCapacity":1258292,
"usedCapacity":0,
"numberOfConnectingServers":1,
"numberOfSnapshots":0,
"protectionDomainId":"4f762205-692a-490e-99af-d2f04dad2942",
"fullAllocated":false,
"volumeType":"Normal",
"statusSummary":"Normal",
"status":"Normal",
"storageControllerId":"174dc06f-5cc2-4626-8ac3-b7973c91b809",
"snapshotAttribute":"-",
"snapshotStatus":null,
"savingSetting":"Disabled",
"savingMode":null,
"dataReductionStatus":"Disabled",
"dataReductionProgressRate":null
},
...
*/
