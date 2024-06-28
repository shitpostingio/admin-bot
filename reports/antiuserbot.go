package reports

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/localization"
)

func PossibleUserbotAttack() string {
	return localization.GetString("antiuserbot_possible_attack")
}

func PossibleUserbotRestriction(userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("antiuserbot_possible_userbot_restriction"),
		userID,
		userName)
}
