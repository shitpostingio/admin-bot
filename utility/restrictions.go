package utility

import (
	"time"

	"github.com/shitpostingio/admin-bot/repository"
)

//GetAppropriateRestrictionEnd returns 1 minute if
//the bot is in testing mode, 30 minutes otherwise
func GetAppropriateRestrictionEnd() int64 {

	if repository.GetTestingStatus() {
		return time.Now().Add(1 * time.Minute).Unix()
	}

	return time.Now().Add(30 * time.Minute).Unix()
}
