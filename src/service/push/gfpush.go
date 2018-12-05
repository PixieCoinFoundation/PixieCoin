package push

import (
	"constants"
	"dao_player"
	. "pixie_contract/api_specification"
	"push"
	"time"
)

// func PushToPlayer(toUsername string, content string, title string) {
// 	push.TPushToOne(toUsername, title, content)
// }

func PushToAll(content string, title string) {
	push.TPushToAll(title, content)
}

func pushToPlayerTC(toUsername string, title string, content string) {
	push.TPushToOne(toUsername, title, content)
}

func PushKorTC(toUsername string, title string, content string) {
	iph, _, _, _, _, _, _, _, inp := getPlayerPushConfig(toUsername)

	if iph {
		return
	}

	nh := time.Now().Hour()
	if inp && (nh <= constants.KOR_PUSH_DAY_START || nh >= constants.KOR_PUSH_DAY_END) {
		return
	}

	push.TPushToOne(toUsername, title, content)
}

func PushShopCart(username, title, content string, sync bool) {
	if sync {
		pushShopCart1(username, title, content)
	} else {
		go pushShopCart1(username, title, content)
	}
}

func PushGuildWar(username, title, content string, sync bool) {
	if sync {
		pushGuildWar1(username, title, content)
	} else {
		go pushGuildWar1(username, title, content)
	}
}

func PushGuild(username, title, content string, sync bool) {
	if sync {
		pushGuild1(username, title, content)
	} else {
		go pushGuild1(username, title, content)
	}
}

func PushFriendReq(username, title, content string, sync bool) {
	if sync {
		pushFriendReq1(username, title, content)
	} else {
		go pushFriendReq1(username, title, content)
	}
}

func PushHopeDone(username, title, content string, sync bool) {
	if sync {
		pushHopeDone1(username, title, content)
	} else {
		go pushHopeDone1(username, title, content)
	}
}

func PushPartyInvite(username, title, content string, sync bool) {
	if sync {
		pushPartyInvite1(username, title, content)
	} else {
		go pushPartyInvite1(username, title, content)
	}
}

func PushNewCustom(username, title, content string, sync bool) {
	if sync {
		pushNewCustom1(username, title, content)
	} else {
		go pushNewCustom1(username, title, content)
	}
}

func getPlayerPushConfig(username string) (ignorePush, ignoreGuildWarPush, ignoreGuildPush, acceptFriendRequestPush, ignoreHopeDonePush, acceptPartyInvitePush, ignoreNewCustomPush, ignoreShopCartPush, ignoreNightPush bool) {
	e1 := Extra1{}

	if s, err := dao_player.GetPlayer(username); err == nil {
		e1 = s.ExtraInfo1
	}

	return e1.IgnorePush, e1.IgnoreGuildWarPush, e1.IgnoreGuildPush, e1.AcceptFriendRequestPush, e1.IgnoreHopeDonePush, e1.AcceptPartyInvitePush, e1.IgnoreNewCustomPush, e1.IgnoreCartPush, e1.IgnoreNightPush
}

func pushShopCart1(username, title, content string) {
	iph, _, _, _, _, _, _, iscp, _ := getPlayerPushConfig(username)

	if !iph && !iscp {
		pushToPlayerTC(username, title, content)
	}
}

func pushGuildWar1(username, title, content string) {
	iph, igwph, _, _, _, _, _, _, _ := getPlayerPushConfig(username)

	if !iph && !igwph {
		pushToPlayerTC(username, title, content)
	}
}

func pushGuild1(username, title, content string) {
	iph, _, igph, _, _, _, _, _, _ := getPlayerPushConfig(username)

	if !iph && !igph {
		pushToPlayerTC(username, title, content)
	}
}

func pushFriendReq1(username, title, content string) {
	iph, _, _, afph, _, _, _, _, _ := getPlayerPushConfig(username)

	if !iph && afph {
		pushToPlayerTC(username, title, content)
	}
}

func pushHopeDone1(username, title, content string) {
	iph, _, _, _, ihdph, _, _, _, _ := getPlayerPushConfig(username)

	if !iph && !ihdph {
		pushToPlayerTC(username, title, content)
	}
}

func pushPartyInvite1(username, title, content string) {
	iph, _, _, _, _, apph, _, _, _ := getPlayerPushConfig(username)

	if !iph && apph {
		pushToPlayerTC(username, title, content)
	}
}

func pushNewCustom1(username, title, content string) {
	iph, _, _, _, _, _, incph, _, _ := getPlayerPushConfig(username)

	if !iph && !incph {
		pushToPlayerTC(username, title, content)
	}
}
