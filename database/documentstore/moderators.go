package documentstore

import (
	"context"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"

	"github.com/shitpostingio/admin-bot/api/tdlib"
	"github.com/shitpostingio/admin-bot/entities"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														SELECT														   *
 *																													   *
 ***********************************************************************************************************************
 */

// GetAllModerators retrieves all the moderators from the database.
func GetAllModerators(collection *mongo.Collection) (moderators []entities.Moderator, err error) {

	//
	findCtx, cancelFindCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelFindCtx()

	cursor, err := collection.Find(findCtx, bson.D{})
	if err != nil {
		err = xerrors.Errorf("GetAllModerators: unable to find moderators: %s", err)
		return
	}

	decodeCtx, cancelDecodeCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelDecodeCtx()
	err = cursor.All(decodeCtx, &moderators)
	if err != nil {
		err = xerrors.Errorf("GetAllModerators: unable decode moderators: %s", err)
	}

	return

}

// GetModeratorByTelegramID retrieves the moderator with the given telegram id.
func GetModeratorByTelegramID(telegramID int64, collection *mongo.Collection) (moderator entities.Moderator, err error) {

	if telegramID == 0 {
		err = xerrors.New("GetModeratorByTelegramID: telegramID 0")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"telegramid": telegramID}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&moderator)
	return

}

// GetModeratorByUsername retrieves the moderator with the given username.
func GetModeratorByUsername(username string, collection *mongo.Collection) (moderator entities.Moderator, err error) {

	if username == "" {
		err = xerrors.New("GetModeratorByUsername: empty username")
		return
	}

	//
	//TODO: aggiungere cache
	chat, err := tdlib.ResolveUsername(username)
	if err != nil {
		return
	}

	if chat.Type.ChatTypeType() != client.TypeChatTypePrivate {
		err = xerrors.Errorf("GetModeratorByUsername: the username %s does not belong to a person", username)
		return
	}

	return GetModeratorByTelegramID(chat.ID, collection)

}

/*
 ***********************************************************************************************************************
 *																													   *
 *													INSERT															   *
 *																													   *
 ***********************************************************************************************************************
 */

// AddModerator adds the given user to the moderators.
func AddModerator(user *tgbotapi.ChatMember, moddedByTelegramID int64, collection *mongo.Collection) (generatedID string, err error) {

	if user == nil {
		err = xerrors.New("AddModerator: user nil")
		return
	}

	moderator := entities.Moderator{
		TelegramID:         user.User.ID,
		CanChangeInfo:      user.CanChangeInfo,
		CanDeleteMessages:  user.CanDeleteMessages,
		CanInviteUsers:     user.CanInviteUsers,
		CanRestrictMembers: user.CanRestrictMembers,
		CanPinMessages:     user.CanPinMessages,
		CanPromoteMembers:  user.CanPromoteMembers,
		PromotedBy:         moddedByTelegramID,
		PromotionDate:      time.Now(),
	}

	// Telegram will return all false values for the creator
	if user.IsCreator() {
		moderator.CanChangeInfo = true
		moderator.CanDeleteMessages = true
		moderator.CanInviteUsers = true
		moderator.CanRestrictMembers = true
		moderator.CanPinMessages = true
		moderator.CanPromoteMembers = true
		moderator.IsAdmin = true
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, moderator)
	if err != nil {
		err = xerrors.Errorf("AddModerator: unable to add moderator into the document store: %s", err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return generatedID, err

}

func RemoveModerator(unmoddedUserID int64, collection *mongo.Collection) (err error) {

	if unmoddedUserID == 0 {
		err = xerrors.New("RemoveModerator: unmoddedUserID 0")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	filter := bson.M{"telegramid": unmoddedUserID}
	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		err = xerrors.Errorf("AddModerator: unable to add moderator into the document store: %s", err)
	}

	return

}

/*
 ***********************************************************************************************************************
 *																													   *
 *													UPDATE														   *
 *																													   *
 ***********************************************************************************************************************
 */

// UpdateModeratorsDetails updates the details of the moderators in the database.
func UpdateModeratorsDetails(chatID int64, bot *tgbotapi.BotAPI, collection *mongo.Collection) error {

	chatAdministratorsConfig := tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
	}

	chatAdministrators, err := bot.GetChatAdministrators(chatAdministratorsConfig)
	if err != nil {
		return xerrors.Errorf("UpdateModeratorsDetails: GetChatAdministrators: %s", err)
	}

	for _, chatAdministrator := range chatAdministrators {

		if IsModerator(chatAdministrator.User.ID, collection) {

			err = UpdateModeratorDetails(&chatAdministrator, collection)
			if err != nil {
				return xerrors.Errorf("UpdateModeratorsDetails: %s", err)
			}

			continue

		}

		_, err = AddModerator(&chatAdministrator, 0, collection)
		if err != nil {
			return xerrors.Errorf("UpdateModeratorsDetails: %s", err)
		}
	}

	return nil

}

// UpdateModeratorDetailsByTelegramUser updates the details of a user in the database.
func UpdateModeratorDetailsByTelegramUser(user *tgbotapi.User, chatID int64, bot *tgbotapi.BotAPI, collection *mongo.Collection) error {

	chatMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: user.ID,
		},
	})

	if err != nil {
		return xerrors.Errorf("UpdateModeratorDetailsByTelegramUser: %s", err)
	}

	return UpdateModeratorDetails(&chatMember, collection)

}

// UpdateModeratorDetails updates the details of a moderator in the database.
func UpdateModeratorDetails(chatMember *tgbotapi.ChatMember, collection *mongo.Collection) error {

	if chatMember.IsCreator() {
		return nil
	}

	//
	filter := bson.M{"telegramid": chatMember.User.ID}
	update := bson.D{
		{
			"$set", bson.M{
				"canchangeinfo":      chatMember.CanChangeInfo,
				"candeletemessages":  chatMember.CanDeleteMessages,
				"caninviteusers":     chatMember.CanInviteUsers,
				"canrestrictmembers": chatMember.CanRestrictMembers,
				"canpinmessages":     chatMember.CanPinMessages,
				"canpromotemembers":  chatMember.CanPromoteMembers,
			},
		},
	}
	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	result, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("UpdateModeratorDetails: unable to update moderator with id %d: %s", chatMember.User.ID, err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("UpdateModeratorDetails: no match for moderator with id %d", chatMember.User.ID)
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

// IsModeratorUsername returns true if the given username
// belongs to a moderator.
func IsModeratorUsername(username string, collection *mongo.Collection) bool {
	_, err := GetModeratorByUsername(username, collection)
	return err == nil
}

// IsModerator returns true if the given telegram id belongs to a moderator.
func IsModerator(telegramID int64, collection *mongo.Collection) bool {
	_, err := GetModeratorByTelegramID(telegramID, collection)
	return err == nil
}
