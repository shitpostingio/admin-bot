package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/repository"
)

// GetChatMember gets a chat member using the bot API.
func GetChatMember(userID int64, groupID int64) (tgbotapi.ChatMember, error) {

	chatMemberConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: groupID,
			UserID: userID,
		},
	}

	return bot.GetChatMember(chatMemberConfig)
}

// GetUserProfilePhotos returns the user's profile photos using the bot API.
func GetUserProfilePhotos(userID int64, maxPhotos int) (tgbotapi.UserProfilePhotos, error) {

	userPhotoConfig := tgbotapi.UserProfilePhotosConfig{
		UserID: userID,
		Limit:  maxPhotos,
	}

	return repository.Bot.GetUserProfilePhotos(userPhotoConfig)
}
