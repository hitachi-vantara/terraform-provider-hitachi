package utils

import (
	"fmt"
	"strings"
	"time"

	config "terraform-provider-hitachi/hitachi/common/config"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries int
	Delay      int
}

// getRetryConfig returns retry configuration using defaults or overrides from config.ConfigData
func GetRetryConfig() RetryConfig {
	maxRetries := config.DEFAULT_MAX_RETRIES
	retryDelaySec := config.DEFAULT_RETRY_DELAY_SEC

	if config.ConfigData != nil {
		if config.ConfigData.MaxRetries > 0 {
			maxRetries = config.ConfigData.MaxRetries
		}
		if config.ConfigData.RetryDelaySec > 0 {
			retryDelaySec = config.ConfigData.RetryDelaySec
		}
	}

	return RetryConfig{
		MaxRetries: maxRetries,
		Delay:      retryDelaySec,
	}
}

// ExecuteWithRetry wraps HTTP job operations and retries on transient errors
func ExecuteWithRetry(operationName string, fn func() (*string, error)) (*string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	retryCfg := GetRetryConfig()
	var lastErr error

	for attempt := 1; attempt <= retryCfg.MaxRetries; attempt++ {
		log.WriteDebug(fmt.Sprintf("Attempt %d/%d for %s", attempt, retryCfg.MaxRetries, operationName))

		resIDs, err := fn()
		if err == nil {
			return resIDs, nil
		}

		lastErr = err
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "wait for a while and try again") ||  strings.Contains(errStr, "503") {
			log.WriteWarn(fmt.Sprintf("%s failed: %v — retrying in %v", operationName, err, retryCfg.Delay))
			time.Sleep(time.Duration(retryCfg.Delay) * time.Second)
			continue
		}

		// Non-retriable error
		log.WriteError(fmt.Errorf("%s failed: %v (not retriable)", operationName, err))
		return nil, err
	}

	log.WriteError(fmt.Errorf("%s failed after %d attempts: %v", operationName, retryCfg.MaxRetries, lastErr))
	return nil, lastErr
}
