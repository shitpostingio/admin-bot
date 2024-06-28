package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
)

// kickUser kicks a user in a supergroup.
func kickUser(msg *tgbotapi.Message) {

	/* WE WANT TO BE STEALTHY */
	adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)

	//We want to re apply the restrictions to the users
	// chatMemberRestrictions, userIsRestricted, _ := utility.GetChatMemberRestrictions(msg.ReplyToMessage.From.ID, msg.Chat.ID, admin)
	adminbot.KickUser(msg.ReplyToMessage.From.ID, msg.From.ID, msg.Chat.ID)
	// if userIsRestricted {
	// 	utility.restrictUser(msg.ReplyToMessage.From, chatMemberRestrictions, &admin.Bot.Self, admin)
	// }

}
