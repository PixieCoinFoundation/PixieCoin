package types

import ()

type TimerCode int32

const (
	TIMER_DEBUG                  TimerCode = -3
	TIMER_UPDATE_SERVER_STATUS   TimerCode = -2
	TIMER_GC                     TimerCode = -1
	RESEARVED                    TimerCode = 0
	TIMER_PLAYER_TIMEOUT         TimerCode = 1
	TIMER_REFRESH_DB_SESSION     TimerCode = 2
	TIMER_RECONNECT_LOGIN_SERVER TimerCode = 3
	TIMER_CHECK_COSPLAY          TimerCode = 4
	TIMER_PROCESS_COSPLAY        TimerCode = 5
	TIMER_CHECK_COS_ITEMS        TimerCode = 6
	TIMER_REFRESH_COS_ITEM_RANK  TimerCode = 7
	TIMER_REMOVE_TOP_COS_ITEM    TimerCode = 8
	TIMER_REFRESH_TOP_COS_ITEM   TimerCode = 9

	TIMER_REFRESH_DESIGNER_TOP_AMOUNT  TimerCode = 10
	TIMER_REFRESH_DESIGNER_TOP_SALE    TimerCode = 11
	TIMER_REFRESH_DESIGNER_TOP_GOLD    TimerCode = 12
	TIMER_REFRESH_DESIGNER_TOP_DIAMOND TimerCode = 13

	TIMER_SERVER_HEARTBEAT TimerCode = 14

	TIMER_SEND_MAIL_TO_ALL TimerCode = 15

	TIMER_REFRESH_COSPLAY        TimerCode = 16
	TIMER_REFRESH_FINISH_COSPLAY TimerCode = 17

	// help
	TIMER_REFRESH_HELP TimerCode = 18
	// pk
	TIMER_REFRESH_PK_CACHE TimerCode = 19
	// cos comment
	TIMER_REFRESH_COS_COMMENT TimerCode = 20
	// player timer work
	TIMER_PLAYER_TIME_WORK TimerCode = 1

	// 保留项
	TIMER_RESERVED_0_PLAYER_CLOSE TimerCode = -1000 // gfplayer关闭timer
	TIMER_RESERVED_0_AGENT_CLOSE  TimerCode = -1001 // agent主关闭timer
)
