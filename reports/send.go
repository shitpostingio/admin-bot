package reports

import (
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/repository"
)

//nolint
const (
	URGENT     = true
	NON_URGENT = false
)

func ReportWithMarkup(message string, markup interface{}, isUrgent bool) error {
	_, err := api.SendTextMessageWithMarkup(repository.GetTelegramConfiguration().ReportChannelID,
		message,
		markup,
		isUrgent)
	return err
}

func ReportInPlaintext(message string, isUrgent bool) error {
	_, err := api.SendPlainTextMessage(repository.GetTelegramConfiguration().ReportChannelID,
		message,
		isUrgent)
	return err
}

func Report(message string, isUrgent bool) error {
	_, err := api.SendTextMessage(repository.GetTelegramConfiguration().ReportChannelID,
		message,
		isUrgent)
	return err
}
