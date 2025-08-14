package telemetry

import (
	"os"
)

const (
	TERRAFORM_ROOT_DIR                    = "/opt/hitachi/terraform"
	TERRAFORM_TELEMETRY_DIR               = TERRAFORM_ROOT_DIR + "/telemetry"
	TERRAFORM_USER_CONSENT_FILE           = TERRAFORM_ROOT_DIR + "/user_consent.json"
	TERRAFORM_TELEMETRY_STATS_AWS_FILE    = TERRAFORM_TELEMETRY_DIR + "/telemetry.json"
	TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE = TERRAFORM_TELEMETRY_DIR + "/usages.json"
	// also see hitachi/common/config/config_model.go
)

// getEnv retrieves the environment variable value or returns the default value if not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
