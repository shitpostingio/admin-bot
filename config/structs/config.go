package structs

// Config represents the bot configuration
type Config struct {
	//nolint: maligned
	Telegram      TelegramConfiguration
	Database      DatabaseConfiguration
	DocumentStore DocumentStoreConfiguration
	Loglog        LoglogConfiguration
	Webhook       WebhookConfiguration
	AntiNSFW      AntiNSFWConfiguration
	AntiSpam      AntiSpamConfiguration
	AntiFlood     AntiFloodConfiguration
	AntiUserbot   AntiUserbotConfiguration
	RateLimiter   RateLimiterConfiguration
	FPServer      FPServerConfiguration
	Tdlib         TdlibConfiguration
	AdminBot      AdminBotConfiguration
}
