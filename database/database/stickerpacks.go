package database

import (
	"github.com/shitpostingio/admin-bot/database/documentstore"
)

//BlacklistStickerPack blacklists a sticker pack by its set_name. Also logs the user who added it to the blacklist
func BlacklistStickerPack(setName string, blacklisterTelegramID int64) (generatedID string, err error) {
	return documentstore.BlacklistStickerPack(setName, blacklisterTelegramID, documentstore.StickerPackCollection)
}

//PardonStickerPack removes a sticker pack from the blacklist. Also logs the user who did it
func PardonStickerPack(setName string) error {
	return documentstore.PardonStickerPack(setName, documentstore.StickerPackCollection)
}

// StickerPackIsBlacklisted returns true if a sticker pack is not blacklisted
func StickerPackIsBlacklisted(setName string) bool {
	return documentstore.StickerPackIsBlacklisted(setName, documentstore.StickerPackCollection)
}
