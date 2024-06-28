package automod

import (
	"github.com/shitpostingio/admin-bot/reports"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/telegram"
)

func performNSFWChecks(uniqueFileID, fileID string) (isNSFW bool, nsfwTableID string, description string, score float64) {

	//isNSFW, score, description = analysisadapter.GetNSFWScores(uniqueFileID, fileID)
	//if !isNSFW {
	//	return
	//}
	//
	//nsfwTableID, _ = database.BlacklistNSFWMedia(uniqueFileID, fileID, description, score, repository.Bot.Self.ID)
	return

}

// performNSFWChecksOnPhoto gets the NSFW scores for the photo from FPServer.
// If the photo is NSFW, it'll be blacklisted and added to the nsfw table.
func performNSFWChecksOnMedia(uniqueFileID, fileID string, msg *tgbotapi.Message) (isNSFW bool) {

	isNSFW, nsfwTableID, description, score := performNSFWChecks(uniqueFileID, fileID)
	if nsfwTableID != "" {

		reportText := reports.RemovedNSFWMedia(msg.From.ID, telegram.GetName(msg.From), description, score)
		reportNSFWMessage(msg.Chat.ID, msg.MessageID, reportText, nsfwTableID, msg.Chat.Type)
		adminbot.DeleteMessageAndLog(reportText, msg.Chat.ID, msg.MessageID)

	}

	return isNSFW

}
