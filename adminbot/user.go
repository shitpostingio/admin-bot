package adminbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api"
)

// GetChatMember gets a chat member using the bot API in a rate limited fashion.
func GetChatMember(userID int64, groupID int64) (tgbotapi.ChatMember, error) {
	return api.GetChatMember(userID, groupID)
}

// GetUserProfilePhotos returns the user's profile photos in a rate limited fashion.
func GetUserProfilePhotos(userID int64, maxPhotos int) (tgbotapi.UserProfilePhotos, error) {
	return api.GetUserProfilePhotos(userID, maxPhotos)
}
