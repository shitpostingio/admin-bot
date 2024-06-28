package botapi

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RestrictUser restricts a user using the bot API.
func RestrictUser(userID int64, chatID int64, untilDate int64,
	canSendMessages,
	canSendMediaMessages,
	canSendOtherMessages,
	canAddWebPagePreviews bool) (*tgbotapi.APIResponse, error) {

	restrictionConfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{UserID: userID, ChatID: chatID},
		UntilDate:        untilDate,
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       canSendMessages,
			CanSendMediaMessages:  canSendMediaMessages,
			CanSendOtherMessages:  canSendOtherMessages,
			CanAddWebPagePreviews: canAddWebPagePreviews,
		},
	}

	response, err := bot.Request(restrictionConfig)
	return response, err
}

// UnrestrictUser unrestricts a user using the bot API.
func UnrestrictUser(userID int64, chatID int64) (*tgbotapi.APIResponse, error) {

	unrestrictionConfig := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{UserID: userID, ChatID: chatID},
		UntilDate:        0,
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         true,
			CanInviteUsers:        true,
			CanPinMessages:        true,
		},
	}

	response, err := bot.Request(unrestrictionConfig)
	return response, err
}
