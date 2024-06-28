package database

import (
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
)

// GetSourceByUsername gets a source from the document store by its username
func GetSourceByUsername(username string) (source entities.Source, err error) {
	return documentstore.GetSourceByUsername(username, documentstore.SourcesCollection)
}

// GetSourceByTelegramID gets a source from the document store by its telegram id
func GetSourceByTelegramID(telegramID int64) (source entities.Source, err error) {
	return documentstore.GetSourceByTelegramID(telegramID, documentstore.SourcesCollection)
}

// GetSource gets a source from the document store
func GetSource(telegramID int64, username string) (source entities.Source, err error) {
	return documentstore.GetSource(telegramID, username, documentstore.SourcesCollection)
}

// BlacklistSource blacklists a source in the document store
func BlacklistSource(sourceTelegramID int64, sourceUsername string, adderTelegramID int64) (generatedID string, err error) {
	return documentstore.BlacklistSource(sourceTelegramID, sourceUsername, adderTelegramID, documentstore.SourcesCollection)
}

// WhitelistSource whitelists a source in the document store
func WhitelistSource(sourceTelegramID int64, sourceUsername string, adderTelegramID int64) (generatedID string, err error) {
	return documentstore.WhitelistSource(sourceTelegramID, sourceUsername, adderTelegramID, documentstore.SourcesCollection)
}

// RemoveSource removes a source from the document store
func RemoveSource(telegramID int64, username string) error {
	return documentstore.RemoveSource(telegramID, username, documentstore.SourcesCollection)
}

// UpdateSource updates a source in the document store
func UpdateSource(telegramID int64, username string) error {
	return documentstore.UpdateSource(telegramID, username, documentstore.SourcesCollection)
}

// SourceIsWhitelisted returns true if the source is whitelisted
func SourceIsWhitelisted(telegramID int64, username string) bool {
	return documentstore.SourceIsWhitelisted(telegramID, username, documentstore.SourcesCollection)
}

// SourceIsBlacklisted returns true if the source is blacklisted
func SourceIsBlacklisted(telegramID int64, username string) bool {
	return documentstore.SourceIsBlacklisted(telegramID, username, documentstore.SourcesCollection)
}
