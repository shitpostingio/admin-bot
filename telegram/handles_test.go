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
 *														SET UP														   *
 *																													   *
 ***********************************************************************************************************************
 */

/*
 ***********************************************************************************************************************
 *																													   *
 *														TESTS														   *
 *																													   *
 ***********************************************************************************************************************
 */
//
//func TestGetAllMentions(t *testing.T) {
//
//	db := SetupTests(false)
//	text := `Hey shitposters, you should really join the channels in the shitposting.io network, like @Sushiporn, they're really nice and most definitely the best ones on @teleGram!!!`
//	entities := []tgbotapi.MessageEntity{
//		{
//			Offset: 4,
//			Length: 11,
//			Type:   "text_link",
//			URL:    "http://t.me/shitpost",
//		},
//		{
//			Offset: 60,
//			Length: 14,
//			Type:   "url",
//		},
//		{
//			Offset: 89,
//			Length: 10,
//			Type:   "mention",
//		},
//		{
//			Offset: 158,
//			Length: 9,
//			Type:   "mention",
//		},
//	}
//
//	type args struct {
//		entities []tgbotapi.MessageEntity
//		text     string
//	}
//	tests := []struct {
//		name        string
//		args        args
//		wantHandles []string
//	}{
//		{
//			name: "TestGetAllHandles",
//			args: args{
//				entities: entities,
//				text:     text,
//				db:       db,
//			},
//			wantHandles: []string{"shitpost", "sushiporn", "telegram"},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if gotHandles := GetAllMentions(tt.args.text, tt.args.entities, nil, tt.args.db); !reflect.DeepEqual(gotHandles, tt.wantHandles) {
//				t.Errorf("GetAllMentions() = %v, want %v", gotHandles, tt.wantHandles)
//			}
//		})
//	}
//}

func TestGetMentions(t *testing.T) {

	text := "You guys should really join @shitpOst and @sushiporn, the best channels on @teLegram"
	entities := []tgbotapi.MessageEntity{
		{
			Offset: 28,
			Length: 9,
			Type:   "mention",
		},
		{
			Offset: 42,
			Length: 10,
			Type:   "mention",
		},
		{
			Offset: 75,
			Length: 9,
			Type:   "mention",
		},
	}
	tUTF16 := utf16.Encode([]rune(text))

	type args struct {
		tUTF16   []uint16
		entities []tgbotapi.MessageEntity
	}

	tests := []struct {
		name        string
		args        args
		wantHandles []string
	}{
		{
			name: "TestGetHandles",
			args: args{
				tUTF16:   tUTF16,
				entities: entities,
			},
			wantHandles: []string{"shitpost", "sushiporn", "telegram"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHandles := GetMentions(tt.args.tUTF16, tt.args.entities); !reflect.DeepEqual(gotHandles, tt.wantHandles) {
				t.Errorf("GetMentions() = %v, want %v", gotHandles, tt.wantHandles)
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

//func BenchmarkGetAllMentions(b *testing.B) {
//
//	db := SetupTests(false)
//	dbReply := []map[string]interface{}{{"id": 1, "hostname": "t.me", "is_banworthy": false, "is_telegram": true}}
//	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "hostnames"  WHERE (hostname = t.me)`).WithReply(dbReply)
//
//	text := `Hey shitposters, you should really join the channels in the shitposting.io network, like @sushiporn, they're really nice and most definitely the best ones on @telegram!!!`
//	entities := []tgbotapi.MessageEntity{
//		{
//			Offset: 4,
//			Length: 11,
//			Type:   "text_link",
//			URL:    "http://t.me/shitpost",
//		},
//		{
//			Offset: 60,
//			Length: 14,
//			Type:   "url",
//		},
//		{
//			Offset: 89,
//			Length: 10,
//			Type:   "mention",
//		},
//		{
//			Offset: 158,
//			Length: 9,
//			Type:   "mention",
//		},
//	}
//
//	for i := 0; i < b.N; i++ {
//		GetAllMentions(text, entities, nil, db)
//	}
//}

func BenchmarkGetMentions(b *testing.B) {
	text := "You guys should really join @shitpost and @sushiporn, the best channels on @telegram"
	entities := []tgbotapi.MessageEntity{
		{
			Offset: 28,
			Length: 9,
			Type:   "mention",
		},
		{
			Offset: 42,
			Length: 10,
			Type:   "mention",
		},
		{
			Offset: 75,
			Length: 9,
			Type:   "mention",
		},
	}
	tUTF16 := utf16.Encode([]rune(text))

	for i := 0; i < b.N; i++ {
		GetMentions(tUTF16, entities)
	}
}
