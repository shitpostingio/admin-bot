package automod

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/analysisadapter"
	"github.com/shitpostingio/admin-bot/reports"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf16"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/entities"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/utility"
)

// performChecksOnMessageText performs various messge checks
// and returns true if the message has been deleted.
func performChecksOnMessageText(msg *tgbotapi.Message) bool {

	if checkMessageOrigin(msg) {
		return true
	}

	messageText := telegram.GetMessageText(msg)
	messageTextLength := len(messageText)

	// We can skip additional checks on empty messages.
	if messageTextLength == 0 {
		return false
	}

	if checkMessageLength(messageText, messageTextLength, msg) {
		return true
	}

	// Additional checks need to be performed on the
	// UTF-16 representation of the text.
	tUTF16 := utf16.Encode([]rune(messageText))
	if checkMessageText(tUTF16, msg) {
		return true
	}

	if checkUnwantedHostnames(tUTF16, msg) {
		return true
	}

	return false
}

// checkMessageOrigin checks if a message has been forwarded. If it has, it then
// checks if it has been forwarded from a non-whitelisted channel or from a
// blacklisted handle. If this is the case, the message is deleted and a report
// is sent to the report channel.
func checkMessageOrigin(msg *tgbotapi.Message) bool {

	var reportText string

	// Forwards from channels must be explicitly allowed.
	// In case they aren't, we will try to blacklist the source.
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {

		if database.SourceIsWhitelisted(msg.ForwardFromChat.ID, msg.ForwardFromChat.UserName) {
			return false
		}

		reportText = reports.ForwardFromChannel(msg.From.ID, telegram.GetName(msg.From))
		_, _ = database.BlacklistSource(msg.ForwardFromChat.ID, msg.ForwardFromChat.UserName, msg.From.ID)
		handleMessageDeletionAndReport(reportText, msg)
		return true
	}

	// Remove forwards from blacklisted handles.
	if msg.ForwardFrom != nil {

		//handles for users and bots must be explicitly blacklisted to be removed
		if !database.SourceIsBlacklisted(msg.ForwardFrom.ID, msg.ForwardFrom.UserName) {
			return false
		}

		reportText = reports.ForwardFromBlacklistedHandle(msg.From.ID, telegram.GetName(msg.From))
		handleMessageDeletionAndReport(reportText, msg)
		_ = database.UpdateSource(msg.ForwardFrom.ID, msg.ForwardFrom.UserName)
		return true
	}

	// Remove messages sent via blacklisted inline bots
	if msg.ViaBot != nil {

		if !database.SourceIsBlacklisted(msg.ViaBot.ID, msg.ViaBot.UserName) {
			return false
		}

		reportText = reports.MessageSentViaBlacklistedInlineBot(msg.From.ID, telegram.GetName(msg.From))
		handleMessageDeletionAndReport(reportText, msg)
		_ = database.UpdateSource(msg.ViaBot.ID, msg.ViaBot.UserName)
		return true
	}

	return false

}

func checkMessageLength(text string, textLength int, msg *tgbotapi.Message) bool {

	// Admins are not subject to these limitations
	if repository.Admins[msg.From.ID] {
		return false
	}

	if textLength > 800 || strings.Count(text, "\n") > 15 {
		logText := fmt.Sprintf("Removed long message posted by %s", telegram.GetNameOrUsername(msg.From))
		adminbot.DeleteMessageAndLog(logText, msg.Chat.ID, msg.MessageID)
		return true
	}

	return false

}

// checkMessageText deletes messages over 800 bytes sent by people that are not db admins
// and unwanted handles or links.
func checkMessageText(tUTF16 []uint16, msg *tgbotapi.Message) bool {

	/* UNWANTED HANDLES OR LINKS */
	if messageHasUnwantedHandles(tUTF16, msg) {
		reportText := reports.RemovedUnwantedHandle(msg.From.ID, telegram.GetName(msg.From))
		handleMessageDeletionAndReport(reportText, msg)
		return true
	}

	return false
}

// messageHasUnwantedHandles returns true if the text has unwanted handles.
func messageHasUnwantedHandles(tUTF16 []uint16, msg *tgbotapi.Message) bool {

	handles := telegram.GetAllMentionsUTF16(tUTF16, telegram.GetMessageEntities(msg), msg.ReplyMarkup)
	for _, handle := range handles {

		if handle == "joinchat" {
			return true
		}

		if strings.HasPrefix(handle, "admin") {
			HandleAtAdmin(msg)
			continue
		}

		source, err := database.GetSource(0, handle)
		if err == nil {
			if source.IsWhitelisted {
				continue
			} else {
				return true
			}
		}

		if tdlib.IsGroupOrChannelUsername(handle) {
			_, _ = database.BlacklistSource(0, handle, msg.From.ID)
			return true
		}

		if strings.HasSuffix(handle, "bot") {

			g, err := analysisadapter.GetGibberishValues(handle)
			if err != nil || !g.IsGibberish {
				continue
			}

			bot := repository.Bot
			adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
			_ = adminbot.RestrictMessages(msg.From, msg.Chat.ID, 0, &bot.Self)
			markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnrestrictButton(msg.From.ID), buttons.CreateBanForGibberishHandleButton(msg.From.ID, handle))
			text := reports.UserMutedForUnwantedLink(bot.Self.ID, telegram.GetName(&bot.Self), msg.From.ID, telegram.GetName(msg.From), "posting a gibberish bot handle")
			_ = reports.ReportWithMarkup(text, markup, reports.NON_URGENT)
			//_ = adminbot.BanUser(msg.From, &repository.Bot.Self, "posting a gibberish bot handle", repository.GetTelegramConfiguration().GroupID)
			return true
		}
	}

	return false
}

// checkUnwantedHostnames checks if the message contains links from unwanted hostnames.
// Chat moderators and db admins are immune to these checks.
func checkUnwantedHostnames(tUTF16 []uint16, msg *tgbotapi.Message) bool {

	if !repository.GetTestingStatus() && utility.IsChatAdminByMessage(msg) {
		return false
	}

	urls := telegram.GetURLs(tUTF16, telegram.GetMessageEntities(msg))
	for _, textURL := range urls {

		if !strings.HasPrefix(textURL, "http") {
			textURL = fmt.Sprintf("http://%s", textURL)
		}

		parsedURL, err := url.Parse(textURL)
		if err == nil {

			if strings.Contains(parsedURL.Host, "bitly") ||
				strings.Contains(parsedURL.Host, "bit.ly") ||
				strings.Contains(parsedURL.Host, "tinyurl") {

				data, err := http.Get(textURL) // nolint: gosec
				if err != nil {
					continue
				}

				res, err := ioutil.ReadAll(data.Body)
				if err != nil {
					utility.CloseSafely(data.Body)
					continue
				}

				if strings.Contains(string(res), `toNumbers("f655ba9d09a112d4968c63579db590b4")`) {
					_, _, _ = database.BlacklistHostName(textURL, true, false, repository.Bot.Self.ID)
				}

				utility.CloseSafely(data.Body)

			}

		}

		dbHostname, err := database.GetHostName(textURL)

		// Telegram links have already been checked
		// in previous functions, we can skip them.
		if err != nil || dbHostname.IsTelegram {
			continue
		}

		adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
		punishUserForUnwantedHostname(&dbHostname, msg)
		return true

	}

	return false
}

func punishUserForUnwantedHostname(dbHostname *entities.HostName, msg *tgbotapi.Message) {

	if !dbHostname.IsBanworthy {
		_ = reports.Report(reports.RemovedUnwantedLink(dbHostname.Host, msg.From.ID, telegram.GetName(msg.From)), reports.NON_URGENT)
		return
	}

	motivation := fmt.Sprintf(localization.GetString("automod_unwanted_link_ban_reason"), dbHostname.Host)
	bot := repository.Bot

	if time.Since(documentstore.GetUserJoinDate(msg.From)) <= time.Hour {
		banUserForUnwantedHostname(motivation, dbHostname, msg, bot)
	} else {
		muteUserForUnwantedHostname(motivation, dbHostname, msg, bot)
	}

}

func banUserForUnwantedHostname(motivation string, dbHostname *entities.HostName, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	err := adminbot.BanUser(msg.From, &bot.Self, motivation, msg.Chat.ID)
	if err != nil {

		logText := fmt.Sprintf(localization.GetString("automod_unwanted_link_unable_to_ban"), msg.From.ID, telegram.GetName(msg.From), msg.From.ID, dbHostname.Host, err)
		log.Error(logText)

	}

}

func muteUserForUnwantedHostname(motivation string, dbHostname *entities.HostName, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	err := adminbot.RestrictMessages(msg.From, msg.Chat.ID, 0, &bot.Self)
	if err != nil {

		logText := fmt.Sprintf(localization.GetString("automod_unwanted_link_unable_to_mute"), msg.From.ID, telegram.GetName(msg.From), msg.From.ID, dbHostname.Host, err)
		log.Error(logText)
		return

	}

	reportText := reports.UserMutedForUnwantedLink(bot.Self.ID, telegram.GetName(&bot.Self), msg.From.ID, telegram.GetName(msg.From), motivation)
	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnrestrictButton(msg.From.ID), buttons.CreateHandleButton())
	_ = reports.ReportWithMarkup(reportText, markup, reports.NON_URGENT)

}
