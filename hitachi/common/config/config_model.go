package config

type Config struct {
	UserConsentMessage string `json:"user_consent_message"`
	APITimeout         int    `json:"api_timeout"` // in seconds
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
)
