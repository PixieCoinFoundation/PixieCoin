package rao

import (
	"constants"
	"fmt"
	"strconv"
	"time"
)

func GetUserDesignMsgTimeInfoKey(token string) string {
	return fmt.Sprintf("dmi_%s", token)
}

func GetUserDesignMsgListKey(username string, typee int, timeToken string) string {
	return fmt.Sprintf("%s_%d_%s", username, typee, timeToken)
}

func GetNewDesignMsgListkey() string {
	return constants.NEW_GLOBAL_DESIGN_MSG_LIST_KEY
}

func GetTheaterRankKey() string {
	return constants.THEATER_RANK_KEY_PREFIX + fmt.Sprintf("%d", time.Now().UnixNano())
}

func GetTempSaleRankKey(st, et int64) string {
	return constants.TEMP_SALE_RANK_KEY_PREFIX + fmt.Sprintf("%d_", st) + fmt.Sprintf("%d", et)
}

func GetSinglePushJobListKey() string {
	return constants.SINGLE_PUSH_JOB_LIST_KEY
}

func GetHotDeisngerRespKey() string {
	return constants.GLOBAL_HOT_DESIGNER_RESP_KEY
}

func GetDayNewCustomListKey(ds string) string {
	return constants.DAY_NEW_CUSTOM_LIST_KEY_PREFIX + ds
}

func GetDayOrderCntAndAmountKey(t time.Time) (string, string) {
	tt := t.Format("20060102")
	return constants.DAY_ORDER_CNT_KEY_PREFIX + tt, constants.DAY_ORDER_AMOUNT_KEY_PREFIX + tt
}

func GetDayPlayerNewCntKey(t time.Time) string {
	return constants.DAY_NEW_PLAYER_CNT_KEY_PREFIX + t.Format("20060102")
}
func GetHourChannelNewPlayerCntKey(t time.Time, channel string) string {
	return constants.CHANNEL_NEW_PLAYER_CNT_KEY_PREFIX + channel + "_" + t.Format("15")
}

func GetHourAllChannelNewPlayerCntKey(t time.Time) string {
	return constants.ALL_CHANNEL_NEW_PLAYER_CNT_KEY_PREFIX + "_" + t.Format("15")
}

//GetActivationCode1Key 获取List队列名称
func GetActivationCode1Key() string {
	return constants.ACTIVATION_CODE_GAME1
}

//GetActivationCode2Key 获取List队列名称
func GetActivationCode2Key() string {
	return constants.ACTIVATION_CODE_GAME2
}

//GetActivationCode3Key 获取List队列名称
func GetActivationCode3Key() string {
	return constants.ACTIVATION_CODE_GAME3
}

func GetHourChannelNewOrderCntKey(t time.Time, channel string) string {
	return constants.CHANNEL_NEW_ORDER_CNT_KEY_PREFIX + channel + "_" + t.Format("15")
}

func GetHourAllChannelNewOrderCntKey(t time.Time) string {
	return constants.ALL_CHANNEL_NEW_ORDER_CNT_KEY_PREFIX + "_" + t.Format("15")
}

func GetHourChannelNewOrderAmountKey(t time.Time, channel, currency string) string {
	return constants.CHANNEL_NEW_ORDER_AMOUNT_KEY_PREFIX + channel + "_" + t.Format("15") + "_" + currency
}

func GetHourAllChannelNewOrderAmountKey(t time.Time, currency string) string {
	return constants.ALL_CHANNEL_NEW_ORDER_AMOUNT_KEY_PREFIX + "_" + t.Format("15") + "_" + currency
}

func GetGuildMedalRankKey(gwt string) string {
	return constants.GUILD_GWEEK_MEDAL_RANK_CACHE_KEY_PREFIX + gwt
}

func GetGuildActivityRankKey(wt string) string {
	return constants.GUILD_WEEK_ACTIVITY_RANK_CACHE_KEY_PREFIX + wt
}

func GetGuildMedalGWeekKey(gwt string) string {
	return constants.GUILD_GWEEK_MEDAL_RANK_SET_KEY_PREFIX + gwt
}

func GetGuildActivityWeekKey(wt string) string {
	return constants.GUILD_WEEK_ACTIVITY_RANK_SET_KEY_PREFIX + wt
}

func GetDesignerPassedMonthKey(monthToken string) string {
	return constants.DESIGNER_PASSED_MONTH_KEY_PREFIX + monthToken
}

func GetDesignerPassedExpMonthKey(monthToken string) string {
	return constants.DESIGNER_PASSED_EXP_MONTH_KEY_PREFIX + monthToken
}

func GetDesignerTotalDiamondMonthKey(monthToken string) string {
	return constants.DESIGNER_TOTAL_DIAMOND_MONTH_KEY_PREFIX + monthToken
}

func GetDesignerTotalGoldMonthKey(monthToken string) string {
	return constants.DESIGNER_TOTAL_GOLD_MONTH_KEY_PREFIX + monthToken
}

//GetDesignerTotalCoinMonthKey 获取设计师币的key
func GetDesignerTotalCoinMonthKey(monthToken string) string {
	return constants.DESIGNER_TOTAL_COIN_MONTH_KEY_PREFIX + monthToken
}

func GetDesignerTotalSaleMonthKey(monthToken string) string {
	return constants.DESIGNER_TOTAL_SALE_MONTH_KEY_PREFIX + monthToken
}

func GetPKWeekRankKey(token string) string {
	return constants.PK_WEEK_RANK_KEY_PREFIX + token
}

func GetPKGWeekRankKey(token string) string {
	return constants.PK_GWEEK_RANK_KEY_PREFIX + token
}

func GetBoardFirstPageKey(username string) string {
	return constants.BOARD_FISR_PAGE_KEY_PREFIX + username
}

func GetPartyInviteFlowerKey(partyID int64) string {
	return constants.PARTY_INVITE_FLOWER_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetPartyInviteAttendKey(partyID int64) string {
	return constants.PARTY_INVITE_ATTEND_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetGlobalUserSetKey() string {
	return constants.GLOBAL_USER_SET_KEY
}

func GetBoardNewSignKey(username string) string {
	return constants.BOARD_NEW_SIGN_KEY_PREFIX + "_" + username
}

func GetDesignerNewKey(username string) string {
	return constants.DESIGNER_NEW_KEY_PREFIX + username
}

func GetDesignerInfoKey(username string) string {
	return constants.GLOBAL_USER_HASH_KEY_PREFIX + username
}

// func GetUserCloPieceKey(username string) string {
// 	return constants.USER_CLO_PIECE_KEY_PREFIX + username
// }

func GetCustomKey(cid int64) string {
	return constants.CUSTOM_INFO_KEY_PREFIX + strconv.FormatInt(cid, 10)
}

func GetMailListKey(username string) string {
	return constants.MAIL_SET_PREFIX + username
}

func GetPKKey(lid string, pkType int) string {
	return constants.GF_PK_LIST_KEY_PREFIX + lid + "_" + strconv.Itoa(pkType)
}

func GetTopCosItemListKey(cosplayID int64) string {
	return constants.COS_ITEM_LIST_TOP_KEY_PREFIX + strconv.FormatInt(cosplayID, 10)
}

func GetTopCosItemKey(cosplayID int64) string {
	return constants.COS_TOP_ITEM_KEY_PREFIX + strconv.FormatInt(cosplayID, 10)
}

func GetCosItemRankSetKey(cosplayID int64) string {
	return constants.COS_ITEM_LIST_SCORE_KEY_PREFIX + strconv.FormatInt(cosplayID, 10)
}

func GetCosItemUserListKey(cosplayID int64, username string) string {
	return constants.COS_ITEM_USER_KEY_PREFIX + strconv.FormatInt(cosplayID, 10) + "_" + username
}

func GetCosItemUserMarkListKey(cosplayID int64, itemID int64) string {
	return constants.COS_ITEM_USER_MARK_LIST_KEY_PREFIX + strconv.FormatInt(cosplayID, 10) + "_" + strconv.FormatInt(itemID, 10)
}

func GetCosItemIdSetKey(cosplayID int64) string {
	return constants.COS_ITEM_ID_SET_PREFIX + strconv.FormatInt(cosplayID, 10)
}

func GetCosItemSelfKey(itemID int64) string {
	return constants.COS_ITEM_SELF_KEY_PREFIX + strconv.FormatInt(itemID, 10)
}

func GetHelpSelfKey(helpID int64) string {
	return constants.HELP_KEY_PREFIX + strconv.FormatInt(helpID, 10)
}

func GetHelpCmtListKey(helpID int64) string {
	return constants.HELP_COMMENT_LIST_KEY_PREFIX + strconv.FormatInt(helpID, 10)
}

func GetHelpCmtCntKey(helpID int64) string {
	return constants.HELP_CMT_CNT_KEY_PREFIX + strconv.FormatInt(helpID, 10)
}

// func GetUserCustomListKey(username string) string {
// 	return constants.CUSTOM_USER_LIST_KEY_PREFIX + username
// }

// func GetGlobalCustomPassedListKey(cloTypeAndMoelNo string) string {
// 	return constants.CUSTOM_PASSED_LIST_KEY + "_" + cloTypeAndMoelNo
// }

func GetCosItemCmtListKey(itemID int64) string {
	return constants.COS_ITEM_CMT_KEY_PREFIX + strconv.FormatInt(itemID, 10)
}

func GetCosItemHasCmtKey(itemID int64) string {
	return constants.COS_ITEM_HAS_CMT_KEY_PREFIX + strconv.FormatInt(itemID, 10)
}

func GetHelpHasCmtKey(helpID int64) string {
	return constants.HELP_HAS_CMT_KEY_PREFIX + strconv.FormatInt(helpID, 10)
}

func GetNoticeExistKey(nid int64) string {
	return constants.NOTICE_EXIST_KEY_PREFIX + strconv.FormatInt(nid, 10)
}

func GetMailExistKey(mailID int64) string {
	return constants.MAIL_EXIST_KEY_PREFIX + strconv.FormatInt(mailID, 10)
}

// func GetPartyListKey(typee int) string {
// if typee == constants.PARTY_CASUAL_TYPE {
// nowh := time.Now().Format("15")
//key name with hour info
// return constants.PARTY_LIST_KEY_PREFIX + strconv.Itoa(typee) + "_" + nowh, constants.PARTY_CNT_KEY_PREFIX + nowh
// } else {
// return constants.PARTY_LIST_KEY_PREFIX + strconv.Itoa(typee)
// }
// }

func GetPartyCasualCntKey() string {
	return constants.PARTY_CNT_KEY_PREFIX + time.Now().Format("15")
}

// func GetCasultPartyListKey(gapHour int) string {
// 	return constants.PARTY_LIST_KEY_PREFIX + strconv.FormatInt(constants.PARTY_CASUAL_TYPE, 10) + "_" + time.Now().Add(time.Duration(gapHour)*time.Hour).Format("15")
// }

func GetGlobalCasualPartyListKey() string {
	return constants.PARTY_CASUAL_GLOBAL_LIST
}

func GetGlobalPrizePartyListKey() string {
	return constants.PARTY_PRIZE_GLOBAL_LIST
}

func GetGlobalCosplayPartyListKey() string {
	return constants.PARTY_COSPLAY_GLOBAL_LIST
}

func GetPartySelfKey(partyID int64) string {
	return constants.PARTY_SELF_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetPartyItemListKey(partyID int64) string {
	return constants.PARTY_ITEM_LIST_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetPartyItemRankListKey(partyID int64) string {
	return constants.PARTY_ITEM_RANK_LIST_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetPartyItemRankListCacheKey(partyID int64) string {
	return constants.PARTY_ITEM_RANK_LIST_CACHE_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetPartyItemSelfKey(partyID int64, username string) string {
	return constants.PARTY_ITEM_SELF_KEY_PREFIX + strconv.FormatInt(partyID, 10) + "_" + username
}

func GetPartyItemSelfKeyByToken(token string) string {
	return constants.PARTY_ITEM_SELF_KEY_PREFIX + token
}

func GetPartyItemFlowerSetKey(partyID int64, username string) string {
	return constants.PARTY_ITEM_FLOWER_SET_PREFIX + strconv.FormatInt(partyID, 10) + "_" + username
}

func GetPartyItemCmtListKey(partyID int64, username string) string {
	return constants.PARTY_ITEM_CMT_LIST_PREFIX + strconv.FormatInt(partyID, 10) + "_" + username
}

func GetPartyFriendJoinCacheKey(username string) string {
	return constants.PARTY_FRIEND_JOIN_CACHE_KEY_PREFIX + username
}

func GetPartyFriendHostCacheKey(username string) string {
	return constants.PARTY_FRIEND_HOST_CACHE_KEY_PREFIX + username
}

// func GetUserNoticeListKey(username string) string {
// 	return constants.NOTICE_SET_PREFIX + username
// }

func GetUserNewFriendNoticeListKey(username string) string {
	return constants.NEW_FRIEND_SET_PREFIX + username
}

// func GetNoticeSelfKey(nid int64) string {
// 	return constants.NOTICE_KEY_PREFIX + strconv.FormatInt(nid, 10)
// }

func GetPartyInviteListKey(partyID int64) string {
	return constants.PARTY_INVITE_LIST_KEY_PREFIX + strconv.FormatInt(partyID, 10)
}

func GetGlobalPartyInviteSetKey() string {
	return constants.PARTY_GLOBAL_INVITE_SET_KEY
}

func GetDesignMsgSelfKey(msgID int64) string {
	return constants.DESIGN_MSG_CNT_KEY_PREFIX + strconv.FormatInt(msgID, 10)
}

func GetGuildApplyListKey(gid int64) string {
	return constants.GUILD_APPLY_LIST_KEY_PREFIX + strconv.FormatInt(gid, 10)
}

func GetGuildNewApplyCntKey(gid int64) string {
	return constants.GUILD_NEW_APPLY_CNT_KEY_PREFIX + strconv.FormatInt(gid, 10)
}

func GetGuildBoardNewestKey(gid int64) string {
	return constants.GUILD_NEWEST_BOARD_ID_KEY_PREFIX + strconv.FormatInt(gid, 10)
}
