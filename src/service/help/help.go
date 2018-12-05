package help

import (
	"common"
	"sync"
	"time"
)

import (
	"db/help"
	. "types"
)

var (
	inited bool // 是否初始化过吐槽

	helps     []*Help      // 缓存help，key:页数，value:当前页的所有help
	helpsLock sync.RWMutex // lock

	helpComments map[string]([]*HelpComment) // 缓存helpComment，key:helpID，value:所有的对应help的comment
	mapLock      sync.RWMutex                // 修改comments时的lock
)

func init() {
	inited = false

	helps = make([]*Help, 0)
	helpComments = make(map[string]([]*HelpComment))
}

func GetHelpCmtCnt(helpID int64) int {
	return help.GetHelpCmtCnt(helpID)
}

func AddHelp(username string, nickname string, levelID string, boyClothes string, girlClothes string, desc string, previewImage string) (item Help, err error) {
	item = common.GetHelpInRedis(username, nickname, levelID, desc, time.Now().Unix(), previewImage, boyClothes, girlClothes)
	err = help.AddHelp(&item, nil)

	return
}

func GetHelp(page int, pageCount int) (h []*Help, err error) {
	return help.GetRecentHelp(page, pageCount)
}

func AddHelpComment(username string, nickname string, helpID int64, content string) (item HelpComment, err error) {
	item = HelpComment{
		HelpID:    helpID,
		Nickname:  nickname,
		Content:   content,
		Username:  username,
		ReplyTime: time.Now().Unix(),
	}

	if err = help.AddComment(item); err != nil {
		return
	}

	return
}

func GetHelpComment(helpID int64, page int, pageCount int) ([]*HelpComment, error) {
	return help.GetHelpComment(helpID, page, pageCount)
}
