package callback

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
)

func blacklistSource(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	_, err := database.BlacklistSource(0, callbackFields[1], callbackQuery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackQuery.ID, "The operation wasn't successful, please retry")
		return
	}

	markup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonSourceButton(callbackFields[1]), buttons.CreatePrivateWhitelistSourceButton(callbackFields[1]))
	_, _ = adminbot.EditMessageReplyMarkup(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, &markup)

}

func pardonSource(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	err := database.RemoveSource(0, callbackFields[1])
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackQuery.ID, "The operation wasn't successful, please retry")
		return
	}

	markup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistSourceButton(callbackFields[1]), buttons.CreatePrivateWhitelistSourceButton(callbackFields[1]))
	_, _ = adminbot.EditMessageReplyMarkup(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, &markup)

}

func whitelistSource(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	_, err := database.WhitelistSource(0, callbackFields[1], callbackQuery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackQuery.ID, "The operation wasn't successful, please retry")
		return
	}

	markup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistSourceButton(callbackFields[1]))
	_, _ = adminbot.EditMessageReplyMarkup(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, &markup)
}
