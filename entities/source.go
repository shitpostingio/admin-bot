package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Source represents a source in the document store
type Source struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TelegramID    int64              `bson:",omitempty"`
	Username      string             `bson:",omitempty"`
	IsWhitelisted bool               `bson:",omitempty"`
	AddedBy       int64              `bson:",omitempty"`
	LastModified  time.Time
}
