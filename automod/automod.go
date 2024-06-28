package automod

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/reports"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/commands"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/defense/antiflood"
	"github.com/shitpostingio/admin-bot/defense/antispam"
	"github.com/shitpostingio/admin-bot/defense/antiuserbot"
	"github.com/shitpostingio/admin-bot/defense/emergencymode"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
	"github.com/shitpostingio/admin-bot/utility"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//StartDefensiveRoutines starts defensive mechanisms.
func StartDefensiveRoutines() {
	antiflood.Start()
	antispam.Start()
	antiuserbot.Start()
}

// HandleText executes commands and performs checks
// on the text of the message.
func HandleText(msg *tgbotapi.Message) {

	if msg.SenderChat != nil && msg.SenderChat.IsChannel() {

		if !database.SourceIsWhitelisted(msg.From.ID, msg.SenderChat.UserName) {
			err := api.BanChannel(msg.SenderChat.ID, repository.Bot.Self.ID, msg.Chat.ID)
			if err != nil {
				log.Debugf("Unable to ban channel %s: %s", msg.SenderChat.UserName, err)
			}
			adminbot.DeleteMessageAndLog(fmt.Sprintf("Banned channel %s", msg.SenderChat.UserName), msg.Chat.ID, msg.MessageID)
		}

		return
	}

	if msg.IsCommand() {
		HandleCommand(msg)
	}

	performChecksOnMessageText(msg)
}

// HandleMedia handles media messages and checks their content accordingly.
func HandleMedia(msg *tgbotapi.Message, checkNSFW bool) {

	if msg.SenderChat != nil && msg.SenderChat.IsChannel() {

		if !database.SourceIsWhitelisted(msg.From.ID, msg.SenderChat.UserName) {
			err := api.BanChannel(msg.SenderChat.ID, repository.Bot.Self.ID, msg.Chat.ID)
			if err != nil {
				log.Debugf("Unable to ban channel %s: %s", msg.SenderChat.UserName, err)
			}
			adminbot.DeleteMessageAndLog(fmt.Sprintf("Banned channel %s", msg.SenderChat.UserName), msg.Chat.ID, msg.MessageID)
		}

		return
	}

	if performChecksOnMessageText(msg) {
		return
	}

	uniqueID, fileID := telegram.GetFileIDFromMessage(msg)
	performAnalysis(uniqueID, fileID, msg)
}

// HandleSticker checks if a sticker has been forwarded from an unwanted
// source or if it's blacklisted or part of a blacklisted pack.
func HandleSticker(msg *tgbotapi.Message) {

	if msg.SenderChat != nil && msg.SenderChat.IsChannel() {

		if !database.SourceIsWhitelisted(msg.From.ID, msg.SenderChat.UserName) {
			err := api.BanChannel(msg.SenderChat.ID, repository.Bot.Self.ID, msg.Chat.ID)
			if err != nil {
				log.Debugf("Unable to ban channel %s: %s", msg.SenderChat.UserName, err)
			}
			adminbot.DeleteMessageAndLog(fmt.Sprintf("Banned channel %s", msg.SenderChat.UserName), msg.Chat.ID, msg.MessageID)
		}

		return
	}

	if checkMessageOrigin(msg) {
		return
	}

	stickerPackIsBlacklisted := database.StickerPackIsBlacklisted(msg.Sticker.SetName)
	stickerIsBlacklisted := database.MediaIsBlacklisted(msg.Sticker.FileUniqueID, msg.Sticker.FileID)
	if stickerPackIsBlacklisted || stickerIsBlacklisted {
		logText := fmt.Sprintf("Removed a blacklisted sticker posted by %s", telegram.GetNameOrUsername(msg.From))
		adminbot.DeleteMessageAndLog(logText, msg.Chat.ID, msg.MessageID)
	}
}

// HandleGame removes games
func HandleGame(msg *tgbotapi.Message) {

	if msg.SenderChat != nil && msg.SenderChat.IsChannel() {

		if !database.SourceIsWhitelisted(msg.From.ID, msg.SenderChat.UserName) {
			err := api.BanChannel(msg.SenderChat.ID, repository.Bot.Self.ID, msg.Chat.ID)
			if err != nil {
				log.Debugf("Unable to ban channel %s: %s", msg.SenderChat.UserName, err)
			}
			adminbot.DeleteMessageAndLog(fmt.Sprintf("Banned channel %s", msg.SenderChat.UserName), msg.Chat.ID, msg.MessageID)
		}

		return
	}

	logText := fmt.Sprintf("Removed a game posted by %s", telegram.GetNameOrUsername(msg.From))
	adminbot.DeleteMessageAndLog(logText, msg.Chat.ID, msg.MessageID)
}

// HandleCommand removes commands that aren't replies and passes the other
// ones to the command execution.
func HandleCommand(msg *tgbotapi.Message) {

	if !utility.IsChatAdminByMessage(msg) {

		if msg.ReplyToMessage == nil {
			logText := fmt.Sprintf("Removed a command by %s", telegram.GetNameOrUsername(msg.From))
			adminbot.DeleteMessageAndLog(logText, msg.Chat.ID, msg.MessageID)
		}

		return
	}

	commands.ExecuteCommand(msg)

}

// HandleNewChatMember promotes database admins automatically and forwards
// every join to the antiuserbot module.
// In case the EmergencyMode is toggled, it'll also mute every new member
// without any picture or handle.
func HandleNewChatMember(msg *tgbotapi.Message) {

	user := msg.NewChatMembers[len(msg.NewChatMembers)-1]
	if repository.Admins[user.ID] {
		adminbot.PromoteToAdmin(user.ID, msg.Chat.ID)
		_ = database.UpdateModeratorDetailsByTelegramUser(&user, msg.Chat.ID, repository.Bot)
		return
	}

	// Don't send admin joins to the antiuserbot module.
	antiuserbot.HandleUser(&user)

	// Bots definitely can't press buttons
	if user.IsBot {
		return
	}

	// During emergency mode moderators need to approve the user
	if emergencymode.IsEmergency() {
		emergencymode.RestrictUserForEmergency(&user, msg)
		return
	}

	// Sleep a bit to make sure eventual consistency has kicked in
	time.Sleep(200 * time.Millisecond)

	// Check if the user was already restricted
	chatMember, err := adminbot.GetChatMember(user.ID, msg.Chat.ID)
	//chatMember, err := tdlib.GetChatMember(msg.Chat.ID, user.ID)

	// If the user had "serious" restrictions, don't give them a chance to unrestrict themselves
	if err != nil {
		log.Error("Unable to get chat member stats for user with ID ", user.ID, ":", err)
	} else if chatMember.Status == "restricted" {
		log.Info("User with ID", user.ID, " was marked as restricted!")
		return
	}

	//if err == nil && chatMember.Status.ChatMemberStatusType() == "chatMemberStatusRestricted" {
	//	log.Info("User with ID", user.ID, " was marked as restricted!")
	//	return
	//}

	/* RESTRICT THE USER AND SHOW THE RULES */
	_ = adminbot.RestrictMessages(&user, msg.Chat.ID, 0, &repository.Bot.Self)
	verificationMessageRes, err := api.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, repository.Configuration.AdminBot.WelcomeText, false)
	if err != nil {
		return
	}

	// Give the user time to read the rules
	time.Sleep(10 * time.Second)

	if antiuserbot.IsAttack() {
		return
	}

	replyMarkup := buttons.CreateKeyboardWithOneRow(buttons.CreateHumanVerificationButton(msg.From.ID + repository.Bot.Self.ID))
	_, err = adminbot.EditMessageReplyMarkup(verificationMessageRes.MessageID, verificationMessageRes.Chat.ID, &replyMarkup)

	go func() {

		if repository.GetTestingStatus() {
			time.Sleep(15 * time.Second)
		} else {
			time.Sleep(5 * time.Minute)
		}

		deleteErr := api.DeleteMessage(verificationMessageRes.Chat.ID, verificationMessageRes.MessageID)
		if deleteErr == nil {
			adminbot.KickUser(user.ID, repository.Bot.Self.ID, msg.Chat.ID)
		}

	}()

}

//HandleAtAdmin handles @admin mentions
func HandleAtAdmin(msg *tgbotapi.Message) {

	antiflood.IncreaseFloodCounter(1)
	if chatIsUnderAttack() {
		return
	}

	var err error
	var reportedMessageID, backupMessageID int
	var reportedUserID int64
	var reportedUserName string
	reportMessageID := msg.MessageID

	// We will first try to backup the reported message.
	if msg.ReplyToMessage != nil {

		reportedMessageID = msg.ReplyToMessage.MessageID
		reportedUserID = msg.ReplyToMessage.From.ID
		reportedUserName = telegram.GetName(msg.ReplyToMessage.From)
		backupMessageID, err = adminbot.ForwardMessage(repository.GetTelegramConfiguration().BackupChannelID, msg.Chat.ID, reportedMessageID)
		if err != nil {
			log.Error(fmt.Sprintf("Can't forward reported message: %s", err))
		}

	}

	// If we didn't manage to backup the reported message
	// we will at least try to back up the report.
	if backupMessageID == 0 {
		backupMessageID, err = adminbot.ForwardMessage(repository.GetTelegramConfiguration().BackupChannelID, msg.Chat.ID, reportMessageID)
	}

	var reportText string
	if reportedMessageID != 0 {
		reportText = reports.ChatMessageReported(msg.From.ID, telegram.GetName(msg.From), reportedUserID, reportedUserName)
	} else {
		reportText = reports.ChatMessageReport(msg.From.ID, telegram.GetName(msg.From))
	}

	markup := buttons.CreateAtAdminReportMarkup(reportedMessageID, reportMessageID, backupMessageID, msg.Chat.Type)
	_ = reports.ReportWithMarkup(reportText, markup, reports.URGENT)
	log.Info(reportText)
}
