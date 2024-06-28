package buttons

import (
	"fmt"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *												CALLBACK BUTTONS													   *
 *																													   *
 ***********************************************************************************************************************
 */

// CreateHandleButton creates a button to mark an action as handled
func CreateHandleButton() tgbotapi.InlineKeyboardButton {
	markAsHandledString := "2s ok"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_mark_as_handled"), markAsHandledString)
}

// CreateWhitelistMediaButton creates a button to whitelist in 2 steps a media
func CreateWhitelistMediaButton(nsfwTableID string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s whitelist %s", nsfwTableID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_whitelist_media"), actionToPerform)
}

// CreateTgUnbanButton creates a button to unban in two steps a user
// whose ban was not in the database
func CreateTgUnbanButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa tgunban %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_unban_user"), actionToPerform)
}

//CreateUnbanButton creates a button to unban in two steps a user
func CreateUnbanButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa unban %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_unban_user"), actionToPerform)
}

// CreateUnrestrictButton creates a button to unrestrict in two steps a user
func CreateUnrestrictButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa unrestrict %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_unrestrict_user"), actionToPerform)
}

// CreateModUserButton creates a button to mod a user
func CreateModUserButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa mod %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_mod_user"), actionToPerform)
}

// CreateBlacklistBanworthyHostnameButton creates a button to add a banworthy hostname to the blacklist
func CreateBlacklistBanworthyHostnameButton(hostname string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa bbh %s", hostname)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_banworthy"), actionToPerform)
}

// CreateBlacklistTelegramHostnameButton creates a button to add a telegram hostname to the blacklist
func CreateBlacklistTelegramHostnameButton(hostname string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa bth %s", hostname)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_telegram"), actionToPerform)
}

// CreateBlacklistHostnameButton creates a button to add a hostname to the blacklist
func CreateBlacklistHostnameButton(hostname string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2sa bh %s", hostname)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_remove_only"), actionToPerform)
}

// CreateBanForGibberishHandleButton creates a button to ban a user for sending a gibberish bot handle
func CreateBanForGibberishHandleButton(userID int64, handle string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s bg %d %s", userID, handle)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_ban_user"), actionToPerform)
}

// CreateHumanVerificationButton creates a button for users to verify that they've read the rules
func CreateHumanVerificationButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("verify %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("buttons_rules_read"), actionToPerform)
}

// CreateApproveUserButton creates a button to approve a user in an emergency
func CreateApproveUserButton(userID int64) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s unrestrict %d", userID)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("emergencymode_approve_user"), actionToPerform)
}

func CreatePrivateBlacklistMediaButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s blm"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_blacklist"), actionToPerform)
}

func CreatePrivateWhitelistMediaButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s wlm"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_whitelist"), actionToPerform)
}

func CreatePrivatePardonMediaButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s parm"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_remove"), actionToPerform)
}

func CreatePrivateBlacklistStickerPackButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s blsp"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_blacklist_pack"), actionToPerform)
}

func CreatePrivatePardonStickerPackButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s parsp"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_remove_pack"), actionToPerform)
}

func CreatePrivateBlacklistStickerButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s blms"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_blacklist"), actionToPerform)
}

func CreatePrivateWhitelistStickerButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s wlms"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_whitelist"), actionToPerform)
}

func CreatePrivatePardonStickerButton() tgbotapi.InlineKeyboardButton {
	actionToPerform := "2s parms"
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_remove"), actionToPerform)
}

func CreatePrivateBlacklistSourceButton(source string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s bs %s", source)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_blacklist"), actionToPerform)
}

func CreatePrivatePardonSourceButton(source string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s ps %s", source)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_remove"), actionToPerform)
}

func CreatePrivateWhitelistSourceButton(source string) tgbotapi.InlineKeyboardButton {
	actionToPerform := fmt.Sprintf("2s ws %s", source)
	return tgbotapi.NewInlineKeyboardButtonData(localization.GetString("private_button_whitelist"), actionToPerform)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													 URL BUTTONS													   *
 *																													   *
 ***********************************************************************************************************************
 */

// CreateReportedMessageButton creates a button with an URL for a reported message
func CreateReportedMessageButton(messageID int, chatType string) tgbotapi.InlineKeyboardButton {

	var url string
	if repository.GetTelegramConfiguration().GroupLink != "" {
		url = fmt.Sprintf("%s/%d", repository.GetTelegramConfiguration().GroupLink, messageID)
	} else {
		url = fmt.Sprintf("https://t.me/c/%s/%d", getPrivateChatIDString(repository.GetTelegramConfiguration().GroupID, chatType), messageID)
	}

	return tgbotapi.NewInlineKeyboardButtonURL(localization.GetString("buttons_reported_message"), url)
}

// CreateReportMessageButton creates a button with an URL for a report message
func CreateReportMessageButton(messageID int, chatType string) tgbotapi.InlineKeyboardButton {

	var url string
	if repository.GetTelegramConfiguration().GroupLink != "" {
		url = fmt.Sprintf("%s/%d", repository.GetTelegramConfiguration().GroupLink, messageID)
	} else {
		url = fmt.Sprintf("https://t.me/c/%s/%d", getPrivateChatIDString(repository.GetTelegramConfiguration().GroupID, chatType), messageID)
	}

	return tgbotapi.NewInlineKeyboardButtonURL(localization.GetString("buttons_report"), url)
}

// CreateBackupMessageButton creates a button with an URL for a backup message
func CreateBackupMessageButton(messageID int, chatType string) tgbotapi.InlineKeyboardButton {

	var url string
	if repository.GetTelegramConfiguration().BackupChannelLink != "" {
		url = fmt.Sprintf("%s/%d", repository.GetTelegramConfiguration().BackupChannelLink, messageID)
	} else {
		url = fmt.Sprintf("https://t.me/c/%s/%d", getPrivateChatIDString(repository.GetTelegramConfiguration().BackupChannelID, chatType), messageID)
	}

	return tgbotapi.NewInlineKeyboardButtonURL(localization.GetString("buttons_backup"), url)
}

// getPrivateChatIDString returns the chatID converted for private chat links.
// As per Anime Sex Storm:
// This value is modified from a normal integer into this value based on the chat type.
// A "private" chat will always be a normal int,
// A "group" chat will be an int in the negatives,
// A "supergroup" or "channel" chats will be negative and prepended with 100.
func getPrivateChatIDString(originalChatID int64, chatType string) string {

	chatIDStr := strconv.FormatInt(originalChatID, 10)

	if chatType == "group" {
		return chatIDStr[1:]
	}

	return chatIDStr[4:]
}
