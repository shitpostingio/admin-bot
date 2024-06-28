package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LeaveChat leaves a chat using the bot API.
func LeaveChat(chatID int64) (*tgbotapi.APIResponse, error) {
	response, err := bot.Request(tgbotapi.LeaveChatConfig{ChatID: chatID})
	return response, err
}
