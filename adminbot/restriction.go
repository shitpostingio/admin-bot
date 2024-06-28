package adminbot

import (
	"github.com/shitpostingio/admin-bot/reports"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/telegram"
	"github.com/shitpostingio/admin-bot/utility"
	"github.com/shitpostingio/admin-bot/utility/cache"
)

// RestrictUser requests the restriction of a user in a chat.
// If the user was a mod, they'll be also removed from the mods map.
func RestrictUser(user *tgbotapi.User, chatID int64, untilDate int64,
	canSendMessages,
	canSendMediaMessages,
	canSendOtherMessages,
	canAddWebPagePreviews bool,
	admin *tgbotapi.User) error {

	_, err := api.RestrictUser(user.ID, chatID, untilDate,
		canSendMessages,
		canSendMediaMessages,
		canSendOtherMessages,
		canAddWebPagePreviews)
	if err != nil {
		log.Error("Unable to restrict", telegram.GetNameOrUsername(user), ":", err)
		return err
	}

	log.Info(telegram.GetNameOrUsername(admin), "restricted", telegram.GetNameOrUsername(user), "until", utility.FormatUnixDate(untilDate))
	wasMod := cache.RemoveFromMods(user.ID)
	if wasMod {

		demotionMessage := reports.ModeratorDemoted(admin.ID, telegram.GetName(admin), user.ID, telegram.GetName(user))
		_ = reports.Report(demotionMessage, reports.URGENT)
		log.Warn(demotionMessage)

		err := database.RemoveModerator(user.ID)
		if err != nil {
			_ = reports.Report(reports.ModeratorCannotBeRemovedFromTable(user.ID, telegram.GetName(user)), reports.NON_URGENT)
		}
	}

	return nil

}

// UnrestrictUser requests the unrestriction of a user in a chat.
func UnrestrictUser(user *tgbotapi.User, chatID int64, admin *tgbotapi.User) error {

	_, err := api.UnrestrictUser(user.ID, chatID)
	if err != nil {
		log.Error("Unable to unrestrict", telegram.GetNameOrUsername(user), ":", err)
		return err
	}

	log.Info(telegram.GetNameOrUsername(admin), "unrestricted", telegram.GetNameOrUsername(user))
	return nil
}

//RestrictMessages restricts an user from sending messages
func RestrictMessages(user *tgbotapi.User, chatID int64, restrictionEndTime int64, admin *tgbotapi.User) error {
	return RestrictUser(user, chatID, restrictionEndTime,
		false,
		false,
		false,
		false,
		admin)
}

//RestrictMedia restricts an user from sending media
func RestrictMedia(user *tgbotapi.User, chatID int64, restrictionEndTime int64, admin *tgbotapi.User) error {
	return RestrictUser(user, chatID, restrictionEndTime,
		true,
		false,
		false,
		false,
		admin)
}

//RestrictOther restricts an user from sending stickers and gifs
func RestrictOther(user *tgbotapi.User, chatID int64, restrictionEndTime int64, admin *tgbotapi.User) error {
	return RestrictUser(user, chatID, restrictionEndTime,
		true,
		true,
		false,
		false,
		admin)
}
