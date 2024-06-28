package telegram

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

func GetFileType(fileID string) string {
	return fileID[:3]
}

func MediaCanBeAnalyzed(fileID string) bool {

	switch GetFileType(fileID) {
	case STICKER:
		return true
	case PHOTO:
		return true
	case VIDEO:
		return true
	case ANIMATION:
		return true
	}

	return false

}
