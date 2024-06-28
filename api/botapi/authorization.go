package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

// Authorize logs the bot into the provided account using the bot API.
func Authorize(botToken string, debugFlag bool) (*tgbotapi.BotAPI, error) {

	var err error
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = debugFlag
	return bot, nil
}
