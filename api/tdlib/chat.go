package tdlib

import (
	"github.com/shitpostingio/go-tdlib/client"
)

// GetChatMember gets the chat member info via tdlib
func GetChatMember(chatID int64, userID int) (*client.ChatMember, error) {

	chatMember, err := tdlibClient.GetChatMember(&client.GetChatMemberRequest{
		ChatID: chatID,
		UserID: int64(userID),
	})

	return chatMember, err
}
