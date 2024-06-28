package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// KickUser kicks a user using the bot API.
func KickUser(kickedUserID int64, chatID int64) (*tgbotapi.APIResponse, error) {

	response, err := bot.Request(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			UserID: kickedUserID,
			ChatID: chatID,
		}})

	return response, err
}
