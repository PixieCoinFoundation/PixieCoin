package pushlist

import (
	"appcfg"
	"constants"
	"db/pushlist"
	. "logger"
	"service/push"
	"time"
	. "types"
)

func init() {
	if appcfg.GetServerType() == "" && appcfg.GetLanguage() == "" {
		//send single pushes
		go sendPushLoop()
	}
}

func AddBatchSinglePushJob(jobs []*SinglePushJob) {
	if appcfg.GetLanguage() == "" {
		pushlist.AddBatchSinglePushJob(jobs)
	}
}

func AddSinglePushJob(to, title, content string, ty int) {
	if appcfg.GetLanguage() == "" {
		job := &SinglePushJob{
			To:      to,
			Title:   title,
			Content: content,
			Type:    ty,
		}
		pushlist.AddSinglePushJob(job)
	}
}

func sendPushLoop() {
	for {
		job, err := pushlist.PopSinglePushJob()
		for err == nil && job != nil {
			switch job.Type {
			case constants.PUSH_TYPE_GUILD_WAR:
				push.PushGuildWar(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_GUILD:
				push.PushGuild(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_FRIEND_REQUEST:
				push.PushFriendReq(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_HOPE_DONE:
				push.PushHopeDone(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_PARTY_INVITE:
				push.PushPartyInvite(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_CUSTOM_NEW:
				push.PushNewCustom(job.To, job.Title, job.Content, true)
			case constants.PUSH_TYPE_SHOP_CART:
				push.PushShopCart(job.To, job.Title, job.Content, true)
			default:
				Info("unknown push type:", job.Type)
			}

			job, err = pushlist.PopSinglePushJob()
		}

		time.Sleep(5 * time.Second)
	}
}
