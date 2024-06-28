package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

//UnbanUser unbans a user in a chat and marks them as unbanned in the database.
// It will also use a rate limiter not to get restricted by Telegram.
func UnbanUser(userID int64, chatID int64) (*tgbotapi.APIResponse, error) {
	limiter.AuthorizeAction()
	return botapi.UnbanUser(userID, chatID)
}
