package structs

// AntiSpamConfiguration represents the antispam configuration
type AntiSpamConfiguration struct {
	RoutineLifeSpan int `reloadable:"true"`
	TextThreshold   int `reloadable:"true"`
	MediaThreshold  int `reloadable:"true"`
	OtherThreshold  int `reloadable:"true"`
}
