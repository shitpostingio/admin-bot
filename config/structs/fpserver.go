package structs

// FPServerConfiguration represents the fpserver configuration
type FPServerConfiguration struct {
	Address                  string `reloadable:"true"`
	AnalysisImageEndpoint    string `reloadable:"true"`
	AnalysisVideoEndpoint    string `reloadable:"true"`
	FingerprintImageEndpoint string `reloadable:"true"`
	FingerprintVideoEndpoint string `reloadable:"true"`
	GibberishEndpoint        string `reloadable:"true"`
	AuthorizationHeaderName  string `reloadable:"true"`
	AuthorizationHeaderValue string `reloadable:"true"`
	CallerAPIKeyHeaderName   string `reloadable:"true"`
	GibberishInputHeaderName string `reloadable:"true"`
	FilePathHeaderName       string `reloadable:"true"`
	FileSizeThreshold        int    `type:"optional" reloadable:"true"`
}
