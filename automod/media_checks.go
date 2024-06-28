package automod

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/analysisadapter"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/entities"
	"github.com/shitpostingio/admin-bot/reports"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
	log "github.com/sirupsen/logrus"
)

func performAnalysis(fileUniqueID, fileID string, msg *tgbotapi.Message) {

	// Every media can be looked up via FileUniqueID
	media, err := database.FindMediaByFileUniqueID(fileUniqueID)
	if err == nil {
		handleKnownMedia(media, msg)
		return
	}

	// We can perform feature-based analysis for only
	// certain file types.
	if !telegram.MediaCanBeAnalyzed(fileID) {
		return
	}

	// Feature based search
	analysis, analysisErr := analysisadapter.GetAnalysis(fileUniqueID, fileID)
	if analysisErr == nil || !errors.Is(analysisErr, analysisadapter.FingerprintError) {
		media, err = database.FindMediaByFeatures(analysis.Fingerprint.Histogram, analysis.Fingerprint.PHash, 0.08)
		if err == nil {
			handleKnownMedia(media, msg)
			return
		}
	}

	// NSFW
	if analysisErr == nil || !errors.Is(analysisErr, analysisadapter.NSFWError) {
		if analysis.NSFW.IsNSFW {

			nsfwTableID, _ := database.BlacklistNSFWMedia(fileUniqueID, fileID, analysis.NSFW.Label, analysis.NSFW.Confidence, repository.Bot.Self.ID)
			reportText := reports.RemovedNSFWMedia(msg.From.ID, telegram.GetName(msg.From), analysis.NSFW.Label, analysis.NSFW.Confidence)
			reportNSFWMessage(msg.Chat.ID, msg.MessageID, reportText, nsfwTableID, msg.Chat.Type)
			adminbot.DeleteMessageAndLog(reportText, msg.Chat.ID, msg.MessageID)

		}
	}
}

func handleKnownMedia(media entities.Media, msg *tgbotapi.Message) {

	if media.IsWhitelisted {
		return
	}

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
	log.Info(fmt.Sprintf("Removed a blacklisted media posted by %s", telegram.GetNameOrUsername(msg.From)))
}
