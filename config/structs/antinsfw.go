package structs

// AntiNSFWConfiguration represents the antinsfw configuration
type AntiNSFWConfiguration struct {
	APIKey            string
	APIEndpoint       string
	ExplicitThreshold int `reloadable:"true"`
	RacyThreshold     int `reloadable:"true"`
}
