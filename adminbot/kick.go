package adminbot

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api"
)

// KickUser kicks a user from a group.
func KickUser(kickedUserID, moderatorUserID int64, chatID int64) {
	response, err := api.KickUser(kickedUserID, moderatorUserID, chatID)
	if err != nil {

		kickMsg := fmt.Sprintf("Unable to kick user with id %d: %s", kickedUserID, err)

		if response != nil {
			log.Error(kickMsg, "(", response.ErrorCode, ":", response.Description)
		} else {
			log.Error(kickMsg)
		}

	}
}
