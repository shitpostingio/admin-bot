package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ForwardMessage forwards a message using the bot API.
func ForwardMessage(toChatID, fromChatID int64, fromMessageID int) (*tgbotapi.Message, error) {
	forward, err := bot.Send(tgbotapi.NewForward(toChatID, fromChatID, fromMessageID))
	return &forward, err
}
