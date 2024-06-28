package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// EditMessageText edits a text message.
// It will also use a rate limiter not to get restricted by Telegram.
func EditMessageText(messageID int, chatID int64, text, parseMode string, replyMarkup *tgbotapi.InlineKeyboardMarkup) (*tgbotapi.Message, error) {
	limiter.AuthorizeAction()
	return botapi.EditMessageText(messageID, chatID, text, parseMode, replyMarkup)
}

// EditMessageReplyMarkup edits the reply markup in a message.
// It will also use a rate limiter not to get restricted by Telegram.
func EditMessageReplyMarkup(messageID int, chatID int64, replyMarkup *tgbotapi.InlineKeyboardMarkup) (*tgbotapi.Message, error) {
	limiter.AuthorizeAction()
	return botapi.EditMessageReplyMarkup(messageID, chatID, replyMarkup)
}
