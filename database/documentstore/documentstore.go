package documentstore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/utility"
)

const (

	/* TIMEOUT */
	opDeadline = 10 * time.Second

	/* COLLECTION NAMES */
	bansCollectionName        = "bans"
	hostNamesCollectionName   = "hostnames"
	mediaCollectionName       = "media"
	moderatorsCollectionName  = "moderators"
	sourcesCollectionName     = "sources"
	messagesCollectionName    = "messages"
	stickerPackCollectionName = "stickerpacks"
)

var (

	//
	dsCtx    context.Context
	database *mongo.Database

	// BansCollection represents the bans collection in the document store
	BansCollection *mongo.Collection

	// HostNameCollection represents the hostname collection in the document store
	HostNameCollection *mongo.Collection

	// MediaCollection represents the media collection in the document store
	MediaCollection *mongo.Collection

	// ModeratorsCollection represents the moderators collection in the document store
	ModeratorsCollection *mongo.Collection

	// SourcesCollection represents the sources collection in the document store
	SourcesCollection *mongo.Collection

	// MessagesCollection represents the messages collection in the document store
	MessagesCollection *mongo.Collection

	// StickerPackCollection represents the stickerpack collection in the document store
	StickerPackCollection *mongo.Collection
)

// Connect connects to the document store
func Connect(cfg *structs.DocumentStoreConfiguration) {

	client, err := mongo.Connect(context.Background(), cfg.MongoDBConnectionOptions())
	if err != nil {
		utility.LogFatal("Unable to connect to document store:", err)
	}

	pingCtx, cancelPingCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelPingCtx()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		utility.LogFatal("Unable to ping document store:", err)
	}

	//
	dsCtx = context.TODO()

	/* SAVE COLLECTIONS */
	database = client.Database(cfg.DatabaseName)
	BansCollection = database.Collection(bansCollectionName)
	HostNameCollection = database.Collection(hostNamesCollectionName)
	MediaCollection = database.Collection(mediaCollectionName)
	ModeratorsCollection = database.Collection(moderatorsCollectionName)
	SourcesCollection = database.Collection(sourcesCollectionName)
	MessagesCollection = database.Collection(messagesCollectionName)
	StickerPackCollection = database.Collection(stickerPackCollectionName)

}
