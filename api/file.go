package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api/botapi"
	"github.com/shitpostingio/admin-bot/api/cache"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// GetTelegramFile gets a Telegram file given its fileID.
// It will also use a rate limiter not to get restricted by Telegram.
func GetTelegramFile(uniqueFileID, fileID string) (*tgbotapi.File, error) {

	element, found := cache.CheckFilePathCache(uniqueFileID)
	if found {
		element.Mutex.RLock()
		defer element.Mutex.RUnlock()
		if element.Performed {
			return element.File, nil
		}
	}

	element.Mutex.Lock()
	limiter.AuthorizeAction()
	result, err := botapi.GetTelegramFile(fileID)
	if err == nil {
		element.Performed = true
		element.File = result
	} else {
		log.Error("GetTelegramFile: error while performing request", err)
	}

	element.Mutex.Unlock()
	return result, err
}
