package adminbot

import (
	"github.com/pkg/errors"
	"github.com/shitpostingio/admin-bot/reports"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/telegram"
	utility "github.com/shitpostingio/admin-bot/utility/cache"
	log "github.com/sirupsen/logrus"
)

// BanUser bans a user given their tgbotapi.User.
// It will also add the ban data to the database and report it to the report channel.
func BanUser(bannedUser, moderator *tgbotapi.User, reason string, chatID int64) error {

	err := api.BanUser(bannedUser.ID, moderator.ID, chatID)
	if err != nil {
		return errors.Errorf("BanUser: unable to ban user with id %d: %s", bannedUser.ID, err)
	}

	updateBanData(bannedUser, moderator, reason)
	reportBan(bannedUser, moderator, reason)
	return nil

}

// BanUserByID bans a user given their userID.
// It will also add the ban data to the database and report it to the report channel.
func BanUserByID(bannedUserID int64, moderator *tgbotapi.User, reason string, chatID int64) (*tgbotapi.User, error) {

	tdlibUser, err := api.BanUserByID(bannedUserID, moderator.ID, chatID)
	if err != nil {
		return nil, errors.Errorf("BanUserByID: unable to ban user with id %d: %s", bannedUserID, err)
	}

	user := tdlib.GetTgbotapiUserFromTdlibUser(tdlibUser)
	updateBanData(user, moderator, reason)
	reportBan(user, moderator, reason)
	return user, nil

}

// BanUserByUsername bans a user their its username.
// It will also add the ban data to the database and report it to the report channel.
func BanUserByUsername(username string, moderator *tgbotapi.User, reason string, chatID int64) (*tgbotapi.User, error) {

	tdlibUser, err := api.BanUserByUsername(username, moderator.ID, chatID)
	if err != nil {
		return nil, errors.Errorf("BanUserByUsername: unable to ban user with username %s: %s", username, err)
	}

	user := tdlib.GetTgbotapiUserFromTdlibUser(tdlibUser)
	updateBanData(user, moderator, reason)
	reportBan(user, moderator, reason)
	return user, nil

}

/* ------------------------------------------------------------------------------ */

// updateBanData removes the user from the mods if they were one
// and adds the ban data to the database.
func updateBanData(bannedUser, moderator *tgbotapi.User, reason string) {

	wasMod := utility.RemoveFromMods(bannedUser.ID)
	if wasMod {

		demotionMessage := reports.ModeratorDemoted(moderator.ID, telegram.GetName(moderator), bannedUser.ID, telegram.GetName(bannedUser))
		_ = reports.Report(demotionMessage, reports.URGENT)
		log.Warn(demotionMessage)

		err := database.RemoveModerator(bannedUser.ID)
		if err != nil {
			_ = reports.Report(reports.ModeratorCannotBeRemovedFromTable(bannedUser.ID, telegram.GetName(bannedUser)), reports.NON_URGENT)
		}
	}

	_, err := database.AddBan(bannedUser.ID, moderator.ID, reason)
	if err != nil {
		log.Error("updateBanData:", err)
	}

}

// reportBan reports the ban to the report channel.
func reportBan(bannedUser, moderator *tgbotapi.User, reason string) {

	reportText := reports.UserBanned(moderator.ID, telegram.GetName(moderator), bannedUser.ID, telegram.GetName(bannedUser), reason)
	log.Info(reportText)

	reportMarkup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnbanButton(bannedUser.ID))
	err := reports.ReportWithMarkup(reportText, reportMarkup, reports.URGENT)
	if err != nil {
		log.Error("reportBan: unable to send ban report:", err)
	}
}
