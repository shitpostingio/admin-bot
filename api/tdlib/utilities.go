package tdlib

func getTdlibMessageID(botApiMessageID int) int64 {
	return int64(botApiMessageID * tdlibMessageConst)
}
