package gift

import (
	"appcfg"
	"constants"
	"db/gift"
	"sync"
	"time"
	. "types"
)

var gifts map[int64]*GiftPackContent
var gns map[int64]string

var giftLock sync.Mutex

func init() {
	if appcfg.GetServerType() == "" {
		var err error
		if gifts, gns, err = gift.GetGictPacks(); err != nil {
			panic(err)
		}

		go refreshGifts()
	}
}

func UseGiftCode(code string, typee int, packID int64, username string) bool {
	return gift.UseCode(code, typee, packID, username)
}

func GetPackByCode(code string) (int64, string, int) {
	giftLock.Lock()
	defer giftLock.Unlock()
	for id, c := range gifts {
		if c.Type == constants.ONE_CODE_TYPE && c.Code == code {
			return id, gns[id], constants.ONE_CODE_TYPE
		}
	}

	id, typee := gift.GetGiftPackIDByCode(code)
	name := gns[id]
	return id, name, typee
}

func GetGiftPackRewards(packID int64) (ok bool, g int64, d int64, cs []ClothesInfo, i1 int, i2 int, i3 int, i4 int, i5 int, i6 int, tili int, thirdChannel string) {
	giftLock.Lock()
	defer giftLock.Unlock()

	var gpc *GiftPackContent
	if gpc, ok = gifts[packID]; ok {
		g = int64(gpc.Gold)
		d = int64(gpc.Diamond)
		cs = gpc.Clothes
		i1 = gpc.Item101
		i2 = gpc.Item102
		i3 = gpc.Item103
		i4 = gpc.Item104
		i5 = gpc.Item105
		i6 = gpc.Item106
		tili = gpc.Tili
		thirdChannel = gpc.ThirdChannel
		return
	}

	return
}

func refreshGifts() {
	for {
		time.Sleep(60 * time.Second)

		if gfs, gfns, err := gift.GetGictPacks(); err == nil {
			giftLock.Lock()
			gifts = gfs
			gns = gfns
			giftLock.Unlock()
		}
	}
}
