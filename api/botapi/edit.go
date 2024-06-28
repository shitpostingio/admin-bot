package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// EditMessageText edits a text message using the bot API.
func EditMessageText(messageID int, chatID int64, text, parseMode string, replyMarkup *tgbotapi.InlineKeyboardMarkup) (*tgbotapi.Message, error) {

	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: replyMarkup,
		},
		Text:      text,
		ParseMode: parseMode,
	}

	msg, err := bot.Send(edit)
	return &msg, err
}

// EditMessageReplyMarkup edits the reply markup in a message using the bot API.
func EditMessageReplyMarkup(messageID int, chatID int64, replyMarkup *tgbotapi.InlineKeyboardMarkup) (*tgbotapi.Message, error) {

	edit := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: replyMarkup,
		},
	}

	msg, err := bot.Send(edit)
	return &msg, err
}
