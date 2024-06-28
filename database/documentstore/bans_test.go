package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shitpostingio/admin-bot/entities"
)

func TestAddBan(t *testing.T) {

	type args struct {
		bannedUserID    int
		moderatorUserID int
		reason          string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "AddBanWithReason - No Error",
			args: args{
				bannedUserID:    1,
				moderatorUserID: 2,
				reason:          "this is a test",
			},
		},
		{
			name: "AddBanWithoutReason - No Error",
			args: args{
				bannedUserID:    1,
				moderatorUserID: 2,
				reason:          "",
			},
		},
		{
			name: "AddBanBannedUserID 0 - Error",
			args: args{
				bannedUserID:    0,
				moderatorUserID: 2,
				reason:          "this is a test",
			},
			wantErr: true,
		},
		{
			name: "AddBanModeratorID 0 - Error",
			args: args{
				bannedUserID:    1,
				moderatorUserID: 0,
				reason:          "this is a test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			_, err := AddBan(tt.args.bannedUserID, tt.args.moderatorUserID, tt.args.reason, testCollection)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetBansForTelegramID(t *testing.T) {

	tests := []struct {
		name       string
		telegramID int
		wantBans   []entities.Ban
		wantErr    bool
	}{
		{
			name:       "UserID 1 - Two bans",
			telegramID: 1,
			wantBans: []entities.Ban{
				{
					User:     1,
					BannedBy: 2,
					Reason:   "",
				},
				{
					User:     1,
					BannedBy: 2,
					Reason:   "this is a test",
				},
			},
		},
		{
			name:       "UserID 2 - No bans",
			telegramID: 2,
			wantBans:   []entities.Ban{},
		},
		{
			name:       "UserID 0 - Error",
			telegramID: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotBans, err := GetBansForTelegramID(tt.telegramID, testCollection)

			// error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// slice size
			assert.Len(t, gotBans, len(tt.wantBans))

			// ban content
			for i, ban := range gotBans {
				assert.Equal(t, ban.User, tt.wantBans[i].User)
				assert.Equal(t, ban.BannedBy, tt.wantBans[i].BannedBy)
				assert.Equal(t, ban.Reason, tt.wantBans[i].Reason)
				assert.Equal(t, ban.UnbanDate, tt.wantBans[i].UnbanDate)
			}

		})
	}
}

func TestMarkUserAsUnbanned(t *testing.T) {

	tests := []struct {
		name       string
		telegramID int
		wantErr    bool
	}{
		{
			name:       "UserID 1 - No error",
			telegramID: 1,
		},
		{
			name:       "UserID 2 - No matches",
			telegramID: 2,
			wantErr:    true,
		},
		{
			name:       "UserID 3 - Not found",
			telegramID: 3,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MarkUserAsUnbanned(tt.telegramID, testCollection); (err != nil) != tt.wantErr {
				t.Errorf("MarkUserAsUnbanned() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
