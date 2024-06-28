package repository

import (
	"github.com/shitpostingio/admin-bot/config/structs"
)

//GetTestingStatus gets the testing status from the repository.
func GetTestingStatus() bool {
	return Testing
}

//GetDebugStatus gets the debug status from the repository.
func GetDebugStatus() bool {
	return Debug
}

// GetTelegramConfiguration gets the telegram configuration
func GetTelegramConfiguration() *structs.TelegramConfiguration {
	return &Configuration.Telegram
}

// GetAntiSpamConfiguration gets the antispam configuration
func GetAntiSpamConfiguration() *structs.AntiSpamConfiguration {
	return &Configuration.AntiSpam
}

// GetAntiFloodConfiguration gets the antiflood configuration
func GetAntiFloodConfiguration() *structs.AntiFloodConfiguration {
	return &Configuration.AntiFlood
}

// GetAntiUserbotConfiguration gets the antiuserbot configuration
func GetAntiUserbotConfiguration() *structs.AntiUserbotConfiguration {
	return &Configuration.AntiUserbot
}
