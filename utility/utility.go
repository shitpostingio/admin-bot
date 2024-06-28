package utility

import (
	"io"
	"os"
	"strconv"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/repository"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//UnixTimeIn parses a string and returns Time.Now() + the parsed duration in UNIX time
func UnixTimeIn(durationString string) int64 {

	//"e" MEANS FOREVER
	if durationString == "e" {
		return 0
	}

	/* DEFAULT DURATION */
	restrictionDuration := 43200 //12hrs

	//SUPPORTED DURATIONS: w (weeks), d (days), h (hours), m (minutes)
	for index, durationRune := range durationString {

		//WE ONLY SUPPORT SINGLE-RUNE DURATIONS
		if !unicode.IsNumber(durationRune) {
			number, err := strconv.Atoi(durationString[0:index])
			if err != nil {
				break
			}

			switch durationRune {
			case 'w':
				restrictionDuration = number * 604800 //weeks
			case 'd':
				restrictionDuration = number * 86400 //days
			case 'h':
				restrictionDuration = number * 3600 //hours
			case 'm':
				restrictionDuration = number * 60 //minutes
			}
		}
	}

	return time.Now().Unix() + int64(restrictionDuration)
}

//IsChatAdminByMessage returns true if the user is an admin or the creator
func IsChatAdminByMessage(msg *tgbotapi.Message) bool {
	return repository.Admins[msg.From.ID] || repository.Mods[msg.From.ID]
}

// IsChatAdmin returns true if the telegram id belongs to an admin or a mod
func IsChatAdmin(telegramID int64) bool {
	return repository.Admins[telegramID] || repository.Mods[telegramID]
}

//FormatDate formats a date
func FormatDate(date time.Time) string {
	return date.Format("Mon _2 Jan 2006 15:04:05")
}

//FormatUnixDate formats a date in unix format
func FormatUnixDate(unixDate int64) string {

	if unixDate == 0 {
		return "indefinitely"
	}

	return time.Unix(unixDate, 0).Format("Mon _2 Jan 2006 15:04:05")
}

//EmojifyBool returns ✅ for true and ❌ for false
func EmojifyBool(value bool) string {
	if value {
		return "✅"
	}

	return "❌"
}

//CloseSafely closes an entity and logs in case of errors
func CloseSafely(toClose io.Closer) {
	err := toClose.Close()
	if err != nil {
		log.Println(err)
	}
}

// LogFatal logs an error and exits
func LogFatal(format ...interface{}) {
	log.Error(format...)
	os.Exit(1)
}
