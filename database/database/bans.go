package database

import (
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
)

// GetBansForTelegramID gets bans from the document store given a user's telegram id
func GetBansForTelegramID(telegramID int64) (bans []entities.Ban, err error) {
	return documentstore.GetBansForTelegramID(telegramID, documentstore.BansCollection)
}

// AddBan adds a ban to the document store
func AddBan(bannedUserID, moderatorUserID int64, reason string) (generatedID string, err error) {
	return documentstore.AddBan(bannedUserID, moderatorUserID, reason, documentstore.BansCollection)
}

// MarkUserAsUnbanned marks a user as unbanned in the document store
func MarkUserAsUnbanned(telegramID int64) error {
	return documentstore.MarkUserAsUnbanned(telegramID, documentstore.BansCollection)
}
