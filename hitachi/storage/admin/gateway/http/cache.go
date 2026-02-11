package admin

import (
	"time"

	// log "github.com/romana/rlog"
	"github.com/patrickmn/go-cache"
)

var DEFAULT_CACHE_STORAGE_DURATION = 270 * time.Second  // slightly less than 300s backend timeout
var cacheStorageDuration = DEFAULT_CACHE_STORAGE_DURATION
var cacheStorage = cache.New(cacheStorageDuration, 2*cacheStorageDuration)

// func SetCacheStorageDuration(duration time.Duration) {
// 	cacheStorageDuration = duration
// 	log.Infof("NEW CACHE STORAGE DURATION: %v\n", cacheStorageDuration)
// }

// func DeleteItemsCacheStorage() {
// 	cacheStorageDuration = DEFAULT_CACHE_STORAGE_DURATION
// 	cacheStorage.Flush()
// 	log.Infof("CACHE STORAGE: flushed and reset default duration to %v\n", cacheStorageDuration)
// }
