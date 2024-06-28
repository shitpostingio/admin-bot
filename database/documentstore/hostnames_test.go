package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shitpostingio/admin-bot/entities"
)

func TestAddHostNameToBlacklist(t *testing.T) {

	type args struct {
		url             string
		isBanworthy     bool
		isTelegram      bool
		adderTelegramID int
	}

	tests := []struct {
		name     string
		args     args
		wantHost string
		wantErr  bool
	}{
		{
			name: "URL https, no banworthy, no telegram, no telegram id",
			args: args{
				url:             "https://docs.mongodb.com/manual/core/index-case-insensitive/",
				isBanworthy:     false,
				isTelegram:      false,
				adderTelegramID: 0,
			},
			wantHost: "docs.mongodb.com",
		},
		{
			name: "URL https, yes banworthy, no telegram, yes telegram id",
			args: args{
				url:             "https://github.com/mongodb/mongo-go-driver/blob/master/examples/documentation_examples/examples.go",
				isBanworthy:     true,
				isTelegram:      false,
				adderTelegramID: 10,
			},
			wantHost: "github.com",
		},
		{
			name: "URL http, no banworthy, yes telegram, yes telegram id",
			args: args{
				url:             "http://t.me/shitpost",
				isBanworthy:     false,
				isTelegram:      true,
				adderTelegramID: 10,
			},
			wantHost: "t.me",
		},
		{
			name: "URL no protocol, no banworthy, yes telegram, yes telegram id",
			args: args{
				url:             "telegram.dog/shitpost",
				isBanworthy:     false,
				isTelegram:      true,
				adderTelegramID: 10,
			},
			wantHost: "telegram.dog",
		},
		{
			name: "Empty url",
			args: args{
				url:             "",
				isBanworthy:     false,
				isTelegram:      false,
				adderTelegramID: 10,
			},
			wantErr: true,
		},
		{
			name: "Not a url",
			args: args{
				url:             "ciao",
				isBanworthy:     false,
				isTelegram:      false,
				adderTelegramID: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, gotHost, err := BlacklistHostName(tt.args.url, tt.args.isBanworthy, tt.args.isTelegram, tt.args.adderTelegramID, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, gotHost, tt.wantHost)

		})
	}
}

func TestUpdateHostName(t *testing.T) {

	type args struct {
		url               string
		isBanworthy       bool
		isTelegram        bool
		updaterTelegramID int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "docs.mongodb.com - all true no error",
			args: args{
				url:               "https://docs.mongodb.com/manual/core/index-case-insensitive/",
				isBanworthy:       true,
				isTelegram:        true,
				updaterTelegramID: 13,
			},
		},
		{
			name: "github.com - all false no error",
			args: args{
				url:               "https://github.com/mongodb/mongo-go-driver/blob/master/examples/documentation_examples/examples.go",
				isBanworthy:       false,
				isTelegram:        false,
				updaterTelegramID: 10,
			},
		},
		{
			name: "telegram.me - url not found",
			args: args{
				url:               "http://telegram.me/shitpost",
				isBanworthy:       false,
				isTelegram:        true,
				updaterTelegramID: 10,
			},
			wantErr: true,
		},
		{
			name: "Empty url",
			args: args{
				url:               "",
				isBanworthy:       false,
				isTelegram:        false,
				updaterTelegramID: 10,
			},
			wantErr: true,
		},
		{
			name: "Not a url",
			args: args{
				url:               "ciao",
				isBanworthy:       false,
				isTelegram:        false,
				updaterTelegramID: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := UpdateHostName(tt.args.url, tt.args.isBanworthy, tt.args.isTelegram, tt.args.updaterTelegramID, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestGetHostName(t *testing.T) {

	tests := []struct {
		name     string
		url      string
		wantHost entities.HostName
		wantErr  bool
	}{
		{
			name: "docs.mongodb.com - no error",
			url:  "docs.mongodb.org",
			wantHost: entities.HostName{
				Host:         "docs.mongodb.com",
				IsBanworthy:  true,
				IsTelegram:   true,
				LastEditedBy: 13,
			},
		},
		{
			name:    "docs.mongodb.com - no match",
			url:     "mongodb.org",
			wantErr: true,
		},
		{
			name:    "empty host",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotHost, err := GetHostName(tt.url, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantHost.Host, gotHost.Host)
			assert.Equal(t, tt.wantHost.IsBanworthy, gotHost.IsBanworthy)
			assert.Equal(t, tt.wantHost.IsTelegram, gotHost.IsTelegram)
			assert.Equal(t, tt.wantHost.LastEditedBy, gotHost.LastEditedBy)

		})
	}
}

func TestRemoveHostNameFromBlacklist(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name: "docs.mongodb.com - no error",
			url:  "docs.mongodb.org",
		},
		{
			name:    "google.com - not found",
			url:     "google.com",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := PardonHostName(tt.url, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
