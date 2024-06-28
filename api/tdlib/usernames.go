package tdlib

import (
	"github.com/shitpostingio/go-tdlib/client"
)

// ResolveUsername searches public chats to find the corresponding Chat for the input username.
func ResolveUsername(username string) (chat *client.Chat, err error) {
	chat, err = tdlibClient.SearchPublicChat(&client.SearchPublicChatRequest{Username: username})
	return chat, err
}

// IsGroupOrChannelUsername returns true if the input username belongs
// does not belong to a private chat.
func IsGroupOrChannelUsername(username string) bool {

	chat, err := ResolveUsername(username)
	if err != nil {
		return false
	}

	return chat.Type.ChatTypeType() != client.TypeChatTypePrivate
}
