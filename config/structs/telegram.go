package structs

// TelegramConfiguration represents the Telegram configuration
type TelegramConfiguration struct {
	BotToken          string
	GroupID           int64  `reloadable:"true"`
	ReportChannelID   int64  `reloadable:"true"`
	BackupChannelID   int64  `reloadable:"true"`
	GroupLink         string `type:"optional" reloadable:"true"`
	BackupChannelLink string `type:"optional" reloadable:"true"`
}

// WebHookPath returns only the relative path where Telegram will send updates
func (c TelegramConfiguration) WebHookPath() string {
	return "/" + c.BotToken + "/updates"
}
