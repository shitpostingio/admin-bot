package adminbot

import (
	"github.com/shitpostingio/admin-bot/api"
)

// ForwardMessage forwards a message in a rate limited fashion.
func ForwardMessage(toChatID, fromChatID int64, fromMessageID int) (int, error) {
	msg, err := api.ForwardMessage(toChatID, fromChatID, fromMessageID)
	if err != nil {
		return 0, err
	}

	return msg.MessageID, err
}
