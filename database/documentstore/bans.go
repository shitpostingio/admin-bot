package documentstore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"

	"github.com/shitpostingio/admin-bot/entities"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														FIND														   *
 *																													   *
 ***********************************************************************************************************************
 */

// GetBansForTelegramID gets bans for a telegramId from the document store
//TODO: Considerare la possibilità di aggiungere un limit
func GetBansForTelegramID(telegramID int64, collection *mongo.Collection) (bans []entities.Ban, err error) {

	//
	if telegramID == 0 {
		return bans, xerrors.New("GetBansForTelegramID: telegramID was 0")
	}

	//
	filter := bson.M{"user": telegramID}
	findOptions := options.Find().SetSort(bson.D{{"bandate", -1}})
	findCtx, cancelFindCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelFindCtx()

	cursor, err := collection.Find(findCtx, filter, findOptions)
	if err != nil {
		err = xerrors.Errorf("GetBansForTelegramID: unable to find bans for telegramID %d: %s", telegramID, err)
		return
	}

	decodeCtx, cancelDecodeCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelDecodeCtx()
	err = cursor.All(decodeCtx, &bans)
	if err != nil {
		err = xerrors.Errorf("GetBansForTelegramID: unable decode bans for telegramID %d: %s", telegramID, err)
	}

	return

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														INSERT														   *
 *																													   *
 ***********************************************************************************************************************
 */

// AddBan adds a ban to the document store
func AddBan(bannedUserID, moderatorUserID int64, reason string, collection *mongo.Collection) (generatedID string, err error) {

	//
	if bannedUserID == 0 {
		err = xerrors.New("AddBan: bannedUserID was 0")
		return
	}

	if moderatorUserID == 0 {
		err = xerrors.New("AddBan: moderatorUserID was 0")
		return
	}

	//
	ban := entities.Ban{
		User:     bannedUserID,
		BannedBy: moderatorUserID,
		Reason:   reason,
		BanDate:  time.Now(),
	}

	//
	// TODO: Controllare se context.TODO è appropriato o meno
	ctx, cancelCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, ban)
	if err != nil {
		err = xerrors.Errorf("AddBan: unable to add ban into the document store: %s", err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return generatedID, err

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														UPDATE														   *
 *																													   *
 ***********************************************************************************************************************
 */

// MarkUserAsUnbanned marks a user as unbanned
func MarkUserAsUnbanned(telegramID int64, collection *mongo.Collection) error {

	//
	if telegramID == 0 {
		return xerrors.New("MarkUserAsUnbanned: telegramID was 0")
	}

	//
	filter := bson.M{"user": telegramID, "unbandate": nil}
	update := bson.D{{"$set", bson.M{"unbandate": time.Now()}}}
	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	result, err := collection.UpdateMany(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("MarkUserAsUnbanned: unable to perform updates for telegramID %d: %s", telegramID, err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("MarkUserAsUnbanned: 0 matches for telegramID %d", telegramID)
	}

	return nil

}
