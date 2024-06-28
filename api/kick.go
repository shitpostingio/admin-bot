package api

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
	"github.com/shitpostingio/admin-bot/repository"
)

// KickUser kicks a user in a chat.
// It will make sure the hierarchy allows the kick to be performed.
// It will also use a rate limiter not to get restricted by Telegram.
func KickUser(kickedUserID, moderatorUserID int64, chatID int64) (*tgbotapi.APIResponse, error) {

	err := authorizeKick(kickedUserID, moderatorUserID)
	if err != nil {
		return nil, err
	}

	limiter.AuthorizeAction()
	return botapi.KickUser(kickedUserID, chatID)
}

/* ------------------------------------------------------------------------------ */

// authorizeKick checks if the hierarchy allows the kick to be performed.
func authorizeKick(kickedUserID, moderatorUserID int64) error {

	//Check if the hierarchy allows the kick.
	if repository.Admins[kickedUserID] {
		return errors.New("KickUser: admins can't be kicked")
	}

	if repository.Mods[kickedUserID] && !repository.Admins[moderatorUserID] {
		return errors.New("KickUser: mods can't kick admins")
	}

	return nil
}
