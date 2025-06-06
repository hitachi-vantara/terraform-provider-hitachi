package config

const CONFIG_FILE = "/opt/hitachi/terraform/bin/.internal_config"

// Package config provides functionality to load and manage configuration settings for the Hitachi Terraform provider.
type Config struct {
	UserConsentMessage string `json:"user_consent_message"`
	RunConsentMessage  string `json:"run_consent_message"`
	APITimeout         int    `json:"api_timeout"` // in seconds
	AWSTimeout         int    `json:"aws_timeout"` // in seconds
	AWS_URL            string `json:"aws_url"`
}

// ConfigNoAWS is a reduced version of Config without AWS_URL.
type ConfigNoAWS struct {
	UserConsentMessage string `json:"user_consent_message"`
	RunConsentMessage  string `json:"run_consent_message"`
	APITimeout         int    `json:"api_timeout"` // in seconds
	AWSTimeout         int    `json:"aws_timeout"` // in seconds
}

const DEFAULT_CONSENT_MESSAGE = `
  Hitachi Vantara LLC collects usage data such as storage model, storage serial number, operation name, status (success or failure),
  and duration. This data is collected for product improvement purposes only. It remains confidential and it is not shared with any
  third parties. `

const RUN_CONSENT_MESSAGE = "To provide your consent, run bin/user_consent.sh from /opt/hitachi/terraform."

const (
	DEFAULT_API_TIMEOUT = 300 // seconds
	DEFAULT_AWS_TIMEOUT = 300 // seconds
	// DEFAULT_AWS_URL     = "https://5v56roefvl.execute-api.us-west-2.amazonaws.com/api/update_telemetry/ZZZ"
	DEFAULT_AWS_URL = "" // TODO: set this to the actual AWS URL if needed, or leave it empty for local testing purposes.
)
