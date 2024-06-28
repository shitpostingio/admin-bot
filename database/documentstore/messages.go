package documentstore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//StoreMessage serializes and saves a message in the database
func StoreMessage(message *tgbotapi.Message) {

	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	_, err := MessagesCollection.InsertOne(ctx, message)
	if err != nil {
		log.Error(fmt.Sprintf("Error while inserting the update into the document store: %s", err))
		return
	}

}

//GetUserJoinDate returns the join timestamp of a given user, returns time.Time{} if the join isn't in the database
func GetUserJoinDate(user *tgbotapi.User) time.Time {

	filter := bson.M{"newchatmembers.id": user.ID}
	findOneOptions := options.FindOne().SetSort(bson.D{{"_id", -1}})
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	var msg tgbotapi.Message
	err := MessagesCollection.FindOne(ctx, filter, findOneOptions).Decode(&msg)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(int64(msg.Date), 0)
}
