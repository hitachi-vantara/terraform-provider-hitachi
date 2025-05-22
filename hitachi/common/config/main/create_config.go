package main

import (
	"fmt"
	"terraform-provider-hitachi/hitachi/common/config"
)

func main() {
	configPath := "./config.json"
	err := config.CreateDefaultConfigFile(configPath)
	if err != nil {
		fmt.Println("Failed to create config:", err)
	} else {
		fmt.Println("Default config file created successfully at", configPath)
	}
}
