package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

//PromoteUser promotes a user to admin.
// It will also use a rate limiter not to get restricted by Telegram.
func PromoteUser(userID int64, chatID int64,
	canChangeInfo,
	canDeleteMessages,
	canInviteUsers,
	canRestrictMembers,
	canPinMessages,
	canPromoteMembers bool) (*tgbotapi.APIResponse, error) {

	limiter.AuthorizeUrgentAction()
	return botapi.PromoteUser(userID, chatID,
		canChangeInfo,
		canDeleteMessages,
		canInviteUsers,
		canRestrictMembers,
		canPinMessages,
		canPromoteMembers)
}
