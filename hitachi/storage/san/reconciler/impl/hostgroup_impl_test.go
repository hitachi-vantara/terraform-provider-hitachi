package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/reconciler/model"
	"testing"
)

// newHgTestManager is for Testing and provide structure information for connection
func newHgTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   40014,
		Username: "bXNfdm13YXJl",
		Password: "SGl0YWNoaTE=",
		MgmtIP:   "172.25.47.115",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newSanStorageManagerEx %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetHostGroup
func xTestGetHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	resp, err := psm.GetHostGroup("CL1-A", 24)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestUpdateLun
func xTestUpdateLun(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING-HOSTGROUP"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	//  25, 12 ==> TF:  25, 13 : Delete 25, 12 and Add 25, 13 (if 13 already there then error)
	//  25, 12 ==> TF:  25 ==> Delete 25, 12 and Create new 25, XX
	//  25, 12 ==> 25, 12 | 21, 13 ==> Only Add 21, 13 and skip 23, 12
	//  25, 12 => TF: 21, 12    ( Another LDEV is already mapped to LUN) Give Error

	//Ldev0 := 25
	Ldev1 := 25
	Lun1 := 12
	//Lun1 := 13
	// For Single HG Lun ID will be unique
	Ldev2 := 21
	Lun2 := 13
	//Ldev0 := 21
	Ldev3 := 22
	Lun3 := 14
	myLdevIds := []sanmodel.Luns{
		//{LdevId: &Ldev0}, // Automatic Lun ID Assing
		{LdevId: &Ldev3, Lun: &Lun3},
		{LdevId: &Ldev1, Lun: &Lun1},
		{LdevId: &Ldev2, Lun: &Lun2},
	}

	/*
		myWwn := []sanmodel.Wwn{
			{Wwn: "100000109b3dfbbb", Name: "test-wwn1b"},
			{Wwn: "100000109b3dfbbc", Name: "test-wwn1c"},
		}
	*/

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Ldevs:           myLdevIds,
		//Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in ReconcileHostGroup %v", err)
		return
	}
	t.Logf("Response: %v", resp)

}

// go test -v -run TestUpdateWwn
func xTestUpdateWwn(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING-HOSTGROUP"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 32}
	myHostMode := "AIX"

	myWwn := []sanmodel.Wwn{
		{Wwn: "100000109b3dfbbb", Name: "test-wwn1b"},
		{Wwn: "100000109b3dfbbc", Name: "test-wwn1c"},
	}

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
	/*
					// Current WWN Info
						"data" : [ {
						"hostWwnId" : "CL1-A,23,100000109b3dfbbb",
						"portId" : "CL1-A",
						"hostGroupNumber" : 23,
						"hostGroupName" : "TESTING-HOSTGROUP",
						"hostWwn" : "100000109b3dfbbb",
						"wwnNickname" : "test-wwn1"
						}, {
						"hostWwnId" : "CL1-A,23,100000109b3dfaaa",
						"portId" : "CL1-A",
						"hostGroupNumber" : 23,
						"hostGroupName" : "TESTING-HOSTGROUP",
						"hostWwn" : "100000109b3dfaaa",
						"wwnNickname" : "test-wwn2"
						} ]

						============
					//Input Data
						{Wwn: "100000109b3dfbbb", Name: "test-wwn1b"}, => Modified
						{Wwn: "100000109b3dfbbc", Name: "test-wwn1c"}, => Add
						100000109b3dfaaa ==> Delete
				======
				OUTPUT:
				--
				"data" : [ {
			"hostWwnId" : "CL1-A,23,100000109b3dfbbb",
			"portId" : "CL1-A",
			"hostGroupNumber" : 23,
			"hostGroupName" : "TESTING-HOSTGROUP",
			"hostWwn" : "100000109b3dfbbb",
			"wwnNickname" : "test-wwn1b"
		}, {
			"hostWwnId" : "CL1-A,23,100000109b3dfbbc",
			"portId" : "CL1-A",
			"hostGroupNumber" : 23,
			"hostGroupName" : "TESTING-HOSTGROUP",
			"hostWwn" : "100000109b3dfbbc",
			"wwnNickname" : "test-wwn1c"
		}
	*/
	/*
		Curl commands:
		Auth:
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -u ms_vmware:Hitachi1 -X POST https://172.25.47.115/ConfigurationManager/v1/objects/sessions/ -d ""

		Delete
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X DELETE https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns/CL1-A,23,100000109b3dfbbb
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X DELETE https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns/CL1-A,23,100000109b3dfbbc

		Add:
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X POST https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns -d '{"hostWwn": "100000109b3dfbbb", "portId": "CL1-A", "hostGroupNumber": 23 }'
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X POST https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns -d '{"hostWwn": "100000109b3dfaaa", "portId": "CL1-A", "hostGroupNumber": 23 }'
			NickName:
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X PATCH https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns/CL1-A,23,100000109b3dfbbb -d '{"wwnNickname": "test-wwn1"}'
			curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X PATCH https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns/CL1-A,23,100000109b3dfaaa -d '{"wwnNickname": "test-wwn2"}'

		Get:
			curl -v -k -H "Accept:application/json" -H "Authorization:Session e085d3e0-7187-4839-a440-4b5296b7440e" -X GET "https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns?portId=CL1-A&hostGroupNumber=23"
	*/
}

// go test -v -run TestUpdateHostGroupMode
func xTestUpdateHostGroupMode(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	/*
		1) Hostgroup Mode and Option changes
		2) Hostgroup Mode removed from TF but available in existing HG
		3) Hostgroup Option removed from TF but available in existing HG
		4) Both removed
	*/

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 13}
	myHostMode := "LINUX/IRIX"

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestCreateHostGroup
func xTestCreateHostGroup(t *testing.T) {
	psm, err := newHgTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	myHostGroupNumber := 23
	myHostModeOptions := []int{12, 13}
	myHostMode := "AIX"

	//Ldev0 := 25
	Ldev1 := 25
	Lun1 := 12
	// For Single HG Lun ID will be unique
	//Ldev2 := 21
	//Lun2 := 13
	myLdevIds := []sanmodel.Luns{
		//{LdevId: &Ldev0},  // Automatic Lun ID Assing
		{LdevId: &Ldev1, Lun: &Lun1},
		//{LdevId: &Ldev2, Lun: &Lun2},
	}

	myWwn := []sanmodel.Wwn{
		{Wwn: "100000109b3dfbbb", Name: "test-name"},
	}

	//==============
	crReq := sanmodel.CreateHostGroupRequest{
		PortID:          &myPortID,
		HostGroupName:   &myHostGroupName,
		HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
		Ldevs:           myLdevIds,
		Wwns:            myWwn,
	}

	resp, err := psm.ReconcileHostGroup(&crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

/*
DELETE HG Will DELETE Mapping Also:
 curl -v -k -H "Accept:application/json" -H "Content-Type:application/json" -H "Authorization:Session 8148e588-8d1b-43da-86b8-6a0356cbe585" -X DELETE https://172.25.47.115/ConfigurationManager/v1/objects/host-groups/CL1-A,23
 --
GET ALL VOLUME
curl -v -k -H "Accept:application/json" -H "Authorization:Session f876ebcf-4dbd-4d1f-83f9-107b9b11a75c" -X GET https://172.25.47.115/ConfigurationManager/v1/objects/ldevs

 {
    "ldevId" : 25,
    "clprId" : 0,
    "emulationType" : "OPEN-V-CVS",
    "byteFormatCapacity" : "1.00 T",
    "blockCapacity" : 2147483648,
    "numOfPorts" : 1,
    "ports" : [ {
      "portId" : "CL1-A",
      "hostGroupNumber" : 7,
      "hostGroupName" : "test-nvme",
      "lun" : 0
    } ],
    "attributes" : [ "CVS", "HDT" ],
    "label" : "SV-569",
    "status" : "NML",
    "mpBladeId" : 3,
    "ssid" : "0004",
    "poolId" : 1,
    "numOfUsedBlock" : 86016,
    "isRelocationEnabled" : true,
    "tierLevel" : "all",
    "usedCapacityPerTierLevel1" : 42,
    "tierLevelForNewPageAllocation" : "M",
    "isFullAllocationEnabled" : false,
    "resourceGroupId" : 4,
    "dataReductionStatus" : "ENABLED",
    "dataReductionMode" : "compression",
    "dataReductionProcessMode" : "inline",
    "isAluaEnabled" : false,
    "isCompressionAccelerationEnabled" : true,
    "compressionAccelerationStatus" : "ENABLED"
  },
  ============================
  AFTER LDEV-LUN Adding:
  ============================
  {
    "ldevId" : 25,
    "clprId" : 0,
    "emulationType" : "OPEN-V-CVS",
    "byteFormatCapacity" : "1.00 T",
    "blockCapacity" : 2147483648,
    "numOfPorts" : 2,
    "ports" : [ {
      "portId" : "CL1-A",
      "hostGroupNumber" : 7,
      "hostGroupName" : "test-nvme",
      "lun" : 0
    }, {
      "portId" : "CL1-A",
      "hostGroupNumber" : 23,
      "hostGroupName" : "TESTING_REST_API",
      "lun" : 12
    } ],
    "attributes" : [ "CVS", "HDT" ],
    "label" : "SV-569",
    "status" : "NML",
    "mpBladeId" : 3,
    "ssid" : "0004",
    "poolId" : 1,
    "numOfUsedBlock" : 86016,
    "isRelocationEnabled" : true,
    "tierLevel" : "all",
    "usedCapacityPerTierLevel1" : 42,
    "tierLevelForNewPageAllocation" : "M",
    "isFullAllocationEnabled" : false,
    "resourceGroupId" : 4,
    "dataReductionStatus" : "ENABLED",
    "dataReductionMode" : "compression",
    "dataReductionProcessMode" : "inline",
    "isAluaEnabled" : false,
    "isCompressionAccelerationEnabled" : true,
    "compressionAccelerationStatus" : "ENABLED"
  }

  ============
{
    "ldevId" : 21,
    "clprId" : 0,
    "emulationType" : "OPEN-V-CVS-CM",
    "byteFormatCapacity" : "52.00 M",
    "blockCapacity" : 106496,
    "numOfPorts" : 2,
    "ports" : [ {
      "portId" : "CL1-A",
      "hostGroupNumber" : 1,
      "hostGroupName" : "Prod-ESXi-Cluste",
      "lun" : 0
    }, {
      "portId" : "CL1-E",
      "hostGroupNumber" : 1,
      "hostGroupName" : "Prod-ESXi-Cluste",
      "lun" : 36
    } ],
    "attributes" : [ "CMD", "CVS" ],
    "raidLevel" : "RAID5",
    "raidType" : "3D+1P",
    "numOfParityGroups" : 1,
    "parityGroupIds" : [ "1-1" ],
    "driveType" : "SLB5G-M7R6SS",
    "driveByteFormatCapacity" : "6.98 T",
    "driveBlockCapacity" : 15000000128,
    "label" : "frank_commdev_SIE-5813",
    "status" : "NML",
    "mpBladeId" : 0,
    "ssid" : "0004",
    "resourceGroupId" : 0,
    "isAluaEnabled" : false
  },
  ===
  AFTER ADDING
  ===========
  {
    "ldevId" : 21,
    "clprId" : 0,
    "emulationType" : "OPEN-V-CVS-CM",
    "byteFormatCapacity" : "52.00 M",
    "blockCapacity" : 106496,
    "numOfPorts" : 3,
    "ports" : [ {
      "portId" : "CL1-A",
      "hostGroupNumber" : 1,
      "hostGroupName" : "Prod-ESXi-Cluste",
      "lun" : 0
    }, {
      "portId" : "CL1-E",
      "hostGroupNumber" : 1,
      "hostGroupName" : "Prod-ESXi-Cluste",
      "lun" : 36
    }, {
      "portId" : "CL1-A",
      "hostGroupNumber" : 23,
      "hostGroupName" : "TESTING_REST_API",
      "lun" : 13
    } ],
    "attributes" : [ "CMD", "CVS" ],
    "raidLevel" : "RAID5",
    "raidType" : "3D+1P",
    "numOfParityGroups" : 1,
    "parityGroupIds" : [ "1-1" ],
    "driveType" : "SLB5G-M7R6SS",
    "driveByteFormatCapacity" : "6.98 T",
    "driveBlockCapacity" : 15000000128,
    "label" : "frank_commdev_SIE-5813",
    "status" : "NML",
    "mpBladeId" : 0,
    "ssid" : "0004",
    "resourceGroupId" : 0,
    "isAluaEnabled" : false
  },
  =====
  curl -v -k -H "Accept:application/json" -H "Authorization:Session 3eca3926-79cd-43d2-9fdb-3a24cfa79e28" -X GET "https://172.25.47.115/ConfigurationManager/v1/objects/host-wwns?porId=CL1-A&hostGroupNumber=23"
  --
  WWN Before:
    "data" : [ ]
  WWN After
    "data" : [ {
    "hostWwnId" : "CL1-A,23,100000109b3dfbbb",
    "portId" : "CL1-A",
    "hostGroupNumber" : 23,
    "hostGroupName" : "TESTING_REST_API",
    "hostWwn" : "100000109b3dfbbb",
    "wwnNickname" : "test-name"
  } ]

*/
