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
	"github.com/shitpostingio/admin-bot/utility"
)

// GetHostName gets a hostname from the document store given a url
func GetHostName(url string, collection *mongo.Collection) (host entities.HostName, err error) {

	//
	if url == "" {
		err = xerrors.New("GetHostName: tried to find empty host")
		return
	}

	//
	hostname, err := utility.GetHostNameFromURL(url)
	if err != nil || hostname == "" {
		err = xerrors.Errorf("GetHostName: unable to find host for input url %s (error: %s)", url, err)
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"host": hostname}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&host)
	return

}

// BlacklistHostName blacklists a hostname
func BlacklistHostName(url string, isBanworthy, isTelegram bool, adderTelegramID int64, collection *mongo.Collection) (generatedID, hostname string, err error) {

	//
	if url == "" {
		err = xerrors.New("BlacklistHostName: tried to add empty host")
		return
	}

	//
	hostname, err = utility.GetHostNameFromURL(url)
	if err != nil || hostname == "" {
		err = xerrors.Errorf("BlacklistHostName: unable to find host for input url %s (error: %s)", url, err)
		return
	}

	//
	host := entities.HostName{
		Host:         hostname,
		IsBanworthy:  isBanworthy,
		IsTelegram:   isTelegram,
		LastEditedBy: adderTelegramID,
		LastModified: time.Now(),
	}

	//
	// TODO: Controllare se context.TODO è appropriato o meno
	ctx, cancelCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, host)
	if err != nil {
		err = xerrors.Errorf("BlacklistHostName: unable to add host into the document store: %s", err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return

}

// UpdateHostName updates a hostname
func UpdateHostName(url string, isBanworthy, isTelegram bool, updaterTelegramID int64, collection *mongo.Collection) error {

	//
	if url == "" {
		return xerrors.New("UpdateHostName: tried to update empty host")
	}

	//
	hostname, err := utility.GetHostNameFromURL(url)
	if err != nil || hostname == "" {
		return xerrors.Errorf("UpdateHostName: unable to find host for input url %s (error: %s)", url, err)
	}

	//
	filter := bson.M{"host": hostname}
	update := bson.D{
		{
			"$set", bson.M{
				"isbanworthy":  isBanworthy,
				"istelegram":   isTelegram,
				"lasteditedby": updaterTelegramID,
				"lastmodified": time.Now(),
			},
		},
	}
	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	result, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("UpdateHostName: unable to update host %s: %s", hostname, err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("UpdateHostName: no match for host %s", hostname)
	}

	return nil

}

// PardonHostName removes a hostname from the document store
func PardonHostName(url string, collection *mongo.Collection) error {

	//
	if url == "" {
		return xerrors.New("PardonHostName: tried to remove empty host")
	}

	//
	hostname, err := utility.GetHostNameFromURL(url)
	if err != nil || hostname == "" {
		return xerrors.Errorf("PardonHostName: unable to find host for input url %s (error: %s)", url, err)
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"host": hostname}
	result, err := collection.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return xerrors.Errorf("PardonHostName: error while deleting host %s: %s", hostname, err)
	}

	if result.DeletedCount == 0 {
		return xerrors.Errorf("PardonHostName: no matches found for host %s", hostname)
	}

	return nil

}
