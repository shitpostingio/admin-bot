package antiuserbot

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/reports"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/agnivade/levenshtein"

	"github.com/shitpostingio/image-fingerprinting/comparer"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/api"
	"github.com/shitpostingio/admin-bot/config/structs"
	"github.com/shitpostingio/admin-bot/repository"
	"github.com/shitpostingio/admin-bot/telegram"

	"github.com/shitpostingio/admin-bot/analysisadapter"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														STRUCTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

// AttackState represents infos about the handling of userbot attacks
type AttackState struct {
	handled bool
	mutex   sync.Mutex
}

// ChatMemberInfos contains the info needed by Anti Userbot
type ChatMemberInfos struct {
	ID                  int64
	Lock                sync.Mutex
	Restricted          bool
	Name                string
	Handle              string
	HasProfilePicture   bool
	ProfilePicturePHash string
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												CONSTS AND VARS														   *
 *																													   *
 ***********************************************************************************************************************
 */

const (
	/* CACHE */
	suspicionKey   = "suspicion"
	userExpiration = 1 * time.Minute
	userCleanup    = 5 * time.Minute

	/* PRINTS */
	unableToAddToCache        = "Unable to add user %s (id %d) to the join cache: %s"
	chatMemberNotRestricted   = "Unable to restrict user %s (@%s, id %d) for being a possible userbot"
	telegramUserNotRestricted = "Unable to restrict user %s (id %d) for being a possible userbot"
)

var (
	cfg            *structs.AntiUserbotConfiguration
	joinChannel    chan *tgbotapi.User
	joinCache      *cache.Cache
	suspicionCache *cache.Cache
	state          AttackState
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														START														   *
 *																													   *
 ***********************************************************************************************************************
 */

// Start starts join monitoring
func Start() {

	cfg = repository.GetAntiUserbotConfiguration()

	joinChannel = make(chan *tgbotapi.User)
	joinCache = cache.New(userExpiration, userCleanup)
	suspicionCache = cache.New(cache.NoExpiration, cache.NoExpiration)

	err := suspicionCache.Add(suspicionKey, 0, cache.NoExpiration)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to add counter to antiuserbot suspicion cache: %s", err.Error()))
		os.Exit(-10)
	}

	go handleNewUsers()
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												JOIN ROUTINES														   *
 *																													   *
 ***********************************************************************************************************************
 */

// handleNewUsers restricts new users if we're being attacked by userbots or starts a routine for each user
func handleNewUsers() {
	for {

		newUser := <-joinChannel

		if IsAttack() {
			go muteTelegramUser(newUser)
			continue
		}

		go handleSingleUser(newUser)

	}
}

// handleSingleUser performs check to determine the similarity of this user compared to
// the ones that joined recently, in order to prevent userbots
func handleSingleUser(user *tgbotapi.User) {

	chatMember := ChatMemberInfos{
		ID:     user.ID,
		Name:   user.FirstName + user.LastName,
		Handle: user.UserName,
	}

	err := joinCache.Add(strconv.FormatInt(user.ID, 10), &chatMember, userExpiration)
	if err != nil {
		log.Warn(fmt.Sprintf(unableToAddToCache, telegram.GetNameOrUsername(user), user.ID, err.Error()))
	}

	// Increment suspicion after adding the user to the cache,
	// so we can restrict them in case of the suspicion going
	// over the threshold.
	suspicionCount := 1
	increaseSuspicion(1)

	suspicionCount += checkTextSimilarities(&chatMember)
	if !chatMember.Restricted && IsAttack() {
		muteChatMember(&chatMember)
	}

	fingerprintUserProfilePicture(&chatMember)
	suspicionCount += checkPictureSimilarities(&chatMember)
	if !chatMember.Restricted && IsAttack() {
		muteChatMember(&chatMember)
	}

	// Wait `joinRoutineLifespan` before decreasing the suspicion.
	time.Sleep(time.Duration(cfg.RoutineLifespan))
	decreaseSuspicion(suspicionCount)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													SIMILARITIES													   *
 *																													   *
 ***********************************************************************************************************************
 */

// checkTextSimilarities checks for similarities in the names and in the handles of the recent joins
func checkTextSimilarities(chatMember *ChatMemberInfos) (similarities int) {

	// 1 suspicion point for no handle.
	if chatMember.Handle == "" {
		similarities += increaseSuspicion(1)
	}

	for _, item := range joinCache.Items() {

		currentUser := item.Object.(*ChatMemberInfos)
		if currentUser.ID == chatMember.ID {
			continue
		}

		similarities += increaseSuspicion(checkNameSimilarities(chatMember, currentUser))
		similarities += increaseSuspicion(checkHandleSimilarities(chatMember, currentUser))
	}

	return similarities
}

// checkNameSimilarities computes the Levenshtein distance between the names of two users
func checkNameSimilarities(chatMember *ChatMemberInfos, toCompare *ChatMemberInfos) (similarities int) {

	nameDistance := levenshtein.ComputeDistance(chatMember.Name, toCompare.Name)

	if nameDistance < len(chatMember.Name)/2 {
		similarities += increaseSuspicion(1)
	}

	return
}

// checkNameSimilarities computes the Levenshtein distance between the handles of two users
func checkHandleSimilarities(chatMember *ChatMemberInfos, toCompare *ChatMemberInfos) (similarities int) {

	if chatMember.Handle != "" {

		handleDistance := levenshtein.ComputeDistance(chatMember.Handle, toCompare.Handle)

		if handleDistance < len(chatMember.Handle)/2 {
			similarities += increaseSuspicion(1)
		}
	}

	return
}

// fingerprintUserProfilePicture gets the fingerprint of the user's current profile picture
func fingerprintUserProfilePicture(chatMember *ChatMemberInfos) {

	userPhotos, err := adminbot.GetUserProfilePhotos(chatMember.ID, 1)
	if err == nil && userPhotos.TotalCount != 0 {

		profilePicture := userPhotos.Photos[0][len(userPhotos.Photos[0])-1]
		fingerprint, err := analysisadapter.GetFingerprint(profilePicture.FileUniqueID, profilePicture.FileID)

		if err == nil {
			chatMember.HasProfilePicture = true
			chatMember.ProfilePicturePHash = fingerprint.PHash
		}
	}
}

// checkNameSimilarities checks the similarity between
// the user's current profile picture and the ones of the recent joins
func checkPictureSimilarities(chatMember *ChatMemberInfos) (similarities int) {

	if !chatMember.HasProfilePicture {
		similarities += increaseSuspicion(1)
		return
	}

	for _, item := range joinCache.Items() {

		currentUser := item.Object.(*ChatMemberInfos)
		if currentUser.ID == chatMember.ID {
			continue
		}

		if fpcompare.PhotosAreSimilarEnough(chatMember.ProfilePicturePHash, currentUser.ProfilePicturePHash) {

			// 2 suspicion points if the photos are similar.
			similarities += increaseSuspicion(2)

		}
	}

	return
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													RESTRICTIONS													   *
 *																													   *
 ***********************************************************************************************************************
 */

func muteAllNewMembers() {
	for _, item := range joinCache.Items() {
		go muteChatMember(item.Object.(*ChatMemberInfos))
	}
}

func muteChatMember(user *ChatMemberInfos) {

	user.Lock.Lock()
	defer user.Lock.Unlock()
	if user.Restricted {
		return
	}

	_, err := api.RestrictMessages(user.ID, repository.GetTelegramConfiguration().GroupID, 0)
	if err == nil {

		user.Restricted = true
		restrictionReport := reports.PossibleUserbotRestriction(user.ID, user.Name)
		_ = reports.Report(restrictionReport, reports.NON_URGENT)
		log.Warn(restrictionReport)

	} else {

		log.Error(fmt.Sprintf(chatMemberNotRestricted, user.Name, user.Handle, user.ID))
	}
}

func muteTelegramUser(user *tgbotapi.User) {

	err := adminbot.RestrictMessages(user, repository.GetTelegramConfiguration().GroupID, 0, &repository.Bot.Self)

	if err == nil {

		restrictionReport := reports.PossibleUserbotRestriction(user.ID, telegram.GetName(user))
		_ = reports.Report(restrictionReport, reports.NON_URGENT)
		log.Warn(restrictionReport)

	} else {

		log.Error(fmt.Sprintf(telegramUserNotRestricted, telegram.GetNameOrUsername(user), user.ID))
	}
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														SUSPICION													   *
 *																													   *
 ***********************************************************************************************************************
 */

func increaseSuspicion(amount int) (amountIncreased int) {

	currentSuspicion, err := suspicionCache.IncrementInt(suspicionKey, amount)
	if err != nil {
		return
	}

	amountIncreased = amount
	if currentSuspicion > cfg.JoinThreshold {

		state.mutex.Lock()
		if !state.handled {
			state.handled = true
			go reportPossibleUserbotAttack()
			muteAllNewMembers()
		}
		state.mutex.Unlock()
	}

	return
}

func decreaseSuspicion(amount int) {

	// In case of userbot attack we should wait
	// before lowering our defenses.
	if IsAttack() {
		if repository.GetTestingStatus() {
			time.Sleep(15 * time.Second)
		} else {
			time.Sleep(5 * time.Minute)
		}
	}

	currentSuspicion, err := suspicionCache.DecrementInt(suspicionKey, amount)
	if err != nil {
		log.Warn("Unable to decrease suspicion")
	}

	if currentSuspicion < cfg.JoinThreshold {
		state.handled = false
	}
}

func getUserbotAttackSuspicion() int {
	valueItf, _ := suspicionCache.Get(suspicionKey)
	return valueItf.(int)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														REPORTING													   *
 *																													   *
 ***********************************************************************************************************************
 */

func reportPossibleUserbotAttack() {
	_ = reports.ReportInPlaintext(reports.PossibleUserbotAttack(), reports.URGENT)
	log.Warn(reports.PossibleUserbotAttack())
}

/*
 ***********************************************************************************************************************
 *																													   *
 *														ACCESSORS													   *
 *																													   *
 ***********************************************************************************************************************
 */

// IsAttack returns true if the chat is under a possible userbot attack
func IsAttack() bool {
	return getUserbotAttackSuspicion() > cfg.JoinThreshold
}

// HandleUser sends an user to the handling routine.
func HandleUser(user *tgbotapi.User) {
	joinChannel <- user
}
