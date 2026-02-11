package sanstorage

import (
	"testing"
)

// go test -v -run TestGetStoragePorts
func xTestGetStoragePorts(t *testing.T) {
	psm, err := newIscsiTargetTestManager()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	/* OUTPIUT:
	{
		"portId" : "CL1-B",
		"portType" : "FIBRE",
		"portAttributes" : [ "TAR" ],
		"portSpeed" : "AUT",
		"loopId" : "D9",
		"fabricMode" : true,
		"portConnection" : "PtoP",
		"lunSecuritySetting" : true,
		"wwn" : "50060e8008757e01"
	}, {
		"portId" : "CL1-C",
		"portType" : "ISCSI",
		"portAttributes" : [ "TAR" ],
		"portSpeed" : "10G",
		"loopId" : "00",
		"fabricMode" : false,
		"lunSecuritySetting" : true
	},
	*/

	resp, err := psm.GetStoragePorts([]string{}, "", "")
	if err != nil {
		t.Errorf("Unexpected error in GetStoragePorts %v", err)
		return
	}
	t.Logf("Response: %v", resp)
}
