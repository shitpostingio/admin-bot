package private

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/defense/emergencymode"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/shitpostingio/admin-bot/database/database"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//HandlePrivateCommand allows admins to give commands to the bot via private messages
func HandlePrivateCommand(msg *tgbotapi.Message) {

	command := strings.ToLower(msg.Command())

	if command == "cancel" {
		text := localization.GetString("private_command_operation_cancelled")
		markup := tgbotapi.NewRemoveKeyboard(true)
		_ = adminbot.SendPlainTextMessageWithMarkup(msg.Chat.ID, text, markup, false)
		return
	}

	if command == "blacklistall" {
		isBlacklistAllMode = !isBlacklistAllMode
		text := fmt.Sprintf(localization.GetString("private_command_blacklistall_status"), isBlacklistAllMode)
		_ = adminbot.SendPlainTextMessage(msg.Chat.ID, text, false)
		return
	}

	if command == "emergencymode" {
		toggleEmergencyMode(msg)
		return
	}

	/* HANDLE COMMANDS */
	var text string

	if msg.CommandArguments() != "" {
		switch {
		case strings.HasPrefix(command, "ban"):
			text = handleBanCommands(command, msg)
		case strings.HasPrefix(command, "pardon"):
			text = handlePardonCommands(command, msg)
		case strings.HasPrefix(command, "whitelist"):
			text = handleWhitelistCommands(command, msg)
		case strings.HasPrefix(command, "remove"):
			text = handleRemoveCommands(command, msg)
		default:
			handleOtherCommands(command, msg)
		}

		if text != "" {
			_ = adminbot.SendPlainTextMessageWithMarkup(msg.Chat.ID, text, tgbotapi.NewRemoveKeyboard(true), false)
		}
	}
}

func handleBanCommands(command string, msg *tgbotapi.Message) (reply string) {

	var err error

	switch command {
	case "banpack":
		_, err = database.BlacklistStickerPack(msg.CommandArguments(), msg.From.ID)
	case "banmedia":
		_, err = database.BlacklistMedia(msg.CommandArguments(), msg.CommandArguments(), msg.From.ID)
	case "banhandle":
		_, err = database.BlacklistSource(0, msg.CommandArguments(), msg.From.ID)
	case "banhost":
		setupHostnameBlacklist(msg)
		return
	default:
		return localization.GetString("feature_unimplemented")
	}

	if err != nil {
		reply = localization.GetString("private_command_unable_to_add_to_blacklist")
	} else {
		reply = localization.GetString("private_command_added_to_blacklist")
	}

	return reply
}

func handlePardonCommands(command string, msg *tgbotapi.Message) (reply string) {
	var err error

	switch command {
	case "pardonpack":
		err = database.PardonStickerPack(msg.CommandArguments())
	case "pardonmedia":
		err = database.RemoveMedia(msg.CommandArguments(), msg.CommandArguments())
	case "pardonhandle":
		err = database.RemoveSource(0, msg.CommandArguments())
	case "pardonhost":
		err = database.PardonHostName(msg.CommandArguments())
	default:
		return localization.GetString("feature_unimplemented")
	}

	if err != nil {
		reply = localization.GetString("private_command_unable_to_remove_from_blacklist")
	} else {
		reply = localization.GetString("private_command_removed_from_blacklist")
	}

	return reply
}

func handleWhitelistCommands(command string, msg *tgbotapi.Message) (reply string) {

	var err error

	switch command {
	case "whitelistmedia":
		_, err = database.WhitelistMedia(msg.CommandArguments(), msg.CommandArguments(), msg.From.ID)
	case "whitelistchannel":
		_, err = database.WhitelistSource(0, msg.CommandArguments(), msg.From.ID)
	default:
		return localization.GetString("feature_unimplemented")
	}

	if err != nil {
		reply = localization.GetString("private_command_unable_to_whitelist")
	} else {
		reply = localization.GetString("private_command_added_to_whitelist")
	}

	return reply
}

func handleRemoveCommands(command string, msg *tgbotapi.Message) (reply string) {

	var err error

	switch command {
	case "removechannel":
		err = database.RemoveSource(0, msg.CommandArguments())
	default:
		return localization.GetString("feature_unimplemented")
	}

	if err != nil {
		reply = localization.GetString("private_command_unable_to_remove")
	} else {
		reply = localization.GetString("private_command_removed_successfully")
	}

	return reply
}

func handleOtherCommands(command string, inputMessage *tgbotapi.Message) {

	switch command {
	case "mod":
		setupModAddition(inputMessage)
	case "check":
		attemptCheckByText(inputMessage)
	default:
		_ = adminbot.SendPlainTextMessage(inputMessage.Chat.ID, localization.GetString("feature_unimplemented"), false)
	}
}

func toggleEmergencyMode(msg *tgbotapi.Message) {

	if emergencymode.IsEmergency() {
		emergencymode.End()
		_ = adminbot.SendPlainTextMessage(repository.GetTelegramConfiguration().ReportChannelID, localization.GetString("private_emergencymode_cancelled"), false)
		return
	}

	var emergencyModeDuration time.Duration
	var err error

	/* THE USER SPECIFIED THE DURATION */
	if msg.CommandArguments() != "" {
		emergencyModeDuration, err = time.ParseDuration(msg.CommandArguments())
	}

	/* USE DEFAULT IF NO DURATION SPECIFIED OR ERROR WHEN PARSING IT */
	if msg.CommandArguments() == "" || err != nil {
		emergencyModeDuration = getAppropriateEmergencyDuration(repository.GetTestingStatus())
	}

	emergencyModeText := fmt.Sprintf(localization.GetString("emergencymode_active"), emergencyModeDuration)
	emergencyUpdateChannel := emergencymode.Start()

	err = adminbot.SendPlainTextMessage(repository.GetTelegramConfiguration().ReportChannelID, emergencyModeText, false)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to send the emergency mode message on the report channel: %s", err.Error()))
	}

	_ = adminbot.SendPlainTextMessage(msg.Chat.ID, emergencyModeText, false)
	log.Warn(fmt.Sprintf("%s TRIGGERED EMERGENCY MODE", telegram.GetNameOrUsername(msg.From)))

	go func() {
		emergencyExpirationTimer := time.NewTimer(emergencyModeDuration)

		var emergencyEndedText string
		select {
		case <-emergencyExpirationTimer.C:
			emergencymode.End()
			emergencyEndedText = localization.GetString("emergencymode_expired")
		case <-emergencyUpdateChannel:
			emergencyEndedText = localization.GetString("emergencymode_cancelled")
			emergencyExpirationTimer.Stop()
		}

		log.Info(emergencyEndedText)
		_ = adminbot.SendPlainTextMessage(repository.GetTelegramConfiguration().ReportChannelID, emergencyEndedText, false)
	}()

}
