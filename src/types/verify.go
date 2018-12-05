package types

type Verify struct {
	Vid              int64
	PlayerId         string
	PlayerName       string
	PlayerLv         string
	ClothesPath      string
	CloSet           CustomFile
	ClothesIDInGame  string
	IconPath         string
	ClothesSex       string
	ClothesModelNo   int
	ClothesType      string
	ClothesName      string
	ClothesDes       string
	ClothesPrice     int
	ClothesPriceType string
	SubmitTime       string
	Status           string
	VerifyUser       string
	VerifyTime       string
	Info             string

	CopyEvidenceFileName string
	DownTime             string

	Hearts       int
	Score1       int
	Score2       int
	Score3       int
	Score4       int
	Score5       int
	Score6       int
	Score7       int
	Score8       int
	Score9       int
	Score10      int
	Score11      int
	Sale         bool
	TotalCount   int
	Platform     string
	OutputIcon   string
	OutputMain   string
	OutputCollar string
	OutputBottom string
	OutputShadow string
	Scale        float64
	IX           float64
	IY           float64

	BottomUploaded bool
	CollarUploaded bool
	ShadowUploaded bool

	ReportID       int64
	ReportUsername string
	Head           string
	Reason         string
	Contact        string

	ETHAccount string

	RecentCustoms []RecentCustom

	Remark     string
	ZValue     int
	ZValuec    int
	ZValues    int
	ZValueName string
	Tag1       string
	Tag2       string
	Tag1i      int
	Tag2i      int
	Typei      int

	GoldProfit    int
	DiamondProfit int
	Exp           int
	PlayerConfig  string
}
