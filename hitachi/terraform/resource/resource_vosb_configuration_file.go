package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncRestoreConfigFileOperation = &sync.Mutex{}

func ResourceVssbConfigurationFile() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block: Create and/or download configuration definition file of the storage system.",
		CreateContext: resourceVssbConfigurationFileCreate,
		ReadContext:   resourceVssbConfigurationFileRead,
		UpdateContext: resourceVssbConfigurationFileUpdate,
		DeleteContext: resourceVssbConfigurationFileDelete,
		Schema:        schemaimpl.ResourceVssbConfigurationFileSchema,
		CustomizeDiff: validateConfigurationFileInputs,
	}
}

func resourceVssbConfigurationFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncRestoreConfigFileOperation.Lock()
	defer syncRestoreConfigFileOperation.Unlock()

	log.WriteInfo("starting restore config definition file")

	finalPath, err := impl.CreateDownloadConfigurationDefinitionFile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to restore configuration definition file: %w", err))
	}

	// Just set a random ID since there's no persisted object.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	d.Set("output_file_path", finalPath)
	d.Set("status", "Configuration file operation successful")
	log.WriteInfo("Configuration file operation successful")
	log.WriteInfo(fmt.Sprintf("Output file: %s", finalPath))

	return nil
}

func resourceVssbConfigurationFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Same logic as Create since restoring is a one-off action.
	return resourceVssbConfigurationFileCreate(ctx, d, m)
}

func resourceVssbConfigurationFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Nothing to delete; this is a one-time operation.
	d.SetId("")
	return nil
}

func resourceVssbConfigurationFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No state to read from backend.
	return nil
}

func validateConfigurationFileInputs(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return ValidateConfigurationFileInputsLogic(d)
}

// type minimalDiff is defined in resource_vosb_storage_credential.go

// testable
func ValidateConfigurationFileInputsLogic(d minimalDiff) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	downloadOnly := d.Get("download_existconfig_only").(bool)
	createOnly := d.Get("create_only").(bool)
	downloadPath := d.Get("download_path").(string)

	// Cross-field validation logic
	if downloadOnly && createOnly {
		log.WriteDebug("`create_only` is ignored when `download_existconfig_only` is true\n")
	}

	// When a download is expected, validate the path
	if downloadOnly || (!downloadOnly && !createOnly) {
		if strings.TrimSpace(downloadPath) == "" {
			return fmt.Errorf("`download_path` must be set when download is performed")
		}

		dir := filepath.Dir(downloadPath)
		if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
			return fmt.Errorf("parent directory of download_path %q does not exist", dir)
		}
	}

	// When a create is expected, validate create parameters
	if createOnly || (!downloadOnly && !createOnly) {
		if err := validateCreateParams(d); err != nil {
			return err
		}
	} else {
		log.WriteDebug("No create operation specified; skipping other parameters\n")
	}

	return nil
}

func validateCreateParams(d minimalDiff) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	paramList, ok := d.Get("create_configuration_file_param").([]interface{})
	if !ok || len(paramList) == 0 {
		log.WriteDebug("create_configuration_file_param not set, assuming baremetal or AWS")
		return nil
	}

	firstParam, ok := paramList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("create_configuration_file_param[0] must be a map")
	}

	expectedCloudProvider, ok := firstParam["expected_cloud_provider"].(string)
	if !ok || strings.TrimSpace(expectedCloudProvider) == "" {
		return fmt.Errorf("expected_cloud_provider must be set in create_configuration_file_param")
	}
	expectedCloudProvider = strings.ToLower(expectedCloudProvider)

	if expectedCloudProvider == "aws" || expectedCloudProvider == "baremetal" {
		log.WriteDebug("skipping create_configuration_file_param for %s", expectedCloudProvider)
		return nil
	}

	// optional, default to "Normal"
	exportType, _ := firstParam["export_file_type"].(string)

	// General validation rules per type
	switch exportType {
	case "Normal":
		// Normal type does not require any specific parameters
		log.WriteDebug("ignoring other parameters in create_configuration_file_param for export_type %s", "Normal")
	case "AddStorageNodes":
		if firstParam["machine_image_id"] == nil || firstParam["machine_image_id"] == "" {
			return fmt.Errorf("machine_image_id must be set for AddStorageNodes")
		}
		rawList, ok := firstParam["address_setting"]
		if !ok {
			return fmt.Errorf("address_setting must be set and contain at least one item for AddStorageNodes")
		}
		list, ok := rawList.([]interface{})
		if !ok || len(list) == 0 {
			return fmt.Errorf("address_setting must be a non-empty list")
		}

		indexSet := make(map[int]int)    // index -> entry#
		ipSet := make(map[string]string) // IP -> location

		for i, raw := range list {
			item, ok := raw.(map[string]interface{})
			if !ok {
				return fmt.Errorf("address_setting[%d] must be a map", i)
			}

			idx, ok := item["index"].(int)
			if !ok {
				return fmt.Errorf("address_setting[%d] missing or invalid `index`", i)
			}
			if prev, found := indexSet[idx]; found {
				return fmt.Errorf("duplicate index %d in `address_setting` (items #%d and #%d)", idx, prev+1, i+1)
			}
			indexSet[idx] = i

			// Check that all three required IP fields are present and are non-empty strings
			var (
				c, cOK   = item["control_port_ipv4_address"].(string)
				iip, iOK = item["internode_port_ipv4_address"].(string)
				cp, cpOK = item["compute_port_ipv4_address"].(string)
			)
			if !cOK || c == "" {
				return fmt.Errorf("address_setting[%d] missing or invalid `control_port_ipv4_address`", i)
			}
			if !iOK || iip == "" {
				return fmt.Errorf("address_setting[%d] missing or invalid `internode_port_ipv4_address`", i)
			}
			if !cpOK || cp == "" {
				return fmt.Errorf("address_setting[%d] missing or invalid `compute_port_ipv4_address`", i)
			}

			// optional compute_port_ipv6_address for azure only

			// Internal uniqueness within one entry
			if c == iip || c == cp || iip == cp {
				return fmt.Errorf("address_setting[%d] has duplicate IPs among control, internode, and compute fields", i)
			}

			// Global uniqueness across all entries
			for field, ip := range map[string]string{
				"control_port_ipv4_address":   c,
				"internode_port_ipv4_address": iip,
				"compute_port_ipv4_address":   cp,
			} {
				if prev, exists := ipSet[ip]; exists {
					return fmt.Errorf("duplicate IP `%s` for field `%s` in address_setting[%d]; already used in %s", ip, field, i, prev)
				}
				ipSet[ip] = fmt.Sprintf("address_setting[%d]/%s", i, field)
			}
		}

	case "ReplaceStorageNode":
		if firstParam["machine_image_id"] == nil || firstParam["machine_image_id"] == "" {
			return fmt.Errorf("machine_image_id must be set for ReplaceStorageNode")
		}
		if expectedCloudProvider == "google" {
			if firstParam["node_id"] == nil || firstParam["node_id"] == "" {
				return fmt.Errorf("node_id must be set for ReplaceStorageNode on Google")
			}
			// recover_single_node (bool) is optional — no validation
		}

	case "ReplaceDrive":
		if expectedCloudProvider != "google" {
			return fmt.Errorf("ReplaceDrive is only supported on Google Cloud")
		}
		recover, _ := firstParam["recover_single_drive"].(bool)
		driveID, _ := firstParam["drive_id"].(string)

		if recover {
			if driveID != "" {
				return fmt.Errorf("drive_id must not be set when recover_single_drive is true")
			}
		} else {
			if driveID == "" {
				return fmt.Errorf("drive_id must be set when recover_single_drive is false")
			}
		}

	case "AddDrives":
		if firstParam["number_of_drives"] == nil || firstParam["number_of_drives"] == 0 {
			return fmt.Errorf("number_of_drives must be set for AddDrives")
		}

	default:
		return fmt.Errorf("unsupported export_file_type %q for expected_cloud_provider %q", exportType, expectedCloudProvider)
	}

	// This check is already done in the schema, but we do it here too
	// Warn about unknown fields
	allowed := map[string]struct{}{
		"export_file_type":        {},
		"expected_cloud_provider": {},
		"machine_image_id":        {},
		"recover_single_node":     {},
		"node_id":                 {},
		"recover_single_drive":    {},
		"drive_id":                {},
		"number_of_drives":        {},
		"address_setting":         {},
	}
	var extras []string
	for k := range firstParam {
		if _, ok := allowed[k]; !ok {
			extras = append(extras, k)
		}
	}
	if len(extras) > 0 {
		log.WriteWarn(fmt.Sprintf("unexpected parameter(s) ignored for export_file_type %s: %v", exportType, extras))
	}

	return nil
}
