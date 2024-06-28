package adminbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api"
)

// EditMessageText edits a text message in a rate limited fashion.
func EditMessageText(messageID int, chatID int64, text, parseMode string, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	_, err := api.EditMessageText(messageID, chatID, text, parseMode, replyMarkup)
	return err
}

// EditMessageReplyMarkup edits the reply markup in a message in a rate limited fashion.
func EditMessageReplyMarkup(messageID int, chatID int64, replyMarkup *tgbotapi.InlineKeyboardMarkup) (*tgbotapi.Message, error) {
	return api.EditMessageReplyMarkup(messageID, chatID, replyMarkup)
}
