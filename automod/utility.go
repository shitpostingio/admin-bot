package automod

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/defense/antiflood"
	"github.com/shitpostingio/admin-bot/defense/antiuserbot"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/utility"
)

//reportNSFWMessage backs up the NSFW message and reports it.
func reportNSFWMessage(chatID int64, messageID int, reportMessage string, nsfwTableID string, chatType string) {

	fwMsgID, err := adminbot.ForwardMessage(repository.GetTelegramConfiguration().BackupChannelID, chatID, messageID)

	antiflood.IncreaseFloodCounter(1)
	if !antiflood.IsFlood() {

		if err == nil {
			backupRow := tgbotapi.NewInlineKeyboardRow(buttons.CreateBackupMessageButton(fwMsgID, chatType))
			actionRow := tgbotapi.NewInlineKeyboardRow(buttons.CreateWhitelistMediaButton(nsfwTableID), buttons.CreateHandleButton())
			markup := tgbotapi.NewInlineKeyboardMarkup(backupRow, actionRow)
			_ = adminbot.SendTextMessageWithMarkup(repository.GetTelegramConfiguration().ReportChannelID, reportMessage, markup, false)
		} else {
			_ = adminbot.SendTextMessage(repository.GetTelegramConfiguration().ReportChannelID, reportMessage, false)
		}
	}
}

//handleMessageDeletionAndReport deletes a message and reports it if possible
func handleMessageDeletionAndReport(reportText string, msg *tgbotapi.Message) {

	antiflood.IncreaseFloodCounter(1)
	if !antiflood.IsFlood() {
		adminbot.DeleteMessageAndReport(reportText, msg.Chat.ID, msg.MessageID)
		return
	}

	adminbot.DeleteMessageAndLog(reportText, msg.Chat.ID, msg.MessageID)
	_ = adminbot.RestrictMessages(msg.From, msg.Chat.ID, utility.GetAppropriateRestrictionEnd(), &repository.Bot.Self)
}

//chatIsUnderAttack returns true if the chat is being flooded or has suspicious joins
func chatIsUnderAttack() bool {
	return antiflood.IsFlood() || antiuserbot.IsAttack()
}
