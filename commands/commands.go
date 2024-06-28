package commands

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/utility"
)

// ExecuteCommand executes commands, if the user is allowed
func ExecuteCommand(msg *tgbotapi.Message) {

	if !utility.IsChatAdminByMessage(msg) {
		adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
		return
	}

	command := strings.ToLower(msg.Command())
	if msg.ReplyToMessage == nil {
		adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
		handleCommandsWithoutReply(command, msg)
		return
	}

	if msg.CommandArguments() == "" {
		handleCommandsWithoutArguments(command, msg)
		return
	}

	handleCommandsWithArguments(command, msg)

}

func handleCommandsWithoutReply(command string, msg *tgbotapi.Message) {

	switch {
	case command == "banh":
		banByUsername(msg)
	case command == "idban":
		banByID(msg)
	}

}

func handleCommandsWithArguments(command string, msg *tgbotapi.Message) {
	switch {
	case command == "ban":
		banUser(msg)
	case command == "mute", strings.HasPrefix(command, "no"):
		restrictUser(msg)
	}
}

func handleCommandsWithoutArguments(command string, msg *tgbotapi.Message) {
	switch {
	case command == "blsp" || command == "bls":
		blacklistStickerPack(msg)
	case command == "blm":
		blacklistMedia(msg)
	case command == "blh":
		blacklistHandle(msg)
	case command == "kick":
		kickUser(msg)
	case command == "mute", strings.HasPrefix(command, "no"):
		restrictUser(msg)
	}
}
