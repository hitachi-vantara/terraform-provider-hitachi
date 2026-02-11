package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

var (
	ConfigData *Config // Final config after env var overrides
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

		var cfg Config

		file, openErr := os.Open(path)
		if openErr != nil {
			if os.IsNotExist(openErr) {
				log.WriteDebug("Config file not found. Creating default config.")

				if createErr := CreateDefaultConfigFile(path); createErr != nil {
					err = fmt.Errorf("failed to create default config: %w", createErr)
					log.WriteDebug("CreateDefaultConfigFile failed: %v", createErr)
					return
				}

				cfg = Config{
					ConfigNoAWS: ConfigNoAWS{
						UserConsentMessage:    DEFAULT_CONSENT_MESSAGE,
						RunConsentMessage:     RUN_CONSENT_MESSAGE,
						AddStorageNodePollMax: DEFAULT_ASN_POLL_MAX,
						APITimeout:            DEFAULT_API_TIMEOUT,
						AWSTimeout:            DEFAULT_AWS_TIMEOUT,
						MaxConcurrentOps:      DEFAULT_MAX_CONCURRENT_OPS,
						ConcurrentApiDelaySec: DEFAULT_CONCURRENT_API_DELAY_SEC,
						MaxRetries:            DEFAULT_MAX_RETRIES,
						RetryDelaySec:         DEFAULT_RETRY_DELAY_SEC,
					},
					AWS_URL: DEFAULT_AWS_URL,
				}

				log.WriteDebug("Default config created and loaded.")
			} else {
				err = fmt.Errorf("failed to open config file: %w", openErr)
				log.WriteDebug("Failed to open config file: %v", openErr)
				return
			}
		} else {
			defer file.Close()
			decoder := json.NewDecoder(file)
			if decodeErr := decoder.Decode(&cfg); decodeErr != nil {
				err = fmt.Errorf("failed to decode config file: %w", decodeErr)
				log.WriteDebug("Failed to decode config file: %v", decodeErr)
				return
			}
			log.WriteDebug("Config successfully loaded from file.")
		}

		// Apply environment overrides
		ConfigData = resolveConfigFromEnv(&cfg)
		log.WriteDebug("Config applied with environment overrides.")
	})

	if err == nil {
		log.WriteDebug("Config load process completed without error.")
		log.WriteDebug("Config: %+v", ConfigData)
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
	// no AWS_URL allowed in default config
	defaultCfg := ConfigNoAWS{
		UserConsentMessage:    DEFAULT_CONSENT_MESSAGE,
		RunConsentMessage:     RUN_CONSENT_MESSAGE,
		AddStorageNodePollMax: DEFAULT_ASN_POLL_MAX,
		APITimeout:            DEFAULT_API_TIMEOUT,
		AWSTimeout:            DEFAULT_AWS_TIMEOUT,
		MaxConcurrentOps:      DEFAULT_MAX_CONCURRENT_OPS,
		ConcurrentApiDelaySec: DEFAULT_CONCURRENT_API_DELAY_SEC,
		MaxRetries:            DEFAULT_MAX_RETRIES,
		RetryDelaySec:         DEFAULT_RETRY_DELAY_SEC,
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

func resolveConfigFromEnv(cfg *Config) *Config {
	if cfg == nil {
		cfg = &Config{}
	}

	addStorageNodePollMax := getEnvInt("ASN_TIMEOUT", cfg.AddStorageNodePollMax, DEFAULT_ASN_POLL_MAX)
	apiTimeout := getEnvInt("API_TIMEOUT", cfg.APITimeout, DEFAULT_API_TIMEOUT)
	awsTimeout := getEnvInt("AWS_TIMEOUT", cfg.AWSTimeout, DEFAULT_AWS_TIMEOUT)
	awsUrl := getEnvStr("AWS_URL", cfg.AWS_URL, DEFAULT_AWS_URL)
	maxConcurrentOps := getEnvInt("MAX_CONCURRENT_OPS", cfg.MaxConcurrentOps, DEFAULT_MAX_CONCURRENT_OPS)
	concurrentApiDelaySec := getEnvInt("CONCURRENT_API_DELAY_SEC", cfg.ConcurrentApiDelaySec, DEFAULT_CONCURRENT_API_DELAY_SEC)
	maxRetries := getEnvInt("MAX_RETRIES", cfg.MaxRetries, DEFAULT_MAX_RETRIES)
	retryDelaySec := getEnvInt("RETRY_DELAY_SEC", cfg.RetryDelaySec, DEFAULT_RETRY_DELAY_SEC)

	userConsent := DEFAULT_CONSENT_MESSAGE
	if cfg.UserConsentMessage != "" {
		userConsent = cfg.UserConsentMessage
	}

	runConsent := RUN_CONSENT_MESSAGE
	if cfg.RunConsentMessage != "" {
		runConsent = cfg.RunConsentMessage
	}

	return &Config{
		ConfigNoAWS: ConfigNoAWS{
			AddStorageNodePollMax: addStorageNodePollMax,
			APITimeout:            apiTimeout,
			AWSTimeout:            awsTimeout,
			UserConsentMessage:    userConsent,
			RunConsentMessage:     runConsent,
			MaxConcurrentOps:      maxConcurrentOps,
			ConcurrentApiDelaySec: concurrentApiDelaySec,
			MaxRetries:            maxRetries,
			RetryDelaySec:         retryDelaySec,
		},
		AWS_URL: awsUrl,
	}
}

func getEnvStr(key, fromConfig, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	if fromConfig != "" {
		return fromConfig
	}
	return fallback
}

func getEnvInt(key string, fromConfig, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			return num
		}
	}
	if fromConfig > 0 {
		return fromConfig
	}
	return fallback
}
