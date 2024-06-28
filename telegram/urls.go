package telegram

import (
	"unicode/utf16"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetURLs returns the urls and the markdown links in a message.
func GetURLs(tUTF16 []uint16, entities []tgbotapi.MessageEntity) (urls []string) {

	if len(entities) == 0 {
		return urls
	}

	for _, entity := range entities {
		switch entity.Type {
		case "url":
			urls = append(urls, string(utf16.Decode(tUTF16[entity.Offset:entity.Offset+entity.Length])))
		case "text_link":
			urls = append(urls, entity.URL)
		}
	}

	return urls
}
