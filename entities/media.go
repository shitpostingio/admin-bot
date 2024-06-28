package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//nolint
const (
	PHOTO     = "AgA"
	VIDEO     = "BAA"
	ANIMATION = "CgA"
	STICKER   = "CAA"
	VOICE     = "AwA"
	DOCUMENT  = "BQA"
	AUDIO     = "CQA"
	VIDEONOTE = "DQA"
)

// Media represents a media in the document store
type Media struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	FileUniqueID     string
	FileID           string
	IsWhitelisted    bool      `bson:",omitempty"`
	Histogram        []float64 `bson:",omitempty"`
	HistogramAverage float64   `bson:",omitempty"`
	HistogramSum     float64   `bson:",omitempty"`
	PHash            string    `bson:",omitempty"`
	NSFWScore        float64   `bson:",omitempty"`
	NSFWDescription  string    `bson:",omitempty"`
	LastEditedBy     int64     `bson:",omitempty"`
	LastModified     time.Time
}

// MediaCanBeFingerprinted returns true if a media can be fingerprinted
func MediaCanBeFingerprinted(fileID string) bool {

	// Telegram prefixes are 3 characters long
	fileIDPrefix := fileID[:3]

	switch fileIDPrefix {
	case STICKER:
		return true
	case PHOTO:
		return true
	case VIDEO:
		return true
	case ANIMATION:
		return true
	}

	return false

}

// GetHistogramAverageAndSum gets the average and the sum of the input histogram values
func GetHistogramAverageAndSum(histogram []float64) (average, sum float64) {

	coefficient := 1.0
	for i := 0; i < 16; i++ {
		sum += histogram[i] * coefficient
		sum += histogram[31-i] * coefficient
		coefficient++
	}

	average = sum / float64(len(histogram))
	return

}
