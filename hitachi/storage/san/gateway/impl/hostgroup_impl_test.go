package sanstorage

import (
	"fmt"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/gateway/model"
	"testing"
)

// newHostgroupTestManager is for Testing and provide structure information for connection
func newHostgroupTestManager() (*sanStorageManager, error) {

	objStorage := sanmodel.StorageDeviceSettings{
		Serial:   40014,
		Username: "bXNfdm13YXJl",
		Password: "SGl0YWNoaTE=",
		MgmtIP:   "172.25.47.115",
	}
	psm, err := newSanStorageManagerEx(objStorage)
	if err != nil {
		return nil, fmt.Errorf("unexpected error while creating newHostgroupTestManager %v", err)
	}
	return psm, nil
}

// go test -v -run TestGetHostGroup
func xTestGetHostGroup(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// OUTPUT
		[{
			"hostGroupId" : "CL1-A,0",
			"portId" : "CL1-A",
			"hostGroupNumber" : 0,
			"hostGroupName" : "1A-G00",
			"hostMode" : "LINUX/IRIX"
		  }]
	*/
	resp, err := psm.GetHostGroup("CL1-A", 0)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetHostGroupWwns
func xTestGetHostGroupWwns(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// TODO FIX ME - Need to find hostgroup with proper WWNS data. Right now it is coming blank
		// OUTPUT
		[{
		  "data" : [ ]
		}]
	*/
	resp, err := psm.GetHostGroupWwns("CL1-A", 0)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// go test -v -run TestGetHostGroupLuPaths
func xTestGetHostGroupLuPaths(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
			"data" : [ {

		    "lunId" : "CL1-A,0,0",
		    "portId" : "CL1-A",
		    "hostGroupNumber" : 0,
		    "hostMode" : "LINUX/IRIX",
		    "lun" : 0,
		    "ldevId" : 153,
		    "isCommandDevice" : false,

		    "luHostReserve" : {
		      "openSystem" : false,
		      "persistent" : false,
		      "pgrKey" : false,
		      "mainframe" : false,
		      "acaReserve" : false
		    }

				// OUTPUT - Extract upeer data. Not all comes here
				 [{
				  "data" : [ {
				    "lunId" : "CL1-A,0,0",
				    "portId" : "CL1-A",
				    "hostGroupNumber" : 0,
				    "hostMode" : "LINUX/IRIX",
				    "lun" : 0,
				    "ldevId" : 153,
				    "isCommandDevice" : false,
				  } ]
				}]
	*/
	resp, err := psm.GetHostGroupLuPaths("CL1-A", 0)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}

// https://knowledge.hitachivantara.com/Documents/Storage/VSP_5000_Series/90-01-xx/REST_API_Reference_Guide/05_Volume_allocation
// go test -v -run TestCreateHostGroup
func xTestCreateHostGroup(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	myPortID := "CL1-A"
	myHostGroupName := "TESTING_REST_API"
	//myHostGroupNumber := 0
	myHostModeOptions := []int{12, 13}
	myHostMode := "AIX"
	crReq := sanmodel.CreateHostGroupReqGwy{
		PortID:        &myPortID,
		HostGroupName: &myHostGroupName,
		//HostGroupNumber: &myHostGroupNumber,
		HostModeOptions: myHostModeOptions,
		HostMode:        &myHostMode,
	}
	portid, hgnumber, err := psm.CreateHostGroup(crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("portid: %v, hgnumber: %v", portid, hgnumber)
	/* In GET ALL
		 {
	    "hostGroupId" : "CL1-A,24",
	    "portId" : "CL1-A",
	    "hostGroupNumber" : 24,
	    "hostGroupName" : "TESTING_REST_API",
	    "hostMode" : "AIX"
	  },
	*/
}

// go test -v -run TestDeleteGroup
func xTestDeleteGroup(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	HostgroupNumber := 24

	err = psm.DeleteHostGroup(PortId, HostgroupNumber)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully deleted PortId: %s, HostgroupNumber: %d", PortId, HostgroupNumber)
}

// go test -v -run TestDeleteWwn
func xTestDeleteWwn(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	HostgroupNumber := 23
	Wwn := "100000109b3dfbbb"

	err = psm.DeleteWwn(PortId, HostgroupNumber, Wwn)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully deleted wwn: %s on PortId: %s, HostgroupNumber: %d", Wwn, PortId, HostgroupNumber)
}

// go test -v -run TestGetAllHostGroup
func xTestGetAllHostGroup(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/*
		// OUTPUT
		[ {
		 HostGroupID:CL5-G,6
		 PortID:CL5-G
		 HostGroupNumber:6
		 HostGroupName:esxi-204-nvme-5G
		 HostMode:VMWARE_EX
		 HostModeOptions:[]
		 }
		 {
		 HostGroupID:CL5-G,7
		 PortID:CL5-G
		 HostGroupNumber:7
		 HostGroupName:ksh-test-hg7
		 HostMode:VMWARE
		 HostModeOptions:[]
		 }
		]
	*/
	resp, err := psm.GetAllHostGroups()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
}

// go test -v -run TestGetHostGroupModeAndOptions
func xTestGetHostGroupModeAndOptions(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	resp, err := psm.GetHostGroupModeAndOptions()
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Response: %+v", resp)
	/*
			{
		  "hostModes" : [ {
		    "hostModeId" : 0,
		    "hostModeName" : "Standard",
		    "hostModeDisplay" : "LINUX/IRIX"
		  }, {
		    "hostModeId" : 1,
		    "hostModeName" : "(Deprecated) VMware",
		    "hostModeDisplay" : "VMWARE"
		  }, {
		    "hostModeId" : 3,
		    "hostModeName" : "HP",
		    "hostModeDisplay" : "HP-UX"
		  }, {
		    "hostModeId" : 5,
		    "hostModeName" : "OpenVMS",
		    "hostModeDisplay" : "OVMS"
		  }, {
		    "hostModeId" : 7,
		    "hostModeName" : "Tru64",
		    "hostModeDisplay" : "TRU64"
		  }, {
		    "hostModeId" : 9,
		    "hostModeName" : "Solaris",
		    "hostModeDisplay" : "SOLARIS"
		  }, {
		    "hostModeId" : 10,
		    "hostModeName" : "NetWare",
		    "hostModeDisplay" : "NETWARE"
		  }, {
		    "hostModeId" : 12,
		    "hostModeName" : "(Deprecated) Windows",
		    "hostModeDisplay" : "WIN"
		  }, {
		    "hostModeId" : 15,
		    "hostModeName" : "AIX",
		    "hostModeDisplay" : "AIX"
		  }, {
		    "hostModeId" : 33,
		    "hostModeName" : "VMware Extension",
		    "hostModeDisplay" : "VMWARE_EX"
		  }, {
		    "hostModeId" : 44,
		    "hostModeName" : "Windows Extension",
		    "hostModeDisplay" : "WIN_EX"
		  } ],

		  "hostModeOptions" : [ {
		    "hostModeOptionId" : 2,
		    "hostModeOptionDescription" : "VERITAS Database Edition/Advanced Cluster"
		  }, {
		    "hostModeOptionId" : 6,
		    "hostModeOptionDescription" : "TPRLO"
		  }, {
		    "hostModeOptionId" : 7,
		    "hostModeOptionDescription" : "Automatic recognition function of LUN"
		  }, {
		    "hostModeOptionId" : 12,
		    "hostModeOptionDescription" : "No display for ghost LUN"
		  }, {
		    "hostModeOptionId" : 13,
		    "hostModeOptionDescription" : "SIM report at link failure"
		  }, {
		    "hostModeOptionId" : 14,
		    "hostModeOptionDescription" : "HP TruCluster with TrueCopy function"
		  }, {
		    "hostModeOptionId" : 15,
		    "hostModeOptionDescription" : "HACMP"
		  }, {
		    "hostModeOptionId" : 22,
		    "hostModeOptionDescription" : "Veritas Cluster Server"
		  }, {
		    "hostModeOptionId" : 23,
		    "hostModeOptionDescription" : "REC command support"
		  }, {
		    "hostModeOptionId" : 25,
		    "hostModeOptionDescription" : "Support SPC-3 behavior on Persistent Reservation"
		  }, {
		    "hostModeOptionId" : 33,
		    "hostModeOptionDescription" : "Set/Report Device Identifier enable"
		  }, {
		    "hostModeOptionId" : 39,
		    "hostModeOptionDescription" : "Change the nexus specified in the SCSI Target Reset"
		  }, {
		    "hostModeOptionId" : 40,
		    "hostModeOptionDescription" : "V-Vol expansion"
		  }, {
		    "hostModeOptionId" : 43,
		    "hostModeOptionDescription" : "Queue Full Response"
		  }, {
		    "hostModeOptionId" : 51,
		    "hostModeOptionDescription" : "Round Trip Set Up Option"
		  }, {
		    "hostModeOptionId" : 54,
		    "hostModeOptionDescription" : "(VAAI) Support Option for the EXTENDED COPY command"
		  }, {
		    "hostModeOptionId" : 60,
		    "hostModeOptionDescription" : "LUN0 Change Guard"
		  }, {
		    "hostModeOptionId" : 63,
		    "hostModeOptionDescription" : "(VAAI) Support option for vStorage APIs based on T10 standards"
		  }, {
		    "hostModeOptionId" : 68,
		    "hostModeOptionDescription" : "Support Page Reclamation for Linux"
		  }, {
		    "hostModeOptionId" : 71,
		    "hostModeOptionDescription" : "Change the Unit Attention for Blocked Pool-VOLs"
		  }, {
		    "hostModeOptionId" : 73,
		    "hostModeOptionDescription" : "Support Option for WS2012"
		  }, {
		    "hostModeOptionId" : 78,
		    "hostModeOptionDescription" : "The non-preferred path option"
		  }, {
		    "hostModeOptionId" : 80,
		    "hostModeOptionDescription" : "Multi Text OFF Mode"
		  }, {
		    "hostModeOptionId" : 81,
		    "hostModeOptionDescription" : "NOP-In Suppress Mode"
		  }, {
		    "hostModeOptionId" : 82,
		    "hostModeOptionDescription" : "Discovery CHAP Mode"
		  }, {
		    "hostModeOptionId" : 83,
		    "hostModeOptionDescription" : "Report iSCSI Full Portal List Mode"
		  }, {
		    "hostModeOptionId" : 88,
		    "hostModeOptionDescription" : "Port Consolidation"
		  }, {
		    "hostModeOptionId" : 91,
		    "hostModeOptionDescription" : "Disable I/O wait for OpenStack Option"
		  }, {
		    "hostModeOptionId" : 96,
		    "hostModeOptionDescription" : "Change the nexus specified in the SCSI Logical Unit Reset"
		  }, {
		    "hostModeOptionId" : 97,
		    "hostModeOptionDescription" : "Proprietary ANCHOR command support"
		  }, {
		    "hostModeOptionId" : 105,
		    "hostModeOptionDescription" : "Task Set Full response in the event of I/O overload"
		  }, {
		    "hostModeOptionId" : 110,
		    "hostModeOptionDescription" : "ODX support for WS2012"
		  }, {
		    "hostModeOptionId" : 113,
		    "hostModeOptionDescription" : "iSCSI CHAP Authentication Log"
		  }, {
		    "hostModeOptionId" : 114,
		    "hostModeOptionDescription" : "The automatic asynchronous reclamation on ESXi6.5 or later"
		  }, {
		    "hostModeOptionId" : 122,
		    "hostModeOptionDescription" : "Task Set Full response after reach QoS upper limit"
		  }, {
		    "hostModeOptionId" : 124,
		    "hostModeOptionDescription" : "Guaranteed response during controller failure"
		  }, {
		    "hostModeOptionId" : 131,
		    "hostModeOptionDescription" : "WCE bit OFF mode"
		  } ]

	*/
}

// go test -v -run TestSetHostGroupModeAndOptions
func xTestSetHostGroupModeAndOptions(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	HostgroupNumber := 23
	//HostgroupNumber := 24
	myHostModeOptions := []int{12, 13}
	//myHostModeOptions := []int{12, 14}
	myHostMode := "AIX"
	//myHostMode := "LINUX/IRIX"
	// Both Set: { "hostMode": "WIN", "hostModeOptions": [12,33]}
	// RestHostMode: { "hostMode": "WIN", "hostModeOptions": [-1]}
	// Only HostMode: { "hostMode": "AIX"}

	crReq := sanmodel.SetHostModeAndOptions{
		HostMode:        myHostMode,
		HostModeOptions: &myHostModeOptions,
	}

	err = psm.SetHostGroupModeAndOptions(PortId, HostgroupNumber, crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully set HostModeAndOption: %+v, for PortId: %s, HostgroupNumber: %d", crReq, PortId, HostgroupNumber)
	/**
	BEFORE:
		"hostGroupId" : "CL1-A,23",
		"portId" : "CL1-A",
		"hostGroupNumber" : 23,
		"hostGroupName" : "TEST_REST_API_HOST",
		"hostMode" : "AIX",
		"hostModeOptions" : [ 12, 33 ]
	-----------
	AFTER:
	    "hostGroupId" : "CL1-A,23",
		"portId" : "CL1-A",
		"hostGroupNumber" : 23,
		"hostGroupName" : "TEST_REST_API_HOST",
		"hostMode" : "LINUX/IRIX",   /// MODIFIED <====
		"hostModeOptions" : [ 12, 14 ]  // MODIFIED <====
		**/
}

// go test -v -run TestAddLdevToHG
func xTestAddLdevToHG(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	//PortIds := []string{"CL1-A"}
	HostgroupNumber := 23
	LdevId := 0
	LunId := 0

	crReq := sanmodel.AddLdevToHgReqGwy{
		//PortIds:         PortIds,
		PortID:          &PortId,
		HostGroupNumber: &HostgroupNumber,
		LdevID:          &LdevId,
		Lun:             &LunId,
	}

	err = psm.AddLdevToHG(crReq)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully Added HostModeAndOption: %+v, for PortId: %s, HostgroupNumber: %d", crReq, PortId, HostgroupNumber)
}

// go test -v -run TestRemoveLdevFromHG
func xTestRemoveLdevFromHG(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	HostgroupNumber := 0
	lunID := 0

	err = psm.RemoveLdevFromHG(PortId, HostgroupNumber, lunID)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully Removed Ldev from Hostgroup")
}

// go test -v -run TestSetHostWwnNickName
func xTestSetHostWwnNickName(t *testing.T) {
	psm, err := newHostgroupTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	PortId := "CL1-A"
	HostgroupNumber := 23
	hostWwn := "100000109b3dfbbb"
	//wwnNickname := "TEST-RESTAPI"
	//When deleting the nickname from the WWN
	wwnNickname := ""

	err = psm.SetHostWwnNickName(PortId, HostgroupNumber, hostWwn, wwnNickname)
	if err != nil {
		t.Errorf("Unexpected error in Get %v", err)
		return
	}
	t.Logf("Sucessfully Set NickName for wwn")

	/**

		  "data" : [ {
	    "hostWwnId" : "CL1-A,23,100000109b3dfbbb",
	    "portId" : "CL1-A",
	    "hostGroupNumber" : 23,
	    "hostGroupName" : "TESTING_REST_API",
	    "hostWwn" : "100000109b3dfbbb",
	    "wwnNickname" : "TEST-RESTAPI"
	  } ]

	  --
	  After Delete nick name
	  --
	    "data" : [ {
	    "hostWwnId" : "CL1-A,23,100000109b3dfbbb",
	    "portId" : "CL1-A",
	    "hostGroupNumber" : 23,
	    "hostGroupName" : "TESTING_REST_API",
	    "hostWwn" : "100000109b3dfbbb",
	    "wwnNickname" : "-"
	  } ]

		**/
}
