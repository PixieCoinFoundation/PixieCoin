package model

type RecommendSubjectGM struct {
	ID           int64  `xorm:"id"`
	Name         string `xorm:"name"`
	SuitDetail   string `xorm:"suit_detail"`
	SuitOwnerID  string `xorm:"suit_owner"`
	HotSpot      string `xorm:"hot_spot"`
	OpenTime     int64  `xorm:"open_time"`
	BannerFileID string `xorm:"banner_file_id"`
}

//推荐单品跑马灯
type RecommendItemGM struct {
	ID            int64  `xorm:"id"`
	SubjectID     int64  `xorm:"subject_id"`
	ClothesIDList string `xorm:"clothes_id"`
	BannerFileID  string `xorm:"banner_file_id"`
	HtmlZipFileID string `xorm:"html_zip_file_id"`
}

//每周专题
type RecommendItemsGM struct {
	ID            int64  `xorm:"id"`
	SubjectID     int64  `xorm:"subject_id"`
	ItemList      string `xorm:"item_list"` //专题服装列表
	BannerFileID  string `xorm:"banner_file_id"`
	HtmlZipFileID string `xorm:"html_zip_file_id"`
	AddTime       int64  `xorm:"add_time"`
}

type TransRecommendItemGM struct {
	ID            int64
	SubjectName   string
	ClothesIDList string
	BannerImg     string
	Link          string
}

//推荐单品池
type RecommendPoolGM struct {
	ID            int64  `xorm:"id"`
	ClothesID     string `xorm:"clothes_id"`
	OwnerUsername string `xorm:"owner_username"`
	AddTime       int64  `xorm:"add_time"`
	Status        int    `xorm:"status"`
}
