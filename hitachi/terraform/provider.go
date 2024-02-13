package terraform

import (
	"context"
	"time"

	// "fmt"
	// "strconv"
	// "time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"

	// "github.com/hashicorp/terraform-plugin-log/tflog"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: schemaimpl.ProviderSchema,
		ResourcesMap: map[string]*schema.Resource{

			"hitachi_vsp_volume":                resourceimpl.ResourceStorageLun(),
			"hitachi_vsp_hostgroup":             resourceimpl.ResourceStorageHostGroup(),
			"hitachi_vsp_iscsi_target":          resourceimpl.ResourceStorageIscsiTarget(),
			"hitachi_vsp_iscsi_chap_user":       resourceimpl.ResourceStorageIscsiChapUser(),
			"hitachi_vss_block_compute_node":    resourceimpl.ResourceVssbStorageComputeNode(),
			"hitachi_vss_block_volume":          resourceimpl.ResourceVssbStorageCreateVolume(),
			"hitachi_vss_block_iscsi_chap_user": resourceimpl.ResourceVssbStorageChapUser(),
			"hitachi_vss_block_compute_port":    resourceimpl.ResourceVssbStorageComputePort(),

			//Infra resources

			"hitachi_infra_hostgroup":      resourceimpl.ResourceInfraHostGroup(),
			"hitachi_infra_storage_device": resourceimpl.ResourceInfraStorageDevice(),
			"hitachi_infra_iscsi_target":   resourceimpl.ResourceInfraIscsiTarget(),
			"hitachi_infra_vsp_volume":     resourceimpl.ResourceInfraStorageVOlume(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hitachi_vsp_storage":                datasourceimpl.DataSourceStorageSystem(),
			"hitachi_vsp_volume":                 datasourceimpl.DataSourceStorageLun(),
			"hitachi_vsp_volumes":                datasourceimpl.DataSourceStorageLuns(),
			"hitachi_vsp_hostgroup":              datasourceimpl.DataSourceStorageHostGroup(),
			"hitachi_vsp_hostgroups":             datasourceimpl.DataSourceStorageHostGroups(),
			"hitachi_vsp_iscsi_target":           datasourceimpl.DataSourceStorageIscsiTarget(),
			"hitachi_vsp_iscsi_targets":          datasourceimpl.DataSourceStorageIscsiTargets(),
			"hitachi_vsp_iscsi_chap_user":        datasourceimpl.DataSourceStorageChapUser(),
			"hitachi_vsp_iscsi_chap_users":       datasourceimpl.DataSourceStorageChapUsers(),
			"hitachi_vsp_storage_ports":          datasourceimpl.DataSourceStoragePorts(),
			"hitachi_vsp_dynamic_pools":          datasourceimpl.DataSourceStorageDynamicPools(),
			"hitachi_vsp_parity_groups":          datasourceimpl.DataSourceStorageParityGroups(),
			"hitachi_vss_block_storage_pools":    datasourceimpl.DataSourceStoragePools(),
			"hitachi_vss_block_volumes":          datasourceimpl.DataSourceVssbVolumes(),
			"hitachi_vss_block_compute_nodes":    datasourceimpl.DataSourceVssbComputeNodes(),
			"hitachi_vss_block_volume":           datasourceimpl.DataSourceVssbVolumeNodes(),
			"hitachi_vss_block_storage_ports":    datasourceimpl.DataSourceVssbStoragePorts(),
			"hitachi_vss_block_iscsi_chap_users": datasourceimpl.DataSourceVssbChapUsers(),
			"hitachi_vss_block_iscsi_port_auth":  datasourceimpl.DataSourceVssbComputePort(),
			"hitachi_vss_block_dashboard":        datasourceimpl.DataSourceVssbDashboard(),

			//Infra data-resources
			"hitachi_infra_storage_ports":   datasourceimpl.DataSourceInfraStoragePorts(),
			"hitachi_infra_parity_groups":   datasourceimpl.DataSourceInfraParityGroups(),
			"hitachi_infra_hostgroup":       datasourceimpl.DataSourceInfraHostGroup(),
			"hitachi_infra_hostgroups":      datasourceimpl.DataSourceInfraHostGroups(),
			"hitachi_infra_storage_devices": datasourceimpl.DataSourceInfraStorageDevices(),
			"hitachi_infra_storage_pools":   datasourceimpl.DataSourceInfraStoragePools(),
			"hitachi_infra_iscsi_target":    datasourceimpl.DataSourceInfraIscsiTarget(),
			"hitachi_infra_iscsi_targets":   datasourceimpl.DataSourceInfraIscsiTargets(),
			"hitachi_infra_chap_users":      datasourceimpl.DataSourceInfraChapUsers(),
			"hitachi_infra_volumes":         datasourceimpl.DataSourceInfraVolumes(),
			"hitachi_infra_volume":          datasourceimpl.DataSourceInfraVolume(),
			"hitachi_infra_ucp_systems":     datasourceimpl.DataSourceInfraUcpSystems(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	// var diags diag.Diagnostics
	// tflog.Info(ctx, "THIS IS JUST TESTING TFLOG")
	// // example to append to diags
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Error,
	// 		Summary:  "Unable to create HashiCups client",
	// 		Detail:   "Unable to create anonymous HashiCups client",
	// 	})

	// ============
	// check storage with creds then saves the input given and minimal storage info in a file.

	sanList := []map[string]interface{}{}

	// Uncomment following line if you want to debug Terraform, also update processId in launch.json file
	time.Sleep(15 * time.Second)

	ssarray, err := impl.RegisterStorageSystem(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	for _, pss := range ssarray.VspStorageSystem {
		ss := *pss
		san := map[string]interface{}{
			"storage_device_id":      ss.StorageDeviceID,
			"storage_serial_number":  ss.SerialNumber,
			"storage_device_model":   ss.Model,
			"dkc_micro_code_version": ss.MicroVersion,
			"management_ip":          ss.MgmtIP,
			"svp_ip":                 ss.SvpIP,
			"controller1_ip":         ss.ControllerIP1,
			"controller2_ip":         ss.ControllerIP2,
		}

		log.WriteDebug("san: %+v\n", san)
		sanList = append(sanList, san)
	}

	for _, pss := range ssarray.VssbStorageVersionInfo {
		ss := *pss
		vssb := map[string]interface{}{
			"vssb_storage_api_version":  ss.ApiVersion,
			"vssb_storage_product_name": ss.ProductName,
		}

		log.WriteDebug("vssb: %+v\n", vssb)
		sanList = append(sanList, vssb)
	}

	for _, pss := range ssarray.InfraGwInfo {
		ss := *pss
		infra_gw := map[string]interface{}{
			"infra_gw_address": ss.Address,
		}

		log.WriteDebug("infra_gw: %+v\n", infra_gw)
		sanList = append(sanList, infra_gw)
	}
	return ssarray, nil
}
