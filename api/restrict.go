package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// RestrictUser requests the restriction of a user in a chat.
// If the user was a mod, they'll be also removed from the mods map.
// It will also use a rate limiter not to get restricted by Telegram.
func RestrictUser(userID int64, chatID int64, untilDate int64,
	canSendMessages,
	canSendMediaMessages,
	canSendOtherMessages,
	canAddWebPagePreviews bool) (*tgbotapi.APIResponse, error) {

	limiter.AuthorizeUrgentAction()
	return botapi.RestrictUser(userID, chatID, untilDate,
		canSendMessages,
		canSendMediaMessages,
		canSendOtherMessages,
		canAddWebPagePreviews)
}

// UnrestrictUser requests the unrestriction of a user in a chat.
// It will also use a rate limiter not to get restricted by Telegram.
func UnrestrictUser(userID int64, chatID int64) (*tgbotapi.APIResponse, error) {
	limiter.AuthorizeAction()
	return botapi.UnrestrictUser(userID, chatID)
}

//RestrictMessages restricts an user from sending messages.
func RestrictMessages(userID int64, chatID int64, restrictionEndTime int64) (*tgbotapi.APIResponse, error) {
	return RestrictUser(userID, chatID, restrictionEndTime,
		false,
		false,
		false,
		false)
}

//RestrictMedia restricts an user from sending media.
func RestrictMedia(userID int64, chatID int64, restrictionEndTime int64) (*tgbotapi.APIResponse, error) {
	return RestrictUser(userID, chatID, restrictionEndTime,
		true,
		false,
		false,
		false)
}

//RestrictOther restricts an user from sending stickers and gifs.
func RestrictOther(userID int64, chatID int64, restrictionEndTime int64) (*tgbotapi.APIResponse, error) {
	return RestrictUser(userID, chatID, restrictionEndTime,
		true,
		true,
		false,
		false)
}
