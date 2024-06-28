package structs

// AdminBotConfiguration represents the admin-bot configuration.
type AdminBotConfiguration struct {
	Language         string
	LocalizationPath string
	WelcomeText      string `reloadable:"true"`
	EmergencyText    string `reloadable:"true"`
}
