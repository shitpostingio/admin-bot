package commands

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/utility"
)

//restrictUser restricts a user
func restrictUser(msg *tgbotapi.Message) {

	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)

	// Database admins cannot be restricted.
	// Moderators can be restricted only by admins.
	if repository.Admins[msg.ReplyToMessage.From.ID] {
		return
	}

	if repository.Mods[msg.ReplyToMessage.From.ID] && !repository.Admins[msg.From.ID] {
		return
	}

	command := strings.ToLower(msg.Command())
	restrictionEndTime := utility.UnixTimeIn(msg.CommandArguments())
	switch {
	case isMute(command):
		_ = adminbot.RestrictMessages(msg.ReplyToMessage.From, msg.Chat.ID, restrictionEndTime, msg.From)
	case isNoMedia(command):
		_ = adminbot.RestrictMedia(msg.ReplyToMessage.From, msg.Chat.ID, restrictionEndTime, msg.From)
	case isNoOther(command):
		_ = adminbot.RestrictOther(msg.ReplyToMessage.From, msg.Chat.ID, restrictionEndTime, msg.From)
	}
}

func isMute(command string) bool {

	if command == "mute" ||
		strings.HasPrefix(command, "nomessage") {
		return true
	}

	return false
}

func isNoMedia(command string) bool {

	if strings.HasPrefix(command, "nomedia") ||
		strings.HasPrefix(command, "nopic") {
		return true
	}

	return false
}

func isNoOther(command string) bool {

	if strings.HasPrefix(command, "noother") ||
		strings.HasPrefix(command, "nosticker") ||
		strings.HasPrefix(command, "nogif") {
		return true
	}

	return false
}
