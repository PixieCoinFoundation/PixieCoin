package guild

import (
	"common"
	"constants"
	"math/rand"
	"time"
)

func GetGuildWarMedalRankReward(rank int) int {
	if rank == 1 {
		return constants.GUILD_WAR_RANK_1_REWARD_DMD
	} else if rank == 2 {
		return constants.GUILD_WAR_RANK_2_REWARD_DMD
	} else if rank == 3 {
		return constants.GUILD_WAR_RANK_3_REWARD_DMD
	} else if rank <= 10 {
		return constants.GUILD_WAR_RANK_10_REWARD_DMD
	} else if rank <= 50 {
		return constants.GUILD_WAR_RANK_50_REWARD_DMD
	} else if rank <= 100 {
		return constants.GUILD_WAR_RANK_100_REWARD_DMD
	}
	return 0
}

func GuildWarSeasonLegalDay(t time.Time) bool {
	if t.Weekday() != time.Sunday {
		return true
	}

	//current gwt
	_, gwt := common.GetGuildYWTokenByTime(t, 0, 0)

	//yester gwt
	// yt := t.AddDate(0,0,-1)
	// _,ygwt := common.GetGuildYWTokenByTime(yt,0,0)

	//tomorrow gwt
	tt := t.AddDate(0, 0, 1)
	_, tgwt := common.GetGuildYWTokenByTime(tt, 0, 0)

	if gwt != tgwt {
		return false
	}

	return true
}

func ShouldSettle(t time.Time) bool {
	if t.Weekday() != time.Monday || t.Hour() >= 10 || t.Hour() <= 1 {
		return false
	}

	//current gwt
	_, gwt := common.GetGuildYWTokenByTime(t, 0, 0)

	//yesterday gwt
	yt := t.AddDate(0, 0, -1)
	_, ygwt := common.GetGuildYWTokenByTime(yt, 0, 0)

	if gwt != ygwt {
		return true
	}

	return false
}

func WarDateLegal(wd string, nwd string, ywd string, nh int) bool {
	if wd == nwd {
		return true
	}

	if wd != nwd && wd != ywd {
		return false
	}

	if wd == ywd && nh >= constants.GUILD_WAR_CLAIM_START_HOUR {
		return false
	}

	return true
}

func CanRewardShare(index int, shareCnt int) bool {
	switch index {
	case 0:
		if shareCnt >= 20 {
			return true
		}
	case 1:
		if shareCnt >= 30 {
			return true
		}
	case 2:
		if shareCnt >= 50 {
			return true
		}
	case 3:
		if shareCnt >= 80 {
			return true
		}
	default:
		return false
	}
	return false
}

func GetShareClothesReward(index int) (gold int, diamond int) {
	v := rand.Intn(100)

	switch index {
	case 0:
		if v <= 29 {
			gold = 500
		} else if v <= 69 {
			gold = 1000
		} else {
			gold = 1500
		}
		return
	case 1:
		if v <= 29 {
			gold = 1000
		} else if v <= 69 {
			gold = 2000
		} else {
			gold = 3000
		}
		return
	case 2:
		if v <= 29 {
			diamond = 10
		} else if v <= 69 {
			diamond = 20
		} else {
			diamond = 30
		}
		return
	case 3:
		if v <= 29 {
			diamond = 15
		} else if v <= 69 {
			diamond = 30
		} else {
			diamond = 45
		}
		return
	default:
		return
	}
	return
}
