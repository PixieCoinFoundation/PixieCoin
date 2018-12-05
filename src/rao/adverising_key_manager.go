package rao

import (
	"fmt"
	"time"
)

//GetDesignerDayKey 获取设计师 日榜的key
func GetDesignerDayKey(designerUsername string, now time.Time) string {
	return fmt.Sprintf("%s-day-%s", now.Format("20060102"), designerUsername)
}

//GetDesignerYearKey 获取 年榜的key
func GetDesignerYearKey(cid string, now time.Time) string {
	return fmt.Sprintf("%s-year-%s", now.Format("2006"), cid)
}

//GetDesignerMonthKey 月榜的key
func GetDesignerMonthKey(cid string, now time.Time) string {
	return fmt.Sprintf("%s-month-%s", now.Format("200601"), cid)
}
func DayDesignerListKey(now time.Time) string {
	return fmt.Sprintf("%s-dayDesignerList", now.Format("20060102"))
}
func MonthDesignerListKey(now time.Time) string {
	return fmt.Sprintf("%s-MonthDesignerList", now.Format("200601"))
}

//GetDesignerDayJoinKey  设计师竞每天参加的竞拍
func GetDesignerDayJoinKey(username string, now time.Time) string {
	return fmt.Sprintf("%s-DesignerDayJoin-%s", now.Format("20060102"), username)
}

//GetJoinAuctionKey 设计师参加的竞拍的key
func GetJoinAuctionKey(prefix string, Model, Part, AuctionPos int) string {
	return fmt.Sprintf("%s:%d:%d=%d", prefix, Model, Part, AuctionPos)
}

//GetGoldKey 获取设计师金币的Key
func GetGoldKey(prefix string, Model, Part, AuctionPos int) string {
	return fmt.Sprintf("gold:%s:%d:%d=%d", prefix, Model, Part, AuctionPos)
}

//获取世纪时钻石的key
func GetDiamKey(prefix string, Model, Part, AuctionPos int) string {
	return fmt.Sprintf("diam:%s:%d:%d=%d", prefix, Model, Part, AuctionPos)

}
