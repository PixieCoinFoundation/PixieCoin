package mails

import (
	"appcfg"
	"constants"
	"encoding/json"
	"sync"
	"time"
)

import (
	"db/mail"
	. "logger"
	. "pixie_contract/api_specification"
	"tools"
	. "types"
)

//send to all mails
var cronMails []CronMail
var cronMailsMutex sync.Mutex

func init() {
	if appcfg.GetServerType() == "" {
		go refreshCronMail()
	}
}

func CheckCronMail(username string, regTime int64, ccms []int64, channel string) (changed bool, res []int64) {
	res = make([]int64, 0)
	now := time.Now().Unix()

	cronMailsMutex.Lock()
	defer cronMailsMutex.Unlock()
	for _, cm := range cronMails {
		if cm.ChannelLimit && cm.ThirdChannel != channel {
			continue
		}

		if cm.CronTime > now {
			continue
		} else if cm.MinRegTime > 0 && regTime < cm.MinRegTime {
			continue
		} else if cm.MaxRegTime > 0 && regTime > cm.MaxRegTime {
			continue
		}

		got := false
		for _, rid := range ccms {
			if cm.Id == rid {
				got = true
				break
			}
		}
		if !got {
			var clo []ClothesInfo
			if err := json.Unmarshal([]byte(cm.Clothes), &clo); err != nil {
				Err(err)
				continue
			} else if err := SendToOneC("", username, cm.Title, cm.Content, cm.Diamond, cm.Gold, clo, true); err == nil {
				changed = true
				res = append(res, cm.Id)
			}
		} else {
			res = append(res, cm.Id)
		}
	}

	return
}

func refreshCronMail() {
	for {
		c := mail.GetNotSendCronMails()
		cronMailsMutex.Lock()
		cronMails = c
		cronMailsMutex.Unlock()

		time.Sleep(60 * time.Second)
	}
}

// func GetMails(username string, pageID int, pageCount int) ([]*Mail, error) {
// 	return mail.Find(username, pageID, pageCount)
// }

func GetOneMail(username string, mailID int64) (Mail, error) {
	return mail.FindOne(username, mailID)
}

func DeleteMail(username string, mailID int64) error {
	return mail.DeleteMail(username, mailID)
}

//SetMailRead 将(没有礼物的)邮件设置为已读
func SetMailRead(username string, mailID int64) error {
	return mail.SetRead(username, mailID)
}

// func GetUnreadMailCount(username string) (int, error) {
// 	return mail.GetUnreadCount(username)
// }

func SendToOne(from string, to string, title string, content string, diamond int64, gold int64, clothesID string, shouldDelete bool, expireTime int64) (err error) {
	clothes := make([]ClothesInfo, 0)

	if tools.ClothesIDLegal(clothesID) {
		c := ClothesInfo{
			ClothesID: clothesID,
			Count:     1,
		}
		clothes = append(clothes, c)
	}

	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Diamond:    diamond,
		Gold:       gold,
		Clothes:    clothes,
		Delete:     shouldDelete,
		Time:       time.Now().Unix(),
		ExpireTime: expireTime,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneS(from string, to string, title string, content string, expireTime int64, mt int) (err error) {
	clothes := make([]ClothesInfo, 0)

	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     true,
		Time:       time.Now().Unix(),
		Type:       mt,
		ExpireTime: expireTime,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneSS(from string, to string, title string, content string) (err error) {
	clothes := make([]ClothesInfo, 0)

	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     true,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneSSS(from string, to string, title string, content string, delete bool, et int64) (err error) {
	clothes := make([]ClothesInfo, 0)

	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     delete,
		Time:       now,
		ExpireTime: et,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneSST(from string, to string, title string, content string, mt int) (err error) {
	clothes := make([]ClothesInfo, 0)

	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     true,
		Time:       now,
		Type:       mt,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneT(from string, to string, title string, content string, diamond int64, gold int64, clothesID string, shouldDelete bool, mt int) (err error) {
	clothes := make([]ClothesInfo, 0)

	if tools.ClothesIDLegal(clothesID) {
		c := ClothesInfo{
			ClothesID: clothesID,
			Count:     1,
		}
		clothes = append(clothes, c)
	}

	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Diamond:    diamond,
		Gold:       gold,
		Clothes:    clothes,
		Delete:     shouldDelete,
		Time:       now,
		Type:       mt,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneC(from string, to string, title string, content string, diamond int64, gold int64, clo []ClothesInfo, shouldDelete bool) (err error) {
	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Diamond:    diamond,
		Gold:       gold,
		Clothes:    clo,
		Delete:     shouldDelete,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneCT(from string, to string, title string, content string, diamond int64, gold int64, clo []ClothesInfo, shouldDelete bool, et int64) (err error) {
	now := time.Now().Unix()
	item := Mail{
		From:       from,
		To:         to,
		Title:      title,
		Content:    content,
		Diamond:    diamond,
		Gold:       gold,
		Clothes:    clo,
		Delete:     shouldDelete,
		Time:       now,
		ExpireTime: et,
	}
	return mail.SendToOne(item, nil)
}

func SendToOneClos(to, title, content string, clos []string) error {
	clothes := make([]ClothesInfo, 0)

	for _, c := range clos {
		if tools.ClothesIDLegal(c) {
			c := ClothesInfo{
				ClothesID: c,
				Count:     1,
			}
			clothes = append(clothes, c)
		}
	}

	now := time.Now().Unix()
	item := Mail{
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     true,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendADReturnMail(to, title, content string, diamond, gold int64, logDiamond int) error {
	clothes := make([]ClothesInfo, 0)
	mc := MailContent{
		Content:      content,
		ADLogDiamond: logDiamond,
	}
	mcb, _ := json.Marshal(mc)

	now := time.Now().Unix()
	item := Mail{
		Type:       constants.MAIL_TYPE_AD_RETURN,
		To:         to,
		Title:      title,
		Content:    string(mcb),
		Clothes:    clothes,
		Delete:     true,
		Diamond:    diamond,
		Gold:       gold,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendPXCClothes(to, title, content string, clothesID string, clothesCount int) error {
	clothes := []ClothesInfo{ClothesInfo{ClothesID: clothesID, Count: clothesCount}}

	now := time.Now().Unix()
	item := Mail{
		To:         to,
		Title:      title,
		Content:    content,
		Clothes:    clothes,
		Delete:     true,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}

func SendPXCBuyFailMail(to, title, content string, failClos []SimpleClothes1) error {
	mc := MailContent{
		Content:         content,
		FailClothesList: failClos,
	}
	mcb, _ := json.Marshal(mc)

	now := time.Now().Unix()
	item := Mail{
		Type:       constants.MAIL_TYPE_PXC_BUY_FAIL,
		To:         to,
		Title:      title,
		Content:    string(mcb),
		Clothes:    make([]ClothesInfo, 0),
		Delete:     true,
		Time:       now,
		ExpireTime: now + constants.DEFAULT_MAIL_EXPIRE_TIME,
	}
	return mail.SendToOne(item, nil)
}
