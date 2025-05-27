package config

type Config struct {
	UserConsentMessage string `json:"user_consent_message"`
	APITimeout         int    `json:"api_timeout"` // in seconds
	AWSTimeout         int    `json:"aws_timeout"` // in seconds
	AWS_URL            string `json:"aws_url"`
}

const DEFAULT_CONSENT_MESSAGE = `Hitachi Vantara LLC collects usage data such as:

  • Storage model
  • Storage serial number
  • Operation name
  • Status (success or failure)
  • Duration

This data is collected for product improvement purposes only.
It remains confidential and is not shared with any third parties.`

const (
	DEFAULT_API_TIMEOUT = 300 // seconds
	DEFAULT_AWS_TIMEOUT = 300 // seconds
	// DEFAULT_AWS_URL     = "https://5v56roefvl.execute-api.us-west-2.amazonaws.com/api/update_telemetry/ZZZ"
	DEFAULT_AWS_URL = ""
)
