package adminbot

import (
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/database/database"
)

// PromoteToMod promotes a user to a moderator.
func PromoteToMod(userID int64, chatID int64) error {
	_, err := api.PromoteUser(userID, chatID, false, true, false, false, false, false)
	return err
}

//PromoteToAdmin promotes a database admin to admin.
func PromoteToAdmin(adminUserID int64, chatID int64) {

	a, err := database.GetModeratorByTelegramID(adminUserID)
	if err != nil {
		log.Error("PromoteToAdmin:", err)
		return
	}

	_, err = api.PromoteUser(adminUserID, chatID,
		a.CanChangeInfo, a.CanDeleteMessages, a.CanInviteUsers, a.CanRestrictMembers, a.CanPinMessages, a.CanPromoteMembers)
	if err != nil {
		log.Error("Unable to promote admin with telegramID ", a.TelegramID, err)
	}

}
