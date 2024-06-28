package callback

import (
	"encoding/json"
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														TESTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

func Test_removeActionButtons(t *testing.T) {

	var msg1, msg2 tgbotapi.Message
	_ = json.Unmarshal([]byte(`{"reply_markup":{"inline_keyboard":[[{"text":"Reported message","url":"https://t.me/shitpost"}],[{"text":"Report","url":"https://t.me/shitpost"},{"text":"Backup","url":"https://t.me/shitpost"}],[{"text":"Mark as handled","callback_data":"2s ok"}]]}}`), &msg1)
	_ = json.Unmarshal([]byte(`{"reply_markup":{"inline_keyboard":[[{"text":"Backup","url":"https://t.me/shitpost"}],[{"text":"Whitelist","callback_data":"2s whitelist 4748 v"},{"text":"Mark as handled","callback_data":"2s ok"}]]}}`), &msg2)
	url := "https://t.me/shitpost"

	type args struct {
		replyMarkup *tgbotapi.InlineKeyboardMarkup
	}
	tests := []struct {
		name string
		args args
		want tgbotapi.InlineKeyboardMarkup
	}{
		{
			name: "@admin report 3 url buttons",
			args: args{replyMarkup: msg1.ReplyMarkup},
			want: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.InlineKeyboardButton{
							Text: "Reported message",
							URL:  &url,
						},
					},
					{
						tgbotapi.InlineKeyboardButton{
							Text: "Report",
							URL:  &url,
						},
						tgbotapi.InlineKeyboardButton{
							Text: "Backup",
							URL:  &url,
						},
					},
				},
			},
		},
		{
			name: "NSFW deletion 1 url button",
			args: args{replyMarkup: msg2.ReplyMarkup},
			want: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.InlineKeyboardButton{
							Text: "Backup",
							URL:  &url,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := removeActionButtons(tt.args.replyMarkup)
			if len(got.InlineKeyboard) != len(tt.want.InlineKeyboard) {
				t.Errorf("Different numbers of rows: want %d, got %d",
					len(tt.want.InlineKeyboard), len(got.InlineKeyboard))
				return
			}

			for rowID, row := range tt.want.InlineKeyboard {
				for columnID := range row {

					if got.InlineKeyboard[rowID][columnID].Text != tt.want.InlineKeyboard[rowID][columnID].Text ||
						*got.InlineKeyboard[rowID][columnID].URL != *tt.want.InlineKeyboard[rowID][columnID].URL {

						t.Errorf("Different items: want: {text: %s, url: %s}, got {text: %s, url: %s}",
							tt.want.InlineKeyboard[rowID][columnID].Text, *tt.want.InlineKeyboard[rowID][columnID].URL,
							got.InlineKeyboard[rowID][columnID].Text, *got.InlineKeyboard[rowID][columnID].URL)
					}
				}
			}
		})
	}
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													BENCHMARKS														   *
 *																													   *
 ***********************************************************************************************************************
 */

func Benchmark_getOriginalMessageText(b *testing.B) {

	var msg tgbotapi.Message
	err := json.Unmarshal([]byte(`{"text":"� Alessandro Pomponio banned Telegram Bot Raw.\nFor: test","entities":[{"type":"text_mention","offset":3,"length":19,"url":"","user":{"id":56800135,"first_name":"Alessandro","last_name":"Pomponio","username":"AlessandroPomponio","language_code":"it","is_bot":false}},{"type":"text_mention","offset":30,"length":16,"url":"","user":{"id":211246197,"first_name":"Telegram Bot Raw","last_name":"","username":"RawDataBot","language_code":"","is_bot":true}}]}`), &msg)
	if err != nil {
		b.Error("unable to parse message", err)
	}

	for n := 0; n < b.N; n++ {
		getOriginalMessageText(&msg)
	}
}
