package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendCallbackWithAlert sends a callback response with an alert using the bot API.
func SendCallbackWithAlert(id, text string) error {
	_, err := bot.Request(tgbotapi.NewCallbackWithAlert(id, text))
	return err
}
