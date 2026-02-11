package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ------------------- Server HBA Info Schema -------------------
func serverHbaInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Server ID",
		},
		"hba_wwn": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HBA World Wide Name (may not be present in all responses)",
		},
		"iscsi_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "iSCSI initiator name",
		},
		"port_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of port IDs associated with the HBA",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// ------------------- Server HBAs Info List Schema -------------------
func serverHbasInfoListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_hba_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of server HBAs",
		},
		"server_hba_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of server HBA information",
			Elem: &schema.Resource{
				Schema: serverHbaInfoSchema(),
			},
		},
	}
}

// ------------------- Datasource Get Multiple Server HBAs Schema -------------------
func datasourceAdminServerHbasInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of storage",
		},
		"server_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Server ID",
		},
	}
}

func datasourceAdminServerHbasOutputSchema() map[string]*schema.Schema {
	return serverHbasInfoListSchema()
}

func DatasourceAdminServerHbasSchema() map[string]*schema.Schema {
	schema := datasourceAdminServerHbasInputSchema()

	for k, v := range datasourceAdminServerHbasOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Datasource Get One Server HBA Schema -------------------
func datasourceAdminServerHbaInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Serial number of storage",
		},
		"server_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Server ID",
		},
		"initiator_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "HBA World Wide Name or iSCSI Initiator Name",
		},
	}
}

func datasourceAdminServerHbaOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_hba_info": { // returns only one HBA inside a list
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Server HBA Info returned from API",
			Elem: &schema.Resource{
				Schema: serverHbaInfoSchema(),
			},
		},
	}
}

func DatasourceAdminServerHbaSchema() map[string]*schema.Schema {
	schema := datasourceAdminServerHbaInputSchema()

	for k, v := range datasourceAdminServerHbaOutputSchema() {
		schema[k] = v
	}

	return schema
}

// ------------------- Resource Server HBA Schema -------------------

func resourceAdminServerHbaInputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"serial": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Serial number of storage",
		},
		"server_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "Server ID to add HBA to",
		},
		"hbas": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			Description: "List of HBAs to add to the server",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hba_wwn": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "HBA World Wide Name",
					},
					"iscsi_name": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "iSCSI initiator name",
					},
				},
			},
		},
	}
}

func resourceAdminServerHbaOutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_hba_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of server HBA information returned from API",
			Elem: &schema.Resource{
				Schema: serverHbaInfoSchema(),
			},
		},
		"server_hba_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Total number of server HBAs",
		},
	}
}

func ResourceAdminServerHbaSchema() map[string]*schema.Schema {
	schema := resourceAdminServerHbaInputSchema()

	for k, v := range resourceAdminServerHbaOutputSchema() {
		schema[k] = v
	}

	return schema
}
