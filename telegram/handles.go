package telegram

import (
	"net/url"
	"strings"
	"unicode/utf16"

	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/utility"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetAllMentions returns all handles in a message.
// It'll return explicit handles in the form of @mentions,
// but also handles extracted from telegram links in the form
// of <telegramdomain>/<handle>.
func GetAllMentions(text string, entities []tgbotapi.MessageEntity, markup *tgbotapi.InlineKeyboardMarkup) (handles []string) {

	if len(entities) == 0 && markup == nil {
		return handles
	}

	tUTF16 := utf16.Encode([]rune(text))
	return GetAllMentionsUTF16(tUTF16, entities, markup)
}

// GetAllMentionsUTF16 returns all handles in a message.
// It'll return explicit handles in the form of @mentions,
// but also handles extracted from telegram links in the form
// of <telegramdomain>/<handle>.
func GetAllMentionsUTF16(tUTF16 []uint16, entities []tgbotapi.MessageEntity, markup *tgbotapi.InlineKeyboardMarkup) (handles []string) {

	if len(entities) == 0 && markup == nil {
		return handles
	}

	urls := GetURLs(tUTF16, entities)
	for _, textURL := range urls {

		fullTextURL, err := utility.UnshortenURL(textURL)
		if err != nil {
			continue
		}

		fullTextURLLowercase := strings.ToLower(fullTextURL)
		parsedURL, err := url.Parse(fullTextURLLowercase)
		if err != nil {
			continue
		}

		// Extract handles from known telegram URLs.
		dbHostname, err := database.GetHostName(fullTextURLLowercase)
		if err == nil && dbHostname.IsTelegram {
			//parsedURL.Path WILL RETURN EVERYTHING AFTER THE HOST
			//THAT'S NOT PART OF THE QUERY
			parts := strings.SplitN(parsedURL.Path, "/", 3)
			if len(parts) >= 2 {
				handles = append(handles, parts[1])
			}
		}
	}

	messageHandles := GetMentions(tUTF16, entities)
	handles = append(handles, messageHandles...)
	markupHandles := GetInlineKeyboardMentions(markup)
	handles = append(handles, markupHandles...)

	return handles
}

// GetMentions returns the mentions in a message.
func GetMentions(tUTF16 []uint16, entities []tgbotapi.MessageEntity) (handles []string) {

	if len(entities) == 0 {
		return handles
	}

	for _, entity := range entities {
		if entity.Type == "mention" {
			handle := strings.ToLower(string(utf16.Decode(tUTF16[entity.Offset+1 : entity.Offset+entity.Length])))
			handles = append(handles, handle)
		}
	}

	return handles
}

// GetInlineKeyboardMentions returns the mentions contained in an inline keyboard
func GetInlineKeyboardMentions(markup *tgbotapi.InlineKeyboardMarkup) (handles []string) {

	if markup == nil {
		return handles
	}

	for _, rows := range markup.InlineKeyboard {

		for _, column := range rows {

			if column.URL != nil {

				fullTextURL, err := utility.UnshortenURL(*column.URL)
				if err != nil {
					continue
				}

				fullTextURLLowercase := strings.ToLower(fullTextURL)
				parsedURL, err := url.Parse(fullTextURLLowercase)
				if err != nil {
					continue
				}

				// Extract handles from known telegram URLs.
				dbHostname, err := database.GetHostName(fullTextURLLowercase)
				if err == nil && dbHostname.IsTelegram {
					//parsedURL.Path WILL RETURN EVERYTHING AFTER THE HOST
					//THAT'S NOT PART OF THE QUERY
					parts := strings.SplitN(parsedURL.Path, "/", 3)
					if len(parts) >= 2 {
						handles = append(handles, parts[1])
					}
				}

			}

		}

	}

	return handles

}
