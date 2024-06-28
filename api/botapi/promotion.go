package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//PromoteUser promotes user using the bot API.
func PromoteUser(userID int64, chatID int64,
	canChangeInfo,
	canDeleteMessages,
	canInviteUsers,
	canRestrictMembers,
	canPinMessages,
	canPromoteMembers bool) (*tgbotapi.APIResponse, error) {

	promotionConfig := tgbotapi.PromoteChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		CanChangeInfo:      canChangeInfo,
		CanDeleteMessages:  canDeleteMessages,
		CanInviteUsers:     canInviteUsers,
		CanRestrictMembers: canRestrictMembers,
		CanPinMessages:     canPinMessages,
		CanPromoteMembers:  canPromoteMembers,
	}

	response, err := bot.Request(promotionConfig)
	return response, err
}
