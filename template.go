package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ThemeDir    = "themes"
	TplVarLeft  = "{{"
	TplVarRight = "}}"
	T           *template.Template
)

func init() {
	InitTemplate()
}

func RebuildTemplates(theme string) {
	InitTemplate()
	BuildTemplates(theme)
}

func InitTemplate() {
	T = template.New("__top__")
	funcMap := make(template.FuncMap)
	funcMap["str2html"] = func(raw string) template.HTML {
		return template.HTML(raw)
	}

	T.Funcs(funcMap)
}

func BuildTemplates(theme string) {
	root := filepath.Join(workPath, ThemeDir, theme)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			source, err := ioutil.ReadFile(path)
			checkError(err)
			_, err = T.New(info.Name()).Parse(string(source))
			checkError(err)
		}
		return nil
	})
}

func RenderTemplate(w io.Writer, tpl string, data interface{}) error {
	return T.ExecuteTemplate(w, tpl, data)
}
