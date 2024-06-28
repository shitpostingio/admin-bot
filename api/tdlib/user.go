package tdlib

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shitpostingio/go-tdlib/client"
)

// GetTgbotapiUserFromTdlibUser returns the equivalent tgbotapi.User structure
// of the input tdlib.client.User structure.
func GetTgbotapiUserFromTdlibUser(tlu *client.User) *tgbotapi.User {

	if tlu == nil {
		return nil
	}

	return &tgbotapi.User{
		ID:        int64(tlu.ID),
		FirstName: tlu.FirstName,
		LastName:  tlu.LastName,
		UserName:  tlu.Username,
	}
}

// GetUserByID returns a tdlib.client.User given a userID.
func GetUserByID(userID int64) (*client.User, error) {

	user, err := tdlibClient.GetUser(&client.GetUserRequest{UserID: userID})
	if err != nil {
		return nil, errors.Errorf("BanUserByID.GetUser: %s", err)
	}

	return user, nil
}
