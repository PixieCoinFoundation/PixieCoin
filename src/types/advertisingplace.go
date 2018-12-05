package types

//AdvertisingPlace 广告位
type AdvertisingPlace struct {
	ID           string   //广告位ID
	Name         string   //广告位模特和部位
	Heat         float64  //广告位热度
	Money        int      //广告位原价
	Diamond      int      //钻石价值
	PayType      string   //付费类型
	AddPrice     int      //广告位加价
	NowPrice     int      //广告位现价
	Frequency    []string //频率
	LastAddPrice string   //最后一次加价
	FansAddPrice []string //粉丝加价
	StartTime    int64    //广告位开始时间
	EndTime      int64    // 广告位结束时间
}
