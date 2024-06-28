package database

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
)

// GetAllModerators retrieves all the moderators from the database.
func GetAllModerators() (moderators []entities.Moderator, err error) {
	return documentstore.GetAllModerators(documentstore.ModeratorsCollection)
}

// GetModeratorByTelegramID retrieves the moderator with the given telegram id.
func GetModeratorByTelegramID(telegramID int64) (moderator entities.Moderator, err error) {
	return documentstore.GetModeratorByTelegramID(telegramID, documentstore.ModeratorsCollection)
}

// GetModeratorByUsername retrieves the moderator with the given username.
func GetModeratorByUsername(username string) (moderator entities.Moderator, err error) {
	return documentstore.GetModeratorByUsername(username, documentstore.ModeratorsCollection)
}

// AddModerator adds the given user to the moderators.
func AddModerator(user *tgbotapi.ChatMember, moddedByTelegramID int64) (generatedID string, err error) {
	return documentstore.AddModerator(user, moddedByTelegramID, documentstore.ModeratorsCollection)
}

func RemoveModerator(unmoddedUserID int64) (err error) {
	return documentstore.RemoveModerator(unmoddedUserID, documentstore.ModeratorsCollection)
}

// UpdateModeratorsDetails updates the details of the moderators in the database.
func UpdateModeratorsDetails(chatID int64, bot *tgbotapi.BotAPI) error {
	return documentstore.UpdateModeratorsDetails(chatID, bot, documentstore.ModeratorsCollection)
}

// UpdateModeratorDetailsByTelegramUser updates the details of a user in the database.
func UpdateModeratorDetailsByTelegramUser(user *tgbotapi.User, chatID int64, bot *tgbotapi.BotAPI) error {
	return documentstore.UpdateModeratorDetailsByTelegramUser(user, chatID, bot, documentstore.ModeratorsCollection)
}

// UpdateModeratorDetails updates the details of a moderator in the database.
func UpdateModeratorDetails(chatMember *tgbotapi.ChatMember) error {
	return documentstore.UpdateModeratorDetails(chatMember, documentstore.ModeratorsCollection)
}

// IsModeratorUsername returns true if the given username belongs to a moderator.
func IsModeratorUsername(username string) bool {
	return documentstore.IsModeratorUsername(username, documentstore.ModeratorsCollection)
}

// IsModerator returns true if the given telegram id belongs to a moderator.
func IsModerator(telegramID int64) bool {
	return documentstore.IsModerator(telegramID, documentstore.ModeratorsCollection)
}
