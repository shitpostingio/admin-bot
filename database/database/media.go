package database

import (
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindMediaByFileID finds a media in the document store given its fileid
func FindMediaByFileID(uniqueFileID, fileID string) (media entities.Media, err error) {
	return documentstore.FindMediaByFileID(uniqueFileID, fileID, documentstore.MediaCollection)
}

func FindMediaByFileUniqueID(fileUniqueID string) (media entities.Media, err error) {
	return documentstore.FindMediaByFileUniqueID(fileUniqueID, documentstore.MediaCollection)
}

// FindMediaByFeatures finds a media in the document store given its features
func FindMediaByFeatures(histogram []float64, pHash string, approximation float64) (media entities.Media, err error) {
	return documentstore.FindMediaByFeatures(histogram, pHash, approximation, documentstore.MediaCollection)
}

// FindMediaByID finds a media in the document store given an ObjectID
func FindMediaByID(ID *primitive.ObjectID) (media entities.Media, err error) {
	return documentstore.FindMediaByID(ID, documentstore.MediaCollection)
}

// BlacklistMedia blacklists a media
func BlacklistMedia(uniqueFileID, fileID string, userID int64) (generatedID string, err error) {
	return documentstore.BlacklistMedia(uniqueFileID, fileID, userID, documentstore.MediaCollection)
}

// WhitelistMedia whitelists a media
func WhitelistMedia(uniqueFileID, fileID string, userID int64) (generatedID string, err error) {
	return documentstore.WhitelistMedia(uniqueFileID, fileID, userID, documentstore.MediaCollection)
}

// RemoveMedia removes a media from the document store
func RemoveMedia(uniqueFileID, fileID string) error {
	return documentstore.RemoveMedia(uniqueFileID, fileID, documentstore.MediaCollection)
}

// BlacklistNSFWMedia blacklists a media for being nsfw
func BlacklistNSFWMedia(uniqueFileID, fileID string, description string, score float64, userID int64) (generatedID string, err error) {
	return documentstore.BlacklistNSFWMedia(uniqueFileID, fileID, description, score, userID, documentstore.MediaCollection)
}

// MediaIsBlacklisted returns true if a media is blacklisted
func MediaIsBlacklisted(uniqueFileID, fileID string) bool {
	return documentstore.MediaIsBlacklisted(uniqueFileID, fileID, documentstore.MediaCollection)
}

// MediaIsWhitelisted returns true if a media is whitelisted
func MediaIsWhitelisted(uniqueFileID, fileID string) bool {
	return documentstore.MediaIsWhitelisted(uniqueFileID, fileID, documentstore.MediaCollection)
}
