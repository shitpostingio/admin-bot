package utility

import (
	"testing"
	"time"

	"github.com/shitpostingio/admin-bot/repository"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestUnixTimeIn(t *testing.T) {

	type args struct {
		durationString string
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "5 minutes",
			args: args{"5m"},
			want: time.Now().Add(5 * time.Minute).Unix(),
		},
		{
			name: "3 hours",
			args: args{"3h"},
			want: time.Now().Add(3 * time.Hour).Unix(),
		},
		{
			name: "1 day",
			args: args{"1d"},
			want: time.Now().Add(24 * time.Hour).Unix(),
		},
		{
			name: "1 week",
			args: args{"1w"},
			want: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
		{
			name: "Eternal",
			args: args{"e"},
			want: 0,
		},
		{
			name: "Default 12 hours for duration under 2 characters",
			args: args{"2"},
			want: time.Now().Add(12 * time.Hour).Unix(),
		},
		{
			name: "Default 12 hours for no numbers in the duration",
			args: args{"d"},
			want: time.Now().Add(12 * time.Hour).Unix(),
		},
		{
			name: "Default 12 hours for wrong duration",
			args: args{"3r"},
			want: time.Now().Add(12 * time.Hour).Unix(),
		},
		{
			name: "Default 12 hours for missing duration",
			args: args{""},
			want: time.Now().Add(12 * time.Hour).Unix(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnixTimeIn(tt.args.durationString); got-tt.want >= 100 {
				t.Errorf("UnixTimeIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsChatAdmin(t *testing.T) {

	type args struct {
		msg *tgbotapi.Message
	}

	adminMap := map[int]bool{1: true}
	modMap := map[int]bool{2: true}
	repository.SetAdmins(adminMap)
	repository.SetMods(modMap)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Admin",
			args: args{
				msg: &tgbotapi.Message{From: &tgbotapi.User{ID: 1}},
			},
			want: true,
		},
		{
			name: "Mod",
			args: args{
				msg: &tgbotapi.Message{From: &tgbotapi.User{ID: 1}},
			},
			want: true,
		},
		{
			name: "Nothing",
			args: args{
				msg: &tgbotapi.Message{From: &tgbotapi.User{ID: 3}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChatAdminByMessage(tt.args.msg); got != tt.want {
				t.Errorf("IsChatAdminByMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
