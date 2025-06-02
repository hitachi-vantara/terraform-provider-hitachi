package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

var (
	ConfigData *Config
	configOnce sync.Once
)

// Load loads config from file, or creates a default config if missing.
func Load(path string) error {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var err error

	configOnce.Do(func() {
		log.WriteDebug("Attempting to load config from path: %s", path)

		if _, statErr := os.Stat(path); statErr == nil {
			log.WriteDebug("Config file exists, reading...")

			data, readErr := os.ReadFile(path)
			if readErr != nil {
				err = readErr
				log.WriteDebug("Failed to read config file: %v", readErr)
				return
			}

			var cfg Config
			if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
				err = unmarshalErr
				log.WriteDebug("Failed to parse config JSON: %v", unmarshalErr)
				return
			}

			ConfigData = &cfg
			log.WriteDebug("Config successfully loaded from file.")
			return
		}

		log.WriteDebug("Config file not found. Creating default config.")

		defaultCfg := Config{
			UserConsentMessage: DEFAULT_CONSENT_MESSAGE,
			APITimeout:         DEFAULT_API_TIMEOUT,
		}
		ConfigData = &defaultCfg

		if mkdirErr := os.MkdirAll(filepath.Dir(path), 0755); mkdirErr != nil {
			err = mkdirErr
			log.WriteDebug("Failed to create directory: %v", mkdirErr)
			return
		}

		jsonData, marshalErr := json.MarshalIndent(defaultCfg, "", "  ")
		if marshalErr != nil {
			err = marshalErr
			log.WriteDebug("Failed to marshal default config: %v", marshalErr)
			return
		}

		if writeErr := os.WriteFile(path, jsonData, 0644); writeErr != nil {
			err = writeErr
			log.WriteDebug("Failed to write default config to file: %v", writeErr)
			return
		}

		log.WriteDebug("Default config file created at %s", path)
	})

	if err == nil {
		log.WriteDebug("Config load process completed without error.")
	}

	return err
}

// Get returns the loaded config instance. Panics if Load() hasn't been called.
func Get() *Config {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	if ConfigData == nil {
		log.WriteError("config.Load() must be called before config.Get()")
		panic("config.Load() must be called before config.Get()")
	}
	return ConfigData
}

func CreateDefaultConfigFile(path string) error {
	defaultCfg := Config{
		UserConsentMessage: DEFAULT_CONSENT_MESSAGE,
		RunConsentMessage:  RUN_CONSENT_MESSAGE,
		APITimeout:         DEFAULT_API_TIMEOUT,
		AWSTimeout:         DEFAULT_AWS_TIMEOUT,
		AWS_URL:            DEFAULT_AWS_URL,
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(defaultCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
