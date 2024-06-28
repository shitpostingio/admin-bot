package structs

// AntiUserbotConfiguration represents the antiuserbot configuration
type AntiUserbotConfiguration struct {
	JoinThreshold   int `reloadable:"true"`
	RoutineLifespan int `reloadable:"true"`
}
