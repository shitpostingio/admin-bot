package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Moderator represents a moderator in the document store
type Moderator struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	TelegramID         int64
	IsAdmin            bool
	CanChangeInfo      bool
	CanDeleteMessages  bool
	CanInviteUsers     bool
	CanRestrictMembers bool
	CanPinMessages     bool
	CanPromoteMembers  bool
	PromotedBy         int64 `bson:",omitempty"`
	PromotionDate      time.Time
}
