package adminbot

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/reports"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"
)

const (
	groupNotAllowed   = "I'm sorry but @%s is not allowed to be in this group.\nTo have me here contact us: @shitpost"
	channelNotAllowed = "I'm sorry but @%s is not allowed to be in this channel.\nTo have me here contact us: @shitpost"
)

// LeaveUnauthorizedGroup makes the bot leave a group it wasn't authorized to be in.
// Before doing so, the bot will send a message saying it's not authorized
// and link to the channel, so that the users can contact us.
// It will also log where it has been added, so we can check.
func LeaveUnauthorizedGroup(msg *tgbotapi.Message) {

	groupNotAllowedText := fmt.Sprintf(groupNotAllowed, repository.Bot.Self.UserName)
	_ = SendPlainTextMessage(msg.Chat.ID, groupNotAllowedText, false)

	response, err := api.LeaveChat(msg.Chat.ID)
	if err != nil {

		report := fmt.Sprintf("Unable to leave unauthorized group with ID %d", msg.Chat.ID)

		if response != nil {
			log.Error(report, "(", response.ErrorCode, ":", response.Description, ")")
		} else {
			log.Error(report)
		}

		return

	}

	var reportText string
	if msg.Chat.UserName != "" {
		reportText = reports.UnauthorizedPublicGroup(msg.Chat.Title, msg.Chat.UserName, msg.Chat.ID,
			telegram.GetName(msg.From), msg.From.ID)
	} else {
		reportText = reports.UnauthorizedPrivateGroup(msg.Chat.Title, msg.Chat.ID,
			telegram.GetName(msg.From), msg.From.ID)
	}

	_ = reports.Report(reportText, reports.NON_URGENT)
	log.Warn(reportText)
}

// LeaveUnauthorizedChannel makes the bot leave a channel it wasn't authorized to be in.
// Before doing so, the bot will send a message saying it's not authorized
// and link to the channel, so that the users can contact us.
// It will also log where it has been added, so we can check.
func LeaveUnauthorizedChannel(post *tgbotapi.Message) {

	channelNotAllowedText := fmt.Sprintf(channelNotAllowed, repository.Bot.Self.UserName)
	_ = SendPlainTextMessage(post.Chat.ID, channelNotAllowedText, false)

	response, err := api.LeaveChat(post.Chat.ID)
	if err != nil {

		report := fmt.Sprintf("Unable to leave unauthorized channel with ID %d", post.Chat.ID)

		if response != nil {
			log.Error(report, "(", response.ErrorCode, ":", response.Description, ")")
		} else {
			log.Error(report)
		}

		return

	}

	var reportText string
	if post.Chat.UserName != "" {
		reportText = reports.UnauthorizedPublicChannel(post.Chat.Title, post.Chat.UserName, post.Chat.ID)
	} else {
		reportText = reports.UnauthorizedPrivateChannel(post.Chat.Title, post.Chat.ID)
	}

	_ = reports.Report(reportText, reports.NON_URGENT)
	log.Warn(reportText)

}
