package antispam

import (
	"fmt"
	"time"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/callback/buttons"
	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/defense/antiflood"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/utility"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														STRUCTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

type userActions struct {
	texts int
	media int
	other int
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												CONSTS AND VARS														   *
 *																													   *
 ***********************************************************************************************************************
 */

const (
	/* SPAM TYPES */
	textspam = iota
	mediaspam
	otherspam
	mixspam

	/* SPAM REPORT MESSAGES */
	textSpamReport  = `🚨The user <a href="tg://user?id=%d">%s</a> has been temporarily limited for spamming text messages🚨`
	mediaSpamReport = `🚨The user <a href="tg://user?id=%d">%s</a> has been temporarily limited for spamming media🚨`
	otherSpamReport = `🚨The user <a href="tg://user?id=%d">%s</a> has been temporarily limited for spamming other types of messages🚨`
	mixSpamReport   = `🚨The user <a href="tg://user?id=%d">%s</a> has been temporarily limited for spamming messages🚨`

	/* CHAN SIZES */
	inputChannelSize = 30
)

var (
	/* CHANNELS */
	inputChannel            chan *tgbotapi.Message
	endAntiSpamCycleChannel chan int64
	userChannels            map[int64]chan *tgbotapi.Message

	/* CONFIGURATION */
	cfg *structs.AntiSpamConfiguration
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														START														   *
 *																													   *
 ***********************************************************************************************************************
 */

//Start starts anti spam checks
func Start() {

	cfg = repository.GetAntiSpamConfiguration()
	inputChannel = make(chan *tgbotapi.Message, inputChannelSize)
	endAntiSpamCycleChannel = make(chan int64)
	userChannels = make(map[int64]chan *tgbotapi.Message)

	go detectSpam()
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												USER ROUTINES														   *
 *																													   *
 ***********************************************************************************************************************
 */

//detectSpam receives all incoming messages and sends them to the appropriate routine
func detectSpam() {
	for {
		select {

		case newMessage := <-inputChannel:

			targetChannel, userIsHandled := userChannels[newMessage.From.ID]
			if userIsHandled {
				targetChannel <- newMessage
				continue
			}

			userChannel := make(chan *tgbotapi.Message, inputChannelSize)
			go handleUserActions(userChannel, newMessage.From.ID)
			userChannels[newMessage.From.ID] = userChannel
			userChannel <- newMessage

		case userToRemove := <-endAntiSpamCycleChannel:
			if _, found := userChannels[userToRemove]; found {
				delete(userChannels, userToRemove)
			}
		}
	}
}

//handleUserActions handles the actions of a single user
func handleUserActions(inputChannel <-chan *tgbotapi.Message, userID int64) {

	var actions userActions
	var hasSpammed bool
	var spamType int

	timer := time.After(time.Duration(cfg.RoutineLifeSpan) * time.Second)
	messageIDs := make([]int, 0, 3*cfg.OtherThreshold)

	for {
		select {

		case <-timer:

			endAntiSpamCycleChannel <- userID
			return

		case msg := <-inputChannel:

			if hasSpammed {
				adminbot.DeleteMessage(msg.Chat.ID, msg.MessageID)
				continue
			}

			messageIDs = append(messageIDs, msg.MessageID)

			switch {
			case textMessage(msg):
				actions.texts++
			case mediaMessage(msg):
				actions.media++
			default:
				actions.other++
			}

			hasSpammed, spamType = actions.userIsSpamming()
			if hasSpammed {

				// if the user has spammed, we add 5 minutes to the timer
				timer = time.After(5 * time.Minute)
				antiflood.IncreaseFloodCounter(1)
				_ = adminbot.RestrictMessages(msg.From, msg.Chat.ID, utility.GetAppropriateRestrictionEnd(), &repository.Bot.Self)
				reportSpam(msg.From, spamType)
				adminbot.DeleteMultipleMessages(msg.Chat.ID, messageIDs...)

			}
		}
	}
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												CHANNEL WRAPPERS														   *
 *																													   *
 ***********************************************************************************************************************
 */

// HandleMessage sends a message to the AntiSpam handler.
func HandleMessage(msg *tgbotapi.Message) {
	inputChannel <- msg
}

// EndAntiSpamRoutineForUser ends the antispam routine
// for the input `userID`.
func EndAntiSpamRoutineForUser(userID int64) {
	endAntiSpamCycleChannel <- userID
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													SPAM CHECKS														   *
 *																													   *
 ***********************************************************************************************************************
 */

//userIsSpamming returns true if the user is spamming and the type of spam
func (actions userActions) userIsSpamming() (isSpamming bool, spamType int) {

	/* TEXT SPAM */
	if actions.texts > cfg.TextThreshold {
		return true, textspam
	}

	/* MEDIA SPAM */
	if actions.media > cfg.MediaThreshold {
		return true, mediaspam
	}

	/* OTHER SPAM */
	if actions.other > cfg.OtherThreshold {
		return true, otherspam
	}

	/* MIXED SPAM */
	totalActions := actions.texts + actions.media + actions.other
	if totalActions > 3*cfg.OtherThreshold {
		return true, mixspam
	}

	return false, 0

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														REPORTING													   *
 *																													   *
 ***********************************************************************************************************************
 */

//reportSpam logs and reports the spam to the report channel
func reportSpam(spammer *tgbotapi.User, spamType int) {

	var report string

	switch spamType {
	case textspam:
		report = fmt.Sprintf(textSpamReport, spammer.ID, telegram.GetName(spammer))
	case mediaspam:
		report = fmt.Sprintf(mediaSpamReport, spammer.ID, telegram.GetName(spammer))
	case otherspam:
		report = fmt.Sprintf(otherSpamReport, spammer.ID, telegram.GetName(spammer))
	default:
		report = fmt.Sprintf(mixSpamReport, spammer.ID, telegram.GetName(spammer))
	}

	markup := buttons.CreateKeyboardWithOneRow(buttons.CreateUnrestrictButton(spammer.ID), buttons.CreateHandleButton())
	_ = adminbot.SendTextMessageWithMarkup(repository.GetTelegramConfiguration().ReportChannelID, report, markup, true)
	log.Warn(report)

}

/*
 ***********************************************************************************************************************
 *																													   *
 *													UTILITIES													   	   *
 *																													   *
 ***********************************************************************************************************************
 */

//textMessage returns true if the message is textual
func textMessage(message *tgbotapi.Message) bool {
	return message.Text != ""
}

//mediaMessage returns true if the message contains media
func mediaMessage(message *tgbotapi.Message) bool {
	return message.Photo != nil ||
		message.Video != nil ||
		message.Voice != nil ||
		message.VideoNote != nil ||
		message.Audio != nil
}
