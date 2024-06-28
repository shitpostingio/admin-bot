package tdlib

import (
	"github.com/pkg/errors"
	"strconv"

	"github.com/shitpostingio/go-tdlib/client"
)

// SendCallback sends a callback response using the Tdlib.
func SendCallback(id, text string, showAlert bool) error {

	queryID, err := strconv.Atoi(id)
	if err != nil {
		return errors.Errorf("SendCallbackWithAlert: unable to parse callback query ID: %s", err)
	}

	_, err = tdlibClient.AnswerCallbackQuery(&client.AnswerCallbackQueryRequest{
		CallbackQueryID: client.JsonInt64(queryID),
		Text:            text,
		ShowAlert:       showAlert,
	})

	return err
}
