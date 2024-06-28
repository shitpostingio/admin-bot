package reports

import (
	"fmt"
	"github.com/shitpostingio/admin-bot/localization"
)

func ModeratorDemoted(adminID int64, adminName string, moderatorID int64, moderatorName string) string {
	return fmt.Sprintf(localization.GetString("moderator_demoted"),
		adminID, adminName, adminID,
		moderatorID, moderatorName, moderatorID)
}

func ModeratorCannotBeRemovedFromTable(moderatorID int64, moderatorName string) string {
	return fmt.Sprintf(localization.GetString("moderator_cannot_be_removed"),
		moderatorID, moderatorName, moderatorID)
}

func UserBanned(moderatorID int64, moderatorName string, userID int64, userName string, reason string) string {
	return fmt.Sprintf(localization.GetString("user_banned"),
		moderatorID, moderatorName,
		userID, userName, userID,
		reason)
}

func UnauthorizedPrivateGroup(groupName string, groupID int64, adderName string, adderID int64) string {
	return fmt.Sprintf(localization.GetString("unauthorized_group_report"),
		groupName, "", groupID,
		adderName, adderID, adderID)
}

func UnauthorizedPublicGroup(groupName string, groupHandle string, groupID int64, adderName string, adderID int64) string {
	handlePart := fmt.Sprintf(localization.GetString("unauthorized_report_handle_part"), groupHandle)
	return fmt.Sprintf(localization.GetString("unauthorized_group_report"),
		groupName, handlePart, groupID,
		adderName, adderID, adderID)
}

func UnauthorizedPrivateChannel(channelName string, channelID int64) string {
	return fmt.Sprintf(localization.GetString("unauthorized_channel_report"),
		channelName, "", channelID)
}

func UnauthorizedPublicChannel(channelName string, channelHandle string, channelID int64) string {
	handlePart := fmt.Sprintf(localization.GetString("unauthorized_report_handle_part"), channelHandle)
	return fmt.Sprintf(localization.GetString("unauthorized_channel_report"),
		channelName, handlePart, channelID)
}
