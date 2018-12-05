package pixie_contract

var PIXIE_ERR_UPLOAD_TYPE_WRONG = PixieRespInfo{
	RespCodeStr: "1001",
	RespDesc:    "upload file type wrong",
	RespDescCN:  "上传类型错误",
}

var PIXIE_ERR_UPLOAD_FILE_FAIL = PixieRespInfo{
	RespCodeStr: "1002",
	RespDesc:    "upload file fail",
	RespDescCN:  "上传失败",
}

var PIXIE_ERR_UPLOAD_LIMIT = PixieRespInfo{
	RespCodeStr: "1003",
	RespDesc:    "upload file limit",
	RespDescCN:  "上传次数限制",
}
