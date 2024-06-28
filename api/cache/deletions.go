package cache

import (
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	deletionCacheExpiration = 5 * time.Minute
	deletionCacheCleanup    = 10 * time.Minute
)

var (
	deletionCache = cache.New(deletionCacheExpiration, deletionCacheCleanup)
	deletionMutex sync.RWMutex
)

func CheckDeletionCache(messageIDs []int) bool {

	deletionMutex.RLock()
	defer deletionMutex.RUnlock()

	for _, messageID := range messageIDs {
		_, found := deletionCache.Get(strconv.Itoa(messageID))
		if !found {
			return true
		}
	}

	return false

}

func AddToDeletionCache(messageIDs []int) {
	deletionMutex.Lock()
	for _, messageID := range messageIDs {
		_ = deletionCache.Add(strconv.Itoa(messageID), true, cache.DefaultExpiration)
	}

	deletionMutex.Unlock()
}
