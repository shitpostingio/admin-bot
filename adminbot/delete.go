package adminbot

import (
	"github.com/shitpostingio/admin-bot/reports"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api"
)

// DeleteMessage deletes a `tgbotapi.Message`, performing the request in
// a non-blocking, rate-limited way.
func DeleteMessage(chatID int64, messageID int) {
	go func() {
		err := api.DeleteMessage(chatID, messageID)
		if err != nil {
			log.Error("DeleteMessage", err)
		}
	}()
}

func DeleteMultipleMessages(chatID int64, messageIDs ...int) {
	go func() {
		err := api.DeleteMultipleMessages(chatID, messageIDs)
		if err != nil {
			log.Error("DeleteMultipleMessages", err)
		}
	}()
}

// DeleteMessageAndLog deletes a message using `DeleteMessage` and logs the `logText`
func DeleteMessageAndLog(logText string, chatID int64, messageID int) {
	DeleteMessage(chatID, messageID)
	log.Info(logText)
}

// DeleteMessageAndReport deletes a message and logs the result using `DeleteMessageAndLog`.
// the `reportText` is then sent on the report channel.
func DeleteMessageAndReport(reportText string, chatID int64, messageID int) {
	DeleteMessageAndLog(reportText, chatID, messageID)
	_ = reports.Report(reportText, reports.NON_URGENT)
}
