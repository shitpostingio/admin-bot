package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/analysisadapter"
	"github.com/shitpostingio/admin-bot/api/botapi"
	"github.com/shitpostingio/admin-bot/api/cache"
	"github.com/shitpostingio/admin-bot/api/tdlib"
	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/database/database"
	"github.com/shitpostingio/admin-bot/database/documentstore"
	"github.com/shitpostingio/admin-bot/defense/antispam"
	"github.com/shitpostingio/admin-bot/localization"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/updates"
	"github.com/shitpostingio/admin-bot/utility"
	modCache "github.com/shitpostingio/admin-bot/utility/cache"

	"github.com/bykovme/gotrans"

	"github.com/shitpostingio/admin-bot/automod"
	"github.com/shitpostingio/admin-bot/callback"
	configuration "github.com/shitpostingio/admin-bot/config"
	"github.com/shitpostingio/admin-bot/private"
)

var (
	//configFilePath is the path to the config file
	configFilePath string

	//Version represents the current admin-bot version, a compile-time value
	Version string

	//Build is the git tag for the current version
	Build string

	//debug tells whether the bot is in debug mode or not
	debug bool

	//polling tells whether to retrieve updates via polling or webhook
	polling bool

	//testing tells whether we're in testing mode or not
	testing bool
)

func main() {

	setCLIParams()

	/*************************************************
	 *				CONNECT TO SERVICES				 *
	 *************************************************/
	cfg, err := configuration.Load(configFilePath, testing)
	if err != nil {
		log.Fatal("Unable to load configuration:", err)
	}

	err = os.RemoveAll(cfg.Tdlib.FilesDirectory)
	if err != nil {
		log.Error("Unable to delete tdlib files directory", err)
	}

	err = os.RemoveAll(cfg.Tdlib.DatabaseDirectory)
	if err != nil {
		log.Error("Unable to delete tdlib database directory", err)
	}

	bot, err := botapi.Authorize(cfg.Telegram.BotToken, debug)
	if err != nil {
		log.Fatal("Unable to connect to the bot apis", err)
	}

	_, err = tdlib.Authorize(cfg.Telegram.BotToken, &cfg.Tdlib)
	if err != nil {
		utility.LogFatal("Unable to log into the bot via Tdlib:", err)
	}

	documentstore.Connect(&cfg.DocumentStore)

	/*************************************************
	 *		  	LOAD LOCALIZATION FILES				 *
	 *************************************************/
	err = gotrans.InitLocales(cfg.AdminBot.LocalizationPath)
	if err != nil {
		utility.LogFatal("Unable to load language files:", err)
	}

	localization.SetLanguage(cfg.AdminBot.Language)

	/*************************************************
	 *				  CACHE MODS DATA				 *
	 *************************************************/
	_ = database.UpdateModeratorsDetails(cfg.Telegram.GroupID, bot)
	adminMap, err := modCache.CreateAdminsCache(documentstore.ModeratorsCollection)
	if err != nil {
		utility.LogFatal("Unable to cache admins:", err)
	}

	modMap, err := modCache.CreateModsCache(adminMap, bot, &cfg.Telegram)
	if err != nil {
		utility.LogFatal("Unable to cache mods:", err)
	}

	/*************************************************
	 *				POPULATE REPOSITORY				 *
	 *************************************************/
	repository.SetBot(bot)
	repository.SetConfig(&cfg)
	repository.SetAdmins(adminMap)
	repository.SetMods(modMap)
	repository.SetTestingStatus(testing)
	repository.SetDebugStatus(debug)

	/*************************************************
	 *				 START MONITORING				 *
	 *************************************************/
	cache.CreateActionsCache()
	analysisadapter.Start(cfg.Telegram.BotToken, &cfg.FPServer)
	limiter.StartRateLimiter(cfg.RateLimiter.MaxActionsPerSecond)
	automod.StartDefensiveRoutines()
	configuration.WatchConfig(&cfg)

	/*************************************************
	 *			  START HANDLING UPDATES			 *
	 *************************************************/
	updatesChannel := updates.GetUpdatesChannel(polling, bot, &cfg)
	if updatesChannel == nil {
		utility.LogFatal("Update channel nil")
	}

	authorizationText := fmt.Sprintf(localization.GetString("shitposting_bot_active"), Version, Build, bot.Self.ID, bot.Self.UserName)
	_ = adminbot.SendSilentTextMessage(cfg.Telegram.ReportChannelID, authorizationText, false)
	log.Info(authorizationText)
	handleUpdates(updatesChannel, &cfg)

}

//setCLIParams parses the command line parameters and sets defaults in case they're missing
func setCLIParams() {
	flag.BoolVar(&polling, "polling", false, "use polling instead of webhooks")
	flag.BoolVar(&testing, "testing", false, "use testing mode")
	flag.BoolVar(&debug, "debug", false, "activate all the debug features")
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.Parse()
}

//handleUpdates iterates on the updates and passes them onto the handlers
func handleUpdates(updates tgbotapi.UpdatesChannel, cfg *structs.Config) {
	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			go callback.HandleCallback(update.CallbackQuery)
		case update.EditedMessage != nil:
			go handleMessage(update.EditedMessage, cfg)
		case update.Message != nil:
			go handleMessage(update.Message, cfg)
		case update.ChannelPost != nil:
			go handleChannelPost(update.ChannelPost, cfg)
		}
	}
}

//handleMessage handles `Message`s and EditedMessages
func handleMessage(msg *tgbotapi.Message, cfg *structs.Config) {

	/* PRIVATE MESSAGES ARE NOT SPAM */
	if msg.Chat.IsPrivate() {
		private.HandlePrivateChat(msg)
		return
	}

	/* SAVE ALL PUBLIC MESSAGES */
	go documentstore.StoreMessage(msg)

	/* LEAVE THE GROUP IF THE BOT SHOULDN'T BE IN IT */
	if msg.Chat.ID != cfg.Telegram.GroupID && msg.LeftChatMember == nil {
		adminbot.LeaveUnauthorizedGroup(msg)
		return
	}

	/* SEND EVERY MESSAGE TO THE ANTI SPAM */
	antispam.HandleMessage(msg)

	/* SEND THE MESSAGE TO THE APPROPRIATE HANDLER */
	switch {
	case msg.Text != "":
		automod.HandleText(msg)
	case msg.Photo != nil:
		automod.HandleMedia(msg, true)
	case msg.Video != nil:
		automod.HandleMedia(msg, true)
	case msg.Animation != nil:
		automod.HandleMedia(msg, true)
	case msg.Sticker != nil:
		automod.HandleSticker(msg)
	case msg.Document != nil:
		automod.HandleMedia(msg, false)
	case msg.NewChatMembers != nil:
		automod.HandleNewChatMember(msg)
	case msg.Voice != nil:
		automod.HandleMedia(msg, false)
	case msg.Audio != nil:
		automod.HandleMedia(msg, false)
	case msg.VideoNote != nil:
		automod.HandleMedia(msg, false)
	case msg.Game != nil:
		automod.HandleGame(msg)
	}
}

//handleChannelPost handles `ChannelPost`s
func handleChannelPost(post *tgbotapi.Message, cfg *structs.Config) {
	if post.Chat.ID != cfg.Telegram.ReportChannelID && post.Chat.ID != cfg.Telegram.BackupChannelID {
		adminbot.LeaveUnauthorizedChannel(post)
	}
}
