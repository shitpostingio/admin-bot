package database

import (
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
)

// GetHostName gets a hostname from the document store
func GetHostName(url string) (host entities.HostName, err error) {
	return documentstore.GetHostName(url, documentstore.HostNameCollection)
}

// BlacklistHostName blacklists a host name
func BlacklistHostName(url string, isBanworthy, isTelegram bool, adderTelegramID int64) (generatedID, hostname string, err error) {
	return documentstore.BlacklistHostName(url, isBanworthy, isTelegram, adderTelegramID, documentstore.HostNameCollection)
}

// UpdateHostName updates a hostname in the document store
func UpdateHostName(url string, isBanworthy, isTelegram bool, updaterTelegramID int64) error {
	return documentstore.UpdateHostName(url, isBanworthy, isTelegram, updaterTelegramID, documentstore.HostNameCollection)
}

// PardonHostName removes a hostname from the document store
func PardonHostName(url string) error {
	return documentstore.PardonHostName(url, documentstore.HostNameCollection)
}
