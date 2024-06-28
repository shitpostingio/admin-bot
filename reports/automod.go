package reports

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/localization"
)

func ChatMessageReported(reporteeUserID int64, reporteeUserName string, reportedUserID int64, reportedUserName string) string {
	return fmt.Sprintf(localization.GetString("automod_new_report_by_with_reply"),
		reporteeUserID, reporteeUserName, reporteeUserID,
		reportedUserID, reportedUserName, reportedUserID)
}

func ChatMessageReport(reporteeUserID int64, reporteeUserName string) string {
	return fmt.Sprintf(localization.GetString("automod_new_report_by"),
		reporteeUserID, reporteeUserName, reporteeUserID)
}

func ForwardFromChannel(userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("automod_removed_message_forwarded_channel"),
		userID, userName, userID)
}

func ForwardFromBlacklistedHandle(userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("automod_removed_message_forwared_blacklisted_handle"),
		userID, userName, userID)
}

func MessageSentViaBlacklistedInlineBot(userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("automod_removed_message_via_blacklsited_inline_bot"),
		userID, userName, userID)
}

func RemovedNSFWMedia(userID int64, userName string, label string, confidence float64) string {
	return fmt.Sprintf(localization.GetString("automod_removed_nsfw_media"),
		userID, userName, userID,
		label, confidence)
}

func RemovedUnwantedHandle(userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("automod_removed_unwanted_handle"),
		userID, userName, userID)
}

func RemovedUnwantedLink(host string, userID int64, userName string) string {
	return fmt.Sprintf(localization.GetString("automod_removed_unwanted_link"),
		host,
		userID, userName, userID)
}

func UserMutedForUnwantedLink(botID int64, botName string,
	userID int64, userName string,
	reason string) string {
	return fmt.Sprintf(localization.GetString("automod_unwanted_link_user_muted"),
		botID, botName,
		userID, userName, userID,
		reason)
}
