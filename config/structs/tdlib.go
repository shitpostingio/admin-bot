package structs

// TdlibConfiguration represents the tdlib configuration
type TdlibConfiguration struct {
	LogVerbosityLevel int32 `type:"optional"`
	UseTestDc         bool  `type:"optional"`
	DatabaseDirectory string
	FilesDirectory    string
	APIID             int64
	APIHash           string
}
