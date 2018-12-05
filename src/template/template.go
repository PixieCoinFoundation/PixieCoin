package template

import (
	"appcfg"
	"constants"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var Templates *template.Template
var templateDelims = []string{"{{%", "%}}"}

func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GV && appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	var basePath string

	if appcfg.GetServerType() == constants.SERVER_TYPE_GM {
		basePath = "gm_pages/"
		if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
			basePath = "gm_pages_en/"
		}
	} else if appcfg.GetServerType() == constants.SERVER_TYPE_GV {
		basePath = "gv_pages/"
		if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
			basePath = "gv_pages_en/"
		}
	}

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't process folders themselves
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".html") {
			return nil
		}
		// Info("parsing template:", path)
		templateName := path[len(basePath) : len(path)-5]
		if Templates == nil {
			Templates = template.New(templateName)
			Templates.Delims(templateDelims[0], templateDelims[1])
			_, err = Templates.ParseFiles(path)
		} else {
			_, err = Templates.New(templateName).ParseFiles(path)
		}
		// Info("Processed template %s\n", templateName)
		return err
	})

	if err != nil {
		panic(err)
	}

	if Templates == nil {
		panic("no template found:" + basePath)
	}
}
