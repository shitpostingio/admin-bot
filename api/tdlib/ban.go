package tdlib

import (
	"github.com/pkg/errors"
	"github.com/shitpostingio/go-tdlib/client"
)

func BanUser(bannedUserID int64, chatID int64) error {

	_, err := tdlibClient.SetChatMemberStatus(&client.SetChatMemberStatusRequest{
		ChatID:   chatID,
		MemberID: &client.MessageSenderUser{UserID: bannedUserID},
		Status:   &client.ChatMemberStatusBanned{BannedUntilDate: 0},
	})

	return err
}

// BanUserByID bans a user given their userID.
func BanUserByID(bannedUserID int64, chatID int64) (*client.User, error) {

	_, err := tdlibClient.GetChat(&client.GetChatRequest{ChatID: chatID})
	if err != nil {
		return nil, errors.Errorf("BanUserByID.GetChat: %s", err)
	}

	user, err := GetUserByID(bannedUserID)
	if err != nil {
		return nil, errors.Errorf("BanUserByID.GetUser: %s", err)
	}

	err = BanUser(bannedUserID, chatID)

	return user, err
}
