package api_specification

type Subject struct {
	ID        int64
	Name      string
	SuitOwner string
	HotSpots  []HSpot
	//推荐单品(1.推荐单品跑马灯 2.推荐单品只显示(icon) 3.推荐单品分页码,显示)

	// //推荐专题
	// TopicList []SubjectTopic
	OpenTime  int64
	BannerImg string
}

type HSpot struct {
	X      int64  `json:"x"`
	Y      int64  `json:"y"`
	Width  int64  `json:"w"`
	Height int64  `json:"h"`
	HType  string `json:"t"`
	HLink  string `json:"hl"`
}

type SubjectItems struct {
	ID         int64
	SubjectID  int64 //期刊ID
	ClothesID  string
	PaperFile  string
	PaperExtra string
}

type SubjectItemsTop struct {
	ID            int64
	SubjectID     int64
	ClothesIDList string //推荐服装ID
	BannerImg     string //图片
	HtmlLink      string //h5链接
}

type SubjectPool struct {
	ID            int64  `xorm:"id"`
	ClothesID     string `xorm:"clothes_id"` //推荐服装ID
	OwnerUsername string
	BannerImg     string `xorm:"banner_file_id"` //图片
	AddTime       int64  `xorm:"add_time"`
	Status        int    `xorm:"status"` //衣服状态
}

type SubjectItemsMedium struct {
	ID        int64
	PaperList []*ItemsPaper //服装信息
}

type SubjectItemsDown struct {
	ID        int64
	PaperList []*ItemsPaper //服装信息
}

type TopPapers struct {
	SubjectID int64
	BannerImg string
	HtmlLink  string
	Papers    []RecommendClothes
}

type ItemsPaper struct {
	BannerImg string
	HtmlLink  string
	AddTime   int64
	RecommendClothes
}

type SubjectTopic struct {
	ID        int64
	SubjectID int64
	PaperList []RecommendClothes
	BannerImg string
	HtmlLink  string
	AddTime   int64
}
