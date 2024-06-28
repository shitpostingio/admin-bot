package structs

// RateLimiterConfiguration represents the rate limiter configuration
type RateLimiterConfiguration struct {
	MaxActionsPerSecond int `reloadable:"true"`
}
