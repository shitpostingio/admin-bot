package buttons

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//CreateAtAdminReportMarkup creates the InlineMarkup for a @admin mention
func CreateAtAdminReportMarkup(reportedMessageID, reportMessageID, backupMessageID int, chatType string) (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = CreateReportAndBackupMarkup(reportedMessageID, reportMessageID, backupMessageID, chatType)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(CreateHandleButton()))
	return
}

//CreateReportAndBackupMarkup creates the InlineMarkup containing links to the reported message, report and backup
func CreateReportAndBackupMarkup(reportedMessageID, reportMessageID, backupMessageID int, chatType string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var reportAndBackupRow []tgbotapi.InlineKeyboardButton

	/* REPORTED MESSAGE ROW (OPTIONAL) */
	if reportedMessageID != 0 {
		reportedMessageRow := tgbotapi.NewInlineKeyboardRow(CreateReportedMessageButton(reportedMessageID, chatType))
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, reportedMessageRow)
	}

	/* REPORT AND BACKUP ROW */
	if reportMessageID != 0 {
		reportMessageButton := CreateReportMessageButton(reportMessageID, chatType)
		reportAndBackupRow = append(reportAndBackupRow, reportMessageButton)
	}

	if backupMessageID != 0 {
		backupMessageButton := CreateBackupMessageButton(backupMessageID, chatType)
		reportAndBackupRow = append(reportAndBackupRow, backupMessageButton)
	}

	/* IF WE ADD AN EMPTY ROW TO THE KEYBOARD IT MIGHT CAUSE ISSUES */
	if len(reportAndBackupRow) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, reportAndBackupRow)
	}

	return keyboard
}

// CreateKeyboardWithOneRow creates an inline markup keyboard with one row
func CreateKeyboardWithOneRow(buttons ...tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}
