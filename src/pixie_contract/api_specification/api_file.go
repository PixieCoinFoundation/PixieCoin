package api_specification

//上传文件接口
const HTTP_PLAYER_UPLOAD_FILE string = "api/uploadFiles"

type FileUploadReq struct {
	Content [][]byte

	//1:上传作品 包含icon main bottom collar shadow
	//2:上传作品审核期的举报
	//3:上传套装图片
	Type int
}

type FileUploadResp struct {
	Filenames []string
}
