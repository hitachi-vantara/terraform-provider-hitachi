package config

const CONFIG_FILE = "/opt/hitachi/terraform/bin/.internal_config"

// Package config provides functionality to load and manage configuration settings for the Hitachi Terraform provider.
type Config struct {
	ConfigNoAWS        // Embedding ConfigNoAWS to include its fields
	AWS_URL     string `json:"aws_url"`
}

// ConfigNoAWS is a reduced version of Config without AWS_URL.
type ConfigNoAWS struct {
	UserConsentMessage    string `json:"user_consent_message"`
	RunConsentMessage     string `json:"run_consent_message"`
	AddStorageNodePollMax int    `json:"add_storage_node_poll_max"`
	APITimeout            int    `json:"api_timeout"` // in seconds
	AWSTimeout            int    `json:"aws_timeout"` // in seconds
	MaxConcurrentOps      int    `json:"max_concurrent_ops"`
	ConcurrentApiDelaySec int    `json:"concurrent_api_delay_sec"` // in seconds
	MaxRetries            int    `json:"max_retries"`
	RetryDelaySec         int    `json:"retry_delay_sec"` // in seconds
}

const DEFAULT_CONSENT_MESSAGE = `
  Hitachi Vantara LLC collects usage data such as storage model, storage serial number, operation name, status (success or failure),
  and duration. This data is collected for product improvement purposes only. It remains confidential and it is not shared with any
  third parties. `

const RUN_CONSENT_MESSAGE = "To provide your consent, run bin/user_consent.sh from /opt/hitachi/terraform."

const (
	DEFAULT_ASN_POLL_MAX = 90  // minutes for adding storage node
	DEFAULT_API_TIMEOUT  = 300 // seconds
	DEFAULT_AWS_TIMEOUT  = 300 // seconds
	// DEFAULT_AWS_URL     = "https://5v56roefvl.execute-api.us-west-2.amazonaws.com/api/update_telemetry"
	DEFAULT_AWS_URL = "https://guuc07ks4j.execute-api.us-west-2.amazonaws.com/api/update_telemetry"

	DEFAULT_MAX_CONCURRENT_OPS       = 10
	DEFAULT_CONCURRENT_API_DELAY_SEC = 2
	DEFAULT_MAX_RETRIES              = 8
	DEFAULT_RETRY_DELAY_SEC          = 10
)
