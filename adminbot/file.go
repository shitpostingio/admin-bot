package adminbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api"
)

// GetTelegramFile gets a telegram file given its file ID in a rate limited fashion.
func GetTelegramFile(uniqueFileID, fileID string) (*tgbotapi.File, error) {
	return api.GetTelegramFile(uniqueFileID, fileID)
}
