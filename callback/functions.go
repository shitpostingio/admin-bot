package callback

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/defense/antispam"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/shitpostingio/admin-bot/consts"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/utility"
)

// markReportAsHandled marks a report as handled and saves
// who handled it.
func markReportAsHandled(callbackQuery *tgbotapi.CallbackQuery) {

	handledBy := telegram.GetName(callbackQuery.From)
	originalText := getOriginalMessageText(callbackQuery.Message)
	updatedText := fmt.Sprintf(localization.GetString("callback_report_handled"), callbackQuery.From.ID, handledBy, originalText)

	updatedReplyMarkup := removeActionButtons(callbackQuery.Message.ReplyMarkup)
	err := adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, updatedText, consts.ReportParseMode, &updatedReplyMarkup)
	if err != nil {
		log.Error("Unable to update message after marking report as handled:", err)
	}

}

//whitelistMedia whitelists media via callback
func whitelistMedia(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	nsfwTableID, _ := primitive.ObjectIDFromHex(callbackFields[1])
	nsfwEntity, err := database.FindMediaByID(&nsfwTableID)
	if err != nil {
		log.Error("whitelistMedia:", err)
		return
	}

	_, _ = database.WhitelistMedia(nsfwEntity.FileID, nsfwEntity.FileID, callbackQuery.From.ID)

	updatedText := fmt.Sprintf(localization.GetString("callback_media_whitelisted"), callbackQuery.From.ID, telegram.GetName(callbackQuery.From), getOriginalMessageText(callbackQuery.Message))
	updatedReplyMarkup := removeActionButtons(callbackQuery.Message.ReplyMarkup)
	err = adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, updatedText, consts.ReportParseMode, &updatedReplyMarkup)
	if err != nil {
		log.Error("Unable to update message after whitelisting media:", err)
	}

}

//tgunbanUser unbans a user that wasn't in the database after adding them
func tgunbanUser(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {
	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)
	_, _ = database.AddBan(userID, repository.Bot.Self.ID, "Manual ban")
	unbanUser(callbackFields, callbackQuery)
}

//unbanUser unbans a user
func unbanUser(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)
	err := adminbot.UnbanUser(userID, repository.GetTelegramConfiguration().GroupID, callbackQuery.From)
	if err != nil {
		log.Error("Unable to unban user with ID", userID, "(requested by", telegram.GetNameOrUsername(callbackQuery.From), "):", err)
		return
	}

	log.Info("User with user ID", callbackFields[1], "has been unbanned by", telegram.GetNameOrUsername(callbackQuery.From))

	updatedText := fmt.Sprintf(localization.GetString("callback_user_unbanned"), callbackQuery.From.ID, telegram.GetName(callbackQuery.From), getOriginalMessageText(callbackQuery.Message))
	err = adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, updatedText, consts.ReportParseMode, nil)
	if err != nil {
		log.Error("Unable to update message after unbanning user:", err)
	}

}

func banUserForGibberish(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)
	handle := callbackFields[2]

	_, err := adminbot.BanUserByID(userID, callbackQuery.From, "sending a gibberish handle", repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		log.Error("Unable to ban user with ID", userID, "for sending a gibberish handle")
		return
	}

	adminbot.DeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	_, _ = database.BlacklistSource(0, handle, callbackQuery.From.ID)

}

//unrestrictUser unrestricts a user and terminates the antispam routine, if active
func unrestrictUser(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)

	chatMember, err := adminbot.GetChatMember(userID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		log.Error("Unable to find user with ID", userID, "in the group for the unrestriction requested by",
			telegram.GetNameOrUsername(callbackQuery.From), ": ", err)
		return

	}

	err = adminbot.UnrestrictUser(chatMember.User, repository.GetTelegramConfiguration().GroupID, callbackQuery.From)
	if err != nil {

		log.Error("Unable to unrestrict user with ID", userID, "(requested by", telegram.GetNameOrUsername(callbackQuery.From), "):", err)
		return

	}

	antispam.EndAntiSpamRoutineForUser(userID)

	updatedText := fmt.Sprintf(localization.GetString("callback_user_unrestricted"), callbackQuery.From.ID, telegram.GetName(callbackQuery.From), getOriginalMessageText(callbackQuery.Message))
	err = adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, updatedText, consts.ReportParseMode, nil)
	if err != nil {
		log.Error("Unable to update message after unrestricting user:", err)
	}

}

//modUser mods a user (only CanDeleteMessages)
func modUser(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)
	err := adminbot.PromoteToMod(userID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		log.Error(fmt.Sprintf(localization.GetString("callback_mods_unable_to_mod"), userID, telegram.GetNameOrUsername(callbackQuery.From), err.Error()))
		return
	}

	// We must also add the user to the mod map.
	repository.Mods[userID] = true
	log.Info("User with user ID", callbackFields[1], "has been promoted by", telegram.GetNameOrUsername(callbackQuery.From))

	updatedText := fmt.Sprintf("🛃 User promoted 🛃\n\n%s", callbackQuery.Message.Text)
	err = adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, updatedText, consts.ReportParseMode, nil)
	if err != nil {
		log.Error("Unable to update message after promoting user to mod:", err)
	}

	chatMember, err := adminbot.GetChatMember(userID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		// Try to add it in another way
		_ = database.UpdateModeratorsDetails(repository.GetTelegramConfiguration().GroupID, repository.Bot)
	} else {
		_, _ = database.AddModerator(&chatMember, callbackQuery.From.ID)
	}

}

//handleBlacklistHostname displays buttons to handle certain hostnames
func handleBlacklistHostname(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	//TODO: MIGLIORARE MOLTO

	var err error

	switch callbackFields[0] {
	case "bbh":
		_, _, err = database.BlacklistHostName(callbackFields[1], true, false, callbackQuery.From.ID)
	case "bth":
		_, _, err = database.BlacklistHostName(callbackFields[1], false, true, callbackQuery.From.ID)
	case "bh":
		_, _, err = database.BlacklistHostName(callbackFields[1], false, false, callbackQuery.From.ID)
	default:
		return
	}

	dbEntity, _ := database.GetHostName(callbackFields[1])
	text := fmt.Sprintf(localization.GetString("callback_hosts_result"), err != nil, dbEntity.Host, utility.EmojifyBool(dbEntity.IsBanworthy), utility.EmojifyBool(dbEntity.IsTelegram))
	err = adminbot.EditMessageText(callbackQuery.Message.MessageID, callbackQuery.Message.Chat.ID, text, consts.ReportParseMode, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to send message to blacklist hostname: %s", err.Error()))
	}
}

func getOriginalMessageText(msg *tgbotapi.Message) string {

	if msg.Entities == nil {
		return msg.Text
	}

	builder := strings.Builder{}
	mRunesUTF16 := utf16.Encode([]rune(msg.Text))
	previousIndex := 0
	normalizedIndex := 0

	for _, entity := range msg.Entities {

		normalizedIndex = entity.Offset

		switch entity.Type {
		case "text_mention":

			builder.WriteString(fmt.Sprintf("%s<a href=\"tg://user?id=%d\">%s</a>",
				string(utf16.Decode(mRunesUTF16[previousIndex:normalizedIndex])),
				entity.User.ID, telegram.GetName(entity.User)))

			previousIndex = normalizedIndex + entity.Length

		case "code":

			builder.WriteString(fmt.Sprintf("%s<code>%s</code>",
				string(utf16.Decode(mRunesUTF16[previousIndex:normalizedIndex])),
				string(utf16.Decode(mRunesUTF16[normalizedIndex:normalizedIndex+entity.Length]))))

			previousIndex = normalizedIndex + entity.Length
		}
	}

	builder.WriteString(string(utf16.Decode(mRunesUTF16[previousIndex:])))
	return builder.String()
}

func removeActionButtons(replyMarkup *tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {

	if replyMarkup == nil {
		return tgbotapi.InlineKeyboardMarkup{}
	}

	// The `url` buttons are separated from the others
	// we can truncate the keyboard when we first see
	// a button with some CallbackData in it.
	targetRow := 0
	for rowID, row := range replyMarkup.InlineKeyboard {
		if row[0].CallbackData != nil {
			targetRow = rowID
			break
		}
	}

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: replyMarkup.InlineKeyboard[:targetRow]}
}
