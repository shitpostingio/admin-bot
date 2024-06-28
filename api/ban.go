package api

import (
	"github.com/pkg/errors"
	"github.com/shitpostingio/admin-bot/api/botapi"
	"github.com/shitpostingio/admin-bot/api/cache"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/go-tdlib/client"
)

// BanUser bans a user using the bot API.
// It will make sure no ban for the same person has already been requested.
// It will also use a rate limiter not to get restricted by Telegram.
func BanUser(bannedUserID, moderatorID int64, chatID int64) error {

	ban, err := getBanMutex(bannedUserID, moderatorID)
	if err != nil {
		return err
	}

	//
	defer ban.Mutex.Unlock()
	limiter.AuthorizeUrgentAction()

	//
	err = tdlib.BanUser(bannedUserID, chatID)
	if err != nil {
		return err
	}

	ban.Performed = true
	return nil
}

// BanUserByID bans a user using tdlib.
// It will make sure no ban for the same person has already been requested.
// It will also use a rate limiter not to get restricted by Telegram.
func BanUserByID(bannedUserID, moderatorID int64, chatID int64) (*client.User, error) {

	ban, err := getBanMutex(bannedUserID, moderatorID)
	if err != nil {
		return nil, err
	}

	//
	defer ban.Mutex.Unlock()
	limiter.AuthorizeUrgentAction()

	//
	user, err := tdlib.BanUserByID(bannedUserID, chatID)
	if err != nil {
		return nil, err
	}

	ban.Performed = true
	return user, nil
}

// BanUserByUsername bans a user by their username using tdlib.
// It will make sure no ban for the same person has already been requested.
// It will also use a rate limiter not to get restricted by Telegram.
func BanUserByUsername(username string, moderatorID int64, chatID int64) (*client.User, error) {

	chat, err := tdlib.ResolveUsername(username)
	if err != nil {
		return nil, errors.Errorf("BanUserByUsername.ResolveUsername: %s", err)
	}

	if chat.Type.ChatTypeType() != client.TypeChatTypePrivate {
		return nil, errors.Errorf("BanUserByUsername: %s is not a user", username)
	}

	bannedUserID := chat.ID
	ban, err := getBanMutex(bannedUserID, moderatorID)
	if err != nil {
		return nil, err
	}

	//
	defer ban.Mutex.Unlock()
	limiter.AuthorizeUrgentAction()
	user, err := tdlib.BanUserByID(bannedUserID, chatID)
	if err != nil {
		return nil, err
	}

	ban.Performed = true
	return user, nil
}

func BanChannel(bannedSenderChatID, moderatorID, chatID int64) error {
	ban, err := getBanMutex(bannedSenderChatID, moderatorID)
	if err != nil {
		return err
	}

	//
	defer ban.Mutex.Unlock()
	limiter.AuthorizeUrgentAction()
	err = botapi.BanChannel(bannedSenderChatID, chatID)
	if err != nil {
		return err
	}

	ban.Performed = true
	return nil
}

/* ------------------------------------------------------------------------------ */

// getBanMutex makes sure only one ban action can be performed per bannedUserID.
func getBanMutex(bannedUserID, moderatorUserID int64) (*cache.Action, error) {

	//Check if the hierarchy allows the ban.
	if repository.Admins[bannedUserID] {
		return nil, errors.New("getBanMutex: admins can't be banned")
	}

	if repository.Mods[bannedUserID] && !repository.Admins[moderatorUserID] {
		return nil, errors.New("getBanMutex: mods can't ban other mods")
	}

	//Get data from the cache to perform the ban
	banAction, err := cache.AddBanToCache(bannedUserID)
	if err != nil {
		banAction, err = cache.GetBanFromCache(bannedUserID)
		if err != nil {
			return nil, errors.Errorf("getBanMutex: ban cache error")
		}
	}

	banAction.Mutex.Lock()
	if banAction.Performed {
		return nil, errors.Errorf("getBanMutex: cache hit for user with id %d", bannedUserID)
	}

	return banAction, nil
}
