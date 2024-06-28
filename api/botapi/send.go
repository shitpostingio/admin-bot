package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/consts"
)

/* ------------------------------------------------------------------------------------------------------------------ */

// SendPlainTextMessage sends a plaintext message using the bot API.
func SendPlainTextMessage(chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendPlainTextMessageWithMarkup sends a plaintext message with markup using the bot API.
func SendPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: replyMarkup,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendTextMessage sends a text message with ParseMode enabled using the bot API.
func SendTextMessage(chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendTextMessageWithMarkup sends a text message with markup and ParseMode enabled using the bot API.
func SendTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: replyMarkup,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

/* ------------------------------------------------------------------------------------------------------------------ */

// SendReplyPlainTextMessage sends a plaintext reply message using the bot API.
func SendReplyPlainTextMessage(replyToMessageID int, chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: replyToMessageID,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendReplyPlainTextMessageWithMarkup sends a plaintext reply message with markup using the bot API.
func SendReplyPlainTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: replyToMessageID,
			ReplyMarkup:      replyMarkup,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendReplyTextMessage sends a reply text message with ParseMode enabled using the bot API.
func SendReplyTextMessage(replyToMessageID int, chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: replyToMessageID,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendReplyTextMessageWithMarkup sends a reply text message with markup and ParseMode enabled using the bot API.
func SendReplyTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: replyToMessageID,
			ReplyMarkup:      replyMarkup,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

/* ------------------------------------------------------------------------------------------------------------------ */

// SendSilentPlainTextMessage sends a plaintext message with notifications disabled using the bot API.
func SendSilentPlainTextMessage(chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:              chatID,
			DisableNotification: true,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendSilentPlainTextMessageWithMarkup sends a plaintext message with markup and notifications disabled using the bot API.
func SendSilentPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:              chatID,
			DisableNotification: true,
			ReplyMarkup:         replyMarkup,
		},
		Text: text,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendSilentTextMessage sends a text message with ParseMode enabled and notifications disabled using the bot API.
func SendSilentTextMessage(chatID int64, text string) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:              chatID,
			DisableNotification: true,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}

// SendSilentTextMessageWithMarkup sends a text message with markup, ParseMode enabled and notifications disabled using the bot API.
func SendSilentTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}) (*tgbotapi.Message, error) {

	message := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:              chatID,
			DisableNotification: true,
			ReplyMarkup:         replyMarkup,
		},
		Text:      text,
		ParseMode: consts.ReportParseMode,
	}

	msg, err := bot.Send(message)
	return &msg, err
}
