package api

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/api/botapi"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// GetChatMember gets a chat member using the bot API.
// It will also use a rate limiter not to get restricted by Telegram.
func GetChatMember(userID int64, groupID int64) (tgbotapi.ChatMember, error) {
	limiter.AuthorizeAction()
	return botapi.GetChatMember(userID, groupID)
}

// GetUserProfilePhotos returns the user's profile photos.
// It will also use a rate limiter not to get restricted by Telegram.
func GetUserProfilePhotos(userID int64, maxPhotos int) (tgbotapi.UserProfilePhotos, error) {
	limiter.AuthorizeAction()
	return botapi.GetUserProfilePhotos(userID, maxPhotos)
}
