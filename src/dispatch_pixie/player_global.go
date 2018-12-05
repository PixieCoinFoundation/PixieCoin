package dispatch_pixie

import (
	"appcfg"
	"constants"
	"db/board"
	"db_pixie/land"
	"db_pixie/status"
	. "event_dispatch"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service/guild"
	"service/params"
	"service_pixie/level"
	"time"
)

func LoginHandle(player *GFPlayer, data string, remoteAddr string) (result string, resp PixieRespInfo, parseErr error) {
	var input PixieLoginReq
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if input.Username == "" || input.UID <= 0 || input.AccessToken == "" {
			resp = PIXIE_ERR_LOGIN_PARAM_EMPTY
		} else {
			generatedPlayer, err := GeneratePlayerByToken(input.Username, input.DeviceToken, input.UID, input.Platform, input.DeviceId, "", input.ThirdUsername, input.AccessToken, input.ThirdChannel, input.Channel, input.KorRevertUnregister)
			if err == nil && generatedPlayer != nil {
				NotifyPlayerLogin(&input, remoteAddr, generatedPlayer)
			} else if err.Error() == constants.PLAYER_LIMITED_MSG {
				resp = PIXIE_ERR_CHANNEL_LIMIT
			} else if err.Error() == constants.PLAYER_BAN_MSG {
				resp = PIXIE_ERR_PLAYER_BAN
			} else if err.Error() == constants.PLAYER_DELETING_MSG {
				resp = PIXIE_ERR_PLAYER_DELETING
			} else if err.Error() == constants.PLAYER_DELED_MSG {
				resp = PIXIE_ERR_PLAYER_DELETED
			} else {
				resp = PIXIE_ERR_PLAYER_LOGIN_ERR
			}
		}
	}

	return
}

func GetAllUserDataHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var output GetAllUserDataResp
	now := time.Now()

	output.UserAttributes = player.GetStatus()

	ls, _ := land.ListPlayerLand(player.Username)
	if len(ls) <= 0 {
		status.InitPlayerLand(player.Username, player.Nickname, player.GetHead(), player.GetSex())
		output.Lands, _ = land.ListPlayerLand(player.Username)
	} else {
		output.Lands = ls
	}

	if appcfg.GetBool("check_rt_party_time", true) {
		output.RTPartyTime = []string{"1200-1400", "1900-2300"}
	}

	output.PrizePartyClose = params.PrizePartyClose()
	output.CurrentTime = now.Unix()
	_, output.UTCOffset = now.Zone()
	// output.MCReward = params.GetMonthCardReward(now.Format("200601"))
	output.ShopCart = player.ShopCart
	// output.GotedClothes = player.GetClothes()
	output.MailCount = player.GetUnreadMailCount()
	// output.Events = event.GetOpenEvents(player.ThirdChannel)
	output.BoardNew = board.GetBoardNew(player.Username)
	// output.TheaterID, output.TheaterStartTime, output.TheaterCloseTime, _, output.TheaterClo = params.GetFirstTheaterInfo()
	// output.Theaters = params.GetTheaters()

	output.SpecialAccount = params.IsInWhiteList(player.Username)
	output.OrderOpen = params.OrderOpen()
	output.GiftPackOpen = params.GiftPackOpen()
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		if player.ThirdChannel == constants.THIRD_CHANNEL_KOR_GOOGLE_PLAY || player.ThirdChannel == constants.THIRD_CHANNEL_KOR_ONE_STORE {
			output.GiftPackOpen = true
		}
	} else {
		if (player.ThirdChannel == "" && player.OldChannel == "bili_android") || (player.ThirdChannel != "0" && player.ThirdChannel != "") {
			output.GiftPackOpen = true
		}
	}
	output.IsContractDesigner = params.InDesignCornList(player.Username)
	output.MonthCardCanBuy = params.MonthCardCanBuy()
	output.LevelMap = level.GetLevelExpMap()
	output.ChannelShareURL = params.GetChannelShareUrl(player.ThirdChannel)

	if player.GuildInfo.GID > 0 {
		if fid := guild.GetGuildBoardNewestIDR(player.GuildInfo.GID); fid > 0 {
			if player.GetGuildBoardReadID() < fid {
				output.GuildBoardNewMsg = true
			}
		}
	}

	output.ShouldUploadLog = params.PlayerShouldUploadLog(player.Username)

	result, parseErr = PixieMarshal(output)
	return
}

func HeartBeatHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var output HeartBeatResp
	var input CommonReq

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		now := time.Now()
		output.Time = now.Unix()
		_, output.UTCOffset = now.Zone()
		player.RefreshHeartbeatTime()

		changeStreet := player.RefreshStreet(false)
		changeTili := player.RecoverTili(input.DeviceId)

		if changeStreet || changeTili {
			player.SaveStatusToDB("")
		}

		output.Tili = player.GetTili()
		output.MaxTili = player.GetMaxTili()
		output.TiliRTime = player.GetTiliRTime()

		output.StreetDetail = player.GetStreetDetail()

		result, parseErr = PixieMarshal(output)
	}

	return
}
