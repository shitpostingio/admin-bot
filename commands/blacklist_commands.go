package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
)

// blacklistMedia removes a media and adds it to the blacklist, if not already present
func blacklistMedia(msg *tgbotapi.Message) {

	adminbot.DeleteMultipleMessages(msg.Chat.ID, msg.MessageID, msg.ReplyToMessage.MessageID)

	// TODO: Temporary fix against crashing with animated stickers
	if msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.IsAnimated {
		blacklistStickerPack(msg)
		return
	}

	uniqueID, fileID := telegram.GetFileIDFromMessage(msg.ReplyToMessage)
	if fileID != "" {
		_, _ = database.BlacklistMedia(uniqueID, fileID, msg.From.ID)
	}

}

// blacklistStickerPack removes a sticker and adds its sticker pack to the blacklist, if not already present
func blacklistStickerPack(msg *tgbotapi.Message) {

	adminbot.DeleteMultipleMessages(msg.Chat.ID, msg.MessageID, msg.ReplyToMessage.MessageID)

	// Blacklist the the media as well, just to be safe
	// TODO: Temporary fix against crashing with animated stickers
	//blacklistMedia(msg)

	// A moderator may have used the blacklist sticker
	// command on something that is not a sticker.
	// We have already blacklisted the media, we can return.
	if msg.ReplyToMessage.Sticker == nil {
		return
	}

	// Some stickers may not have a sticker pack
	if msg.ReplyToMessage.Sticker.SetName == "" {
		return
	}

	if !database.StickerPackIsBlacklisted(msg.ReplyToMessage.Sticker.SetName) {
		_, _ = database.BlacklistStickerPack(msg.ReplyToMessage.Sticker.SetName, msg.From.ID)
	}

}

//blacklistHandle blacklists all handles in a message.
func blacklistHandle(msg *tgbotapi.Message) {

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)

	// The command is restricted to admins only
	if !repository.Admins[msg.From.ID] {
		return
	}

	adminbot.DeleteMessage(msg.ReplyToMessage.Chat.ID, msg.ReplyToMessage.MessageID)

	messageEntities := telegram.GetMessageEntities(msg.ReplyToMessage)
	text := telegram.GetMessageText(msg.ReplyToMessage)
	handles := telegram.GetAllMentions(text, messageEntities, msg.ReplyMarkup)

	//Try to blacklist handles in the message first.
	//If there are none, fallback to blacklisting the
	//user from whom the message was forwarded.
	if len(handles) != 0 {

		for _, handle := range handles {

			// Don't blacklist handles belonging to moderators
			// or those that have been whitelisted.
			if !database.IsModeratorUsername(handle) && !database.SourceIsWhitelisted(0, handle) {
				_, _ = database.BlacklistSource(0, handle, msg.From.ID)
			}

		}

		return
	}

	// Check if there's a forward handle to blacklist
	if msg.ReplyToMessage.ForwardFrom != nil && msg.ReplyToMessage.ForwardFrom.UserName != "" {

		// Don't blacklist handles belonging to moderators
		// or those that have been whitelisted.
		forwardHandle := msg.ReplyToMessage.ForwardFrom.UserName
		if !database.IsModeratorUsername(forwardHandle) && !database.SourceIsWhitelisted(0, forwardHandle) {
			_, _ = database.BlacklistSource(0, forwardHandle, msg.From.ID)
		}

	}
}
