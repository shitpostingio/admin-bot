package emergencymode

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/defense/antiflood"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
)

var (
	emergency                   bool
	emergencyStateUpdateChannel chan bool
)

// Start starts an emergency mode
func Start() chan bool {
	emergency = true
	emergencyStateUpdateChannel = make(chan bool)
	return emergencyStateUpdateChannel
}

// End ends an emergency mode
func End() {
	emergency = false
	emergencyStateUpdateChannel <- true
}

// RestrictUserForEmergency restricts user for an emergency
func RestrictUserForEmergency(user *tgbotapi.User, msg *tgbotapi.Message) {

	_ = adminbot.RestrictMessages(user, msg.Chat.ID, 0, &repository.Bot.Self)

	// Send a message in the report chat
	userRestrictionReportText := fmt.Sprintf(localization.GetString("emergencymode_user_restricted"), user.ID, telegram.GetName(user))
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateApproveUserButton(user.ID), buttons.CreateHandleButton())
	_ = adminbot.SendTextMessageWithMarkup(repository.GetTelegramConfiguration().ReportChannelID, userRestrictionReportText, markup, false)

	// Increase flood counter since we're in an emergency
	antiflood.IncreaseFloodCounter(1)

	// Tell the user they'll have to wait to be unrestricted
	_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, repository.Configuration.AdminBot.EmergencyText, false)

}

// PerformEmergencyModeChecks performs checks on a user and restricts them if they have
// no handle or no profile pictures during an emergency.
func PerformEmergencyModeChecks(user *tgbotapi.User, msg *tgbotapi.Message) {

	if user.UserName == "" {
		restrictUserForEmergency(user, msg)
		return
	}

	userPhotos, err := adminbot.GetUserProfilePhotos(user.ID, 1)
	if err != nil {

		unableToGetPicturesReportText := fmt.Sprintf(localization.GetString("emegencymode_unable_to_get_pictures"), user.ID, telegram.GetName(user))
		_ = adminbot.SendTextMessage(repository.GetTelegramConfiguration().ReportChannelID, unableToGetPicturesReportText, false)
		log.Warn(unableToGetPicturesReportText)
		return
	}

	if userPhotos.TotalCount == 0 {
		restrictUserForEmergency(user, msg)
	}
}

// restrictUserForEmergency mutes users for not having the requisites during an emergency.
func restrictUserForEmergency(user *tgbotapi.User, msg *tgbotapi.Message) {

	_ = adminbot.RestrictMessages(user, msg.Chat.ID, 0, &repository.Bot.Self)
	userRestrictionReportText := fmt.Sprintf(localization.GetString("emergencymode_user_restricted"), user.ID, telegram.GetName(user))

	if !antiflood.IsFlood() {

		antiflood.IncreaseFloodCounter(1)
		_ = adminbot.SendTextMessage(repository.GetTelegramConfiguration().ReportChannelID, userRestrictionReportText, false)

	}

	log.Warn(userRestrictionReportText)
}

// IsEmergency returns true if emergency mode is active.
func IsEmergency() bool {
	return emergency
}
