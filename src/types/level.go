package types

type GFLevelS struct {
	LevelID string  `json:"id"`
	No      int     `json:"no"`
	S       string  `json:"s"`
	A       string  `json:"a"`
	B       string  `json:"b"`
	F       string  `json:"f"`
	Keyword string  `json:"keyword"`
	Type    int     `json:"type"`
	Char1   int     `json:"char1"`
	Char2   int     `json:"char2"`
	SLine   int     `json:"s_line"`
	ALine   int     `json:"a_line"`
	BLine   int     `json:"b_line"`
	Adjust  int     `json:"adjust"`
	Warm    float64 `json:"warm"`
	Formal  float64 `json:"formal"`
	Tight   float64 `json:"tight"`
	Bright  float64 `json:"bright"`
	Dark    float64 `json:"dark"`
	Cute    float64 `json:"cute"`
	Man     float64 `json:"man"`
	Tough   float64 `json:"tough"`
	Noble   float64 `json:"noble"`
	Strange float64 `json:"strange"`
	Sexy    float64 `json:"sexy"`
	Sport   float64 `json:"sport"`
	Level   int     `json:"level"`
}

type GFLevel struct {
	LevelID     string `json:"id"`
	No          int    `json:"no"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Story       string `json:"story"`
	SSuite      string `json:"s_suite"`
	Planet      int    `json:"planet"`
	Comments    string `json:"comments"`
	Type        int    `json:"type"`
	Content     string `json:"content1"`
	GoldDropS   int    `json:"gold_drop_s"`
	GoldDropA   int    `json:"gold_drop_a"`
	GoldDropB   int    `json:"gold_drop_b"`
	GoldDropF   int    `json:"gold_drop_f"`
	ClothesDrop int    `json:"clothes_drop"`
	SLine       int    `json:"s_line"`
	ALine       int    `json:"a_line"`
	BLine       int    `json:"b_line"`
	Adjust      int    `json:"adjust"`
	Branch      int    `json:"branch"`
}

type GFLevelP struct {
	LevelID string `json:"id"`
	No      int    `json:"no"`
	Level   int    `json:"level"`
	// Name         string `json:"name"`
	// Title        string `json:"title"`
	// Story        string `json:"story"`
	// SSuite string `json:"s_suite"`
	Planet int `json:"planet"`
	// Comments     string `json:"comments"`
	Type int `json:"type"`
	// Content      string `json:"content1"`
	GoldDropS    int    `json:"gold_drop_s"`
	GoldDropA    int    `json:"gold_drop_a"`
	GoldDropB    int    `json:"gold_drop_b"`
	GoldDropF    int    `json:"gold_drop_f"`
	ClothesDrop1 string `json:"clothes_drop1"`
	ClothesDrop2 string `json:"clothes_drop2"`
	// SLine        int    `json:"s_line"`
	// ALine        int    `json:"a_line"`
	// BLine        int    `json:"b_line"`
	// Adjust int `json:"adjust"`
	Branch    int `json:"branch"` //1 normal 2 planet end 3 ending
	TheaterID int
}
