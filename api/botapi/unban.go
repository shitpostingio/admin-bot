package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UnbanUser unbans a user in a chat using the bot API.
func UnbanUser(userID int64, chatID int64) (*tgbotapi.APIResponse, error) {

	response, err := bot.Request(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			UserID: userID,
			ChatID: chatID,
		}})

	return response, err
}
