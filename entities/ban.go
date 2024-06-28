package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ban represents a ban in the document store
type Ban struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	User      int64
	BannedBy  int64
	Reason    string `bson:",omitempty"`
	BanDate   time.Time
	UnbanDate time.Time `bson:",omitempty"`
}
