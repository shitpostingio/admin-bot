package documentstore

import (
	"context"
	"fmt"
	"strings"
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
 *														SELECT														   *
 *																													   *
 ***********************************************************************************************************************
 */

// GetSourceByUsername gets a source from its username
func GetSourceByUsername(username string, collection *mongo.Collection) (source entities.Source, err error) {

	if username == "" {
		err = xerrors.New("GetSourceByUsername: empty username")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"username": username}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&source)
	return

}

// GetSourceByTelegramID gets a source from its telegram id
func GetSourceByTelegramID(telegramID int64, collection *mongo.Collection) (source entities.Source, err error) {

	if telegramID == 0 {
		err = xerrors.New("GetSourceByTelegramID: telegramID 0")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"telegramid": telegramID}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&source)
	return

}

// GetSource gets a source from the document store
func GetSource(telegramID int64, username string, collection *mongo.Collection) (source entities.Source, err error) {

	if telegramID != 0 {
		source, err = GetSourceByTelegramID(telegramID, collection)
		if err == nil {
			return
		}
	}

	if username != "" {
		source, err = GetSourceByUsername(username, collection)
		if err == nil {
			return
		}
	}

	return source, fmt.Errorf("GetSource: no source found for telegram ID %d or username %s", telegramID, username)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														INSERT														   *
 *																													   *
 ***********************************************************************************************************************
 */

// BlacklistSource blacklists a source
func BlacklistSource(sourceTelegramID int64, sourceUsername string, adderTelegramID int64, collection *mongo.Collection) (generatedID string, err error) {

	if sourceTelegramID == 0 && sourceUsername == "" {
		err = xerrors.New("BlacklistSource: telegram id 0 and empty sourceUsername")
		return
	}

	source, err := GetSource(sourceTelegramID, sourceUsername, collection)
	if err == nil {

		generatedID = source.ID.Hex()

		if !source.IsWhitelisted {
			err = updateSourceByID(&source.ID, sourceTelegramID, sourceUsername, collection)
			return
		}

		err = markSourceAsBlacklisted(&source.ID, collection)
		if err != nil {
			err = xerrors.Errorf("BlacklistSource: %s", err)
		}

		return

	}

	source = entities.Source{
		AddedBy:      adderTelegramID,
		LastModified: time.Now(),
	}

	if sourceTelegramID != 0 {
		source.TelegramID = sourceTelegramID
	}

	if sourceUsername != "" {
		source.Username = strings.ToLower(sourceUsername)
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, source)
	if err != nil {
		err = xerrors.Errorf("BlacklistSource: unable to add source into the document store: %s", err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return generatedID, err

}

// WhitelistSource whitelist a source
func WhitelistSource(sourceTelegramID int64, sourceUsername string, adderTelegramID int64, collection *mongo.Collection) (generatedID string, err error) {

	if sourceTelegramID == 0 && sourceUsername == "" {
		err = xerrors.New("BlacklistSource: telegram id 0 and empty sourceUsername")
		return
	}

	source, err := GetSource(sourceTelegramID, sourceUsername, collection)
	if err == nil {

		generatedID = source.ID.Hex()

		if source.IsWhitelisted {
			err = updateSourceByID(&source.ID, sourceTelegramID, sourceUsername, collection)
			return
		}

		err = markSourceAsWhitelisted(&source.ID, collection)
		if err != nil {
			err = xerrors.Errorf("WhitelistSource: %s", err)
		}

		return

	}

	source = entities.Source{
		AddedBy:       adderTelegramID,
		LastModified:  time.Now(),
		IsWhitelisted: true,
	}

	if sourceTelegramID != 0 {
		source.TelegramID = sourceTelegramID
	}

	if sourceUsername != "" {
		source.Username = strings.ToLower(sourceUsername)
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, source)
	if err != nil {
		err = xerrors.Errorf("BlacklistSource: unable to add source into the document store: %s", err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return generatedID, err

}

// RemoveSource removes a source from the document store
func RemoveSource(telegramID int64, username string, collection *mongo.Collection) error {

	if telegramID == 0 && username == "" {
		return xerrors.New("RemoveSource: telegram id 0 and empty username")
	}

	source, err := GetSource(telegramID, username, collection)
	if err != nil {
		return xerrors.Errorf("RemoveSource: source (%d, %s) is not in the database", telegramID, username)
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"_id": source.ID}
	_, err = collection.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return xerrors.Errorf("RemoveSource: error while deleting source with id %s: %s", source.ID.Hex(), err)
	}

	return nil

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														UPDATES														   *
 *																													   *
 ***********************************************************************************************************************
 */

// UpdateSource updates the source
func UpdateSource(telegramID int64, username string, collection *mongo.Collection) error {

	source, err := GetSource(telegramID, username, collection)
	if err != nil {
		return xerrors.Errorf("UpdateSource: %s", err)
	}

	return updateSourceByID(&source.ID, telegramID, username, collection)

}

func updateSourceByID(ID *primitive.ObjectID, telegramID int64, username string, collection *mongo.Collection) error {

	if ID == nil {
		return xerrors.New("updateSourceByID: ID nil")
	}

	//
	filter := bson.M{"_id": ID}
	updateMap := bson.M{"lastModified": time.Now()}

	if telegramID != 0 {
		updateMap["telegramid"] = telegramID
	}

	if username != "" {
		updateMap["username"] = strings.ToLower(username)
	}

	//
	update := bson.D{
		{
			"$set", updateMap,
		},
	}

	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	_, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("updateSourceByID: unable to update source with id %s: %s", ID, err)
	}

	return nil

}

func markSourceAsBlacklisted(ID *primitive.ObjectID, collection *mongo.Collection) error {

	if ID == nil {
		return xerrors.New("markSourceAsBlacklisted: ID nil")
	}

	//
	filter := bson.M{"_id": ID}
	update := bson.D{
		{
			"$set", bson.M{
				"iswhitelisted": false,
			},
		},
	}
	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	_, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("markSourceAsBlacklisted: unable to mark source with id %s as blacklisted: %s", ID, err)
	}

	return nil

}

func markSourceAsWhitelisted(ID *primitive.ObjectID, collection *mongo.Collection) error {

	if ID == nil {
		return xerrors.New("markSourceAsWhitelisted: ID nil")
	}

	//
	filter := bson.M{"_id": ID}
	update := bson.D{
		{
			"$set", bson.M{
				"iswhitelisted": true,
			},
		},
	}
	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	_, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("markSourceAsWhitelisted: unable to mark source with id %s as whitelisted: %s", ID, err)
	}

	return nil

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														CHECKS														   *
 *																													   *
 ***********************************************************************************************************************
 */

// SourceIsWhitelisted returns true if the source is whitelisted
func SourceIsWhitelisted(telegramID int64, username string, collection *mongo.Collection) bool {
	source, err := GetSource(telegramID, username, collection)
	return err == nil && source.IsWhitelisted
}

// SourceIsBlacklisted returns true if the source is blacklisted
func SourceIsBlacklisted(telegramID int64, username string, collection *mongo.Collection) bool {
	source, err := GetSource(telegramID, username, collection)
	return err == nil && !source.IsWhitelisted
}
