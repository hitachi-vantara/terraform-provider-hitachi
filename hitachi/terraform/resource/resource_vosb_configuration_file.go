package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"path/filepath"
	"strings"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"
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
	d.SetId(finalPath)
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
	return validateConfigurationFileInputsLogic(d)
}

// type minimalDiff interface {
// 	Get(string) interface{}
// }

// testable
func validateConfigurationFileInputsLogic(diff minimalDiff) error {
	downloadOnly := diff.Get("download_existconfig_only").(bool)
	createOnly := diff.Get("create_only").(bool)
	downloadPath := diff.Get("download_path").(string)

	// Cross-field validation logic
	if downloadOnly && createOnly {
		fmt.Printf("`create_only` is ignored when `download_existconfig_only` is true\n")
	}

	// When a download is expected, validate the path
	if downloadOnly || (!downloadOnly && !createOnly) {
		if strings.TrimSpace(downloadPath) == "" {
			return fmt.Errorf("`download_path` must be set when download is performed")
		}

		dir := filepath.Dir(downloadPath)
		if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
			return fmt.Errorf("directory %q does not exist or is not a directory", dir)
		}
	}

	return nil
}
