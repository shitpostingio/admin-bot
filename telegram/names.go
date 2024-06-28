package telegram

import (
	"html"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetNameOrUsername returns the user handle or their first and last name
func GetNameOrUsername(user *tgbotapi.User) string {

	if user.UserName != "" {
		return "@" + user.UserName
	}

	if user.LastName == "" {
		return user.FirstName
	}

	return user.FirstName + " " + user.LastName
}

// GetName returns the user's first and last name
func GetName(user *tgbotapi.User) string {

	var output string

	if user.LastName == "" {
		output = user.FirstName
	} else {
		output = user.FirstName + " " + user.LastName
	}

	return html.EscapeString(output)
}
