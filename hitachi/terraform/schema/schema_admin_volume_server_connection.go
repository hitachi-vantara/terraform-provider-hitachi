package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ------------------- Volume-Server Connections Schema -------------------

func volumeServerConnectionListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connections_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of volume-server connection details.",
			Elem: &schema.Resource{
				Schema: volumeServerConnectionSchema(),
			},
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of entries returned.",
		},
		"total_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of entries available.",
		},
	}
}

// ------------------- Volume-Server Connection Schema -------------------

func volumeServerConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Composite identifier of the connection in the format '{volumeId},{serverId}'.",
		},
		"volume_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Volume ID of the connection.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Server ID of the connection.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"luns": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "List of LUN and port information associated with this volume-server connection.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"lun": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Logical Unit Number (LUN).",
					},
					"port_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Port ID associated with this LUN.",
					},
				},
			},
		},
	}
}

// ------------------- Datasource Multiple Volume-Server Connections Schema -------------------

func datasourceAdminMultipleVolumeServerConnectionsInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Server ID connected to the volume. Either `server_id` or `server_nickname` must be specified. If both are specified, an error occurs.",
			ValidateFunc: validation.IntAtLeast(0),
			ConflictsWith: []string{
				"server_nickname",
			},
		},
		"server_nickname": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Server nickname connected to the volume. Either `server_id` or `server_nickname` must be specified. If both are specified, an error occurs.",
			ValidateFunc: validation.StringLenBetween(1, 64),
			ConflictsWith: []string{
				"server_id",
			},
		},
		"start_volume_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			Description:  "Starting VOLUME ID to display. Default is 0.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"requested_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      2048,
			Description:  "Number of connection information entries to display. Default is 2048. Must be between 1 and 4096.",
			ValidateFunc: validation.IntBetween(1, 4096),
		},
	}
}

func datasourceAdminMultipleVolumeServerConnectionsOutputSchema() map[string]*schema.Schema {
	return volumeServerConnectionListSchema()
}

func DatasourceAdminMultipleVolumeServerConnectionsSchema() map[string]*schema.Schema {
	schema := datasourceAdminMultipleVolumeServerConnectionsInputSchema()
	for k, v := range datasourceAdminMultipleVolumeServerConnectionsOutputSchema() {
		schema[k] = v
	}
	return schema
}

// ------------------- Datasource One Volume-Server Connection Schema -------------------

func datasourceAdminOneVolumeServerConnectionInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"volume_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "ID of the volume.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"server_id": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "ID of the server.",
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}

func datasourceAdminOneVolumeServerConnectionOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Volume-server connection info.",
			Elem: &schema.Resource{
				Schema: volumeServerConnectionSchema(),
			},
		},
	}
}

func DatasourceAdminOneVolumeServerConnectionSchema() map[string]*schema.Schema {
	schema := datasourceAdminOneVolumeServerConnectionInputSchema()
	for k, v := range datasourceAdminOneVolumeServerConnectionOutputSchema() {
		schema[k] = v
	}
	return schema
}

// ------------------- Resource Volume-Server Connection Schema -------------------

func resourceAdminVolumeServerConnectionInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:         schema.TypeInt,
			Required:     true,
			Description:  "Serial number of the storage system.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"volume_ids": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "List of volume IDs to attach to servers.",
			Elem: &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
		"server_ids": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "List of server IDs to connect the specified volumes to.",
			Elem: &schema.Schema{
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceAdminVolumeServerConnectionOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connections_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Volume-server connections info.",
			Elem: &schema.Resource{
				Schema: volumeServerConnectionSchema(),
			},
		},
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Colon-separated list of volume-server connection IDs.",
		},
	}
}

func ResourceAdminVolumeServerConnectionSchema() map[string]*schema.Schema {
	schema := resourceAdminVolumeServerConnectionInputSchema()
	for k, v := range resourceAdminVolumeServerConnectionOutputSchema() {
		schema[k] = v
	}
	return schema
}
