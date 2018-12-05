package language

import (
	"appcfg"
	"constants"
)

func L(key string) string {
	if v, ok := languageMap[key]; ok {
		if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
			return v.Korean
		} else {
			return v.Chinese
		}
	} else {
		return "unknown"
	}
}
