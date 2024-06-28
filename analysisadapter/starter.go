package analysisadapter

import (
	"github.com/shitpostingio/admin-bot/config/structs"
)

var (
	botToken string
	cfg      *structs.FPServerConfiguration
)

// Start starts the fpserver adapter
func Start(telegramBotToken string, fpServerConfig *structs.FPServerConfiguration) {
	botToken = telegramBotToken
	cfg = fpServerConfig
}
