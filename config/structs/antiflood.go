package structs

// AntiFloodConfiguration represents the antiflood configuration
type AntiFloodConfiguration struct {
	Threshold       int `reloadable:"true"`
	RoutineLifeSpan int `reloadable:"true"`
}
