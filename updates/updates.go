package updates

import (
	"net/http"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/config/structs"
)

// GetUpdatesChannel contacts Telegram's servers in order to get updates.
// Updates can be received either via polling or via webhooks.
func GetUpdatesChannel(connectViaPolling bool, bot *tgbotapi.BotAPI, cfg *structs.Config) tgbotapi.UpdatesChannel {

	if !connectViaPolling {
		return useWebhook(bot, cfg)
	}

	_, err := bot.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Error("GetUpdatesChannel: unable to remove webhook:", err)
		return nil
	}

	return usePolling(bot)

}

// usePolling gets an `UpdatesChannel` using polling
func usePolling(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	return bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
}

//useWebhook ets an `UpdatesChannel` using webhooks
func useWebhook(bot *tgbotapi.BotAPI, cfg *structs.Config) tgbotapi.UpdatesChannel {

	go startServer(&cfg.Webhook)

	webhook, err := bot.GetWebhookInfo()
	if err == nil && webhook.IsSet() {

		// A webhook has already been set: we need to make sure
		// it points to the correct domain.
		domainNameStart := strings.Index(webhook.URL, "/") + 2

		// The webhook points to the correct domain
		if strings.HasPrefix(webhook.URL[domainNameStart:], cfg.Webhook.Domain) {
			return bot.ListenForWebhook(cfg.Telegram.WebHookPath())
		}

	}

	// Set up new webhooks
	webhookConfiguration, err := tgbotapi.NewWebhook(cfg.Webhook.WebHookURL(cfg.Telegram.BotToken))
	if err != nil {
		log.Error("useWebhook: unable to create webhook configuration:", err)
		return nil
	}

	webhookConfiguration.MaxConnections = cfg.Webhook.MaxConnections
	_, err = bot.Request(webhookConfiguration)
	if err != nil {
		log.Error("useWebhook: unable to request webhook creation:", err)
		return nil
	}

	return bot.ListenForWebhook(cfg.Telegram.WebHookPath())

}

//startServer starts serving HTTP requests with or without TLS
func startServer(cfg *structs.WebhookConfiguration) {

	if cfg.TLS {
		log.Error(http.ListenAndServeTLS(cfg.BindString(), cfg.TLSCertPath, cfg.TLSKeyPath, nil))
	} else {
		log.Error(http.ListenAndServe(cfg.BindString(), nil))
	}

}
