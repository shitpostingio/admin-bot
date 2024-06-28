package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//GetMessageEntities returns an array of entities. It can be message.Entities,
//message.CaptionEntities or an empty array.
func GetMessageEntities(msg *tgbotapi.Message) (entities []tgbotapi.MessageEntity) {

	if msg.Entities != nil {
		entities = msg.Entities
	} else if msg.CaptionEntities != nil {
		entities = msg.CaptionEntities
	}

	return entities
}

//GetMessageText returns msg.Text or msg.Caption
func GetMessageText(msg *tgbotapi.Message) (text string) {

	if msg.Text != "" {
		text = msg.Text
	} else {
		text = msg.Caption
	}

	return text
}

// GetFileIDFromMessage returns the file id given a message
func GetFileIDFromMessage(msg *tgbotapi.Message) (fileUniqueID, fileID string) {

	switch {
	case msg.Photo != nil:

		fileUniqueID = msg.Photo[len(msg.Photo)-1].FileUniqueID
		fileID = msg.Photo[len(msg.Photo)-1].FileID

	case msg.Video != nil:

		fileUniqueID = msg.Video.FileUniqueID
		fileID = msg.Video.FileID

	case msg.Sticker != nil:

		fileUniqueID = msg.Sticker.FileUniqueID
		fileID = msg.Sticker.FileID

	case msg.Animation != nil:

		fileUniqueID = msg.Animation.FileUniqueID // TODO: RIMUOVERE QUANDO SYFARO AGGIUNGERA' IL CAMPO
		fileID = msg.Animation.FileID

	case msg.Document != nil:

		fileUniqueID = msg.Document.FileUniqueID
		fileID = msg.Document.FileID

	case msg.Voice != nil:

		fileUniqueID = msg.Voice.FileUniqueID
		fileID = msg.Voice.FileID

	case msg.Audio != nil:

		fileUniqueID = msg.Audio.FileUniqueID
		fileID = msg.Audio.FileID

	case msg.VideoNote != nil:

		fileUniqueID = msg.VideoNote.FileUniqueID
		fileID = msg.VideoNote.FileID

	}

	return fileUniqueID, fileID
}
