package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// ForwardMessage forwards a message.
// It will also use a rate limiter not to get restricted by Telegram.
func ForwardMessage(toChatID, fromChatID int64, fromMessageID int) (*tgbotapi.Message, error) {
	limiter.AuthorizeAction()
	return botapi.ForwardMessage(toChatID, fromChatID, fromMessageID)
}
