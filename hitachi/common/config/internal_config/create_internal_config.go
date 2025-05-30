package main

import (
	"fmt"
	"os"
	"terraform-provider-hitachi/hitachi/common/config"
)

func main() {
	// Default path
	configPath := "./.internal_config"

	// If user provides an argument, use it
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	err := config.CreateDefaultConfigFile(configPath)
	if err != nil {
		fmt.Println("Failed to create config:", err)
	} else {
		fmt.Println("Default config file created successfully at", configPath)
	}
}
