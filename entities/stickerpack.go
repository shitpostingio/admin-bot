package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StickerPack represents a sticker pack in the document store.
type StickerPack struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	SetName string
	AddedBy int64 `bson:",omitempty"`
	AddedAt time.Time
}
