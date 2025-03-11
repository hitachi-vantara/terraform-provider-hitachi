package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-hitachi/hitachi/terraform"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {

	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return terraform.Provider()
		},
	}

	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		userPluginDir := "/root/.terraform.d/plugins/localhost/hitachi-vantara/hitachi/2.0/linux_amd64/terraform-provider-hitachi"
		// userPluginDir := "/root/.terraform.d/plugins/20-95.sie.hds.com/hv/hitachi/1.0/linux_amd64/terraform-provider-hitachi"
		err := plugin.Debug(context.Background(), userPluginDir, opts)

		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
