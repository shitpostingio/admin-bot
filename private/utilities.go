package private

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkForHandles(msg *tgbotapi.Message) {

	/* FIND ALL HANDLES */
	handles := telegram.GetAllMentions(msg.Text, telegram.GetMessageEntities(msg), msg.ReplyMarkup)

	/* WE ARE JUST CHECKING FOR HANDLES */
	if len(handles) > 0 {

		/* BLACKLIST EVERYTHING THAT IS FOUND WHEN IN BLACKLISTALL MODE */
		if isBlacklistAllMode {

			for _, handle := range handles {

				if !database.SourceIsBlacklisted(0, handle) {
					_, _ = database.BlacklistSource(0, handle, msg.From.ID)
				}
			}

			text := fmt.Sprintf(localization.GetString("private_blacklistall_handles_added_correctly"), localization.GetString("private_blacklistall_active"))
			_ = adminbot.SendPlainTextMessage(msg.Chat.ID, text, false)
			return
		}

		var markup tgbotapi.InlineKeyboardMarkup
		if database.SourceIsBlacklisted(0, handles[0]) {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivatePardonSourceButton(handles[0]),
				buttons.CreatePrivateWhitelistSourceButton(handles[0]))
		} else {
			markup = buttons.CreateKeyboardWithOneRow(buttons.CreatePrivateBlacklistSourceButton(handles[0]),
				buttons.CreatePrivateWhitelistSourceButton(handles[0]))
		}

		text := fmt.Sprintf("What should we do with the source %s?", handles[0])
		_ = adminbot.SendReplyTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, text, markup, false)

	}
}

func setupModAddition(inputMessage *tgbotapi.Message) {

	userID, err := strconv.ParseInt(inputMessage.CommandArguments(), 10, 64)
	if err != nil {
		_ = adminbot.SendPlainTextMessage(inputMessage.Chat.ID, localization.GetString("private_mods_unable_to_parse_userid"), false)
		return
	}

	chatMember, err := adminbot.GetChatMember(userID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		_ = adminbot.SendPlainTextMessage(inputMessage.Chat.ID, localization.GetString("private_mods_unable_to_find_user"), false)
		return
	}

	text := fmt.Sprintf(localization.GetString("private_mods_are_you_sure"), chatMember.User.ID, telegram.GetName(chatMember.User))
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateModUserButton(userID))
	_ = adminbot.SendTextMessageWithMarkup(inputMessage.Chat.ID, text, markup, false)

}

func setupHostnameBlacklist(inputMessage *tgbotapi.Message) {

	text := localization.GetString("private_hosts_what_to_do")

	/* CREATE INLINE KEYBOARD */
	upperRow := tgbotapi.NewInlineKeyboardRow(buttons.CreateBlacklistBanworthyHostnameButton(inputMessage.CommandArguments()),
		buttons.CreateBlacklistTelegramHostnameButton(inputMessage.CommandArguments()))
	lowerRow := tgbotapi.NewInlineKeyboardRow(buttons.CreateBlacklistHostnameButton(inputMessage.CommandArguments()))
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(upperRow, lowerRow)

	_ = adminbot.SendReplyPlainTextMessageWithMarkup(inputMessage.MessageID, inputMessage.Chat.ID, text, replyMarkup, false)

}

func getAppropriateEmergencyDuration(testing bool) time.Duration {

	if testing {
		return time.Minute * 1
	}

	return time.Hour * 24
}

func createOneTimeReplyKeyboardWithCancelOption(buttonArgument, buttonContent string, createBanButton, createPardonButton, createWhitelistButton bool) tgbotapi.ReplyKeyboardMarkup {

	cancelRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/cancel"))

	if createBanButton {
		banButtonText := fmt.Sprintf("/ban%s %s", strings.Title(buttonArgument), buttonContent)
		banButtonRow := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(banButtonText))
		return tgbotapi.NewReplyKeyboard(banButtonRow, cancelRow)
	}

	var actionRow []tgbotapi.KeyboardButton
	if createPardonButton {
		pardonButtonText := fmt.Sprintf("/pardon%s %s", strings.Title(buttonArgument), buttonContent)
		actionRow = append(actionRow, tgbotapi.NewKeyboardButton(pardonButtonText))
	}

	if createWhitelistButton {
		whitelistButtonText := fmt.Sprintf("/whitelist%s %s", strings.Title(buttonArgument), buttonContent)
		actionRow = append(actionRow, tgbotapi.NewKeyboardButton(whitelistButtonText))
	}

	keyboard := tgbotapi.NewReplyKeyboard(actionRow, cancelRow)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
