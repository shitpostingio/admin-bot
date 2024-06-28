package telegram

import (
	"reflect"
	"testing"
	"unicode/utf16"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														TESTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

func TestGetURLs(t *testing.T) {
	type args struct {
		entities []tgbotapi.MessageEntity
		text     string
	}
	tests := []struct {
		name     string
		args     args
		wantURLs []string
	}{
		{
			name: "Url no http and markdown URL",
			args: args{
				text: "Hi! Check out t.me/shitpost, it's a great channel! You can check out sushiporn too, if you'd like",
				entities: []tgbotapi.MessageEntity{
					{
						Offset: 14,
						Length: 13,
						Type:   "url",
					},
					{
						Offset: 69,
						Length: 9,
						Type:   "text_link",
						URL:    "http://t.me/sushiporn",
					},
				},
			},
			wantURLs: []string{"t.me/shitpost", "http://t.me/sushiporn"},
		},
		{
			name: "Shortened url and google no http",
			args: args{
				text: "check out the shitposting website and follow us on google.com",
				entities: []tgbotapi.MessageEntity{
					{
						Offset: 14,
						Length: 11,
						Type:   "text_link",
						URL:    "https://bit.ly/2WA4TBH",
					},
					{
						Offset: 51,
						Length: 10,
						Type:   "url",
					},
				},
			},
			wantURLs: []string{"https://bit.ly/2WA4TBH", "google.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHandles := GetURLs(utf16.Encode([]rune(tt.args.text)), tt.args.entities); !reflect.DeepEqual(gotHandles, tt.wantURLs) {
				t.Errorf("NewFindURLs() = %v, want %v", gotHandles, tt.wantURLs)
			}
		})
	}
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														BENCHMARKS														   *
 *																													   *
 ***********************************************************************************************************************
 */

func BenchmarkNewFindURLs(b *testing.B) {
	text := "Hi! Check out t.me/shitpost, it's a great channel! You can check out sushiporn too, if you'd like"
	entities := []tgbotapi.MessageEntity{
		{
			Offset: 14,
			Length: 13,
			Type:   "url",
		},
		{
			Offset: 69,
			Length: 9,
			Type:   "text_link",
			URL:    "http://t.me/sushiporn",
		},
	}
	tUTF16 := utf16.Encode([]rune(text))

	for i := 0; i < b.N; i++ {
		GetURLs(tUTF16, entities)
	}
}
