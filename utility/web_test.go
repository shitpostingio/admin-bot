package utility

import (
	"testing"
)

func TestUnshortenURL(t *testing.T) {
	tests := []struct {
		name               string
		URL                string
		wantUnshortenedURL string
		wantErr            bool
	}{
		{
			name:               "Shortened URL",
			URL:                "https://bit.ly/2WA4TBH",
			wantUnshortenedURL: "https://shitposting.io/",
			wantErr:            false,
		},
		{
			name:               "Unshortened URL",
			URL:                "https://shitposting.io/",
			wantUnshortenedURL: "https://shitposting.io/",
			wantErr:            false,
		},
		{
			name:               "Telegram url no http",
			URL:                "t.me/shitpost",
			wantUnshortenedURL: "https://t.me/shitpost",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUnshortenedURL, err := UnshortenURL(tt.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnshortenURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUnshortenedURL != tt.wantUnshortenedURL {
				t.Errorf("UnshortenURL() = %v, want %v", gotUnshortenedURL, tt.wantUnshortenedURL)
			}
		})
	}
}

func TestIsGroupOrChannelHandle(t *testing.T) {
	tests := []struct {
		name   string
		handle string
		want   bool
	}{
		{
			name:   "Channel",
			handle: "shitpost",
			want:   true,
		},

		{
			name:   "Group",
			handle: "thememaly",
			want:   true,
		},

		{
			name:   "User",
			handle: "emaele_",
			want:   false,
		},

		{
			name:   "Bot",
			handle: "levelinebot",
			want:   false,
		},

		{
			name:   "Invalid Handle",
			handle: "rijfioerjgiojerigjoiejgoijeoigrj",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGroupOrChannelHandle(tt.handle); got != tt.want {
				t.Errorf("IsGroupOrChannelHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
