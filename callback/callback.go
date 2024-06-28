package callback

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/shitpostingio/admin-bot/adminbot"
	"github.com/shitpostingio/admin-bot/localization"
	"github.com/shitpostingio/admin-bot/repository"
)

var (
	authorizationCache *cache.Cache
)

func init() {
	authorizationCache = cache.New(5*time.Second, 10*time.Minute)
}

type authorizationRequest struct {
	action  string
	mutex   sync.Mutex
	counter uint8
	handled bool
}

//HandleCallback handles all callback queries
func HandleCallback(callbackQuery *tgbotapi.CallbackQuery) {

	callbackFields := strings.Fields(callbackQuery.Data)

	var response string
	if callbackFields[0] == "verify" {
		handleHumanVerification(callbackFields, callbackQuery)
	} else if callbackFields[0] == "2sa" {
		response = handleAdminOnlyRequests(callbackFields, callbackQuery)
	} else {
		response = handleRequest(callbackFields, callbackQuery)
	}

	if response == "" {
		return
	}

	err := adminbot.SendCallbackWithAlert(callbackQuery.ID, response)
	if err != nil {
		log.Error("Unable to send callback response:", err)
	}

}

func handleHumanVerification(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) string {

	userID, _ := strconv.ParseInt(callbackFields[1], 10, 64)
	userID = userID - repository.Bot.Self.ID
	if callbackQuery.From.ID != userID {
		return ""
	}

	chatMember, err := adminbot.GetChatMember(userID, repository.GetTelegramConfiguration().GroupID)
	if err != nil {
		log.Error("Unable to find user with ID", userID, "in the group for the human verification:", err)
		return localization.GetString("callback_verification_error_occurred")
	}

	err = adminbot.UnrestrictUser(chatMember.User, repository.GetTelegramConfiguration().GroupID, callbackQuery.From)
	if err != nil {

		log.Error("Unable to unrestrict user with ID", userID, "after successfully passing the human verification:", err)
		return localization.GetString("callback_verification_error_occurred")

	}

	adminbot.DeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	return ""

}

func handleAdminOnlyRequests(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) string {

	if !repository.Admins[callbackQuery.From.ID] {
		return localization.GetString("user_unauthorized")
	}

	return handleRequest(callbackFields, callbackQuery)

}

func handleRequest(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) string {

	messageIDStr := strconv.Itoa(callbackQuery.Message.MessageID)
	item, found := authorizationCache.Get(messageIDStr)
	if !found {
		return authorizeFirstStep(messageIDStr, callbackFields)
	}

	var reply string
	request := item.(*authorizationRequest)
	request.mutex.Lock()
	defer request.mutex.Unlock()

	if request.handled {
		return localization.GetString("callback_request_action_already_authorized")
	}

	if request.action == callbackFields[1] {

		request.counter++
		if request.counter >= 2 {
			go handleSecondStep(callbackFields, callbackQuery)
			request.handled = true
		}

		reply = localization.GetString("callback_request_performed_shortly")

	} else {

		if request.counter < 2 {
			reply = localization.GetString("callback_twostep_different_action_alredy_requested")
		} else {
			reply = localization.GetString("callback_twostep_different_action_already_approved")
		}

	}

	return reply

}

//authorizeFirstStep authorizes the first step of a 2fa action
func authorizeFirstStep(messageIDStr string, callbackFields []string) string {

	request := authorizationRequest{
		action:  callbackFields[1],
		counter: 1,
		mutex:   sync.Mutex{},
		handled: false,
	}

	err := authorizationCache.Add(messageIDStr, &request, cache.DefaultExpiration)
	if err != nil {
		return localization.GetString("callback_first_step_error")
	}

	return localization.GetString("callback_first_step_perform_second")

}

func handleSecondStep(callbackFields []string, callbackQuery *tgbotapi.CallbackQuery) {

	/* GET 2ND STEP CALLBACK FIELDS */
	callbackFields = callbackFields[1:]

	switch callbackFields[0] {
	case "ok":
		markReportAsHandled(callbackQuery)
	case "whitelist":
		whitelistMedia(callbackFields, callbackQuery)
	case "tgunban":
		tgunbanUser(callbackFields, callbackQuery)
	case "unban":
		unbanUser(callbackFields, callbackQuery)
	case "unrestrict":
		unrestrictUser(callbackFields, callbackQuery)
	case "mod":
		modUser(callbackFields, callbackQuery)
	case "bbh", "bth", "bh":
		handleBlacklistHostname(callbackFields, callbackQuery)
	case "bg":
		banUserForGibberish(callbackFields, callbackQuery)
	case "blm":
		blacklistPrivateMedia(callbackQuery)
	case "wlm":
		whitelistPrivateMedia(callbackQuery)
	case "parm":
		pardonPrivateMedia(callbackQuery)
	case "blsp":
		blacklistPrivateStickerPack(callbackQuery)
	case "parsp":
		pardonPrivateStickerPack(callbackQuery)
	case "blms":
		blacklistPrivateSticker(callbackQuery)
	case "wlms":
		whitelistPrivateSticker(callbackQuery)
	case "parms":
		pardonPrivateSticker(callbackQuery)
	case "bs":
		blacklistSource(callbackFields, callbackQuery)
	case "ps":
		pardonSource(callbackFields, callbackQuery)
	case "ws":
		whitelistSource(callbackFields, callbackQuery)
	}

}
