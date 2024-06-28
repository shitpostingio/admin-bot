package adminbot

import (
	"github.com/shitpostingio/admin-bot/api"
)

// SendCallbackWithAlert answer a callback query with an alert popup
func SendCallbackWithAlert(id, text string) error {
	return api.SendCallbackWithAlert(id, text)
}
