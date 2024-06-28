package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// LeaveChat leaves a chat given its chatID.
// It will also use a rate limiter not to get restricted by Telegram.
func LeaveChat(chatID int64) (*tgbotapi.APIResponse, error) {
	limiter.AuthorizeAction()
	return botapi.LeaveChat(chatID)
}
