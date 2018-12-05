package dispatch_pixie

import (
	"appcfg"
	"constants"
	. "pixie_contract"
	. "pixie_contract/api_specification"
	. "player_mem"
	"service_pixie/paper"
	"service_pixie/paperInfo"
	"time"
	. "zk_manager"
)

func PaperVerifyHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperVerifyReq
	var output PaperVerifyResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		nt := time.Now()
		date := nt.Format("20060102")
		timeStamp := nt.Unix()

		if input.Style1 == "" && input.Style2 == "" {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if input.Score <= 0 || input.Score > 100 || (!paperInfo.TagLegal(input.Tag1) && input.Tag1 != "") || (!paperInfo.TagLegal(input.Tag2) && input.Tag2 != "") || (input.Style1 != "" && !paperInfo.StyleLegal(input.Style1)) || (input.Style2 != "" && !paperInfo.StyleLegal(input.Style2)) {
			resp = PIXIE_ERR_PARAM_ILLEGAL
			goto END
		}

		if num := player.GetDayVerifyNum(date); num >= appcfg.GetInt("verify_day_limit", constants.PIXIE_VERIFY_DAY_LIMIT) {
			resp = PIXIE_ERR_PAPER_VERIFY_LIMIT
		} else {
			if success, path := LockPixiePaper(input.PaperID); success {
				defer Unlock(path)

				if err, paperMod := paper.GetDesignPaperByID(input.PaperID); err != nil {
					resp = PIXIE_ERR_PAPER_VERIFY_GET_WRONG
				} else {
					if paperMod.VerifiedCount < appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET) {
						if err := paper.PaperVerifyInsert(input.PaperID, paperMod.Cname, paperMod.Extra, player.Username, player.Nickname, input.Tag1, input.Tag2, input.Style1, input.Style2, input.Score, timeStamp); err != nil {
							resp = PIXIE_ERR_PAPER_VERIFY_INSERT
						} else {
							player.IncreVerify(date)
							if err := player.SaveStatusToDB(""); err != nil {
								resp = PIXIE_ERR_FLUSH_PLAYER
							} else {
								output.VerifiedCount = num + 1
								output.VerifyLimit = appcfg.GetInt("verify_day_limit", constants.PIXIE_VERIFY_DAY_LIMIT)
								output.VerifiedNum = paperMod.VerifiedCount + 1
								output.VerifyCountLimit = appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET)

								if err, isReport, randPaperForVerify := paper.GetRandomPaperForVerify(player.Username); err != nil {
									if err == constants.PixieNoPaperForVerify {
										resp = PIXIE_ERR_RANDOM_PAPER_NO_DATA
									} else {
										resp = PIXIE_ERR_RANDOM_PAPER_FOR_VERIFY
									}
								} else {
									output.IsReport = isReport
									output.Paper = randPaperForVerify
								}
							}
						}
					} else {
						resp = PIXIE_ERR_PAPER_VERIFY_COUNT_LIMIT
					}
				}
			} else {
				resp = PIXIE_ERR_LOCKPAPER_FAIL
			}
		}
	END:
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetRandomPaperForVerifyHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetVerifyPaperReq
	var output GetVerifyPaperResp

	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		today := time.Now().Format("20060102")

		if err, isReport, paperOwner := paper.GetRandomPaperForVerify(player.Username); err != nil {
			if err == constants.PixieNoPaperForVerify {
				resp = PIXIE_ERR_RANDOM_PAPER_NO_DATA
			} else {
				resp = PIXIE_ERR_RANDOM_PAPER_FOR_VERIFY
			}
		} else {
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			} else {
				output.VerifiedNum = paperOwner.VerifiedCount
				output.VerifyCountLimit = appcfg.GetInt("fake_verify_target", constants.VERIFY_TARGET)
				output.VerifiedCount = player.GetDayVerifyNum(today)
				output.VerifyLimit = appcfg.GetInt("verify_day_limit", constants.PIXIE_VERIFY_DAY_LIMIT)
				output.IsReport = isReport
				output.Paper = paperOwner
			}
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func PaperNotifyQueryHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperNotifyQueryReq
	var output PaperNotifyQueryResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		tn := time.Now()
		today := tn.Format("20060102")
		yes := tn.AddDate(0, 0, -1).Format("20060102")
		if err, yesterdayProfit, allProfit, res := paper.PaperNotifyQuery(player.Username, input.Page, input.PageSize, tn); err != nil {
			resp = PIXIE_ERR_PAPER_VERIFY_QUERY
		} else {
			output.Res = res
			output.BenefitCount = allProfit
			output.VerifyCount = player.GetVerifyCount()
			output.YesBenefit = yesterdayProfit
			output.YesVerify = player.GetDayVerifyNum(yes)

			output.VerifiedCount = player.GetDayVerifyNum(today)
			output.VerifyLimit = appcfg.GetInt("verify_day_limit", constants.PIXIE_VERIFY_DAY_LIMIT)
		}

		result, parseErr = PixieMarshal(output)
	}
	return
}

func PaperVerifiedDetailQueryHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperVerifiedDetailReq
	var output PaperVerifiedDetailResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		checkDate := time.Unix(input.CheckTime, 0).AddDate(0, 0, -1).Format("20060102")
		if err, list := paper.PaperVerifiedDetailQuery(player.Username, checkDate); err != nil {
			resp = PIXIE_ERR_PAPER_GET_VERIFIED_DETAIL_WRONG
		} else if err = paper.UpdatePaperNotifyRewarded(input.NotifyID); err != nil {
			resp = PIXIE_ERR_UPDATE_PAPER_NOTIFY_STATUS
		} else {
			output.Res = list
			output.RewardCount = constants.VERIFY_REWARD_COUNT
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func PaperReportCopyHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperReportCopyReq
	var output PaperReportCopyResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err := paper.PaperReportCopy(player.Username, input.PaperID, input.Reason, input.Pic, input.Contact); err != nil {
			resp = PIXIE_ERR_REPORT_WRONG
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func PaperCopySupportHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperCopySupportReq
	var output PaperCopySupportResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err := paper.PaperCopySupprot(input.CopyID, input.PaperID); err != nil {
			resp = PIXIE_PAPER_COPY_SUPPORT_WRONG
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func PaperCopyQueryHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input PaperCopyQueryReq
	var output PaperCopyQueryResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err, list := paper.PaperCopyQuery(input.PaperID); err != nil {
			resp = PIXIE_ERR_COPY_QUERY_WRONG
		} else {
			output.List = list
		}
		result, parseErr = PixieMarshal(output)
	}
	return
}

func GetPaperCopyRewardHandle(player *GFPlayer, data string) (result string, resp PixieRespInfo, parseErr error) {
	var input GetPaperCopyRewardReq
	var output GetPaperCopyRewardResp
	if parseErr = PixieUnmarshal(data, &input); parseErr == nil {
		if err, currencyType, currencyVal := paper.GetPaperNotifyReward(input.NotifyID); err != nil {
			resp = PIXIE_ERR_GET_PAPER_VERIFY_REWARD_WRONG
		} else if err = paper.UpdatePaperNotifyRewarded(input.NotifyID); err != nil {
			resp = PIXIE_ERR_UPDATE_PAPER_NOTIFY_STATUS
		} else {
			if currencyType == constants.PIXIE_GOLD_TYPE {
				player.AddGold(int64(currencyVal))
			} else {
				player.AddPxc(currencyVal)
			}
			if err := player.SaveStatusToDB(""); err != nil {
				resp = PIXIE_ERR_FLUSH_PLAYER
			}
		}

		output.CurrentGold = player.GetMoney()
		output.CurrentPxc = player.GetPxc()
		result, parseErr = PixieMarshal(output)
	}
	return
}
