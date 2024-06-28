package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

/* ------------------------------------------------------------------------------------------------------------------ */

// SendPlainTextMessage sends a plaintext message.
// It will also use a rate limiter not to get restricted by Telegram.
func SendPlainTextMessage(chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendPlainTextMessage(chatID, text)
}

// SendPlainTextMessageWithMarkup sends a plaintext message with markup.
// It will also use a rate limiter not to get restricted by Telegram.
func SendPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendPlainTextMessageWithMarkup(chatID, text, replyMarkup)
}

// SendTextMessage sends a text message with ParseMode enabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendTextMessage(chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendTextMessage(chatID, text)
}

// SendTextMessageWithMarkup sends a text message with markup and ParseMode enabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendTextMessageWithMarkup(chatID, text, replyMarkup)
}

/* ------------------------------------------------------------------------------------------------------------------ */

// SendReplyPlainTextMessage sends a plaintext reply message.
// It will also use a rate limiter not to get restricted by Telegram.
func SendReplyPlainTextMessage(replyToMessageID int, chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendReplyPlainTextMessage(replyToMessageID, chatID, text)
}

// SendReplyPlainTextMessageWithMarkup sends a plaintext reply message with markup.
// It will also use a rate limiter not to get restricted by Telegram.
func SendReplyPlainTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendReplyPlainTextMessageWithMarkup(replyToMessageID, chatID, text, replyMarkup)
}

// SendReplyTextMessage sends a reply text message with ParseMode enabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendReplyTextMessage(replyToMessageID int, chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendReplyTextMessage(replyToMessageID, chatID, text)
}

// SendReplyTextMessageWithMarkup sends a reply text message with markup and ParseMode enabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendReplyTextMessageWithMarkup(replyToMessageID int, chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendReplyTextMessageWithMarkup(replyToMessageID, chatID, text, replyMarkup)
}

/* ------------------------------------------------------------------------------------------------------------------ */

// SendSilentPlainTextMessage sends a plaintext message with notifications disabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendSilentPlainTextMessage(chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendSilentPlainTextMessage(chatID, text)
}

// SendSilentPlainTextMessageWithMarkup sends a plaintext message with markup and notifications disabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendSilentPlainTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendSilentPlainTextMessageWithMarkup(chatID, text, replyMarkup)
}

// SendSilentTextMessage sends a text message with ParseMode enabled and notifications disabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendSilentTextMessage(chatID int64, text string, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendSilentTextMessage(chatID, text)
}

// SendSilentTextMessageWithMarkup sends a text message with markup, ParseMode enabled and notifications disabled.
// It will also use a rate limiter not to get restricted by Telegram.
func SendSilentTextMessageWithMarkup(chatID int64, text string, replyMarkup interface{}, urgent bool) (*tgbotapi.Message, error) {
	authorizeWithAppropriatePriority(urgent)
	return botapi.SendSilentTextMessageWithMarkup(chatID, text, replyMarkup)
}

/* ------------------------------------------------------------------------------------------------------------------ */

// authorizeWithAppropriatePriority uses the urgent flag to
// authorize the action with the appropriate priority.
func authorizeWithAppropriatePriority(urgent bool) {
	if urgent {
		limiter.AuthorizeUrgentAction()
	} else {
		limiter.AuthorizeAction()
	}
}
