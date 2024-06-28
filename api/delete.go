package api

import (
	"github.com/pkg/errors"
	"github.com/shitpostingio/admin-bot/api/botapi"
	"github.com/shitpostingio/admin-bot/api/cache"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// DeleteMessage deletes a message in a chat.
// It will also use a rate limiter not to get restricted by Telegram.
func DeleteMessage(chatID int64, messageID int) error {

	messageIDSlice := []int{messageID}

	if !cache.CheckDeletionCache(messageIDSlice) {
		return errors.Errorf("DeleteMessage: cache hit for messageID %d", messageID)
	}

	limiter.AuthorizeAction()
	err := botapi.DeleteMessage(chatID, messageID)
	if err == nil {
		cache.AddToDeletionCache(messageIDSlice)
	}

	return err
}

func DeleteMultipleMessages(chatID int64, messageIDs []int) error {

	if !cache.CheckDeletionCache(messageIDs) {
		return errors.Errorf("DeleteMultipleMessages: cache hit for messageIDs %v", messageIDs)
	}

	limiter.AuthorizeAction()
	err := tdlib.DeleteMultipleMessages(chatID, messageIDs)
	if err == nil {
		cache.AddToDeletionCache(messageIDs)
	}

	return err
}
