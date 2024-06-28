package adminbot

import (
	"github.com/shitpostingio/admin-bot/api"
)

/* ------------------------------ NORMAL MESSAGES ------------------------------ */

// SendPlainTextMessage sends a plain text message.
func SendPlainTextMessage(chatID int64, text string, urgent bool) error {
	_, err := api.SendPlainTextMessage(chatID, text, urgent)
	return err
}

// SendPlainTextMessageWithMarkup sends a text message with an inline keyboard.
func SendPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendPlainTextMessageWithMarkup(chatID, text, replyMarkup, urgent)
	return err
}

// SendTextMessage sends a text message with HTML parse mode.
func SendTextMessage(chatID int64, text string, urgent bool) error {
	_, err := api.SendTextMessage(chatID, text, urgent)
	return err
}

// SendTextMessageWithMarkup sends a text message with HTML parse mode and an inline keyboard.
func SendTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendTextMessageWithMarkup(chatID, text, replyMarkup, urgent)
	return err
}

/* ------------------------------ REPLY MESSAGES ------------------------------ */

// SendReplyPlainTextMessage sends a plain text reply message.
func SendReplyPlainTextMessage(replyToMessageID int, chatID int64, text string, urgent bool) error {
	_, err := api.SendReplyPlainTextMessage(replyToMessageID, chatID, text, urgent)
	return err
}

// SendReplyPlainTextMessageWithMarkup sends a text message with markup and an inline keyboard.
func SendReplyPlainTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendReplyPlainTextMessageWithMarkup(replyToMessageID, chatID, text, replyMarkup, urgent)
	return err
}

// SendReplyTextMessage sends a reply text message with HTML parse mode.
func SendReplyTextMessage(replyToMessageID int, chatID int64, text string, urgent bool) error {
	_, err := api.SendReplyTextMessage(replyToMessageID, chatID, text, urgent)
	return err
}

// SendReplyTextMessageWithMarkup sends a reply text message with HTML parse mode and an inline keyboard.
func SendReplyTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendReplyTextMessageWithMarkup(replyToMessageID, chatID, text, replyMarkup, urgent)
	return err
}

///* ------------------------------ SILENT MESSAGES ------------------------------ */

// SendSilentPlainTextMessage sends a plain text message with notifications disabled.
func SendSilentPlainTextMessage(chatID int64, text string, urgent bool) error {
	_, err := api.SendSilentPlainTextMessage(chatID, text, urgent)
	return err
}

// SendSilentPlainTextMessageWithMarkup sends a text message with an inline keyboard and notifications disabled.
func SendSilentPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendSilentPlainTextMessageWithMarkup(chatID, text, replyMarkup, urgent)
	return err
}

// SendSilentTextMessage sends a text message with HTML parse mode and notifications disabled.
func SendSilentTextMessage(chatID int64, text string, urgent bool) error {
	_, err := api.SendSilentTextMessage(chatID, text, urgent)
	return err
}

// SendSilentTextMessageWithMarkup sends a text message with HTML parse mode and an inline keyboard and notifications disabled.
func SendSilentTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) error {
	_, err := api.SendSilentTextMessageWithMarkup(chatID, text, replyMarkup, urgent)
	return err
}
