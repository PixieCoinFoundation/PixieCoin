package types

//DesignerZoneBanner 设计师社区配置
type DesignerZoneBanner struct {
	ID           int    `json:"id"`           //数据库ID
	BannerTitle  string `json:"bannerTitle"`  //banner 标题
	BannerType   string `json:"bannerType"`   //banner 类型
	URL          string `json:"url"`          //跳转链接
	IndexOrRecom string `json:"indexOrRecom"` //主页banner或者推荐页banner
	Shunxu       int    `json:"shunxu"`       //顺序
	IsStick      bool   `json:"istick"`       //是否置顶
	CreateTime   string `json:"createTime"`   // 创建时间
	StartTime    string `json:"startTime"`    //开始时间
	EndTime      string `json:"endTime"`      //结束时间
	IMGURL       string `json:"img_url"`      //图片地址
}
type Topic struct {
	ID           int    `json:"id"`            //数据库ID
	BannerTitle  string `json:"bannerTitle"`   //banner 标题
	OwerTitle    string `json:"ower_title"`    //对应的属性名称
	URL          string `json:"url"`           //跳转链接
	ClothesCount int    `json:"clothes_count"` //衣服数量
	Shunxu       int    `json:"shunxu"`        //顺序
	IsStick      bool   `json:"istick"`        //是否置顶
	CreateTime   string `json:"createTime"`    // 创建时间
	StartTime    string `json:"startTime"`     //开始时间
	EndTime      string `json:"endTime"`       //结束时间
	TitleIcon    string `json:"title_icon"`    //图片地址
	MainIcon     string `json:"main_icon"`     //图片地址1
	SliderIcon   string `json:"slider_icon"`
}

type BannerInfo struct {
	BannerType string // 1.html 2.designerShop 3.topic
	URL        string
	Order      int
	Params     string
}
type IndexBannerInfo struct {
	BanenrID int
	Order    int
	Config   string
	IconPath string
}
type TopicInfo struct {
	TopicID    int
	Title      string
	TitleIcon  string
	MainIcon   string
	SliderIcon string
}
type ClothesPoolInfo struct {
	ID        int    `json:"id"`
	ClothesID string `json:"clothesID"` //衣服ID
	Cname     string `json:"cname"`     //衣服名称
	UserName  string `json:"username"`  //用户ID
	NickName  string `json:"nickname"`  //用户昵称
	Model     int    `json:"model"`     //角色
	Part      int    `json:"part"`      //衣服类型
	PriceType int    `json:"priceType"` //价格类型
	Price     int    `json:"price"`     //价格
	Star      int    `json:"star"`      //星级
	SaleCount int    `json:"saleCount"` //销售数量
	Status    string `json:"status"`    //库存，下架状态
	Tag1      int    `json:"tag1"`      //标签
	Tag2      int    `json:"tag2"`      //标签

	AdminID    string `json:"adminID"`    //管理员ID
	EntryTime  string `json:"entryTime"`  //进入预选池时间
	ImgURL     string `josn:"imgURL"`     //图片URL
	PoolStatus int    `json:"poolstatus"` //衣服的状态，1代表预选池，2代表推荐池
	Topic      int    `json:"topic"`      //专题
	IsStick    bool   `json:"is_stick"`   //是否置顶
}
type DesignerPoolInfo struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`           //用户名
	Nickname         string `json:"nickname"`           //用户昵称
	FirstClothesTime int64  `json:"first_clothes_time"` //衣服第一次上架时间
	InSaleCount      int    `json:"in_sale_count"`      //正在售卖的衣服总量
	SaleGold         int    `json:"sale_gold"`          //金币销售额
	SaleDiam         int    `json:"sale_diam"`          //钻石销售额
	AdminID          string `json:"adminID"`            //管理员ID
	EntryTime        string `json:"entry_time"`         //进入预选池的时间
	PoolStatus       int    `json:"poolstatus"`         //衣服的状态，1代表预选池，2代表推荐池
	IsStick          bool   `json:"is_stick"`           //是否置顶
}
