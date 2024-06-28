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

//BlacklistStickerPack blacklists a sticker pack by its set_name. Also logs the user who added it to the blacklist
func BlacklistStickerPack(setName string, blacklisterTelegramID int64, collection *mongo.Collection) (generatedID string, err error) {

	if setName == "" {
		err = xerrors.New("BlacklistStickerPack: attempt to blacklist sticker pack with an empty setname")
		return
	}

	//
	stickerPack := entities.StickerPack{
		SetName: setName,
		AddedBy: blacklisterTelegramID,
		AddedAt: time.Now(),
	}

	//
	// TODO: Controllare se context.TODO è appropriato o meno
	ctx, cancelCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, stickerPack)
	if err != nil {
		err = xerrors.Errorf("BlacklistStickerPack: unable to add sticker pack %s into the document store: %s", setName, err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return generatedID, err
}

//PardonStickerPack removes a sticker pack from the blacklist. Also logs the user who did it
func PardonStickerPack(setName string, collection *mongo.Collection) error {

	if setName == "" {
		return xerrors.New("PardonStickerPack: attempt to blacklist sticker pack with an empty setname")
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"setname": setName}
	_, err := collection.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return xerrors.Errorf("PardonStickerPack: error while removing sticker pack %s: %s", setName, err)
	}

	return nil

}

// StickerPackIsBlacklisted returns true if a sticker pack is not blacklisted
func StickerPackIsBlacklisted(setName string, collection *mongo.Collection) bool {

	if setName == "" {
		return false
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"setname": setName}
	result := collection.FindOne(ctx, filter, options.FindOne())
	return result.Err() == nil

}
