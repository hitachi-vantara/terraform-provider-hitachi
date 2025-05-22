package telemetry

import (
	"os"
)

const (
	// Default values
	defaultRootDir                  = "/opt/hitachi/terraform"
	defaultTelemetryDir             = defaultRootDir + "/telemetry"
	defaultUserConsentFile          = defaultRootDir + "/user_consent.json"
	defaultTelemetryStatsAWSFile    = defaultTelemetryDir + "/telemetry.json"
	defaultTelemetryAverageTimeFile = defaultTelemetryDir + "/usages.json"
	defaultAPIGURL                  = "https://5v56roefvl.execute-api.us-west-2.amazonaws.com/api/update_telemetry/ZZZ"
	// DEFAULT_CONSENT_MESSAGE is in common/config/config_model.go

)

var (
	// Use environment variables if present, otherwise fall back to default values
	TERRAFORM_TELEMETRY_DIR               = getEnv("HV_TERRAFORM_TELEMETRY_DIR", defaultTelemetryDir)
	TERRAFORM_USER_CONSENT_FILE           = getEnv("HV_TERRAFORM_USER_CONSENT_FILE", defaultUserConsentFile)
	TERRAFORM_TELEMETRY_STATS_AWS_FILE    = getEnv("HV_TERRAFORM_TELEMETRY_STATS_AWS_FILE", defaultTelemetryStatsAWSFile)
	TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE = getEnv("HV_TERRAFORM_TELEMETRY_AVERAGE_TIME_FILE", defaultTelemetryAverageTimeFile)
	APIG_URL                              = getEnv("HV_APIG_URL", defaultAPIGURL)
)

// getEnv retrieves the environment variable value or returns the default value if not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
