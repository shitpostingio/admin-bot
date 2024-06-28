package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// DeleteMessage deletes a message using the bot API.
func DeleteMessage(chatID int64, messageID int) error {
	_, err := bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
	return err
}
