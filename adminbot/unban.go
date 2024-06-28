package adminbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/api/cache"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/telegram"
	log "github.com/sirupsen/logrus"
)

// UnbanUser unbans a user in a chat and marks them as unbanned in the database.
func UnbanUser(userID int64, chatID int64, moderator *tgbotapi.User) (err error) {

	_, err = api.UnbanUser(userID, chatID)
	if err != nil {
		return errors.Errorf("UnbanUser: could not unban user with ID %d: %s", userID, err)
	}

	updateUnbanData(userID, moderator)
	return nil
}

// updateUnbanData marks the user as unbanned in the database
// and removes them from the ban cache.
func updateUnbanData(userID int64, moderator *tgbotapi.User) {

	// Remove the ban from the cache so
	// the user can be banned again on need
	cache.RemoveBanFromCache(userID)

	//
	_ = database.MarkUserAsUnbanned(userID)
	log.Info(telegram.GetNameOrUsername(moderator), "unbanned user with TelegramID", userID)
}
