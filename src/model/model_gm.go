package model

import (
	. "pixie_contract/api_specification"
)

type GMOfficialPaper struct {
	ID          int64 `xorm:"pk autoincr id "`
	Cname       string
	Desc        string
	ClothesType string //模特
	PartType    string //部位
	PaperExtra  string // 衣服信息
	PaperStyle  string //风格
	PaperFile   string

	PriceType int //货币类型
	Price     int64
	Star      int // 星级

	Tag1        string
	Tag2        string
	UnlockLevel string //解锁关卡
	HiddenTag   string // 暗标签
	SaleTime    int64  //售卖时间
	AdminName   string //上传名称
	UploadTime  int64  //上传时间
}

type GMRecommendClothes struct {
	PaperID        int64 `xorm:"id"`
	AuthorUsername string
	AuthorNickname string
	AuthorHead     string
	AuthorSex      int
	OwnerUsername  string
	OwnerNickname  string
	OwnerHead      string
	OwnerSex       int

	ClothesType string //衣服模特
	PartType    string //衣服类别
	Cname       string //图纸名称
	Desc        string //描述
	File        string //图纸源文件 PaperFile json
	Extra       string //图纸文件处理后 PaperExtra json
	Status      int    //图纸状态
	Star        int    //图纸星级

	//普通标签
	Tag1 string
	Tag2 string

	//风格 内容是PaperStyle的json字符串
	Style string

	//暗标签
	STag      string
	ClothesID string //服装ID
	Price     int64
	PriceType int
}

type GMOfficialPaperTrans struct {
	ID          int64
	Cname       string
	Desc        string
	ClothesType string //模特
	PartType    string //部位
	PaperExtra  string // 衣服信息
	PaperStyle  string //风格

	PriceType string //货币类型
	Price     string
	Star      string // 星级

	Tag1        string
	Tag2        string
	UnlockLevel string //解锁关卡
	HiddenTag   string // 暗标签
	SaleTime    string //售卖时间
	AdminName   string //上传名称
	UploadTime  string //上传时间
}

type Version struct {
	GameName      string `xorm:"game"` //游戏名字
	CoreVersion   string //大版本
	DownLoadUrl   string `xorm:"download_url"`  // 下载地址
	DownLoadUrl2  string `xorm:"download_url2"` //
	ScriptVersion string //小版本
	VersionFile   string //版本地址 map[string]string
	Tishen        int    //提审
}

type Maintenance struct {
	ID        int64 `xorm:"id"`
	Content   string
	StartTime int64
	EndTime   int64
	ShowTime  int64
	URL       string `xorm:"url"`
}

type ConfigInfo struct {
	ID    int64 `xorm:"id"`
	Key   string
	Value string
}

type AdminGM struct {
	ID            int64  `xorm:"id"`
	Username      string `xorm:"username"`
	Password      string `xorm:"password"`
	Nickname      string `xorm:"nickname"`
	LastLoginTime string `xorm:"login_time"`
	Role          int64  `xorm:"role"`
}

type ApiGM struct {
	ID        string `xorm:"id"`
	Name      string `xorm:"name"`
	MenuLayer string `xorm:"menu_layer"`
	MenuUrl   string `xorm:"menu_url"`
	MenuType  string `xorm:"menu_type"`
	ParentID  string `xorm:"parent_id"`
	Order     int    `xorm:"order"`
}

type MenuInfo struct {
	ID       int64      `json:"id"`
	Title    string     `json:"title"`
	Icon     string     `json:"icon"`
	ParentID int64      `json:"parent_id"`
	Spread   bool       `json:"spread"`
	Href     string     `json:"href"`
	Children []MenuInfo `json:"children"`
}

type AdminMenuPriv struct {
	ID      int64 `xorm:"id"`
	AdminID int64 `xorm:"admin_id"`
	MenuID  int64 `xorm:"menu_id"`
}

type RoleApiPriv struct {
	ID         int64  `xorm:"id"`
	RoleID     int64  `xorm:"role_id"`
	MenuIDList string `xorm:"menu_id"`
}

type AdminRolePriv struct {
	ID      int64 `xorm:"id"`
	AdminID int64 `xorm:"admin_id"`
	RoleID  int64 `xorm:"role_id"`
}

type AdminRole struct {
	ID       int64  `xorm:"id"`
	RoleName string `xorm:"role_name"`
	Desc     string `xorm:"desc"`
	Disable  bool   `xorm:"disable"`
}

type DesignPaperGM struct {
	PaperID        int64 `xorm:"id"`
	AuthorUsername string
	AuthorNickname string
	AuthorHead     string

	ClothesType string //衣服模特
	PartType    string //衣服类别
	Cname       string //图纸名称
	Desc        string //描述
	Extra       string `xorm:"paper_extra"` //图纸文件
	File        string `xorm:"paper_file"`
	Star        int    //图纸星级
	Status      int    //
	//普通标签
	Tag1 string
	Tag2 string

	//风格 内容是PaperStyle的json字符串
	Style string `xorm:"paper_style"`

	//暗标签
	STag       string `xorm:"stag"`
	UploadTime int64
	Reason     string
}

type SalePaperGM struct {
	ID            int64  `xorm:"id"`
	OwnerUsername string `xorm:"owner_username"`
	PaperID       string `xorm:"paper_id"`
	Cname         string `xorm:"cname"`
	ClothesType   string //模特名称
	PartType      string //服装类别
	Price         int64
	PriceType     int
	Star          int //图纸星级
	Tag1          string
	Tag2          string
	Style         string `xorm:"paper_style"`
	STag          string `xorm:"stag"`
	Extra         string `xorm:"paper_extra"` //图纸文件
}

type OffPaperGM struct {
	ID          int64  `xorm:"id"`
	Cname       string `xorm:"cname"`
	ClothesType string //模特名称
	PartType    string //服装类别
	Price       int64
	PriceType   int
	Star        int //图纸星级
	Tag1        string
	Tag2        string
	Style       string `xorm:"paper_style"`
	STag        string `xorm:"hidden_tag"`
	Extra       string `xorm:"paper_extra"` //图纸文件
}

type DesignPaperGMTrans struct {
	PaperID        int64
	AuthorUsername string
	ClothesType    string //模特名称
	PartType       string //服装类别
	Desc           string //描述
	ModelPic       string //模特图片

	Status       string  //状态
	ClothesLayer int     //图层
	IconImg      string  //缩略图
	BottomLayer  string  //上层
	MainImg      string  //下层
	CollorImg    string  //领子
	ShadowImg    string  // 阴影
	IconIX       float64 //游戏图标偏移
	IconIY       float64 //游戏图标偏移

	UploadTime    string //上传时间
	Front         string
	Back          string
	RecentlyPaper []PaperSimple
}

type PaperSimple struct {
	PaperID int64
	IconImg string //图标
}

type StatusGM struct {
	Username       string
	Uid            int64
	Nickname       string
	StreetName     string          //街道名称
	Head           string          //头像
	Sex            int             //性别
	Level          int             //等级
	Exp            int             //经验
	Money          int64           //金币
	Tili           int             //体力
	MaxTili        int             //最大体力
	TiliRTime      string          //体力
	OffClothesList []ClothesInfoGM //官方服装列表
	ClothesList    []ClothesInfoGM //设计师服装列表
	HomeShow       HomeShowDetail
	SuitMap        string               //套装内容
	RecordList     []PixieRecordGM      // 关卡
	NpcFriendship  []FriendshipDetailGM //Npc亲密度
	NpcBuff        string
	StreetOperate  string //街道运营信息
	DayExtraInfo   string //日常信息
	ExtraInfo1     string //
	BanStartTime   int64
	BanEndTime     int64
	BanReason      string
}

type PixieRecordGM struct {
	ID          string  //关卡ID
	Score       float64 //分数
	BestScore   float64 //最高分
	Rank        string  //排名
	Clothes     string  //衣服
	BestClothes string  //最佳搭配
	HistoryRank string  //历史排名
}

type ClothesInfoGM struct {
	ClothesID string
	Count     int
	Time      string // 最近一次购买的时间
}

type FriendshipDetailGM struct {
	NpcID string
	Level int //等级
	Exp   int //经验
}
