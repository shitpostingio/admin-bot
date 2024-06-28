package documentstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/shitpostingio/admin-bot/entities"
)

func TestBlacklistMedia(t *testing.T) {

	type args struct {
		fileID string
		userID int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GIF - No error",
			args: args{
				fileID: "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
				userID: 11,
			},
		},
		{
			name: "photo - no error",
			args: args{
				fileID: "AgADBAAD7rYxG3vG4VMRKAh8HpBPUh0mthsABAEAAwIAA20AA6KhAAIWBA",
				userID: 11,
			},
		},
		{
			name: "video - no error",
			args: args{
				fileID: "BAADAQADjQADBSjQR5EcvphDJJnwFgQ",
				userID: 11,
			},
		},
		{
			name: "document - no error",
			args: args{
				fileID: "BQADBAADtQYAAnvG4VNbb4oMumoR4xYE",
				userID: 11,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BlacklistMedia(tt.args.fileID, tt.args.userID, testCollection)
			if (err != nil) != tt.wantErr {
				t.Errorf("BlacklistMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFindMediaByFileID(t *testing.T) {

	tests := []struct {
		name      string
		fileID    string
		wantMedia entities.Media
		wantErr   bool
	}{
		{
			name:   "GIF - No error",
			fileID: "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
			wantMedia: entities.Media{
				FileID:           "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
				Histogram:        []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 95, 1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
				HistogramAverage: 37.25,
				HistogramSum:     1192,
				PHash:            "p:c3320afcd36f3c81",
			},
		},
		{
			name:    "File not found - Error, no fileID match",
			fileID:  "aCgADBAAD6gEAAop",
			wantErr: true,
		},
		{
			name:    "File not found - Error, unable to get fingerprint",
			fileID:  "CgADBAAD6gEAAop",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMedia, err := FindMediaByFileID(tt.fileID, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantMedia.FileID, gotMedia.FileID)
			assert.Equal(t, tt.wantMedia.Histogram, gotMedia.Histogram)
			assert.Equal(t, tt.wantMedia.HistogramAverage, gotMedia.HistogramAverage)
			assert.Equal(t, tt.wantMedia.HistogramSum, gotMedia.HistogramSum)
			assert.Equal(t, tt.wantMedia.PHash, gotMedia.PHash)

		})
	}
}

func TestFindMediaByFeatures(t *testing.T) {

	type args struct {
		histogram     []float64
		pHash         string
		approximation float64
	}

	tests := []struct {
		name      string
		args      args
		wantMedia entities.Media
		wantErr   bool
	}{
		{
			name: "GIF - No error",
			args: args{
				histogram:     []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 95, 1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
				pHash:         "p:c3320afcd36f3c81",
				approximation: 0.04,
			},
			wantMedia: entities.Media{
				FileID:           "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
				Histogram:        []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 95, 1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
				HistogramAverage: 37.25,
				HistogramSum:     1192,
				PHash:            "p:c3320afcd36f3c81",
			},
		},
		{
			name: "No match - histogram nil",
			args: args{
				pHash:         "p:c3320afcd36f3c81",
				approximation: 0.04,
			},
			wantErr: true,
		},
		{
			name: "No match - no phash",
			args: args{
				histogram:     []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 95, 1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
				approximation: 0.04,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMedia, err := FindMediaByFeatures(tt.args.histogram, tt.args.pHash, tt.args.approximation, testCollection)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantMedia.FileID, gotMedia.FileID)
			assert.Equal(t, tt.wantMedia.Histogram, gotMedia.Histogram)
			assert.Equal(t, tt.wantMedia.HistogramAverage, gotMedia.HistogramAverage)
			assert.Equal(t, tt.wantMedia.HistogramSum, gotMedia.HistogramSum)
			assert.Equal(t, tt.wantMedia.PHash, gotMedia.PHash)

		})
	}
}

func TestFindMediaByID(t *testing.T) {

	tests := []struct {
		name      string
		fileID    string
		wantMedia entities.Media
		wantErr   bool
	}{
		{
			name:   "GIF - No error",
			fileID: "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
			wantMedia: entities.Media{
				FileID:           "CgADBAAD6gEAAop-_VJFQHZLB74cdBYE",
				Histogram:        []float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 95, 1, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0},
				HistogramAverage: 37.25,
				HistogramSum:     1192,
				PHash:            "p:c3320afcd36f3c81",
			},
		},
		{
			name:    "No fileID - Error",
			fileID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ID := &primitive.ObjectID{}
			ID = nil

			if tt.fileID != "" {
				wantMedia, err := FindMediaByFileID(tt.fileID, testCollection)
				require.NoError(t, err)
				ID = &wantMedia.ID
			}

			gotMedia, err := FindMediaByID(ID, testCollection)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantMedia.FileID, gotMedia.FileID)
			assert.Equal(t, tt.wantMedia.Histogram, gotMedia.Histogram)
			assert.Equal(t, tt.wantMedia.HistogramAverage, gotMedia.HistogramAverage)
			assert.Equal(t, tt.wantMedia.HistogramSum, gotMedia.HistogramSum)
			assert.Equal(t, tt.wantMedia.PHash, gotMedia.PHash)

		})
	}
}
