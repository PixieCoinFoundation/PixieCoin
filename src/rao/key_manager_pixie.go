package rao

import (
	"fmt"
)

func GetPaperKey(paperID int64) string {
	return fmt.Sprintf("p_%d", paperID)
}
