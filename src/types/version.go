package types

type GFVersion struct {
	Game               string `bson:"game" json:"game"`
	CoreVersion        string `bson:"coreVersion" json:"coreVersion"`
	CoreDownloadURL    string `bson:"downloadUrl" json:"downloadUrl"`
	CoreDownloadURL2   string
	ScriptVersion      string `bson:"username"  json:"scriptVersion"`
	VersionFile        string `bson:"versionFile"  json:"versionFile"`
	FileServerAddress  string `bson:"fileServerAddr"  json:"fileServerAddr"`
	FileServerAddress2 string `bson:"fileServerAddr2"  json:"fileServerAddr2"`
	FileServerAddress3 string `bson:"fileServerAddr3"  json:"fileServerAddr3"`
	Tishen             bool   `bson:"tishen"  json:"tishen"` //true代表大版本不同时连接提审服 false表示进行强更
}
