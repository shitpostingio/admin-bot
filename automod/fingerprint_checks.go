package automod

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/telegram"
)

func performFingerprintChecks(uniqueFileID, fileID string, msg *tgbotapi.Message) bool {

	media, err := database.FindMediaByFileID(uniqueFileID, fileID)
	if err != nil {
		return false
	}

	if !media.IsWhitelisted {
		adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
		log.Info(fmt.Sprintf("Removed a blacklisted media posted by %s", telegram.GetNameOrUsername(msg.From)))
	}

	return true
}
