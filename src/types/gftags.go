package types

type GFTags struct {
	UnlockIndex int `json:"unlockIndex"` // 已解锁到的位置
	// BoyClothesMap  map[string](*TagClothes) `json:"boyClothesMap"`  // key-value: 服装类型-对应的5个列表
	// GirlClothesMap map[string](*TagClothes) `json:"girlClothesMap"` // key-value: 服装类型-对应的5个列表

	TagsMap map[string]map[string](*TagClothes) `json:"tagsMap"`
}

type TagClothes struct {
	Names       [5]string   `json:"names"`   // 5个标签的名字
	Covers      [5]string   `json:"covers"`  // 5个封面
	ClothesList [5][]string `json:"clothes"` // 5个list
}
