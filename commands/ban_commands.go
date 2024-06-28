package commands

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
)

// banUser bans a user in a supergroup. Telegram command is /ban
func banUser(msg *tgbotapi.Message) {

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
	err := adminbot.BanUser(msg.ReplyToMessage.From, msg.From, msg.CommandArguments(), msg.Chat.ID)
	if err != nil {
		log.Error("banUser: ", err)
	}

}

func banByUsername(msg *tgbotapi.Message) {

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)

	mentions := telegram.GetMentions(utf16.Encode([]rune(msg.Text)), telegram.GetMessageEntities(msg))
	if len(mentions) == 0 {
		log.Error("banByUsername: attempt to ban by handle with no handle")
		return
	}

	username := mentions[0]
	reason := strings.ReplaceAll(strings.ToLower(msg.CommandArguments()), "@"+username, "")
	_, err := adminbot.BanUserByUsername(username, msg.From, reason, msg.Chat.ID)
	if err != nil {
		log.Error("banByUsername: ", err)
	}

}

func banByID(msg *tgbotapi.Message) {

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)

	if msg.CommandArguments() == "" {
		log.Error("banByID: no id or reason")
		return
	}

	words := strings.Fields(msg.CommandArguments())
	bannedUserID, err := strconv.ParseInt(words[0], 10, 64)
	if err != nil {
		log.Error("banByID: couldn't parse user ID", words[0])
		return
	}

	reason := strings.ReplaceAll(msg.CommandArguments(), words[0], "")
	_, err = adminbot.BanUserByID(bannedUserID, msg.From, reason, msg.Chat.ID)
	if err != nil {
		log.Error("banByID: ", err)
	}

}

func reportBan(bannedUser, moderator *tgbotapi.User, reason string, chatID int64) {
	reportText := fmt.Sprintf(localization.GetString("user_banned"), moderator.ID, telegram.GetName(moderator), bannedUser.ID, telegram.GetName(bannedUser), reason)
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnbanButton(bannedUser.ID))
	_ = adminbot.SendTextMessageWithMarkup(repository.GetTelegramConfiguration().ReportChannelID, reportText, markup, true)
}
