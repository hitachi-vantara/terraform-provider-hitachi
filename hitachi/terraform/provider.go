package terraform

import (
	"context"
	"fmt"
	"sync"
	config "terraform-provider-hitachi/hitachi/common/config"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	telemetry "terraform-provider-hitachi/hitachi/common/telemetry"
	datasourceimpl "terraform-provider-hitachi/hitachi/terraform/datasource"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	resourceimpl "terraform-provider-hitachi/hitachi/terraform/resource"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	configOnce sync.Once
	configErr  error

	userConsentOnce    sync.Once
	isUserConsentExist bool
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: schemaimpl.ProviderSchema,
		ResourcesMap: map[string]*schema.Resource{
			"hitachi_vsp_storage_maintenance": resourceimpl.ResourceStorageMaintenance(),
			"hitachi_vsp_volume":              resourceimpl.ResourceStorageLun(),
			"hitachi_vsp_hostgroup":           resourceimpl.ResourceStorageHostGroup(),
			"hitachi_vsp_iscsi_target":        resourceimpl.ResourceStorageIscsiTarget(),
			"hitachi_vsp_iscsi_chap_user":     resourceimpl.ResourceStorageIscsiChapUser(),
			// "hitachi_vsp_snapshot":                     resourceimpl.ResourceVspSnapshot(),
			// "hitachi_vsp_snapshot_group":               resourceimpl.ResourceVspSnapshotGroup(),
			"hitachi_vsp_pav_ldev":                     resourceimpl.ResourceVspPavLdev(),
			"hitachi_vosb_compute_node":                resourceimpl.ResourceVssbStorageComputeNode(),
			"hitachi_vosb_storage_node":                resourceimpl.ResourceVssbStorageNode(),
			"hitachi_vosb_volume":                      resourceimpl.ResourceVssbStorageCreateVolume(),
			"hitachi_vosb_iscsi_chap_user":             resourceimpl.ResourceVssbStorageChapUser(),
			"hitachi_vosb_compute_port":                resourceimpl.ResourceVssbStorageComputePort(),
			"hitachi_vosb_change_user_password":        resourceimpl.ResourceVssbChangeUserPassword(),
			"hitachi_vosb_add_drives_to_pool":          resourceimpl.ResourceVssbAddDrivesToPool(),
			"hitachi_vosb_configuration_file":          resourceimpl.ResourceVssbConfigurationFile(),
			"hitachi_vsp_one_volume":                   resourceimpl.ResourceAdminVolume(),
			"hitachi_vsp_one_volume_qos":               resourceimpl.ResourceAdminVolumeQos(),
			"hitachi_vsp_one_iscsi_target":             resourceimpl.ResourceAdminIscsiTarget(),
			"hitachi_vsp_one_server":                   resourceimpl.ResourceAdminServer(),
			"hitachi_vsp_one_port":                     resourceimpl.ResourceAdminPort(),
			"hitachi_vsp_one_server_path":              resourceimpl.ResourceAdminServerPath(),
			"hitachi_vsp_one_server_hba":               resourceimpl.ResourceAdminServerHba(),
			"hitachi_vsp_one_pool":                     resourceimpl.ResourceAdminPool(),
			"hitachi_vsp_one_volume_server_connection": resourceimpl.ResourceAdminVolumeServerConnection(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"hitachi_vsp_storage":              datasourceimpl.DataSourceStorageSystem(),
			"hitachi_vsp_volume":               datasourceimpl.DataSourceStorageLun(),
			"hitachi_vsp_volumes":              datasourceimpl.DataSourceStorageLuns(),
			"hitachi_vsp_hostgroup":            datasourceimpl.DataSourceStorageHostGroup(),
			"hitachi_vsp_hostgroups":           datasourceimpl.DataSourceStorageHostGroups(),
			"hitachi_vsp_iscsi_target":         datasourceimpl.DataSourceStorageIscsiTarget(),
			"hitachi_vsp_iscsi_targets":        datasourceimpl.DataSourceStorageIscsiTargets(),
			"hitachi_vsp_iscsi_chap_user":      datasourceimpl.DataSourceStorageChapUser(),
			"hitachi_vsp_iscsi_chap_users":     datasourceimpl.DataSourceStorageChapUsers(),
			"hitachi_vsp_storage_ports":        datasourceimpl.DataSourceStoragePorts(),
			"hitachi_vsp_dynamic_pool":         datasourceimpl.DataSourceStorageDynamicPool(),
			"hitachi_vsp_dynamic_pools":        datasourceimpl.DataSourceStorageDynamicPools(),
			"hitachi_vsp_pav_alias":            datasourceimpl.DataSourceStoragePavAlias(),
			"hitachi_vsp_supported_host_modes": datasourceimpl.DataSourceStorageSupportedHostModes(),
			"hitachi_vsp_parity_groups":        datasourceimpl.DataSourceStorageParityGroups(),
			"hitachi_vsp_parity_group":         datasourceimpl.DataSourceStorageParityGroup(),
			// "hitachi_vsp_snapshot":                      datasourceimpl.DatasourceVspSnapshot(),
			// "hitachi_vsp_snapshots":                     datasourceimpl.DatasourceVspSnapshotRange(),
			// "hitachi_vsp_snapshot_group":                datasourceimpl.DatasourceVspSnapshotGroup(),
			// "hitachi_vsp_snapshot_groups":               datasourceimpl.DatasourceVspMultipleSnapshotGroups(),
			// "hitachi_vsp_snapshot_family":               datasourceimpl.DatasourceVspSnapshotFamily(),
			// "hitachi_vsp_vclone_parent_vols":            datasourceimpl.DatasourceVspVirtualCloneParentVolume(),
			"hitachi_vosb_storage_pools":                datasourceimpl.DataSourceVssbStoragePools(),
			"hitachi_vosb_volumes":                      datasourceimpl.DataSourceVssbVolumes(),
			"hitachi_vosb_compute_nodes":                datasourceimpl.DataSourceVssbComputeNodes(),
			"hitachi_vosb_volume":                       datasourceimpl.DataSourceVssbVolumeNodes(),
			"hitachi_vosb_storage_nodes":                datasourceimpl.DataSourceVssbStorageNodes(),
			"hitachi_vosb_storage_ports":                datasourceimpl.DataSourceVssbStoragePorts(),
			"hitachi_vosb_iscsi_chap_users":             datasourceimpl.DataSourceVssbChapUsers(),
			"hitachi_vosb_iscsi_port_auth":              datasourceimpl.DataSourceVssbComputePort(),
			"hitachi_vosb_dashboard":                    datasourceimpl.DataSourceVssbDashboard(),
			"hitachi_vosb_storage_drives":               datasourceimpl.DataSourceVssbStorageDrives(),
			"hitachi_vsp_one_storage":                   datasourceimpl.DataSourceStorageSystemAdmin(),
			"hitachi_vsp_one_volume_qos":                datasourceimpl.DataSourceStorageVolumeQos(),
			"hitachi_vsp_one_volume":                    datasourceimpl.DatasourceAdminVolume(),
			"hitachi_vsp_one_volumes":                   datasourceimpl.DatasourceAdminVolumes(),
			"hitachi_vsp_one_iscsi_target":              datasourceimpl.DatasourceAdminIscsiTarget(),
			"hitachi_vsp_one_iscsi_targets":             datasourceimpl.DatasourceAdminIscsiTargets(),
			"hitachi_vsp_one_servers":                   datasourceimpl.DataSourceAdminServerList(),
			"hitachi_vsp_one_server":                    datasourceimpl.DataSourceAdminServerInfo(),
			"hitachi_vsp_one_ports":                     datasourceimpl.DatasourceAdminPorts(),
			"hitachi_vsp_one_port":                      datasourceimpl.DatasourceAdminPort(),
			"hitachi_vsp_one_server_hbas":               datasourceimpl.DatasourceAdminServerHbas(),
			"hitachi_vsp_one_server_hba":                datasourceimpl.DatasourceAdminServerHba(),
			"hitachi_vsp_one_server_path":               datasourceimpl.DatasourceAdminServerPath(),
			"hitachi_vsp_one_volume_server_connection":  datasourceimpl.DatasourceAdminVolumeServerConnection(),
			"hitachi_vsp_one_volume_server_connections": datasourceimpl.DatasourceAdminVolumeServerConnections(),
			"hitachi_vsp_one_pools":                     datasourceimpl.DatasourceAdminPools(),
			"hitachi_vsp_one_pool":                      datasourceimpl.DatasourceAdminPool(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure sets up the provider's configuration.
// providerConfigure is executed each time Terraform runs plan, apply, etc., and it's the canonical place to:
// Load provider settings,
// Initialize clients,
// Validate environment or file presence (like consent file),
// Return diagnostics if required files are missing or invalid.
// Recommended:
// This function must not make any backend API calls or perform side effects.
// It should only validate and return the configured client or settings.

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	diags := diag.Diagnostics{}

	configOnce.Do(func() {
		configErr = config.Load(config.CONFIG_FILE) // data saved in config.ConfigData global var
	})

	if configErr != nil {
		log.WriteInfo("Could not load %s. A default config may have been created. Details: %v", config.CONFIG_FILE, configErr)
		return nil, diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Default config file created",
				Detail:   fmt.Sprintf("Could not read %s — a new default config was created. Details: %v", config.CONFIG_FILE, configErr),
			},
		}
	}

	consentMessage(&diags)

	// check storage with creds then saves the input given and minimal storage info in a file.

	sanList := []map[string]interface{}{}

	// Uncomment following line if you want to debug Terraform, also update processId in launch.json file
	// time.Sleep(15 * time.Second)

	ssarray, err := impl.RegisterStorageSystem(d)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return nil, diags
	}

	for _, pss := range ssarray.VspStorageSystem {
		ss := *pss
		san := map[string]interface{}{
			"storage_device_id":      ss.StorageDeviceID,
			"storage_serial_number":  ss.SerialNumber,
			"storage_device_model":   ss.Model,
			"dkc_micro_code_version": ss.MicroVersion,
			"management_ip":          ss.MgmtIP,
			"ip":                     ss.IP,
			"controller1_ip":         ss.ControllerIP1,
			"controller2_ip":         ss.ControllerIP2,
		}
		log.WriteDebug("san: %+v\n", san)
		sanList = append(sanList, san)
	}

	for _, pss := range ssarray.VssbStorageVersionInfo {
		ss := *pss
		vssb := map[string]interface{}{
			"vosb_storage_api_version":  ss.ApiVersion,
			"vosb_storage_product_name": ss.ProductName,
		}
		log.WriteDebug("vosb: %+v\n", vssb)
		sanList = append(sanList, vssb)
	}

	for _, pss := range ssarray.AdminStorageSystem {
		ss := *pss
		log.WriteDebug("vsp_one: %+v\n", ss)
		log.WriteDebug("MGMT IP: %+v\n", ss.Serial)

	}

	return ssarray, diags
}

func consentMessage(diags *diag.Diagnostics) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	userConsentOnce.Do(func() {
		isUserConsentExist = telemetry.IsUserConsentExist()
	})

	if !isUserConsentExist {
		log.WriteInfo("User has not run bin/user_consent.sh")

		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "User consent is requested for telemetry data collection. This is optional.",
			Detail:   config.ConfigData.UserConsentMessage + config.ConfigData.RunConsentMessage,
		})
	}
}
