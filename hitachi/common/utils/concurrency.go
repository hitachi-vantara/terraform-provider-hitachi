package utils

import (
	"fmt"
	"sync"
	"time"

	config "terraform-provider-hitachi/hitachi/common/config"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

type ConcurrencyConfig struct {
	MaxConcurrentOps int
	ApiDelaySec      int
}

// GetConcurrencyConfig returns the configured concurrency and cfg.ApiDelaySec values.
// It falls back to default constants if ConfigData is nil or fields are unset.
func GetConcurrencyConfig() ConcurrencyConfig {
	maxConcurrentOps := config.DEFAULT_MAX_CONCURRENT_OPS
	apiDelay := config.DEFAULT_CONCURRENT_API_DELAY_SEC

	if config.ConfigData != nil {
		if config.ConfigData.MaxConcurrentOps > 0 {
			maxConcurrentOps = config.ConfigData.MaxConcurrentOps
		}
		if config.ConfigData.ConcurrentApiDelaySec > 0 {
			apiDelay = config.ConfigData.ConcurrentApiDelaySec
		}
	}

	return ConcurrencyConfig{
		MaxConcurrentOps: maxConcurrentOps,
		ApiDelaySec:      apiDelay,
	}
}

// RunConcurrentOperations provides reusable concurrency and throttling control
func RunConcurrentOperations[T any](
	label string, // "Read" or "Delete"
	items []T, // items to process
	fn func(T, int) error, // operation
) []error {
	log := commonlog.GetLogger()

	if len(items) == 0 {
		log.WriteInfo("[%s] No items to process.", label)
		return nil
	}

	cfg := GetConcurrencyConfig()
	fmt.Printf("Max ops: %d, API cfg.ApiDelaySec: %d sec\n", cfg.MaxConcurrentOps, cfg.ApiDelaySec)

	log.WriteInfo("[%s] Starting %d operations (cfg.MaxConcurrentOps=%d, cfg.ApiDelaySec=%dms)", label, len(items), cfg.MaxConcurrentOps, cfg.ApiDelaySec)

	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.MaxConcurrentOps)
	errs := make([]error, len(items))

	start := time.Now()

	for i, item := range items {
		sem <- struct{}{}
		wg.Add(1)

		go func(idx int, it T) {
			defer wg.Done()
			defer func() { <-sem }()

			errs[idx] = fn(it, idx)
			if cfg.ApiDelaySec > 0 && idx < len(items)-1 {
				time.Sleep(time.Duration(cfg.ApiDelaySec) * time.Millisecond)
			}
		}(i, item)
	}

	wg.Wait()
	log.WriteInfo("[%s] All operations completed in %v", label, time.Since(start))
	return errs
}
