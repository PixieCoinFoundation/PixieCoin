package types

import (
	. "pixie_contract/api_specification"
)

type UploadPaperLogExtra struct {
	*UploadPaperReq
	PaperID int64
}

type DeletePaperLogExtra struct {
	*DeletePaperReq
}
