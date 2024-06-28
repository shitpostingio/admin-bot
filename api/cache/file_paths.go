package cache

import (
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	filePathCacheExpiration = 1 * time.Hour
	filePathCacheCleanup    = 2 * time.Hour
)

var (
	fpCache = cache.New(filePathCacheExpiration, filePathCacheCleanup)
)

type FilePathCacheElement struct {
	File      *tgbotapi.File
	Mutex     sync.RWMutex
	Performed bool
}

func CheckFilePathCache(fileUniqueID string) (*FilePathCacheElement, bool) {

	value, found := fpCache.Get(fileUniqueID)
	if found {
		return value.(*FilePathCacheElement), true
	}

	var fp FilePathCacheElement
	err := fpCache.Add(fileUniqueID, &fp, cache.DefaultExpiration)
	if err != nil {
		log.Error("FPCache: error", err, "for fileUniqueID", fileUniqueID)
	}

	return &fp, found

}
