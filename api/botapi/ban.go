package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BanUser bans a user in a chat using the bot API.
func BanUser(bannedUserID int64, chatID int64) error {

	//
	userToBan := tgbotapi.ChatMemberConfig{UserID: bannedUserID, ChatID: chatID}
	banConfig := tgbotapi.KickChatMemberConfig{ChatMemberConfig: userToBan}

	//
	_, err := bot.Request(banConfig)
	return err
}

func BanChannel(bannedSenderChatID, chatID int64) error {

	//
	chatToBan := tgbotapi.BanChatSenderChatConfig{SenderChatID: bannedSenderChatID, ChatID: chatID}

	//
	_, err := bot.Request(chatToBan)
	return err
}
