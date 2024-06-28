package documentstore

import (
	"context"
	"log"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/shitpostingio/admin-bot/analysisadapter"
	"github.com/shitpostingio/admin-bot/api/botapi"
	"github.com/shitpostingio/admin-bot/api/cache"
	"github.com/shitpostingio/admin-bot/config"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
	"github.com/shitpostingio/admin-bot/repository"
)

var (
	testCollection *mongo.Collection
)

func TestMain(m *testing.M) {

	cfg, err := config.Load("../../config.toml", false)
	if err != nil {
		log.Fatal("Unable to load configuration:", err)
	}

	err = loglog.Setup(cfg.Loglog.ApplicationName)
	if err != nil {
		log.Fatal("Unable to set up loglog:", err)
	}

	bot, err := botapi.Authorize(cfg.Telegram.BotToken, false)
	if err != nil {
		log.Fatal("Unable to connect to the bot apis", err)
	}

	//_, err = tdlib.Authorize(cfg.Telegram.BotToken, &cfg.Tdlib)
	//if err != nil {
	//	utility.LogFatal("Unable to log into the bot via Tdlib:", err)
	//}

	Connect(&cfg.DocumentStore)
	testCollection = database.Collection("a_test_collection")
	err = testCollection.Drop(context.Background())
	if err != nil {
		log.Fatal("Unable to drop test collection")
	}

	repository.SetBot(bot)
	repository.SetConfig(&cfg)
	cache.CreateActionsCache()
	analysisadapter.Start(cfg.Telegram.BotToken, &cfg.FPServer)
	limiter.StartRateLimiter(cfg.RateLimiter.MaxActionsPerSecond)

	//TODO: droppare contenuti test collection a fine test
	os.Exit(m.Run())

}
