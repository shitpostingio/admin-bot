package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/config/structs"
)

//SetBot sets the bot in the repository.
func SetBot(inputBot *tgbotapi.BotAPI) {
	Bot = inputBot
}

//SetConfig sets the configuration in the repository.
func SetConfig(inputCfg *structs.Config) {
	Configuration = inputCfg
}

//SetMods sets the mod map in the repository.
func SetMods(inputMap map[int64]bool) {
	Mods = inputMap
}

//SetAdmins sets the admin map in the repository.
func SetAdmins(inputMap map[int64]bool) {
	Admins = inputMap
}

//SetTestingStatus sets the testing status in the repository.
func SetTestingStatus(testingStatus bool) {
	Testing = testingStatus
}

//SetDebugStatus sets the debug status in the repository.
func SetDebugStatus(debugStatus bool) {
	Debug = debugStatus
}
