package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HostName represents a hostname in the document store
type HostName struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Host         string
	IsBanworthy  bool  `bson:",omitempty"`
	IsTelegram   bool  `bson:",omitempty"`
	LastEditedBy int64 `bson:",omitempty"`
	LastModified time.Time
}
