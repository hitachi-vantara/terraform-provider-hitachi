package admin

import (
	"encoding/json"
	"testing"

	gwymodel "terraform-provider-hitachi/hitachi/storage/admin/gateway/model"
)

// Helper functions for pointers
func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

// go test -v -run TestReconcileAddHostGroupsToServer
func xTestReconcileAddHostGroupsToServer(t *testing.T) {
	psm, err := newTestManager()
	if err != nil {
		t.Fatalf("Unexpected error creating adminStorageManager: %v", err)
	}

	serverId := 1234 // Replace with mock or valid test server ID

	tests := []struct {
		name   string
		params gwymodel.AddHostGroupsToServerParam
	}{
		{
			name: "No HostGroups - Should Skip",
			params: gwymodel.AddHostGroupsToServerParam{
				HostGroups: []gwymodel.HostGroupForAddToServerParam{},
			},
		},
		{
			name: "Single HostGroup Add and Sync",
			params: gwymodel.AddHostGroupsToServerParam{
				HostGroups: []gwymodel.HostGroupForAddToServerParam{
					{
						PortID:        "CL1-A",
						HostGroupID:   intPtr(100),
						HostGroupName: nil,
					},
				},
			},
		},
		{
			name: "Multiple HostGroups Add and Sync",
			params: gwymodel.AddHostGroupsToServerParam{
				HostGroups: []gwymodel.HostGroupForAddToServerParam{
					{
						PortID:        "CL1-A",
						HostGroupID:   intPtr(100),
						HostGroupName: strPtr("HG-CL1-A"),
					},
					{
						PortID:        "CL3-A",
						HostGroupID:   intPtr(101),
						HostGroupName: strPtr("HG-CL3-A"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paramsJSON, _ := json.MarshalIndent(tt.params, "", "  ")
			t.Logf("Params for %s:\n%s", tt.name, string(paramsJSON))

			err := psm.AddHostGroupsToServer(serverId, tt.params)
			if err != nil {
				t.Errorf("Unexpected error in %s: %v", tt.name, err)
				return
			}

			t.Logf("%s succeeded for ServerID %d", tt.name, serverId)
		})
	}
}
