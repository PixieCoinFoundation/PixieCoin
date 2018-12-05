package rao

import (
	"constants"
	"fmt"
)

func GetHostChatListKey(hostID int64) string {
	return fmt.Sprintf("%s%d", constants.RT_PARTY_CHAT_LIST_KEY_PREFIX, hostID)
}

//根据段位获得实时舞会房间列表 按人数排序
func GetHostListKey(ds string) string {
	return constants.RT_PARTY_HOST_LIST_KEY_PREFIX + ds
}

func GetHostSelfKey(hostID int64) string {
	return fmt.Sprintf("%s%d", constants.RT_PARTY_HOST_KEY_PREFIX, hostID)
}

func GetHostJoinerKey(hostID int64, username string) string {
	return fmt.Sprintf("%s%d_%s", constants.RT_PARTY_HOST_JOINER_KEY_PREFIX, hostID, username)
}

// func GetJoinerDressKey(username string, hostID int64) string {
// 	return fmt.Sprintf("%s%d_%s", constants.RT_PARTY_JOINER_DRESS_KEY_PREFIX, hostID, username)
// }

// func GetJoinerSyncTimeKey(username string, hostID int64) string {
// 	return fmt.Sprintf("%s%d_%s", constants.RT_PARTY_JOINER_SYNC_TIME_KEY_PREFIX, hostID, username)
// }

func GetHostJoinerHashKey(username string, hostID int64) string {
	return fmt.Sprintf("%s%d_%s", constants.RT_PARTY_JOINER_HASH_KEY_PREFIX, hostID, username)
}
