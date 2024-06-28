package analysisadapter

import "fmt"

//nolint
const (
	PHOTO     = "AgA"
	VIDEO     = "BAA"
	ANIMATION = "CgA"
	STICKER   = "CAA"
	VOICE     = "AwA"
	DOCUMENT  = "BQA"
	AUDIO     = "CQA"
	VIDEONOTE = "DQA"
)

func getAnalysisEndpoint(fileID, fileUniqueID string) string {

	// Telegram prefixes are 3 characters long
	fileIDPrefix := fileID[:3]

	switch fileIDPrefix {
	case PHOTO:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.AnalysisImageEndpoint, fileUniqueID)
	case VIDEO:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.AnalysisVideoEndpoint, fileUniqueID)
	case ANIMATION:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.AnalysisVideoEndpoint, fileUniqueID)
	case STICKER:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.AnalysisImageEndpoint, fileUniqueID)
	case VOICE:
		fallthrough
	case DOCUMENT:
		fallthrough
	case AUDIO:
		fallthrough
	case VIDEONOTE:
		fallthrough
	default:
		return ""
	}

}

func getFingerprintEndpoint(fileID, fileUniqueID string) string {

	// Telegram prefixes are 3 characters long
	fileIDPrefix := fileID[:3]

	switch fileIDPrefix {
	case PHOTO:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.FingerprintImageEndpoint, fileUniqueID)
	case VIDEO:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.FingerprintVideoEndpoint, fileUniqueID)
	case ANIMATION:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.FingerprintVideoEndpoint, fileUniqueID)
	case STICKER:
		return fmt.Sprintf("%s/%s/%s", cfg.Address, cfg.FingerprintImageEndpoint, fileUniqueID)
	case VOICE:
		fallthrough
	case DOCUMENT:
		fallthrough
	case AUDIO:
		fallthrough
	case VIDEONOTE:
		fallthrough
	default:
		return ""
	}

}

func getGibberishEndpoint() string {
	return fmt.Sprintf("%s/%s", cfg.Address, cfg.GibberishEndpoint)
}
