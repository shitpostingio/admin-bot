package api

import (
	"github.com/shitpostingio/admin-bot/api/tdlib"
	limiter "github.com/shitpostingio/admin-bot/ratelimiter"
)

// SendCallbackWithAlert sends a callback response that shows an alert.
// It will also use a rate limiter not to get restricted by Telegram.
func SendCallbackWithAlert(id, text string) error {
	limiter.AuthorizeAction()
	return tdlib.SendCallback(id, text, true)
}

// SendCallback sends a callback response.
// It will also use a rate limiter not to get restricted by Telegram.
func SendCallback(id, text string) error {
	limiter.AuthorizeAction()
	return tdlib.SendCallback(id, text, false)
}
