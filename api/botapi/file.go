package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetTelegramFile gets a telegram File using the bot API.
func GetTelegramFile(fileID string) (*tgbotapi.File, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	return &file, err
}
