package telegram

import (
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

func TestGetNameOrUsername(t *testing.T) {

	type args struct {
		user *tgbotapi.User
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Handle and name",
			args: args{&tgbotapi.User{FirstName: "Admin", LastName: "Bot", UserName: "AdminBot"}},
			want: "@AdminBot",
		},
		{
			name: "Handle no name",
			args: args{&tgbotapi.User{UserName: "AdminBot"}},
			want: "@AdminBot",
		},
		{
			name: "Name no handle",
			args: args{&tgbotapi.User{FirstName: "Admin", LastName: "Bot"}},
			want: "Admin Bot",
		},
		{
			name: "Nothing",
			args: args{&tgbotapi.User{}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNameOrUsername(tt.args.user); got != tt.want {
				t.Errorf("GetHandleOrUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
