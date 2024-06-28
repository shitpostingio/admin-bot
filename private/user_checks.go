package private

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/entities"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
	"github.com/shitpostingio/admin-bot/utility"
)

func attemptCheckByText(msg *tgbotapi.Message) {

	// Try checking usernames first
	mentions := telegram.GetAllMentions(msg.Text, telegram.GetMessageEntities(msg), msg.ReplyMarkup)
	if len(mentions) != 0 {
		checkUserByUsername(msg, mentions[0])
		return
	}

	//check by id
	userID, err := strconv.ParseInt(msg.CommandArguments(), 10, 64)
	if err == nil {
		checkUserByID(msg, userID)
		return
	}

	_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, localization.GetString("private_checks_could_not_understand"), false)

}

func checkUserByUsername(msg *tgbotapi.Message, username string) {

	if tdlib.IsGroupOrChannelUsername(username) {
		result := fmt.Sprintf(localization.GetString("private_checks_username_is_group_or_channel"), username)
		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, result, false)
		return
	}

	chat, err := tdlib.ResolveUsername(username)
	if err != nil {
		result := fmt.Sprintf(localization.GetString("private_checks_unable_to_resolve_username"), err)
		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, result, false)
		return
	}

	tdlUser, err := tdlib.GetUserByID(chat.ID)
	if err != nil {
		result := fmt.Sprintf(localization.GetString("private_checks_unable_to_get_user_info"), err)
		_ = adminbot.SendReplyPlainTextMessage(msg.MessageID, msg.Chat.ID, result, false)
		return
	}

	tgUser := tdlib.GetTgbotapiUserFromTdlibUser(tdlUser)
	checkUserStatus(msg, tgUser)

}

func checkUserByID(msg *tgbotapi.Message, userID int64) {

	tdlibUser, err := tdlib.GetUserByID(userID)
	if err != nil {
		result := fmt.Sprintf(localization.GetString("private_checks_unable_to_get_user_info"), err)
		_ = adminbot.SendPlainTextMessage(msg.Chat.ID, result, false)
		return
	}

	tgUser := tdlib.GetTgbotapiUserFromTdlibUser(tdlibUser)
	checkUserStatus(msg, tgUser)

}

func checkUserStatus(msg *tgbotapi.Message, user *tgbotapi.User) {

	// Retrieve user status from Telegram and check the data in our database.
	chatMember, err := adminbot.GetChatMember(user.ID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		result := fmt.Sprintf(localization.GetString("private_checks_unable_to_get_user_info"), err)
		_ = adminbot.SendReplyTextMessage(msg.MessageID, msg.Chat.ID, result, false)
	}

	bans, _ := database.GetBansForTelegramID(user.ID)

	if telegram.ChatMemberIsBanned(&chatMember) {
		banStatus(msg, user, bans)
		return
	}

	// If we have bans related to the user in our database
	// but the user has been unbanned, soft-delete them.
	if len(bans) != 0 {
		_ = database.MarkUserAsUnbanned(user.ID)
	}

	if telegram.ChatMemberIsRestricted(&chatMember) {
		restrictionStatus(msg, &chatMember)
		return
	}

	result := fmt.Sprintf(localization.GetString("private_checks_not_banned_or_restricted"), user.ID, telegram.GetName(user))
	_ = adminbot.SendReplyTextMessage(msg.MessageID, msg.Chat.ID, result, false)

}

func banStatus(msg *tgbotapi.Message, user *tgbotapi.User, bans []entities.Ban) {

	timesBanned := len(bans)
	if timesBanned == 0 {

		result := fmt.Sprintf(localization.GetString("private_checks_no_db_info"), user.ID, telegram.GetName(user))
		markup := buttons.CreateKeyboardWithOneRow(buttons.CreateTgUnbanButton(user.ID))
		_ = adminbot.SendReplyTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, result, markup, false)
		return

	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(localization.GetString("private_checks_ban_preamble"), user.ID, telegram.GetName(user), timesBanned))

	for _, ban := range bans {

		builder.WriteString(fmt.Sprintf(localization.GetString("private_checks_ban_entry"), ban.BannedBy, ban.BannedBy,
			utility.FormatDate(ban.BanDate), ban.Reason))

	}

	result := builder.String()
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnbanButton(user.ID))
	_ = adminbot.SendReplyTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, result, markup, false)

}

func restrictionStatus(msg *tgbotapi.Message, chatMember *tgbotapi.ChatMember) {

	restrictions := telegram.GetPermissionsFromChatMember(chatMember)
	restrictedUser := telegram.GetName(chatMember.User)
	restrictionEnd := utility.FormatUnixDate(chatMember.UntilDate)

	result := fmt.Sprintf(localization.GetString("private_checks_restriction_report"), chatMember.User.ID, restrictedUser, restrictionEnd, emojifyRestrictions(restrictions))
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnrestrictButton(chatMember.User.ID))
	_ = adminbot.SendReplyTextMessageWithMarkup(msg.MessageID, msg.Chat.ID, result, markup, false)

}

func emojifyRestrictions(r *tgbotapi.ChatPermissions) string {

	if r == nil {
		return ""
	}

	sb := strings.Builder{}

	// Send Messages
	sb.WriteString(utility.EmojifyBool(r.CanSendMessages))
	sb.WriteString("\tCan send messages\n")

	// Send Media Messages
	sb.WriteString(utility.EmojifyBool(r.CanSendMediaMessages))
	sb.WriteString("\tCan send media\n")

	// Send Polls
	sb.WriteString(utility.EmojifyBool(r.CanSendPolls))
	sb.WriteString("\tCan send polls\n")

	// Send Other Messages
	sb.WriteString(utility.EmojifyBool(r.CanSendOtherMessages))
	sb.WriteString("\tCan send stickers and gifs\n")

	// Add Web Page Previews
	sb.WriteString(utility.EmojifyBool(r.CanAddWebPagePreviews))
	sb.WriteString("\tCan add web page previews\n")

	// Change info
	sb.WriteString(utility.EmojifyBool(r.CanChangeInfo))
	sb.WriteString("\tCan change group info\n")

	// Invite users
	sb.WriteString(utility.EmojifyBool(r.CanInviteUsers))
	sb.WriteString("\tCan invite users\n")

	// Pin messages
	sb.WriteString(utility.EmojifyBool(r.CanPinMessages))
	sb.WriteString("\tCan pin messages\n")

	return sb.String()

}
