package cache

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/repository"
)

// CreateAdminsCache caches the database administrators in a map.
func CreateAdminsCache(collection *mongo.Collection) (map[int64]bool, error) {

	dbModerators, err := database.GetAllModerators()
	if err != nil {
		return nil, fmt.Errorf("CreateAdminsCache: no moderators in the database: %s", err)
	}

	isAdmin := make(map[int64]bool)
	for _, moderator := range dbModerators {

		if moderator.IsAdmin {
			isAdmin[moderator.TelegramID] = true
		}

	}

	return isAdmin, nil

}

// CreateModsCache caches the mods in a map.
func CreateModsCache(isAdmin map[int64]bool, bot *tgbotapi.BotAPI, cfg *structs.TelegramConfiguration) (map[int64]bool, error) {

	chatAdministratorsConfig := tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: cfg.GroupID}}
	chatAdministrators, err := bot.GetChatAdministrators(chatAdministratorsConfig)
	if err != nil {
		return nil, fmt.Errorf("CreateModsCache: unable to retrieve chat administrators: %s", err)
	}

	// We don't want database admins in the mod map.
	isMod := make(map[int64]bool)
	for _, moderator := range chatAdministrators {
		if !isAdmin[moderator.User.ID] && !moderator.User.IsBot {
			isMod[moderator.User.ID] = true
		}
	}

	return isMod, nil

}

//RemoveFromMods removes the `userID` from the mod map, if present.
func RemoveFromMods(userID int64) bool {

	chatMods := repository.Mods

	_, wasMod := chatMods[userID]
	if wasMod {
		delete(chatMods, userID)
	}

	return wasMod

}
