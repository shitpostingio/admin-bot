package documentstore

import (
	"context"
	"fmt"
	"math"
	"time"

	fpcompare "github.com/shitpostingio/image-fingerprinting/comparer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"

	"github.com/shitpostingio/admin-bot/analysisadapter"
	"github.com/shitpostingio/admin-bot/entities"
)

const (
	mediaApproximation = 0.08 //TODO: rendere variabile
)

// FindMediaByFileID finds a media in the database via its file id
func FindMediaByFileID(uniqueFileID, fileID string, collection *mongo.Collection) (media entities.Media, err error) {

	media, err = FindMediaByFileUniqueID(uniqueFileID, collection)
	if err == nil {
		return
	}

	if fileID == "" {
		err = xerrors.New("FindMediaByFileID: empty fileID")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"fileid": fileID}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&media)
	if err == nil { //match via fileID
		return
	}

	//
	if !entities.MediaCanBeFingerprinted(fileID) {
		err = xerrors.Errorf("FindMediaByFileID no match for fileID %s", fileID)
		return
	}

	fingerprint, err := analysisadapter.GetFingerprint(uniqueFileID, fileID)
	fmt.Println("Fingerprint", fingerprint)
	if err != nil {
		err = xerrors.Errorf("FindMediaByFileID could not get fingerprint values for fileID %s", fileID)
		return
	}

	media, err = FindMediaByFeatures(fingerprint.Histogram, fingerprint.PHash, mediaApproximation, collection)
	return

}

func FindMediaByFileUniqueID(fileUniqueID string, collection *mongo.Collection) (media entities.Media, err error) {

	if fileUniqueID == "" {
		err = xerrors.New("FindMediaByFileUniqueID: empty fileID")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"fileuniqueid": fileUniqueID}
	err = collection.FindOne(ctx, filter, options.FindOne()).Decode(&media)
	return media, err

}

// FindMediaByFeatures finds a media by its features
func FindMediaByFeatures(histogram []float64, pHash string, approximation float64, collection *mongo.Collection) (media entities.Media, err error) {

	//
	if histogram == nil {
		err = xerrors.New("FindMediaByFeatures: histogram was nil")
		return
	}

	if pHash == "" {
		err = xerrors.New("FindMediaByFeatures: pHash was empty")
		return
	}

	//
	average, sum := entities.GetHistogramAverageAndSum(histogram)
	minAvg := math.Trunc(average - 1)
	maxAvg := math.Ceil(average + 1)
	minSum := math.Trunc(sum - (sum * approximation))
	maxSum := math.Ceil(sum + (sum * approximation))

	//
	filter := bson.D{
		{
			Key: "histogramaverage",
			Value: bson.D{
				{"$gte", minAvg},
				{"$lte", maxAvg},
			},
		},
		{Key: "histogramsum",
			Value: bson.D{
				{"$gte", minSum},
				{"$lte", maxSum},
			},
		},
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//TODO: ordinare secondo qualcosa i dati

	//
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		err = xerrors.Errorf("FindMediaByFeatures: unable to retrieve media: %s", err)
		return
	}

	media, err = findBestMatch(pHash, cursor)
	if err != nil {
		err = xerrors.Errorf("FindMediaByFeatures: %s", err)
		return
	}

	return

}

// FindMediaByID finds a media given its ObjectID
func FindMediaByID(ID *primitive.ObjectID, collection *mongo.Collection) (media entities.Media, err error) {

	//
	if ID == nil {
		err = xerrors.New("FindMediaByID: ID was nil")
		return
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"_id": ID}
	err = collection.FindOne(ctx, filter).Decode(&media)
	return

}

/*
 ***********************************************************************************************************************
 *																													   *
 *														INSERT														   *
 *																													   *
 ***********************************************************************************************************************
 */

// BlacklistMedia blacklists a media
func BlacklistMedia(uniqueFileID, fileID string, userID int64, collection *mongo.Collection) (generatedID string, err error) {

	//
	if fileID == "" && uniqueFileID == "" {
		err = xerrors.New("BlacklistMedia: empty fileID and unique id")
		return
	}

	//
	foundMedia, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	if err == nil {

		err = markMediaAsBlacklisted(&foundMedia.ID, collection)
		if err != nil {

			err = xerrors.Errorf("BlacklistMedia: %s", err)
			return

		}

		generatedID = foundMedia.ID.Hex()
		return

	}

	//
	media := entities.Media{
		FileUniqueID: uniqueFileID,
		FileID:       fileID,
		LastEditedBy: userID,
		LastModified: time.Now(),
	}

	//
	if entities.MediaCanBeFingerprinted(fileID) {

		fingerprint, err := analysisadapter.GetFingerprint(uniqueFileID, fileID)
		if err == nil {

			average, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)
			media.Histogram = fingerprint.Histogram
			media.HistogramAverage = average
			media.HistogramSum = sum
			media.PHash = fingerprint.PHash

		}

	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, media)
	if err != nil {
		err = xerrors.Errorf("BlacklistMedia: unable to add media with fileID %s into the document store: %s", fileID, err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return

}

// WhitelistMedia whitelists a media
func WhitelistMedia(uniqueFileID, fileID string, userID int64, collection *mongo.Collection) (generatedID string, err error) {

	if fileID == "" {
		err = xerrors.New("WhitelistMedia: empty fileID")
		return
	}

	foundMedia, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	if err == nil {

		err = markMediaAsWhitelisted(&foundMedia.ID, collection)
		if err != nil {

			err = xerrors.Errorf("WhitelistMedia: %s", err)
			return

		}

		generatedID = foundMedia.ID.Hex()
		return

	}

	//
	media := entities.Media{
		FileUniqueID:  uniqueFileID,
		FileID:        fileID,
		IsWhitelisted: true,
		LastEditedBy:  userID,
		LastModified:  time.Now(),
	}

	//
	if entities.MediaCanBeFingerprinted(fileID) {

		fingerprint, err := analysisadapter.GetFingerprint(uniqueFileID, fileID)
		if err == nil {

			average, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)
			media.Histogram = fingerprint.Histogram
			media.HistogramAverage = average
			media.HistogramSum = sum
			media.PHash = fingerprint.PHash

		}

	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, media)
	if err != nil {
		err = xerrors.Errorf("WhitelistMedia: unable to add media with fileID %s into the document store: %s", fileID, err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return

}

// RemoveMedia removes a media from the document store
func RemoveMedia(uniqueFileID, fileID string, collection *mongo.Collection) error {

	//
	if fileID == "" {
		return xerrors.New("RemoveMedia: empty file ID")
	}

	// Find matching media
	foundMedia, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	if err != nil {
		return xerrors.Errorf("RemoveMedia: no matching media found for fileID %s: %s", fileID, err)
	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"_id": foundMedia.ID}
	result, err := collection.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return xerrors.Errorf("RemoveMedia: error while removing media %s: %s", fileID, err)
	}

	if result.DeletedCount == 0 {
		return xerrors.Errorf("RemoveMedia: no matches found for media %s", fileID)
	}

	return nil
}

// BlacklistNSFWMedia blacklists a media for being NSFW
func BlacklistNSFWMedia(uniqueFileID, fileID string, description string, score float64, userID int64, collection *mongo.Collection) (generatedID string, err error) {

	//
	if fileID == "" {
		err = xerrors.New("BlacklistNSFWMedia: empty fileID")
		return
	}

	//
	foundMedia, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	if err == nil {

		err = markMediaAsNSFW(&foundMedia.ID, description, score, collection)
		if err != nil {

			err = xerrors.Errorf("BlacklistNSFWMedia: %s", err)
			return

		}

		generatedID = foundMedia.ID.Hex()
		return

	}

	//
	media := entities.Media{
		FileUniqueID:    uniqueFileID,
		FileID:          fileID,
		NSFWDescription: description,
		NSFWScore:       score,
		LastEditedBy:    userID,
		LastModified:    time.Now(),
	}

	//
	if entities.MediaCanBeFingerprinted(fileID) {

		fingerprint, err := analysisadapter.GetFingerprint(uniqueFileID, fileID)
		if err == nil {

			average, sum := entities.GetHistogramAverageAndSum(fingerprint.Histogram)
			media.Histogram = fingerprint.Histogram
			media.HistogramAverage = average
			media.HistogramSum = sum
			media.PHash = fingerprint.PHash

		}

	}

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), opDeadline)
	defer cancelCtx()

	result, err := collection.InsertOne(ctx, media)
	if err != nil {
		err = xerrors.Errorf("BlacklistNSFWMedia: unable to add media with fileID %s into the document store: %s", fileID, err)
		return
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		generatedID = objectID.Hex()
	}

	return

}

/*
***********************************************************************************************************************
*																													   *
*														UPDATE														   *
*																													   *
***********************************************************************************************************************
 */

func markMediaAsBlacklisted(ID *primitive.ObjectID, collection *mongo.Collection) error {

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

	result, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("markMediaAsBlacklisted: unable to mark media with _id %s as blacklisted: %s", ID.Hex(), err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("markMediaAsBlacklisted: no match for media with _id %s", ID.Hex())
	}

	return nil

}

func markMediaAsWhitelisted(ID *primitive.ObjectID, collection *mongo.Collection) error {

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

	result, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("markMediaAsWhitelisted: unable to mark media with _id %s as whitelisted: %s", ID.Hex(), err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("markMediaAsWhitelisted: no match for media with _id %s", ID.Hex())
	}

	return nil

}

func markMediaAsNSFW(ID *primitive.ObjectID, description string, score float64, collection *mongo.Collection) error {

	//
	filter := bson.M{"_id": ID}
	update := bson.D{
		{
			"$set", bson.M{
				"iswhitelisted":   false,
				"nsfwdescription": description,
				"nsfwscore":       score,
			},
		},
	}

	updCtx, cancelUpdCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelUpdCtx()

	result, err := collection.UpdateOne(updCtx, filter, update)
	if err != nil {
		return xerrors.Errorf("markMediaAsNSFW: unable to add nsfw data to media with _id %s: %s", ID.Hex(), err)
	}

	if result.MatchedCount == 0 {
		return xerrors.Errorf("markMediaAsNSFW: no match for media with _id %s", ID.Hex())
	}

	return nil

}

/*
***********************************************************************************************************************
*																													   *
*														CHECK														   *
*																													   *
***********************************************************************************************************************
 */

// MediaIsBlacklisted returns true if the media is blacklisted
func MediaIsBlacklisted(uniqueFileID, fileID string, collection *mongo.Collection) bool {

	if fileID == "" {
		return false
	}

	media, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	return err == nil && !media.IsWhitelisted

}

// MediaIsWhitelisted returns true if the media is whitelisted
func MediaIsWhitelisted(uniqueFileID, fileID string, collection *mongo.Collection) bool {

	if fileID == "" {
		return false
	}

	media, err := FindMediaByFileID(uniqueFileID, fileID, collection)
	return err == nil && media.IsWhitelisted

}

func findBestMatch(referencePHash string, cursor *mongo.Cursor) (media entities.Media, err error) {

	defer func() {
		_ = cursor.Close(dsCtx)
	}()

	i := 0
	for cursor.Next(context.TODO()) {

		i++
		// Support variable. If we deserialize directly in media,
		// since IsWhitelisted is an omitempty field, it won't be
		// deserialized in case of it being missing. This way, if
		// a document with it set to true has already been retrieved,
		// it will always keep being true.
		var res entities.Media
		err = cursor.Decode(&res)
		if err == nil && fpcompare.PhotosAreSimilarEnough(referencePHash, res.PHash) {
			media = res
			fmt.Println("match in ", i, "iterations. FileID", media.FileID, "whitelist", media.IsWhitelisted, "_id", media.ID)
			return
		}

	}

	err = xerrors.New("no match found")
	return

}
