package structs

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//GroupUserRestrictions represents the restrictions that a group member can have
type GroupUserRestrictions struct {
	UntilDate             int64
	CanSendMessages       bool
	CanSendMediaMessages  bool
	CanSendOtherMessages  bool
	CanAddWebPagePreviews bool
}

//AntiSpam maps userIDs to goroutine channels
type AntiSpam struct {
	InputChannel    chan *tgbotapi.Message
	EndCycleChannel chan int
	UserChannels    map[int]chan *tgbotapi.Message
}

//CloudmersiveNSFWResult represents the data returned by Cloudmersive
type CloudmersiveNSFWResult struct {
	Successful            bool    `json:"Successful"`
	Score                 float64 `json:"Score"`
	ClassificationOutcome string  `json:"ClassificationOutcome"`
}

//AWSNSFWResult represents the data returned by Amazon Rekognition
type AWSNSFWResult struct {
	ExplicitNudityConfidence float64
	SuggestiveConfidence     float64
	BestMatch                string
}
