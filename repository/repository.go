package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/config/structs"
)

var (

	// Bot represents the bot in the Telegram Bot API.
	Bot *tgbotapi.BotAPI

	// Configuration represents the configuration
	Configuration *structs.Config

	// Mods holds the user IDs of mods
	Mods map[int64]bool

	// Admins holds the user ids of admins
	Admins map[int64]bool

	// Testing represents the testing status of the bot
	Testing bool

	// Debug represents the debug status of the bot
	Debug bool
)
