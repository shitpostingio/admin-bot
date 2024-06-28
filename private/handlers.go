package private

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
)

var isBlacklistAllMode bool

func init() {
	isBlacklistAllMode = false
}

// HandlePrivateChat handles private messages sent to the bot
func HandlePrivateChat(msg *tgbotapi.Message) { //nolint:gocyclo

	if !repository.Admins[msg.From.ID] {
		return
	}

	switch {
	case msg.Photo != nil:
		HandlePrivateMedia(msg)
	case msg.Video != nil:
		HandlePrivateMedia(msg)
	case msg.Animation != nil:
		HandlePrivateMedia(msg)
	case msg.Document != nil:
		HandlePrivateMedia(msg)
	case msg.Voice != nil:
		HandlePrivateMedia(msg)
	case msg.Audio != nil:
		HandlePrivateMedia(msg)
	case msg.VideoNote != nil:
		HandlePrivateMedia(msg)
	case msg.Sticker != nil:
		HandlePrivateSticker(msg)
	case msg.IsCommand():
		HandlePrivateCommand(msg)
	case msg.Text != "":
		HandlePrivateText(msg)
	}
}

//HandlePrivateText allows admins to blacklist or pardon handles via private messages
func HandlePrivateText(msg *tgbotapi.Message) {

	if msg.ForwardSenderName != "" {
		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, localization.GetString("private_checks_profile_hidden_from_forwards"), false)
		return
	}

	switch {
	case msg.ForwardFrom != nil:
		checkUserStatus(msg, msg.ForwardFrom)
	default:
		checkForHandles(msg)
	}

}

// HandlePrivateMedia handles media sent to the bot in private
func HandlePrivateMedia(msg *tgbotapi.Message) {

	var text string
	var markup tgbotapi.InlineKeyboardMarkup

	uniqueID, fileID := telegram.GetFileIDFromMessage(msg)
	media, err := database.FindMediaByFileID(uniqueID, fileID)

	if isBlacklistAllMode {

		// blacklist it in case of error, better safe than sorry!
		if err != nil || media.IsWhitelisted {

			_, err := database.BlacklistMedia(uniqueID, fileID, msg.From.ID)
			if err != nil {
				text = fmt.Sprintf(localization.GetString("private_blacklistall_media_unable_to_add"), localization.GetString("blacklistall_active"))
			} else {
				text = fmt.Sprintf(localization.GetString("private_blacklistall_media_added_correctly"), localization.GetString("blacklistall_active"))
			}

		} else {

			text = fmt.Sprintf(localization.GetString("private_blacklistall_media_already_blacklisted"), localization.GetString("blacklistall_active"))

		}

		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, text, false)

	} else {

		text = localization.GetString("private_media_what_to_do")
		if err == nil {
			fileID = media.FileID
		}

		if err != nil || media.IsWhitelisted {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistMediaButton())
		} else {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonMediaButton(), buttons.CreatePrivateWhitelistMediaButton())
		}

		_ = adminbot.SendReplyTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, text, markup, false)

	}

}

//HandlePrivateSticker allows admins to blacklist or pardon stickers/sticker packs via private messages
func HandlePrivateSticker(msg *tgbotapi.Message) {

	var text string

	// In BlacklistAllMode we will blacklist the sticker pack
	if isBlacklistAllMode {

		if !database.StickerPackIsBlacklisted(msg.Sticker.SetName) {

			_, err := database.BlacklistStickerPack(msg.Sticker.SetName, msg.From.ID)
			if err != nil {
				text = fmt.Sprintf(localization.GetString("private_blacklistall_stickerpack_unable_to_add"), localization.GetString("blacklistall_active"))
			} else {
				text = fmt.Sprintf(localization.GetString("private_blacklistall_stickerpack_added_correctly"), localization.GetString("blacklistall_active"))
			}

		} else {

			text = fmt.Sprintf(localization.GetString("private_blacklistall_stickerpack_already_blacklisted"), localization.GetString("blacklistall_active"))

		}

		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, text, false)

	} else {

		text = localization.GetString("private_sticker_what_to_do")

		var markup tgbotapi.InlineKeyboardMarkup
		var stickerPackAction tgbotapi.InlineKeyboardButton

		// If the sticker belongs to no sticker pack we won't show the button
		if msg.Sticker.SetName != "" {

			if !database.StickerPackIsBlacklisted(msg.Sticker.SetName) {
				stickerPackAction = buttons.CreatePrivateBlacklistStickerPackButton()
			} else {
				stickerPackAction = buttons.CreatePrivatePardonStickerPackButton()
			}

		}

		media, err := database.FindMediaByFileID(msg.Sticker.FileUniqueID, msg.Sticker.FileID)
		if err != nil || media.IsWhitelisted {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistStickerButton(), stickerPackAction)
		} else {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonStickerButton(), buttons.CreatePrivateWhitelistStickerButton(), stickerPackAction)
		}

		_ = adminbot.SendReplyPlainTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, text, markup, false)

	}

}
