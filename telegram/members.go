package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetPermissionsFromChatMember returns the user permissions.
func GetPermissionsFromChatMember(chatMember *tgbotapi.ChatMember) *tgbotapi.ChatPermissions {

	if chatMember == nil {
		return nil
	}

	return &tgbotapi.ChatPermissions{
		CanSendMessages:       chatMember.CanSendMessages,
		CanSendMediaMessages:  chatMember.CanSendMediaMessages,
		CanSendPolls:          chatMember.CanSendPolls,
		CanSendOtherMessages:  chatMember.CanSendOtherMessages,
		CanAddWebPagePreviews: chatMember.CanAddWebPagePreviews,
		CanChangeInfo:         chatMember.CanChangeInfo,
		CanInviteUsers:        chatMember.CanInviteUsers,
		CanPinMessages:        chatMember.CanPinMessages,
	}
}

// ChatMemberIsRestricted returns true if the user is restricted
func ChatMemberIsRestricted(chatMember *tgbotapi.ChatMember) bool {

	if chatMember == nil {
		return false
	}

	return chatMember.Status == "restricted"

}

// ChatMemberIsBanned returns true if the user was kicked
func ChatMemberIsBanned(chatMember *tgbotapi.ChatMember) bool {

	if chatMember == nil {
		return false
	}

	return chatMember.Status == "kicked"

}
