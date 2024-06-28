package callback

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/telegram"
)

func blacklistPrivateMedia(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	_, err := database.BlacklistMedia(uniqueID, fileID, callbackquery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	newMarkup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonMediaButton(), buttons.CreatePrivateWhitelistMediaButton())
	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func whitelistPrivateMedia(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	_, err := database.WhitelistMedia(uniqueID, fileID, callbackquery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	newMarkup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistMediaButton())
	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func pardonPrivateMedia(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	err := database.RemoveMedia(uniqueID, fileID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	newMarkup := buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistMediaButton())
	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func blacklistPrivateStickerPack(callbackquery *tgbotapi.CallbackQuery) {

	_, err := database.BlacklistStickerPack(callbackquery.Message.ReplyToMessage.Sticker.SetName, callbackquery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	var newMarkup tgbotapi.InlineKeyboardMarkup
	uniqueFileID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	if database.MediaIsBlacklisted(uniqueFileID, fileID) {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonStickerButton(), buttons.CreatePrivateWhitelistStickerButton(), buttons.CreatePrivatePardonStickerPackButton())
	} else {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivatePardonStickerPackButton())
	}

	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func pardonPrivateStickerPack(callbackquery *tgbotapi.CallbackQuery) {

	err := database.PardonStickerPack(callbackquery.Message.ReplyToMessage.Sticker.SetName)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	var newMarkup tgbotapi.InlineKeyboardMarkup
	uniqueFileID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	if database.MediaIsBlacklisted(uniqueFileID, fileID) {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonStickerButton(), buttons.CreatePrivateWhitelistStickerButton(), buttons.CreatePrivateBlacklistStickerPackButton())
	} else {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivateBlacklistStickerPackButton())
	}

	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func blacklistPrivateSticker(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	_, err := database.BlacklistMedia(uniqueID, fileID, callbackquery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	var newMarkup tgbotapi.InlineKeyboardMarkup
	if database.StickerPackIsBlacklisted(callbackquery.Message.ReplyToMessage.Sticker.SetName) {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonStickerButton(), buttons.CreatePrivateWhitelistStickerButton(), buttons.CreatePrivatePardonStickerPackButton())
	} else {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonStickerButton(), buttons.CreatePrivateWhitelistStickerButton(), buttons.CreatePrivateBlacklistStickerPackButton())
	}

	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func whitelistPrivateSticker(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	_, err := database.WhitelistMedia(uniqueID, fileID, callbackquery.From.ID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	var newMarkup tgbotapi.InlineKeyboardMarkup
	if database.StickerPackIsBlacklisted(callbackquery.Message.ReplyToMessage.Sticker.SetName) {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivatePardonStickerPackButton())
	} else {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivateBlacklistStickerPackButton())
	}

	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}

func pardonPrivateSticker(callbackquery *tgbotapi.CallbackQuery) {

	uniqueID, fileID := telegram.GetFileIDFromMessage(callbackquery.Message.ReplyToMessage)
	err := database.RemoveMedia(uniqueID, fileID)
	if err != nil {
		_ = adminbot.SendCallbackWithAlert(callbackquery.ID, "The operation wasn't successful, please retry")
		return
	}

	var newMarkup tgbotapi.InlineKeyboardMarkup
	if database.StickerPackIsBlacklisted(callbackquery.Message.ReplyToMessage.Sticker.SetName) {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivatePardonStickerPackButton())
	} else {
		newMarkup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), buttons.CreatePrivateBlacklistStickerPackButton())
	}

	_, _ = adminbot.EditMessageReplyMarkup(callbackquery.Message.MessageID, callbackquery.Message.Chat.ID, &newMarkup)

}
