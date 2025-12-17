package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-hitachi/hitachi/terraform"
)

var version = "2.2.0" // Will be overridden with -ldflags during build

// main.go is only used for plugin startup via plugin.Serve(...).
// It runs once when Terraform loads the provider binary, and it's not meant for handling provider-specific configuration logic.
func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Hitachi Terraform Provider version: %s\n", version)
		os.Exit(0)
	}

	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "Enable plugin debug mode (for use with delve)")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return terraform.Provider()
		},
		Debug: debugMode,
	})
}
