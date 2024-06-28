package tdlib

import (
	"github.com/shitpostingio/go-tdlib/client"
)

// DeleteMessage deletes a message using the bot API.
func DeleteMessage(chatID int64, messageID int) error {

	tdlibMessageID := getTdlibMessageID(messageID)

	_, err := tdlibClient.DeleteMessages(&client.DeleteMessagesRequest{
		ChatID:     chatID,
		MessageIDs: []int64{tdlibMessageID},
		Revoke:     true,
	})

	return err
}

func DeleteMultipleMessages(chatID int64, messageIDs []int) error {

	messages := make([]int64, len(messageIDs))

	for i, messageID := range messageIDs {
		messages[i] = getTdlibMessageID(messageID)
	}

	_, err := tdlibClient.DeleteMessages(&client.DeleteMessagesRequest{
		ChatID:     chatID,
		MessageIDs: messages,
		Revoke:     true,
	})

	return err
}
